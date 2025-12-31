<template>
  <v-sheet class="fill-height pa-6 bg-background overflow-y-auto">
    <div class="d-flex flex-column h-100">
      <!-- 顶部: 核心KPI -->
      <div class="mb-6 d-flex align-center justify-space-between">
        <div>
          <h1 class="text-h5 font-weight-bold text-high-emphasis">教务运营驾驶舱</h1>
          <p class="text-subtitle-2 text-medium-emphasis mt-1">数据截止: {{ currentDate }}</p>
        </div>
        <v-btn color="primary" variant="flat" prepend-icon="mdi-refresh" @click="handleRefresh" :loading="loading">
          刷新数据
        </v-btn>
      </div>

      <div v-if="loading && firstLoad" class="d-flex justify-center align-center" style="height: 400px;">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
      </div>

      <div v-else class="pb-6">
        <KpiCards :summaryData="summaryData" :warningData="warningData" @navigate="navigateTo" />

        <v-row dense class="mb-4">
          <v-col cols="12" md="6">
            <ChartFinance ref="chartFinanceRef" />
          </v-col>
          <v-col cols="12" md="6">
            <ChartHeatmap />
          </v-col>
        </v-row>

        <v-row dense class="mb-4">
          <v-col cols="12" md="6">
            <ChartEngagement />
          </v-col>
          <v-col cols="12" md="6">
            <ChartGrowth />
          </v-col>
        </v-row>

        <v-row dense>
          <v-col cols="12" md="6">
            <ChartTeacher />
          </v-col>
          <v-col cols="12" md="6">
            <ChartBalance @navigate="navigateTo" />
          </v-col>
        </v-row>
      </div>
    </div>
  </v-sheet>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useDashboard } from './Dashboard.logic';
import KpiCards from './KpiCards.vue';
import ChartFinance from './ChartFinance.vue';
import ChartTrend from './ChartTrend.vue';
import ChartGrowth from './ChartGrowth.vue';
import ChartTeacher from './ChartTeacher.vue';
import ChartBalance from './ChartBalance.vue';
import ChartHeatmap from './ChartHeatmap.vue';
import ChartEngagement from './ChartEngagement.vue';

const {
  loading,
  summaryData,
  warningData,
  currentDate,
  loadDashboardData,
  navigateTo
} = useDashboard();

const chartFinanceRef = ref();
const firstLoad = ref(true);

// 监听 loading 状态，第一次加载完成后将 firstLoad 置为 false
watch(loading, (newVal) => {
  if (!newVal) {
    firstLoad.value = false;
  }
});

const handleRefresh = async () => {
  // 刷新核心数据
  const dashboardPromise = loadDashboardData();

  // 刷新图表数据 (如果组件已挂载)
  if (chartFinanceRef.value) {
    chartFinanceRef.value.loadData();
  }

  await dashboardPromise;
};
</script>