package repository

import (
	"context"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order entity.RechargeOrder) error
}

type OrderRepositoryImpl struct {
	dao dao.RechargeOrderDao
}

func NewOrderRepository(dao dao.RechargeOrderDao) OrderRepository {
	return &OrderRepositoryImpl{dao: dao}
}

func (or *OrderRepositoryImpl) CreateOrder(ctx context.Context, order entity.RechargeOrder) error {
	o := model.RechargeOrder{
		StudentCourseID: order.StudentCourse.ID,
		Hours:           order.Hours,
		Amount:          order.Amount,
		Remark:          order.Remark,
	}

	return or.dao.CreateRechargeRecord(ctx, &o)
}
