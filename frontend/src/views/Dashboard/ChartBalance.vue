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
      <div ref="chartRef" style="flex: 1; height: 280px;"></div>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useChart } from '../../composables/useChart';

defineEmits(['navigate']);
const chartRef = ref<HTMLElement | null>(null);

const getOption = (isDark: boolean) => ({
  tooltip: { trigger: 'item' },
  legend: { top: 'middle', left: 'right', orient: 'vertical' },
  series: [{
    name: '课时', type: 'pie', radius: ['50%', '80%'], center: ['35%', '50%'], avoidLabelOverlap: false,
    label: { show: false, position: 'center' },
    emphasis: { label: { show: true, fontSize: 18, fontWeight: 'bold', formatter: '{b}\n{c}人' } },
    data: [
      { value: 110, name: '充足', itemStyle: { color: '#66BB6A' } },
      { value: 15, name: '预警', itemStyle: { color: '#FFA726' } },
      { value: 3, name: '欠费', itemStyle: { color: '#EF5350' } }
    ]
  }]
});

useChart(chartRef, getOption);
</script>