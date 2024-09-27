package models

import "time"

type Appointment struct {
	ID             int        `json:"id"`
    DoctorID       int        `json:"doctor_id"`
    ClientID       int        `json:"client_id"`
	AppointmentTime time.Time `json:"appointment_time"`
	Description string        `json:"description"`
}
