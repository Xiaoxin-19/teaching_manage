package dao

import (
	"context"

	"gorm.io/gorm"
)

type StudentDao interface {
	CreateStudent(ctx context.Context, stu *Student) error
	UpdateStudent(ctx context.Context, stu *Student) error
	DeleteStudent(ctx context.Context, id uint) error
	GetStudentByID(ctx context.Context, id uint) (*Student, error)
	GetStudentList(ctx context.Context, key string, offset int, limit int) ([]Student, error)
}

type StudentGormDao struct {
	db *gorm.DB
}

func NewStudentDao(db *gorm.DB) StudentDao {
	return &StudentGormDao{db: db}
}

type Student struct {
	gorm.Model
	Name      string `gorm:"column:name;not null;comment:学生姓名" json:"name"`
	Gender    string `gorm:"column:gender;comment:学生性别" json:"gender"`
	Hours     int    `gorm:"column:hours;default:0;comment:课时数" json:"hours"`
	Phone     string `gorm:"column:phone;comment:学生电话号码" json:"phone"`
	TeacherID uint   `gorm:"column:teacher_id;not null;comment:授课老师" json:"teacher_id"`
}

func (s StudentGormDao) CreateStudent(ctx context.Context, stu *Student) error {
	return gorm.G[Student](s.db).Create(ctx, stu)
}

func (s StudentGormDao) UpdateStudent(ctx context.Context, stu *Student) error {
	_, err := gorm.G[Student](s.db).Where("id = ?", stu.ID).Updates(ctx, *stu)
	if err != nil {
		return err
	}
	return nil
}

func (s StudentGormDao) DeleteStudent(ctx context.Context, id uint) error {
	_, err := gorm.G[Student](s.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s StudentGormDao) GetStudentByID(ctx context.Context, id uint) (*Student, error) {
	stu, err := gorm.G[Student](s.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &stu, nil
}

func (s StudentGormDao) GetStudentList(ctx context.Context, key string, offset int, limit int) ([]Student, error) {
	var students []Student
	query := gorm.G[Student](s.db).Offset(offset).Limit(limit)
	if key != "" {
		query = query.Where("name LIKE ?", "%"+key+"%")
	}
	students, err := query.Find(ctx)
	if err != nil {
		return nil, err
	}
	return students, nil
}
