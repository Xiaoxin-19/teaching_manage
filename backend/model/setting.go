package model

import "gorm.io/gorm"

type Setting struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex;not null"`
	Value string `gorm:"not null"`
}

func (s *Setting) TableName() string {
	return "settings"
}
