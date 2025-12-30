package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	StudentID uint `gorm:"column:student_id;not null;comment:'学生主键';index"`
	TeacherID uint `gorm:"column:teacher_id;not null;comment:'教师主键';index"`
	// 关联学生与教师，添加外键约束：更新级联、删除受限（避免误删学生或教师导致记录丢失）
	Student Student `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Teacher Teacher `gorm:"foreignKey:TeacherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	// TeachingDateMs 用 Unix 毫秒 (UTC) 表示上课时刻，便于精确排序（int64）
	TeachingDate   time.Time `gorm:"column:teaching_date;type:date;not null;comment:'上课日期'"`
	TeachingDateMs int64     `gorm:"column:teaching_date_ms;index;comment:'上课时间 Unix 毫秒(UTC) 整数表示'"`
	StartTime      string    `gorm:"column:start_time;not null;type:datetime;comment:'上课开始时间'"`
	EndTime        string    `gorm:"column:end_time;not null;type:datetime;comment:'上课结束时间'"`
	Active         bool      `gorm:"column:active;not null;default:false;comment:'是否生效'"`
	Comment        string    `gorm:"column:comment;size:255;comment:'备注字段'"`
}

type RecordDAO interface {
	CreateRecord(ctx context.Context, record Record) error
}

func NewRecordDao(db *gorm.DB) RecordDAO {
	return &RecordGormDAO{db: db}
}

type RecordGormDAO struct {
	db *gorm.DB
}

func (r *RecordGormDAO) CreateRecord(ctx context.Context, record Record) error {
	// 若未提供 TeachingDateMs，则从 TeachingDate 生成 Unix 毫秒（UTC）
	if record.TeachingDateMs == 0 && !record.TeachingDate.IsZero() {
		record.TeachingDateMs = record.TeachingDate.UTC().UnixMilli()
	}
	return gorm.G[Record](r.db).Create(ctx, &record)
}

func ConvertRecordTimeToUnixMs(r *Record) {
	// 若未提供 TeachingDateMs，则从 TeachingDate 生成 Unix 毫秒（UTC）
	if r.TeachingDateMs == 0 && !r.TeachingDate.IsZero() {
		r.TeachingDateMs = r.TeachingDate.UTC().UnixMilli()
	}
}
