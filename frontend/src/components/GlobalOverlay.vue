<template>
  <v-overlay v-model="visible" class="global-overlay" :persistent="true" scroll-strategy="block" contained
    :scrim="'rgba(0, 0, 0, 0.45)'">
    <div class="overlay-card">
      <v-progress-circular indeterminate color="primary" size="42" width="4" class="mb-3" />
      <div class="overlay-text">{{ text }}</div>
      <div class="overlay-sub">{{ subText }}</div>
    </div>
  </v-overlay>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'

const props = defineProps<{
  defaultMessage?: string
  subMessage?: string
}>()

const visible = ref(false)
const text = ref(props.defaultMessage || '正在处理，请稍候...')
const subText = computed(() => props.subMessage || '请稍候，操作完成前请勿关闭窗口')

const show = (msg?: string) => {
  text.value = msg || props.defaultMessage || '正在处理，请稍候...'
  visible.value = true
}

const setMessage = (msg: string) => {
  text.value = msg
}

const hide = () => {
  visible.value = false
}

// 防止意外挂起时文本被置空
watch(visible, (val) => {
  if (!val) {
    text.value = props.defaultMessage || '正在处理，请稍候...'
  }
})

// 暴露给父组件调用
defineExpose({ show, hide, setMessage })
</script>

<style scoped>
.global-overlay {
  display: flex;
  align-items: center;
  justify-content: center;
}

.global-overlay :deep(.v-overlay__scrim) {
  backdrop-filter: blur(2px);
}

.overlay-card {
  min-width: 260px;
  max-width: min(90vw, 360px);
  padding: 20px 24px;
  border-radius: 14px;
  background: rgba(var(--v-theme-surface), 0.9);
  color: rgb(var(--v-theme-on-surface));
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.25);
  text-align: center;
}

.overlay-text {
  font-weight: 600;
  font-size: 15px;
  margin-bottom: 6px;
}

.overlay-sub {
  font-size: 13px;
  color: rgba(var(--v-theme-on-surface), 0.7);
  text-align: center;
  line-height: 1.4;
}

@media (max-width: 640px) {
  .overlay-card {
    min-width: 220px;
    padding: 16px 18px;
    border-radius: 12px;
  }

  .overlay-text {
    font-size: 14px;
  }

  .overlay-sub {
    font-size: 12px;
  }
}
</style>
