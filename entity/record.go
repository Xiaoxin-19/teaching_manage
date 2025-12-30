package entity

import "time"

type Record struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Student      Student
	Teacher      Teacher
	TeachingDate time.Time
	StartTime    string
	EndTime      string
	Active       bool
	Remark       string
}
