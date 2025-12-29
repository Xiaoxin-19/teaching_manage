// 定义前端通用的数据模型接口

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

export interface StudentItem {
  id: number;
  name: string;
  phone: string;
  balance: number;
  gender: string;
  teacher_id: number | string | null;
  note: string;
  lastModified?: string;
}

export interface StudentData {
  id?: number;
  name: string;
  phone: string;
  balance: number;
  gender: string;
  teacher_id: number | string | null;
  note: string;
  lastModified?: string;
}

export interface RecordItem {
  id: number;
  date: string;
  type: string;
  amount: number;
  remark: string;
}

export interface ResponseWrapper<T> {
  code: number
  message: string
  data: T
}

export type FetchDetailsFn = (studentId: number) => Promise<{ data: RecordItem[] }>;

// 分类标签类型
export interface OrderTag {
  label: string;
  color: string;
}

export interface OrderDTO {
  id: number;
  active: boolean;
  type: string;
  hours: number;
  comment: string;
  created_at: number;
  updated_at: number;
}