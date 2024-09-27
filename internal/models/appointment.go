package models

import "time"

type Appointment struct {
	ID             int        `json:"ID"`
    DoctorID       int        `json:"DoctorID"`
    ClientID       int        `json:"ClientID"`
	AppointmentTime time.Time `json:"AppointmentTime"`
	Description string        `json:"Description"`
}
