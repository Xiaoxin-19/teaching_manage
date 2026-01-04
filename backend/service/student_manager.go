package service

import (
	"context"
	"fmt"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	"teaching_manage/backend/repository"
	requestx "teaching_manage/backend/service/request"
	responsex "teaching_manage/backend/service/response"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type StudentManager struct {
	Ctx        context.Context
	repo       repository.StudentRepository
	repoCourse repository.CourseRepository
}

func NewStudentManager(repo repository.StudentRepository, repoCourse repository.CourseRepository) *StudentManager {
	return &StudentManager{repo: repo, repoCourse: repoCourse}
}

func (sm StudentManager) GetStudentList(ctx context.Context, req *requestx.GetStudentListRequest) (*responsex.GetStudentListResponse, error) {

	studentDs := []entity.Student{}
	var total int64
	var err error
	if req.ShowDeleted {
		// 如果包含已删除的学生，则调用 Unscoped 版本的方法
		studentDs, total, err = sm.repo.ListStudentsWithStatusUnscoped(ctx, req.Keyword, req.Offset, req.Limit, req.StatusLevel, req.StatusTarget)
		if err != nil {
			return nil, err
		}
	} else {
		studentDs, total, err = sm.repo.ListStudentsWithStatus(ctx, req.Keyword, req.Offset, req.Limit, req.StatusLevel, req.StatusTarget)
		if err != nil {
			return nil, err
		}
	}

	studentDTOs := make([]responsex.StudentDTO, 0, len(studentDs))
	for _, s := range studentDs {
		deleteAt := int64(0)
		if !s.DeletedAt.IsZero() {
			deleteAt = s.DeletedAt.UnixMilli()
		}
		studentDTOs = append(studentDTOs, responsex.StudentDTO{
			ID:            s.ID,
			StudentNumber: s.StudentNumber,
			Name:          s.Name,
			Gender:        s.Gender,
			Phone:         s.Phone,
			Status:        s.Status,
			Remark:        s.Remark,
			CreatedAt:     s.CreatedAt.UnixMilli(),
			UpdatedAt:     s.UpdatedAt.UnixMilli(),
			DeletedAt:     deleteAt,
		})
	}

	return &responsex.GetStudentListResponse{
		Students: studentDTOs,
		Total:    total,
	}, nil
}

func (sm StudentManager) CreateStudent(ctx context.Context, req *requestx.CreateStudentRequest) (string, error) {
	logger.Info("Creating one student",
		logger.String("student_name", req.Name),
		logger.String("phone", req.Phone),
		logger.String("gender", req.Gender),
		logger.String("remark", req.Remark),
	)

	err := sm.repo.CreateStudent(ctx, &entity.Student{
		Name:   req.Name,
		Gender: req.Gender,
		Phone:  req.Phone,
		Remark: req.Remark,
	})

	if err != nil {
		logger.Error("failed to create student", logger.ErrorType(err))
		return "", fmt.Errorf("failed to create student: %w", err)
	}
	return "student created", nil
}

func (sm StudentManager) UpdateStudent(ctx context.Context, req *requestx.UpdateStudentRequest) (string, error) {
	logger.Info("Updating one student",
		logger.UInt("id", req.ID),
		logger.String("student_name", req.Name),
		logger.String("phone", req.Phone),
		logger.String("gender", req.Gender),
		logger.String("remark", req.Remark),
		logger.Int("status", req.Status),
	)

	student := &entity.Student{
		ID:     req.ID,
		Name:   req.Name,
		Gender: req.Gender,
		Phone:  req.Phone,
		Remark: req.Remark,
		Status: req.Status,
	}

	// 如果状态改为退学(3)，自动设置退学时间
	if req.Status == 3 {
		student.WithdrawAt = time.Now()
	}

	err := sm.repo.UpdateStudent(ctx, student)

	if err != nil {
		logger.Error("failed to update student", logger.ErrorType(err))
		return "", fmt.Errorf("failed to update student: %w", err)
	}
	return "student updated", nil
}

func (sm StudentManager) DeleteStudent(ctx context.Context, req *requestx.DeleteStudentRequest) (string, error) {
	logger.Info("Deleting one student",
		logger.UInt("id", req.ID),
	)

	// 检查是否有正在上课的课程
	inProgressCourses, err := sm.repoCourse.GetByStudentID(ctx, req.ID)
	if err != nil {
		logger.Error("failed to check student courses before deletion", logger.ErrorType(err))
		return "", fmt.Errorf("failed to check student courses: %w", err)
	}

	countEffective := 0
	for _, c := range inProgressCourses {
		if c.Status != entity.StudentSubjectStatusCompleted {
			countEffective++
		}
	}

	if countEffective > 0 {
		return "", fmt.Errorf("该学生有正在进行的课程，无法删除")
	}
	db := dao.GetDB()

	err = db.Transaction(func(tx *gorm.DB) error {
		txStuRepo := repository.NewStudentRepository(dao.NewStudentDao(tx))
		txCourseRepo := repository.NewCourseRepository(dao.NewStudentCourseDao(tx))

		// 删除学生的选课记录
		logger.Debug("Deleting student courses")
		err = txCourseRepo.DeleteByStudentID(ctx, req.ID)
		if err != nil {
			logger.Error("failed to delete student courses in transaction", logger.ErrorType(err))
			return fmt.Errorf("failed to delete student courses: %w", err)
		}

		logger.Debug("Deleting student record")
		// 删除学生记录
		err = txStuRepo.DeleteStudent(ctx, req.ID)
		if err != nil {
			logger.Error("failed to delete student in transaction", logger.ErrorType(err))
			return fmt.Errorf("failed to delete student: %w", err)
		}

		// 不删除关联的课程记录和充值记录，保留历史数据
		return nil
	})

	if err != nil {
		return "", err
	}

	return "deleted successfully", nil
}

func (sm StudentManager) Export2Excel(ctx context.Context) (string, error) {
	filepath, err := wails.SaveFileDialog(ctx, wails.SaveDialogOptions{
		Title:           "选择导出文件位置",
		DefaultFilename: fmt.Sprintf("students_%s.xlsx", time.Now().Format("20060102_150405")),
		Filters:         []wails.FileFilter{{DisplayName: "Excel 文件", Pattern: "*.xlsx"}},
	})
	if err != nil {
		return "", err
	}
	if filepath == "" {
		return "cancel", nil
	}

	// Get all students with status <= 3 (正常，停课，退出)
	stus, _, err := sm.repo.ListStudentsWithStatus(ctx, "", 0, -1, 3, 0)
	if err != nil {
		return "", err
	}

	// export to excel
	err = sm.exportToExcel(filepath, stus)
	if err != nil {
		return "", fmt.Errorf("导出失败:请检查文件是否被占用或有读写权限")
	}

	return filepath, nil
}

func (sm StudentManager) exportToExcel(path string, students []entity.Student) error {
	headers := []string{"编号", "学生姓名", "性别", "电话号码", "状态", "备注"}
	rows := make([][]string, 0, len(students))
	for _, s := range students {
		rows = append(rows, []string{
			s.StudentNumber,
			s.Name,
			pkg.Gender(s.Gender).ZhString(),
			s.Phone,
			model.StudentStatusToString(s.Status),
			s.Remark,
		})
	}
	return pkg.ExportToExcel(path, headers, rows)
}

func (sm StudentManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "student_manager/get_student_list", sm.GetStudentList)
	dispatcher.RegisterTyped(d, "student_manager/create_student", sm.CreateStudent)
	dispatcher.RegisterTyped(d, "student_manager/update_student", sm.UpdateStudent)
	dispatcher.RegisterTyped(d, "student_manager/delete_student", sm.DeleteStudent)
	dispatcher.RegisterNoReq(d, "student_manager/export_students", sm.Export2Excel)
}
