package service

import (
	"context"
	"fmt"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	requestx "teaching_manage/backend/service/request"
	responsex "teaching_manage/backend/service/response"
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

	// 1. 在读学员总数 (去重: 学生状态为正常(1)且至少有一个未完成的科目)
	if err := db.Model(&model.Student{}).
		Joins("JOIN student_subjects ON students.id = student_subjects.student_id").
		Where("students.status = 1 AND student_subjects.status IN (1, 2) AND students.deleted_at IS NULL").
		Distinct("students.id").
		Count(&summary.TotalStudents).Error; err != nil {
		logger.Error("Failed to count active students", logger.ErrorType(err))
	}

	// 2. 本月新增学员 (根据student_subjects的创建时间)
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if err := db.Model(&model.Student{}).
		Joins("JOIN student_subjects ON students.id = student_subjects.student_id").
		Where("student_subjects.created_at >= ? AND students.deleted_at IS NULL", startOfMonth).
		Distinct("students.id").
		Count(&summary.NewStudentsThisMonth).Error; err != nil {
		logger.Error("Failed to count new students", logger.ErrorType(err))
	}

	// 3. 剩余总课时 (StudentSubject.balance求和)
	if err := db.Model(&model.StudentSubject{}).
		Where("status IN (1, 2) AND deleted_at IS NULL").
		Select("COALESCE(SUM(balance), 0)").
		Scan(&summary.TotalRemainingHours).Error; err != nil {
		logger.Error("Failed to sum remaining hours", logger.ErrorType(err))
	}

	// 4. 欠费与预警人数 (基于StudentSubject)
	// 欠费学生(某科目欠费)
	var arrearsCount int64
	if err := db.Model(&model.StudentSubject{}).
		Where("balance < 0 AND status IN (1, 2) AND deleted_at IS NULL").
		Select("COUNT(DISTINCT student_id)").
		Scan(&arrearsCount).Error; err != nil {
		logger.Error("Failed to count arrears students", logger.ErrorType(err))
	}
	summary.TotalArrears = arrearsCount

	// 预警学生(0 <= balance < 5的课程)
	var warningCount int64
	if err := db.Model(&model.StudentSubject{}).
		Where("balance >= 0 AND balance < 5 AND status IN (1, 2) AND deleted_at IS NULL").
		Select("COUNT(DISTINCT student_id)").
		Scan(&warningCount).Error; err != nil {
		logger.Error("Failed to count warning students", logger.ErrorType(err))
	}
	summary.TotalWarning = warningCount
	// 5. 本月消课数 (Records active=true)
	nextMonth := startOfMonth.AddDate(0, 1, 0)

	// 本月消课数
	var currentMonthCount int64
	if err := db.Model(&model.Record{}).
		Where("active = 1 AND deleted_at IS NULL AND teaching_date >= ? AND teaching_date < ?", startOfMonth, nextMonth).
		Count(&currentMonthCount).Error; err != nil {
		logger.Error("Failed to count monthly records", logger.ErrorType(err))
	}
	summary.MonthlyHours = currentMonthCount

	// 上月消课数 (用于计算环比)
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)
	var lastMonthCount int64
	if err := db.Model(&model.Record{}).
		Where("active = 1 AND deleted_at IS NULL AND teaching_date >= ? AND teaching_date < ?", startOfLastMonth, startOfMonth).
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

	// 根据 rangeType 动态生成 SQL 时间范围
	var xAxis []string
	var sqlFormat string
	var queryStartDate string

	now := time.Now()

	// Helper to generate monthly axis
	generateMonthly := func(months int) {
		y, m, _ := now.Date()
		// Start from (months-1) ago. e.g. 6m -> start 5 months ago
		startMonth := time.Date(y, m, 1, 0, 0, 0, 0, now.Location()).AddDate(0, -(months - 1), 0)

		queryStartDate = startMonth.Format("2006-01-02")
		sqlFormat = "%Y-%m"

		for i := 0; i < months; i++ {
			xAxis = append(xAxis, startMonth.AddDate(0, i, 0).Format("2006-01"))
		}
	}

	switch rangeType.Type {
	case "1m":
		// Daily, last 1 month
		start := now.AddDate(0, -1, 0)
		queryStartDate = start.Format("2006-01-02")
		sqlFormat = "%Y-%m-%d"
		for d := start; !d.After(now); d = d.AddDate(0, 0, 1) {
			xAxis = append(xAxis, d.Format("2006-01-02"))
		}
	case "12m":
		generateMonthly(12)
	case "all":
		// 从数据库中找出最早的订单或记录日期作为起点
		var minOrderTimeStr *string
		var minRecordTimeStr *string

		db.Model(&model.RechargeOrder{}).Select("MIN(created_at)").Scan(&minOrderTimeStr)
		db.Model(&model.Record{}).Select("MIN(teaching_date)").Scan(&minRecordTimeStr)

		startTime := now
		found := false

		parseTime := func(timeStr *string) *time.Time {
			if timeStr == nil {
				return nil
			}
			// 尝试解析常见格式，SQLite 默认存储可能是 "2006-01-02 15:04:05.999999999-07:00" 或简单的字符串
			// 这里尝试几种常见格式
			formats := []string{
				"2006-01-02 15:04:05-07:00",
				"2006-01-02 15:04:05.999999999-07:00",
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05Z07:00",
				"2006-01-02",
			}
			for _, f := range formats {
				if t, err := time.ParseInLocation(f, *timeStr, time.Local); err == nil {
					return &t
				}
			}
			return nil
		}

		if t := parseTime(minOrderTimeStr); t != nil {
			startTime = *t
			found = true
		}
		if t := parseTime(minRecordTimeStr); t != nil {
			if !found || t.Before(startTime) {
				startTime = *t
				found = true
			}
		}

		if !found {
			generateMonthly(12)
		} else {
			months := (now.Year()-startTime.Year())*12 + int(now.Month()-startTime.Month()) + 1
			if months < 1 {
				months = 1
			}
			generateMonthly(months)
		}
	case "6m":
		fallthrough
	default:
		generateMonthly(6)
	}
	result.XAxis = xAxis

	type ChartStat struct {
		Label  string
		Hours  int64
		Amount float64
	}

	// 1. 充值数据 (RechargeOrder: 课时 + 金额)
	var rechargeStats []ChartStat
	err := db.Model(&model.RechargeOrder{}).
		Select(fmt.Sprintf("strftime('%s', created_at) as label, SUM(hours) as hours, SUM(amount) as amount", sqlFormat)).
		Where("deleted_at IS NULL AND created_at >= ?", queryStartDate).
		Group("label").
		Order("label").
		Scan(&rechargeStats).Error
	if err != nil {
		return result, err
	}

	// 2. 消课数据 (Record: 课时数 + 金额=课时数*单价)
	var consumeStats []ChartStat
	err = db.Model(&model.Record{}).
		Select(fmt.Sprintf("strftime('%s', teaching_date) as label, COUNT(id) as hours, COALESCE(SUM(price_snapshot), 0) as amount", sqlFormat)).
		Where("active = 1 AND deleted_at IS NULL AND teaching_date >= ?", queryStartDate).
		Group("label").
		Order("label").
		Scan(&consumeStats).Error
	if err != nil {
		return result, err
	}

	// 3. 数据填充 (Map to Slice)
	rechargeMap := make(map[string]ChartStat)
	for _, s := range rechargeStats {
		rechargeMap[s.Label] = s
	}

	consumeMap := make(map[string]ChartStat)
	for _, s := range consumeStats {
		consumeMap[s.Label] = s
	}

	for _, label := range xAxis {
		rechargeStat := rechargeMap[label]
		consumeStat := consumeMap[label]

		result.RechargeData = append(result.RechargeData, rechargeStat.Hours)
		result.RechargeAmount = append(result.RechargeAmount, rechargeStat.Amount)

		result.ConsumeData = append(result.ConsumeData, consumeStat.Hours)
		result.ConsumeAmount = append(result.ConsumeAmount, consumeStat.Amount)

		result.NetData = append(result.NetData, rechargeStat.Hours-consumeStat.Hours)
	}

	return result, nil
}

