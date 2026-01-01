package entity

import "time"

type Order struct {
	Student
	Id        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Hours     int       `json:"hours"`
	Comment   string    `json:"comment"`
	Active    bool      `json:"active"`
}
