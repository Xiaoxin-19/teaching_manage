package entity

import (
	"teaching_manage/backend/pkg"
	"time"
)

type Teacher struct {
	ID            uint       `json:"id"`
	TeacherNumber string     `json:"teacher_number"`
	Name          string     `json:"name"`
	Gender        pkg.Gender `json:"gender"`
	Phone         string     `json:"phone"`
	Remark        string     `json:"remark"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     time.Time  `json:"deleted_at"`
}