// // GetHeatmapData 获取热力图数据 (适配 ChartHeatmap.vue 组件)
// func (m *DashboardManager) GetHeatmapData(ctx context.Context) ([][]int, error) {
// 	db := dao.GetDB()

// 	// 结构体接收数据库聚合结果
// 	type HeatStat struct {
// 		DayOfWeek int    // SQLite strftime('%w') 返回 0-6 (0是周日)
// 		Hour      string // 格式如 "08", "14"
// 		Count     int
// 	}
// 	var stats []HeatStat

// 	// 分析最近 6 个月的数据
// 	startDate := time.Now().AddDate(0, -6, 0)

// 	// SQLite 查询: 聚合 星期几 和 小时
// 	// 注意: start_time 在数据库中可能是 "08:00:00" 或 "08:00"，substr(start_time, 1, 2) 取前两位
// 	err := db.Model(&model.Record{}).
// 		Select("CAST(strftime('%w', teaching_date) AS INTEGER) as day_of_week, "+
// 			"substr(start_time, 1, 2) as hour, "+
// 			"COUNT(id) as count").
// 		Where("active = 1 AND deleted_at IS NULL AND teaching_date >= ?", startDate).
// 		Group("day_of_week, hour").
// 		Scan(&stats).Error

// 	if err != nil {
// 		logger.Error("Failed to get heatmap data", logger.ErrorType(err))
// 		return nil, err
// 	}

