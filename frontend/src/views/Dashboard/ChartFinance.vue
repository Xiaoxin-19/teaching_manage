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
import { ref } from 'vue';
import * as echarts from 'echarts';
import { useChart } from '../../composables/useChart';

const range = ref('6m');
const chartRef = ref<HTMLElement | null>(null);

const getOption = (isDark: boolean) => {
  let xAxisData: string[] = [], rechargeData: number[] = [], consumeData: number[] = [], netData: number[] = [];

  // TODO: 这里应替换为真实的后端 API 数据
  if (range.value === '1m') {
    xAxisData = Array.from({ length: 30 }, (_, i) => `${i + 1}日`);
    rechargeData = Array.from({ length: 30 }, () => Math.floor(Math.random() * 40) + 5);
    consumeData = Array.from({ length: 30 }, () => Math.floor(Math.random() * 30) + 10);
  } else if (range.value === '6m') {
    xAxisData = ['7月', '8月', '9月', '10月', '11月', '12月'];
    rechargeData = [320, 150, 400, 200, 550, 280];
    consumeData = [220, 232, 201, 234, 290, 345];
  } else if (range.value === '12m') {
    xAxisData = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'];
    rechargeData = [300, 280, 450, 320, 150, 400, 200, 550, 280, 310, 420, 500];
    consumeData = [200, 180, 350, 220, 232, 201, 234, 290, 345, 300, 320, 380];
  } else {
    xAxisData = ['2023-Q1', '2023-Q2', '2023-Q3', '2023-Q4', '2024-Q1', '2024-Q2', '2024-Q3', '2024-Q4'];
    rechargeData = [1200, 1500, 1100, 1800, 1600, 1400, 1700, 1900];
    consumeData = [1000, 1200, 1300, 1400, 1500, 1550, 1600, 1750];
  }
  netData = rechargeData.map((val, i) => val - consumeData[i]);

  return {
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { data: ['充值课时 (收入)', '消课课时 (营收)', '净增库存'], bottom: 0 },
    grid: { top: 30, left: 40, right: 40, bottom: 30, containLabel: true },
    xAxis: { type: 'category', data: xAxisData, axisLine: { show: false }, axisTick: { show: false } },
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

const handleRangeChange = () => {
  refresh();
};
</script>

<style scoped>
.chart-box {
  width: 100%;
  height: 320px;
}
</style>