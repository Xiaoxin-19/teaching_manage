package service

import (
	"context"
	"teaching_manage/entity"
	"teaching_manage/pkg/dispatcher"
	"teaching_manage/repository"
	requestx "teaching_manage/service/request"
)

type StudentManager struct {
	Ctx  context.Context
	repo repository.StudentRepository
}

func NewStudentManager(repo repository.StudentRepository) *StudentManager {
	return &StudentManager{repo: repo}
}

func (sm StudentManager) GetStudentList(ctx context.Context, req *requestx.GetStudentListRequest) ([]entity.Student, error) {
	studentDs, err := sm.repo.GetStudentList(ctx, req.Key, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	var result []entity.Student
	for _, stu := range studentDs {
		result = append(result, entity.Student{
			ID:        stu.ID,
			CreatedAt: stu.CreatedAt.UnixMilli(),
			UpdatedAt: stu.UpdatedAt.UnixMilli(),
			Name:      stu.Name,
			Gender:    stu.Gender,
			Hours:     stu.Hours,
			Phone:     stu.Phone,
			TeacherID: stu.TeacherID,
		})
	}
	return result, nil
}

func (sm StudentManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "student_manager:get_student_list", sm.GetStudentList)
}
