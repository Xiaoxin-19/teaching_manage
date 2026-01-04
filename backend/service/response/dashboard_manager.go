package responsex

// DashboardSummaryResponse 核心指标卡
type DashboardSummaryResponse struct {
	TotalStudents        int64  `json:"total_students"`
	NewStudentsThisMonth int64  `json:"new_students_this_month"` // 本月新增
	MonthlyHours         int64  `json:"monthly_hours"`           // 本月消课数 (节)
	MonthOverMonth       string `json:"month_over_month"`        // 环比增长
	TotalRemainingHours  int64  `json:"total_remaining_hours"`   // 剩余总课时
	TotalArrears         int64  `json:"total_arrears"`           // 欠费人数
	TotalWarning         int64  `json:"total_warning"`           // 预警人数
}

// ChartDataDTO 通用图表数据结构
type ChartDataDTO struct {
	XAxis  []string `json:"x_axis"`
	Series []int64  `json:"series"`
}

// FinanceChartDTO 资金/课时流转图表
type FinanceChartDTO struct {
	XAxis          []string  `json:"x_axis"`
	RechargeData   []int64   `json:"recharge_data"`   // 充值课时
	RechargeAmount []float64 `json:"recharge_amount"` // 充值金额
	ConsumeData    []int64   `json:"consume_data"`    // 消课课时
	ConsumeAmount  []float64 `json:"consume_amount"`  // 消课金额
	NetData        []int64   `json:"net_data"`        // 净增库存(课时)
}

// TeacherRankDTO 教师排行
type TeacherRankDTO struct {
	Names  []string `json:"names"`
	Values []int64  `json:"values"`
}

type EngagementStat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type GetStudentEngagementDataResponse struct {
	Stats []EngagementStat `json:"stats"`
}

type BalanceStat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type GetStudentBalanceDataResponse struct {
	Stats []BalanceStat `json:"stats"`
}

// SubjectRank 科目排行数据项
type SubjectRank struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// GetSubjectRankDataResponse 科目消课占比响应
type GetSubjectRankDataResponse struct {
	Data []SubjectRank `json:"data"`
}

// StudentGrowthTrendResponse 学员增长和流失趋势
type StudentGrowthTrendResponse struct {
	XAxis  []string `json:"x_axis"` // 月份标签 (YYYY-MM)
	Growth []int64  `json:"growth"` // 新增学员数
	Loss   []int64  `json:"loss"`   // 流失学员数 (删除)
	Net    []int64  `json:"net"`    // 净增 (growth - loss)
}
