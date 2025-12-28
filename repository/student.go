package repository

import (
	"context"
	"teaching_manage/dao"
)

type StudentRepository interface {
	GetStudentList(ctx context.Context, key string, offset int, limit int) ([]dao.Student, error)
	GetStudentByID(ctx context.Context, id uint) (*dao.Student, error)
	UpdateStudentByID(ctx context.Context, stu *dao.Student) error
}

type StudentRepositoryImpl struct {
	dao dao.StudentDao
}

func NewStudentRepository(dao dao.StudentDao) StudentRepository {
	return &StudentRepositoryImpl{dao: dao}
}
func (sr StudentRepositoryImpl) GetStudentList(ctx context.Context, key string, offset int, limit int) ([]dao.Student, error) {
	students, err := sr.dao.GetStudentList(ctx, key, offset, limit)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (sr StudentRepositoryImpl) GetStudentByID(ctx context.Context, id uint) (*dao.Student, error) {
	student, err := sr.dao.GetStudentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (sr StudentRepositoryImpl) UpdateStudentByID(ctx context.Context, stu *dao.Student) error {
	return sr.dao.UpdateStudent(ctx, stu)
}
