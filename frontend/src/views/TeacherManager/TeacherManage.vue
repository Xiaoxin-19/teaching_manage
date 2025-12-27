<template>
  <v-sheet class="fill-height pa-6 ">
    <div class="d-flex flex-column pa-6">
      <!-- 页面头部 -->
      <div class="d-flex justify-space-between align-center mb-4">
        <!-- 左侧：搜索框 -->
        <!-- 移除 bg-white -->
        <v-text-field v-model="search" density="compact" variant="outlined" label="搜索教师姓名"
          prepend-inner-icon="mdi-magnify" hide-details single-line class="rounded"
          style="max-width: 300px;"></v-text-field>

        <!-- 右侧：按钮组 -->
        <div class="d-flex">
          <!-- 导出按钮 -->
          <!-- 移除 bg-white -->
          <v-btn color="success" variant="outlined" prepend-icon="mdi-microsoft-excel" class="mr-3" @click="exportData">
            导出 Excel
          </v-btn>

          <!-- 新增按钮 -->
          <v-btn color="primary" prepend-icon="mdi-plus" elevation="2" @click="openAdd">
            新增教师
          </v-btn>
        </div>
      </div>

      <!-- 教师数据表格 -->
      <v-card class="rounded-lg elevation-2 border">
        <v-data-table :headers="headers" :items="teachers" :search="search" density="compact"
          v-model:items-per-page="itemsPerPage" v-model:page="page" @update:options="loadItems" hover>
          <!-- 姓名列 -->
          <template v-slot:item.name="{ item }">
            <span class="font-weight-medium text-body-2">{{ item.name }}</span>
          </template>

          <!-- 性别列 -->
          <template v-slot:item.gender="{ item }">
            <v-chip :color="getGenderColor(item.gender)" variant="flat" class="text-uppercase" size="small" label>
              {{ getGenderLabel(item.gender) }}
            </v-chip>
          </template>

          <!-- 备注列 -->
          <template v-slot:item.remark="{ item }">
            <!-- 如果有内容，显示截断文本和 Tooltip -->
            <v-tooltip location="top" v-if="item.remark">
              <template v-slot:activator="{ props }">
                <div v-bind="props" class="text-truncate text-body-2" style="max-width: 150px; cursor: default;">
                  {{ item.remark }}
                </div>
              </template>
              <span>{{ item.remark }}</span>
            </v-tooltip>
            <!-- 如果无内容，显示占位符 -->
            <span v-else class="text-disabled text-caption">-</span>
          </template>


          <!-- 操作列 -->
          <template v-slot:item.actions="{ item }">
            <div class="d-flex justify-end">
              <v-btn icon="mdi-pencil-outline" size="small" variant="text" color="primary" @click="openEdit(item)"
                title="编辑详情"></v-btn>
              <v-btn icon="mdi-delete-outline" size="small" variant="text" color="error" @click="deleteItem(item)"
                title="删除"></v-btn>
            </div>
          </template>

          <!-- 无数据状态提示 -->
          <template v-slot:no-data>
            <div class="pa-4 text-center text-medium-emphasis">暂无教师数据</div>
          </template>
        </v-data-table>
      </v-card>


      <!-- 弹窗组件 -->
      <ModifyTeacher v-model="dialogVisible" :is-edit="isEdit" :initial-data="currentData" @save="handleSave" />
    </div>
  </v-sheet>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import ModifyTeacher from './ModifyTeacher.vue'
import { useTeacherManage } from './TeacherManage.logic'
import { LogDebug } from '../../../wailsjs/runtime/runtime';



const {
  search,
  dialogVisible,
  isEdit,
  currentData,
  headers,
  teachers,
  itemsPerPage,
  page,
  openAdd,
  openEdit,
  deleteItem,
  loadItems,
  exportData,
  handleSave,
  fetchTeachers,
  getGenderColor,
  getGenderLabel,
} = useTeacherManage()

onMounted(() => {
  LogDebug("TeacherManage.vue mounted")
  // 初始加载教师数据
  fetchTeachers()
})
</script>
