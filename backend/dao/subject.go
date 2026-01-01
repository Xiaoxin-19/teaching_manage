package dao

import (
	"context"
	"errors"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type SubjectDao interface {
	CreateSubject(ctx context.Context, subject *model.Subject) error
	GetSubjectByName(ctx context.Context, name string) (*model.Subject, error)
	GetSubjectByID(ctx context.Context, id uint) (*model.Subject, error)
	UpdateSubject(ctx context.Context, subject *model.Subject) error
	DeleteSubject(ctx context.Context, id uint) error
	GetSubjectList(ctx context.Context, keyword string, offset int, limit int) ([]model.Subject, int64, error)
}

type SubjectGormDao struct {
	db *gorm.DB
}

func NewSubjectDao(db *gorm.DB) SubjectDao {
	return &SubjectGormDao{db: db}
}

func (d *SubjectGormDao) CreateSubject(ctx context.Context, subject *model.Subject) error {
	err := gorm.G[model.Subject](d.db).Create(ctx, subject)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	}
	return err
}

func (d *SubjectGormDao) GetSubjectByName(ctx context.Context, name string) (*model.Subject, error) {
	var subject model.Subject
	subject, err := gorm.G[model.Subject](d.db).Where("name = ?", name).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &subject, nil
}

func (d *SubjectGormDao) GetSubjectByID(ctx context.Context, id uint) (*model.Subject, error) {
	var subject model.Subject
	subject, err := gorm.G[model.Subject](d.db).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &subject, nil
}

func (d *SubjectGormDao) UpdateSubject(ctx context.Context, subject *model.Subject) error {
	_, err := gorm.G[model.Subject](d.db).Select("name").Where("id = ?", subject.ID).
		Updates(ctx, model.Subject{
			Name: subject.Name,
		})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	}
	return err
}

func (d *SubjectGormDao) DeleteSubject(ctx context.Context, id uint) error {
	_, err := gorm.G[model.Subject](d.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (d *SubjectGormDao) GetSubjectList(ctx context.Context, keyword string, offset int, limit int) ([]model.Subject, int64, error) {
	var subjects []model.Subject
	query := gorm.G[model.Subject](d.db).Preload("StudentSubjects", nil)
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	total, err := query.Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}
	subjects, err = query.Offset(offset).Limit(limit).Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return subjects, total, nil
}
