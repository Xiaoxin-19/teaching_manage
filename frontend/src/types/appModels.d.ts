// Shared TypeScript interfaces for frontend components

export interface StudentData {
  id?: number | string
  name: string
  gender: string
  phone: string
  teacherId: number | string | null
  remark: string
  lastModified?: string
}

export interface TeacherOption {
  id: number | string
  name: string
}

export interface TeacherData {
  id?: number | string
  name: string
  phone: string
  remark: string
  gender: string
  lastModified?: string
}

export interface RecordItem {
  date: string
  type: string
  amount: number
  balanceAfter: number
  remark: string
}

export type FetchDetailsFn = (studentId: number) => Promise<{ studentName: string, records: RecordItem[] }>


export interface ResponseWrapper<T> {
  code: number
  message: string
  data: T
}