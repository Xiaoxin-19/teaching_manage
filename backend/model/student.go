package model

import (
	"fmt"
	"teaching_manage/backend/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type StudentStatus int

const (
	StudentStatusEnrolled  StudentStatus = 1 // 正常上课
	StudentStatusSuspended StudentStatus = 2 // 暂时停课（如请假、欠费、保留学籍）
	StudentStatusWithdrawn StudentStatus = 3 // 退出
	StudentStatusNone      StudentStatus = 0 // 空状态
)

func IntToStudentStatus(status int) StudentStatus {
	switch status {
	case 1:
		return StudentStatusEnrolled
	case 2:
		return StudentStatusSuspended
	case 3:
		return StudentStatusWithdrawn
	}
	return StudentStatusNone
}

func StudentStatusToString[T StudentStatus | int](status T) string {
	switch StudentStatus(status) {
	case StudentStatusEnrolled:
		return "正常"
	case StudentStatusSuspended:
		return "停课"
	case StudentStatusWithdrawn:
		return "退学"
	default:
		return "未知"
	}
}

// Student 学生表
type Student struct {
	gorm.Model
	StudentNumber string        `json:"student_number" gorm:"uniqueIndex;size:20;comment:学号(业务主键)"` // 新增
	Name          string        `json:"name" gorm:"index;size:100;not null;comment:学生姓名"`
	Gender        string        `json:"gender" gorm:"size:10;comment:学生性别"`
	Phone         string        `json:"phone" gorm:"index;size:20;comment:学生电话号码"`
	Remark        string        `json:"remark" gorm:"type:text;comment:备注"`
	Status        StudentStatus `json:"status" gorm:"size:20;default:1;comment:学生状态，如在读、毕业、休学等"`
}

func (Student) TableName() string {
	return "students"
}

// AfterCreate GORM 钩子：在创建记录后自动生成学号
// 逻辑：插入数据 -> 获取生成的 ID -> 生成学号 -> 更新学号字段
func (s *Student) AfterCreate(tx *gorm.DB) (err error) {
	// 策略：前缀(S) + 年份(2024) + ID补零(00001)
	// 例如 ID=5，年份=2024 -> 学号: S202400005
	currentYear := time.Now().Format("2006")
	s.StudentNumber = fmt.Sprintf("S%s%05d", currentYear, s.ID)

	// 使用 UpdateColumn 仅更新 StudentNumber 字段，避免触发其他钩子或更新时间戳
	return tx.Model(s).UpdateColumn("student_number", s.StudentNumber).Error
}

// AfterDelete GORM 钩子：在删除学生后，级联软删除关联的课程记录
func (s *Student) AfterDelete(tx *gorm.DB) (err error) {
	logger.Info("Cascading soft delete for StudentSubject records",
		logger.UInt("student_id", s.ID), logger.String("student_name", s.Name),
	)
	// 级联软删除关联的 StudentSubject
	// 由于 StudentSubject 包含 gorm.Model，Delete 操作会执行软删除（更新 deleted_at）
	if err := tx.Where("student_id = ?", s.ID).Delete(&StudentSubject{}).Error; err != nil {
		return err
	}
	return nil
}
