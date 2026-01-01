<template>
  <v-row dense class="mb-4">
    <!-- 1. 在读学员 -->
    <v-col cols="12" sm="6" md="3">
      <v-card class="rounded-lg elevation-2 h-100 pa-4 kpi-card" hover @click="$emit('navigate', 'students')">
        <div class="d-flex justify-space-between align-start mb-2">
          <div>
            <div class="text-subtitle-2 text-medium-emphasis font-weight-bold mb-1">在读学员</div>
            <div class="text-h4 font-weight-bold text-primary">{{ summaryData.totalStudents }}</div>
          </div>
          <v-avatar color="primary" variant="tonal" rounded size="48">
            <v-icon size="28">mdi-account-school</v-icon>
          </v-avatar>
        </div>
        <div class="d-flex align-center mt-2">
          <v-chip size="small" :color="summaryData.newStudentsThisMonth >= 0 ? 'success' : 'error'" variant="tonal"
            label class="font-weight-bold">
            <v-icon start size="small">{{ summaryData.newStudentsThisMonth >= 0 ? 'mdi-arrow-up' : 'mdi-arrow-down'
              }}</v-icon> {{ summaryData.newStudentsThisMonth }} 本月新增
          </v-chip>
        </div>
      </v-card>
    </v-col>

    <!-- 2. 本月消课 -->
    <v-col cols="12" sm="6" md="3">
      <v-card class="rounded-lg elevation-2 h-100 pa-4 kpi-card" hover @click="$emit('navigate', 'records')">
        <div class="d-flex justify-space-between align-start mb-2">
          <div>
            <div class="text-subtitle-2 text-medium-emphasis font-weight-bold mb-1">本月消课</div>
            <div class="text-h4 font-weight-bold text-success">{{ summaryData.monthlyHours }}</div>
          </div>
          <v-avatar color="success" variant="tonal" rounded size="48">
            <v-icon size="28">mdi-chart-line</v-icon>
          </v-avatar>
        </div>
        <div class="d-flex align-center mt-2">
          <v-chip size="small" :color="summaryData.monthOverMonth.includes('-') ? 'error' : 'success'" variant="tonal"
            label class="font-weight-bold">
            <v-icon start size="small">{{ summaryData.monthOverMonth.includes('-') ? 'mdi-trending-down' :
              'mdi-trending-up' }}</v-icon> {{ summaryData.monthOverMonth }} 环比{{
                summaryData.monthOverMonth.includes('-') ? '下降' : '增长' }}
          </v-chip>
        </div>
      </v-card>
    </v-col>

    <!-- 3. 剩余总课时 -->
    <v-col cols="12" sm="6" md="3">
      <v-card class="rounded-lg elevation-2 h-100 pa-4 kpi-card" hover @click="$emit('navigate', 'students')">
        <div class="d-flex justify-space-between align-start mb-2">
          <div>
            <div class="text-subtitle-2 text-medium-emphasis font-weight-bold mb-1">剩余总课时 (存量)</div>
            <div class="text-h4 font-weight-bold text-teal">{{ summaryData.totalRemainingHours }}</div>
          </div>
          <v-avatar color="teal" variant="tonal" rounded size="48">
            <v-icon size="28">mdi-database-clock</v-icon>
          </v-avatar>
        </div>
        <div class="d-flex align-center mt-2">
          <v-chip size="small" color="teal" variant="tonal" label class="font-weight-bold">
            人均 {{ Math.round(summaryData.totalRemainingHours / (summaryData.totalStudents || 1)) }} 节待修
          </v-chip>
        </div>
      </v-card>
    </v-col>

    <!-- 4. 欠费/预警 -->
    <v-col cols="12" sm="6" md="3">
      <v-card class="rounded-lg elevation-2 h-100 pa-4 kpi-card border-error" hover
        @click="$emit('navigate', 'students')">
        <div class="d-flex justify-space-between align-start mb-2">
          <div>
            <div class="text-subtitle-2 text-medium-emphasis font-weight-bold mb-1">欠费预警</div>
            <div class="text-h4 font-weight-bold text-error">
              {{ warningData.balanceNegative }} <span
                class="text-body-1 text-medium-emphasis font-weight-medium">人</span>
            </div>
          </div>
          <v-avatar color="error" variant="tonal" rounded size="48">
            <v-icon size="28">mdi-bell-ring</v-icon>
          </v-avatar>
        </div>
        <div class="d-flex align-center mt-2">
          <v-chip size="small" color="warning" variant="tonal" label class="font-weight-bold">
            <v-icon start size="small">mdi-alert-circle-outline</v-icon> {{ warningData.balanceLow }} 人余额不足( &lt;5 )
          </v-chip>
        </div>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
defineProps<{
  summaryData: {
    totalStudents: number;
    monthlyHours: number;
    totalRemainingHours: number;
    monthOverMonth: string;
    newStudentsThisMonth: number;
  };
  warningData: {
    balanceLow: number;
    balanceNegative: number;
  };
}>();

defineEmits(['navigate']);
</script>

<style scoped>
.kpi-card {
  min-height: 140px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.border-error {
  border-left: 4px solid rgb(var(--v-theme-error)) !important;
}
</style>