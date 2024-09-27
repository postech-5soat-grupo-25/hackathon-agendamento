package models

import "time"

type WorkingHours struct {
    ID         int       `json:"ID"`
    DoctorID   int       `json:"DoctorID"`    
    StartTime  time.Time `json:"StartTime"`   
    EndTime    time.Time `json:"EndTime"`     
    DaysOfWeek []bool     `json:"DaysOfWeek"`
}