package repository

import (
	"context"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order entity.RechargeOrder) error
	GetOrderList(ctx context.Context, studentID uint, subjectIDs []uint, orderType []string, dateStart string, dateEnd string, offset int, limit int) ([]entity.RechargeOrder, int64, error)
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

func (or *OrderRepositoryImpl) GetOrderList(ctx context.Context, studentID uint, subjectIDs []uint, orderType []string, dateStart string, dateEnd string, offset int, limit int) ([]entity.RechargeOrder, int64, error) {
	orders, total, err := or.dao.GetRechargeOrderList(ctx, studentID, subjectIDs, orderType, dateStart, dateEnd, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []entity.RechargeOrder
	for _, order := range orders {
		result = append(result, entity.RechargeOrder{
			ID:          order.ID,
			OrderNumber: order.OrderNumber,
			StudentCourse: entity.StudentSubject{
				ID: order.StudentCourse.ID,
				Student: entity.Student{
					ID:            order.StudentCourse.Student.ID,
					Name:          order.StudentCourse.Student.Name,
					StudentNumber: order.StudentCourse.Student.StudentNumber,
				},
				Subject: entity.Subject{
					ID:   order.StudentCourse.Subject.ID,
					Name: order.StudentCourse.Subject.Name,
				},
				Teacher: entity.Teacher{
					ID:   order.StudentCourse.Teacher.ID,
					Name: order.StudentCourse.Teacher.Name,
				},
			},
			Hours:     order.Hours,
			Amount:    order.Amount,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
			Remark:    order.Remark,
		})
	}
	return result, total, nil
}
