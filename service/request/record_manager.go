package requestx

type CreateRecordRequest struct {
	StudentID    uint   `json:"student_id" validate:"required"`
	TeachingDate string `json:"teaching_date" validate:"required,datetime=2006-01-02"`
	StartTime    string `json:"start_time" validate:"required,datetime=15:04"`
	EndTime      string `json:"end_time" validate:"required,datetime=15:04"`
	Remark       string `json:"remark" validate:"max=255"`
}

type GetRecordListRequest struct {
	StudentKey string `json:"student_key" validate:"max=100"`
	TeacherKey string `json:"teacher_key" validate:"max=100"`
	StartDate  string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	Offset     int    `json:"offset" validate:"gte=0"`
	Limit      int    `json:"limit" validate:"oneof=10 25 50 100 -1"`
}

type ActivateRecordRequest struct {
	RecordID uint `json:"record_id" validate:"required"`
}

type BatchActivateRecordsRequest struct {
	RecordIDs []uint `json:"record_ids" validate:"required,dive,gt=0"`
}

type DeleteRecordRequest struct {
	RecordID uint `json:"record_id" validate:"required,gt=0"`
}

type ExportRecordsRequest struct {
	StudentKey string `json:"student_key" validate:"max=100"`
	TeacherKey string `json:"teacher_key" validate:"max=100"`
	StartDate  string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

type ImportRecordsRequest struct {
	Filepath string `json:"filepath" validate:"max=2048"`
}
