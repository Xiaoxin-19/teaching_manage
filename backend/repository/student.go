package repository

import (
	"context"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"

	"gorm.io/gorm"
)

type StudentRepository interface {
	GetStudentList(ctx context.Context, key string, offset int, limit int) ([]entity.Student, int64, error)
	GetStudentByName(ctx context.Context, name string) (*entity.Student, error)
	GetStudentByID(ctx context.Context, id uint) (*entity.Student, error)
	UpdateStudentByID(ctx context.Context, stu *entity.Student) error
	CreateStudent(ctx context.Context, stu *entity.Student) error
	DeleteStudentByID(ctx context.Context, id uint) error
	UpdateStudentHoursByID(ctx context.Context, id uint, diff int) error
	UpdateStudentHoursByIDWithDeleted(ctx context.Context, id uint, diff int) error
	GetStudentByIdWithDeleted(ctx context.Context, id uint) (*entity.Student, error)
}

type StudentRepositoryImpl struct {
	dao dao.StudentDao
}

func NewStudentRepository(dao dao.StudentDao) StudentRepository {
	return &StudentRepositoryImpl{dao: dao}
}
func (sr StudentRepositoryImpl) GetStudentList(ctx context.Context, key string, offset int, limit int) ([]entity.Student, int64, error) {
	students, total, err := sr.dao.GetStudentList(ctx, key, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	var result []entity.Student
	for _, stu := range students {
		result = append(result, entity.Student{
			ID:        stu.ID,
			CreatedAt: stu.CreatedAt,
			UpdatedAt: stu.UpdatedAt,
			Name:      stu.Name,
			Gender:    stu.Gender,
			Hours:     stu.Hours,
			Phone:     stu.Phone,
			Remark:    stu.Remark,
			TeacherID: stu.TeacherID,
		})
	}
	return result, total, nil
}

func (sr StudentRepositoryImpl) GetStudentByName(ctx context.Context, name string) (*entity.Student, error) {
	student, err := sr.dao.GetStudentByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &entity.Student{
		ID:        student.ID,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
		Name:      student.Name,
		Gender:    student.Gender,
		Hours:     student.Hours,
		Phone:     student.Phone,
		Remark:    student.Remark,
		TeacherID: student.TeacherID,
		Teacher: entity.Teacher{
			ID:        student.Teacher.ID,
			Name:      student.Teacher.Name,
			DeletedAt: student.Teacher.DeletedAt.Time,
		},
	}, nil
}

func (sr StudentRepositoryImpl) GetStudentByID(ctx context.Context, id uint) (*entity.Student, error) {
	student, err := sr.dao.GetStudentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Student{
		ID:        student.ID,
		Name:      student.Name,
		Gender:    student.Gender,
		Hours:     student.Hours,
		Phone:     student.Phone,
		TeacherID: student.TeacherID,
		Remark:    student.Remark,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
		Teacher: entity.Teacher{
			ID:        student.Teacher.ID,
			Name:      student.Teacher.Name,
			DeletedAt: student.Teacher.DeletedAt.Time,
		},
	}, nil
}

func (sr StudentRepositoryImpl) GetStudentByIdWithDeleted(ctx context.Context, id uint) (*entity.Student, error) {
	student, err := sr.dao.GetStudentByIdWithDeleted(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Student{
		ID:        student.ID,
		Name:      student.Name,
		Gender:    student.Gender,
		Hours:     student.Hours,
		Phone:     student.Phone,
		TeacherID: student.TeacherID,
		Remark:    student.Remark,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
		Teacher: entity.Teacher{
			ID:        student.Teacher.ID,
			Name:      student.Teacher.Name,
			DeletedAt: student.Teacher.DeletedAt.Time,
		},
	}, nil
}

func (sr StudentRepositoryImpl) UpdateStudentByID(ctx context.Context, stu *entity.Student) error {
	return sr.dao.UpdateStudent(ctx, &dao.Student{
		Model:     gorm.Model{ID: stu.ID},
		Name:      stu.Name,
		Gender:    stu.Gender,
		Hours:     stu.Hours,
		Phone:     stu.Phone,
		TeacherID: stu.TeacherID,
		Remark:    stu.Remark,
	})
}

func (sr StudentRepositoryImpl) UpdateStudentHoursByIDWithDeleted(ctx context.Context, id uint, diff int) error {
	return sr.dao.UpdateStudentHoursWithDeleted(ctx, id, diff)
}

func (sr StudentRepositoryImpl) UpdateStudentHoursByID(ctx context.Context, id uint, diff int) error {
	return sr.dao.UpdateStudentHours(ctx, id, diff)
}

func (sr StudentRepositoryImpl) CreateStudent(ctx context.Context, stu *entity.Student) error {
	return sr.dao.CreateStudent(ctx, &dao.Student{
		Model:     gorm.Model{},
		Name:      stu.Name,
		Gender:    stu.Gender,
		Hours:     stu.Hours,
		Phone:     stu.Phone,
		TeacherID: stu.TeacherID,
		Remark:    stu.Remark,
	})
}

func (sr StudentRepositoryImpl) DeleteStudentByID(ctx context.Context, id uint) error {
	return sr.dao.DeleteStudent(ctx, id)
}
