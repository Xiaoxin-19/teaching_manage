package responsex

type CourseDTO struct {
	ID        uint       `json:"id"`
	Student   StudentDTO `json:"student"`
	Subject   SubjectDTO `json:"subject"`
	Teacher   TeacherDTO `json:"teacher"`
	Balance   int        `json:"balance"` // 剩余课时数
	Remark    string     `json:"remark"`  // 备注
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
	Status    int        `json:"status"` // 课程状态
}

type GetCourseListResponse struct {
	Courses []CourseDTO `json:"courses"`
	Total   int64       `json:"total"`
}
