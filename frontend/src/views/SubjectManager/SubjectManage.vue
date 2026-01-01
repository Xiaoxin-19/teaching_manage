<template>
  <v-sheet class="fill-height pa-6 bg-background">
    <div class="d-flex flex-column h-100">

      <!-- 1. 头部区域 -->
      <div class="d-flex justify-space-between align-center mb-6">
        <div>
          <h2 class="text-h5 font-weight-bold text-high-emphasis">科目管理</h2>
          <!-- 已移除解释文字 -->
        </div>

        <!-- 右侧操作区 -->
        <div class="d-flex align-center gap-4 ml-auto" style="gap: 16px">
          <!-- 搜索 -->
          <v-text-field v-model="search" prepend-inner-icon="mdi-magnify" label="搜索科目" single-line hide-details
            density="compact" variant="outlined" class="rounded" style="width: 240px" bg-color="surface"
            @update:model-value="loadData"></v-text-field>

          <!-- 新增按钮 -->
          <v-btn color="primary" prepend-icon="mdi-plus" elevation="2" @click="openAdd">
            新增科目
          </v-btn>
        </div>
      </div>

      <!-- 2. 数据列表 -->
      <v-card elevation="2" class="flex-grow-1 rounded-lg overflow-hidden border">
        <v-data-table-server :headers="headers" :items="subjects" :loading="loading" :page="page"
          :items-per-page="itemsPerPage" :items-length="totalItems" @update:options="loadItems" hover
          density="comfortable" class="h-100">
          <!-- 科目名称：纯文字居中 -->
          <template v-slot:item.name="{ item }">
            <span class="font-weight-medium text-body-1">{{ item.name }}</span>
          </template>

          <!-- 在读学员：格式化显示 -->
          <template v-slot:item.student_count="{ item }">
            <div class="d-flex align-center justify-center">
              <!-- 使用固定宽度的 Chip 确保数字对齐 -->
              <v-chip size="small" :color="item.student_count > 0 ? 'primary' : 'grey'"
                :variant="item.student_count > 0 ? 'tonal' : 'flat'"
                class="font-weight-bold mr-2 tabular-nums justify-center" style="min-width: 54px;">
                {{ formatCount(item.student_count) }}
              </v-chip>
              <span class="text-caption text-medium-emphasis">人</span>
            </div>
          </template>

          <!-- 操作栏：内容居中对齐 -->
          <template v-slot:item.actions="{ item }">
            <div class="d-flex justify-center gap-2" style="gap: 8px">
              <v-btn size="small" variant="text" color="primary" prepend-icon="mdi-pencil" @click="openEdit(item)">
                编辑
              </v-btn>

              <v-tooltip location="top" :disabled="item.student_count === 0">
                <template v-slot:activator="{ props }">
                  <span v-bind="props">
                    <v-btn size="small" variant="text" color="error" prepend-icon="mdi-delete"
                      :disabled="item.student_count > 0" @click="handleDelete(item)">
                      删除
                    </v-btn>
                  </span>
                </template>
                <span>该科目下仍有学员，禁止删除</span>
              </v-tooltip>
            </div>
          </template>

          <template v-slot:no-data>
            <div class="pa-8 text-center text-medium-emphasis">
              <v-icon size="64" class="mb-2 text-disabled">mdi-playlist-remove</v-icon>
              <div class="text-body-1">暂无科目数据</div>
              <v-btn color="primary" variant="text" class="mt-2" @click="openAdd">添加第一个科目</v-btn>
            </div>
          </template>
        </v-data-table-server>
      </v-card>

      <!-- 3. 新增/编辑弹窗 -->
      <v-dialog v-model="dialogVisible" max-width="400px">
        <v-card class="rounded-lg">
          <!-- 弹窗标题栏：增加了动态图标 -->
          <v-card-title class="d-flex justify-space-between align-center py-3 px-4 bg-surface">
            <div class="d-flex align-center">
              <v-icon :icon="isEdit ? 'mdi-notebook-edit-outline' : 'mdi-notebook-plus-outline'" color="primary"
                class="mr-2"></v-icon>
              <span class="text-subtitle-1 font-weight-bold">{{ isEdit ? '编辑科目' : '新增科目' }}</span>
            </div>
            <v-btn icon="mdi-close" variant="text" size="small" @click="dialogVisible = false"></v-btn>
          </v-card-title>

          <v-divider></v-divider>

          <v-card-text class="pa-6">
            <v-form ref="formRef" @submit.prevent="handleSave">
              <v-text-field v-model="formData.name" label="科目名称" placeholder="例如：钢琴、乐理" variant="outlined"
                density="comfortable" :rules="[v => !!v || '科目名称不能为空']" autofocus
                prepend-inner-icon="mdi-bookshelf"></v-text-field>

              <div class="text-caption text-medium-emphasis mt-2">
                <v-icon size="small" class="mr-1">mdi-information-outline</v-icon>
                科目名称将用于全校课程关联。
              </div>
            </v-form>
          </v-card-text>

          <!-- 底部操作栏：左侧显示时间，右侧按钮 -->
          <v-card-actions class="pa-4 pt-0 d-flex align-center">
            <!-- 上次编辑时间 -->
            <div v-if="isEdit && formData.lastModified"
              class="text-caption text-medium-emphasis d-flex align-center mr-auto">
              <v-icon icon="mdi-clock-outline" size="x-small" class="mr-1"></v-icon>
              上次编辑: {{ formData.lastModified }}
            </div>
            <v-spacer v-else></v-spacer> <!-- 占位符，确保按钮靠右 -->

            <v-btn variant="text" @click="dialogVisible = false">取消</v-btn>
            <v-btn color="primary" variant="elevated" @click="handleSave" :loading="loading"
              :disabled="!formData.name">保存</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

    </div>
  </v-sheet>
</template>

<script setup lang="ts">
import { useSubjectManage } from './SubjectManage.logic'

const {
  loading,
  search,
  subjects,
  headers,
  dialogVisible,
  isEdit,
  formData,
  formRef,
  page,
  itemsPerPage,
  totalItems,
  loadItems,
  loadData,
  openAdd,
  openEdit,
  handleSave,
  handleDelete,
  formatCount
} = useSubjectManage()
</script>

<style scoped>
/* 数字等宽字体，保证对齐 */
.tabular-nums {
  font-feature-settings: "tnum";
  font-variant-numeric: tabular-nums;
}
</style>