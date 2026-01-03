

import { Dispatch } from "../../wailsjs/go/main/App";
import { Order, ResponseWrapper } from "../types/appModels";
import { GetOrderListRequest } from "../types/request";
import { GetOrderListResponse } from "../types/response";


export async function GetOrderList(req: GetOrderListRequest): Promise<GetOrderListResponse> {
  try {
    const reqStr = JSON.stringify(req);
    const resp = await Dispatch("order_manager/get_order_list", reqStr);
    const respData = JSON.parse(resp) as ResponseWrapper<GetOrderListResponse>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  } catch (error) {
    console.error("Failed to fetch order list:", error);
    throw error;
  }
}

export async function ExportOrdersToExcel(req: GetOrderListRequest): Promise<string> {
  try {
    const reqStr = JSON.stringify(req);
    const resp = await Dispatch("order_manager/export_orders_to_excel", reqStr);
    const respData = JSON.parse(resp) as ResponseWrapper<string>;
    if (respData.code !== 200) {
      throw new Error(respData.message);
    }
    return respData.data;
  } catch (error) {
    console.error("Failed to export orders to Excel:", error);
    throw error;
  }
}