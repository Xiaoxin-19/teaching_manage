package entity

import "time"

type StudentSubject struct {
	ID        uint `json:"id"`
	Student   Student
	Teacher   Teacher
	Subject   Subject
	Balance   int                  `json:"balance"` // 剩余课时数
	Remark    string               `json:"remark"`  // 备注
	Status    StudentSubjectStatus `json:"status"`  // 课程状态
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type StudentSubjectStatus int

const (
	StudentSubjectStatusActive    StudentSubjectStatus = 1 // 正常
	StudentSubjectStatusSuspended StudentSubjectStatus = 2 // 停学
	StudentSubjectStatusCompleted StudentSubjectStatus = 3 // 结课
)

func (ss StudentSubjectStatus) String() string {
	switch ss {
	case StudentSubjectStatusActive:
		return "Active"
	case StudentSubjectStatusSuspended:
		return "Suspended"
	case StudentSubjectStatusCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}
func (ss StudentSubjectStatus) ZhString() string {
	switch ss {
	case StudentSubjectStatusActive:
		return "正常"
	case StudentSubjectStatusSuspended:
		return "停学"
	case StudentSubjectStatusCompleted:
		return "结课"
	default:
		return "未知"
	}
}

func ParseStudentSubjectStatus(status int) StudentSubjectStatus {
	switch status {
	case 1:
		return StudentSubjectStatusActive
	case 2:
		return StudentSubjectStatusSuspended
	case 3:
		return StudentSubjectStatusCompleted
	default:
		return StudentSubjectStatusActive
	}
}
