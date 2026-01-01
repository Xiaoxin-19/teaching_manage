<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="d-flex justify-space-between align-center mb-2">
      <div>
        <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
          <v-icon color="indigo" class="mr-2">mdi-scale-balance</v-icon>
          课时流转健康度
        </h3>
        <div class="text-caption text-medium-emphasis mt-1">分析“进水”与“出水”速度</div>
      </div>
      <v-btn-toggle v-model="range" density="compact" variant="outlined" color="indigo" mandatory
        @update:model-value="handleRangeChange">
        <v-btn value="1m" size="small">近一月</v-btn>
        <v-btn value="6m" size="small">近半年</v-btn>
        <v-btn value="12m" size="small">近一年</v-btn>
        <v-btn value="all" size="small">全部</v-btn>
      </v-btn-toggle>
    </div>
    <div ref="chartRef" class="chart-box"></div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import * as echarts from 'echarts';
import { useChart } from '../../composables/useChart';
import { Dispatch } from '../../../wailsjs/go/main/App';
import { ResponseWrapper } from '../../types/appModels';
import { GetFinanceChartDataResponse } from '../../types/response';

const range = ref('6m');
const chartRef = ref<HTMLElement | null>(null);

// 定义数据状态
const chartData = ref({
  xAxis: [] as string[],
  rechargeData: [] as number[],
  consumeData: [] as number[],
  netData: [] as number[]
});

const getOption = (isDark: boolean) => {
  const { xAxis, rechargeData, consumeData, netData } = chartData.value;

  return {
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { data: ['充值课时 (收入)', '消课课时 (营收)', '净增库存'], top: 0 },
    grid: { top: 40, left: 20, right: 20, bottom: 60, containLabel: true },
    // 添加 dataZoom 组件以支持大量数据滚动查看
    dataZoom: [
      {
        type: 'inside',
        start: 0,
        end: 100
      },
      {
        start: 0,
        end: 100,
        bottom: 10
      }
    ],
    xAxis: { type: 'category', data: xAxis, axisLine: { show: false }, axisTick: { show: false } },
    yAxis: [
      { type: 'value', splitLine: { lineStyle: { type: 'dashed', opacity: 0.3 } } },
      { type: 'value', show: false }
    ],
    series: [
      { name: '充值课时 (收入)', type: 'bar', barGap: 0, itemStyle: { color: '#5C6BC0', borderRadius: [4, 4, 0, 0] }, data: rechargeData },
      { name: '消课课时 (营收)', type: 'bar', itemStyle: { color: '#66BB6A', borderRadius: [4, 4, 0, 0] }, data: consumeData },
      { name: '净增库存', type: 'line', yAxisIndex: 1, smooth: true, itemStyle: { color: '#FFA726' }, lineStyle: { width: 3, type: 'dashed' }, data: netData }
    ]
  };
};

const { refresh } = useChart(chartRef, getOption);

const loadData = async () => {
  try {
    const res = await Dispatch("dashboard_manager:get_finance_chart", JSON.stringify({ type: range.value }));
    const response = JSON.parse(res) as ResponseWrapper<GetFinanceChartDataResponse>;

    if (response.code === 200 && response.data) {

      const data = response.data;
      chartData.value = {
        xAxis: data.x_axis || [],
        rechargeData: data.recharge_data || [],
        consumeData: data.consume_data || [],
        netData: data.net_data || []
      };
      refresh();
    } else {
      console.error("Failed to load finance chart data:", response.message);
    }
  } catch (e) {
    console.error("Failed to load finance chart data", e);
  }
};

const handleRangeChange = () => {
  loadData();
};

defineExpose({
  loadData
});

onMounted(() => {
  loadData();
});
</script>

<style scoped>
.chart-box {
  width: 100%;
  height: 380px;
}
</style>