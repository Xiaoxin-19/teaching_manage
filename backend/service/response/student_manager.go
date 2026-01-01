package responsex

type GetStudentListResponse struct {
	Students []StudentDTO `json:"students"`
	Total    int64        `json:"total"`
}

type StudentDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Hours       int    `json:"hours"`
	Phone       string `json:"phone"`
	TeacherID   uint   `json:"teacher_id"`
	Remark      string `json:"remark"`
	TeacherName string `json:"teacher_name"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	DeletedAt   int64  `json:"deleted_at"`
}
