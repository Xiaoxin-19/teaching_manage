package model

import "gorm.io/gorm"

// RechargeOrder 充值记录表，记录每次给学生课时充值或扣减的明细
type RechargeOrder struct {
	gorm.Model
	StudentCourseID uint           `gorm:"index;not null;comment:关联的学籍ID" `
	StudentCourse   StudentSubject `gorm:"foreignKey:StudentCourseID" `
	Hours           int            `gorm:"comment:变动课时数(正数充值,负数扣减)" `
	Type            string         `gorm:"size:20;comment:类型(充值/赠送/退费/消课)" `
	Operator        string         `gorm:"size:50;comment:操作人" `
	Remark          string         `gorm:"size:255;comment:备注" `
}

func (RechargeOrder) TableName() string {
	return "recharge_records"
}
