package dao

import (
	"context"
	"errors"
	"teaching_manage/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	StudentID uint `gorm:"column:student_id;not null;comment:'学生主键';index;uniqueIndex:idx_stu_teach_date_time"`
	TeacherID uint `gorm:"column:teacher_id;not null;comment:'教师主键';index;uniqueIndex:idx_stu_teach_date_time"`
	// 关联学生与教师，添加外键约束：更新级联、删除受限（避免误删学生或教师导致记录丢失）
	Student Student `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Teacher Teacher `gorm:"foreignKey:TeacherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	// TeachingDateMs 用 Unix 毫秒 (UTC) 表示上课时刻，便于精确排序（int64）
	TeachingDate   time.Time `gorm:"column:teaching_date;type:date;not null;comment:'上课日期';uniqueIndex:idx_stu_teach_date_time"`
	TeachingDateMs int64     `gorm:"column:teaching_date_ms;index;comment:'上课时间 Unix 毫秒(UTC) 整数表示'"`
	StartTime      string    `gorm:"column:start_time;not null;comment:'上课开始时间';uniqueIndex:idx_stu_teach_date_time"`
	EndTime        string    `gorm:"column:end_time;not null;comment:'上课结束时间';uniqueIndex:idx_stu_teach_date_time"`
	Active         bool      `gorm:"column:active;not null;default:false;comment:'是否生效'"`
	Remark         string    `gorm:"column:remark;size:255;comment:'备注字段'"`
}

type RecordDAO interface {
	CreateRecord(ctx context.Context, record Record) error
	GetRecordList(ctx context.Context, stuKey string, teachKey string,
		startDate string, endDate string, offset int, limit int) ([]Record, int64, int64, error)
	ActivateRecord(ctx context.Context, recordID uint) error
	GetRecordByID(ctx context.Context, d uint) (*Record, error)
	DeleteRecordByID(ctx context.Context, id uint) error
	GetAllPendingRecordList(ctx context.Context) ([]Record, error)
}

func NewRecordDao(db *gorm.DB) RecordDAO {
	return &RecordGormDAO{db: db}
}

type RecordGormDAO struct {
	db *gorm.DB
}

func (r *RecordGormDAO) CreateRecord(ctx context.Context, record Record) error {
	convertRecordTimeToUnixMs(&record)
	err := gorm.G[Record](r.db).Create(ctx, &record)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicatedKey
		}
		return err
	}
	return nil
}

func convertRecordTimeToUnixMs(r *Record) {
	// 若未提供 TeachingDateMs，则从 TeachingDate 生成 Unix 毫秒（UTC）
	if r.TeachingDateMs == 0 && !r.TeachingDate.IsZero() {
		r.TeachingDateMs = r.TeachingDate.UTC().UnixMilli()
	}
}

func (r *RecordGormDAO) GetRecordList(ctx context.Context, stuKey string, teachKey string,
	startDate string, endDate string, offset int, limit int) ([]Record, int64, int64, error) {
	var records []Record

	// Unscoped 用于包含记录中学生和老师被软删除的记录
	// 构建查询，关联学生和教师表以进行模糊搜索
	query := r.db.WithContext(ctx).Model(&Record{}).Unscoped().Where("records.deleted_at is null")
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
	err = r.db.WithContext(ctx).Model(&Record{}).Where("active = ?", false).Count(&pendingTotal).Error

	// 应用分页参数并执行查询
	err = query.Offset(offset).Limit(limit).Order("records.teaching_date_ms DESC").Find(&records).Error
	if err != nil {
		return nil, 0, 0, err
	}
	return records, total, pendingTotal, nil
}

func (r *RecordGormDAO) ActivateRecord(ctx context.Context, recordID uint) error {
	_, err := gorm.G[Record](r.db).Where("id = ?", recordID).Update(ctx, "active", true)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecordGormDAO) GetRecordByID(ctx context.Context, d uint) (*Record, error) {
	var record Record
	record, err := gorm.G[Record](r.db).Where("id = ?", d).First(ctx)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RecordGormDAO) DeleteRecordByID(ctx context.Context, id uint) error {
	_, err := gorm.G[Record](r.db).Where("id = ?", id).Delete(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("delete record but id not found", logger.ErrorType(err), logger.UInt("record_id", id))
		return nil
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *RecordGormDAO) GetAllPendingRecordList(ctx context.Context) ([]Record, error) {
	var records []Record
	records, err := gorm.G[Record](r.db).Where("active = ?", false).Find(ctx)
	if err != nil {
		return nil, err
	}
	return records, nil
}
