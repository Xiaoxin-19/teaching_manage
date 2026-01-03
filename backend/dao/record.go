package dao

import (
	"context"
	"errors"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type RecordDAO interface {
	CreateRecord(ctx context.Context, record model.Record) error
	GetRecordList(ctx context.Context, stuIDs []uint, teacherIDs []uint, subjectIDs []uint,
		startDate string, endDate string, offset int, limit int, active *bool) ([]model.Record, int64, int64, error)
	ActivateRecord(ctx context.Context, recordID uint) error
	GetRecordByID(ctx context.Context, d uint) (*model.Record, error)
	DeleteRecordByID(ctx context.Context, id uint) error
	GetAllPendingRecordList(ctx context.Context) ([]model.Record, error)
}

func NewRecordDao(db *gorm.DB) RecordDAO {
	return &RecordGormDAO{db: db}
}

type RecordGormDAO struct {
	db *gorm.DB
}

func (r *RecordGormDAO) CreateRecord(ctx context.Context, record model.Record) error {
	convertRecordTimeToUnixMs(&record)
	err := gorm.G[model.Record](r.db).Create(ctx, &record)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicatedKey
		}
		return err
	}
	return nil
}

func convertRecordTimeToUnixMs(r *model.Record) {
	// 若未提供 TeachingDateMs，则从 TeachingDate 生成 Unix 毫秒（UTC）
	if r.TeachingDateMs == 0 && !r.TeachingDate.IsZero() {
		r.TeachingDateMs = r.TeachingDate.UTC().UnixMilli()
	}
}

func (r *RecordGormDAO) GetRecordList(ctx context.Context, stuIDs []uint, teacherIDs []uint, subjectIDs []uint,
	startDate string, endDate string, offset int, limit int, active *bool) ([]model.Record, int64, int64, error) {
	var records []model.Record

	// 构建查询，关联学生和教师表以进行模糊搜索
	query := r.db.WithContext(ctx).Model(&model.Record{}).Where("records.deleted_at is null")
	query = query.Joins("Teacher").Joins("Student").Joins("Subject")

	// 过滤学生ID
	if len(stuIDs) > 0 {
		query = query.Where("records.student_id IN ?", stuIDs)
	}
	// 过滤教师ID
	if len(teacherIDs) > 0 {
		query = query.Where("records.teacher_id IN ?", teacherIDs)
	}
	// 过滤科目ID
	if len(subjectIDs) > 0 {
		query = query.Where("records.subject_id IN ?", subjectIDs)
	}

	// 过滤教学日期范围
	if startDate != "" {
		query = query.Where("teaching_date >= ?", startDate)
	}
	if endDate != "" {
		// 增加时间部分以包含当天的记录
		query = query.Where("teaching_date <= ?", endDate+" 23:59:59")
	}

	// 过滤激活状态
	if active != nil {
		query = query.Where("records.active = ?", *active)
	}

	// 获取总记录数
	total := int64(0)
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, 0, err
	}

	pendingTotal := int64(0)
	err = r.db.WithContext(ctx).Model(&model.Record{}).Where("active = ?", false).Count(&pendingTotal).Error

	// 应用分页参数并执行查询
	err = query.Offset(offset).Limit(limit).Order("records.teaching_date_ms DESC").Find(&records).Error
	if err != nil {
		return nil, 0, 0, err
	}
	return records, total, pendingTotal, nil
}

func (r *RecordGormDAO) ActivateRecord(ctx context.Context, recordID uint) error {
	_, err := gorm.G[model.Record](r.db).Where("id = ?", recordID).Update(ctx, "active", true)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecordGormDAO) GetRecordByID(ctx context.Context, d uint) (*model.Record, error) {
	var record model.Record
	record, err := gorm.G[model.Record](r.db).Where("id = ?", d).First(ctx)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RecordGormDAO) DeleteRecordByID(ctx context.Context, id uint) error {
	_, err := gorm.G[model.Record](r.db).Where("id = ?", id).Delete(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *RecordGormDAO) GetAllPendingRecordList(ctx context.Context) ([]model.Record, error) {
	var records []model.Record
	records, err := gorm.G[model.Record](r.db).Where("active = ?", false).Find(ctx)
	if err != nil {
		return nil, err
	}
	return records, nil
}
