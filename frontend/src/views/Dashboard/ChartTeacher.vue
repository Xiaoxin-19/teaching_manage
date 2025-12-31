<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="d-flex justify-space-between align-center mb-2">
      <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
        <v-icon color="purple" class="mr-2">mdi-medal-outline</v-icon> 本月教师消课排行 (Top 5)
      </h3>
    </div>
    <div ref="chartRef" class="chart-box-sm"></div>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useChart } from '../../composables/useChart';

const chartRef = ref<HTMLElement | null>(null);

const getOption = (isDark: boolean) => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  grid: { top: 10, left: 10, right: 30, bottom: 10, containLabel: true },
  xAxis: { type: 'value', show: false },
  yAxis: { 
    type: 'category', 
    data: ['王老师', '赵老师', '钱老师', '李老师', '张老师'], 
    axisLine: { show: false }, axisTick: { show: false } 
  },
  series: [{
    type: 'bar', data: [45, 52, 68, 70, 85], barWidth: 16, label: { show: true, position: 'right' },
    itemStyle: { 
      borderRadius: 4, 
      color: (params: any) => {
        const colors = ['#CE93D8', '#BA68C8', '#AB47BC', '#9C27B0', '#7B1FA2'];
        return colors[params.dataIndex] || '#7B1FA2';
      }
    }
  }]
});

useChart(chartRef, getOption);
</script>

<style scoped>
.chart-box-sm { width: 100%; height: 280px; }
</style>