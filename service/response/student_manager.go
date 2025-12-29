package responsex

import "teaching_manage/entity"

type GetStudentListResponse struct {
	Students []entity.Student `json:"students"`
	Total    int64            `json:"total"`
}
