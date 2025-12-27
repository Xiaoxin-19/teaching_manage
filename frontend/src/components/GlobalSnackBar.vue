<template>
  <!-- 
    全局消息容器
    注意：如果将来使用 <Teleport to="body">，请务必在内部包裹 <v-theme-provider :theme="theme.global.name.value"> 
    以确保主题样式能正确穿透。当前默认放置在 v-app 内部，因此能自动继承主题。
  -->
  <div class="global-toast-layer">

    <!-- 1. 顶部居右 (Top Right) -->
    <div class="toast-container top-right">
      <transition-group name="slide-x-reverse-transition">
        <div v-for="item in getItems('top-right')" :key="item.id" class="toast-wrapper">
          <v-alert :color="item.color" :icon="item.icon" :model-value="true" :variant="alertVariant" density="compact"
            class="mb-2 toast-alert" closable elevation="4" border="start" @click:close="remove(item.id)">
            {{ item.text }}
          </v-alert>
        </div>
      </transition-group>
    </div>

    <!-- 2. 底部居右 (Bottom Right) - 默认位置 -->
    <div class="toast-container bottom-right">
      <transition-group name="slide-x-reverse-transition">
        <div v-for="item in getItems('bottom-right')" :key="item.id" class="toast-wrapper">
          <v-alert :color="item.color" :icon="item.icon" :model-value="true" :variant="alertVariant" density="compact"
            class="mt-2 toast-alert" closable elevation="4" border="start" @click:close="remove(item.id)">
            {{ item.text }}
          </v-alert>
        </div>
      </transition-group>
    </div>

    <!-- 3. 顶部居中 (Top Center) -->
    <div class="toast-container top-center">
      <transition-group name="slide-y-transition">
        <div v-for="item in getItems('top-center')" :key="item.id" class="toast-wrapper">
          <v-alert :color="item.color" :icon="item.icon" :model-value="true" :variant="alertVariant" density="compact"
            class="mb-2 toast-alert" closable elevation="4" border="start" @click:close="remove(item.id)">
            {{ item.text }}
          </v-alert>
        </div>
      </transition-group>
    </div>

    <!-- 4. 底部居中 (Bottom Center) -->
    <div class="toast-container bottom-center">
      <transition-group name="slide-y-reverse-transition">
        <div v-for="item in getItems('bottom-center')" :key="item.id" class="toast-wrapper">
          <v-alert :color="item.color" :icon="item.icon" :model-value="true" :variant="alertVariant" density="compact"
            class="mt-2 toast-alert" closable elevation="4" border="start" @click:close="remove(item.id)">
            {{ item.text }}
          </v-alert>
        </div>
      </transition-group>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTheme } from 'vuetify'

// --- 类型定义 ---
export type SnackbarType = 'success' | 'error' | 'info' | 'warning'
export type SnackbarLocation = 'top-right' | 'top-center' | 'bottom-right' | 'bottom-center'

interface SnackbarItem {
  id: number
  text: string
  color: string
  icon: string
  location: SnackbarLocation
  timer?: number
}

// --- 状态 ---
const items = ref<SnackbarItem[]>([])
let idCounter = 0

// --- 主题适配 ---
const theme = useTheme()

// 计算属性：根据当前主题模式动态调整 Alert 的风格
// 亮色模式下使用 'elevated' (高亮实心)，深色模式下使用 'tonal' (柔和半透明)，视觉体验更佳
const alertVariant = computed(() => {
  return theme.global.current.value.dark ? 'tonal' : 'elevated'
})

// --- 方法 ---

/**
 * 根据位置获取对应的消息列表
 */
const getItems = (loc: SnackbarLocation) => {
  return items.value.filter(item => item.location === loc)
}

/**
 * 移除消息
 */
const remove = (id: number) => {
  const index = items.value.findIndex(item => item.id === id)
  if (index !== -1) {
    // 如果有定时器，清除它
    if (items.value[index].timer) {
      clearTimeout(items.value[index].timer)
    }
    items.value.splice(index, 1)
  }
}

/**
 * 打开消息提示 (对外暴露的核心方法)
 * @param msg 消息文本
 * @param type 消息类型 (默认 success)
 * @param location 显示位置 (默认 bottom-right)
 * @param timeout 自动关闭时间 (毫秒，默认 3000)
 */
const open = (
  msg: string,
  type: SnackbarType = 'success',
  location: SnackbarLocation = 'bottom-right',
  timeout: number = 3000
) => {
  const id = ++idCounter

  // 确定颜色和图标
  let color = 'success'
  let icon = 'mdi-check-circle'

  switch (type) {
    case 'success':
      color = 'success'
      icon = 'mdi-check-circle'
      break
    case 'error':
      color = 'error'
      icon = 'mdi-alert-circle'
      break
    case 'info':
      color = 'info'
      icon = 'mdi-information'
      break
    case 'warning':
      color = 'warning'
      icon = 'mdi-alert'
      break
  }

  // 创建新消息对象
  const newItem: SnackbarItem = {
    id,
    text: msg,
    color,
    icon,
    location
  }

  // 设置自动关闭定时器
  if (timeout > 0) {
    newItem.timer = window.setTimeout(() => {
      remove(id)
    }, timeout)
  }

  // 添加到队列
  items.value.push(newItem)
}

// 暴露给父组件
defineExpose({
  open,
  remove
})
</script>

<style scoped>
/* 全局层：确保 pointer-events 不会阻挡页面其他操作 */
.global-toast-layer {
  pointer-events: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999;
}

/* 通用容器样式 */
.toast-container {
  position: absolute;
  display: flex;
  flex-direction: column;
  padding: 16px;
  width: 100%;
  max-width: 400px;
  /* 限制 Toast 最大宽度 */
}

/* 消息框本身的样式 */
.toast-wrapper {
  pointer-events: auto;
  /* 恢复点击事件，以便可以点击关闭按钮 */
  transition: all 0.3s ease;
}

/* Alert 样式微调：让它看起来更像 Snackbar */
.toast-alert {
  font-size: 0.875rem;
  font-weight: 500;
  /* 确保文字颜色在两种模式下都有良好的对比度 */
  transition: all 0.3s ease;
}

/* --- 各个方位定位 --- */

/* 顶部居右 */
.top-right {
  top: 0;
  right: 0;
  align-items: flex-end;
  /* 内容靠右对齐 */
}

/* 底部居右 */
.bottom-right {
  bottom: 0;
  right: 0;
  align-items: flex-end;
  /* 关键：底部容器反向排列，新消息出现在最底部，旧消息被顶上去 */
  flex-direction: column-reverse;
}

/* 顶部居中 */
.top-center {
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  align-items: center;
}

/* 底部居中 */
.bottom-center {
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  align-items: center;
  flex-direction: column-reverse;
}
</style>