import { Dispatch } from '../../wailsjs/go/main/App'
import { ResponseWrapper } from '../types/appModels'
import { OpenFileDialogResponse, SetWebDavConfigRequest } from '../types/request'
import { WebDavBackupItem, WebDavConfigResponse } from '../types/response'

export interface SetBackupLocalPathRequest {
  local_path: string
}

/**
 * 设置本地备份路径
 */
export async function SetBackupLocalPath(req: SetBackupLocalPathRequest): Promise<string> {
  try {
    const resultStr = await Dispatch('backup_manager/set_backup_local_path', JSON.stringify(req))
    const result = JSON.parse(resultStr) as ResponseWrapper<string>
    if (result.code !== 200) {
      throw new Error(result.message || '设置备份路径失败')
    }
    return result.data as string
  }
  catch (error: any) {
    console.error("API Error [SetBackupLocalPath]:", error);
    throw error
  }
}

/**
 * 获取本地备份路径
 */
export async function GetBackupLocalPath(): Promise<string> {
  try {
    const result = await Dispatch('backup_manager/get_backup_local_path', '{}')
    const resultParsed = JSON.parse(result) as ResponseWrapper<string>
    if (resultParsed.code !== 200) {
      throw new Error(resultParsed.message || '获取备份路径失败')
    }
    return resultParsed.data as string
  } catch (error: any) {
    console.error("API Error [GetBackupLocalPath]:", error);
    throw error
  }
}


export async function OpenFileDialog(req: OpenFileDialogResponse): Promise<string> {
  try {
    const resultStr = await Dispatch('backup_manager/open_path_selector_dialog', JSON.stringify(req))
    const result = JSON.parse(resultStr) as ResponseWrapper<string>
    if (result.code !== 200) {
      throw new Error(result.message || '打开文件对话框失败')
    }
    return result.data as string
  } catch (error: any) {
    console.error("API Error [OpenFileDialog]:", error);
    throw error
  }
}


export async function SetWebDavConfig(req: SetWebDavConfigRequest): Promise<string> {
  try {
    const resultStr = await Dispatch('backup_manager/set_webdav_config', JSON.stringify(req))
    const result = JSON.parse(resultStr) as ResponseWrapper<string>
    if (result.code !== 200) {
      throw new Error(result.message || '设置 WebDav 配置失败')
    }
    return result.data as string
  } catch (error: any) {
    console.error("API Error [SetWebDavConfig]:", error);
    throw error
  }
}

export async function GetWebDavConfig(): Promise<WebDavConfigResponse> {
  try {
    const resultStr = await Dispatch('backup_manager/get_webdav_config', '{}')
    const result = JSON.parse(resultStr) as ResponseWrapper<WebDavConfigResponse>
    if (result.code !== 200) {
      throw new Error(result.message || '获取 WebDav 配置失败')
    }
    return result.data as WebDavConfigResponse
  } catch (error: any) {
    console.error("API Error [GetWebDavConfig]:", error);
    throw error
  }
}

export async function TestWebDavConnection(req: SetWebDavConfigRequest): Promise<string> {
  try {
    const resultStr = await Dispatch('backup_manager/test_webdav_connection', JSON.stringify(req))
    const result = JSON.parse(resultStr) as ResponseWrapper<string>
    if (result.code !== 200) {
      throw new Error(result.message || '测试 WebDav 连接失败')
    }
    return result.data as string
  } catch (error: any) {
    console.error("API Error [TestWebDavConnection]:", error);
    throw error
  }
}

export async function ListWebDavBackups(): Promise<WebDavBackupItem[]> {
  try {
    const resultStr = await Dispatch('backup_manager/list_webdav_backups', '{}')
    const result = JSON.parse(resultStr) as ResponseWrapper<WebDavBackupItem[]>
    if (result.code !== 200) {
      throw new Error(result.message || '获取 WebDav 备份列表失败')
    }
    return result.data as WebDavBackupItem[]
  } catch (error: any) {
    console.error("API Error [ListWebDavBackups]:", error);
    throw error
  }
}