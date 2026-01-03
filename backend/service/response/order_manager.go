package responsex

type GetOrdersListResponse struct {
	Orders []OrderDTO `json:"orders"`
	Total  int64      `json:"total"`
}

type OrderDTO struct {
	ID        uint       `json:"id"`
	OrderNo   string     `json:"order_number"`
	Student   StudentDTO `json:"student"`
	Subject   SubjectDTO `json:"subject"`
	Teacher   TeacherDTO `json:"teacher"`
	Hours     int        `json:"hours"`
	Amount    float64    `json:"amount"`
	Type      string     `json:"type" validate:"oneof=increase decrease"`
	Remark    string     `json:"remark"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
}

func OrderDTOTypeToString(hours int) string {
	if hours >= 0 {
		return "increase"
	}
	return "decrease"
}

func OrderDTOTypeToZhString(hours int) string {
	if hours >= 0 {
		return "充值"
	}
	return "扣费"
}
