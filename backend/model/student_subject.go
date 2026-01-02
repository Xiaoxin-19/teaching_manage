package model

import (
	"teaching_manage/backend/pkg/logger"

	"gorm.io/gorm"
)

// StudentSubject 学生-科目关联表，记录学生所学科目及对应的授课老师和课时余额等信息
type StudentSubject struct {
	gorm.Model
	StudentID uint    `gorm:"uniqueIndex:idx_stu_sub;not null;comment:学生ID" `
	Student   Student `gorm:"foreignKey:StudentID" `

	SubjectID uint    `gorm:"uniqueIndex:idx_stu_sub;not null;comment:科目ID" `
	Subject   Subject `gorm:"foreignKey:SubjectID" `

	TeacherID uint    `gorm:"not null;comment:该科目的授课老师" `
	Teacher   Teacher `gorm:"foreignKey:TeacherID" `

	Balance  int    `gorm:"default:0;comment:剩余课时" `
	TotalBuy int    `gorm:"default:0;comment:累计购买课时" `
	Remark   string `gorm:"type:text;comment:备注" `
	Status   int    `gorm:"default:1;comment:状态，1-正常，2-停学 3-结课" `
}

func (StudentSubject) TableName() string {
	return "student_subjects"
}

// AfterDelete GORM 钩子：删除关联的订单记录
func (ss *StudentSubject) AfterDelete(tx *gorm.DB) (err error) {
	logger.Debug("AfterDelete hook triggered for StudentSubject", logger.UInt("student_subject_id", ss.ID))
	if err := tx.Where("student_course_id = ?", ss.ID).Debug().Delete(&RechargeOrder{}).Error; err != nil {
		logger.Error("Failed to cascade delete RechargeOrder records", logger.ErrorType(err), logger.UInt("student_subject_id", ss.ID))
		return err
	}
	logger.Info("Successfully cascaded delete of RechargeOrder records", logger.UInt("student_subject_id", ss.ID))
	return nil
}
