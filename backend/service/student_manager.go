package service

import (
	"context"
	"errors"
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
)

type StudentManager struct {
	Ctx  context.Context
	repo repository.StudentRepository
}

func NewStudentManager(repo repository.StudentRepository) *StudentManager {
	return &StudentManager{repo: repo}
}

func (sm StudentManager) GetStudentList(ctx context.Context, req *requestx.GetStudentListRequest) (*responsex.GetStudentListResponse, error) {
	studentDs, total, err := sm.repo.ListStudentsWithStatus(ctx, req.Keyword, req.Offset, req.Limit, 3)
	if err != nil {
		return nil, err
	}

	studentDTOs := make([]responsex.StudentDTO, 0, len(studentDs))
	for _, s := range studentDs {
		if req.Status != 0 && s.Status != req.Status {
			continue
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
	)

	err := sm.repo.UpdateStudent(ctx, &entity.Student{
		ID:     req.ID,
		Name:   req.Name,
		Gender: req.Gender,
		Phone:  req.Phone,
		Remark: req.Remark,
		Status: req.Status,
	})

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
	err := sm.repo.DeleteStudent(ctx, req.ID)
	if err != nil && !errors.Is(err, dao.ErrRecordNotFound) {
		logger.Error("failed to delete student", logger.ErrorType(err))
		return "", fmt.Errorf("failed to delete student: %w", err)
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
	stus, _, err := sm.repo.ListStudentsWithStatus(ctx, "", 0, -1, 3)
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
