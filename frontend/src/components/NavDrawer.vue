<script setup lang="ts">
import { useTheme } from 'vuetify/lib/composables/theme';

const theme = useTheme()
// 定义 Props：接收父组件的状态
const props = defineProps<{
  modelValue: boolean // 控制 drawer 显示/隐藏 (v-model)
  rail: boolean       // 控制迷你模式
}>()

// 定义 Emits：向父组件发送事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'update:rail', value: boolean): void
  (e: 'toggleTheme'): void
}>()

// 菜单配置数据，使用 router 路径 (或 name)
const menuItems = [
  { title: '数据看板', icon: 'mdi-view-dashboard', to: { name: 'dashboard' } },
  { title: '学生档案', icon: 'mdi-account-group', to: { name: 'students' } },
  { title: '教师管理', icon: 'mdi-human-male-board', to: { name: 'teachers' } },
  { title: '科目管理', icon: 'mdi-bookshelf', to: { name: 'subjects' } },
  { title: '课程管理', icon: 'mdi-school-outline', to: { name: 'courses' } },
  { title: '上课记录', icon: 'mdi-clipboard-text-clock', to: { name: 'records' } },
  { title: '财务管理', icon: 'mdi-currency-usd', to: { name: 'finance' } },
  { divider: true },
  { title: '系统设置', icon: 'mdi-cog', to: { name: 'settings' } },
]
</script>

<template>
  <v-navigation-drawer :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)" :rail="rail"
    width="210" permanent elevation="2" class="bg-surface" @click="emit('update:rail', false)">

    <!-- 2. 菜单列表区（使用 :to + link 由 vue-router 管理激活状态） -->
    <v-list nav density="compact" class="mt-2">
      <template v-for="(item, index) in menuItems" :key="index">
        <v-divider v-if="item.divider" class="my-2"></v-divider>

        <v-list-item v-else :to="item.to" link color="primary" rounded="lg" class="mb-1">
          <template v-slot:prepend>
            <v-icon :icon="item.icon"></v-icon>
          </template>
          <v-list-item-title class="font-weight-medium">
            {{ item.title }}
          </v-list-item-title>

          <!-- Rail 模式下的悬浮提示 -->
          <v-tooltip v-if="rail" activator="parent" location="end">{{ item.title }}</v-tooltip>
        </v-list-item>
      </template>
    </v-list>

    <!-- 3. 底部操作区 -->
    <template v-slot:append>
      <div class="pa-2">
        <!-- 折叠/展开 切换按钮 -->
        <v-list-item nav rounded="lg" @click.stop="emit('update:rail', !rail)" class="mb-2" color="primary">
          <template v-slot:prepend>
            <v-icon>{{ rail ? 'mdi-chevron-right' : 'mdi-chevron-left' }}</v-icon>
          </template>
          <v-list-item-title>
            {{ rail ? '展开' : '收起侧边栏' }}
          </v-list-item-title>
        </v-list-item>
      </div>
    </template>
  </v-navigation-drawer>
</template>