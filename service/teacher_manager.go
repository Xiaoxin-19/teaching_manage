package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"teaching_manage/dao"
	"teaching_manage/entity"
	"teaching_manage/pkg"
	"teaching_manage/pkg/wraper"
	"teaching_manage/repository"
	requestx "teaching_manage/service/request"
	responsex "teaching_manage/service/response"
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

func (tm TeacherManager) CreateTeacher(r string) string {
	var teacher entity.Teacher
	err := json.NewDecoder(strings.NewReader(r)).Decode(&teacher)
	if err != nil {
		return wraper.NewBadResponse("unkown json format").ToJSON()
	}

	err = tm.repo.CreateTeacher(context.Background(), teacher)
	if err != nil {
		return wraper.NewBadResponse("internal server error").ToJSON()
	}
	return wraper.NewSuccessResponse("teacher created").ToJSON()
}

func (tm TeacherManager) GetTeacherList(r string) string {
	ctx := context.Background()
	// wails.LogDebugf(tm.Ctx, "service GetTeacherList request: %s", r)
	var req requestx.GetTeacherListRequest
	err := json.Unmarshal([]byte(r), &req)
	if err != nil {
		return wraper.NewBadResponse("unkown json format").ToJSON()
	}

	teachers, total, err := tm.repo.GetTeacherList(ctx, req.Key, req.Offset, req.Limit)
	if err != nil {
		return wraper.NewBadResponse("internal server error").ToJSON()
	}

	var result []entity.Teacher
	for _, t := range teachers {
		result = append(result, entity.Teacher{
			Id:        t.ID,
			CreatedAt: t.CreatedAt.UnixMilli(),
			UpdatedAt: t.UpdatedAt.UnixMilli(),
			Name:      t.Name,
			Gender:    pkg.Gender(t.Gender),
			Phone:     t.Phone,
			Remark:    t.Remark,
		})
	}
	return wraper.NewSuccessResponse(responsex.GetTeacherListResponse{
		Teachers: result,
		Total:    total,
	}).ToJSON()
}

func (tm TeacherManager) DeleteTeacher(r string) string {
	var req requestx.DeleteTeacherRequest
	err := json.Unmarshal([]byte(r), &req)
	if err != nil {
		return wraper.NewBadResponse("unkown json format").ToJSON()
	}

	err = tm.repo.DeleteTeacher(context.Background(), req.Id)
	if err != nil {
		return wraper.NewBadResponse("internal server error").ToJSON()
	}
	return wraper.NewSuccessResponse("teacher deleted").ToJSON()
}

func (tm TeacherManager) UpdateTeacher(r string) string {
	var req requestx.UpdateTeacherRequest
	err := json.Unmarshal([]byte(r), &req)
	if err != nil {
		return wraper.NewBadResponse("unkown json format").ToJSON()
	}
	teacher := entity.Teacher{
		Id:     req.Id,
		Name:   req.Name,
		Phone:  req.Phone,
		Gender: pkg.Gender(req.Gender),
		Remark: req.Remark,
	}
	err = tm.repo.UpdateTeacher(context.Background(), teacher)
	if err != nil {
		return wraper.NewBadResponse("internal server error").ToJSON()
	}
	return wraper.NewSuccessResponse("teacher updated").ToJSON()
}

func (tm TeacherManager) ExportTeacher2Excel() string {
	filepath, err := wails.SaveFileDialog(tm.Ctx, wails.SaveDialogOptions{
		Title:           "选择导出文件位置",
		DefaultFilename: fmt.Sprintf("teachers_%s.xlsx", time.Now().Format("20060102_150405")),
		Filters: []wails.FileFilter{
			{DisplayName: "Excel 文件", Pattern: "*.xlsx"},
		},
	})
	if err != nil {
		return wraper.NewBadResponse("file dialog error").ToJSON()
	}

	if filepath == "" {
		return wraper.NewBadResponse("export cancelled").ToJSON()
	}

	teachers, _, err := tm.repo.GetTeacherList(context.Background(), "", 0, 1000000)
	if err != nil {
		return wraper.NewBadResponse("internal server error").ToJSON()
	}

	// 导出到 Excel using service-level conversion then generic exporter
	if err := tm.exportTeachersToExcel(filepath, teachers); err != nil {
		return wraper.NewBadResponse("failed to save excel: " + err.Error()).ToJSON()
	}

	return wraper.NewSuccessResponse("exported").ToJSON()
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
