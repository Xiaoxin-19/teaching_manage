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
	GetSubjectListUnscoped(ctx context.Context, keyword string, offset int, limit int) ([]model.Subject, int64, error)
}

type SubjectGormDao struct {
	getDB DBGetter
}

func NewSubjectDao(getDB DBGetter) SubjectDao {
	return &SubjectGormDao{getDB: getDB}
}

func (d *SubjectGormDao) CreateSubject(ctx context.Context, subject *model.Subject) error {
	db := GetDBFromCtx(ctx, d.getDB())
	err := gorm.G[model.Subject](db).Create(ctx, subject)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	}
	return err
}

func (d *SubjectGormDao) GetSubjectByName(ctx context.Context, name string) (*model.Subject, error) {
	var subject model.Subject
	db := GetDBFromCtx(ctx, d.getDB())
	subject, err := gorm.G[model.Subject](db).Where("name = ?", name).First(ctx)
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
	db := GetDBFromCtx(ctx, d.getDB())
	subject, err := gorm.G[model.Subject](db).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &subject, nil
}

func (d *SubjectGormDao) UpdateSubject(ctx context.Context, subject *model.Subject) error {
	db := GetDBFromCtx(ctx, d.getDB())
	_, err := gorm.G[model.Subject](db).Select("name").Where("id = ?", subject.ID).
		Updates(ctx, model.Subject{
			Name: subject.Name,
		})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	}
	return err
}

func (d *SubjectGormDao) DeleteSubject(ctx context.Context, id uint) error {
	db := GetDBFromCtx(ctx, d.getDB())
	_, err := gorm.G[model.Subject](db).Where("id = ?", id).Delete(ctx)
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
	db := GetDBFromCtx(ctx, d.getDB())
	query := gorm.G[model.Subject](db).Preload("StudentSubjects", nil)
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	total, err := query.Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	subjects, err = query.Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return subjects, total, nil
}

func (d *SubjectGormDao) GetSubjectListUnscoped(ctx context.Context, keyword string, offset int, limit int) ([]model.Subject, int64, error) {
	var subjects []model.Subject
	db := GetDBFromCtx(ctx, d.getDB())
	query := db.Unscoped().WithContext(ctx).Model(&model.Subject{}).Preload("StudentSubjects", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	total := int64(0)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err = query.Find(&subjects).Error
	if err != nil {
		return nil, 0, err
	}
	return subjects, total, nil
}
