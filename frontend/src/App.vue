<script setup lang="ts">
import { ref } from 'vue'
import NavDrawer from './components/NavDrawer.vue'

// --- 全局状态 ---
const theme = ref('light')
const drawer = ref(true)
const rail = ref(false)

// --- 方法 ---
const toggleTheme = () => {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
}
</script>

<template>
  <v-app :theme="theme">

    <!-- 侧边导航组件 -->
    <NavDrawer v-model="drawer" v-model:rail="rail" :theme="theme" @toggle-theme="toggleTheme" />

    <!-- 主内容区域 -->
    <v-main>
      <div class="main-layout-container">

        <!-- 动态内容视口 -->
        <div class="layout-content bg-grey-lighten-4">
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
  </v-app>
</template>

<style scoped>
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