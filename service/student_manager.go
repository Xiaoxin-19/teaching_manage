package service

import (
	"context"
	"encoding/json"
	"teaching_manage/entity"
	"teaching_manage/pkg/wraper"
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

func (sm StudentManager) GetStudentList(r string) string {
	var req requestx.GetStudentListRequest
	json.Unmarshal([]byte(r), &req)

	studentDs, err := sm.repo.GetStudentList(context.Background(), req.Key, req.Offset, req.Limit)
	if err != nil {
		return wraper.NewBadResponse("internal server error").ToJSON()
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

	// return json string
	return wraper.NewSuccessResponse(result).ToJSON()
}
