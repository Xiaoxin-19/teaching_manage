import { Dispatch } from "../../wailsjs/go/main/App";
import { ResponseWrapper } from "../types/appModels";
import { CreateTeacherRequest, DeleteTeacherRequest, GetTeacherListRequest, UpdateTeacherRequest } from "../types/request";
import { GetTeacherListResponse } from "../types/response";


export async function CreateTeacher(req: CreateTeacherRequest): Promise<void> {
  try {
    const resultStr = await Dispatch('teacher_manager/create_teacher', JSON.stringify(req))
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>
    if (resp.code !== 200) {
      throw new Error(resp.message || '创建教师失败')
    }
    return;
  } catch (error: any) {
    console.error("API Error [CreateTeacher]:", error);
    throw error
  }
}

export async function GetTeacherList(req: GetTeacherListRequest): Promise<GetTeacherListResponse> {
  try {
    const resultStr = await Dispatch('teacher_manager/get_teacher_list', JSON.stringify(req))
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetTeacherListResponse>
    if (resp.code !== 200) {
      throw new Error(resp.message || '获取教师列表失败')
    }
    return resp.data;
  } catch (error: any) {
    console.error("API Error [GetTeacherList]:", error);
    throw error
  }
}

export async function UpdateTeacher(req: UpdateTeacherRequest): Promise<void> {
  try {
    const resultStr = await Dispatch('teacher_manager/update_teacher', JSON.stringify(req))
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>
    if (resp.code !== 200) {
      throw new Error(resp.message || '更新教师失败')
    }
    return;
  } catch (error: any) {
    console.error("API Error [UpdateTeacher]:", error);
    throw error
  }
}

export async function DeleteTeacher(req: DeleteTeacherRequest): Promise<void> {
  try {
    const resultStr = await Dispatch('teacher_manager/delete_teacher', JSON.stringify(req))
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>
    if (resp.code !== 200) {
      throw new Error(resp.message || '删除教师失败')
    }
    return;
  } catch (error: any) {
    console.error("API Error [DeleteTeacher]:", error);
    throw error
  }
}

export async function ExportTeachersToExcel(): Promise<string> {
  try {
    const resultStr = await Dispatch('teacher_manager/export_teacher_to_excel', '')
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>
    if (resp.code !== 200) {
      throw new Error(resp.message || '导出教师数据失败')
    }
    return resp.data;
  }
  catch (error: any) {
    console.error("API Error [ExportTeachersToExcel]:", error);
    throw error
  }
}