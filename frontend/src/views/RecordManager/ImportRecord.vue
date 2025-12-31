<template>
  <v-dialog :model-value="modelValue" @update:model-value="close" max-width="500px">
    <v-card class="rounded-lg elevation-4">
      <v-card-title class="d-flex justify-space-between align-center py-3 px-4">
        <div class="d-flex align-center text-subtitle-1 font-weight-bold">
          <v-icon icon="mdi-file-excel" color="success" class="mr-2"></v-icon>
          批量导入教学记录
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" @click="close"></v-btn>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text class="pa-6 text-center">
        <!-- 添加 --wails-drop-target 样式变量，标记此区域为拖拽目标 -->
        <div class="upload-zone d-flex flex-column align-center justify-center rounded-lg py-10 mb-4 cursor-pointer"
          style="--wails-drop-target: drop" @click="triggerFileInput">
          <v-avatar :color="selectedFile ? 'success' : 'primary'" variant="tonal" size="64" class="mb-4">
            <v-icon size="32">{{ selectedFile ? 'mdi-file-check' : 'mdi-cloud-upload' }}</v-icon>
          </v-avatar>
          <div class="text-h6 font-weight-bold text-high-emphasis">
            {{ selectedFile ? '已选择文件' : '点击或拖拽 Excel 文件至此' }}
          </div>
          <div class="text-body-2 text-medium-emphasis mt-1">
            {{ selectedFile || '支持 .xlsx, .xls 格式文件' }}
          </div>
        </div>
        <div class="text-caption text-warning text-left px-1 font-weight-medium d-flex align-start">
          <v-icon size="small" color="warning" class="mr-1 mt-1">
            mdi-alert-circle-outline
          </v-icon>
          <span>
            注意：批量导入的记录默认为“未生效”状态，需要手动点击“立即生效”之后才会减扣学生课时。
          </span>
        </div>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions class="pa-3">
        <v-btn variant="outlined" color="info" prepend-icon="mdi-file-download-outline" @click="downloadTemplate">
          下载模板
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn variant="text" class="mr-2" @click="close">取消</v-btn>
        <v-btn color="success" variant="flat" @click="startImport">开始导入</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { useImportRecord } from './ImportRecord.logic'

const props = defineProps<{
  modelValue: boolean;
}>();

// 定义事件，用于通知父组件关闭对话框，以及告知导入结果是否成功，携带导入结果信息
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'import-success'): void;
  (e: 'import-failed', errorInfo: any): void;
}>();

const {
  selectedFile,
  close,
  triggerFileInput,
  downloadTemplate,
  startImport
} = useImportRecord(props, emit)
</script>

<style scoped>
.upload-zone {
  border: 2px dashed rgba(var(--v-border-color), var(--v-border-opacity));
  transition: all 0.3s ease;
  background-color: rgba(var(--v-theme-on-surface), 0.02);
}

.upload-zone:hover,
.upload-zone.wails-drop-target-active {
  border-color: rgb(var(--v-theme-primary));
  background-color: rgba(var(--v-theme-primary), 0.08);
  transform: scale(1.01);
}
</style>