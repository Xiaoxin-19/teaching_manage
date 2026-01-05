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
	"gorm.io/gorm"
)

type TeacherManager struct {
	Ctx        context.Context
	repo       repository.TeacherRepository
	repoCourse repository.CourseRepository
}

func NewTeacherManager(repo repository.TeacherRepository, repoCourse repository.CourseRepository) *TeacherManager {
	return &TeacherManager{repo: repo, repoCourse: repoCourse}
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
		Status: int(pkg.TeacherStatusActive),
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
	var teachers []entity.Teacher
	var total int64
	var err error
	if req.ShowDeleted {
		teachers, total, err = tm.repo.GetTeacherListUnscoped(ctx, req.KeyWord, req.Status, req.Offset, req.Limit)
		if err != nil {
			return responsex.GetTeacherListResponse{}, fmt.Errorf("internal server error")
		}
	} else {
		teachers, total, err = tm.repo.GetTeacherList(ctx, req.KeyWord, req.Status, req.Offset, req.Limit)
		if err != nil {
			return responsex.GetTeacherListResponse{}, fmt.Errorf("internal server error")
		}
	}

	teacherDtos := make([]responsex.TeacherDTO, len(teachers))
	for i, t := range teachers {
		deletedAt := int64(0)
		if !t.DeletedAt.IsZero() {
			deletedAt = t.DeletedAt.UnixMilli()
		}
		teacherDtos[i] = responsex.TeacherDTO{
			ID:            t.ID,
			Name:          t.Name,
			TeacherNumber: t.TeacherNumber,
			Gender:        pkg.Gender(t.Gender).String(),
			Phone:         t.Phone,
			Status:        t.Status,
			Remark:        t.Remark,
			CreatedAt:     t.CreatedAt.UnixMilli(),
			UpdatedAt:     t.UpdatedAt.UnixMilli(),
			DeletedAt:     deletedAt,
		}
	}
	return responsex.GetTeacherListResponse{
		Teachers: teacherDtos,
		Total:    total,
	}, nil
}

func (tm TeacherManager) DeleteTeacher(ctx context.Context, req *requestx.DeleteTeacherRequest) (string, error) {
	logger.Warn("deleting teacher",
		logger.UInt("teacher_id", req.Id),
	)
	// 检查是否有正在进行中的课程

	courses, err := tm.repoCourse.GetByTeacherID(ctx, req.Id)
	if err != nil {
		logger.Error("failed to get courses by teacher id", logger.ErrorType(err))
		return "", err
	}
	countEffectiveCourses := 0
	for _, c := range courses {
		if c.Status != 3 { // not "已取消"
			countEffectiveCourses++
		}
	}
	if countEffectiveCourses > 0 {
		logger.Warn("teacher has effective courses, cannot delete", logger.UInt("teacher_id", req.Id), logger.Int("effective_course_count", countEffectiveCourses))
		return "", fmt.Errorf("该教师有正在进行中的课程，无法删除")
	}

	db := dao.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewTeacherRepository(dao.NewTeacherDao(dao.GetDBTarget(tx)))
		txCourseRepo := repository.NewCourseRepository(dao.NewStudentCourseDao(dao.GetDBTarget(tx)))

		// 删除教师的选课记录
		err = txCourseRepo.DeleteByTeacherID(ctx, req.Id)
		if err != nil {
			logger.Error("failed to delete teacher courses in transaction", logger.ErrorType(err))
			return fmt.Errorf("failed to delete teacher courses: %w", err)
		}

		// 删除教师记录
		err = txRepo.DeleteTeacher(ctx, req.Id)
		if err != nil {
			logger.Error("failed to delete teacher in transaction", logger.ErrorType(err))
			return fmt.Errorf("failed to delete teacher: %w", err)
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	return "teacher deleted", nil
}

func (tm TeacherManager) UpdateTeacher(ctx context.Context, req *requestx.UpdateTeacherRequest) (string, error) {
	logger.Info("updating teacher",
		logger.UInt("teacher_id", req.Id),
		logger.String("teacher_name", req.Name),
		logger.String("phone", req.Phone),
		logger.String("remark", req.Remark),
	)

	// 如果需要修改老师离职，则需要检测是否有正在进行中的课程
	if req.Status == int(pkg.TeacherStatusResigned) {
		course, err := tm.repoCourse.GetByTeacherID(ctx, req.Id)

		if err != nil {
			logger.Error("failed to get courses by teacher id", logger.ErrorType(err))
			return "", err
		}
		effectiveCourseCount := 0
		for _, c := range course {
			if c.Status != 3 { // not "已取消"
				effectiveCourseCount++
			}
		}
		if effectiveCourseCount > 0 {
			logger.Warn("teacher has effective courses, cannot resign", logger.UInt("teacher_id", req.Id), logger.Int("effective_course_count", effectiveCourseCount))
			return "", fmt.Errorf("该教师有正在进行中的课程，无法设置为离职状态")
		}
	}

	teacher := entity.Teacher{
		ID:     req.Id,
		Name:   req.Name,
		Phone:  req.Phone,
		Gender: pkg.Gender(req.Gender),
		Status: req.Status,
		Remark: req.Remark,
	}

	if err := tm.repo.UpdateTeacher(ctx, teacher); err != nil {
		logger.Error("failed to update teacher", logger.ErrorType(err))
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

	teachers, _, err := tm.repo.GetTeacherList(ctx, "", 0, 0, 1000000)
	if err != nil {
		return "", err
	}

	if err := tm.exportTeachersToExcel(filepath, teachers); err != nil {
		return "", fmt.Errorf("导出失败:请检查文件是否被占用或有读写权限")
	}
	return filepath, nil
}

// exportTeachersToExcel converts dao.Teacher to generic rows and calls pkg.ExportToExcel.
func (tm TeacherManager) exportTeachersToExcel(path string, teachers []entity.Teacher) error {
	headers := []string{"编号", "姓名", "性别", "电话", "状态", "备注", "创建时间", "更新时间"}
	rows := make([][]string, 0, len(teachers))
	for _, t := range teachers {
		rows = append(rows, []string{
			t.TeacherNumber,
			t.Name,
			pkg.Gender(t.Gender).ZhString(),
			t.Phone,
			pkg.TeacherStatus(t.Status).String(),
			t.Remark,
			t.CreatedAt.Format(time.RFC3339),
			t.UpdatedAt.Format(time.RFC3339),
		})
	}
	return pkg.ExportToExcel(path, headers, rows)
}

func (tm TeacherManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "teacher_manager/create_teacher", tm.CreateTeacher)
	dispatcher.RegisterTyped(d, "teacher_manager/get_teacher_list", tm.GetTeacherList)
	dispatcher.RegisterTyped(d, "teacher_manager/delete_teacher", tm.DeleteTeacher)
	dispatcher.RegisterTyped(d, "teacher_manager/update_teacher", tm.UpdateTeacher)
	dispatcher.RegisterNoReq(d, "teacher_manager/export_teacher_to_excel", tm.ExportTeacher2Excel)
}
