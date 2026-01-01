package entity

import "time"

type Subject struct {
	ID              uint             `json:"id"`
	Name            string           `json:"name"`           // 如：钢琴、声乐
	SubjectNumber   string           `json:"subject_number"` // 如：SUBJ0001
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	StudentSubjects []StudentSubject `json:"student_subjects"` // 选修该科目的学生列表
}
