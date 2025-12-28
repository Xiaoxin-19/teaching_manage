package entity

type Order struct {
	Student
	Id        uint   `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Hours     int    `json:"hours"`
	Comment   string `json:"comment"`
	Active    bool   `json:"active"`
}
