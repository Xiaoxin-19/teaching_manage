import { Dispatch } from "../../wailsjs/go/main/App";
import { ResponseWrapper } from "../types/appModels";
import { GetDashboardSummaryResponse, GetFinanceChartDataResponse, GetSubjectRankDataResponse } from "../types/response";
import type { GetFinanceChartRequest } from "../types/request";

/**
 * 获取仪表板摘要数据
 * @returns Dashboard摘要信息
 */
export async function GetDashboardSummary(): Promise<GetDashboardSummaryResponse> {
  try {
    const resultStr = await Dispatch('dashboard_manager/get_summary', '');
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetDashboardSummaryResponse>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '获取概要数据失败');
    }

    return resp.data;
  } catch (error: any) {
    console.error('API Error [GetDashboardSummary]:', error);
    throw error;
  }
}

/**
 * 获取财务流转数据（课时和金额）
 * @param rangeType - 时间范围: "1m" (近一月), "6m" (近半年), "12m" (近一年), "all" (全部)
 * @returns 财务图表数据
 */
export async function GetFinanceChartData(rangeType: GetFinanceChartRequest['type'] = '6m'): Promise<GetFinanceChartDataResponse> {
  try {
    const req = JSON.stringify({ type: rangeType });
    const resultStr = await Dispatch('dashboard_manager/get_finance_chart', req);
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetFinanceChartDataResponse>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '获取财务数据失败');
    }

    return resp.data;
  } catch (error: any) {
    console.error('API Error [GetFinanceChartData]:', error);
    throw error;
  }
}

/**
 * 获取科目消课排行数据
 * @returns 科目排行数据 (本月各科目消课占比)
 */
export async function GetSubjectRankData(): Promise<GetSubjectRankDataResponse> {
  try {
    const resultStr = await Dispatch('dashboard_manager/get_subject_rank', '');
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetSubjectRankDataResponse>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '获取科目排行数据失败');
    }

    return resp.data;
  } catch (error: any) {
    console.error('API Error [GetSubjectRankData]:', error);
    throw error;
  }
}
