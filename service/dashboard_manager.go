package service

import (
	"context"
	"fmt"
	"teaching_manage/dao"
	"teaching_manage/pkg/dispatcher"
	"teaching_manage/pkg/logger"
	requestx "teaching_manage/service/request"
	responsex "teaching_manage/service/response"
	"time"
)

type DashboardManager struct {
	Ctx context.Context
}

func NewDashboardManager() *DashboardManager {
	return &DashboardManager{}
}

// GetSummaryData 获取顶部核心指标卡数据
func (m *DashboardManager) GetSummaryData(ctx context.Context) (responsex.DashboardSummaryResponse, error) {
	db := dao.GetDB()
	var summary responsex.DashboardSummaryResponse

	// 1. 在读学员总数
	if err := db.Model(&dao.Student{}).Count(&summary.TotalStudents).Error; err != nil {
		logger.Error("Failed to count students", logger.ErrorType(err))
	}

	// 2. 本月新增学员
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if err := db.Model(&dao.Student{}).Where("created_at >= ?", startOfMonth).Count(&summary.NewStudentsThisMonth).Error; err != nil {
		logger.Error("Failed to count new students", logger.ErrorType(err))
	}

	// 3. 剩余总课时
	if err := db.Model(&dao.Student{}).Select("COALESCE(SUM(hours), 0)").Scan(&summary.TotalRemainingHours).Error; err != nil {
		logger.Error("Failed to sum remaining hours", logger.ErrorType(err))
	}

	// 4. 欠费与预警人数
	// 欠费: hours < 0
	if err := db.Model(&dao.Student{}).Where("hours < 0").Count(&summary.TotalArrears).Error; err != nil {
		logger.Error("Failed to count arrears", logger.ErrorType(err))
	}
	// 预警: 0 <= hours < 5
	if err := db.Model(&dao.Student{}).Where("hours >= 0 AND hours < 5").Count(&summary.TotalWarning).Error; err != nil {
		logger.Error("Failed to count warning", logger.ErrorType(err))
	}

	// 5. 本月消课 (Records active=true)
	// 计算本月第一天和下个月第一天
	nextMonth := startOfMonth.AddDate(0, 1, 0)

	// 本月消课数
	var currentMonthCount int64
	if err := db.Model(&dao.Record{}).
		Where("active = 1 AND teaching_date >= ? AND teaching_date < ?", startOfMonth, nextMonth).
		Count(&currentMonthCount).Error; err != nil {
		logger.Error("Failed to count monthly records", logger.ErrorType(err))
	}
	summary.MonthlyHours = currentMonthCount

	// 上月消课数 (用于计算环比)
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)
	var lastMonthCount int64
	if err := db.Model(&dao.Record{}).
		Where("active = 1 AND teaching_date >= ? AND teaching_date < ?", startOfLastMonth, startOfMonth).
		Count(&lastMonthCount).Error; err != nil {
		logger.Error("Failed to count last month records", logger.ErrorType(err))
	}

	// 计算环比
	if lastMonthCount > 0 {
		diff := float64(currentMonthCount - lastMonthCount)
		ratio := (diff / float64(lastMonthCount)) * 100
		prefix := "+"
		if ratio < 0 {
			prefix = "" // 负数自带符号
		}
		summary.MonthOverMonth = fmt.Sprintf("%s%.1f%%", prefix, ratio)
	} else {
		summary.MonthOverMonth = "+0%" // 或 N/A
	}

	return summary, nil
}

// GetFinanceChartData 获取资金/课时流转数据
func (m *DashboardManager) GetFinanceChartData(ctx context.Context, rangeType *requestx.GetFinanceDataRequest) (responsex.FinanceChartDTO, error) {
	db := dao.GetDB()
	var result responsex.FinanceChartDTO

	// 默认实现：按月聚合 (最近 6 个月)
	// TODO 实际项目中需要根据 rangeType 动态生成 SQL 时间范围

	months := []string{}
	now := time.Now()
	// 生成最近6个月的月份标签
	for i := 5; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		months = append(months, t.Format("2006-01"))
	}
	result.XAxis = months

	// SQLite 查询：按月分组统计充值
	// 注意：SQLite 的 strftime 用法
	type MonthlyStat struct {
		Month string
		Total int64
	}

	// 1. 充值数据 (Orders)
	var rechargeStats []MonthlyStat
	// 筛选最近6个月的数据
	startDate := now.AddDate(0, -5, 0).Format("2006-01") + "-01"

	err := db.Model(&dao.Order{}).
		Select("strftime('%Y-%m', created_at) as month, SUM(hours) as total").
		Where("active = 1 AND created_at >= ?", startDate).
		Group("month").
		Order("month").
		Scan(&rechargeStats).Error
	if err != nil {
		return result, err
	}

	// 2. 消课数据 (Records)
	var consumeStats []MonthlyStat
	err = db.Model(&dao.Record{}).
		Select("strftime('%Y-%m', teaching_date) as month, COUNT(id) as total"). // 假设每条记录1课时
		Where("active = 1 AND teaching_date >= ?", startDate).
		Group("month").
		Order("month").
		Scan(&consumeStats).Error
	if err != nil {
		return result, err
	}

	// 3. 数据填充 (Map to Slice)
	rechargeMap := make(map[string]int64)
	for _, s := range rechargeStats {
		rechargeMap[s.Month] = s.Total
	}

	consumeMap := make(map[string]int64)
	for _, s := range consumeStats {
		consumeMap[s.Month] = s.Total
	}

	for _, m := range months {
		rVal := rechargeMap[m]
		cVal := consumeMap[m]
		result.RechargeData = append(result.RechargeData, rVal)
		result.ConsumeData = append(result.ConsumeData, cVal)
		result.NetData = append(result.NetData, rVal-cVal)
	}

	return result, nil
}

