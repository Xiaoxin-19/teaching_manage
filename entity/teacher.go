package entity

import "teaching_manage/pkg"

type Teacher struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Gender    pkg.Gender `json:"gender"`
	Phone     string     `json:"phone"`
	Remark    string     `json:"remark"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
	DeletedAt int64      `json:"deleted_at"`
}
