<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="mb-2">
      <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
        <v-icon color="orange" class="mr-2">mdi-grid</v-icon> 高峰时段分布
      </h3>
      <div class="text-caption text-medium-emphasis pl-8">颜色越深代表排课越密集</div>
    </div>
    <div ref="chartRef" class="chart-box"></div>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useChart } from '../../composables/useChart';

const chartRef = ref<HTMLElement | null>(null);

const getOption = (isDark: boolean) => {
  const hours = ['08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00'];
  const days = ['周一', '周二', '周三', '周四', '周五', '周六', '周日'];
  const data = [[0, 0, 1], [0, 1, 0], [0, 2, 0], [0, 3, 0], [0, 4, 2], [0, 5, 8], [0, 6, 5], [1, 0, 0], [1, 1, 0], [1, 2, 0], [1, 3, 0], [1, 4, 3], [1, 5, 9], [1, 6, 6], [2, 0, 0], [2, 1, 0], [2, 2, 0], [2, 3, 0], [2, 4, 2], [2, 5, 8], [2, 6, 4], [3, 0, 0], [3, 1, 0], [3, 2, 0], [3, 3, 0], [3, 4, 4], [3, 5, 9], [3, 6, 6], [4, 0, 0], [4, 1, 0], [4, 2, 0], [4, 3, 0], [4, 4, 5], [4, 5, 10], [4, 6, 8], [5, 0, 8], [5, 1, 12], [5, 2, 6], [5, 3, 10], [5, 4, 12], [5, 5, 6], [5, 6, 2], [6, 0, 6], [6, 1, 10], [6, 2, 5], [6, 3, 8], [6, 4, 10], [6, 5, 4], [6, 6, 1]];

  return {
    tooltip: { position: 'top' },
    grid: { height: '70%', top: '10%', bottom: '10%' },
    xAxis: { type: 'category', data: hours, splitArea: { show: true } },
    yAxis: { type: 'category', data: days, splitArea: { show: true } },
    visualMap: { show: false, min: 0, max: 12, inRange: { color: ['#FFF3E0', '#FF9800', '#BF360C'] } },
    series: [{
      name: '排课密度',
      type: 'heatmap',
      data: data.map(item => [item[1], item[0], item[2] || '-']),
      label: { show: true },
      itemStyle: { emphasis: { shadowBlur: 10, shadowColor: 'rgba(0, 0, 0, 0.5)' } }
    }]
  };
};

useChart(chartRef, getOption);
</script>

<style scoped>
.chart-box {
  width: 100%;
  height: 320px;
}
</style>