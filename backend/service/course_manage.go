package service

import (
	"context"
	"errors"
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

type CourseManager struct {
	Ctx     context.Context
	repo    repository.CourseRepository
	stuRepo repository.StudentRepository
}

func NewCourseManager(repo repository.CourseRepository, stuRepo repository.StudentRepository) *CourseManager {
	return &CourseManager{repo: repo, stuRepo: stuRepo}
}

func (cm CourseManager) CreateCourse(ctx context.Context, req *requestx.CreateCourseRequest) (string, error) {
	logger.Info("Creating one course",
		logger.UInt("student_id", req.StudentId),
		logger.UInt("subject_id", req.SubjectId),
		logger.UInt("teacher_id", req.TeacherId),
		logger.String("remark", req.Remark),
	)

	// check student status
	// Student Status: 3 = 退学
	student, err := cm.stuRepo.GetStudentByID(ctx, req.StudentId)
	if err != nil {
		logger.Error("failed to get student", logger.ErrorType(err))
		return "", fmt.Errorf("failed to get student: %w", err)
	}
	if student.Status >= 2 {
		logger.Warn("try to create course for a wrong student", logger.UInt("student_id", req.StudentId), logger.Int("student_status", student.Status))
		return "", fmt.Errorf("学员状态异常(停课或退学)，无法选课")
	}

	err = cm.repo.CreateCourse(ctx, entity.StudentSubject{
		Student: entity.Student{ID: req.StudentId},
		Subject: entity.Subject{ID: req.SubjectId},
		Teacher: entity.Teacher{ID: req.TeacherId},
		Remark:  req.Remark,
	})

	if errors.Is(err, dao.ErrDuplicatedKey) {
		logger.Error("failed to create course: duplicated key", logger.ErrorType(err))
		return "", fmt.Errorf("该学生已选修此科目，无法重复选课")
	}

	if err != nil {
		return "", err
	}
	return "course created", nil
}

func (cm CourseManager) GetCourseList(ctx context.Context, req *requestx.GetCourseListRequest) (*responsex.GetCourseListResponse, error) {
	ptrToStr := func(p *int) string {
		if p == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", *p)
	}
	logger.Info("Getting course list",
		logger.String("student_id", fmt.Sprintf("%v", req.StudentIds)),
		logger.String("subject_id", fmt.Sprintf("%v", req.SubjectIds)),
		logger.String("teacher_id", fmt.Sprintf("%v", req.TeacherIds)),
		logger.String("status", fmt.Sprintf("%v", req.Statuses)),
		logger.String("balance_min", ptrToStr(req.BalanceMin)),
		logger.String("balance_max", ptrToStr(req.BalanceMax)),
		logger.Int("offset", req.Offset),
		logger.Int("limit", req.Limit),
	)

	courses, total, err := cm.repo.GetCourseList(ctx,
		req.StudentIds, req.SubjectIds, req.TeacherIds, req.BalanceMin, req.BalanceMax, req.Statuses, req.Keyword, req.Offset, req.Limit)
	if err != nil {
		logger.Error("failed to get course list", logger.ErrorType(err))
		return nil, fmt.Errorf("failed to get course list: %w", err)
	}
	courseDTOs := make([]responsex.CourseDTO, 0, len(courses))
	for _, c := range courses {
		courseDTOs = append(courseDTOs, responsex.CourseDTO{
			ID:        c.ID,
			Student:   responsex.StudentDTO{ID: c.Student.ID, Name: c.Student.Name, StudentNumber: c.Student.StudentNumber, Status: c.Student.Status},
			Subject:   responsex.SubjectDTO{ID: c.Subject.ID, Name: c.Subject.Name},
			Teacher:   responsex.TeacherDTO{ID: c.Teacher.ID, Name: c.Teacher.Name, TeacherNumber: c.Teacher.TeacherNumber},
			Status:    int(c.Status),
			Balance:   c.Balance,
			Remark:    c.Remark,
			CreatedAt: c.CreatedAt.UnixMilli(),
			UpdatedAt: c.UpdatedAt.UnixMilli(),
		})

	}
	return &responsex.GetCourseListResponse{
		Courses: courseDTOs,
		Total:   total,
	}, nil
}

func (cm CourseManager) ToggleStatus(ctx context.Context, req *requestx.ToggleCourseStatusRequest) (string, error) {
	logger.Info("Toggling course status", logger.UInt("course_id", req.CourseId))
	err := cm.repo.ToggleStatus(ctx, req.CourseId)
	if err != nil {
		logger.Error("failed to toggle course status", logger.ErrorType(err))
		return "", fmt.Errorf("failed to toggle course status: %w", err)
	}
	return "status updated", nil
}

func (cm CourseManager) DeleteCourse(ctx context.Context, req *requestx.DeleteCourseRequest) (string, error) {
	logger.Info("Deleting course", logger.UInt("course_id", req.CourseId), logger.String("is_hard_delete", fmt.Sprintf("%v", req.IsHardDelete)))
	err := cm.repo.DeleteCourse(ctx, req.CourseId, req.IsHardDelete, req.Remark)
	if errors.Is(err, dao.ErrRecordNotFound) {
		logger.Error("failed to delete course: record not found", logger.ErrorType(err))
		return "course deleted", nil
	}

	if err != nil {
		logger.Error("failed to delete course", logger.ErrorType(err))
		return "", fmt.Errorf("failed to delete course: %w", err)
	}
	return "course deleted", nil
}

func (cm CourseManager) UpdateCourse(ctx context.Context, req *requestx.UpdateCourseRequest) (string, error) {
	logger.Info("Updating course", logger.UInt("course_id", req.ID))

	// 1. 获取课程信息
	course, err := cm.repo.GetCourseByID(ctx, req.ID)
	if err != nil {
		logger.Error("failed to get course", logger.ErrorType(err))
		return "", fmt.Errorf("failed to get course: %w", err)
	}

	// 2. 验证学员状态
	// Student Status: 3 = 退学
	if course.Student.Status == 3 {
		return "", fmt.Errorf("学员已退学，无法更新课程信息")
	}

	// 3. 验证课程状态
	// Course Status: 3 = 已结课
	if course.Status == 3 {
		return "", fmt.Errorf("课程已结课，无法更新信息")
	}

	// 4. 更新课程信息
	err = cm.repo.UpdateCourse(ctx, req.ID, req.TeacherId, req.Remark)
	if err != nil {
		logger.Error("failed to update course", logger.ErrorType(err))
		return "", fmt.Errorf("failed to update course: %w", err)
	}

	return "course updated", nil
}

func (cm CourseManager) RechargeCourse(ctx context.Context, req *requestx.RechargeCourseRequest) (string, error) {
	logger.Info("Recharging course", logger.UInt("course_id", req.CourseId), logger.Int("hours", req.Hours))

	// 1. Get existing course
	course, err := cm.repo.GetCourseByID(ctx, req.CourseId)
	if err != nil {
		logger.Error("failed to get course", logger.ErrorType(err))
		return "", fmt.Errorf("failed to get course: %w", err)
	}

	// 2. Validate Status
	if course.Student.Status == 3 {
		return "", fmt.Errorf("学员已退学，无法充值/扣费")
	}
	if course.Status == 3 {
		return "", fmt.Errorf("课程已结课，无法充值/扣费")
	}

	db := dao.GetDB()
	// 3. Recharge (Transaction)
	err = db.Transaction(func(tx *gorm.DB) error {
		txCourseDao := dao.NewStudentCourseDao(tx)
		txCourseRepo := repository.NewCourseRepository(txCourseDao)

		txRechargeDao := dao.NewRechargeOrderDao(tx)
		txRechargeRepo := repository.NewOrderRepository(txRechargeDao)

		if err := txCourseRepo.RechargeCourse(ctx, req.CourseId, req.Hours); err != nil {
			return err
		}

		record := &entity.RechargeOrder{
			StudentCourse: entity.StudentSubject{ID: req.CourseId},
			Hours:         req.Hours,
			Amount:        req.Amount,
			Remark:        req.Remark,
		}
		if err := txRechargeRepo.CreateOrder(ctx, *record); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.Error("failed to recharge course", logger.ErrorType(err))
		return "", fmt.Errorf("failed to recharge course: %w", err)
	}

	return "course recharged", nil
}

func (cm CourseManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "course_manager/create_course", cm.CreateCourse)
	dispatcher.RegisterTyped(d, "course_manager/get_course_list", cm.GetCourseList)
	dispatcher.RegisterTyped(d, "course_manager/toggle_status", cm.ToggleStatus)
	dispatcher.RegisterTyped(d, "course_manager/delete", cm.DeleteCourse)
	dispatcher.RegisterTyped(d, "course_manager/update", cm.UpdateCourse)
	dispatcher.RegisterTyped(d, "course_manager/recharge", cm.RechargeCourse)
}
