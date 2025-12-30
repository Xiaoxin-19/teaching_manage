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