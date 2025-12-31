<template>
  <v-app :theme="theme">

    <!-- 侧边导航组件 -->
    <NavDrawer v-model="drawer" v-model:rail="rail" :theme="theme" @toggle-theme="toggleTheme" />

    <!-- 主内容区域 -->
    <v-main>
      <div class="main-layout-container">

        <!-- 动态内容视口 -->
        <div class="layout-content bg-background">
          <v-container fluid class="pa-0 h-100">

            <router-view v-slot="{ Component }">
              <transition name="fade-transition" mode="out-in">
                <keep-alive>
                  <component :is="Component" />
                </keep-alive>
              </transition>
            </router-view>

          </v-container>
        </div>

      </div>
    </v-main>

    <!-- ref="globalSnackbarRef" 用于在 js 中获取组件实例 -->
    <GlobalSnackBar ref="globalSnackbarRef" />

    <!-- 2. 全局确认对话框 (Confirm Dialog) -->
    <GlobalConfirmDialog ref="globalConfirmRef" />
  </v-app>
</template>

<script setup lang="ts">
import { ref, provide } from 'vue'
import NavDrawer from './components/NavDrawer.vue'
import GlobalSnackBar from './components/GlobalSnackBar.vue'
import { registerToast } from './composables/useToast'
import GlobalConfirmDialog, { ConfirmOptions } from './components/GlobalConfirmDialog.vue'
import { registerConfirm } from './composables/useConfirm'

// --- 全局状态 ---
const theme = ref('light')
const drawer = ref(true)
const rail = ref(true)

// --- 方法 ---
const toggleTheme = () => {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
}


// ----------------------------------------------------------------
// 1. Snackbar (消息提示) 配置逻辑
// ----------------------------------------------------------------

// 定义组件引用的类型
const globalSnackbarRef = ref<InstanceType<typeof GlobalSnackBar> | null>(null)

/**
 * 2. 定义一个通用的调用函数
 * 这个函数会去调用组件内部暴露的 open 方法
 */
const showToast = (
  msg: string,
  type: 'success' | 'error' | 'info' | 'warning' = 'success',
  location: 'top-right' | 'top-center' | 'bottom-right' | 'bottom-center' = 'bottom-right',
  timeout = 3000
) => {
  if (globalSnackbarRef.value) {
    globalSnackbarRef.value.open(msg, type, location, timeout)
  } else {
    console.warn('GlobalSnackbar 组件未挂载！')
  }
}

// 3. 将这个函数注入给所有后代组件
// 'showToast' 是注入的 key，子组件通过这个 key 来获取函数
provide('showToast', showToast)
registerToast(showToast)


// ----------------------------------------------------------------
// 2. Confirm Dialog (确认对话框) 配置逻辑
// ----------------------------------------------------------------

// 获取组件实例引用
const globalConfirmRef = ref<InstanceType<typeof GlobalConfirmDialog> | null>(null)

/**
 * 定义全局显示确认框的函数 (返回 Promise)
 * @param title 标题
 * @param message 详细内容
 * @param options 配置项 { type, confirmText, cancelText }
 * @returns Promise<boolean> (true=确认, false=取消)
 */
const showConfirm = (
  title: string,
  message: string,
  options?: ConfirmOptions
): Promise<boolean> => {
  if (globalConfirmRef.value) {
    // 调用组件内部的 open 方法，它返回一个 Promise
    return globalConfirmRef.value.open(title, message, options)
  }

  console.warn('GlobalConfirmDialog 组件尚未挂载')
  // 如果组件未挂载，默认返回 false (取消操作) 以保证安全
  return Promise.resolve(false)
}

// 注入给所有后代组件 (key: 'showConfirm')
// 配合 src/composables/useConfirm.ts 使用
provide('showConfirm', showConfirm)
registerConfirm(showConfirm)
</script>
<style scoped>
/* 隐藏滚动条 */
div::-webkit-scrollbar {
  display: none;
}

/* 右侧主内容容器 */
.main-layout-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: rgb(var(--v-theme-background));
  transition: background-color 0.3s;
}

/* 内容视图区 */
.layout-content {
  flex: 1;
  overflow-y: auto;
  position: relative;
}

/* 占位符样式 */
.placeholder-box {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  border: 2px dashed #ccc;
  border-radius: 8px;
}

/* 过渡动画 */
.fade-transition-enter-active,
.fade-transition-leave-active {
  transition: opacity 0.2s ease;
}

.fade-transition-enter-from,
.fade-transition-leave-to {
  opacity: 0;
}

/* 滚动条美化 */
.layout-content::-webkit-scrollbar {
  width: 8px;
}

.layout-content::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}
</style>