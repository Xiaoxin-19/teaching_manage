package responsex

type GetTeacherListResponse struct {
	Teachers []TeacherDTO `json:"teachers"`
	Total    int64        `json:"total"`
}

type TeacherDTO struct {
	ID            uint   `json:"id"`
	TeacherNumber string `json:"teacher_number"`
	Name          string `json:"name"`
	Gender        string `json:"gender"`
	Phone         string `json:"phone"`
	Remark        string `json:"remark"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	DeletedAt     int64  `json:"deleted_at"`
}
