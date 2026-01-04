<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="d-flex justify-space-between align-center mb-2">
      <div>
        <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
          <v-icon color="blue-accent-2" class="mr-2">mdi-trending-up</v-icon>
          学员增长和流失趋势 (近6个月)
        </h3>
        <div class="text-caption text-medium-emphasis mt-1">
          绿色为新增学员，红色为流失学员，蓝色为净增人数
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
import { GetStudentGrowthTrendData } from '../../api/dashboard';
import type { StudentGrowthTrendResponse } from '../../types/response';

const chartRef = ref<HTMLElement | null>(null);
const chartData = ref<StudentGrowthTrendResponse>({
  x_axis: [],
  growth: [],
  loss: [],
  net: []
});

const getOption = (isDark: boolean) => {
  const data = chartData.value;

  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' },
      formatter: (params: any) => {
        if (!Array.isArray(params)) {
          params = [params];
        }
        let result = `<div>${params[0].axisValue}</div>`;
        params.forEach((param: any) => {
          result += `<div><span style="color:${param.color}">●</span> ${param.seriesName}: <strong>${param.value}</strong>人</div>`;
        });
        return result;
      }
    },
    legend: {
      data: ['新增学员', '流失学员', '净增学员'],
      top: 'top',
      textStyle: {
        color: isDark ? '#eee' : '#333'
      }
    },
    grid: {
      left: '5%',
      right: '5%',
      bottom: '10%',
      top: '12%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.x_axis,
      boundaryGap: false,
      axisLine: {
        lineStyle: {
          color: isDark ? '#444' : '#ddd'
        }
      },
      axisLabel: {
        color: isDark ? '#aaa' : '#666',
        fontSize: 11
      }
    },
    yAxis: {
      type: 'value',
      name: '人数',
      nameTextStyle: {
        color: isDark ? '#aaa' : '#666'
      },
      splitLine: {
        lineStyle: {
          type: 'dashed',
          color: isDark ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)'
        }
      },
      axisLabel: {
        color: isDark ? '#aaa' : '#666'
      }
    },
    series: [
      {
        name: '新增学员',
        type: 'line',
        data: data.growth,
        smooth: true,
        itemStyle: {
          color: '#4CAF50'
        },
        lineStyle: {
          color: '#4CAF50',
          width: 2
        },
        areaStyle: {
          color: 'rgba(76, 175, 80, 0.2)'
        },
        emphasis: {
          itemStyle: {
            borderColor: '#4CAF50',
            borderWidth: 2
          }
        },
        symbol: 'circle',
        symbolSize: 6
      },
      {
        name: '流失学员',
        type: 'line',
        data: data.loss,
        smooth: true,
        itemStyle: {
          color: '#F44336'
        },
        lineStyle: {
          color: '#F44336',
          width: 2
        },
        areaStyle: {
          color: 'rgba(244, 67, 54, 0.2)'
        },
        emphasis: {
          itemStyle: {
            borderColor: '#F44336',
            borderWidth: 2
          }
        },
        symbol: 'circle',
        symbolSize: 6
      },
      {
        name: '净增学员',
        type: 'bar',
        data: data.net,
        itemStyle: {
          color: '#2196F3',
          borderRadius: [6, 6, 0, 0]
        },
        emphasis: {
          itemStyle: {
            color: '#1976D2'
          }
        },
        barWidth: '40%'
      }
    ]
  };
};

const { refresh } = useChart(chartRef, getOption);

const loadData = async () => {
  try {
    const data = await GetStudentGrowthTrendData();
    if (data) {
      chartData.value = data;
      refresh();
    }
  } catch (e) {
    console.error('Failed to load growth trend data', e);
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