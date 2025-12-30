package entity

type Student struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Hours       int     `json:"hours"`
	Phone       string  `json:"phone"`
	TeacherID   uint    `json:"teacher_id"`
	Remark      string  `json:"remark"`
	TeacherName string  `json:"teacher_name"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
	Teacher     Teacher `json:"teacher"`
}
