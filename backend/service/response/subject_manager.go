package responsex

type GetSubjectListResponse struct {
	Subjects []SubjectDTO `json:"subjects"`
	Total    int64        `json:"total"`
}
type SubjectDTO struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	SubjectNumber string `json:"subject_number"`
	StudentCount  int64  `json:"student_count"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}
