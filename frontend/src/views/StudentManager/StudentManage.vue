<script setup lang="ts">
import { useStudentManage } from './StudentManage.logic'
import DetailsDialog from './DetailsDialog.vue'
import ModifyStudent from './ModifyStudent.vue'
import RechargeDialog from './RechargeDialog.vue'

const {
  search,
  page,
  itemsPerPage,
  totalItems,
  loading,
  dialog,
  dialogRecharge,
  dialogDetails,
  editedIndex,
  editedItem,
  rechargeItem,
  students,
  teacherOptions, // 传入给 ModifyStudent
  headers,
  getStatusColor,
  getStatusText,
  openEdit,
  openDelete,
  saveStudent,
  openRecharge,
  saveRecharge,
  openAdd,
  exportStudents,
  openDetails,
  loadItems,
} = useStudentManage()
</script>

<template>
  <v-sheet class="fill-height pa-6 bg-background">
    <div class="d-flex flex-column h-100">

      <!-- 1. 顶部操作栏 -->
      <div class="d-flex justify-space-between align-center mb-4 flex-shrink-0">
        <v-text-field v-model="search" prepend-inner-icon="mdi-magnify" label="搜索学生姓名" single-line hide-details
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
          <template v-slot:item.balance="{ item }">
            <span :class="item.balance < 0 ? 'text-error font-weight-bold' : 'font-weight-medium'">
              {{ item.balance }}
            </span>
          </template>

          <template v-slot:item.status="{ item }">
            <v-chip :color="getStatusColor(item.balance)" size="small" class="font-weight-bold" label variant="tonal">
              {{ getStatusText(item.balance) }}
            </v-chip>
          </template>

          <template v-slot:item.actions="{ item }">
            <div class="d-flex justify-center align-center gap-2" style="gap: 8px;">
              <v-tooltip location="top" text="充值/调整">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="primary" v-bind="props" @click="openRecharge(item)">
                    <v-icon>mdi-wallet-plus</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>

              <v-tooltip location="top" text="明细记录">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="info" v-bind="props" @click="openDetails(item)">
                    <v-icon>mdi-file-document-outline</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>

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
        :teacher-options="teacherOptions" @save="saveStudent" />

      <!-- 4. 充值弹窗 -->
      <RechargeDialog v-model="dialogRecharge" :student="rechargeItem" @save="saveRecharge" />

      <!-- 6. 明细弹窗 -->
      <DetailsDialog v-model:modelValue="dialogDetails" :student-id="editedItem.id" :student-name="editedItem.name">
      </DetailsDialog>
    </div>
  </v-sheet>
</template>

<style scoped>
/* 隐藏滚动条 */
div::-webkit-scrollbar {
  display: none;
}
</style>