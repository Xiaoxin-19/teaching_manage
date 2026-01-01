package dao

import (
	"context"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type RechargeOrderDao interface {
	CreateRechargeRecord(ctx context.Context, record *model.RechargeOrder) error
	GetRechargeRecords(ctx context.Context, studentID uint, offset, limit int) ([]model.RechargeOrder, int64, error)
}

type RechargeOrderGormDao struct {
	db *gorm.DB
}

func NewRechargeOrderDao(db *gorm.DB) RechargeOrderDao {
	return &RechargeOrderGormDao{db: db}
}

func (d *RechargeOrderGormDao) CreateRechargeRecord(ctx context.Context, record *model.RechargeOrder) error {
	return d.db.WithContext(ctx).Create(record).Error
}

func (d *RechargeOrderGormDao) GetRechargeRecords(ctx context.Context, studentID uint, offset, limit int) ([]model.RechargeOrder, int64, error) {
	var records []model.RechargeOrder
	var total int64

	query := d.db.WithContext(ctx).Model(&model.RechargeOrder{})
	if studentID > 0 {
		// 需要关联查询，因为 recharge_record 只有 student_course_id
		query = query.Joins("JOIN student_courses ON student_courses.id = recharge_records.student_course_id").
			Where("student_courses.student_id = ?", studentID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("StudentCourse.Subject").
		Preload("StudentCourse.Student").
		Order("created_at desc").
		Offset(offset).Limit(limit).
		Find(&records).Error

	return records, total, err
}
