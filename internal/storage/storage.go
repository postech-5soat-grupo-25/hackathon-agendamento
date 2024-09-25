package storage

import "github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"

type Storage interface {
	GetAgendamento(id string) (*models.Appointment, error)
	GetAgendamentos() ([]*models.Appointment, error)
	CreateAgendamento(agendamento *models.Appointment) error
	ExcluirAgendamento(id string) error
	CreateOrEditWorkingHours(workhours *models.WorkingHours) error
	GetWorkingHours(id_doctor int) (*models.WorkingHours, error)
}
