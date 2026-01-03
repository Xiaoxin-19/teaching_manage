package requestx

type GetOrderListRequest struct {
	StudentID  uint     `json:"student_id" validate:"omitempty,gt=0"`
	SubjectIDs []uint   `json:"subject_ids" validate:"omitempty,dive,gt=0"`
	Type       []string `json:"type" validate:"omitempty,dive,oneof=increase decrease"`
	DateStart  string   `json:"date_start" validate:"omitempty,datetime=2006-01-02"`
	DateEnd    string   `json:"date_end" validate:"omitempty,datetime=2006-01-02"`
	Offset     int      `json:"offset" validate:"gte=0"`
	Limit      int      `json:"limit" validate:"oneof=10 25 50 100 -1"`
}

type ExportOrdersRequest = GetOrderListRequest
