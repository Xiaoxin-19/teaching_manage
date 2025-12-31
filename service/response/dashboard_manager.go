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
	XAxis        []string `json:"x_axis"`
	RechargeData []int64  `json:"recharge_data"`
	ConsumeData  []int64  `json:"consume_data"`
	NetData      []int64  `json:"net_data"`
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
