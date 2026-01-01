package requestx

type CreateOrderRequest struct {
	StudentID uint `json:"student_id" validate:"required"`
	// 变更课时数不能为0
	Hours   int    `json:"hours" validate:"required,ne=0"`
	Comment string `json:"comment" validate:"max=255"`
}

type GetOrdersByStudentIDRequest struct {
	StudentID uint `json:"student_id" validate:"required"`
	Offset    int  `json:"offset" validate:"gte=0"`
	Limit     int  `json:"limit" validate:"oneof=10 25 50 100 -1"`
}

type Export2ExcelByIDRequest struct {
	StudentID uint `json:"student_id" validate:"required"`
}
