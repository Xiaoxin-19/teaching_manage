package requestx

type CreateRecordRequest struct {
	StudentID    uint   `json:"student_id" validate:"required"`
	SubjectID    uint   `json:"subject_id" validate:"required"`
	TeachingDate string `json:"teaching_date" validate:"required,datetime=2006-01-02"`
	StartTime    string `json:"start_time" validate:"required,datetime=15:04"`
	EndTime      string `json:"end_time" validate:"required,datetime=15:04"`
	Remark       string `json:"remark" validate:"max=255"`
}

type GetRecordListRequest struct {
	StudentIDs []uint `json:"student_ids" validate:"omitempty,dive,gt=0"`
	TeacherIDs []uint `json:"teacher_ids" validate:"omitempty,dive,gt=0"`
	SubjectIDs []uint `json:"subject_ids" validate:"omitempty,dive,gt=0"`
	StartDate  string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	Offset     int    `json:"offset" validate:"gte=0"`
	Limit      int    `json:"limit" validate:"oneof=10 25 50 100 -1"`
	Active     *bool  `json:"active" validate:"omitempty"` // nil: all, true: active, false: inactive
}

type ActivateRecordRequest struct {
	ID uint `json:"id" validate:"required"`
}

type BatchActivateRecordsRequest struct {
	IDs []uint `json:"ids" validate:"required,dive,gt=0"`
}

type DeleteRecordRequest struct {
	ID uint `json:"id" validate:"required,gt=0"`
}

type ExportRecordsRequest struct {
	StudentIDs []uint `json:"student_ids" validate:"omitempty,dive,gt=0"`
	TeacherIDs []uint `json:"teacher_ids" validate:"omitempty,dive,gt=0"`
	SubjectIDs []uint `json:"subject_ids" validate:"omitempty,dive,gt=0"`
	StartDate  string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	Active     *bool  `json:"active" validate:"omitempty"`
}

type ImportRecordsRequest struct {
	Filepath string `json:"filepath" validate:"required,max=2048,filepath"`
}
