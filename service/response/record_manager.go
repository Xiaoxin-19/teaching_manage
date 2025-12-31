package responsex

type GetRecordListResponse struct {
	Records      []RecordDTO `json:"records"`
	Total        int64       `json:"total"`
	TotalPending int64       `json:"total_pending"`
}

type RecordDTO struct {
	ID           uint   `json:"id"`
	StudentID    uint   `json:"student_id"`
	TeacherID    uint   `json:"teacher_id"`
	StudentName  string `json:"student_name"`
	TeacherName  string `json:"teacher_name"`
	TeachingDate string `json:"teaching_date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Active       bool   `json:"active"`
	Remark       string `json:"remark"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type ImportFromExcelResponse struct {
	Filepath   string     `json:"filepath"`
	TotalRows  int        `json:"total_rows"`
	ErrorInfos [][]string `json:"error_infos"`
}

type SelectFileResponse struct {
	Filepath string `json:"filepath"`
}
