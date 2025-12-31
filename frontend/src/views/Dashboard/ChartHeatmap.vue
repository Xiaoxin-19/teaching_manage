<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="mb-2">
      <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
        <v-icon color="orange" class="mr-2">mdi-grid</v-icon>
        高峰时段热力分布
      </h3>
      <div class="text-caption text-medium-emphasis pl-8">
        颜色越深代表排课越密集，建议错峰排课
      </div>
    </div>
    <div ref="chartRef" class="chart-box"></div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useChart } from '../../composables/useChart';
import { Dispatch } from '../../../wailsjs/go/main/App';
import { ResponseWrapper } from '../../types/appModels';

const chartRef = ref<HTMLElement | null>(null);
const heatmapData = ref<[number, number, number][]>([]);

const getOption = (isDark: boolean) => {
  // 定义核心营业时段：08:00 - 22:00
  const hours = [
    '08:00', '09:00', '10:00', '11:00', '12:00',
    '13:00', '14:00', '15:00', '16:00', '17:00',
    '18:00', '19:00', '20:00', '21:00'
  ];

  // 星期定义 (从下往上: 周日 -> 周一，或者符合阅读习惯 从上往下: 周一 -> 周日)
  // ECharts Heatmap Y轴默认是从下往上的，所以数据里的 Y=0 对应 Y轴最底部
  // 我们希望周一在最上面，所以 Y轴数据应该是 ['周日', '周六', ..., '周一']
  const days = ['周日', '周六', '周五', '周四', '周三', '周二', '周一'];

  const data = heatmapData.value;
  const maxVal = data.length > 0 ? Math.max(...data.map(item => item[2])) : 10;

  return {
    tooltip: {
      position: 'top',
      formatter: (params: any) => {
        return `${params.name}<br />${days[params.value[1]]} ${hours[params.value[0]]}: <b>${params.value[2]}节课</b>`;
      }
    },
    grid: {
      top: '10%',
      bottom: '25%',
      left: '12%', // 留出空间显示星期
      right: '2%'
    },
    xAxis: {
      type: 'category',
      data: hours,
      splitArea: { show: true },
      axisLabel: { interval: 0, rotate: 45, fontSize: 10 } // 旋转标签防止重叠
    },
    yAxis: {
      type: 'category',
      data: days,
      splitArea: { show: true }
    },
    visualMap: {
      min: 0,
      max: maxVal || 10,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '0%',
      inRange: {
        // 使用更直观的热力颜色：浅黄 -> 橙 -> 红
        color: isDark
          ? ['#424242', '#FFB74D', '#E65100']
          : ['#FFF3E0', '#FF9800', '#BF360C']
      },
      textStyle: {
        color: isDark ? '#eee' : '#333'
      }
    },
    series: [{
      name: '排课密度',
      type: 'heatmap',
      data: data,
      label: { show: false }, // 格子太小，隐藏数字
      itemStyle: {
        emphasis: {
          shadowBlur: 10,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        },
        borderColor: isDark ? '#333' : '#fff',
        borderWidth: 1
      }
    }]
  };
};

const { refresh } = useChart(chartRef, getOption);

const loadData = async () => {
  try {
    const res = await Dispatch("dashboard_manager:get_heatmap", "");
    const response = JSON.parse(res) as ResponseWrapper<number[][]>;

    if (response.code === 200 && response.data) {

      heatmapData.value = response.data.map((item: number[]) => {
        const dayOfWeek = item[0]; // 0(Sun) - 6(Sat)
        const hour = item[1];      // 8 - 21
        const value = item[2];

        // 1. 映射 X 轴 (Hour): 08:00 对应索引 0
        const xIndex = hour - 8;

        // 2. 映射 Y 轴 (Day): 
        // 前端 days: ['周日', '周六', ..., '周一'] (index 0 是底部)
        // 目标: Sun(0)->0, Mon(1)->6, Sat(6)->1
        const yIndex = (7 - dayOfWeek) % 7;

        return [xIndex, yIndex, value] as [number, number, number];
      });

      refresh();
    }
  } catch (e) {
    console.error("Failed to load heatmap data", e);
  }
};

onMounted(() => {
  loadData();
});

defineExpose({
  loadData
});
</script>

<style scoped>
.chart-box {
  width: 100%;
  height: 400px;
  /* 热力图需要更高一点 */
}
</style>