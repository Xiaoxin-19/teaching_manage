package service

import (
	"context"
	"fmt"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/pkg"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	"teaching_manage/backend/repository"
	requestx "teaching_manage/backend/service/request"
	responsex "teaching_manage/backend/service/response"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type OrderManager struct {
	Ctx     context.Context
	repo    repository.OrderRepository
	stuRepo repository.StudentRepository
}

func NewOrderManager(repo repository.OrderRepository, stuRepo repository.StudentRepository) *OrderManager {
	return &OrderManager{repo: repo, stuRepo: stuRepo}
}

func (om OrderManager) GetOrderList(ctx context.Context, req *requestx.GetOrderListRequest) (responsex.GetOrdersListResponse, error) {
	logger.Info(
		"GetOrderList request",
		logger.String("student_id", fmt.Sprint(req.StudentID)),
		logger.String("subject_ids", fmt.Sprint(req.SubjectIDs)),
		logger.String("type", fmt.Sprint(req.Type)),
		logger.String("date_start", req.DateStart),
		logger.String("date_end", req.DateEnd),
		logger.String("offset", fmt.Sprint(req.Offset)),
		logger.String("limit", fmt.Sprint(req.Limit)),
	)
	var orders []entity.RechargeOrder
	var total int64
	var err error

	orders, total, err = om.repo.GetOrderList(ctx, req.StudentID, req.SubjectIDs, req.Type, req.DateStart, req.DateEnd, req.Offset, req.Limit)
	if err != nil {
		logger.Error("GetOrderList error", logger.ErrorType(err))
		return responsex.GetOrdersListResponse{}, err
	}

	var orderDTOs []responsex.OrderDTO
	for _, order := range orders {
		orderDTOs = append(orderDTOs, responsex.OrderDTO{
			ID:        order.ID,
			OrderNo:   order.OrderNumber,
			Student:   responsex.StudentDTO{ID: order.StudentCourse.Student.ID, Name: order.StudentCourse.Student.Name, StudentNumber: order.StudentCourse.Student.StudentNumber},
			Subject:   responsex.SubjectDTO{ID: order.StudentCourse.Subject.ID, Name: order.StudentCourse.Subject.Name},
			Teacher:   responsex.TeacherDTO{ID: order.StudentCourse.Teacher.ID, Name: order.StudentCourse.Teacher.Name},
			Hours:     order.Hours,
			Amount:    order.Amount,
			Type:      responsex.OrderDTOTypeToString(order.Hours),
			Remark:    order.Remark,
			CreatedAt: order.CreatedAt.UnixMilli(),
			UpdatedAt: order.UpdatedAt.UnixMilli(),
		})
	}

	return responsex.GetOrdersListResponse{
		Orders: orderDTOs,
		Total:  total,
	}, nil
}

func (om OrderManager) ExportOrdersToExcel(ctx context.Context, req *requestx.ExportOrdersRequest) (string, error) {
	filepath, err := wails.SaveFileDialog(ctx, wails.SaveDialogOptions{
		Title:           "选择导出文件位置",
		DefaultFilename: fmt.Sprintf("student_orders_%s.xlsx", time.Now().Format("20060102_150405")),
		Filters:         []wails.FileFilter{{DisplayName: "Excel 文件", Pattern: "*.xlsx"}},
	})
	if err != nil {
		return "", err
	}
	if filepath == "" {
		return "cancel", nil
	}

	// get orders
	orders, _, err := om.repo.GetOrderList(ctx, req.StudentID, req.SubjectIDs, req.Type, req.DateStart, req.DateEnd, 0, -1)
	if err != nil {
		return "", err
	}

	// export to excel
	err = om.exportToExcel(filepath, orders)
	if err != nil {
		return "", fmt.Errorf("导出失败:请检查文件是否被占用或有读写权限")
	}
	return filepath, nil
}

func (om OrderManager) exportToExcel(path string, orders []entity.RechargeOrder) error {

	headers := []string{"订单号", "学生姓名", "科目", "类别", "课时数", "实际金额", "操作日期", "备注"}
	rows := make([][]string, 0, len(orders))
	for _, order := range orders {
		rows = append(rows, []string{
			order.OrderNumber,
			fmt.Sprintf("%s(%s)", order.StudentCourse.Student.Name, order.StudentCourse.Student.StudentNumber),
			order.StudentCourse.Subject.Name,
			responsex.OrderDTOTypeToZhString(order.Hours),
			fmt.Sprintf("%d", order.Hours),
			fmt.Sprintf("%.2f", order.Amount),
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			order.Remark,
		})
	}
	return pkg.ExportToExcel(path, headers, rows)
}

func (om OrderManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterTyped(d, "order_manager/get_order_list", om.GetOrderList)
	dispatcher.RegisterTyped(d, "order_manager/export_orders_to_excel", om.ExportOrdersToExcel)
}
