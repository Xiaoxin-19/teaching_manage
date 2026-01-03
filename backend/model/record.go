package model

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	StudentID uint `gorm:"column:student_id;not null;comment:'学生主键';index;uniqueIndex:idx_stu_teach_date_time"`
	TeacherID uint `gorm:"column:teacher_id;not null;comment:'教师主键';index;uniqueIndex:idx_stu_teach_date_time"`
	SubjectID uint `gorm:"column:subject_id;not null;comment:'科目主键';index;uniqueIndex:idx_stu_teach_date_time"`

	// 关联学生与教师，添加外键约束：更新级联、删除受限（避免误删学生或教师导致记录丢失）
	Student Student `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Teacher Teacher `gorm:"foreignKey:TeacherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Subject Subject `gorm:"foreignKey:SubjectID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	// TeachingDateMs 用 Unix 毫秒 (UTC) 表示上课时刻，便于精确排序（int64）
	TeachingDate   time.Time `gorm:"column:teaching_date;type:date;not null;comment:'上课日期';uniqueIndex:idx_stu_teach_date_time"`
	TeachingDateMs int64     `gorm:"column:teaching_date_ms;index;comment:'上课时间 Unix 毫秒(UTC) 整数表示'"`
	StartTime      string    `gorm:"column:start_time;not null;comment:'上课开始时间';uniqueIndex:idx_stu_teach_date_time"`
	EndTime        string    `gorm:"column:end_time;not null;comment:'上课结束时间';uniqueIndex:idx_stu_teach_date_time"`
	Active         bool      `gorm:"column:active;not null;default:false;comment:'是否生效'"`
	Remark         string    `gorm:"column:remark;size:255;comment:'备注字段'"`
}

func (Record) TableName() string {
	return "records"
}
