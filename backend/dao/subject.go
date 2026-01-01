package dao

import (
	"context"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type SubjectDao interface {
	CreateSubject(ctx context.Context, subject *model.Subject) error
	GetSubjectByName(ctx context.Context, name string) (*model.Subject, error)
	GetAllSubjects(ctx context.Context) ([]model.Subject, error)
}

type SubjectGormDao struct {
	db *gorm.DB
}

func NewSubjectDao(db *gorm.DB) SubjectDao {
	return &SubjectGormDao{db: db}
}

func (d *SubjectGormDao) CreateSubject(ctx context.Context, subject *model.Subject) error {
	return d.db.WithContext(ctx).Create(subject).Error
}

func (d *SubjectGormDao) GetSubjectByName(ctx context.Context, name string) (*model.Subject, error) {
	var subject model.Subject
	err := d.db.WithContext(ctx).Where("name = ?", name).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

func (d *SubjectGormDao) GetAllSubjects(ctx context.Context) ([]model.Subject, error) {
	var subjects []model.Subject
	err := d.db.WithContext(ctx).Find(&subjects).Error
	return subjects, err
}
