package controllers

import (
    "fmt"
    "github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
    "github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/storage"
)

type AppointmentsController struct {
    storage storage.Storage
}

func (s *AppointmentsController) CreateOrEditWorkingHours(hours *models.WorkingHours) error {
    if !validateWorkingHours(*hours) {
        return fmt.Errorf("invalid working hours")
    }

    err := s.storage.CreateOrEditWorkingHours(hours)
    if err != nil {
        return err
    }

    return nil
}

func (s *AppointmentsController) GetWorkingHours(doctorID int) (*models.WorkingHours, error) {
    wkhrs, err := s.storage.GetWorkingHours(doctorID) 

    if err != nil {
        return nil, err
    }

    return wkhrs, nil
}

func validateWorkingHours(hours models.WorkingHours) bool {
    if len(hours.DaysOfWeek) != 7 {
        return false
    }

    if hours.EndTime.Before(hours.StartTime) {
        return false
    }

    return true
}

func NewController(storage storage.Storage) AppointmentsInterface {
    return &AppointmentsController{storage: storage}
}
