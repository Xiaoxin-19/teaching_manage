import { Dispatch } from "../../wailsjs/go/main/App";
import { ResponseWrapper } from "../types/appModels";
import { CreateCourseRequest, DeleteCourseRequest, GetCourseListRequest, RechargeRequest, UpdateCourseRequest } from "../types/request";
import { GetCourseListResponse } from "../types/response";






// --- API 方法 ---

export async function GetCourseList(req: GetCourseListRequest): Promise<GetCourseListResponse> {
  try {
    const reqStr = JSON.stringify(req);
    const resultStr = await Dispatch('course_manager/get_course_list', reqStr);
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetCourseListResponse>;

    if (resp.code !== 200) throw new Error(resp.message || '获取课程列表失败');
    console.log("GetCourseList Response:", JSON.stringify(resp.data));
    return resp.data;
  } catch (error: any) {
    console.error("API Error [GetCourseList]:", error);
    throw error;
  }
}

export async function CreateCourse(data: CreateCourseRequest): Promise<void> {
  try {
    const req = JSON.stringify(data);
    console.log("CreateCourse Request:", req);
    const resultStr = await Dispatch('course_manager/create_course', req);
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
    const req = JSON.stringify({ course_id: courseId });
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