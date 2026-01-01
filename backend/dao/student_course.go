package dao

import (
	"context"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type StudentCourseDao interface {
	CreateStudentCourse(ctx context.Context, sc *model.StudentSubject) error
	GetStudentCourse(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error)
	GetStudentCourseWithDeleted(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error)
	UpdateBalance(ctx context.Context, id uint, delta int) error
	RestoreStudentCourse(ctx context.Context, id uint) error
	GetCoursesByStudentID(ctx context.Context, studentID uint) ([]model.StudentSubject, error)
}

type StudentCourseGormDao struct {
	db *gorm.DB
}

func NewStudentCourseDao(db *gorm.DB) StudentCourseDao {
	return &StudentCourseGormDao{db: db}
}

func (d *StudentCourseGormDao) CreateStudentCourse(ctx context.Context, sc *model.StudentSubject) error {
	return d.db.WithContext(ctx).Create(sc).Error
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

func (d *StudentCourseGormDao) RestoreStudentCourse(ctx context.Context, id uint) error {
	return d.db.Unscoped().WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

func (d *StudentCourseGormDao) GetCoursesByStudentID(ctx context.Context, studentID uint) ([]model.StudentSubject, error) {
	var courses []model.StudentSubject
	err := d.db.WithContext(ctx).
		Preload("Subject").
		Preload("Teacher").
		Where("student_id = ?", studentID).
		Find(&courses).Error
	return courses, err
}
