<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="mb-2">
      <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
        <v-icon color="pink" class="mr-2">mdi-piano</v-icon> 热门科目排行
      </h3>
      <div class="text-caption text-medium-emphasis pl-8">本月各科目消课占比</div>
    </div>
    <div ref="chartRef" class="chart-box"></div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useChart } from '../../composables/useChart';
import { GetSubjectRankData } from '../../api/dashboard';
import type { GetSubjectRankDataResponse } from '../../types/response';

const chartRef = ref<HTMLElement | null>(null);
const chartData = ref<GetSubjectRankDataResponse['data']>([]);

// 动态生成颜色函数，确保颜色均匀分布且不重复
const generateColor = (index: number): string => {
  // 使用HSL色彩空间，色相均匀分布 (0-360度)
  // 饱和度65%，亮度55% 确保颜色鲜艳且易辨识
  const hue = (index * 45) % 360;
  const saturation = 65;
  const lightness = 55;
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
};

const loadData = async () => {
  try {
    const resp = await GetSubjectRankData();
    chartData.value = resp.data || [];
  } catch (error) {
    console.error('Failed to load subject rank data:', error);
    chartData.value = [];
  }
};

const seriesData = computed(() => {
  return chartData.value.map((item, index) => ({
    value: item.value,
    name: item.name,
    itemStyle: { color: generateColor(index) }
  }));
});

const getOption = (isDark: boolean) => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} 节 ({d}%)' },
  legend: {
    orient: 'vertical', left: 'right', top: 'center',
    textStyle: { color: isDark ? '#ccc' : '#333' }
  },
  series: [
    {
      name: '科目消课',
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['35%', '50%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 5,
        borderColor: isDark ? '#1e1e1e' : '#fff',
        borderWidth: 2
      },
      label: { show: false, position: 'center' },
      emphasis: {
        itemStyle: {
          borderRadius: 5,
          borderColor: isDark ? '#1e1e1e' : '#fff',
          borderWidth: 2
        },
        label: { show: true, fontSize: 16, fontWeight: 'bold' }
      },
      data: seriesData.value
    }
  ]
});

const { setOption } = useChart(chartRef, getOption);

// 当数据变化时，更新图表
watch(() => seriesData.value, () => {
  if (setOption) {
    const isDark = document.documentElement.classList.contains('dark');
    setOption(getOption(isDark));
  }
}, { deep: true });

// 页面加载时获取数据
loadData();

defineExpose({ loadData });
</script>

<style scoped>
.chart-box {
  width: 100%;
  height: 320px;
}
</style>