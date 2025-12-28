package requestx

type GetStudentListRequest struct {
	Key    string `json:"key" validate:"max=100"`
	Offset int    `json:"offset" validate:"gte=0"`
	Limit  int    `json:"limit" validate:"oneof=10 25 50 100 -1"`
}
