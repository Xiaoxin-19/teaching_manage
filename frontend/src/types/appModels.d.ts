// ... existing code ...
export interface TeacherOption {
  id: number | string
  name: string
}

export interface Student {
  id: number;
  student_number: string;
  name: string;
  phone: string;
  gender: string;
  remark: string;
  status: number;
  lastModified?: string;
}

export interface StudentData {
  id?: number;
  name: string;
  phone: string;
  gender: string;
  status: number;
  remark: string;
  lastModified?: string;
}


export interface Teacher {
  id: number;
  teacher_number: string;
  name: string;
  phone: string;
  gender: string;
  remark: string;
  updated_at: string;
  created_at: string;
  lastModified?: string;
}

export interface Subject {
  id: number;
  subject_number: string;
  name: string;
  student_count: number; // 关联学员数
  created_at?: number;
  updated_at?: number;
  lastModified?: string;
}

export interface Course {
  id: number;
  student: Student;
  subject: Subject;
  teacher: Teacher;
  balance: number;
  remark: string;
  status: number;
  created_at: number;
  updated_at: number;
}


// 新增：教学记录状态类型
export type RecordStatus = 'active' | 'pending';

// 新增：教学记录接口定义
export interface TeachingRecord {
  id: number;
  status: RecordStatus;
  date: string;
  time: string; // 显示用的时间段字符串，如 "10:00-12:00"
  startTime: string;
  endTime: string;
  studentId: number | null;
  studentName: string;
  teacherId: number | null;
  teacherName: string;
  remark: string;
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

export interface Order {
  id: number
  order_number: string
  student: Student
  subject: Subject
  teacher: Teacher
  type: 'increase' | 'decrease'
  amount: number
  hours: number
  created_at: number
  updated_at: number
  remark: string
  tags: OrderTag[]
  // 仅用于前端模拟筛选的字段，对接真实后端时不需要
  _studentId?: number
  _subjectId?: number
  _subjectName?: string
}