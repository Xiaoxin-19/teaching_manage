package repository

import (
	"context"
	"fmt"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg/logger"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, sc entity.StudentSubject) error

	// Get
	GetCourseList(ctx context.Context,
		students []uint, subjects []uint, teachers []uint, min *int, max *int, statuses []int, keyword string, offset int, limit int) ([]entity.StudentSubject, int64, error)
	GetCourseByID(ctx context.Context, id uint) (*entity.StudentSubject, error)
	GetByStudentIDAndSubjectID(ctx context.Context, studentID uint, subjectID uint) (*entity.StudentSubject, error)
	GetByStudentID(ctx context.Context, d uint) ([]entity.StudentSubject, error)
	GetByTeacherID(ctx context.Context, id uint) ([]entity.StudentSubject, error)
	GetBySubjectID(ctx context.Context, subjectID uint) ([]entity.StudentSubject, error)
	// Update
	UpdateStatus(ctx context.Context, id uint, status int, remark string) error
	UpdateCourseTeacher(ctx context.Context, id uint, teacherID uint, remark string) error
	UpdateBalance(ctx context.Context, id uint, hours int) error
	RechargeCourse(ctx context.Context, id uint, hours int, amount float64) error
	ToggleStatus(ctx context.Context, id uint) error

	// Delete
	DeleteByStudentID(ctx context.Context, stuID uint) error
	DeleteByTeacherID(ctx context.Context, teacherID uint) error
	DeleteBySubjectID(ctx context.Context, subjectID uint) error
	DeleteCourse(ctx context.Context, id uint) error
}

type CourseRepositoryImpl struct {
	dao dao.StudentCourseDao
}

func NewCourseRepository(d dao.StudentCourseDao) CourseRepository {
	return &CourseRepositoryImpl{
		dao: d,
	}
}

func (cr CourseRepositoryImpl) CreateCourse(ctx context.Context, sc entity.StudentSubject) error {
	modelSc := model.StudentSubject{
		StudentID: sc.Student.ID,
		SubjectID: sc.Subject.ID,
		TeacherID: sc.Teacher.ID,
		Remark:    sc.Remark,
		Balance:   0,
		TotalBuy:  0,
	}
	return cr.dao.CreateStudentCourse(ctx, &modelSc)
}

func (cr CourseRepositoryImpl) GetCourseList(ctx context.Context,
	students []uint, subjects []uint, teachers []uint, min *int, max *int, statuses []int, keyword string, offset int, limit int) ([]entity.StudentSubject, int64, error) {
	modelScs, total, err := cr.dao.GetStudentCourseList(ctx, students, subjects, teachers, min, max, statuses, keyword, offset, limit)

	if err != nil {
		return nil, 0, err
	}
	var res []entity.StudentSubject
	for _, sc := range modelScs {
		res = append(res, entity.StudentSubject{
			ID:        sc.ID,
			Student:   entity.Student{ID: sc.StudentID, Name: sc.Student.Name, StudentNumber: sc.Student.StudentNumber, Status: int(sc.Student.Status)},
			Subject:   entity.Subject{ID: sc.SubjectID, Name: sc.Subject.Name},
			Teacher:   entity.Teacher{ID: sc.TeacherID, Name: sc.Teacher.Name, TeacherNumber: sc.Teacher.TeacherNumber},
			Balance:   sc.Balance,
			AvgPrice:  sc.AvgPrice,
			Remark:    sc.Remark,
			Status:    entity.ParseStudentSubjectStatus(sc.Status),
			CreatedAt: sc.CreatedAt,
			UpdatedAt: sc.UpdatedAt,
		})
	}

	return res, total, nil
}

func (cr CourseRepositoryImpl) ToggleStatus(ctx context.Context, id uint) error {
	sc, err := cr.dao.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var newStatus int
	switch sc.Status {
	case 1:
		newStatus = 2
	case 2:
		newStatus = 1
	default:
		return nil
	}

	return cr.dao.UpdateStatus(ctx, id, newStatus, "")
}

func (cr CourseRepositoryImpl) DeleteCourse(ctx context.Context, id uint) error {
	return cr.dao.Delete(ctx, id)
}

func (cr CourseRepositoryImpl) GetCourseByID(ctx context.Context, id uint) (*entity.StudentSubject, error) {
	sc, err := cr.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.StudentSubject{
		ID:        sc.ID,
		Student:   entity.Student{ID: sc.Student.ID, Name: sc.Student.Name, StudentNumber: sc.Student.StudentNumber, Status: int(sc.Student.Status)},
		Subject:   entity.Subject{ID: sc.Subject.ID, Name: sc.Subject.Name},
		Teacher:   entity.Teacher{ID: sc.Teacher.ID, Name: sc.Teacher.Name, TeacherNumber: sc.Teacher.TeacherNumber},
		Status:    entity.StudentSubjectStatus(sc.Status),
		Balance:   sc.Balance,
		AvgPrice:  sc.AvgPrice,
		Remark:    sc.Remark,
		CreatedAt: sc.CreatedAt,
		UpdatedAt: sc.UpdatedAt,
	}, nil
}

