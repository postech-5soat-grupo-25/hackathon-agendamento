package storage

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/lib/pq"

	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/config"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
)

var (
	db  *sql.DB
	err error
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) GetAgendamento(id string) (*models.Appointment, error) {
	return nil, nil
}
func (p *Postgres) GetAgendamentos() ([]*models.Appointment, error) {
	return nil, nil
}

func (p *Postgres) CreateAgendamento(agendamento *models.Appointment) error {
	return nil
}

func (p *Postgres) ExcluirAgendamento(id string) error {
	return nil
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
