import { Teacher, Course, Subject } from "./appModels"

// 后端返回的原始数据结构
export interface StudentDTO {
  id: number
  student_number: string
  name: string
  gender: string
  phone: string
  remark: string
  status: number
  created_at: number
  updated_at: number
  deleted_at?: number
  lastModified?: string
}

export interface GetTeacherListResponse {
  teachers: Teacher[]
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
  student: StudentDTO
  teacher: Teacher
  subject: Subject
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
  recharge_data: number[]      // 充值课时
  recharge_amount: number[]    // 充值金额 (可选显示)
  consume_data: number[]       // 消课课时
  consume_amount: number[]     // 消课金额 (可选显示)
  net_data: number[]
}

export interface SubjectRank {
  name: string
  value: number
}

export interface GetSubjectRankDataResponse {
  data: SubjectRank[]
}

export interface EngagementStat {
  code: string  // 枚举键：Dormant, Lazy, Regular, High
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

export interface StudentGrowthTrendResponse {
  x_axis: string[]  // 月份标签 (YYYY-MM)
  growth: number[]  // 新增学员数
  loss: number[]    // 流失学员数 (删除)
  net: number[]     // 净增 (growth - loss)
}

export interface GetSubjectListResponse {
  subjects: Subject[];
  total: number;
}

// Student Course relate

export interface GetCourseListResponse {
  courses: Course[];
  total: number;
}

export interface GetOrderListResponse {
  orders: Order[];
  total: number;
}