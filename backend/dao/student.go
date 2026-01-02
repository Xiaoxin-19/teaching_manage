package dao

import (
	"context"
	"errors"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg/logger"

	"gorm.io/gorm"
)

type StudentDao interface {
	CreateStudent(ctx context.Context, stu *model.Student) error
	UpdateStudent(ctx context.Context, stu *model.Student) error
	DeleteStudent(ctx context.Context, id uint) error
	GetStudentByID(ctx context.Context, id uint) (*model.Student, error)
	GetStudentByIdWithDeleted(ctx context.Context, id uint) (*model.Student, error)
	GetStudentListWithStatus(ctx context.Context, key string, offset int, limit int, status model.StudentStatus, targetStatus model.StudentStatus) ([]model.Student, int64, error)
	GetStudentByName(ctx context.Context, name string) (*model.Student, error)
	UpdateStudentHours(ctx context.Context, id uint, hours int) error
	UpdateStudentHoursWithDeleted(ctx context.Context, id uint, hours int) error
}

type StudentGormDao struct {
	db *gorm.DB
}

func NewStudentDao(db *gorm.DB) StudentDao {
	return &StudentGormDao{db: db}
}

func (s StudentGormDao) CreateStudent(ctx context.Context, stu *model.Student) error {
	err := gorm.G[model.Student](s.db).Create(ctx, stu)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	}
	return err
}

func (s StudentGormDao) UpdateStudent(ctx context.Context, stu *model.Student) error {
	_, err := gorm.G[model.Student](s.db).Where("id = ?", stu.ID).Select(
		"name",
		"gender",
		"phone",
		"remark",
		"status",
	).Updates(ctx, *stu)
	if err != nil {
		return err
	}
	return nil
}

func (s StudentGormDao) UpdateStudentHours(ctx context.Context, id uint, diff int) error {
	_, err := gorm.G[model.Student](s.db).Where("id = ?", id).Update(ctx, "hours", gorm.Expr("hours + ?", diff))
	if err != nil {
		return err
	}
	return nil
}

func (s StudentGormDao) UpdateStudentHoursWithDeleted(ctx context.Context, id uint, diff int) error {
	_, err := s.db.Unscoped().WithContext(ctx).Model(&model.Student{}).
		Where("id = ?", id).Update("hours", gorm.Expr("hours + ?", diff)).RowsAffected, s.db.Error
	if err != nil {
		return err
	}
	return nil
}

func (s StudentGormDao) DeleteStudent(ctx context.Context, id uint) error {
	// 使用 ID 初始化 struct，确保 AfterDelete 钩子能获取到 ID
	stu := model.Student{Model: gorm.Model{ID: id}}

	result := s.db.WithContext(ctx).Delete(&stu)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return result.Error
	}
	return nil
}

func (s StudentGormDao) GetStudentByID(ctx context.Context, id uint) (*model.Student, error) {
	stu, err := gorm.G[model.Student](s.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &stu, nil
}

// GetStudentByIdWithDeleted retrieves a student by ID, including those that have been soft-deleted.
func (s StudentGormDao) GetStudentByIdWithDeleted(ctx context.Context, id uint) (*model.Student, error) {
	stu := model.Student{}

	err := s.db.Unscoped().WithContext(ctx).Where("id = ?", id).Preload("Teacher", nil).First(&stu).Error
	if err != nil {
		return nil, err
	}
	return &stu, nil
}

func (s StudentGormDao) GetStudentListWithStatus(ctx context.Context, key string, offset int, limit int,
	statusLevel model.StudentStatus, statusTarget model.StudentStatus) ([]model.Student, int64, error) {

	logger.Info("Fetching student list with status filter:",
		logger.Int("input_status", int(statusLevel)), logger.Int("target_status", int(statusTarget)))
	var students []model.Student
	var total int64

	// student status filtering
	query := gorm.G[model.Student](s.db).Where("status <= ?", statusLevel)
	if statusTarget != 0 {
		query = query.Where("status = ?", statusTarget)
	}

	// keyword filtering
	if key != "" {
		query = query.Where("name LIKE ?", "%"+key+"%").Or("phone LIKE ?", "%"+key+"%").Or("student_number LIKE ?", "%"+key+"%")
	}

	// get total count
	total, err := query.Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}

	// get paginated results
	students, err = query.Offset(offset).Limit(limit).Order("created_at desc").Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return students, total, nil
}

func (s StudentGormDao) GetStudentByName(ctx context.Context, name string) (*model.Student, error) {

	stu, err := gorm.G[model.Student](s.db).Where("name = ?", name).Preload("Teacher", nil).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}
	return &stu, nil
}