// 	var chartData [][]int

// 	for _, s := range stats {
// 		hourInt, err := strconv.Atoi(s.Hour)
// 		if err != nil {
// 			continue
// 		}

// 		// 过滤营业时间以外的数据 (08:00 - 21:00)
// 		if hourInt < 8 || hourInt > 21 {
// 			continue
// 		}

// 		// 直接返回原始数据 [DayOfWeek, Hour, Count]，由前端负责视图映射
// 		// DayOfWeek: 0(Sun) - 6(Sat)
// 		// Hour: 8 - 21
// 		chartData = append(chartData, []int{s.DayOfWeek, hourInt, s.Count})
// 	}

// 	return chartData, nil
// }

// // GetStudentEngagementData 获取学员活跃度分布 (基于过去30天课次)
// // Dormant: 0, Lazy: 1-3, Regular: 4-8, High: >8
// func (m *DashboardManager) GetStudentEngagementData(ctx context.Context) (responsex.GetStudentEngagementDataResponse, error) {
// 	db := dao.GetDB()

// 	type EngagementStat struct {
// 		FrequencyLevel string
// 		StudentCount   int
// 	}
// 	var stats []EngagementStat

// 	// 使用 Go 计算时间，避免数据库方言差异和时区问题
// 	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

// 	// 优化查询：使用 LEFT JOIN 替代子查询，提高性能
// 	query := `
// 		SELECT
// 			CASE
// 				WHEN lesson_count = 0 THEN 'Dormant'
// 				WHEN lesson_count >= 1 AND lesson_count <= 3 THEN 'Lazy'
// 				WHEN lesson_count >= 4 AND lesson_count <= 8 THEN 'Regular'
// 				ELSE 'High'
// 			END as frequency_level,
// 			COUNT(*) as student_count
// 		FROM (
// 			SELECT
// 				s.id,
// 				COUNT(r.id) as lesson_count
// 			FROM students s
// 			LEFT JOIN records r ON s.id = r.student_id
// 				AND r.active = 1
// 				AND r.deleted_at IS NULL
// 				AND r.teaching_date >= ?
// 			WHERE s.deleted_at IS NULL
// 			GROUP BY s.id
// 		)
// 		GROUP BY frequency_level
// 	`

