package repository

import (
	"context"
	"teaching_manage/dao"
	"teaching_manage/entity"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order entity.Order) error
	GetOrdersByStudentID(ctx context.Context, studentID uint, offset int, limit int) ([]entity.Order, int64, error)
}

type OrderRepositoryImpl struct {
	dao dao.OrderDAO
	// Implement the order repository
}

func NewOrderRepository(dao dao.OrderDAO) OrderRepository {
	return &OrderRepositoryImpl{dao: dao}
}

func (or *OrderRepositoryImpl) CreateOrder(ctx context.Context, order entity.Order) error {
	o := dao.Order{
		StudentID: order.Student.ID,
		Hours:     order.Hours,
		Comment:   order.Comment,
		Active:    order.Active,
	}

	return (or.dao).CreateOrder(ctx, o)
}

func (or *OrderRepositoryImpl) GetOrdersByStudentID(ctx context.Context, studentID uint, offset int, limit int) ([]entity.Order, int64, error) {
	orders, total, err := (or.dao).GetOrdersByStudentID(ctx, studentID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	var result []entity.Order
	for _, o := range orders {
		result = append(result, entity.Order{
			Id:        o.ID,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
			Student:   entity.Student{ID: o.StudentID},
			Hours:     o.Hours,
			Comment:   o.Comment,
			Active:    o.Active,
		})
	}

	return result, total, nil
}
