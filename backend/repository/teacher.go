package repository

import (
	"context"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg"
)

type TeacherRepository interface {
	GetTeacherList(ctx context.Context, key string, offset int, limit int) ([]entity.Teacher, int64, error)
	CreateTeacher(ctx context.Context, teacher entity.Teacher) error
	DeleteTeacher(ctx context.Context, id uint) error
	UpdateTeacher(ctx context.Context, teacher entity.Teacher) error
}

type TeacherRepositoryImpl struct {
	dao dao.TeacherDao
}

func NewTeacherRepository(dao dao.TeacherDao) TeacherRepository {
	return &TeacherRepositoryImpl{dao: dao}
}

func (tr TeacherRepositoryImpl) GetTeacherList(ctx context.Context, key string, offset int, limit int) ([]entity.Teacher, int64, error) {
	teachers, total, err := tr.dao.GetTeacherList(ctx, key, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var entities []entity.Teacher
	for _, t := range teachers {
		entities = append(entities, entity.Teacher{
			ID:            t.ID,
			Name:          t.Name,
			TeacherNumber: t.TeacherNumber,
			Gender:        pkg.Gender(t.Gender),
			Phone:         t.Phone,
			Remark:        t.Remark,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
		})
	}
	return entities, total, nil
}

func (tr TeacherRepositoryImpl) CreateTeacher(ctx context.Context, teacher entity.Teacher) error {
	err := tr.dao.CreateTeacher(ctx, &model.Teacher{
		Name:   teacher.Name,
		Phone:  teacher.Phone,
		Gender: string(teacher.Gender),
		Remark: teacher.Remark,
	})
	return err
}

func (tr TeacherRepositoryImpl) DeleteTeacher(ctx context.Context, id uint) error {
	err := tr.dao.DeleteTeacher(ctx, id)
	return err
}

func (tr TeacherRepositoryImpl) UpdateTeacher(ctx context.Context, teacher entity.Teacher) error {
	t := model.Teacher{
		Name:   teacher.Name,
		Phone:  teacher.Phone,
		Gender: string(teacher.Gender),
		Remark: teacher.Remark,
	}
	t.ID = teacher.ID
	err := tr.dao.UpdateTeacher(ctx, &t)
	return err
}