// 	if err := db.Raw(query, startDate).Scan(&stats).Error; err != nil {
// 		logger.Error("Failed to get student engagement data", logger.ErrorType(err))
// 		return responsex.GetStudentEngagementDataResponse{}, err
// 	}

// 	// 转换为前端需要的格式，确保所有类型都有值（即使是0）
// 	resultMap := map[string]int{
// 		"Dormant": 0,
// 		"Lazy":    0,
// 		"Regular": 0,
// 		"High":    0,
// 	}
// 	for _, s := range stats {
// 		resultMap[s.FrequencyLevel] = s.StudentCount
// 	}

// 	// 映射到前端展示名称
// 	// 顺序：沉睡 -> 消极 -> 达标 -> 高频
// 	return responsex.GetStudentEngagementDataResponse{
// 		Stats: []responsex.EngagementStat{
// 			{Name: "沉睡 (0次)", Value: resultMap["Dormant"]},
// 			{Name: "消极 (1-3次)", Value: resultMap["Lazy"]},
// 			{Name: "达标 (4-8次)", Value: resultMap["Regular"]},
// 			{Name: "高频 (>8次)", Value: resultMap["High"]},
// 		},
// 	}, nil
// }

// // GetStudentGrowthData 获取学员增长趋势数据 (最近 6 个月)
// func (m *DashboardManager) GetStudentGrowthData(ctx context.Context) (responsex.ChartDataDTO, error) {
// 	db := dao.GetDB()
// 	var result responsex.ChartDataDTO

// 	// 生成最近 6 个月的月份标签
// 	months := []string{}
// 	now := time.Now()
// 	// 规范化到月初，避免月末(如31号)计算月份偏移时出现跳跃或重复
// 	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
// 	limit := 6

// 	for i := limit - 1; i >= 0; i-- {
// 		t := startOfMonth.AddDate(0, -i, 0)
// 		months = append(months, t.Format("2006-01"))
// 	}
// 	result.XAxis = months

// 	// 计算起始日期
// 	startDate := startOfMonth.AddDate(0, -limit+1, 0).Format("2006-01") + "-01"

// 	type MonthlyStat struct {
// 		Month string
// 		Total int64
// 	}
// 	var stats []MonthlyStat

// 	// 查询数据库: 按月统计 created_at
// 	err := db.Model(&model.Student{}).
// 		Select("strftime('%Y-%m', created_at) as month, COUNT(id) as total").
// 		Where("deleted_at IS NULL AND created_at >= ?", startDate).
// 		Group("month").
// 		Order("month").
// 		Scan(&stats).Error

// 	if err != nil {
// 		logger.Error("Failed to get student growth data", logger.ErrorType(err))
// 		return result, err
// 	}

// 	// 填充数据 (Map to Slice)
// 	dataMap := make(map[string]int64)
// 	for _, s := range stats {
// 		dataMap[s.Month] = s.Total
// 	}

// 	for _, m := range months {
// 		result.Series = append(result.Series, dataMap[m])
// 	}

// 	return result, nil
// }

// // GetStudentBalanceData 获取学员账户健康度分布 (基于剩余课时)
// // Arrears: <0, Warning: 0-5, Sufficient: >=5
// func (m *DashboardManager) GetStudentBalanceData(ctx context.Context) (responsex.GetStudentBalanceDataResponse, error) {
// 	db := dao.GetDB()

// 	type BalanceStat struct {
// 		BalanceLevel string
// 		StudentCount int
// 	}
// 	var stats []BalanceStat

