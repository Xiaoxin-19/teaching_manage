package repository

import (
	"context"
	"fmt"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type StudentRepository interface {
	ListStudentsWithStatus(ctx context.Context, key string, offset int, limit int, levelStatus int, targetStatus int) ([]entity.Student, int64, error)
	GetStudentByName(ctx context.Context, name string) ([]entity.Student, error)
	GetStudentByID(ctx context.Context, id uint) (*entity.Student, error)
	UpdateStudent(ctx context.Context, stu *entity.Student) error
	CreateStudent(ctx context.Context, stu *entity.Student) error
	DeleteStudent(ctx context.Context, id uint) error
	// UpdateStudentHoursByID(ctx context.Context, id uint, diff int) error
	// UpdateStudentHoursByIDWithDeleted(ctx context.Context, id uint, diff int) error
	// GetStudentByIdWithDeleted(ctx context.Context, id uint) (*entity.Student, error)
}

type StudentRepositoryImpl struct {
	dao dao.StudentDao
}

func NewStudentRepository(dao dao.StudentDao) StudentRepository {
	return &StudentRepositoryImpl{dao: dao}
}
func (sr StudentRepositoryImpl) ListStudentsWithStatus(ctx context.Context, key string, offset int, limit int, levelStatus int, targetStatus int) ([]entity.Student, int64, error) {
	if levelStatus <= 0 || levelStatus > 3 {
		return nil, 0, fmt.Errorf("invalid status: %d", levelStatus)
	}
	students, total, err := sr.dao.GetStudentListWithStatus(ctx, key, offset, limit, model.StudentStatus(levelStatus), model.StudentStatus(targetStatus))
	if err != nil {
		return nil, 0, err
	}
	var result []entity.Student
	for _, stu := range students {
		result = append(result, entity.Student{
			ID:            stu.ID,
			StudentNumber: stu.StudentNumber,
			CreatedAt:     stu.CreatedAt,
			UpdatedAt:     stu.UpdatedAt,
			Name:          stu.Name,
			Gender:        stu.Gender,
			Phone:         stu.Phone,
			Status:        int(stu.Status),
			Remark:        stu.Remark,
		})
	}
	return result, total, nil
}

func (sr StudentRepositoryImpl) GetStudentByName(ctx context.Context, name string) ([]entity.Student, error) {
	students, err := sr.dao.GetStudentByName(ctx, name)
	if err != nil {
		return nil, err
	}
	var result []entity.Student
	for _, stu := range students {
		result = append(result, entity.Student{
			ID:            stu.ID,
			StudentNumber: stu.StudentNumber,
			CreatedAt:     stu.CreatedAt,
			UpdatedAt:     stu.UpdatedAt,
			Name:          stu.Name,
			Gender:        stu.Gender,
			Phone:         stu.Phone,
			Status:        int(stu.Status),
			Remark:        stu.Remark,
		})
	}
	return result, nil
}

func (sr StudentRepositoryImpl) GetStudentByID(ctx context.Context, id uint) (*entity.Student, error) {
	student, err := sr.dao.GetStudentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Student{
		ID:            student.ID,
		StudentNumber: student.StudentNumber,
		CreatedAt:     student.CreatedAt,
		UpdatedAt:     student.UpdatedAt,
		Name:          student.Name,
		Gender:        student.Gender,
		Phone:         student.Phone,
		Status:        int(student.Status),
		Remark:        student.Remark,
	}, nil
}

// func (sr StudentRepositoryImpl) GetStudentByIdWithDeleted(ctx context.Context, id uint) (*entity.Student, error) {
// 	student, err := sr.dao.GetStudentByIdWithDeleted(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &entity.Student{
// 		ID:        student.ID,
// 		Name:      student.Name,
// 		Gender:    student.Gender,
// 		Hours:     student.Hours,
// 		Phone:     student.Phone,
// 		TeacherID: student.TeacherID,
// 		Remark:    student.Remark,
// 		CreatedAt: student.CreatedAt,
// 		UpdatedAt: student.UpdatedAt,
// 		Teacher: entity.Teacher{
// 			ID:        student.Teacher.ID,
// 			Name:      student.Teacher.Name,
// 			DeletedAt: student.Teacher.DeletedAt.Time,
// 		},
// 	}, nil
// }

func (sr StudentRepositoryImpl) UpdateStudent(ctx context.Context, stu *entity.Student) error {
	return sr.dao.UpdateStudent(ctx, &model.Student{
		Model:  gorm.Model{ID: stu.ID},
		Name:   stu.Name,
		Gender: stu.Gender,
		Phone:  stu.Phone,
		Remark: stu.Remark,
		Status: model.StudentStatus(stu.Status),
	})
}

// func (sr StudentRepositoryImpl) UpdateStudentHoursByIDWithDeleted(ctx context.Context, id uint, diff int) error {
// 	return sr.dao.UpdateStudentHoursWithDeleted(ctx, id, diff)
// }

// func (sr StudentRepositoryImpl) UpdateStudentHoursByID(ctx context.Context, id uint, diff int) error {
// 	return sr.dao.UpdateStudentHours(ctx, id, diff)
// }

func (sr StudentRepositoryImpl) CreateStudent(ctx context.Context, stu *entity.Student) error {
	return sr.dao.CreateStudent(ctx, &model.Student{
		Name:   stu.Name,
		Gender: stu.Gender,
		Phone:  stu.Phone,
		Remark: stu.Remark,
	})
}

func (sr StudentRepositoryImpl) DeleteStudent(ctx context.Context, id uint) error {
	return sr.dao.DeleteStudent(ctx, id)
}
