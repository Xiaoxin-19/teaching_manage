<template>
  <v-dialog
    v-model="show"
    max-width="400px"
    persistent
    transition="dialog-bottom-transition"
  >
    <v-card class="rounded-lg elevation-4">
      <!-- 标题栏：使用 variant="tonal" 的 Avatar 自动适配深色模式 -->
      <v-card-title class="d-flex align-center py-4 px-4">
        <v-avatar :color="typeColor" variant="tonal" size="40" class="mr-3">
          <v-icon :icon="typeIcon" size="24" :color="typeColor"></v-icon>
        </v-avatar>
        <span class="text-h6 font-weight-bold">{{ title }}</span>
      </v-card-title>

      <!-- 内容区：使用 text-medium-emphasis 适配文字颜色，ml-14 保持缩进对齐 -->
      <v-card-text class="px-4 pb-2 text-body-2 text-medium-emphasis ml-14 mt-n2">
        {{ message }}
      </v-card-text>

      <!-- 按钮区 -->
      <v-card-actions class="pa-4 d-flex justify-end">
        <v-btn
          variant="text"
          @click="cancel"
        >
          {{ cancelText }}
        </v-btn>
        <v-btn
          :color="typeColor"
          variant="flat"
          elevation="1"
          @click="confirm"
        >
          {{ confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// --- 类型定义 ---
export type ConfirmType = 'success' | 'info' | 'warning' | 'error'

export interface ConfirmOptions {
  type?: ConfirmType
  confirmText?: string
  cancelText?: string
}

// --- 状态 ---
const show = ref(false)
const title = ref('')
const message = ref('')

const type = ref<ConfirmType>('warning')
const confirmText = ref('确认')
const cancelText = ref('取消')

// 存储 Promise 的 resolve 函数
let resolvePromise: ((value: boolean) => void) | null = null

// --- 配置 ---
const typeConfig: Record<ConfirmType, { icon: string; color: string }> = {
  success: { icon: 'mdi-check-circle-outline', color: 'success' },
  info:    { icon: 'mdi-information-outline', color: 'info' },
  warning: { icon: 'mdi-alert-outline',       color: 'warning' },
  error:   { icon: 'mdi-alert-circle-outline', color: 'error' }
}

// --- 计算属性 ---
const typeIcon = computed(() => typeConfig[type.value]?.icon || 'mdi-alert')
const typeColor = computed(() => typeConfig[type.value]?.color || 'primary')

// --- 方法 ---

/**
 * 打开对话框
 * @param titleText 标题
 * @param messageText 内容
 * @param options 配置项 { type, confirmText, cancelText }
 * @returns Promise<boolean> true=确认, false=取消
 */
const open = (
  titleText: string, 
  messageText: string, 
  options: ConfirmOptions = {}
): Promise<boolean> => {
  title.value = titleText
  message.value = messageText
  
  // 设置默认值
  type.value = options.type || 'warning'
  confirmText.value = options.confirmText || '确认'
  cancelText.value = options.cancelText || '取消'

  show.value = true

  return new Promise((resolve) => {
    resolvePromise = resolve
  })
}

const confirm = () => {
  show.value = false
  if (resolvePromise) resolvePromise(true)
}

const cancel = () => {
  show.value = false
  if (resolvePromise) resolvePromise(false)
}

// 暴露 open 方法
defineExpose({ open })
</script>