package requestx

type UpdateCourseRequest struct {
	ID        uint   `json:"id" validate:"required"`
	TeacherId uint   `json:"teacher_id" validate:"required"`
	Remark    string `json:"remark" validate:"max=255"`
}

type CreateCourseRequest struct {
	StudentId uint   `json:"student_id" validate:"required"`
	SubjectId uint   `json:"subject_id" validate:"required"`
	TeacherId uint   `json:"teacher_id" validate:"required"`
	Remark    string `json:"remark" validate:"max=255"`
}

type ToggleCourseStatusRequest struct {
	CourseId uint `json:"course_id" validate:"required"`
}

type DeleteCourseRequest struct {
	CourseId     uint   `json:"course_id" validate:"required"`
	IsHardDelete bool   `json:"is_hard_delete"`
	Remark       string `json:"remark"`
}

type RechargeCourseRequest struct {
	CourseId uint    `json:"course_id" validate:"required"`
	Hours    int     `json:"hours" validate:"required,ne=0"` // 正数为充值，负数为扣除，不允许为 0
	Amount   float64 `json:"amount"`                         // 实付/退费金额
	Remark   string  `json:"remark"`
}

type GetCourseListRequest struct {
	StudentIds []uint `json:"students" validate:"omitempty"`
	SubjectIds []uint `json:"subjects" validate:"omitempty"`
	TeacherIds []uint `json:"teachers" validate:"omitempty"`
	BalanceMin *int   `json:"balance_min" validate:"omitempty"`
	BalanceMax *int   `json:"balance_max" validate:"omitempty"`
	Offset     int    `json:"offset" validate:"gte=0"`
	Limit      int    `json:"limit" validate:"oneof=10 25 50 100 -1"`
	Statuses   []int  `json:"status" validate:"omitempty,dive,oneof=1 2 3 4 5"`
	Keyword    string `json:"keyword" validate:"max=100"`
}
