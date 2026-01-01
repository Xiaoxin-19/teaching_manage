package model

import "gorm.io/gorm"

// Teacher 教师表
type Teacher struct {
	gorm.Model
	Name   string `gorm:"column:name;not null;comment:教师姓名;index;unique" `
	Gender string `gorm:"column:gender;comment:教师性别" `
	Phone  string `gorm:"column:phone;comment:电话号码" `
	Remark string `gorm:"column:remark;comment:备注" `
}

func (Teacher) TableName() string {
	return "teachers"
}
