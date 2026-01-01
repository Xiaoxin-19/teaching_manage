package entity

type StudentSubject struct {
	ID       uint `json:"id"`
	Student  Student
	Teacher  Teacher
	Subject  Subject
	Balance  int    `json:"balance"`   // 剩余课时数
	TotalBuy int    `json:"total_buy"` // 累计购买课时数
	Remark   string `json:"remark"`    // 备注
}
