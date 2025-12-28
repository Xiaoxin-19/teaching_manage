package responsex

import "teaching_manage/entity"

type GetOrdersByStudentIDResponse struct {
	Orders []entity.Order `json:"orders"`
	Total  int64          `json:"total"`
}
