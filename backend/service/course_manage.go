package service

import (
	"context"
	"errors"
	"fmt"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/pkg"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	"teaching_manage/backend/repository"
	requestx "teaching_manage/backend/service/request"
	responsex "teaching_manage/backend/service/response"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type CourseManager struct {
	Ctx       context.Context
	uw        repository.UnitWork
	repo      repository.CourseRepository
	stuRepo   repository.StudentRepository
	orderRepo repository.OrderRepository
}

func NewCourseManager(uw repository.UnitWork, repo repository.CourseRepository, stuRepo repository.StudentRepository, orderRepo repository.OrderRepository) *CourseManager {
	return &CourseManager{uw: uw, repo: repo, stuRepo: stuRepo, orderRepo: orderRepo}
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
	var err error
	if req.IsHardDelete {
		err = cm.repo.DeleteCourse(ctx, req.CourseId)
	} else {
		err = cm.repo.UpdateStatus(ctx, req.CourseId, 3, req.Remark)
	}
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

	// 4. 更新教师信息
	err = cm.repo.UpdateCourseTeacher(ctx, req.ID, req.TeacherId, req.Remark)
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

	err = cm.uw.Execute(ctx, func(txCtx context.Context) error {
		// 处理金额符号：充值时为正数，退费时为负数
		adjustedAmount := req.Amount
		if req.Hours < 0 {
			// 退费时，amount 变为负数，表示从总价值中扣除
			adjustedAmount = -req.Amount
		}

		// 更新课程的剩余课时和总价值
		if err := cm.repo.RechargeCourse(txCtx, req.CourseId, req.Hours, adjustedAmount); err != nil {
			return err
		}

		record := &entity.RechargeOrder{
			StudentCourse: entity.StudentSubject{ID: req.CourseId},
			Hours:         req.Hours,
			Amount:        req.Amount,
			Remark:        req.Remark,
		}
		// 创建充值订单
		if err := cm.orderRepo.CreateOrder(txCtx, *record); err != nil {
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

func (cm CourseManager) Export2Excel(ctx context.Context, req *requestx.ExportCourseListRequest) (string, error) {
	filepath, err := wails.SaveFileDialog(ctx, wails.SaveDialogOptions{
		Title:           "选择导出文件位置",
		DefaultFilename: fmt.Sprintf("courses_%s.xlsx", time.Now().Format("20060102_150405")),
		Filters:         []wails.FileFilter{{DisplayName: "Excel 文件", Pattern: "*.xlsx"}},
	})
	if err != nil {
		return "", err
	}
	if filepath == "" {
		return "cancel", nil
	}

	// Override Limit to -1 to get all records
	req.Limit = -1
	req.Offset = 0

	courses, _, err := cm.repo.GetCourseList(ctx,
		req.StudentIds, req.SubjectIds, req.TeacherIds, req.BalanceMin, req.BalanceMax, req.Statuses, req.Keyword, req.Offset, req.Limit)
	if err != nil {
		return "", fmt.Errorf("failed to get course list for export: %w", err)
	}

	// Export
	err = cm.exportToExcel(filepath, courses)
	if err != nil {
		return "", fmt.Errorf("导出失败:请检查文件是否被占用或有读写权限")
	}

	return filepath, nil
}

func (cm CourseManager) exportToExcel(path string, courses []entity.StudentSubject) error {
	headers := []string{"学员姓名", "学号", "科目", "授课老师", "剩余课时", "状态", "备注"}
	rows := make([][]string, 0, len(courses))
	for _, c := range courses {
		rows = append(rows, []string{
			c.Student.Name,
			c.Student.StudentNumber,
			c.Subject.Name,
			c.Teacher.Name,
			fmt.Sprintf("%d", c.Balance),
			c.Status.ZhString(),
			c.Remark,
		})
	}
	return pkg.ExportToExcel(path, headers, rows)
}

func (cm CourseManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "course_manager/create_course", cm.CreateCourse)
	dispatcher.RegisterTyped(d, "course_manager/get_course_list", cm.GetCourseList)
	dispatcher.RegisterTyped(d, "course_manager/toggle_status", cm.ToggleStatus)
	dispatcher.RegisterTyped(d, "course_manager/delete", cm.DeleteCourse)
	dispatcher.RegisterTyped(d, "course_manager/update", cm.UpdateCourse)
	dispatcher.RegisterTyped(d, "course_manager/recharge", cm.RechargeCourse)
	dispatcher.RegisterTyped(d, "course_manager/export_courses", cm.Export2Excel)
}
