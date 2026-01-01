package dao

import (
	"context"
	"errors"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg/logger"

	"gorm.io/gorm"
)

type RecordDAO interface {
	CreateRecord(ctx context.Context, record model.Record) error
	GetRecordList(ctx context.Context, stuKey string, teachKey string,
		startDate string, endDate string, offset int, limit int) ([]model.Record, int64, int64, error)
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

func (r *RecordGormDAO) GetRecordList(ctx context.Context, stuKey string, teachKey string,
	startDate string, endDate string, offset int, limit int) ([]model.Record, int64, int64, error) {
	var records []model.Record

	// Unscoped 用于包含记录中学生和老师被软删除的记录
	// 构建查询，关联学生和教师表以进行模糊搜索
	query := r.db.WithContext(ctx).Model(&model.Record{}).Unscoped().Where("records.deleted_at is null")
	query = query.Joins("Teacher").Joins("Student")

	if stuKey != "" {
		query = query.Where("Student.name LIKE ?", "%"+stuKey+"%")
	}
	if teachKey != "" {
		query = query.Where("Teacher.name LIKE ?", "%"+teachKey+"%")
	}

	// 过滤教学日期范围
	if startDate != "" {
		query = query.Where("teaching_date >= ?", startDate)
	}
	if endDate != "" {
		// 增加时间部分以包含当天的记录
		query = query.Where("teaching_date <= ?", endDate+" 23:59:59")
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
		logger.Error("delete record but id not found", logger.ErrorType(err), logger.UInt("record_id", id))
		return nil
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
