import { Dispatch } from "../../wailsjs/go/main/App"
import { ResponseWrapper } from "../types/appModels"
import { CreateStudentRequest, DeleteStudentRequest, GetStudentListRequest, UpdateStudentRequest } from "../types/request"
import { GetStudentListResponse } from "../types/response"

// 建议定义请求参数接口


export async function GetStudentList(reqData: GetStudentListRequest): Promise<GetStudentListResponse> {
  try {

    const resultStr = await Dispatch('student_manager/get_student_list', JSON.stringify(reqData))


    const resp = JSON.parse(resultStr) as ResponseWrapper<GetStudentListResponse>


    if (resp.code !== 200) {

      throw new Error(resp.message || '获取学生列表失败')
    }

    // 4. 成功：直接返回数据载荷 (Data Payload)
    // 注意：这里返回原始数据，数据格式化（如日期转换）建议在 Logic 层处理
    return resp.data;

  } catch (error: any) {
    console.error('API Error [GetStudentList]:', error)
    throw error
  }
}

export async function CreateStudent(data: CreateStudentRequest): Promise<void> {
  try {

    const resultStr = await Dispatch('student_manager/create_student', JSON.stringify(data))

    const resp = JSON.parse(resultStr) as ResponseWrapper<string>

    if (resp.code !== 200) {
      throw new Error(resp.message || '创建学生失败')
    }

    return;
  } catch (error: any) {

    console.error('API Error [CreateStudent]:', error)
    throw error
  }
}

export async function UpdateStudent(data: UpdateStudentRequest): Promise<void> {
  try {
    const resultStr = await Dispatch('student_manager/update_student', JSON.stringify(data))
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>
    if (resp.code !== 200) {
      throw new Error(resp.message || '更新学生失败')
    }
    return;
  } catch (error: any) {
    console.error('API Error [UpdateStudent]:', error)
    throw error
  }
}

export async function DeleteStudent(data: DeleteStudentRequest): Promise<void> {
  try {
    const resultStr = await Dispatch('student_manager/delete_student', JSON.stringify(data))
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>
    if (resp.code !== 200) {
      throw new Error(resp.message || '删除学生失败')
    }
    return;
  } catch (error: any) {
    console.error('API Error [DeleteStudent]:', error)
    throw error
  }
}

export async function ExportStudentsToExcel(): Promise<string> {
  try {
    const resultStr = await Dispatch('student_manager/export_students', '') // 无需参数
    const resp = JSON.parse(resultStr) as ResponseWrapper<string> // 假设返回文件路径字符串
    if (resp.code !== 200) {
      throw new Error(resp.message || '导出学生数据失败')
    }
    return resp.data; // 返回文件路径
  } catch (error: any) {
    console.error('API Error [ExportStudentsToExcel]:', error)
    throw error
  }
}