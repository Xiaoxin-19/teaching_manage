package requestx

type GetStudentListRequest struct {
	Keyword string `json:"keyword" validate:"max=100"`
	Offset  int    `json:"offset" validate:"gte=0"`
	Limit   int    `json:"limit" validate:"oneof=10 25 50 100 -1"`
	Status  int    `json:"status" validate:"oneof=0 1 2 3"` // 0表示不筛选状态
}

type CreateStudentRequest struct {
	Name   string `json:"name" validate:"required,max=100"`
	Gender string `json:"gender" validate:"required,oneof=male female"`
	Phone  string `json:"phone" validate:"max=20"`
	Remark string `json:"remark" validate:"max=255"`
}

type UpdateStudentRequest struct {
	ID     uint   `json:"id" validate:"required"`
	Name   string `json:"name" validate:"required,max=100"`
	Gender string `json:"gender" validate:"required,oneof=male female"`
	Phone  string `json:"phone" validate:"max=20"`
	Remark string `json:"remark" validate:"max=255"`
	Status int    `json:"status" validate:"oneof=1 2 3,required"`
}

type DeleteStudentRequest struct {
	ID uint `json:"id" validate:"required"`
}
