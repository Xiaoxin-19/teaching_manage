package model

import (
	"fmt"

	"gorm.io/gorm"
)

// Teacher 教师表
type Teacher struct {
	gorm.Model
	TeacherNumber string `gorm:"column:teacher_number;size:20;uniqueIndex;comment:教师编号(业务主键)" `
	Name          string `gorm:"column:name;not null;comment:教师姓名;index" `
	Gender        string `gorm:"column:gender;comment:教师性别" `
	Phone         string `gorm:"column:phone;comment:电话号码" `
	Remark        string `gorm:"column:remark;comment:备注" `
}

func (Teacher) TableName() string {
	return "teachers"
}

// AfterCreate GORM 钩子：在创建记录后自动生成教师编号
// 逻辑：插入数据 -> 获取生成的 ID -> 生成教师编号 -> 更新教师编号字段
func (t *Teacher) AfterCreate(tx *gorm.DB) (err error) {
	teacherNumber := GenerateTeacherNumber(t.ID)
	err = tx.Model(t).Where("id = ?", t.ID).Update("teacher_number", teacherNumber).Error
	return err
}

// GenerateTeacherNumber 根据教师 ID 生成教师编号:格式 T00000001
func GenerateTeacherNumber(id uint) string {
	return "T" + fmt.Sprintf("%08d", id)
}
