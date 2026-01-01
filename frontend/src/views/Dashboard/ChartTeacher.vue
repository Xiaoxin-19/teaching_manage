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
import { ref, onMounted } from 'vue';
import { useChart } from '../../composables/useChart';
import { Dispatch } from '../../../wailsjs/go/main/App';
import { ResponseWrapper } from '../../types/appModels';
import { GetTeacherRankDataResponse } from '../../types/response';

const chartRef = ref<HTMLElement | null>(null);
const chartData = ref<GetTeacherRankDataResponse>({ names: [], values: [] });

const getOption = (isDark: boolean) => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  grid: { top: 10, left: 10, right: 30, bottom: 10, containLabel: true },
  xAxis: { type: 'value', show: false },
  yAxis: {
    type: 'category',
    data: chartData.value.names,
    axisLine: { show: false }, axisTick: { show: false },
    axisLabel: { color: isDark ? '#eee' : '#666' }
  },
  series: [{
    type: 'bar',
    data: chartData.value.values,
    barWidth: 16,
    label: {
      show: true,
      position: 'right',
      color: isDark ? '#eee' : '#666'
    },
    itemStyle: {
      borderRadius: 4,
      color: (params: any) => {
        // 颜色从浅到深，对应 Rank 5 -> Rank 1
        const colors = ['#CE93D8', '#BA68C8', '#AB47BC', '#9C27B0', '#7B1FA2'];
        // 如果数据少于5个，确保颜色也能正确映射 (取最后几个颜色)
        const offset = 5 - chartData.value.values.length;
        const index = params.dataIndex + (offset > 0 ? offset : 0);
        return colors[index] || '#7B1FA2';
      }
    }
  }]
});

const { refresh } = useChart(chartRef, getOption);

const loadData = async () => {
  try {
    const res = await Dispatch("dashboard_manager:get_teacher_rank", "");
    const response = JSON.parse(res) as ResponseWrapper<GetTeacherRankDataResponse>;
    if (response.code === 200 && response.data) {
      chartData.value = response.data;
      refresh();
    }
  } catch (e) {
    console.error("Failed to load teacher rank data", e);
  }
};

onMounted(() => {
  loadData();
});

defineExpose({ loadData });
</script>

<style scoped>
.chart-box-sm {
  width: 100%;
  height: 280px;
}
</style>