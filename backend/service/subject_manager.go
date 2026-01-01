package service

import (
	"context"
	"fmt"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	"teaching_manage/backend/repository"
	requestx "teaching_manage/backend/service/request"
	responsex "teaching_manage/backend/service/response"
)

type SubjectManager struct {
	Ctx  context.Context
	repo repository.SubjectRepository
}

func NewSubjectManager(repo repository.SubjectRepository) *SubjectManager {
	return &SubjectManager{repo: repo}
}

func (sm SubjectManager) GetSubjectList(ctx context.Context, req *requestx.GetSubjectListRequest) (responsex.GetSubjectListResponse, error) {
	subjects, total, err := sm.repo.GetSubjectList(ctx, req.KeyWord, req.Offset, req.Limit)
	if err != nil {
		return responsex.GetSubjectListResponse{}, err
	}

	dtoSubjects := make([]responsex.SubjectDTO, 0, len(subjects))
	for _, s := range subjects {
		dtoSubjects = append(dtoSubjects, responsex.SubjectDTO{
			ID:            s.ID,
			SubjectNumber: s.SubjectNumber,
			Name:          s.Name,
			CreatedAt:     s.CreatedAt.UnixMilli(),
			UpdatedAt:     s.UpdatedAt.UnixMilli(),
			StudentCount:  int64(len(s.StudentSubjects)),
		})
	}
	return responsex.GetSubjectListResponse{
		Subjects: dtoSubjects,
		Total:    total,
	}, nil
}

func (sm SubjectManager) CreateSubject(ctx context.Context, req *requestx.CreateSubjectRequest) (string, error) {
	logger.Info("Creating one subject",
		logger.String("subject_name", req.Name),
	)
	err := sm.repo.CreateSubject(ctx, entity.Subject{
		Name: req.Name,
	})
	if err != nil {
		logger.Error("failed to create subject", logger.ErrorType(err))
		return "", fmt.Errorf("科目 %s 已存在，无法重复创建", req.Name)
	}
	return "subject created", nil
}

func (sm SubjectManager) UpdateSubject(ctx context.Context, req *requestx.UpdateSubjectRequest) (string, error) {
	logger.Info("Updating one subject",
		logger.UInt("subject_id", req.ID),
		logger.String("subject_name", req.Name),
	)
	err := sm.repo.UpdateSubject(ctx, entity.Subject{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		logger.Error("failed to update subject", logger.ErrorType(err))
		return "", err
	}
	return "subject updated", nil
}

func (sm SubjectManager) DeleteSubject(ctx context.Context, req *requestx.DeleteSubjectRequest) (string, error) {
	logger.Info("Deleting one subject",
		logger.UInt("subject_id", req.ID),
	)
	err := sm.repo.DeleteSubject(ctx, req.ID)
	if err != nil {
		logger.Error("failed to delete subject", logger.ErrorType(err))
		return "", err
	}
	return "subject deleted", nil
}

func (sm SubjectManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "subject_manager/get_subject_list", sm.GetSubjectList)
	dispatcher.RegisterTyped(d, "subject_manager/create_subject", sm.CreateSubject)
	dispatcher.RegisterTyped(d, "subject_manager/update_subject", sm.UpdateSubject)
	dispatcher.RegisterTyped(d, "subject_manager/delete_subject", sm.DeleteSubject)
}
