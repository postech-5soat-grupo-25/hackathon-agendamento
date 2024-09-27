package controllers

import (
    "fmt"
    "time"
    "encoding/json"
    "github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
    "github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/storage"
)

type AppointmentsController struct {
    storage storage.Storage
}

func (s *AppointmentsController) CreateOrEditWorkingHours(hours *models.WorkingHours) models.Response {
    if !validateWorkingHours(*hours) {
        return errorResponse(400, "Invalid working hours")
    }

    wkhrs, err := s.storage.CreateOrEditWorkingHours(hours)
    if err != nil {
        return errorResponse(500, "Internal Server Error")
    }

    wkhrsJson, err := json.Marshal(wkhrs)
    if err != nil {
        return errorResponse(500, "Internal Server Error")
    }

    return models.Response{
        StatusCode: 200,    
        Body:       wkhrsJson,
    }
}

func (s *AppointmentsController) GetWorkingHours(doctorID int) models.Response {
    wkhrs, err := s.storage.GetWorkingHours(doctorID) 

    if err != nil {
        return errorResponse(500, "Undefined error")
    }

    wkhrsJson, err := json.Marshal(wkhrs)
    if err != nil {
        return errorResponse(500, "Internal Server Error")
    }

    return models.Response{
        StatusCode: 200,    
        Body:       wkhrsJson,
    }
}

func (s *AppointmentsController) ScheduleAppointment(appointment *models.Appointment) models.Response {
    // now := time.Now()
    // oneHourLater := now.Add(1 * time.Hour)
    // if (!appointment.AppointmentTime.After(oneHourLater)){
    //     return errorResponse(400, "Consultas devem ser marcadas com 1 hora de antecendencia")
    // }

    wkhrs, err := s.storage.GetWorkingHours(appointment.DoctorID)
    if err != nil {
        fmt.Println("failed to GetWorkingHours: %w", err)
        return errorResponse(500, "Undefined error")
    }

    if !isWithinWorkingHours(wkhrs, appointment.AppointmentTime) {
        return errorResponse(400, "Fora do horario do médico")
    }

    agendamento, err := s.storage.CreateAgendamento(appointment)

    if err != nil {
        return errorResponse(500, "Internal Server Error")
    }

    agendamentoJson, err := json.Marshal(agendamento)
    if err != nil {
        return errorResponse(500, "Internal Server Error")
    }

    return models.Response{
        StatusCode: 200,    
        Body:       agendamentoJson,
    }
}

func (s *AppointmentsController) CancelScheduledAppointment(id int) models.Response {
    err := s.storage.ExcluirAgendamento(id)
    if err != nil {
        return errorResponse(500, "Undefined error")
    }

    return models.Response{
        StatusCode: 200,    
        Body:       nil,
    }
}

func (s *AppointmentsController) GetClientAppointments(id_client int) models.Response {
    var appointments []*models.Appointment
    appointments, err := s.storage.GetAgendamentosCliente(id_client)

    if err != nil {
        return errorResponse(500, "Undefined error")
    }

    appointmentsJson, err := json.Marshal(appointments)
    if err != nil {
        return errorResponse(500, "Internal Server Error")
    }

    return models.Response{
        StatusCode: 200,    
        Body:       appointmentsJson,
    }
}

func isWithinWorkingHours(wkhrs *models.WorkingHours, appointmentTime time.Time) bool {
    weekday := int(appointmentTime.Weekday())  // 0 = Sunday, 1 = Monday, etc.
    appointmentHour := appointmentTime.Hour()
    appointmentMinute := appointmentTime.Minute()

    if !wkhrs.DaysOfWeek[weekday] {
        return false
    }

    startTime := wkhrs.StartTime.Hour()*60 + wkhrs.StartTime.Minute()
    endTime := wkhrs.EndTime.Hour()*60 + wkhrs.EndTime.Minute()
    appointmentTimeInMinutes := appointmentHour*60 + appointmentMinute

    return appointmentTimeInMinutes >= startTime && appointmentTimeInMinutes < endTime
}

func validateWorkingHours(hours models.WorkingHours) bool {
    if len(hours.DaysOfWeek) != 7 {
        fmt.Println("aqui")
        return false
    }

    if hours.EndTime.Before(hours.StartTime) {
        fmt.Println("eu")
        return false
    }

    return true
}

func errorResponse(statusCode int, message string) models.Response {
    errorMessage := map[string]string{
        "error": message,
    }
    body, _ := json.Marshal(errorMessage)

    response := models.Response{
        StatusCode: statusCode,    
        Body:       body,
    }
    return response
}

func NewController(storage storage.Storage) AppointmentsInterface {
    return &AppointmentsController{storage: storage}
}
