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
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var template_excel_headers = []string{"学生姓名", "上课日期", "开始时间", "结束时间", "备注	"}

type RecordManager struct {
	Ctx   context.Context
	repo  repository.RecordRepository
	repoS repository.StudentRepository
}

func NewRecordManager(repo repository.RecordRepository, repoS repository.StudentRepository) *RecordManager {
	return &RecordManager{repo: repo, repoS: repoS}
}

func (rm RecordManager) CreateRecord(ctx context.Context, req *requestx.CreateRecordRequest) (string, error) {
	logger.Info("Creating one record",
		logger.UInt("student_id", req.StudentID),
		logger.String("teaching_date", req.TeachingDate),
		logger.String("start_time", req.StartTime),
		logger.String("end_time", req.EndTime),
	)

	// 解析教学日期
	teachingDate, err := time.Parse("2006-01-02", req.TeachingDate)
	if err != nil {
		return "", err
	}

	// 验证开始时间是否在结束时间之前
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return "", err
	}
	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return "", err
	}
	if !startTime.Before(endTime) {
		return "", fmt.Errorf("start time must be before end time")
	}

	// 检查教师是否存在（假设教师ID通过请求传入，这里暂时使用学生ID作为教师ID示例）
	student, err := rm.repoS.GetStudentByID(ctx, req.StudentID)
	if err != nil {
		logger.Error("failed to get student by ID", logger.UInt("student_id", req.StudentID), logger.ErrorType(err))
		return "cant find student", fmt.Errorf("student not found: %v", err)
	}

	if student.Teacher.ID == 0 || !student.Teacher.DeletedAt.IsZero() {
		logger.Error("associated teacher not found for student", logger.UInt("student_id", req.StudentID),
			logger.UInt("teacher_id", student.Teacher.ID), logger.String("teacher_deleted_at", student.Teacher.DeletedAt.Local().String()))
		return fmt.Sprintf("associated teacher not found for student Name %s, ", student.Name),
			fmt.Errorf("associated teacher not found for student ID %d", req.StudentID)
	}
	teacherID := student.Teacher.ID

	record := &entity.Record{
		Student:      entity.Student{ID: req.StudentID},
		Teacher:      entity.Teacher{ID: teacherID},
		TeachingDate: teachingDate,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Remark:       req.Remark,
	}

	err = rm.repo.CreateRecord(ctx, record)
	if err != nil {
		return "", err
	}
	return "Record created successfully", nil
}

func (rm RecordManager) GetRecordList(ctx context.Context, req *requestx.GetRecordListRequest) (responsex.GetRecordListResponse, error) {
	records, total, pendingTotal, err := rm.repo.GetRecordList(ctx, req.StudentKey, req.TeacherKey,
		req.StartDate, req.EndDate, req.Offset, req.Limit)
	if err != nil {
		logger.Error("failed to get record list", logger.ErrorType(err))
		return responsex.GetRecordListResponse{}, err
	}
	logger.Info("Fetched records",
		logger.Int64("total", total),
		logger.Int("fetched_count", len(records)),
	)

	result := make([]responsex.RecordDTO, len(records))
	for i, rec := range records {
		result[i] = responsex.RecordDTO{
			ID:           rec.ID,
			CreatedAt:    rec.CreatedAt.UnixMilli(),
			UpdatedAt:    rec.UpdatedAt.UnixMilli(),
			StudentID:    rec.Student.ID,
			StudentName:  rec.Student.Name,
			TeacherID:    rec.Teacher.ID,
			TeacherName:  rec.Teacher.Name,
			TeachingDate: rec.TeachingDate.Format("2006-01-02"),
			StartTime:    rec.StartTime,
			EndTime:      rec.EndTime,
			Active:       rec.Active,
			Remark:       rec.Remark,
		}
		if !rec.Student.DeletedAt.IsZero() {
			result[i].StudentName = fmt.Sprintf("%s (已删除)", rec.Student.Name)
		}

		if !rec.Teacher.DeletedAt.IsZero() {
			result[i].TeacherName = fmt.Sprintf("%s (已删除)", rec.Teacher.Name)
		}
	}
	return responsex.GetRecordListResponse{
		Records:      result,
		Total:        total,
		TotalPending: pendingTotal,
	}, nil
}

