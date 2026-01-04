package requestx

type GetSubjectListRequest struct {
	KeyWord     string `json:"keyword"`
	Offset      int    `json:"offset" validate:"gte=0"`
	Limit       int    `json:"limit" validate:"oneof=10 25 50 100 -1"`
	ShowDeleted bool   `json:"show_deleted" validate:""` // 是否包含已删除的科目
}

type CreateSubjectRequest struct {
	Name string `json:"name" validate:"required"`
}
type UpdateSubjectRequest struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
type DeleteSubjectRequest struct {
	ID uint `json:"id" validate:"required"`
}
