// 后端返回的原始数据结构
export interface StudentDTO {
  id: number
  name: string
  gender: string
  hours: number
  phone: string
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