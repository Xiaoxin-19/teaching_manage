package entity

import "time"

type RechargeOrder struct {
	ID            uint           `json:"id"`
	OrderNumber   string         `json:"order_number"`
	StudentCourse StudentSubject `json:"student_course"`
	Amount        float64        `json:"amount"`
	Hours         int            `json:"hours"`
	Status        string         `json:"status"`
	Remark        string         `json:"remark"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     time.Time      `json:"deleted_at"`
}
