import { CreateTeacher } from "../api/teacher";

export interface GetStudentListRequest {
  Offset: number;
  Limit: number;
  keyword: string;
  Status: number;
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