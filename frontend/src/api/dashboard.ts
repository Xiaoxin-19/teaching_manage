import { Dispatch } from "../../wailsjs/go/main/App";
import { ResponseWrapper } from "../types/appModels";
import { GetDashboardSummaryResponse, GetFinanceChartDataResponse, GetSubjectRankDataResponse, GetStudentEngagementDataResponse, StudentGrowthTrendResponse, GetStudentBalanceDataResponse } from "../types/response";
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
 * 获取学员活跃度分布数据
 * @returns 学员活跃度分布数据 (基于过去30天课次，按科目数归一化)
 */
export async function GetStudentEngagementData(): Promise<GetStudentEngagementDataResponse> {
  try {
    const resultStr = await Dispatch('dashboard_manager/get_student_engagement', '');
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetStudentEngagementDataResponse>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '获取学员活跃度数据失败');
    }

    return resp.data;
  } catch (error: any) {
    console.error('API Error [GetStudentEngagementData]:', error);
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
/**
 * 获取热力图数据
 * @returns 热力图数据 (高峰时段热力分布，最近6个月)
 */
export async function GetHeatmapData(): Promise<number[][]> {
  try {
    const resultStr = await Dispatch('dashboard_manager/get_heatmap', '');
    const resp = JSON.parse(resultStr) as ResponseWrapper<number[][]>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '获取热力图数据失败');
    }

    return resp.data || [];
  } catch (error: any) {
    console.error('API Error [GetHeatmapData]:', error);
    throw error;
  }
}

/**
 * 获取学员增长和流失趋势数据
 * @returns 学员增长和流失趋势 (最近6个月，包括新增、流失、净增)
 */
export async function GetStudentGrowthTrendData(): Promise<StudentGrowthTrendResponse> {
  try {
    const resultStr = await Dispatch('dashboard_manager/get_student_growth_trend', '');
    const resp = JSON.parse(resultStr) as ResponseWrapper<StudentGrowthTrendResponse>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '获取学员增长趋势数据失败');
    }

    return resp.data;
  } catch (error: any) {
    console.error('API Error [GetStudentGrowthTrendData]:', error);
    throw error;
  }
}

/**
 * 获取学员账户健康度分布数据
 * @returns 学员账户健康度分布 (基于剩余课时：欠费、预警、充足)
 */
export async function GetStudentBalanceData(): Promise<GetStudentBalanceDataResponse> {
  try {
    const resultStr = await Dispatch('dashboard_manager/get_student_balance', '');
    const resp = JSON.parse(resultStr) as ResponseWrapper<GetStudentBalanceDataResponse>;

    if (resp.code !== 200) {
      throw new Error(resp.message || '获取学员账户健康度数据失败');
    }

    return resp.data;
  } catch (error: any) {
    console.error('API Error [GetStudentBalanceData]:', error);
    throw error;
  }
}