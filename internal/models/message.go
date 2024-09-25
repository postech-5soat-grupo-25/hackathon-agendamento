package models

import (
	"encoding/json"
)

const (
	GetDoctorWorkingHoursMessageType = "getdoctorworkinghours"
	appointment = "appointment"
	workingHours = "workinghours"
)

type Message struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type GetDoctorWorkingHoursMessage struct {
	DoctorID   int       `json:"DoctorID"`
}
