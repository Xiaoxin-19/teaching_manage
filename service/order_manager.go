package service

import (
	"context"
	"fmt"
	"teaching_manage/dao"
	"teaching_manage/entity"
	"teaching_manage/pkg"
	"teaching_manage/pkg/dispatcher"
	"teaching_manage/repository"
	requestx "teaching_manage/service/request"
	responsex "teaching_manage/service/response"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type OrderManager struct {
	Ctx     context.Context
	repo    repository.OrderRepository
	stuRepo repository.StudentRepository
}

func NewOrderManager(repo repository.OrderRepository, stuRepo repository.StudentRepository) *OrderManager {
	return &OrderManager{repo: repo, stuRepo: stuRepo}
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
	ordersEntity := make([]responsex.OrderDTO, 0, len(orders))

	for _, o := range orders {
		ordersEntity = append(ordersEntity, responsex.OrderDTO{
			Id:        o.ID,
			CreatedAt: o.CreatedAt.UnixMilli(),
			UpdatedAt: o.UpdatedAt.UnixMilli(),
			Hours:     o.Hours,
			Comment:   o.Comment,
			Active:    o.Active,
			Type:      responsex.OrderDTOTypeToString(o.Hours),
		})
	}
	return responsex.GetOrdersByStudentIDResponse{
		Orders: ordersEntity,
		Total:  total,
	}, nil
}

func (om OrderManager) Export2ExcelByID(ctx context.Context, req *requestx.Export2ExcelByIDRequest) (string, error) {
	filepath, err := wails.SaveFileDialog(ctx, wails.SaveDialogOptions{
		Title:           "选择导出文件位置",
		DefaultFilename: fmt.Sprintf("student_orders_%s.xlsx", time.Now().Format("20060102_150405")),
		Filters:         []wails.FileFilter{{DisplayName: "Excel 文件", Pattern: "*.xlsx"}},
	})
	if err != nil {
		return "", err
	}
	if filepath == "" {
		return "cancelled", nil
	}

	// get orders
	orders, _, err := om.repo.GetOrdersByStudentID(ctx, req.StudentID, 0, -1)
	if err != nil {
		return "", err
	}

	student, err := om.stuRepo.GetStudentByID(ctx, req.StudentID)
	if err != nil {
		return "", err
	}

	// export to excel
	err = om.exportToExcel(filepath, student.Name, orders)
	if err != nil {
		return "", err
	}
	return filepath, nil
}

func (om OrderManager) exportToExcel(path string, stuName string, orders []entity.Order) error {

	headers := []string{"学生姓名", "类别", "课时数", "操作日期", "备注"}
	rows := make([][]string, 0, len(orders))
	for _, order := range orders {
		rows = append(rows, []string{
			stuName,
			responsex.OrderDTOTypeToString(order.Hours),
			fmt.Sprintf("%d", order.Hours),
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			order.Comment,
		})
	}
	return pkg.ExportToExcel(path, headers, rows)
}

func (om OrderManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "order_manager:create_order", om.CreateOrder)
	dispatcher.RegisterTyped(d, "order_manager:get_orders_by_student_id", om.GetOrdersByStudentID)
	dispatcher.RegisterTyped(d, "order_manager:export_orders_by_student_id", om.Export2ExcelByID)
}
