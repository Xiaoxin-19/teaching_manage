package dao

import (
	"context"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	StudentID uint   `gorm:"column:student_id;type:int(10) unsigned;comment:'学生主键'"`
	Hours     int    `gorm:"column:hours;type:int(11);not null;comment:'充值课时数'"`
	Comment   string `gorm:"column:comment;type:varchar(50);comment:'备注'"`
	Active    bool   `gorm:"column:active;;not null;default:true;comment:'是否生效'"`
}

type OrderDAO interface {
	CreateOrder(ctx context.Context, order Order) error
	GetOrdersByStudentID(ctx context.Context, studentID uint, offset int, limit int) ([]Order, int64, error)
}

func NewOrderDao(db *gorm.DB) OrderDAO {
	return &OrderGormDAO{db: db}
}

type OrderGormDAO struct {
	db *gorm.DB
}

func (o *OrderGormDAO) CreateOrder(ctx context.Context, order Order) error {
	return gorm.G[Order](o.db).Create(ctx, &order)
}

func (o *OrderGormDAO) GetOrdersByStudentID(ctx context.Context, studentID uint, offset int, limit int) ([]Order, int64, error) {
	var orders []Order
	query := gorm.G[Order](o.db).Where("student_id = ?", studentID)
	total, err := query.Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	orders, err = query.Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return orders, total, err
}
