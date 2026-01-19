<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="d-flex justify-space-between align-center mb-2">
      <div>
        <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
          <v-icon color="indigo" class="mr-2">mdi-scale-balance</v-icon>
          课时流转趋势 (最近6个月)
        </h3>
        <div class="text-caption text-medium-emphasis mt-1">主要显示课时，可选显示金额</div>
      </div>
      <div class="d-flex gap-2 align-center">
        <v-btn-toggle v-model="range" density="compact" variant="outlined" color="indigo" mandatory
          @update:model-value="handleRangeChange">
          <v-btn value="1m" size="small">近一月</v-btn>
          <v-btn value="6m" size="small">近半年</v-btn>
          <v-btn value="12m" size="small">近一年</v-btn>
          <v-btn value="all" size="small">全部</v-btn>
        </v-btn-toggle>
        <v-btn size="small" variant="text" icon @click="toggleShowAmount" :title="showAmount ? '隐藏金额' : '显示金额'">
          <v-icon>{{ showAmount ? 'mdi-eye' : 'mdi-eye-off' }}</v-icon>
        </v-btn>
      </div>
    </div>
    <div ref="chartRef" class="chart-box"></div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import * as echarts from 'echarts';
import { useChart } from '../../composables/useChart';
import { GetFinanceChartData } from '../../api/dashboard';
import type { GetFinanceChartRequest } from '../../types/request';

const range = ref<GetFinanceChartRequest['type']>('6m');
const showAmount = ref(false);
const chartRef = ref<HTMLElement | null>(null);

// 定义数据状态
const chartData = ref({
  xAxis: [] as string[],
  rechargeData: [] as number[],
  rechargeAmount: [] as number[],
  consumeData: [] as number[],
  consumeAmount: [] as number[]
});

const getOption = (isDark: boolean) => {
  const { xAxis, rechargeData, rechargeAmount, consumeData, consumeAmount } = chartData.value;

  const baseOption: any = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' },
      formatter: (params: any) => {
        if (!Array.isArray(params)) return ''

        let html = `<div style="padding: 8px"><strong>${params[0]?.axisValue}</strong><br>`
        params.forEach((param: any) => {
          const value = param.value
          // 如果是金额系列，显示为人民币格式
          const displayValue = param.seriesName.includes('金额')
            ? `¥${(value || 0).toFixed(2)}`
            : `${value || 0}`
          html += `<span style="color: ${param.color}">● ${param.seriesName}: ${displayValue}</span><br>`
        })
        html += '</div>'
        return html
      }
    },
    legend: {
      data: ['充值课时', '消课课时', ...(showAmount.value ? ['充值金额', '消课金额'] : [])],
      top: 0
    },
    grid: {
      top: 40,
      left: 20,
      right: showAmount.value ? 60 : 20,
      bottom: 60,
      containLabel: true
    },
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
    xAxis: {
      type: 'category',
      data: xAxis,
      axisLine: { show: false },
      axisTick: { show: false }
    },
    yAxis: [
      {
        type: 'value',
        name: '课时数',
        position: 'left',
        splitLine: { lineStyle: { type: 'dashed', opacity: 0.3 } }
      },
      ...(showAmount.value ? [{
        type: 'value',
        name: '金额 (¥)',
        position: 'right',
        splitLine: { show: false }
      }] : [])
    ],
    series: [
      {
        name: '充值课时',
        type: 'bar',
        barGap: 0,
        itemStyle: { color: '#5C6BC0', borderRadius: [4, 4, 0, 0] },
        yAxisIndex: 0,
        data: rechargeData
      },
      {
        name: '消课课时',
        type: 'bar',
        itemStyle: { color: '#66BB6A', borderRadius: [4, 4, 0, 0] },
        yAxisIndex: 0,
        data: consumeData
      },
      ...(showAmount.value ? [
        {
          name: '充值金额',
          type: 'line',
          smooth: true,
          yAxisIndex: 1,
          itemStyle: { color: '#FF9800' },
          lineStyle: { type: 'dashed', width: 2 },
          data: rechargeAmount
        },
        {
          name: '消课金额',
          type: 'line',
          smooth: true,
          yAxisIndex: 1,
          itemStyle: { color: '#F44336' },
          lineStyle: { type: 'dashed', width: 2 },
          data: consumeAmount
        }
      ] : [])
    ]
  };

  return baseOption;
};

const { refresh } = useChart(chartRef, getOption);

const loadData = async () => {
  try {
    const data = await GetFinanceChartData(range.value);
    chartData.value = {
      xAxis: data.x_axis || [],
      rechargeData: data.recharge_data || [],
      rechargeAmount: data.recharge_amount || [],
      consumeData: data.consume_data || [],
      consumeAmount: data.consume_amount || []
    };
    refresh();
  } catch (e) {
    console.error("Failed to load finance chart data", e);
  }
};

const handleRangeChange = () => {
  loadData();
};

const toggleShowAmount = () => {
  showAmount.value = !showAmount.value;
  refresh();
};

defineExpose({ loadData });

onMounted(() => {
  loadData();
});
</script>

<style scoped>
.chart-box {
  width: 100%;
  height: var(--dashboard-chart-height, 380px);
}

.gap-2 {
  gap: 8px;
}
</style>