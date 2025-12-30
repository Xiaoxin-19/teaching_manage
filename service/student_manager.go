package service

import (
	"context"
	"errors"
	"fmt"
	"teaching_manage/dao"
	"teaching_manage/entity"
	"teaching_manage/pkg"
	"teaching_manage/pkg/dispatcher"
	"teaching_manage/pkg/logger"
	"teaching_manage/repository"
	requestx "teaching_manage/service/request"
	responsex "teaching_manage/service/response"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type StudentManager struct {
	Ctx   context.Context
	repo  repository.StudentRepository
	repoT repository.TeacherRepository
}

func NewStudentManager(repo repository.StudentRepository, repoT repository.TeacherRepository) *StudentManager {
	return &StudentManager{repo: repo, repoT: repoT}
}

func (sm StudentManager) GetStudentList(ctx context.Context, req *requestx.GetStudentListRequest) (*responsex.GetStudentListResponse, error) {
	studentDs, total, err := sm.repo.GetStudentList(ctx, req.Key, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	studentDTOs := make([]responsex.StudentDTO, len(studentDs))
	for i, s := range studentDs {
		studentDTOs[i] = responsex.StudentDTO{
			ID:          s.ID,
			Name:        s.Name,
			Gender:      s.Gender,
			Hours:       s.Hours,
			Phone:       s.Phone,
			TeacherID:   s.TeacherID,
			Remark:      s.Remark,
			TeacherName: s.TeacherName,
			CreatedAt:   s.CreatedAt.UnixMilli(),
			UpdatedAt:   s.UpdatedAt.UnixMilli(),
		}
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
		logger.Int("hours", req.Hours),
		logger.UInt("teacher_id", req.TeacherID),
		logger.String("remark", req.Remark),
	)

	err := sm.repo.CreateStudent(ctx, &entity.Student{
		Name:      req.Name,
		Gender:    req.Gender,
		Hours:     req.Hours,
		Phone:     req.Phone,
		TeacherID: req.TeacherID,
		Remark:    req.Remark,
	})

	if errors.Is(err, dao.ErrDuplicatedKey) {
		logger.Error("duplicate student name", logger.String("student_name", req.Name))
		return "", fmt.Errorf("duplicate : student name [%s] already exists", req.Name)
	}

	if err != nil {
		logger.Error("failed to create student", logger.ErrorType(err))
		return "", fmt.Errorf("failed to create student: %w", err)
	}
	return "student created", nil
}

func (sm StudentManager) UpdateStudent(ctx context.Context, req *requestx.UpdateStudentRequest) (string, error) {
	return "updated successfully", sm.repo.UpdateStudentByID(ctx, &entity.Student{
		ID:        req.ID,
		Name:      req.Name,
		Gender:    req.Gender,
		Phone:     req.Phone,
		TeacherID: req.TeacherID,
		Remark:    req.Remark,
	})
}

func (sm StudentManager) DeleteStudent(ctx context.Context, req *requestx.DeleteStudentRequest) (string, error) {
	return "deleted successfully", sm.repo.DeleteStudentByID(ctx, req.ID)
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

	stus, _, err := sm.repo.GetStudentList(ctx, "", 0, -1)
	if err != nil {
		return "", err
	}

	// Get teachers for mapping
	teachers, _, err := sm.repoT.GetTeacherList(ctx, "", 0, -1)
	if err != nil {
		return "", err
	}
	teacherMap := make(map[uint]string)
	for _, teacher := range teachers {
		teacherMap[teacher.ID] = teacher.Name
	}

	// Map teacher names
	for i, stu := range stus {
		if name, ok := teacherMap[stu.TeacherID]; ok {
			stus[i].TeacherName = name
		}
	}

	// export to excel
	err = sm.exportToExcel(filepath, stus)
	if err != nil {
		return "", fmt.Errorf("导出失败:请检查文件是否被占用或有读写权限")
	}

	return filepath, nil
}

func (sm StudentManager) exportToExcel(path string, students []entity.Student) error {
	headers := []string{"学生姓名", "性别", "课时数", "电话号码", "授课老师", "备注"}
	rows := make([][]string, 0, len(students))
	for _, s := range students {
		rows = append(rows, []string{
			s.Name,
			pkg.Gender(s.Gender).ZhString(),
			fmt.Sprintf("%d", s.Hours),
			s.Phone,
			s.TeacherName,
			s.Remark,
		})
	}
	return pkg.ExportToExcel(path, headers, rows)
}

func (sm StudentManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "student_manager:get_student_list", sm.GetStudentList)
	dispatcher.RegisterTyped(d, "student_manager:create_student", sm.CreateStudent)
	dispatcher.RegisterTyped(d, "student_manager:update_student", sm.UpdateStudent)
	dispatcher.RegisterTyped(d, "student_manager:delete_student", sm.DeleteStudent)
	dispatcher.RegisterNoReq(d, "student_manager:export_students", sm.Export2Excel)
}
