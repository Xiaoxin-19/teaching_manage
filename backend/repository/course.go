package repository

import (
	"context"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, sc entity.StudentSubject) error
	GetCourseList(ctx context.Context,
		students []uint, subjects []uint, teachers []uint, min *int, max *int, statuses []int, keyword string, offset int, limit int) ([]entity.StudentSubject, int64, error)
	GetCourseByID(ctx context.Context, id uint) (*entity.StudentSubject, error)
	UpdateCourse(ctx context.Context, id uint, teacherID uint, remark string) error
	RechargeCourse(ctx context.Context, id uint, hours int) error
	ToggleStatus(ctx context.Context, id uint) error
	DeleteCourse(ctx context.Context, id uint, isHardDelete bool, remark string) error
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

	return cr.dao.UpdateStatus(ctx, id, newStatus)
}

func (cr CourseRepositoryImpl) DeleteCourse(ctx context.Context, id uint, isHardDelete bool, remark string) error {
	if isHardDelete {
		return cr.dao.Delete(ctx, id)
	}
	return cr.dao.FinishCourse(ctx, id, remark)
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
		Remark:    sc.Remark,
		CreatedAt: sc.CreatedAt,
		UpdatedAt: sc.UpdatedAt,
	}, nil
}

func (cr CourseRepositoryImpl) UpdateCourse(ctx context.Context, id uint, teacherID uint, remark string) error {
	return cr.dao.UpdateStudentCourseInfo(ctx, id, teacherID, remark)
}

func (cr CourseRepositoryImpl) RechargeCourse(ctx context.Context, id uint, hours int) error {
	return cr.dao.Recharge(ctx, id, hours)
}
