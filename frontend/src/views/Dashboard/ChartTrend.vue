<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="d-flex justify-space-between align-center mb-2">
      <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
        <v-icon color="primary" class="mr-2">mdi-chart-bar</v-icon> 月度消课趋势
      </h3>
    </div>
    <div ref="chartRef" class="chart-box"></div>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import * as echarts from 'echarts';
import { useChart } from '../../composables/useChart';

const chartRef = ref<HTMLElement | null>(null);

const getOption = (isDark: boolean) => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  grid: { top: 30, left: 40, right: 20, bottom: 20, containLabel: true },
  xAxis: { type: 'category', data: ['7月', '8月', '9月', '10月', '11月', '12月'], axisTick: { show: false } },
  yAxis: { type: 'value', splitLine: { lineStyle: { type: 'dashed', opacity: 0.3 } } },
  series: [{
    name: '消课总数', type: 'bar', barWidth: 20,
    data: [120, 132, 101, 134, 290, 345],
    itemStyle: { 
      borderRadius: [4, 4, 0, 0], 
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: '#42A5F5' }, 
        { offset: 1, color: '#1976D2' }
      ]) 
    }
  }]
});

useChart(chartRef, getOption);
</script>

<style scoped>
.chart-box { width: 100%; height: 320px; }
</style>