func (rm *RecordManager) ActivateRecord(ctx context.Context, req *requestx.ActivateRecordRequest) (string, error) {
	logger.Info("Activating record", logger.UInt("record_id", req.RecordID))
	db := dao.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		return activateRecord(ctx, req.RecordID, tx)
	})
	if err != nil {
		logger.Error("failed to activate record", logger.UInt("record_id", req.RecordID), logger.ErrorType(err))
		return "", fmt.Errorf("fail: record activate %s", err.Error())
	}
	return "Record activated successfully", nil
}

func activateRecord(ctx context.Context, recordID uint, db *gorm.DB) error {
	txStudentRepo := repository.NewStudentRepository(dao.NewStudentDao(db))
	txRecordRepo := repository.NewRecordRepository(dao.NewRecordDao(db))

	// find record
	record, err := txRecordRepo.GetRecordByID(ctx, recordID)
	if err != nil {
		logger.Error("failed to get record by ID", logger.UInt("record_id", recordID), logger.ErrorType(err))
		return err
	}

	// update student hours
	student, err := txStudentRepo.GetStudentByIdWithDeleted(ctx, record.Student.ID)
	if err != nil {
		logger.Error("failed to get student by ID", logger.UInt("student_id", record.Student.ID), logger.ErrorType(err))
		return err
	}

	logger.Debug("student info:", logger.String("name", student.Name), logger.UInt("id", student.ID), logger.Int("current_hours", student.Hours))
	err = txStudentRepo.UpdateStudentHoursByIDWithDeleted(ctx, student.ID, -1)
	if err != nil {
		logger.Error("failed to update student hours", logger.UInt("student_id", student.ID), logger.ErrorType(err))
		return err
	}

	// activate record
	err = txRecordRepo.ActivateRecord(ctx, recordID)
	if err != nil {
		logger.Error("failed to activate record", logger.ErrorType(err))
		return err
	}
	return nil
}

func (rm *RecordManager) ActivateAllPendingRecords(ctx context.Context) (string, error) {
	logger.Info("Activating all pending records")
	db := dao.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		txRecordRepo := repository.NewRecordRepository(dao.NewRecordDao(tx))
		// find all pending records
		pendingRecords, err := txRecordRepo.GetAllPendingRecordList(ctx)
		if err != nil {
			logger.Error("failed to get pending records", logger.ErrorType(err))
			return err
		}
		recordIDs := make([]uint, 0, len(pendingRecords))
		for _, record := range pendingRecords {
			if !record.Active {
				recordIDs = append(recordIDs, record.ID)
			}
		}

		logger.Info("Found pending records", logger.Int("count", len(recordIDs)),
			logger.String("record_ids", fmt.Sprintf("%v", recordIDs)))

		for _, record := range pendingRecords {
			if !record.Active {
				recordIDs = append(recordIDs, record.ID)
				err = activateRecord(ctx, record.ID, tx)
				if err != nil {
					logger.Error("failed to activate record", logger.UInt("record_id", record.ID), logger.ErrorType(err))
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		logger.Error("failed to activate all pending records", logger.ErrorType(err))
		return "", fmt.Errorf("fail: activate all pending records %s", err.Error())
	}
	return "All pending records activated successfully", nil
}

func (rm *RecordManager) DeleteRecordByID(ctx context.Context, req *requestx.DeleteRecordRequest) (string, error) {
	logger.Info("Deleting record", logger.UInt("record_id", req.RecordID))
	db := dao.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		txStuRepo := repository.NewStudentRepository(dao.NewStudentDao(tx))
		txRecRepo := repository.NewRecordRepository(dao.NewRecordDao(tx))

		// if record is active ,need to return hours to student
		// find record
		record, err := txRecRepo.GetRecordByID(ctx, req.RecordID)
		if err != nil {
			logger.Error("failed to get record by ID", logger.UInt("record_id", req.RecordID), logger.ErrorType(err))
			return err
		}
		logger.Info("record info:", logger.UInt("id", record.ID), logger.UInt("student_id", record.Student.ID),
			logger.String("active", fmt.Sprintf("%v", record.Active)))
		// return hours to student
		if record.Active {
			err = txStuRepo.UpdateStudentHoursByIDWithDeleted(ctx, record.Student.ID, 1)
			if err != nil {
				logger.Error("failed to return hours to student before deletion", logger.UInt("student_id", record.Student.ID), logger.ErrorType(err))
				return fmt.Errorf("fail: return hours to student before deletion failed: %s", err.Error())
			}
		}

		// delete record
		err = txRecRepo.DeleteRecordByID(ctx, req.RecordID)
		if err != nil {
			logger.Error("failed to activate record before deletion", logger.UInt("record_id", req.RecordID), logger.ErrorType(err))
			return err
		}
		return nil
	})

	if err != nil {
		logger.Error("failed to delete record", logger.UInt("record_id", req.RecordID), logger.ErrorType(err))
		return "", fmt.Errorf("fail: delete record failed: %s", err.Error())
	}

	return "Record deleted successfully", nil
}

