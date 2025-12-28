package requestx

type GetStudentListRequest struct {
	Key    string `json:"key" validate:"max=100"`
	Offset int    `json:"offset" validate:"gte=0"`
	Limit  int    `json:"limit" validate:"oneof=10 25 50 100 -1"`
}

type CreateStudentRequest struct {
	Name      string `json:"name" validate:"required,max=100"`
	Gender    string `json:"gender" validate:"required,oneof=male female"`
	Hours     int    `json:"hours" validate:"gte=0"`
	Phone     string `json:"phone" validate:"max=20"`
	TeacherID uint   `json:"teacher_id" validate:"required"`
}

type UpdateStudentRequest struct {
	ID        uint   `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required,max=100"`
	Gender    string `json:"gender" validate:"required,oneof=male female"`
	Phone     string `json:"phone" validate:"max=20"`
	TeacherID uint   `json:"teacher_id" validate:"required"`
}

type DeleteStudentRequest struct {
	ID uint `json:"id" validate:"required"`
}
