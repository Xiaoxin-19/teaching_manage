package requestx

type CreateTeacherRequest struct {
	Name   string `json:"name" validate:"required"`
	Phone  string `json:"phone"`
	Gender string `json:"gender" validate:"required,oneof=male female"`
	Remark string `json:"remark"`
}

type UpdateTeacherRequest struct {
	Id     uint   `json:"id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Phone  string `json:"phone"`
	Gender string `json:"gender" validate:"required,oneof=male female"`
	Remark string `json:"remark"`
}

type GetTeacherListRequest struct {
	Key    string `json:"key"`
	Offset int    `json:"offset" validate:"gte=0"`
	Limit  int    `json:"limit" validate:"oneof=10 25 50 100 -1"`
}

type DeleteTeacherRequest struct {
	Id uint `json:"id" validate:"required"`
}
