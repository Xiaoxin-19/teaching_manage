package entity

import "time"

type Student struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	Hours       int       `json:"hours"`
	Phone       string    `json:"phone"`
	TeacherID   uint      `json:"teacher_id"`
	Remark      string    `json:"remark"`
	TeacherName string    `json:"teacher_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Teacher     Teacher   `json:"teacher"`
	DeletedAt   time.Time `json:"deleted_at"`
}
