package dao

import (
	"context"
	"errors"
	"fmt"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StudentCourseDao interface {
	CreateStudentCourse(ctx context.Context, sc *model.StudentSubject) error
	GetStudentCourse(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error)
	GetStudentCourseWithDeleted(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error)
	GetByStudentIDAndSubjectID(ctx context.Context, studentID uint, subjectID uint) (*model.StudentSubject, error)
	GetByStudentID(ctx context.Context, d uint) ([]model.StudentSubject, error)
	GetByTeacherID(ctx context.Context, id uint) ([]model.StudentSubject, error)
	GetBySubjectID(ctx context.Context, d uint) ([]model.StudentSubject, error)

	UpdateBalance(ctx context.Context, id uint, delta int) error
	Recharge(ctx context.Context, id uint, hours int, amount float64) error
	RestoreStudentCourse(ctx context.Context, id uint) error
	GetStudentCourseList(ctx context.Context, students []uint, subjects []uint, teachers []uint, min *int, max *int, statuses []int, keyword string, offset int, limit int) ([]model.StudentSubject, int64, error)
	UpdateStatus(ctx context.Context, id uint, status int, remark string) error
	GetByID(ctx context.Context, id uint) (*model.StudentSubject, error)
	UpdateStudentCourseInfo(ctx context.Context, id uint, teacherID uint, remark string) error
	Delete(ctx context.Context, id uint) error
	DeleteByStudentID(ctx context.Context, stuID uint) error
	DeleteByTeacherID(ctx context.Context, teacherID uint) error
	DeleteBySubjectID(ctx context.Context, subjectID uint) error
}

type StudentCourseGormDao struct {
	getDB DBGetter
}

func NewStudentCourseDao(getDB DBGetter) StudentCourseDao {
	return &StudentCourseGormDao{getDB: getDB}
}

func (d *StudentCourseGormDao) CreateStudentCourse(ctx context.Context, sc *model.StudentSubject) error {
	db := GetDBFromCtx(ctx, d.getDB())
	err := gorm.G[model.StudentSubject](db).Create(ctx, sc)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	}
	return err
}

func (d *StudentCourseGormDao) GetStudentCourseList(ctx context.Context,
	students []uint, subjects []uint, teachers []uint, min *int, max *int, statuses []int, keyword string, offset int, limit int) ([]model.StudentSubject, int64, error) {
	var scs []model.StudentSubject
	db := GetDBFromCtx(ctx, d.getDB())
	query := gorm.G[model.StudentSubject](db).Preload("Teacher", nil).Preload("Subject", nil)

	// Join Student table to allow filtering on its status
	query = query.Joins(clause.JoinTarget{Association: "Student"}, func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
		return nil
	})

	if len(statuses) > 0 {
		query = query.Where(`(CASE 
			WHEN "Student"."status" = 3 THEN 5 
			WHEN "Student"."status" = 2 THEN 4 
			ELSE "student_subjects"."status" 
		END) IN ?`, statuses)
	}

	// apply filters
	if len(students) > 0 {
		query = query.Where("student_subjects.student_id IN ?", students)
	}
	if len(subjects) > 0 {
		query = query.Where("student_subjects.subject_id IN ?", subjects)
	}
	if len(teachers) > 0 {
		query = query.Where("student_subjects.teacher_id IN ?", teachers)
	}
	if min != nil {
		query = query.Where("student_subjects.balance >= ?", *min)
	}
	if max != nil {
		query = query.Where("student_subjects.balance <= ?", *max)
	}

	// count total records
	total, err := query.Count(ctx, "*")
	if err != nil {
		logger.Error("failed to count student courses",
			logger.ErrorType(err),
		)
		return nil, 0, err
	}
	// apply pagination
	scs, err = query.Offset(offset).Limit(limit).Find(ctx)
	if err != nil {
		logger.Error("failed to get student courses",
			logger.ErrorType(err),
		)
		return nil, 0, err
	}
	return scs, total, nil
}

