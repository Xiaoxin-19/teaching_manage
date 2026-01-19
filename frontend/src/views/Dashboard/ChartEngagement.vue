<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="d-flex justify-space-between align-center mb-2">
      <div>
        <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
          <v-icon color="deep-purple-accent-2" class="mr-2">mdi-account-heart</v-icon>
          学员活跃度分布 (近30天)
        </h3>
        <div class="text-caption text-medium-emphasis mt-1">
          基于平均月课次/科统计：<span class="text-warning font-weight-bold">&lt;1/科为缺勤风险</span>，<span
            class="text-info font-weight-bold">1-3/科为达标</span>
        </div>
      </div>
    </div>
    <div ref="chartRef" class="chart-box"></div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import * as echarts from 'echarts';
import { useChart } from '../../composables/useChart';
import { GetStudentEngagementData } from '../../api/dashboard';
import type { EngagementStat } from '../../types/response';

// 前端本地定义的活跃度分类
const ENGAGEMENT_CATEGORIES = {
  Dormant: { name: '沉睡 (0/科)', desc: '需激活' },
  Lazy: { name: '消极 (<1/科)', desc: '缺勤风险' },
  Regular: { name: '达标 (1-3/科)', desc: '每周1-2练' },
  High: { name: '高频 (>3/科)', desc: '集训/多科' }
};

const chartRef = ref<HTMLElement | null>(null);
const engagementData = ref([
  { code: 'Dormant', value: 0, name: '沉睡 (0/科)', desc: '需激活' },
  { code: 'Lazy', value: 0, name: '消极 (<1/科)', desc: '缺勤风险' },
  { code: 'Regular', value: 0, name: '达标 (1-3/科)', desc: '每周1-2练' },
  { code: 'High', value: 0, name: '高频 (>3/科)', desc: '集训/多科' }
]);

const getOption = (isDark: boolean) => {
  const data = engagementData.value;

  const colors = [
    isDark ? '#616161' : '#9E9E9E', // 灰：沉睡
    '#FF9800', // 橙：消极 (警告色)
    '#2196F3', // 蓝：达标 (正常色)
    '#4CAF50'  // 绿：高频 (优质色)
  ];

  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      formatter: (params: any) => {
        const item = params[0];
        const dataItem = data[item.dataIndex];
        return `
          <div style="margin-bottom: 4px;"><b>${item.name}</b></div>
          <div style="font-size: 12px; color: ${item.color}">${dataItem.desc}</div>
          <div style="margin-top: 4px; font-size: 14px;">
            <span style="display:inline-block;margin-right:5px;border-radius:10px;width:10px;height:10px;background-color:${item.color};"></span>
            <b>${item.value} 人</b>
          </div>
        `;
      }
    },
    grid: {
      top: '15%',
      left: '3%',
      right: '4%',
      bottom: '8%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(i => i.name),
      axisTick: { show: false },
      axisLine: { show: false },
      axisLabel: {
        color: isDark ? '#eee' : '#666',
        interval: 0,
        fontSize: 11,
        formatter: (value: string) => {
          // 标签换行显示，避免拥挤
          return value.replace(' (', '\n(');
        }
      }
    },
    yAxis: {
      type: 'value',
      name: '人数',
      splitLine: {
        lineStyle: {
          type: 'dashed',
          opacity: 0.3
        }
      }
    },
    series: [
      {
        name: '学员人数',
        type: 'bar',
        barWidth: '40%',
        data: data.map((item, index) => ({
          value: item.value,
          itemStyle: {
            color: colors[index],
            borderRadius: [6, 6, 0, 0]
          }
        })),
        label: {
          show: true,
          position: 'top',
          color: isDark ? '#fff' : '#666',
          fontWeight: 'bold',
          formatter: '{c}'
        },
        // 添加背景条，让柱状图看起来更现代
        showBackground: true,
        backgroundStyle: {
          color: isDark ? 'rgba(255, 255, 255, 0.05)' : 'rgba(0, 0, 0, 0.03)',
          borderRadius: [6, 6, 0, 0]
        }
      }
    ]
  };
};

const { refresh } = useChart(chartRef, getOption);

const loadData = async () => {
  try {
    const data = await GetStudentEngagementData();
    if (data && data.stats) {
      // 使用 code 进行映射，而不是依赖名称字符串
      // 这样即使后端改变显示名称，映射也不会破裂
      const statsMap = new Map(data.stats.map((s) => [s.code, s.value]));
      engagementData.value = engagementData.value.map((item) => ({
        ...item,
        value: statsMap.get(item.code) || 0
      }));
      refresh();
    }
  } catch (e) {
    console.error('Failed to load engagement data', e);
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
  height: var(--dashboard-chart-height, 380px);
}
</style>