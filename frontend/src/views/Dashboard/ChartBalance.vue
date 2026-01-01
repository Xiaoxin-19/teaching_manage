<template>
  <v-card class="rounded-lg elevation-2 h-100 pa-4">
    <div class="d-flex justify-space-between align-center mb-2">
      <h3 class="text-subtitle-1 font-weight-bold d-flex align-center">
        <v-icon color="warning" class="mr-2">mdi-chart-pie</v-icon> 学员账户健康度
      </h3>
      <v-btn variant="text" size="small" color="primary" append-icon="mdi-arrow-right"
        @click="$emit('navigate', 'students')">
        催费管理
      </v-btn>
    </div>
    <div class="d-flex">
      <div ref="chartRef" style="flex: 1; height: 340px;"></div>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useChart } from '../../composables/useChart';
import { Dispatch } from '../../../wailsjs/go/main/App';
import { ResponseWrapper } from '../../types/appModels';
import { GetStudentBalanceDataResponse, BalanceStat } from '../../types/response';

defineEmits(['navigate']);
const chartRef = ref<HTMLElement | null>(null);
const chartData = ref<BalanceStat[]>([]);

const getOption = (isDark: boolean) => ({
  tooltip: { trigger: 'item' },
  legend: {
    top: 'middle',
    left: 'right',
    orient: 'vertical',
    textStyle: { color: isDark ? '#eee' : '#333' }
  },
  series: [{
    name: '课时', type: 'pie', radius: ['50%', '80%'], center: ['35%', '50%'], avoidLabelOverlap: false,
    label: { show: false, position: 'center' },
    emphasis: {
      label: {
        show: true,
        fontSize: 18,
        fontWeight: 'bold',
        formatter: '{b}\n{c}人',
        color: isDark ? '#fff' : '#333'
      }
    },
    data: chartData.value.length > 0 ? chartData.value.map(item => {
      let color = '#66BB6A'; // 默认充足
      if (item.name.includes('预警')) color = '#FFA726';
      if (item.name.includes('欠费')) color = '#EF5350';
      return {
        value: item.value,
        name: item.name,
        itemStyle: { color }
      };
    }) : [
      { value: 0, name: '充足', itemStyle: { color: '#66BB6A' } },
      { value: 0, name: '预警', itemStyle: { color: '#FFA726' } },
      { value: 0, name: '欠费', itemStyle: { color: '#EF5350' } }
    ]
  }]
});

const { refresh } = useChart(chartRef, getOption);

const loadData = async () => {
  try {
    const res = await Dispatch("dashboard_manager:get_student_balance", "");
    const response = JSON.parse(res) as ResponseWrapper<GetStudentBalanceDataResponse>;
    if (response.code === 200 && response.data && response.data.stats) {
      chartData.value = response.data.stats;
      refresh();
    }
  } catch (e) {
    console.error("Failed to load balance data", e);
  }
};

onMounted(() => {
  loadData();
});

defineExpose({ loadData });
</script>