import { Dispatch } from "../../wailsjs/go/main/App";
import { ResponseWrapper, Subject } from "../types/appModels";
import { CreateSubjectRequest, DeleteSubjectRequest, GetSubjectListRequest, UpdateSubjectRequest } from "../types/request";
import { GetSubjectListResponse } from "../types/response";



// --- API 方法 ---

export async function GetSubjectList(req: GetSubjectListRequest): Promise<GetSubjectListResponse> {
  try {
    // 假设后端接口: subject_manager/get_subject_list
    // 参数: {"keyword": "..."}
    const reqStr = JSON.stringify(req);
    const resultStr = await Dispatch('subject_manager/get_subject_list', reqStr);
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetSubjectListResponse>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '获取科目列表失败');
    }
    return resp.data;
  } catch (error: any) {
    console.error("API Error [GetSubjectList]:", error);
    throw error;
  }
}

export async function CreateSubject(req: CreateSubjectRequest): Promise<void> {
  try {
    const reqStr = JSON.stringify(req);
    const resultStr = await Dispatch('subject_manager/create_subject', reqStr);
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '创建科目失败');
    }
  } catch (error: any) {
    throw error;
  }
}

export async function UpdateSubject(req: UpdateSubjectRequest): Promise<void> {
  try {
    const reqStr = JSON.stringify(req);
    const resultStr = await Dispatch('subject_manager/update_subject', reqStr);
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '更新科目失败');
    }
  } catch (error: any) {
    throw error;
  }
}

export async function DeleteSubject(req: DeleteSubjectRequest): Promise<void> {
  try {
    const reqStr = JSON.stringify(req);
    const resultStr = await Dispatch('subject_manager/delete_subject', reqStr);
    const resp = JSON.parse(resultStr) as ResponseWrapper<string>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '删除科目失败');
    }
  } catch (error: any) {
    throw error;
  }
}