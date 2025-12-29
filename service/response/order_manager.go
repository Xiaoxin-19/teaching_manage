package responsex

type GetOrdersByStudentIDResponse struct {
	Orders []OrderDTO `json:"orders"`
	Total  int64      `json:"total"`
}

type OrderDTO struct {
	Id        uint   `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Hours     int    `json:"hours"`
	Comment   string `json:"comment"`
	Active    bool   `json:"active"`
	Type      string `json:"type" validate:"oneof=increase decrease"`
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
