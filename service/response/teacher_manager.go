package responsex

import "teaching_manage/entity"

type GetTeacherListResponse struct {
	Teachers []entity.Teacher `json:"teachers"`
	Total    int64            `json:"total"`
}
