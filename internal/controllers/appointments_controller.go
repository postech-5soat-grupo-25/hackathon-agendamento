package controllers

import "github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"

type AppointmentsInterface interface {
	GetWorkingHours(id_doctor int) (*models.WorkingHours, error)
	CreateOrEditWorkingHours(workhours *models.WorkingHours) error

	// ScheduleAppointment(agendamento *models.Appointment)
	// CancelScheduledAppointment(id string) error
	// GetClientAppointments(id_client string) ([]*models.Appointment, error)
}
