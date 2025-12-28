package service

import (
	"context"
	"fmt"
	"teaching_manage/dao"
	"teaching_manage/entity"
	"teaching_manage/pkg/dispatcher"
	"teaching_manage/repository"
	requestx "teaching_manage/service/request"
	responsex "teaching_manage/service/response"

	"gorm.io/gorm"
)

type OrderManager struct {
	repo repository.OrderRepository
}

func NewOrderManager(repo repository.OrderRepository) *OrderManager {
	return &OrderManager{repo: repo}
}

func (om OrderManager) CreateOrder(ctx context.Context, order *requestx.CreateOrderRequest) (string, error) {

	db := dao.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		txO := dao.NewOrderDao(tx)
		txS := dao.NewStudentDao(tx)
		oRepo := repository.NewOrderRepository(txO)
		strRepo := repository.NewStudentRepository(txS)

		// update student hours
		student, err := strRepo.GetStudentByID(ctx, order.StudentID)
		if err != nil {
			return err
		}
		student.Hours += order.Hours
		err = strRepo.UpdateStudentByID(ctx, student)
		if err != nil {
			return err
		}

		eOrder := entity.Order{
			Student: entity.Student{ID: order.StudentID},
			Hours:   order.Hours,
			Comment: order.Comment,
			Active:  true,
		}

		// create order record
		err = oRepo.CreateOrder(ctx, eOrder)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("create order failed")
	}
	return "order created", nil
}

func (om OrderManager) GetOrdersByStudentID(ctx context.Context, req *requestx.GetOrdersByStudentIDRequest) (responsex.GetOrdersByStudentIDResponse, error) {
	orders, total, err := om.repo.GetOrdersByStudentID(ctx, req.StudentID, req.Offset, req.Limit)
	if err != nil {
		return responsex.GetOrdersByStudentIDResponse{}, err
	}
	ordersEntity := make([]entity.Order, 0, len(orders))

	for _, o := range orders {
		ordersEntity = append(ordersEntity, entity.Order{
			Id:        o.ID,
			CreatedAt: o.CreatedAt.UnixMilli(),
			UpdatedAt: o.UpdatedAt.UnixMilli(),
			Student:   entity.Student{ID: o.StudentID},
			Hours:     o.Hours,
			Comment:   o.Comment,
			Active:    o.Active,
		})
	}
	return responsex.GetOrdersByStudentIDResponse{
		Orders: ordersEntity,
		Total:  total,
	}, nil
}

func (om OrderManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "order_manager:create_order", om.CreateOrder)
	dispatcher.RegisterTyped(d, "order_manager:get_orders_by_student_id", om.GetOrdersByStudentID)
}
