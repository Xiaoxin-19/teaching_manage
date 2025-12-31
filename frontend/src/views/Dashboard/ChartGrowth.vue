<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="mb-2">
      <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
        <v-icon color="success" class="mr-2">mdi-trending-up</v-icon> 学员增长趋势
      </h3>
      <div class="text-caption text-medium-emphasis pl-8">每月新增报名人数</div>
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
  tooltip: { trigger: 'axis' },
  grid: { top: 30, left: 30, right: 20, bottom: 20, containLabel: true },
  xAxis: { type: 'category', data: ['7月', '8月', '9月', '10月', '11月', '12月'], boundaryGap: false },
  yAxis: { type: 'value', splitLine: { show: false } },
  series: [{
    name: '新增学员', type: 'line', smooth: true,
    data: [5, 8, 12, 7, 15, 20],
    areaStyle: { 
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: 'rgba(76, 175, 80, 0.5)' }, 
        { offset: 1, color: 'rgba(76, 175, 80, 0.0)' }
      ]) 
    },
    itemStyle: { color: '#4CAF50' },
    lineStyle: { width: 3 }
  }]
});

useChart(chartRef, getOption);
</script>

<style scoped>
.chart-box { width: 100%; height: 320px; }
</style>