package requestx

type GetStudentListRequest struct {
	Key    string `json:"key"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
