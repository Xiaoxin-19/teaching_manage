package entity

import "time"

type Student struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	StudentNumber string    `json:"student_number"`
	Gender        string    `json:"gender"`
	Phone         string    `json:"phone"`
	Status        int       `json:"status"`
	Remark        string    `json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}
