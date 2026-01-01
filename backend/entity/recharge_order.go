package entity

import "time"

type RechargeOrder struct {
	ID        uint      `json:"id"`
	StudentID uint      `json:"student_id"`
	Amount    float64   `json:"amount"`
	Hours     int       `json:"hours"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
