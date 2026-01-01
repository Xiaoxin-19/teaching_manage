package model

import (
	"fmt"

	"gorm.io/gorm"
)

// Subject 科目表
type Subject struct {
	gorm.Model
	SubjectNumber string `gorm:"unique;size:20;not null;comment:科目编号(业务主键)"` // 如：SUBJ0001
	Name          string `gorm:"unique;size:50;not null;comment:科目名称"`       // 如：钢琴、声乐

	// 一对多关联：一个科目可以被多个学生选修
	StudentSubjects []StudentSubject `gorm:"foreignKey:SubjectID;references:ID"`
}

func (Subject) TableName() string {
	return "subjects"
}

// AfterCreate GORM 钩子：在创建记录后自动生成科目编号
// 逻辑：插入数据 -> 获取生成的 ID -> 生成科目编号 -> 更新科目编号字段
func (s *Subject) AfterCreate(tx *gorm.DB) (err error) {
	subjectNumber := fmt.Sprintf("SUBJ%04d", s.ID)
	err = tx.Model(s).Where("id = ?", s.ID).Update("subject_number", subjectNumber).Error
	return err
}
