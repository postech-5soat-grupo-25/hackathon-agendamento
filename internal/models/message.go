package models

import (
	"encoding/json"
)

const (
	GetDoctorWorkingHoursMessageType = "getdoctorworkinghours"
	GetClientAppointmentsMessage = "getappointment"
	AppointmentMessage = "appointment"
	CancelScheduledAppointmentMessage = "cancelappointment"
	WorkingHoursMessage = "workinghours"
)

type Message struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type Response struct {
	StatusCode int    `json:"statuscode"`
	Body json.RawMessage `json:"body"`
}

type GetDoctorWorkingHoursMessage struct {
	DoctorID   int       `json:"DoctorID"`
}

type GenericIDMessage struct {
	ID   int       `json:"ID"`
}
