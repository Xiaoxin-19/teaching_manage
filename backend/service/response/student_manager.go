package responsex

type GetStudentListResponse struct {
	Students []StudentDTO `json:"students"`
	Total    int64        `json:"total"`
}

type StudentDTO struct {
	ID            uint   `json:"id"`
	StudentNumber string `json:"student_number"`
	Name          string `json:"name"`
	Gender        string `json:"gender"`
	Phone         string `json:"phone"`
	Remark        string `json:"remark"`
	Status        int    `json:"status"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}