func (rm *RecordManager) ExportRecordToExcel(ctx context.Context, req *requestx.ExportRecordsRequest) (string, error) {
	logger.Info("start export records to excel", logger.String("student_key", req.StudentKey),
		logger.String("teacher_key", req.TeacherKey),
		logger.String("start_date", req.StartDate),
		logger.String("end_date", req.EndDate),
	)

	filepath, err := wails.SaveFileDialog(rm.Ctx, wails.SaveDialogOptions{
		Title:           "选择导出文件位置",
		DefaultFilename: fmt.Sprintf("teaching_records_%s.xlsx", time.Now().Format("20060102_150405")),
		Filters:         []wails.FileFilter{{DisplayName: "Excel 文件", Pattern: "*.xlsx"}},
	})

	if err != nil {
		return "", err
	}
	if filepath == "" {
		return "cancel", nil
	}
	records, _, _, err := rm.repo.GetRecordList(ctx, req.StudentKey, req.TeacherKey, req.StartDate, req.EndDate, 0, -1)
	if err != nil {
		logger.Error("failed to get record list for export", logger.ErrorType(err))
		return "", fmt.Errorf("fail: to get record list: %v", err)
	}

	err = exportRecordsToExcelFile(ctx, records, filepath)
	if err != nil {
		logger.Error("failed to export records to excel", logger.ErrorType(err))
		return "", fmt.Errorf("导出失败:请检查文件是否被占用或有读写权限")
	}
	return "Records exported successfully", nil
}

func exportRecordsToExcelFile(ctx context.Context, records []entity.Record, path string) error {

	headers := []string{"学生姓名", "教师姓名", "上课日期", "上课时间", "状态", "备注	"}
	rows := make([][]string, 0, len(records))
	statusToString := map[bool]string{
		true:  "已激活",
		false: "未激活",
	}
	for _, r := range records {
		rows = append(rows, []string{
			r.Student.Name,
			r.Teacher.Name,
			r.TeachingDate.Format("2006-01-02"),
			fmt.Sprintf("%s - %s", r.StartTime, r.EndTime),
			statusToString[r.Active],
			r.Remark,
		})
	}

	return pkg.ExportToExcel(path, headers, rows)
}

func (rm *RecordManager) DownloadImportTemplate(ctx context.Context) (string, error) {
	logger.Info("start download record import template")
	filepath, err := wails.SaveFileDialog(rm.Ctx, wails.SaveDialogOptions{
		Title:           "选择导出模板文件位置",
		DefaultFilename: "record_import_template.xlsx",
		Filters:         []wails.FileFilter{{DisplayName: "Excel 文件", Pattern: "*.xlsx"}},
	})
	if err != nil {
		return "", err
	}

	if filepath == "" {
		return "cancel", nil
	}

	logger.Info("exporting record import template to", logger.String("filepath", filepath))
	rows := [][]string{
		{"张三", "2024-10-01", "10:00", "11:00", "第一次上课"},
	}

	err = pkg.ExportToExcel(filepath, template_excel_headers, rows)
	if err != nil {
		logger.Error("failed to export record import template", logger.ErrorType(err))
		return "", fmt.Errorf("导出失败:请检查文件是否被占用或有读写权限")
	}
	return filepath, nil
}

func (rm *RecordManager) ShowFilePicker(ctx context.Context) (responsex.SelectFileResponse, error) {
	logger.Info("start open file dialog")
	filepath, err := wails.OpenFileDialog(rm.Ctx, wails.OpenDialogOptions{
		Title:   "选择导入文件位置",
		Filters: []wails.FileFilter{{DisplayName: "Excel 文件", Pattern: "*.xlsx"}},
	})
	if err != nil {
		return responsex.SelectFileResponse{}, err
	}

	if filepath == "" {
		return responsex.SelectFileResponse{Filepath: "cancel"}, nil
	}
	return responsex.SelectFileResponse{Filepath: filepath}, nil
}

