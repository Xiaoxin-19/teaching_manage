package dao

import (
	"context"
	"errors"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type TeacherDao interface {
	CreateTeacher(ctx context.Context, t *model.Teacher) error
	UpdateTeacher(ctx context.Context, t *model.Teacher) error
	DeleteTeacher(ctx context.Context, id uint) error
	GetTeacherByID(ctx context.Context, id uint) (*model.Teacher, error)
	GetTeacherList(ctx context.Context, key string, status int, offset int, limit int) ([]model.Teacher, int64, error)
	GetTeacherListUnscoped(ctx context.Context, key string, status int, offset int, limit int) ([]model.Teacher, int64, error)
}

type TeacherGormDao struct {
	db *gorm.DB
}

func NewTeacherDao(db *gorm.DB) TeacherDao {
	return &TeacherGormDao{db: db}
}

func (s TeacherGormDao) CreateTeacher(ctx context.Context, t *model.Teacher) error {
	err := gorm.G[model.Teacher](s.db).Create(ctx, &model.Teacher{
		Name:   t.Name,
		Gender: t.Gender,
		Phone:  t.Phone,
		Status: t.Status,
		Remark: t.Remark,
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	}
	return err
}

func (s TeacherGormDao) UpdateTeacher(ctx context.Context, t *model.Teacher) error {
	_, err := gorm.G[model.Teacher](s.db).Where("id = ?", t.ID).Select("name", "gender", "phone", "status", "remark").Updates(ctx, model.Teacher{
		Name:   t.Name,
		Gender: t.Gender,
		Phone:  t.Phone,
		Status: t.Status,
		Remark: t.Remark,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s TeacherGormDao) DeleteTeacher(ctx context.Context, id uint) error {
	_, err := gorm.G[model.Teacher](s.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (s TeacherGormDao) GetTeacherByID(ctx context.Context, id uint) (*model.Teacher, error) {
	t, err := gorm.G[model.Teacher](s.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Get teacher list
func (s TeacherGormDao) GetTeacherList(ctx context.Context, key string, status int, offset int, limit int) ([]model.Teacher, int64, error) {
	var teachers []model.Teacher
	query := gorm.G[model.Teacher](s.db).Where("")

	if key != "" {
		query = query.Where("(name LIKE ? OR phone LIKE ? OR teacher_number LIKE ?)", "%"+key+"%", "%"+key+"%", "%"+key+"%")
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}
	total, err := query.Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}

	// 处理没有分页参数的情况
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	teachers, err = query.Find(ctx)
	return teachers, total, nil
}

// Get teacher list including soft-deleted records
func (s TeacherGormDao) GetTeacherListUnscoped(ctx context.Context, key string, status int, offset int, limit int) ([]model.Teacher, int64, error) {
	var teachers []model.Teacher
	query := s.db.Unscoped().WithContext(ctx).Model(&model.Teacher{})
	if key != "" {
		likeKey := "%" + key + "%"
		query = query.Where("(name LIKE ? OR phone LIKE ? OR teacher_number LIKE ?)", likeKey, likeKey, likeKey)
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}
	total := int64(0)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	teachers = []model.Teacher{}
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err = query.Find(&teachers).Error
	if err != nil {
		return nil, 0, err
	}
	return teachers, total, nil
}
