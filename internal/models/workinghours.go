package models

import "time"

type WorkingHours struct {
    ID         int       `json:"id"`
    DoctorID   int       `json:"doctor_id"`    
    StartTime  time.Time `json:"start_time"`   
    EndTime    time.Time `json:"end_time"`     
    DaysOfWeek []bool     `json:"days_of_week"`
}