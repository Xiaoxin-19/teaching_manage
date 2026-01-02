import { Dispatch } from "../../wailsjs/go/main/App";
import { ResponseWrapper } from "../types/appModels";

// --- 模型定义 ---

export interface Course {
  id: number;
  studentName: string;
  studentCode: string;
  studentStatus: number; // 1:正常, 2:停课, 3:退学
  subjectName: string;
  subjectId: number;
  teacherName: string;
  teacherCode: string;
  teacherId: number;
  balance: number;
  totalBuy: number;
  courseStatus: number; // 1:正常, 2:暂停, 3:结课
  updatedAt?: string;
}

export interface GetCourseListRequest {
  search?: string;
  subjects?: string[];
  teachers?: string[];
  balanceMin?: number | null;
  balanceMax?: number | null;
  status?: string[];
  page: number;
  pageSize: number;
}

export interface GetCourseListResponse {
  courses: Course[];
  total: number;
}

export interface EnrollRequest {
  student_Id: number;
  subject_Id: number;
  teacher_Id: number;
  remark?: string;
}

export interface RechargeRequest {
  courseId: number;
  hours: number; // 正数为充值，负数为扣除
  amount: number; // 实付/退费金额 (虽然前端移除了金额输入，但后端接口可能仍需保留字段兼容，传0即可)
  remark: string;
}

export interface DeleteCourseRequest {
  courseId: number;
  isHardDelete: boolean; // false: 软删除(结课+清算), true: 硬删除
  remark?: string; // 软删除时的备注
}

export interface UpdateCourseRequest {
  id: number;
  teacherId: number;
  remark?: string;
}

// --- API 方法 ---

export async function GetCourseList(req: GetCourseListRequest): Promise<GetCourseListResponse> {
  try {
    const reqStr = JSON.stringify(req);
    const resultStr = await Dispatch('course_manager/get_list', reqStr);
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetCourseListResponse>;

    if (resp.code !== 200) throw new Error(resp.message || '获取课程列表失败');
    return resp.data;
  } catch (error: any) {
    console.error("API Error [GetCourseList]:", error);
    throw error;
  }
}

export async function EnrollCourse(data: EnrollRequest): Promise<void> {
  try {
    const req = JSON.stringify(data);
    const resultStr = await Dispatch('course_manager/enroll', req);
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>;
    if (resp.code !== 200) throw new Error(resp.message);
  } catch (error: any) { throw error; }
}

export async function RechargeCourse(data: RechargeRequest): Promise<void> {
  try {
    const req = JSON.stringify(data);
    const resultStr = await Dispatch('course_manager/recharge', req);
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>;
    if (resp.code !== 200) throw new Error(resp.message);
  } catch (error: any) { throw error; }
}

export async function ToggleCourseStatus(courseId: number): Promise<void> {
  try {
    const req = JSON.stringify({ courseId });
    const resultStr = await Dispatch('course_manager/toggle_status', req);
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>;
    if (resp.code !== 200) throw new Error(resp.message);
  } catch (error: any) { throw error; }
}

export async function DeleteCourse(data: DeleteCourseRequest): Promise<void> {
  try {
    const req = JSON.stringify(data);
    const resultStr = await Dispatch('course_manager/delete', req);
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>;
    if (resp.code !== 200) throw new Error(resp.message);
  } catch (error: any) { throw error; }
}

export async function UpdateCourse(data: UpdateCourseRequest): Promise<void> {
  try {
    const req = JSON.stringify(data);
    const resultStr = await Dispatch('course_manager/update', req);
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>;
    if (resp.code !== 200) throw new Error(resp.message);
  } catch (error: any) { throw error; }
}