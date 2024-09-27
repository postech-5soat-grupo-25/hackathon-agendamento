package controllers

import "github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"

type AppointmentsInterface interface {
	GetWorkingHours(id_doctor int) models.Response
	CreateOrEditWorkingHours(workhours *models.WorkingHours) models.Response
	ScheduleAppointment(agendamento *models.Appointment) models.Response
	CancelScheduledAppointment(id int) models.Response
	GetClientAppointments(id_client int) models.Response
}
