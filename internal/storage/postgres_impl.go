package storage

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
	"fmt"
	"github.com/lib/pq"

	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/config"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
)

var (
	db  *sql.DB
	err error
)

const (
    getWorkingHoursQuery = `SELECT id, doctor_id, start_time, end_time, days_of_week 
                            FROM working_hours 
                            WHERE doctor_id = $1;`

    updateWorkingHoursQuery = `UPDATE working_hours 
                               SET start_time = $1, end_time = $2, days_of_week = $3 
                               WHERE doctor_id = $4;`

    insertWorkingHoursQuery = `INSERT INTO working_hours (doctor_id, start_time, end_time, days_of_week) 
                               VALUES ($1, $2, $3, $4);`

	insertAppointmentQuery = `INSERT INTO appointments (doctor_id, client_id, appointment_time, description)
							   VALUES ($1, $2, $3, $4) RETURNING id;`

	checkAvailabilityQuery = `SELECT COUNT(*) 
							   FROM appointments 
							   WHERE (doctor_id = $1 OR client_id = $2) 
							   AND appointment_time = $3;`
	cancelAppointmentQuery = `DELETE FROM appointments 
							   WHERE id = $1;`

	getClientAppointmentsQuery = `SELECT id, doctor_id, client_id, appointment_time, description
							   FROM appointments
							   WHERE client_id = $1;`					   
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) GetAgendamento(id string) (*models.Appointment, error) {
	return nil, nil
}
func (p *Postgres) GetAgendamentosCliente(id_client int) ([]*models.Appointment, error) {
    rows, err := p.db.Query(getClientAppointmentsQuery, id_client)
    if err != nil {
        return nil, fmt.Errorf("failed to query client appointments: %w", err)
    }
    defer rows.Close()

    var appointments []*models.Appointment

    for rows.Next() {
        var appointment models.Appointment

        err := rows.Scan(&appointment.ID, &appointment.DoctorID, &appointment.ClientID, &appointment.AppointmentTime, &appointment.Description)
        if err != nil {
            return nil, fmt.Errorf("failed to scan appointment row: %w", err)
        }

        appointments = append(appointments, &appointment)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
    }

    return appointments, nil
}

func (p *Postgres) CreateAgendamento(agendamento *models.Appointment) (*models.Appointment, error) {
    var count int
    err = p.db.QueryRow(checkAvailabilityQuery, agendamento.DoctorID, agendamento.ClientID, agendamento.AppointmentTime).Scan(&count)
    if err != nil {
        fmt.Println("failed to insert appointment: %w", err)
        return nil, fmt.Errorf("failed to check availability: %w", err)
    }

    if count > 0 {
        return nil, fmt.Errorf("doctor or client already has an appointment at this time")
    }
    var appointmentID int
    err = p.db.QueryRow(insertAppointmentQuery, agendamento.DoctorID, agendamento.ClientID, agendamento.AppointmentTime, agendamento.Description).Scan(&appointmentID)
    if err != nil {
        fmt.Println("failed to insert appointment: %w", err)
        return nil, fmt.Errorf("failed to insert appointment: %w", err)
    }

    agendamento.ID = appointmentID

    return agendamento, nil
}

func (p *Postgres) ExcluirAgendamento(id int) error {
	result, err := p.db.Exec(cancelAppointmentQuery, id)
    if err != nil {
        return fmt.Errorf("failed to cancel appointment with id %s: %w", id, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check affected rows: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no appointment found with id: %s", id)
    }

    return nil
}

func (p *Postgres) CreateOrEditWorkingHours(workhours *models.WorkingHours) (*models.WorkingHours, error) {
    existingWH, err := p.GetWorkingHours(workhours.DoctorID)
    if err != nil && err.Error() != fmt.Sprintf("no working hours found for doctor_id: %d", workhours.DoctorID) {
        return nil, fmt.Errorf("error fetching existing working hours: %w", err)
    }

    if existingWH != nil { //update entry
        _, err := p.db.Exec(updateWorkingHoursQuery, 
            workhours.StartTime, 
            workhours.EndTime, 
            pq.BoolArray(workhours.DaysOfWeek), 
            workhours.DoctorID)
        
        if err != nil {
            return nil, fmt.Errorf("failed to update working hours: %w", err)
        }
    } else { // new entry
        _, err := p.db.Exec(insertWorkingHoursQuery, 
            workhours.DoctorID, 
            workhours.StartTime, 
            workhours.EndTime, 
            pq.BoolArray(workhours.DaysOfWeek))
        
        if err != nil {
            return nil, fmt.Errorf("failed to insert new working hours: %w", err)
        }
    }

    return workhours, nil
}

func (p *Postgres)GetWorkingHours(id_doctor int) (*models.WorkingHours, error) {
    row := p.db.QueryRow(getWorkingHoursQuery, id_doctor)

    var wh models.WorkingHours
	var daysOfWeek pq.BoolArray
    err := row.Scan(&wh.ID, &wh.DoctorID, &wh.StartTime, &wh.EndTime, &daysOfWeek)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("failed to GetWorkingHours: %w", err)
            return nil, fmt.Errorf("no working hours found for doctor_id: %d", id_doctor)
        }
        fmt.Println("failed to GetWorkingHours: %w", err)
        return nil, fmt.Errorf("failed to query working hours: %w", err)
    }
	wh.DaysOfWeek = []bool(daysOfWeek)
    return &wh, nil
}

func NewStorage(ctx context.Context) (Storage, error) {
	connStr := config.GetPostgresConnString()
	// retry to connect to the database
	for {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			break
		}

		slog.Log(ctx, slog.LevelError, "Failed to connect to database", "error", err)
		slog.Log(ctx, slog.LevelError, "Retrying in 2 seconds...")

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	err := db.Ping()
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		return nil, err
	}

	return &Postgres{db: db}, nil
}
