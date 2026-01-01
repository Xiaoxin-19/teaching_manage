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