// 后端返回的原始数据结构
export interface StudentDTO {
  id: number
  name: string
  gender: string
  hours: number
  phone: string
  remark: string
  teacher_id: number
  created_at: number
  updated_at: number
}

export interface GetTeacherListResponse {
  teachers: TeacherDTO[]
  total: number
}

export interface GetStudentListResponse {
  students: StudentDTO[]
  total: number
}

export interface GetOrdersByStudentIdResponse {
  orders: OrderDTO[]
  total: number
}

export interface RecordDTO {
  id: number
  student_id: number
  teacher_id: number
  student_name: string
  teacher_name: string
  teaching_date: string
  start_time: string
  end_time: string
  active: boolean
  remark: string
  created_at: number
  updated_at: number
}
export interface GetRecordListResponse {
  records: RecordDTO[]
  total: number
  total_pending: number
}

export interface ImportExcelResponse {
  filepath: string
  total_rows: number
  error_infos: string[][]
}

export interface SelectFileResponse {
  filepath: string
}

export interface GetDashboardSummaryResponse {
  "total_students": number
  "new_students_this_month": number
  "monthly_hours": number
  "month_over_month": string
  "total_remaining_hours": number
  "total_arrears": number
  "total_warning": number
}

export interface GetFinanceChartDataResponse {
  x_axis: string[]
  recharge_data: number[]
  consume_data: number[]
  net_data: number[]
}

export interface EngagementStat {
  name: string
  value: number
}

export interface GetStudentEngagementDataResponse {
  stats: EngagementStat[]
}

export interface GetStudentGrowthDataResponse {
  x_axis: string[]
  series: number[]
}

export interface GetTeacherRankDataResponse {
  names: string[]
  values: number[]
}

export interface BalanceStat {
  name: string
  value: number
}

export interface GetStudentBalanceDataResponse {
  stats: BalanceStat[]
}
