<template>
  <v-dialog v-model="dialogVisible" max-width="450" :persistent="loading">
    <v-card class="rounded-lg">
      <v-card-title class="d-flex justify-space-between align-center py-3 px-4 bg-surface border-b">
        <div class="d-flex align-center">
          <v-icon :icon="mode === 'charge' ? 'mdi-wallet-plus' : 'mdi-cash-minus'"
            :color="mode === 'charge' ? 'success' : 'warning'" class="mr-2"></v-icon>
          <span class="text-subtitle-1 font-weight-bold">
            {{ mode === 'charge' ? '课时充值' : '退费/扣课时' }}
          </span>
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" @click="dialogVisible = false" :disabled="loading"></v-btn>
      </v-card-title>

      <v-card-text class="pa-6">
        <!-- 顶部信息卡 -->
        <div class="bg-grey-lighten-4 rounded pa-3 mb-4 d-flex justify-space-between align-center">
          <div>
            <div class="text-caption text-medium-emphasis">学员</div>
            <div class="font-weight-bold text-body-2">{{ course.student?.name ?? '未知学员' }}</div>
          </div>
          <div>
            <div class="text-caption text-medium-emphasis">科目</div>
            <div class="font-weight-bold text-body-2">{{ course.subject.name }}</div>
          </div>
          <div class="text-right">
            <div class="text-caption text-medium-emphasis">当前剩余</div>
            <div class="font-weight-bold text-body-2" :class="'text-' + getBalanceColor(course.balance)">
              {{ course.balance }} 节
            </div>
          </div>
        </div>

        <v-form @submit.prevent="handleSubmit">
          <!-- 课时输入 -->
          <v-text-field v-model="form.hours" :label="mode === 'charge' ? '充值课时' : '扣除课时'" type="number"
            variant="outlined" density="comfortable" class="mb-1"
            :prepend-inner-icon="mode === 'charge' ? 'mdi-plus-circle-outline' : 'mdi-minus-circle-outline'" suffix="节"
            autofocus min="0" :error-messages="!isValid && form.hours !== 0 ? ['请输入大于0的课时数'] : []"
            :disabled="loading"></v-text-field>

          <!-- 实际金额 -->
          <v-text-field v-model="form.amount" label="实际金额 (选填)" type="number" variant="outlined" density="comfortable"
            class="mb-1" prepend-inner-icon="mdi-currency-cny" suffix="元" min="0"
            :error-messages="form.amount !== null && form.amount < 0 ? ['金额必须大于等于0'] : []" :disabled="loading"
            placeholder="0.00"></v-text-field>
          <!-- 智能标签 -->
          <div class="d-flex flex-wrap gap-2 mb-2">
            <v-chip v-for="tag in smartTags" :key="tag" size="small" link
              :color="mode === 'charge' ? 'success' : 'warning'" variant="tonal" @click="applyTag(tag)"
              :disabled="loading">
              {{ mode === 'charge' ? '+' : '-' }}{{ tag }}
            </v-chip>
          </div>

          <!-- 备注输入 -->
          <v-textarea v-model="form.remark" label="备注" rows="2" variant="outlined" density="comfortable"
            placeholder="例如：微信支付、活动赠送、打错退费" class="mt-2" hint="输入关键词自动分类" persistent-hint
            :disabled="loading"></v-textarea>

          <!-- 智能分类标签展示区 -->
          <div class="mt-2 d-flex align-center flex-wrap gap-2" style="min-height: 24px">
            <template v-if="inferredTags.length > 0">
              <v-icon size="small" color="primary">mdi-auto-fix</v-icon>
              <v-chip v-for="(tag, i) in inferredTags" :key="tag.label" :color="tag.color" size="x-small" label
                variant="tonal" class="font-weight-bold">
                {{ tag.label }}
              </v-chip>
            </template>
          </div>
        </v-form>
      </v-card-text>

      <v-card-actions class="pa-4 pt-0">
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="dialogVisible = false" :disabled="loading">取消</v-btn>
        <v-btn :color="mode === 'charge' ? 'success' : 'warning'" variant="elevated" :disabled="!isValid"
          :loading="loading" @click="handleSubmit" class="px-6">
          确认{{ mode === 'charge' ? '充值' : '扣除' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { useRechargeDialog } from './RechargeDialog.logic'

const props = defineProps<{
  modelValue: boolean
  course: any
  initialMode?: 'charge' | 'refund'
}>()

const emit = defineEmits(['update:modelValue', 'submit'])

const {
  loading,
  dialogVisible,
  mode,
  form,
  isValid,
  smartTags,
  inferredTags,
  applyTag,
  handleSubmit,
  getBalanceColor,
  onApiSuccess,
  onApiError
} = useRechargeDialog(props, emit)

// 暴露方法给父组件
defineExpose({
  onApiSuccess,
  onApiError
})
</script>