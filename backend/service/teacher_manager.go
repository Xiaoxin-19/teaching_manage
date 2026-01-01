package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

type TeacherManager struct {
	Ctx  context.Context
	repo repository.TeacherRepository
}

func NewTeacherManager(repo repository.TeacherRepository) *TeacherManager {
	return &TeacherManager{repo: repo}
}

func (tm TeacherManager) CreateTeacher(ctx context.Context, teacher *requestx.CreateTeacherRequest) (string, error) {
	logger.Info("Creating one teacher",
		logger.String("teacher_name", teacher.Name),
		logger.String("phone", teacher.Phone),
		logger.String("remark", teacher.Remark),
	)

	err := tm.repo.CreateTeacher(ctx, entity.Teacher{
		Name:   strings.TrimSpace(teacher.Name),
		Phone:  strings.TrimSpace(teacher.Phone),
		Gender: pkg.Gender(teacher.Gender),
		Remark: strings.TrimSpace(teacher.Remark),
	})

	if errors.Is(err, dao.ErrDuplicatedKey) {
		logger.Error("duplicate teacher name", logger.String("teacher_name", teacher.Name))
		return "", fmt.Errorf("duplicate: teacher name [%s] already exists", teacher.Name)
	}

	if err != nil {
		logger.Error("failed to create teacher", logger.ErrorType(err))
		return "", fmt.Errorf("failed to create teacher: %w", err)
	}
	return "teacher created", nil
}

func (tm TeacherManager) GetTeacherList(ctx context.Context, req *requestx.GetTeacherListRequest) (responsex.GetTeacherListResponse, error) {

	teachers, total, err := tm.repo.GetTeacherList(ctx, req.Key, req.Offset, req.Limit)
	if err != nil {
		return responsex.GetTeacherListResponse{}, fmt.Errorf("internal server error")
	}

	teacherDtos := make([]responsex.TeacherDTO, len(teachers))
	for i, t := range teachers {
		teacherDtos[i] = responsex.TeacherDTO{
			ID:        t.ID,
			Name:      t.Name,
			Gender:    pkg.Gender(t.Gender).String(),
			Phone:     t.Phone,
			Remark:    t.Remark,
			CreatedAt: t.CreatedAt.UnixMilli(),
			UpdatedAt: t.UpdatedAt.UnixMilli(),
		}
	}
	return responsex.GetTeacherListResponse{
		Teachers: teacherDtos,
		Total:    total,
	}, nil
}

func (tm TeacherManager) DeleteTeacher(ctx context.Context, req *requestx.DeleteTeacherRequest) (string, error) {
	if err := tm.repo.DeleteTeacher(ctx, req.Id); err != nil {
		return "", err
	}
	return "teacher deleted", nil
}

func (tm TeacherManager) UpdateTeacher(ctx context.Context, req *requestx.UpdateTeacherRequest) (string, error) {
	teacher := entity.Teacher{
		ID:     req.Id,
		Name:   req.Name,
		Phone:  req.Phone,
		Gender: pkg.Gender(req.Gender),
		Remark: req.Remark,
	}
	if err := tm.repo.UpdateTeacher(ctx, teacher); err != nil {
		return "", err
	}
	return "teacher updated", nil
}

func (tm TeacherManager) ExportTeacher2Excel(ctx context.Context) (string, error) {
	filepath, err := wails.SaveFileDialog(tm.Ctx, wails.SaveDialogOptions{
		Title:           "选择导出文件位置",
		DefaultFilename: fmt.Sprintf("teachers_%s.xlsx", time.Now().Format("20060102_150405")),
		Filters:         []wails.FileFilter{{DisplayName: "Excel 文件", Pattern: "*.xlsx"}},
	})
	if err != nil {
		return "", err
	}
	if filepath == "" {
		return "cancel", nil
	}

	teachers, _, err := tm.repo.GetTeacherList(ctx, "", 0, 1000000)
	if err != nil {
		return "", err
	}

	if err := tm.exportTeachersToExcel(filepath, teachers); err != nil {
		return "", fmt.Errorf("导出失败:请检查文件是否被占用或有读写权限")
	}
	return filepath, nil
}

// exportTeachersToExcel converts dao.Teacher to generic rows and calls pkg.ExportToExcel.
func (tm TeacherManager) exportTeachersToExcel(path string, teachers []dao.Teacher) error {
	headers := []string{"姓名", "性别", "电话", "备注", "创建时间", "更新时间"}
	rows := make([][]string, 0, len(teachers))
	for _, t := range teachers {
		rows = append(rows, []string{
			t.Name,
			pkg.Gender(t.Gender).ZhString(),
			t.Phone,
			t.Remark,
			t.CreatedAt.Format(time.RFC3339),
			t.UpdatedAt.Format(time.RFC3339),
		})
	}
	return pkg.ExportToExcel(path, headers, rows)
}

func (tm TeacherManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "teacher_manager:create_teacher", tm.CreateTeacher)
	dispatcher.RegisterTyped(d, "teacher_manager:get_teacher_list", tm.GetTeacherList)
	dispatcher.RegisterTyped(d, "teacher_manager:delete_teacher", tm.DeleteTeacher)
	dispatcher.RegisterTyped(d, "teacher_manager:update_teacher", tm.UpdateTeacher)
	dispatcher.RegisterNoReq(d, "teacher_manager:export_teacher_to_excel", tm.ExportTeacher2Excel)
}
