package storage

import "github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"

type Storage interface {
	GetAgendamento(id string) (*models.Appointment, error)
	GetAgendamentosCliente(id_client int) ([]*models.Appointment, error)
	CreateAgendamento(agendamento *models.Appointment) (*models.Appointment, error)
	ExcluirAgendamento(id int) error
	CreateOrEditWorkingHours(workhours *models.WorkingHours) (*models.WorkingHours, error)
	GetWorkingHours(id_doctor int) (*models.WorkingHours, error)
}