func (cr CourseRepositoryImpl) UpdateCourseTeacher(ctx context.Context, id uint, teacherID uint, remark string) error {
	return cr.dao.UpdateStudentCourseInfo(ctx, id, teacherID, remark)
}

func (cr CourseRepositoryImpl) UpdateBalance(ctx context.Context, id uint, hours int) error {
	return cr.dao.UpdateBalance(ctx, id, hours)
}

func (cr CourseRepositoryImpl) RechargeCourse(ctx context.Context, id uint, hours int, amount float64) error {
	return cr.dao.Recharge(ctx, id, hours, amount)
}

func (cr CourseRepositoryImpl) GetByStudentIDAndSubjectID(ctx context.Context, studentID uint, subjectID uint) (*entity.StudentSubject, error) {
	sc, err := cr.dao.GetByStudentIDAndSubjectID(ctx, studentID, subjectID)
	if err != nil {
		return nil, err
	}

	logger.Debug("GetByStudentIDAndSubjectID", logger.UInt("studentID", studentID), logger.UInt("subjectID", subjectID), logger.String("course", fmt.Sprintf("%v", sc)))
	return &entity.StudentSubject{
		ID:        sc.ID,
		Student:   entity.Student{ID: sc.Student.ID, Name: sc.Student.Name, StudentNumber: sc.Student.StudentNumber, Status: int(sc.Student.Status)},
		Subject:   entity.Subject{ID: sc.Subject.ID, Name: sc.Subject.Name},
		Teacher:   entity.Teacher{ID: sc.Teacher.ID, Name: sc.Teacher.Name, TeacherNumber: sc.Teacher.TeacherNumber, Status: int(sc.Teacher.Status), DeletedAt: sc.Teacher.DeletedAt.Time},
		Status:    entity.StudentSubjectStatus(sc.Status),
		Balance:   sc.Balance,
		AvgPrice:  sc.AvgPrice,
		Remark:    sc.Remark,
		CreatedAt: sc.CreatedAt,
		UpdatedAt: sc.UpdatedAt,
	}, nil
}

func (c CourseRepositoryImpl) GetByStudentID(ctx context.Context, d uint) ([]entity.StudentSubject, error) {
	course, err := c.dao.GetByStudentID(ctx, d)
	if err != nil {
		return nil, err
	}
	var result []entity.StudentSubject
	for _, sc := range course {
		result = append(result, entity.StudentSubject{
			ID:        sc.ID,
			Student:   entity.Student{ID: sc.Student.ID, Name: sc.Student.Name, StudentNumber: sc.Student.StudentNumber, Status: int(sc.Student.Status)},
			Balance:   sc.Balance,
			AvgPrice:  sc.AvgPrice,
			Remark:    sc.Remark,
			Status:    entity.ParseStudentSubjectStatus(sc.Status),
			CreatedAt: sc.CreatedAt,
			UpdatedAt: sc.UpdatedAt,
		})
	}
	return result, nil
}

func (c CourseRepositoryImpl) DeleteByStudentID(ctx context.Context, stuID uint) error {
	return c.dao.DeleteByStudentID(ctx, stuID)
}

func (cr CourseRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status int, remark string) error {
	return cr.dao.UpdateStatus(ctx, id, status, remark)
}

func (cr CourseRepositoryImpl) DeleteByTeacherID(ctx context.Context, teacherID uint) error {
	return cr.dao.DeleteByTeacherID(ctx, teacherID)
}

func (cr CourseRepositoryImpl) GetByTeacherID(ctx context.Context, id uint) ([]entity.StudentSubject, error) {
	courses, err := cr.dao.GetByTeacherID(ctx, id)
	if err != nil {
		return nil, err
	}
	var result []entity.StudentSubject
	for _, sc := range courses {
		result = append(result, entity.StudentSubject{
			ID:        sc.ID,
			Balance:   sc.Balance,
			AvgPrice:  sc.AvgPrice,
			Remark:    sc.Remark,
			Status:    entity.ParseStudentSubjectStatus(sc.Status),
			CreatedAt: sc.CreatedAt,
			UpdatedAt: sc.UpdatedAt,
		})
	}
	return result, nil
}

func (cr CourseRepositoryImpl) GetBySubjectID(ctx context.Context, subjectID uint) ([]entity.StudentSubject, error) {
	courses, err := cr.dao.GetBySubjectID(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	var result []entity.StudentSubject
	for _, sc := range courses {
		result = append(result, entity.StudentSubject{
			ID:        sc.ID,
			Balance:   sc.Balance,
			AvgPrice:  sc.AvgPrice,
			Remark:    sc.Remark,
			Status:    entity.ParseStudentSubjectStatus(sc.Status),
			CreatedAt: sc.CreatedAt,
			UpdatedAt: sc.UpdatedAt,
		})
	}
	return result, nil
}

func (cr CourseRepositoryImpl) DeleteBySubjectID(ctx context.Context, subjectID uint) error {
	return cr.dao.DeleteBySubjectID(ctx, subjectID)
}
