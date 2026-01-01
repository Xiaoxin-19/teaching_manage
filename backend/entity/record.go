package entity

import "time"

type Record struct {
	ID           uint `json:"id"`
	Student      Student
	Teacher      Teacher
	TeachingDate time.Time `json:"teaching_date"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Active       bool      `json:"active"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
