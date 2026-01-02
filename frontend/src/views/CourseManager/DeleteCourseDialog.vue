<template>
  <v-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)" max-width="400">
    <v-card class="rounded-lg">
      <v-card-title class="d-flex align-center py-3 px-4 bg-amber-lighten-5">
        <v-icon color="warning" class="mr-2">mdi-archive-remove-outline</v-icon>
        <span class="text-subtitle-1 font-weight-bold text-warning-darken-2">办理退课</span>
      </v-card-title>
      <v-card-text class="pa-4 pt-5">
        <div class="text-body-1 font-weight-bold mb-2">您确定要为 {{ course.studentName }} 办理 {{ course.subjectName }} 退课吗？
        </div>
        <div class="text-body-2 text-medium-emphasis mb-4">此操作将终止课程服务，状态变更为“已结课”。</div>
        <div v-if="(course.balance || 0) > 0" class="bg-amber-lighten-5 pa-3 rounded border border-warning mb-3">
          <div class="text-subtitle-2 font-weight-bold text-warning-darken-2 d-flex align-center mb-1">
            <v-icon size="small" class="mr-1">mdi-cash-multiple</v-icon>检测到剩余余额: {{ course.balance }} 节
          </div>
          <div class="text-caption text-medium-emphasis mt-1">系统将自动结清剩余课时（生成清算记录）。<br />请在备注中说明余额去向（如：退费/转课）。</div>
        </div>
        <v-textarea :model-value="remark" @update:model-value="$emit('update:remark', $event)" label="退课备注 (必填)"
          rows="2" variant="outlined" density="compact" hide-details="auto" placeholder="例如：已退费500元 / 转入钢琴课"
          :rules="[v => !!v || '请填写退课/余额处理说明']"></v-textarea>
      </v-card-text>
      <v-card-actions class="pa-3 bg-grey-lighten-5">
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="$emit('update:modelValue', false)">取消</v-btn>
        <v-btn color="warning" variant="elevated" :disabled="!isValid" @click="$emit('confirm')">确认退课</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { type Course } from '../../api/course'

defineProps<{
  modelValue: boolean
  course: Partial<Course>
  remark: string
  isValid: boolean
}>()

defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'update:remark', value: string): void
  (e: 'confirm'): void
}>()
</script>
