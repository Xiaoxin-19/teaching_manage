import { CreateTeacher } from "../api/teacher";

export interface GetStudentListRequest {
  Offset: number;
  Limit: number;
  keyword: string;
  Status_Level: number;
  Status_Target: number;
}

export interface CreateStudentRequest {
  Name: string;
  Phone: string;
  Gender: string;
  Remark: string;
}

export interface UpdateStudentRequest {
  ID: number;
  Name: string;
  Phone: string;
  Gender: string;
  Remark: string;
  Status: number;
}

export interface DeleteStudentRequest {
  ID: number;
}

// Teacher related requests
export interface CreateTeacherRequest {
  Name: string;
  Phone: string;
  Gender: string;
  Remark: string;
}


export interface UpdateTeacherRequest {
  ID: number;
  Name: string;
  Phone: string
  Gender: string;
  Remark: string;
}

export interface GetTeacherListRequest {
  Offset: number;
  Limit: number;
  Keyword: string;
}

export interface DeleteTeacherRequest {
  ID: number;
}

// Subject related requests
// --- 类型定义 (建议后续移动到 src/types/request.d.ts 和 response.d.ts) ---


export interface GetSubjectListRequest {
  Offset: number;
  Limit: number;
  Keyword: string;
}

export interface CreateSubjectRequest {
  Name: string;
}

export interface UpdateSubjectRequest {
  ID: number;
  Name: string;
}

export interface DeleteSubjectRequest {
  ID: number;
}


// Course related requests

export interface CreateCourseRequest {
  student_Id: number;
  subject_Id: number;
  teacher_Id: number;
  remark?: string;
}

export interface GetCourseListRequest {
  Students?: number[];
  Subjects?: number[];
  Teachers?: number[];
  Balance_Min?: number | null;
  Balance_Max?: number | null;
  Status?: number[];
  Offset: number;
  Limit: number;
}

export interface RechargeRequest {
  course_id: number;
  hours: number; // 正数为充值，负数为扣除
  amount: number; // 实付/退费金额 (虽然前端移除了金额输入，但后端接口可能仍需保留字段兼容，传0即可)
  remark: string;
}

export interface DeleteCourseRequest {
  course_id: number;
  is_hard_delete: boolean; // false: 软删除(结课+清算), true: 硬删除
  remark?: string; // 软删除时的备注
}

export interface UpdateCourseRequest {
  id: number;
  teacher_id: number;
  remark?: string;
}

export interface GetOrderListRequest {
  student_id?: number;
  subject_ids?: number[];
  teacher_ids?: number[];
  type?: string[];
  date_start?: string;
  date_end?: string;
  Offset: number;
  Limit: number;
}

export interface CreateRecordRequest {
  student_id: number;
  subject_id: number;
  teaching_date: string; // "YYYY-MM-DD"
  start_time: string;    // "HH:mm"
  end_time: string;      // "HH:mm"
  remark?: string;
}

export interface GetRecordListRequest {
  student_ids?: number[];
  teacher_ids?: number[];
  subject_ids?: number[];
  start_date?: string;
  end_date?: string;
  offset: number;
  limit: number;
  active?: boolean | null;
}

export interface ActivateRecordRequest {
  id: number;
}

export interface DeleteRecordRequest {
  id: number;
}

export interface ExportRecordsRequest {
  student_ids?: number[];
  teacher_ids?: number[];
  subject_ids?: number[];
  start_date?: string;
  end_date?: string;
  active?: boolean | null;
}

export interface ImportFromExcelRequest {
  filepath: string;
}