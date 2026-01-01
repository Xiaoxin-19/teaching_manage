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