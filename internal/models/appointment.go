package models

type Appointment struct {
	Id          int
	PatientId   int
	DoctorId    int
	Date        string
	Time        string
	Description string
}
