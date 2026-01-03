import { Dispatch } from "../../wailsjs/go/main/App";
import { ResponseWrapper } from "../types/appModels";
import { ActivateRecordRequest, CreateRecordRequest, DeleteRecordRequest, ExportRecordsRequest, GetRecordListRequest, ImportFromExcelRequest } from "../types/request";
import { GetRecordListResponse, ImportExcelResponse, SelectFileResponse } from "../types/response";

export async function GetRecordList(req: GetRecordListRequest): Promise<GetRecordListResponse> {
  try {
    const reqStr = JSON.stringify(req);
    const resp = await Dispatch("record_manager/get_record_list", reqStr);
    const respData = JSON.parse(resp) as ResponseWrapper<GetRecordListResponse>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  } catch (error) {
    console.error("Failed to get record list:", error);
    throw error;
  }
}

export async function CreateRecord(data: CreateRecordRequest): Promise<string> {
  try {
    const reqStr = JSON.stringify(data);
    const resp = await Dispatch("record_manager/create_record", reqStr);
    const respData = JSON.parse(resp) as ResponseWrapper<string>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  }
  catch (error) {
    console.error("Failed to create record:", error);
    throw error;
  }
}

export async function ActivateRecord(req: ActivateRecordRequest): Promise<string> {
  try {
    const reqStr = JSON.stringify(req);
    const resp = await Dispatch("record_manager/activate_record", reqStr);
    const respData = JSON.parse(resp) as ResponseWrapper<string>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  } catch (error) {
    console.error("Failed to activate record:", error);
    throw error;
  }
}

export async function ActivateAllPendingRecords(): Promise<string> {
  try {
    const resp = await Dispatch("record_manager/activate_all_pending_records", '');
    const respData = JSON.parse(resp) as ResponseWrapper<string>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  } catch (error) {
    console.error("Failed to activate all pending records:", error);
    throw error;
  }
}


export async function DeleteRecordByID(req: DeleteRecordRequest): Promise<string> {
  try {
    const reqStr = JSON.stringify(req);
    const resp = await Dispatch("record_manager/delete_record_by_id", reqStr);
    const respData = JSON.parse(resp) as ResponseWrapper<string>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  } catch (error) {
    console.error("Failed to delete record:", error);
    throw error;
  }
}


export async function ExportRecordToExcel(req: ExportRecordsRequest): Promise<string> {
  try {
    const reqStr = JSON.stringify(req);
    const resp = await Dispatch("record_manager/export_record_to_excel", reqStr);
    const respData = JSON.parse(resp) as ResponseWrapper<string>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  }
  catch (error) {
    console.error("Failed to export records to excel:", error);
    throw error;
  }
}


export async function DownloadRecordImportTemplate(): Promise<string> {
  try {
    const resp = await Dispatch("record_manager/download_import_template", '');
    const respData = JSON.parse(resp) as ResponseWrapper<string>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  }
  catch (error) {
    console.error("Failed to download record import template:", error);
    throw error;
  }
}

export async function SelectFilePath(): Promise<string> {
  try {
    const resp = await Dispatch("record_manager/select_import_file", '');
    const respData = JSON.parse(resp) as ResponseWrapper<SelectFileResponse>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data.filepath;
  } catch (error) {
    console.error("Failed to select file path:", error);
    throw error;
  }
}

export async function ImportFromExcel(req: ImportFromExcelRequest): Promise<ImportExcelResponse> {
  try {
    const reqStr = JSON.stringify(req);
    const resp = await Dispatch("record_manager/import_from_excel", reqStr);
    const respData = JSON.parse(resp) as ResponseWrapper<ImportExcelResponse>;
    if (respData.code !== 200) {
      const error = new Error(respData.message);
      (error as any).data = respData.data;
      throw error;
    }
    return respData.data;
  } catch (error) {
    console.error("Failed to import from excel:", error);
    throw error;
  }
}