// GetTeacherRankData 获取教师排行
func (m *DashboardManager) GetTeacherRankData(ctx context.Context) (responsex.TeacherRankDTO, error) {
	db := dao.GetDB()
	var result responsex.TeacherRankDTO

	type RankStat struct {
		Name  string
		Total int64
	}
	var stats []RankStat

	// 统计本月数据
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	err := db.Table("records").
		Select("teachers.name, COUNT(records.id) as total").
		Joins("JOIN teachers ON records.teacher_id = teachers.id").
		Where("records.active = 1 AND records.deleted_at IS NULL AND records.teaching_date >= ?", startOfMonth).
		Group("teachers.name").
		Order("total DESC").
		Limit(5).
		Scan(&stats).Error

	if err != nil {
		return result, err
	}

	// 教师排行通常是横向柱状图，ECharts 往往需要数据倒序 (数值最大的在最上面)
	// 但 ECharts 的 category y轴默认是从下往上，所以我们把 Top 1 放在数组最后
	for i := len(stats) - 1; i >= 0; i-- {
		result.Names = append(result.Names, stats[i].Name)
		result.Values = append(result.Values, stats[i].Total)
	}

	return result, nil
}

// GetHeatmapData 获取热力图数据
func (m *DashboardManager) GetHeatmapData(ctx context.Context) ([][]int, error) {
	db := dao.GetDB()

	// 我们需要返回一个二维数组 [day_index, hour_index, value]
	// SQLite strftime('%w', teaching_date) 返回 0-6 (0是周日)
	// 我们前端定义的 days 数组通常是 [周一, ..., 周日]，所以周日(0)需要转换成 6，周一(1)转成 0

	type HeatStat struct {
		DayOfWeek int // 0-6 (Sun-Sat)
		Hour      int // 0-23
		Count     int
	}
	var stats []HeatStat

	// 获取最近 3 个月的数据作为样本，样本太少热力图没意义
	startDate := time.Now().AddDate(0, -3, 0)

	// 注意：start_time 是字符串 "HH:MM"，我们需要截取前两位
	err := db.Model(&dao.Record{}).
		Select("CAST(strftime('%w', teaching_date) AS INTEGER) as day_of_week, "+
			"CAST(substr(start_time, 1, 2) AS INTEGER) as hour, "+
			"COUNT(id) as count").
		Where("active = 1 AND teaching_date >= ?", startDate).
		Group("day_of_week, hour").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	// 转换数据格式
	// 前端 ECharts Heatmap data format: [yIndex, xIndex, value]
	// yIndex (Day): 0=Mon, ..., 6=Sun
	// xIndex (Hour): 对应 hours 数组的索引 (08:00, 10:00...)
	// 这里我们需要做一个映射，或者简单点，直接返回原始数据让前端处理
	// 为了方便，我们在后端做简单的映射。假设前端 Hour 轴是 8, 10, 12...20 (每2小时一格)

	var chartData [][]int
	for _, s := range stats {
		// 转换 DayOfWeek: SQLite 0=Sun -> 前端 6; SQLite 1=Mon -> 前端 0
		yIndex := s.DayOfWeek - 1
		if s.DayOfWeek == 0 {
			yIndex = 6
		}

		// 转换 Hour: 映射到最接近的偶数点 (8, 10, 12...)
		// 简单的索引映射： (Hour - 8) / 2。如果 Hour < 8 归为 0
		xIndex := -1
		if s.Hour >= 8 && s.Hour <= 20 {
			xIndex = (s.Hour - 8) / 2
		}

		if xIndex >= 0 {
			chartData = append(chartData, []int{yIndex, xIndex, s.Count})
		}
	}

	return chartData, nil
}

func (m *DashboardManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterNoReq(d, "dashboard_manager:get_summary", m.GetSummaryData)
	dispatcher.RegisterTyped(d, "dashboard_manager:get_finance_chart", m.GetFinanceChartData)
	dispatcher.RegisterNoReq(d, "dashboard_manager:get_teacher_rank", m.GetTeacherRankData)
	dispatcher.RegisterNoReq(d, "dashboard_manager:get_heatmap", m.GetHeatmapData)
}
