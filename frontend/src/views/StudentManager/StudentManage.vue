<script setup lang="ts">
import { onActivated } from 'vue'
import { useStudentManage } from './StudentManage.logic'
import ModifyStudent from './ModifyStudent.vue'

const {
  search,
  filterStatus,
  hasActiveFilters,
  clearAllFilters,
  page,
  itemsPerPage,
  totalItems,
  loading,
  dialog,
  editedIndex,
  editedItem,
  students,
  headers,
  getStatusColor,
  getStatusText,
  getGenderColor,
  getGenderLabel,
  openEdit,
  openDelete,
  saveStudent,
  openAdd,
  exportStudents,
  loadItems,
} = useStudentManage()

onActivated(() => {
  loadItems({ page: page.value, itemsPerPage: itemsPerPage.value })
})
</script>

<template>
  <v-sheet class="fill-height pa-6 bg-background">
    <div class="d-flex flex-column h-100">

      <!-- 1. 顶部操作栏 -->
      <div class="d-flex justify-space-between align-center mb-4 flex-shrink-0">
        <v-text-field v-model="search" prepend-inner-icon="mdi-magnify" label="搜索学生姓名、学号、手机号" single-line hide-details
          density="compact" variant="outlined" class="rounded" style="max-width: 300px"></v-text-field>

        <div class="d-flex">
          <v-btn prepend-icon="mdi-microsoft-excel" color="success" variant="outlined" class="mr-3" height="40"
            @click="exportStudents">
            导出学生列表
          </v-btn>

          <v-btn color="primary" prepend-icon="mdi-plus" elevation="2" height="40" @click="openAdd">
            新增学生
          </v-btn>
        </div>
      </div>

      <!-- 2. 数据表格区域 -->
      <v-card elevation="2" class="flex-grow-1 d-flex flex-column rounded-lg overflow-hidden border">
        <v-data-table-server :headers="headers" :items="students" :search="search" :page="page"
          :items-per-page="itemsPerPage" :loading="loading" hover fixed-header density="comfortable"
          :items-length="totalItems" @update:options="loadItems" class="h-100">

          <!-- 顶部插槽：显示生效的筛选条件 -->
          <template v-slot:top>
            <v-expand-transition>
              <div v-if="hasActiveFilters" class="px-4 py-3 border-b d-flex align-center flex-wrap"
                style="gap: 8px; background-color: rgba(var(--v-theme-surface-variant), 0.1);">
                <div class="text-caption font-weight-bold mr-2 text-medium-emphasis d-flex align-center">
                  <v-icon size="small" class="mr-1">mdi-filter-variant</v-icon>
                  当前筛选:
                </div>

                <v-chip v-if="filterStatus" closable size="small" color="primary" variant="flat"
                  @click:close="clearAllFilters">
                  状态: {{ getStatusText(filterStatus) }}
                </v-chip>

                <v-spacer></v-spacer>

                <v-btn variant="text" size="small" color="error" prepend-icon="mdi-delete-sweep-outline"
                  @click="clearAllFilters">
                  重置筛选
                </v-btn>
              </div>
            </v-expand-transition>
            <v-divider v-if="hasActiveFilters"></v-divider>
          </template>

          <!-- 1. 状态列头筛选 -->
          <template v-slot:header.status="{ column }">
            <div class="header-filter-container d-flex align-center justify-center">
              <span class="font-weight-bold mr-2">{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom end" offset="5">
                <template v-slot:activator="{ props }">
                  <v-icon v-bind="props" :icon="filterStatus ? 'mdi-filter' : 'mdi-filter-outline'" size="small"
                    class="filter-icon" :class="{ active: filterStatus }"
                    :color="filterStatus ? 'primary' : ''"></v-icon>
                </template>
                <v-card min-width="200" class="pa-2 rounded-lg elevation-4">
                  <div class="text-subtitle-2 mb-2 px-2 font-weight-bold d-flex align-center">
                    <v-icon size="small" class="mr-2" color="primary">mdi-filter-check</v-icon>
                    筛选状态
                  </div>
                  <v-list density="compact" nav select-strategy="single-leaf">
                    <v-list-item :active="filterStatus === null"
                      @click="filterStatus = null; loadItems({ page: 1, itemsPerPage })" color="primary" rounded="lg">
                      <template v-slot:prepend>
                        <v-icon icon="mdi-all-inclusive" size="small"></v-icon>
                      </template>
                      <v-list-item-title>全部状态</v-list-item-title>
                    </v-list-item>

                    <v-list-item :active="filterStatus === 1"
                      @click="filterStatus = 1; loadItems({ page: 1, itemsPerPage })" color="success" rounded="lg">
                      <template v-slot:prepend>
                        <v-icon icon="mdi-check-circle-outline" size="small" color="success"></v-icon>
                      </template>
                      <v-list-item-title>正常上课</v-list-item-title>
                    </v-list-item>

                    <v-list-item :active="filterStatus === 2"
                      @click="filterStatus = 2; loadItems({ page: 1, itemsPerPage })" color="warning" rounded="lg">
                      <template v-slot:prepend>
                        <v-icon icon="mdi-pause-circle-outline" size="small" color="warning"></v-icon>
                      </template>
                      <v-list-item-title>暂时停课</v-list-item-title>
                    </v-list-item>

                    <v-list-item :active="filterStatus === 3"
                      @click="filterStatus = 3; loadItems({ page: 1, itemsPerPage })" color="error" rounded="lg">
                      <template v-slot:prepend>
                        <v-icon icon="mdi-close-circle-outline" size="small" color="error"></v-icon>
                      </template>
                      <v-list-item-title>已经退学</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-card>
              </v-menu>
            </div>
          </template>


          <template v-slot:item.phone="{ item }">
            <div class="text-center">{{ item.phone == "" ? "-" : item.phone }}</div>
          </template>

          <!-- 性别列 -->
          <template v-slot:item.gender="{ item }">
            <v-chip :color="getGenderColor(item.gender)" variant="flat" class="text-uppercase" size="small" label>
              {{ getGenderLabel(item.gender) }}
            </v-chip>
          </template>

          <!-- 状态列 -->
          <template v-slot:item.status="{ item }">
            <v-chip :color="getStatusColor(item.status)" variant="outlined" class="text-uppercase" size="small" label>
              {{ getStatusText(item.status) }}
            </v-chip>
          </template>


          <!-- 操作列 -->
          <template v-slot:item.actions="{ item }">
            <div class="d-flex justify-center align-center gap-2" style="gap: 8px;">
              <v-tooltip location="top" text="编辑信息">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" class="text-medium-emphasis" v-bind="props"
                    @click="openEdit(item)">
                    <v-icon>mdi-pencil-outline</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>

              <v-tooltip location="top" text="删除档案">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="error" v-bind="props" @click="openDelete(item)">
                    <v-icon>mdi-delete-outline</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>
            </div>
          </template>

          <template v-slot:no-data>
            <div class="pa-8 text-center text-medium-emphasis">
              <v-icon size="64" class="mb-2 text-disabled">mdi-account-off-outline</v-icon>
              <div class="text-body-1">暂无学生数据</div>
              <v-btn color="primary" variant="text" class="mt-2" @click="openAdd">点击添加第一位学生</v-btn>
            </div>
          </template>
        </v-data-table-server>
      </v-card>

      <!-- 3. 新增/编辑弹窗 -->
      <ModifyStudent v-model:modelValue="dialog" :is-edit="editedIndex >= 0" :initial-data="editedItem"
        @save="saveStudent" />
    </div>
  </v-sheet>
</template>

<style scoped>
/* 隐藏滚动条 */
div::-webkit-scrollbar {
  display: none;
}
</style>