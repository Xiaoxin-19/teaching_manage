package dao

import (
	"context"
	"errors"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StudentCourseDao interface {
	CreateStudentCourse(ctx context.Context, sc *model.StudentSubject) error
	GetStudentCourse(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error)
	GetStudentCourseWithDeleted(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error)
	UpdateBalance(ctx context.Context, id uint, delta int) error
	Recharge(ctx context.Context, id uint, hours int) error
	RestoreStudentCourse(ctx context.Context, id uint) error
	GetStudentCourseList(ctx context.Context, students []uint, subjects []uint, teachers []uint, min *int, max *int, statuses []int, keyword string, offset int, limit int) ([]model.StudentSubject, int64, error)
	UpdateStatus(ctx context.Context, id uint, status int) error
	GetByID(ctx context.Context, id uint) (*model.StudentSubject, error)
	UpdateStudentCourseInfo(ctx context.Context, id uint, teacherID uint, remark string) error
	FinishCourse(ctx context.Context, id uint, remark string) error
	Delete(ctx context.Context, id uint) error
}

type StudentCourseGormDao struct {
	db *gorm.DB
}

func NewStudentCourseDao(db *gorm.DB) StudentCourseDao {
	return &StudentCourseGormDao{db: db}
}

func (d *StudentCourseGormDao) CreateStudentCourse(ctx context.Context, sc *model.StudentSubject) error {
	err := gorm.G[model.StudentSubject](d.db).Create(ctx, sc)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	}
	return err
}

func (d *StudentCourseGormDao) GetStudentCourseList(ctx context.Context,
	students []uint, subjects []uint, teachers []uint, min *int, max *int, statuses []int, keyword string, offset int, limit int) ([]model.StudentSubject, int64, error) {
	var scs []model.StudentSubject
	query := gorm.G[model.StudentSubject](d.db).Preload("Teacher", nil).Preload("Subject", nil)

	// Join Student table to allow filtering on its status
	query = query.Joins(clause.JoinTarget{Association: "Student"}, func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
		return nil
	})

	if len(statuses) > 0 {
		query = query.Where(`(CASE 
			WHEN "Student"."status" = 3 THEN 5 
			WHEN "Student"."status" = 2 THEN 4 
			ELSE "student_subjects"."status" 
		END) IN ?`, statuses)
	}

	// apply filters
	if len(students) > 0 {
		query = query.Where("student_subjects.student_id IN ?", students)
	}
	if len(subjects) > 0 {
		query = query.Where("student_subjects.subject_id IN ?", subjects)
	}
	if len(teachers) > 0 {
		query = query.Where("student_subjects.teacher_id IN ?", teachers)
	}
	if min != nil {
		query = query.Where("student_subjects.balance >= ?", *min)
	}
	if max != nil {
		query = query.Where("student_subjects.balance <= ?", *max)
	}

	// count total records
	total, err := query.Count(ctx, "*")
	if err != nil {
		logger.Error("failed to count student courses",
			logger.ErrorType(err),
		)
		return nil, 0, err
	}
	// apply pagination
	scs, err = query.Offset(offset).Limit(limit).Find(ctx)
	if err != nil {
		logger.Error("failed to get student courses",
			logger.ErrorType(err),
		)
		return nil, 0, err
	}
	return scs, total, nil
}

func (d *StudentCourseGormDao) GetStudentCourse(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error) {
	var sc model.StudentSubject
	err := d.db.WithContext(ctx).
		Where("student_id = ? AND subject_id = ?", studentID, subjectID).
		First(&sc).Error
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (d *StudentCourseGormDao) GetStudentCourseWithDeleted(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error) {
	var sc model.StudentSubject
	err := d.db.Unscoped().WithContext(ctx).
		Where("student_id = ? AND subject_id = ?", studentID, subjectID).
		First(&sc).Error
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (d *StudentCourseGormDao) UpdateBalance(ctx context.Context, id uint, delta int) error {
	return d.db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		UpdateColumn("balance", gorm.Expr("balance + ?", delta)).Error
}

func (d *StudentCourseGormDao) Recharge(ctx context.Context, id uint, hours int) error {
	updates := map[string]interface{}{
		"balance": gorm.Expr("balance + ?", hours),
	}
	if hours > 0 {
		updates["total_buy"] = gorm.Expr("total_buy + ?", hours)
	}
	return d.db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (d *StudentCourseGormDao) RestoreStudentCourse(ctx context.Context, id uint) error {
	return d.db.Unscoped().WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

func (d *StudentCourseGormDao) UpdateStatus(ctx context.Context, id uint, status int) error {
	return d.db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (d *StudentCourseGormDao) GetByID(ctx context.Context, id uint) (*model.StudentSubject, error) {
	var sc model.StudentSubject
	err := d.db.WithContext(ctx).Preload("Student").Preload("Subject").Preload("Teacher").First(&sc, id).Error
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (d *StudentCourseGormDao) UpdateStudentCourseInfo(ctx context.Context, id uint, teacherID uint, remark string) error {
	return d.db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"teacher_id": teacherID,
			"remark":     remark,
		}).Error
}

func (d *StudentCourseGormDao) Delete(ctx context.Context, id uint) error {
	ss := model.StudentSubject{Model: gorm.Model{ID: id}}
	return d.db.WithContext(ctx).Delete(&ss).Error
}

func (d *StudentCourseGormDao) FinishCourse(ctx context.Context, id uint, remark string) error {
	return d.db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": 3,
			"remark": remark,
		}).Error
}
