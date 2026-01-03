package dao

import (
	"context"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type RechargeOrderDao interface {
	CreateRechargeRecord(ctx context.Context, record *model.RechargeOrder) error
	GetRechargeRecords(ctx context.Context, studentID uint, offset, limit int) ([]model.RechargeOrder, int64, error)
	GetRechargeOrderList(ctx context.Context, studentID uint, subjectIDs []uint, orderType []string, dateStart string, dateEnd string, offset int, limit int) ([]model.RechargeOrder, int64, error)
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

func (d *RechargeOrderGormDao) GetRechargeOrderList(ctx context.Context, studentID uint, subjectIDs []uint, orderType []string, dateStart string, dateEnd string, offset int, limit int) ([]model.RechargeOrder, int64, error) {
	query := d.db.WithContext(ctx).Model(&model.RechargeOrder{}).Joins("StudentCourse").
		Preload("StudentCourse.Student").
		Preload("StudentCourse.Subject").
		Preload("StudentCourse.Teacher")

	// 过滤条件
	if studentID > 0 {
		query = query.Where("StudentCourse.student_id = ?", studentID)
	}

	if len(subjectIDs) > 0 {
		query = query.Where("StudentCourse.subject_id IN ?", subjectIDs)
	}

	if len(orderType) > 0 {
		for _, t := range orderType {
			switch t {
			case "increase":
				query = query.Where("hours >= 0")
			case "decrease":
				query = query.Where("hours < 0")
			}
		}
	}

	if dateStart != "" {
		query = query.Where("recharge_orders.created_at >= ?", dateStart)
	}

	if dateEnd != "" {
		query = query.Where("recharge_orders.created_at <= ?", dateEnd)
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var orders []model.RechargeOrder
	err = query.Order("recharge_orders.created_at desc").Offset(offset).Limit(limit).Find(&orders).Error
	return orders, total, err
}