// 	// 使用 CASE WHEN 对StudentSubject的balance进行分桶统计
// 	// 统计有不同balance等级的学生数（去重）
// 	query := `
// 		SELECT
// 			CASE
// 				WHEN balance < 0 THEN 'Arrears'
// 				WHEN balance >= 0 AND balance < 5 THEN 'Warning'
// 				ELSE 'Sufficient'
// 			END as balance_level,
// 			COUNT(DISTINCT student_id) as student_count
// 		FROM student_subjects
// 		WHERE status IN (1, 2) AND deleted_at IS NULL
// 		GROUP BY balance_level
// 	`

// 	if err := db.Raw(query).Scan(&stats).Error; err != nil {
// 		logger.Error("Failed to get student balance data", logger.ErrorType(err))
// 		return responsex.GetStudentBalanceDataResponse{}, err
// 	}

// 	// 初始化 Map 确保所有状态都有值 (即使数据库查出来是 0)
// 	resultMap := map[string]int{
// 		"Arrears":    0,
// 		"Warning":    0,
// 		"Sufficient": 0,
// 	}
// 	for _, s := range stats {
// 		resultMap[s.BalanceLevel] = s.StudentCount
// 	}

// 	// 转换为前端 ECharts Pie 图所需格式
// 	// 顺序建议：充足 -> 预警 -> 欠费
// 	return responsex.GetStudentBalanceDataResponse{
// 		Stats: []responsex.BalanceStat{
// 			{Name: "充足 (>=5课时)", Value: resultMap["Sufficient"]},
// 			{Name: "预警 (<5课时)", Value: resultMap["Warning"]},
// 			{Name: "欠费 (<0课时)", Value: resultMap["Arrears"]},
// 		},
// 	}, nil
// }

// GetSubjectRankData 获取科目消课排行 (本月各科目消课占比)
func (m *DashboardManager) GetSubjectRankData(ctx context.Context) (responsex.GetSubjectRankDataResponse, error) {
	db := dao.GetDB()
	var result responsex.GetSubjectRankDataResponse

	// 本月起始日期
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := startOfMonth.AddDate(0, 1, 0)

	// 查询本月各科目消课数 (join Record 和 Subject)
	type SubjectRank struct {
		Name  string `gorm:"column:name"`
		Count int64  `gorm:"column:count"`
	}

	var ranks []SubjectRank
	err := db.Model(&model.Record{}).
		Select("subjects.name, COUNT(records.id) as count").
		Joins("JOIN subjects ON records.subject_id = subjects.id").
		Where("records.active = 1 AND records.deleted_at IS NULL AND subjects.deleted_at IS NULL AND records.teaching_date >= ? AND records.teaching_date < ?", startOfMonth, nextMonth).
		Group("subjects.id, subjects.name").
		Order("count DESC").
		Scan(&ranks).Error

	if err != nil {
		logger.Error("Failed to get subject rank data", logger.ErrorType(err))
		return result, err
	}

	// 转换为前端所需格式
	for _, rank := range ranks {
		result.Data = append(result.Data, responsex.SubjectRank{
			Name:  rank.Name,
			Value: rank.Count,
		})
	}

	return result, nil
}

func (m *DashboardManager) RegisterRoute(d *dispatcher.Dispatcher) {
	dispatcher.RegisterNoReq(d, "dashboard_manager/get_summary", m.GetSummaryData)
	dispatcher.RegisterTyped(d, "dashboard_manager/get_finance_chart", m.GetFinanceChartData)
	dispatcher.RegisterNoReq(d, "dashboard_manager/get_subject_rank", m.GetSubjectRankData)
	// dispatcher.RegisterNoReq(d, "dashboard_manager/get_heatmap", m.GetHeatmapData)
	// dispatcher.RegisterNoReq(d, "dashboard_manager/get_student_engagement", m.GetStudentEngagementData)
	// dispatcher.RegisterNoReq(d, "dashboard_manager/get_student_growth", m.GetStudentGrowthData)
	// dispatcher.RegisterNoReq(d, "dashboard_manager/get_student_balance", m.GetStudentBalanceData)
}