func (d *StudentCourseGormDao) GetStudentCourse(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error) {
	var sc model.StudentSubject
	db := GetDBFromCtx(ctx, d.getDB())
	err := db.WithContext(ctx).
		Where("student_id = ? AND subject_id = ?", studentID, subjectID).
		First(&sc).Error
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (d *StudentCourseGormDao) GetStudentCourseWithDeleted(ctx context.Context, studentID, subjectID uint) (*model.StudentSubject, error) {
	var sc model.StudentSubject
	db := GetDBFromCtx(ctx, d.getDB())
	err := db.Unscoped().WithContext(ctx).
		Where("student_id = ? AND subject_id = ?", studentID, subjectID).
		First(&sc).Error
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (d *StudentCourseGormDao) UpdateBalance(ctx context.Context, id uint, delta int) error {
	db := GetDBFromCtx(ctx, d.getDB())
	return db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", delta)).Error
}

func (d *StudentCourseGormDao) Recharge(ctx context.Context, id uint, hours int, amount float64) error {
	db := GetDBFromCtx(ctx, d.getDB())
	updates := map[string]interface{}{
		"balance": gorm.Expr("balance + ?", hours),
	}
	if hours > 0 {
		updates["total_buy"] = gorm.Expr("total_buy + ?", hours)
	}

	// Update AvgPrice: (OldBalance * OldAvgPrice + Amount) / (OldBalance + Hours)
	// 该方法仅用于【充值】和【退费】，日常消课通过 UpdateBalance 实现：
	//   - 充值：hours > 0, amount > 0（正数）
	//           新avg_price = (balance * avg_price + amount) / (balance + hours)
	//   - 退费：hours < 0, amount < 0（负数，调用者应在服务层处理符号转换）
	//           新avg_price = (balance * avg_price + amount) / (balance + hours)
	//           例：30节@50元/节，退费10节，退款500元 → amount=-500
	//               = (30*50 - 500) / (30-10) = 1000/20 = 50元/节 ✓
	// 当余额耗尽(balance + hours = 0)时，保留原有的 avg_price 供后续充值计算使用
	// Note: In SQLite, division by zero returns NULL.
	updates["avg_price"] = gorm.Expr("COALESCE((balance * avg_price + ?) / NULLIF(balance + ?, 0), avg_price)", amount, hours)

	return db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (d *StudentCourseGormDao) RestoreStudentCourse(ctx context.Context, id uint) error {
	db := GetDBFromCtx(ctx, d.getDB())
	return db.Unscoped().WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

func (d *StudentCourseGormDao) UpdateStatus(ctx context.Context, id uint, status int, remark string) error {
	db := GetDBFromCtx(ctx, d.getDB())
	return db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": status,
			"remark": remark,
		}).Error
}

func (d *StudentCourseGormDao) GetByID(ctx context.Context, id uint) (*model.StudentSubject, error) {
	var sc model.StudentSubject
	db := GetDBFromCtx(ctx, d.getDB())
	err := db.WithContext(ctx).Preload("Student").Preload("Subject").Preload("Teacher").First(&sc, id).Error
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (d *StudentCourseGormDao) UpdateStudentCourseInfo(ctx context.Context, id uint, teacherID uint, remark string) error {
	db := GetDBFromCtx(ctx, d.getDB())
	return db.WithContext(ctx).Model(&model.StudentSubject{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"teacher_id": teacherID,
			"remark":     remark,
		}).Error
}

func (d *StudentCourseGormDao) Delete(ctx context.Context, id uint) error {
	ss := model.StudentSubject{Model: gorm.Model{ID: id}}
	db := GetDBFromCtx(ctx, d.getDB())
	return db.WithContext(ctx).Delete(&ss).Error
}

func (s StudentCourseGormDao) GetByStudentIDAndSubjectID(ctx context.Context, studentID uint, subjectID uint) (*model.StudentSubject, error) {
	var sc model.StudentSubject
	db := GetDBFromCtx(ctx, s.getDB())
	err := db.WithContext(ctx).
		Where("student_id = ? AND subject_id = ?", studentID, subjectID).
		Preload("Student").Preload("Teacher").Preload("Subject").
		First(&sc).Error

	logger.Debug("teacher info:", logger.String("teacher_name", fmt.Sprintf("%v", sc.Teacher)))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (s StudentCourseGormDao) GetByStudentID(ctx context.Context, d uint) ([]model.StudentSubject, error) {
	var sc []model.StudentSubject
	db := GetDBFromCtx(ctx, s.getDB())
	err := db.WithContext(ctx).
		Where("student_id = ?", d).
		Preload("Student").
		Find(&sc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (d *StudentCourseGormDao) DeleteByStudentID(ctx context.Context, stuID uint) error {
	ss := model.StudentSubject{}
	db := GetDBFromCtx(ctx, d.getDB())
	return db.WithContext(ctx).Where("student_id = ?", stuID).Delete(&ss).Error
}

func (d *StudentCourseGormDao) DeleteByTeacherID(ctx context.Context, teacherID uint) error {
	ss := model.StudentSubject{}
	db := GetDBFromCtx(ctx, d.getDB())
	err := db.WithContext(ctx).Where("teacher_id = ?", teacherID).Delete(&ss).Error
	return err
}

func (d *StudentCourseGormDao) GetByTeacherID(ctx context.Context, id uint) ([]model.StudentSubject, error) {
	var sc []model.StudentSubject
	db := GetDBFromCtx(ctx, d.getDB())
	err := db.WithContext(ctx).Where("teacher_id = ?", id).Find(&sc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (d *StudentCourseGormDao) GetBySubjectID(ctx context.Context, subjectID uint) ([]model.StudentSubject, error) {
	var sc []model.StudentSubject
	db := GetDBFromCtx(ctx, d.getDB())
	err := db.WithContext(ctx).Where("subject_id = ?", subjectID).Find(&sc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (d *StudentCourseGormDao) DeleteBySubjectID(ctx context.Context, subjectID uint) error {
	ss := model.StudentSubject{}
	db := GetDBFromCtx(ctx, d.getDB())
	err := db.WithContext(ctx).Where("subject_id = ?", subjectID).Delete(&ss).Error
	return err
}
