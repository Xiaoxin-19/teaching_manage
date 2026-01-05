<template>
  <v-app>

    <v-system-bar style="--wails-draggable: drag" window>
      <div class="titlebar-left">
        <span class="titlebar-title">Teaching Manager</span>
      </div>
      <v-spacer></v-spacer>
      <v-btn density="compact" variant="text" @click="toggleTheme" :icon="themeIcon" aria-label="切换主题">
      </v-btn>
      <v-btn density="compact" variant="text" @click="handleMinimise" icon="mdi-window-minimize" aria-label="最小化">
        <v-icon size="18" icon="mdi-window-minimize" />
      </v-btn>
      <v-btn density="compact" variant="text" @click="handleToggleMax" icon="mdi-window-maximize" aria-label="最大化/还原">
      </v-btn>
      <v-btn density="compact" variant="text" class=" close" @click="handleClose" icon="mdi-close" aria-label="关闭">
      </v-btn>
    </v-system-bar>

    <!-- 侧边导航组件 -->
    <NavDrawer v-model="drawer" v-model:rail="rail" @toggle-theme="toggleTheme" />



    <!-- 主内容区域 -->
    <v-main>
      <!-- 动态内容视口 -->
      <v-container fluid class="pa-0 h-100">

        <router-view v-slot="{ Component }">
          <transition name="fade-transition" mode="out-in">
            <keep-alive>
              <component :is="Component" />
            </keep-alive>
          </transition>
        </router-view>

      </v-container>

      <!-- ref="globalSnackbarRef" 用于在 js 中获取组件实例 -->
      <GlobalSnackBar ref="globalSnackbarRef" />

      <!-- 2. 全局确认对话框 (Confirm Dialog) -->
      <GlobalConfirmDialog ref="globalConfirmRef" />

      <!-- 3. 全屏遮罩 -->
      <GlobalOverlay ref="globalOverlayRef" />
    </v-main>


  </v-app>
</template>

<script setup lang="ts">
import { onMounted, ref, provide, watch, computed } from 'vue'
import NavDrawer from './components/NavDrawer.vue'
import GlobalSnackBar from './components/GlobalSnackBar.vue'
import { registerToast } from './composables/useToast'
import GlobalConfirmDialog, { ConfirmOptions } from './components/GlobalConfirmDialog.vue'
import { registerConfirm } from './composables/useConfirm'
import GlobalOverlay from './components/GlobalOverlay.vue'
import { registerGlobalOverlay } from './composables/useGlobalOverlay'
import { WindowMinimise, WindowToggleMaximise, WindowIsMaximised, Quit, WindowSetDarkTheme, WindowSetLightTheme } from '../wailsjs/runtime/runtime'
import { useTheme } from 'vuetify/lib/composables/theme'

// --- 全局状态 ---
const theme = useTheme()
const drawer = ref(true)
const rail = ref(true)
const isMaximized = ref(false)
const themeIcon = computed(() => theme.current.value.dark === true ? 'mdi-weather-night' : 'mdi-white-balance-sunny')

// --- 方法 ---
const toggleTheme = () => {
  theme.cycle()
}

const handleMinimise = () => WindowMinimise()

const handleToggleMax = async () => {
  WindowToggleMaximise()
  isMaximized.value = await WindowIsMaximised()
}

const handleClose = () => Quit()


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


// ----------------------------------------------------------------
// 3. 全局遮罩 (全屏禁用操作)
// ----------------------------------------------------------------
const globalOverlayRef = ref<InstanceType<typeof GlobalOverlay> | null>(null)

const overlayController = {
  show: (msg?: string) => {
    if (globalOverlayRef.value) {
      globalOverlayRef.value.show(msg)
    } else {
      console.warn('GlobalOverlay 组件未挂载！')
    }
  },
  hide: () => {
    if (globalOverlayRef.value) {
      globalOverlayRef.value.hide()
    }
  },
  setMessage: (msg: string) => {
    if (globalOverlayRef.value && globalOverlayRef.value.setMessage) {
      globalOverlayRef.value.setMessage(msg)
    }
  }
}

provide('globalOverlay', overlayController)
registerGlobalOverlay(overlayController)

onMounted(async () => {
  isMaximized.value = await WindowIsMaximised()
  // 初始化侧边栏宽度
  const initialWidth = drawer.value ? (rail.value ? '56px' : '210px') : '0px'
  document.documentElement.style.setProperty('--drawer-width', initialWidth)
})


// Watch drawer/rail state to update computed css var for layout margin
watch([drawer, rail], ([d, r]) => {
  const width = d ? (r ? '56px' : '210px') : '0px'
  document.documentElement.style.setProperty('--drawer-width', width)
})
</script>
<style scoped>
:global(:root) {
  --titlebar-height: 36px;
}



.titlebar-left {
  display: flex;
  gap: 8px;
  align-items: baseline;
  font-size: 12px;
  letter-spacing: 0.2px;
}

.titlebar-title {
  font-weight: 700;
  font-size: 13px;
}

.titlebar-subtitle {
  opacity: 0.7;
}

/* 禁用所有输入控件的浏览器自动填充 */
:deep(input),
:deep(textarea) {
  -webkit-autocomplete: off;
  autocomplete: off;
}

/* 侧边栏适配无边框标题栏高度 */
:deep(.v-navigation-drawer) {
  position: fixed !important;
  top: var(--titlebar-height) !important;
  left: 0;
  height: calc(100vh - var(--titlebar-height)) !important;
  z-index: 9 !important;
}

/* 主内容区域整体向下错开标题栏高度 */
:deep(.v-main) {
  padding-top: var(--titlebar-height);
}

/* 主内容与侧边栏间距 */
:global(:root) {
  --drawer-width: 210px;
}

:deep(.main-layout-container) {
  margin-left: var(--drawer-width);
  transition: margin-left 0.12s ease;
}

/* 隐藏滚动条 */
div::-webkit-scrollbar {
  display: none;
}

/* 右侧主内容容器 */
.main-layout-container {
  min-height: calc(100vh - var(--titlebar-height));
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
  min-height: calc(100vh - var(--titlebar-height));
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