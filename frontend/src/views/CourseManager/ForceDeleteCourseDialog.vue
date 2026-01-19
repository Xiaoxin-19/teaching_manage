<template>
  <v-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)" max-width="400">
    <v-card class="rounded-lg">
      <v-card-title class="d-flex align-center py-3 px-4 bg-error-lighten-5">
        <v-icon color="error" class="mr-2">mdi-alert-decagram</v-icon>
        <span class="text-subtitle-1 font-weight-bold text-error">彻底删除</span>
      </v-card-title>
      <v-card-text class="pa-4 pt-5">
        <div class="text-body-1 font-weight-bold mb-2 text-error">高危操作警告</div>
        <div class="text-body-2 mb-4">您正在尝试彻底删除 <strong>{{ course.student?.name }} - {{ course.subject?.name }}</strong>
          的课程记录。</div>
        <div class="bg-red-lighten-5 pa-3 rounded border border-error text-caption text-error">
          <v-icon size="small" color="error"
            class="mr-1">mdi-alert</v-icon>此操作将<strong>永久移除</strong>该条目，包括相关的上课记录和历史数据，且<strong>无法恢复</strong>！通常仅用于清除错误录入的数据。
        </div>
      </v-card-text>
      <v-card-actions class="pa-3 bg-grey-lighten-5">
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="$emit('update:modelValue', false)">取消</v-btn>
        <v-btn color="error" variant="elevated" @click="$emit('confirm')">确认删除</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { Course } from '../../types/appModels'

defineProps<{
  modelValue: boolean
  course: Partial<Course>
}>()

defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
}>()
</script>
