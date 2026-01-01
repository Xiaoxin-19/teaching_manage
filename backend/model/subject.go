package model

import "gorm.io/gorm"

// Subject 科目表
type Subject struct {
	gorm.Model
	Name string `gorm:"unique;size:50;not null;comment:科目名称"` // 如：钢琴、声乐
}

func (Subject) TableName() string {
	return "subjects"
}
