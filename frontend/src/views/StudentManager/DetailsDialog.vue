<template>
  <v-dialog :model-value="modelValue" @update:model-value="updateModelValue" max-width="900px" scrollable
    transition="dialog-bottom-transition">
    <v-card class="rounded-lg elevation-4">
      <v-card-title class="d-flex justify-space-between align-center py-3 px-4">
        <div class="d-flex align-center text-subtitle-1 font-weight-bold">
          <v-icon icon="mdi-file-document-outline" color="primary" class="mr-2"></v-icon>
          课时明细
          <span class="text-caption text-medium-emphasis ml-2 font-weight-normal">({{ studentName }})</span>
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" density="comfortable" @click="close"></v-btn>
      </v-card-title>

      <v-divider></v-divider>

      <v-card-text class="pa-0">
        <v-data-table-server :headers="headers" :items="records" :items-length="totalItems" :page="page"
          :items-per-page="itemsPerPage" :loading="loading" @update:options="loadItems" density="compact" hover
          fixed-header height="400px" class="elevation-0">
          <template v-slot:item.type="{ item }">
            <v-chip :color="getTypeColor(item.type)" size="x-small" label class="font-weight-bold" variant="tonal">
              {{ getTypeLabel(item.type) }}
            </v-chip>
          </template>

          <template v-slot:item.tags="{ item }">
            <div class="d-flex gap-1 flex-wrap" style="gap: 4px">
              <v-chip v-for="(tag, index) in item.tags" :key="index" :color="tag.color" size="x-small" label
                variant="outlined" class="font-weight-medium">
                {{ tag.label }}
              </v-chip>
            </div>
          </template>

          <template v-slot:item.amount="{ item }">
            <span :class="item.amount > 0 ? 'text-success' : 'text-error'" class="font-weight-bold">
              {{ item.amount > 0 ? '+' : '' }}{{ item.amount }}
            </span>
          </template>

          <template v-slot:no-data>
            <div class="pa-4 text-center text-medium-emphasis">暂无记录</div>
          </template>
        </v-data-table-server>
      </v-card-text>

      <v-divider></v-divider>

      <v-card-actions class="pa-3">
        <v-spacer></v-spacer>
        <v-btn variant="text" class="mr-2" @click="close">
          关闭
        </v-btn>
        <v-btn prepend-icon="mdi-microsoft-excel" color="success" size="small" variant="outlined"
          @click="() => onExport(studentId)">
          导出 Excel
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import type { RecordItem } from '../../types/appModels'
import { useDetailsDialog } from './DetailsDialog.logic'

const props = defineProps<{
  modelValue: boolean
  studentId?: number
  studentName?: string
}>()

const emit = defineEmits(['update:modelValue', 'export', 'loaded'])

const {
  studentId, studentName, headers, records, page, itemsPerPage, totalItems, loading,
  close, updateModelValue, onExport, getTypeColor, getTypeLabel, load, loadItems
} = useDetailsDialog(emit)

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      studentId.value = Number(props.studentId)
      load(props.studentId, props.studentName).then(() => emit('loaded'))
    }
  },
)
</script>

<style scoped>
:deep(.v-data-table--density-compact .v-data-table-header__content) {
  font-weight: bold;
}
</style>