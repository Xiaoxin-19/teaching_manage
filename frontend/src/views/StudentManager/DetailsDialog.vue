<template>
  <v-dialog :model-value="modelValue" @update:model-value="updateModelValue" max-width="800px" scrollable
    transition="dialog-bottom-transition">
    <v-card class="rounded-lg elevation-4">
      <!-- 1. 弹窗头部 (极简风格) -->
      <v-card-title class="d-flex justify-space-between align-center py-2 px-4 bg-white border-b">
        <div class="d-flex align-center text-subtitle-2 font-weight-bold text-grey-darken-3">
          <!-- 图标替代姓名 -->
          <v-icon icon="mdi-file-document-outline" color="primary" class="mr-2"></v-icon>
          课时明细
          <!-- 姓名作为辅助信息 -->
          <span class="text-caption text-medium-emphasis ml-1 font-weight-normal">({{ studentName }})</span>
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" color="grey" @click="close"></v-btn>
      </v-card-title>

      <!-- 2. 课时变动记录表格 -->
      <v-card-text class="pa-0">
        <v-data-table :headers="headers" :items="records" density="compact" hover fixed-header height="400px"
          class="elevation-0">
          <!-- 类型列：带颜色的Chip -->
          <template v-slot:item.type="{ item }">
            <v-chip :color="getTypeColor(item.type)" size="x-small" label class="font-weight-bold">
              {{ item.type }}
            </v-chip>
          </template>

          <!-- 变动金额列：带颜色 -->
          <template v-slot:item.amount="{ item }">
            <span :class="item.amount > 0 ? 'text-success' : 'text-error'" class="font-weight-bold">
              {{ item.amount > 0 ? '+' : '' }}{{ item.amount }}
            </span>
          </template>

          <!-- 结余列：加粗 -->
          <template v-slot:item.balanceAfter="{ item }">
            <span class="font-weight-medium">{{ item.balanceAfter }}</span>
          </template>

          <!-- 无数据插槽 -->
          <template v-slot:no-data>
            <div class="pa-4 text-center text-grey">暂无记录</div>
          </template>
        </v-data-table>
      </v-card-text>

      <!-- 3. 底部操作栏 -->
      <v-card-actions class="pa-2 bg-grey-lighten-4">
        <v-spacer></v-spacer>
        <v-btn color="grey-darken-1" size="small" variant="text" class="mr-2" @click="close">
          关闭
        </v-btn>
        <v-btn prepend-icon="mdi-export-variant" color="success" size="small" variant="elevated" elevation="1"
          @click="onExport">
          导出 Excel
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import type { RecordItem, FetchDetailsFn } from '../../types/appModels'
import { headers, useDetailsDialog } from './DetailsDialog.logic'

// === Props 定义 ===
const props = defineProps<{
  modelValue: boolean                // 控制显示/隐藏 (v-model)
  studentId?: number                // 学生 ID（可选）
  // 可选：父组件直接传入用于回退显示的数据
  studentName?: string
  records?: RecordItem[]
  // 可选：父组件注入实际的后端调用函数
  fetchDetails?: FetchDetailsFn
}>()

// === Emits 定义 ===
const emit = defineEmits(['update:modelValue', 'export', 'loaded'])

// 使用逻辑层，传入可能的 fetcher
const { studentName, records, close, updateModelValue, onExport, getTypeColor, load } = useDetailsDialog(emit as unknown as (event: string, ...args: any[]) => void, props.fetchDetails)

// 当 dialog 被打开时，尝试加载数据（优先使用 fetcher，否则使用 props 中的回退数据）
watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      load(props.studentId, props.studentName, props.records).then(() => emit('loaded'))
    }
  },
)
</script>

<style scoped>
.bg-grey-lighten-4 {
  background-color: #f5f5f5 !important;
}

/* 强制让表格表头加粗，提升可读性 */
:deep(.v-data-table--density-compact .v-data-table-header__content) {
  font-weight: bold;
}
</style>