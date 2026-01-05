package service

import (
	"context"
	"fmt"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	"teaching_manage/backend/repository"
	requestx "teaching_manage/backend/service/request"
	responsex "teaching_manage/backend/service/response"

	"gorm.io/gorm"
)

type SubjectManager struct {
	Ctx        context.Context
	repo       repository.SubjectRepository
	repoCourse repository.CourseRepository
}

func NewSubjectManager(repo repository.SubjectRepository, repoCourse repository.CourseRepository) *SubjectManager {
	return &SubjectManager{repo: repo, repoCourse: repoCourse}
}

func (sm SubjectManager) GetSubjectList(ctx context.Context, req *requestx.GetSubjectListRequest) (responsex.GetSubjectListResponse, error) {
	var subjects []entity.Subject
	var total int64
	var err error
	if req.ShowDeleted {
		subjects, total, err = sm.repo.GetSubjectListUnscoped(ctx, req.KeyWord, req.Offset, req.Limit)
		if err != nil {
			return responsex.GetSubjectListResponse{}, err
		}
	} else {
		subjects, total, err = sm.repo.GetSubjectList(ctx, req.KeyWord, req.Offset, req.Limit)
		if err != nil {
			return responsex.GetSubjectListResponse{}, err
		}
	}

	dtoSubjects := make([]responsex.SubjectDTO, 0, len(subjects))
	for _, s := range subjects {
		effectiveCourseCount := 0
		for _, cs := range s.StudentSubjects {
			if cs.Status != 3 { // not "已取消"
				effectiveCourseCount++
			}
		}
		deletedAt := int64(0)
		if !s.DeletedAt.IsZero() {
			deletedAt = s.DeletedAt.UnixMilli()
		}
		dtoSubjects = append(dtoSubjects, responsex.SubjectDTO{
			ID:            s.ID,
			SubjectNumber: s.SubjectNumber,
			Name:          s.Name,
			CreatedAt:     s.CreatedAt.UnixMilli(),
			UpdatedAt:     s.UpdatedAt.UnixMilli(),
			DeletedAt:     deletedAt,
			StudentCount:  int64(effectiveCourseCount),
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

	// 检查是否有正在进行中的课程
	courses, err := sm.repoCourse.GetBySubjectID(ctx, req.ID)
	if err != nil {
		logger.Error("failed to get courses by subject id", logger.ErrorType(err))
		return "", fmt.Errorf("failed to get courses by subject id: %w", err)
	}

	countEffectiveCourses := 0
	for _, c := range courses {
		if c.Status != 3 { // not "已取消"
			countEffectiveCourses++
		}
	}

	if countEffectiveCourses > 0 {
		logger.Error("cannot delete subject with active courses",
			logger.UInt("subject_id", req.ID),
			logger.Int("active_course_count", countEffectiveCourses),
		)
		return "", fmt.Errorf("无法删除，有 %d 个正在进行中的课程的科目", countEffectiveCourses)
	}

	db := dao.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewSubjectRepository(dao.NewSubjectDao(dao.GetDBTarget(tx)))
		txCourseRepo := repository.NewCourseRepository(dao.NewStudentCourseDao(dao.GetDBTarget(tx)))
		// 删除该科目的选课记录
		err = txCourseRepo.DeleteBySubjectID(ctx, req.ID)
		if err != nil {
			logger.Error("failed to delete subject courses in transaction", logger.ErrorType(err))
			return fmt.Errorf("failed to delete subject courses: %w", err)
		}
		// 删除科目记录
		err = txRepo.DeleteSubject(ctx, req.ID)
		if err != nil {
			logger.Error("failed to delete subject in transaction", logger.ErrorType(err))
			return fmt.Errorf("failed to delete subject: %w", err)
		}
		return nil
	})
	if err != nil {
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
