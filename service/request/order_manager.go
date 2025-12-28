package requestx

type CreateOrderRequest struct {
	StudentID uint   `json:"student_id"`
	Hours     int    `json:"hours"`
	Comment   string `json:"comment"`
}

type GetOrdersByStudentIDRequest struct {
	StudentID uint `json:"student_id" validate:"required"`
	Offset    int  `json:"offset" validate:"gte=0"`
	Limit     int  `json:"limit" validate:"oneof=10 25 50 100 -1"`
}