func (rm *RecordManager) ImportFromExcel(ctx context.Context, req *requestx.ImportRecordsRequest) (responsex.ImportFromExcelResponse, error) {
	logger.Info("start import records from excel", logger.String("filepath", req.Filepath))
	importFilePath := req.Filepath

	f, err := excelize.OpenFile(importFilePath)
	if err != nil {
		logger.Error("failed to open excel file", logger.ErrorType(err))
		return responsex.ImportFromExcelResponse{}, fmt.Errorf("fail:open excel file failed: %w", err)
	}
	defer f.Close()

	// Validate dates in excel
	records, errInfo, err := validateTeachingRecords(ctx, f)
	if err != nil {
		logger.Error("failed to validate excel data", logger.ErrorType(err))
		return responsex.ImportFromExcelResponse{}, err
	}

	// Check for validation errors
	hasErrors := false
	for _, rowErr := range errInfo {
		if len(rowErr) > 0 {
			hasErrors = true
			break
		}
	}

	if hasErrors {
		logger.Error("excel data validation failed")
		return responsex.ImportFromExcelResponse{
			Filepath:   importFilePath,
			ErrorInfos: errInfo,
			TotalRows:  0,
		}, fmt.Errorf("数据验证失败，请检查错误信息")
	}

	db := dao.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		txStudentRepo := repository.NewStudentRepository(dao.NewStudentDao(tx))
		txRecordRepo := repository.NewRecordRepository(dao.NewRecordDao(tx))

		for i, record := range records {
			// Find Student
			student, err := txStudentRepo.GetStudentByName(ctx, record.Student.Name)
			if err != nil {
				logger.Error("failed to get student by name", logger.String("student_name", record.Student.Name), logger.ErrorType(err))
				return fmt.Errorf("第 %d 行: 查询学生 '%s' 失败", i+2, record.Student.Name)
			}

			// Check Teacher
			if student.Teacher.ID == 0 || !student.Teacher.DeletedAt.IsZero() {
				logger.Error("associated teacher not found for student", logger.String("student_name", record.Student.Name),
					logger.UInt("teacher_id", student.Teacher.ID), logger.String("teacher_deleted_at", student.Teacher.DeletedAt.Local().String()))
				return fmt.Errorf("第 %d 行: 学生 '%s' 没有关联有效的教师", i+2, record.Student.Name)
			}

			// Complete the record information
			record.Student.ID = student.ID
			record.Teacher.ID = student.Teacher.ID
			record.Active = false // Imported records are pending by default

			// Create Record
			err = txRecordRepo.CreateRecord(ctx, &record)
			if err != nil {
				if errors.Is(err, dao.ErrDuplicatedKey) {
					logger.Warn("duplicate record found during import", logger.String("student_name", record.Student.Name), logger.ErrorType(err))
					return fmt.Errorf("第 %d 行: 学生 '%s' 的该上课记录已存在，重复导入", i+2, record.Student.Name)
				}
				logger.Error("failed to create record", logger.String("student_name", record.Student.Name), logger.ErrorType(err))
				return fmt.Errorf("第 %d 行: 创建记录失败: %w", i+2, err)
			}
		}
		return nil
	})

	if err != nil {
		logger.Error("failed to import records", logger.ErrorType(err))
		return responsex.ImportFromExcelResponse{
			Filepath:  importFilePath,
			TotalRows: len(records),
		}, err
	}

	return responsex.ImportFromExcelResponse{
		Filepath:   importFilePath,
		TotalRows:  len(records),
		ErrorInfos: [][]string{},
	}, nil
}

