package model

import (
	"fmt"
	"teaching_manage/backend/pkg/snowflake"

	"gorm.io/gorm"
)

// RechargeOrder 充值记录表，记录每次给学生课时充值或扣减的明细
type RechargeOrder struct {
	gorm.Model
	OrderNumber     string         `gorm:"column:order_number;size:50;uniqueIndex;comment:'订单编号(业务主键)'"`
	StudentCourseID uint           `gorm:"index;not null;comment:关联的学籍ID" `
	StudentCourse   StudentSubject `gorm:"foreignKey:StudentCourseID" `
	Hours           int            `gorm:"comment:变动课时数(正数充值,负数扣减)" `
	Amount          float64        `gorm:"comment:涉及金额" `
	Remark          string         `gorm:"size:255;comment:备注" `
}

func (RechargeOrder) TableName() string {
	return "recharge_orders"
}

func (r *RechargeOrder) BeforeCreate(tx *gorm.DB) error {
	// 如果没有订单号，则生成
	if r.OrderNumber == "" {
		r.OrderNumber = GenerateOrderNumber("ORD")
	}
	return nil
}

func GenerateOrderNumber(prefix string) string {
	// 使用 Snowflake 算法生成唯一 ID
	id := snowflake.GenerateID()
	return fmt.Sprintf("%s%s", prefix, id)
}
