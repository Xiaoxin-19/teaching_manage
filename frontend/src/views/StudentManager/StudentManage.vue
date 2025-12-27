<script setup lang="ts">

import { useStudentManage } from './StudentManage.logic'
import DetailsDialog from './DetailsDialog.vue'

const {
  search,
  dialog,
  dialogRecharge,
  dialogDelete,
  dialogDetails,
  editedIndex,
  editedItem,
  rechargeItem,
  rechargeForm,
  students,
  headers,
  getStatusColor,
  getStatusText,
  openEdit,
  openDelete,
  deleteItemConfirm,
  closeDialog,
  closeDelete,
  saveStudent,
  openRecharge,
  saveRecharge,
} = useStudentManage()
</script>

<template>
  <v-sheet class="h-100">
    <div class="d-flex flex-column pa-6">

      <!-- 1. 顶部操作栏 -->
      <div class="d-flex justify-end align-center mb-4 flex-shrink-0">

        <div class="d-flex align-center gap-2">
          <v-text-field v-model="search" prepend-inner-icon="mdi-magnify" label="搜索姓名或电话" single-line hide-details
            density="compact" variant="outlined" bg-color="surface" style="width: 260px" class="mr-3"></v-text-field>

          <v-btn color="primary" prepend-icon="mdi-plus" elevation="1" height="40" @click="dialog = true">
            新增学生
          </v-btn>
        </div>
      </div>

      <!-- 2. 数据表格区域 -->
      <v-card elevation="1" class="flex-grow-1 d-flex flex-column rounded-lg overflow-hidden border">
        <v-data-table :headers="headers" :items="students" :search="search" hover fixed-header density="default"
          class="h-100">
          <!-- 自定义列：剩余课时 (高亮负数) -->
          <template v-slot:item.balance="{ item }">
            <span :class="item.balance < 0 ? 'text-error font-weight-bold' : ''">
              {{ item.balance }} 课时
            </span>
          </template>

          <!-- 自定义列：状态 (Chip 标签) -->
          <template v-slot:item.status="{ item }">
            <v-chip :color="getStatusColor(item.balance)" size="small" class="font-weight-medium" label variant="flat">
              {{ getStatusText(item.balance) }}
            </v-chip>
          </template>

          <!-- 自定义列：操作按钮组 -->
          <template v-slot:item.actions="{ item }">
            <div class="d-flex justify-end">
              <v-tooltip location="top" text="课时充值/调整">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="primary" v-bind="props" @click="openRecharge(item)">
                    <v-icon>mdi-cash-plus</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>
              <v-tooltip location="top" text="明细记录">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="primary" v-bind="props"
                    @click="editedItem.id = item.id; dialogDetails = true">
                    <v-icon>mdi-cash-plus</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>

              <v-tooltip location="top" text="编辑信息">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="grey-darken-1" v-bind="props" @click="openEdit(item)">
                    <v-icon>mdi-pencil</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>

              <v-tooltip location="top" text="删除档案">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="error" v-bind="props" @click="openDelete(item)">
                    <v-icon>mdi-delete</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>
            </div>
          </template>

          <!-- 空状态展示 -->
          <template v-slot:no-data>
            <div class="pa-8 text-center">
              <v-icon size="64" color="grey-lighten-2" class="mb-2">mdi-account-off</v-icon>
              <div class="text-body-1 text-grey">暂无学生数据</div>
              <v-btn color="primary" variant="text" class="mt-2" @click="dialog = true">点击添加第一位学生</v-btn>
            </div>
          </template>
        </v-data-table>
      </v-card>

      <!-- 3. 弹窗组件：新增/编辑学生 -->
      <ModifyStudent v-model:modelValue="dialog" :is-edit="editedIndex >= 0" :initial-data="editedItem"
        @save="saveStudent" @update:modelValue="dialog = $event">
      </ModifyStudent>

      <!-- 4. 弹窗组件：课时充值 -->
      <v-dialog v-model="dialogRecharge" max-width="450px">
        <v-card class="rounded-lg">
          <v-card-title class="pa-4 d-flex align-center bg-surface-variant">
            <v-icon color="primary" class="mr-2">mdi-wallet-plus</v-icon>
            <span class="text-h6">课时调整</span>
          </v-card-title>

          <v-card-text class="pt-6">
            <div class="mb-4 text-body-2 bg-grey-lighten-4 pa-3 rounded">
              正在操作对象: <span class="font-weight-bold text-primary text-h6 ml-1">{{ rechargeItem.name }}</span>
              <div class="mt-1 text-grey-darken-1">当前余额: {{ rechargeItem.balance }} 课时</div>
            </div>

            <v-select v-model="rechargeForm.type" :items="['充值', '赠送', '退费']" label="调整类型" variant="outlined"
              density="comfortable"></v-select>
            <v-text-field v-model.number="rechargeForm.amount" label="数量" type="number" variant="outlined"
              density="comfortable" :prefix="rechargeForm.type === '退费' ? '-' : '+'" suffix="课时"></v-text-field>
            <v-text-field v-model="rechargeForm.note" label="操作备注" placeholder="例如：2024春季续费" variant="outlined"
              density="comfortable"></v-text-field>
          </v-card-text>

          <v-card-actions class="pa-4 pt-0">
            <v-spacer></v-spacer>
            <v-btn color="grey" variant="text" @click="dialogRecharge = false">取消</v-btn>
            <v-btn color="success" variant="elevated" @click="saveRecharge" class="px-6">确认调整</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <!-- 5. 弹窗组件：删除确认 -->
      <v-dialog v-model="dialogDelete" max-width="400px">
        <v-card class="rounded-lg">
          <v-card-title class="text-h6 pa-4 bg-error text-white d-flex align-center">
            <v-icon class="mr-2">mdi-alert</v-icon>
            确认删除
          </v-card-title>
          <v-card-text class="pa-4 text-body-1">
            您确定要删除学生 <span class="font-weight-bold">{{ editedItem.name }}</span> 的档案吗？
            <div class="text-caption text-error mt-2">此操作不可撤销，并将删除其所有历史记录。</div>
          </v-card-text>
          <v-card-actions class="pa-4">
            <v-spacer></v-spacer>
            <v-btn color="grey-darken-1" variant="text" @click="closeDelete">再想想</v-btn>
            <v-btn color="error" variant="elevated" @click="deleteItemConfirm">确认删除</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <!-- 弹窗组件：课程充值明细（可导出Excel） -->
      <DetailsDialog v-model:modelValue="dialogDetails" :student-id="editedItem.id"></DetailsDialog>
    </div>
  </v-sheet>
</template>