func validateTeachingRecords(ctx context.Context, f *excelize.File) ([]entity.Record, [][]string, error) {
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		logger.Error("failed to read rows from excel file", logger.ErrorType(err))
		return nil, nil, fmt.Errorf("无法打开Sheet1")
	}

	if len(rows) < 2 {
		logger.Warn("Excel not contain effective data")
		return nil, nil, fmt.Errorf("Excel 不包含有效数据")
	}

	errInfo := make([][]string, len(rows)-1)

	// validate table header
	for i, header := range template_excel_headers {
		if rows[0][i] != header {
			logger.Error("table header not match template",
				logger.String("input_header", rows[0][i]), logger.String("template_header", header))
			return nil, nil, fmt.Errorf("无效的模板格式, 预期表头包含 %s 但是存在 %s", header, rows[0][i])
		}
	}

	// validate data rows
	records := make([]entity.Record, 0, len(rows)-1)
	for i, row := range rows[1:] {
		if len(row) < 4 {
			errInfo[i] = append(errInfo[i], fmt.Sprintf("第 %d 行: 列数不足4列，信息不足", i+2))
			continue
		}

		// remove blank spaces
		stuName := strings.ReplaceAll(row[0], " ", "")
		teachingDate := strings.ReplaceAll(row[1], " ", "")
		startTime := strings.ReplaceAll(row[2], " ", "")
		endTime := strings.ReplaceAll(row[3], " ", "")
		remark := ""
		if len(row) > 4 {
			remark = strings.TrimSpace(row[4])
		}

		// validate date format

		// student name not empty
		if stuName == "" {
			errInfo[i] = append(errInfo[i], fmt.Sprintf("第 %d 行: 学生姓名不能为空", i+2))
		}

		// teaching date format
		logger.Debug("the teaching date is ", logger.String("date", teachingDate))
		teachingDate = strings.ReplaceAll(teachingDate, "－", "-")
		teachingDate = strings.ReplaceAll(teachingDate, "／", "-")
		teachingDate = strings.ReplaceAll(teachingDate, "/", "-")
		teachingDate = strings.ReplaceAll(teachingDate, ".", "-")

		parsedTeachingDate, err := time.Parse("2006-01-02", teachingDate)
		if err != nil {
			// 尝试解析 MM-DD-YYYY 格式 (例如 12-01-2025)
			parsedTeachingDate, err = time.Parse("01-02-2006", teachingDate)
		}

		if err != nil {
			logger.Error("teaching date parse error", logger.String("teaching_date", teachingDate), logger.ErrorType(err))
			errInfo[i] = append(errInfo[i], fmt.Sprintf("第 %d 行: 上课日期格式错误，需为 YYYY-MM-DD 或 MM-DD-YYYY 格式", i+2))
		}

		// start time format
		startTime = strings.ReplaceAll(startTime, "：", ":")
		startTimeParsed, err := time.Parse("15:04", startTime)
		if err != nil {
			errInfo[i] = append(errInfo[i], fmt.Sprintf("第 %d 行: 开始时间格式错误，需为 HH:MM 格式", i+2))
		}
		startTimeValid := err == nil

		// end time format
		endTime = strings.ReplaceAll(endTime, "：", ":")
		endTimeParsed, err := time.Parse("15:04", endTime)
		if err != nil {
			errInfo[i] = append(errInfo[i], fmt.Sprintf("第 %d 行: 结束时间格式错误，需为 HH:MM 格式", i+2))
		}
		endTimeValid := err == nil

		// validate start time before end time
		if startTimeValid && endTimeValid && !startTimeParsed.Before(endTimeParsed) {
			errInfo[i] = append(errInfo[i], fmt.Sprintf("第 %d 行: 开始时间必须早于结束时间", i+2))
		}

		if len(errInfo[i]) > 0 {
			records = append(records, entity.Record{})
			continue
		}

		records = append(records, entity.Record{
			Student:      entity.Student{Name: stuName},
			TeachingDate: parsedTeachingDate,
			StartTime:    startTime,
			EndTime:      endTime,
			Remark:       remark,
		})
	}

	return records, errInfo, nil
}

func (rm *RecordManager) RegisterRoute(d *dispatcher.Dispatcher) {
	// Register routes related to record management
	dispatcher.RegisterTyped(d, "record_manager:create_record", rm.CreateRecord)
	dispatcher.RegisterTyped(d, "record_manager:get_record_list", rm.GetRecordList)
	dispatcher.RegisterTyped(d, "record_manager:activate_record", rm.ActivateRecord)
	dispatcher.RegisterTyped(d, "record_manager:delete_record_by_id", rm.DeleteRecordByID)
	dispatcher.RegisterNoReq(d, "record_manager:activate_all_pending_records", rm.ActivateAllPendingRecords)
	dispatcher.RegisterTyped(d, "record_manager:export_record_to_excel", rm.ExportRecordToExcel)
	dispatcher.RegisterNoReq(d, "record_manager:download_import_template", rm.DownloadImportTemplate)
	dispatcher.RegisterTyped(d, "record_manager:import_from_excel", rm.ImportFromExcel)
	dispatcher.RegisterNoReq(d, "record_manager:select_import_file", rm.ShowFilePicker)
}
