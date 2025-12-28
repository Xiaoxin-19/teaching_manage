<script setup lang="ts">
import { useStudentManage } from './StudentManage.logic'
import DetailsDialog from './DetailsDialog.vue'
import ModifyStudent from './ModifyStudent.vue'

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
  rechargeForm,
  students,
  teacherOptions, // 传入给 ModifyStudent
  headers,
  rechargeType,
  rechargeColor,
  inferredTags,
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
              {{ item.balance }} 课时
            </span>
          </template>

          <template v-slot:item.status="{ item }">
            <v-chip :color="getStatusColor(item.balance)" size="small" class="font-weight-bold" label variant="tonal">
              {{ getStatusText(item.balance) }}
            </v-chip>
          </template>

          <template v-slot:item.actions="{ item }">
            <div class="d-flex justify-end">
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
              <v-btn color="primary" variant="text" class="mt-2" @click="dialog = true">点击添加第一位学生</v-btn>
            </div>
          </template>
        </v-data-table-server>
      </v-card>

      <!-- 3. 新增/编辑弹窗 -->
      <ModifyStudent v-model:modelValue="dialog" :is-edit="editedIndex >= 0" :initial-data="editedItem"
        :teacher-options="teacherOptions" @save="saveStudent" />

      <!-- 4. 充值弹窗 -->
      <v-dialog v-model="dialogRecharge" max-width="450px" transition="dialog-bottom-transition">
        <v-card class="rounded-lg elevation-4">
          <v-card-title class="d-flex justify-space-between align-center py-3 px-4">
            <div class="d-flex align-center">
              <v-icon color="primary" class="mr-3">mdi-wallet-plus</v-icon>
              <div>
                <div class="text-subtitle-1 font-weight-bold" style="line-height: 1.2;">课时调整</div>
                <div class="text-caption text-medium-emphasis mt-1 font-weight-normal">
                  操作对象: {{ rechargeItem.name }} (当前: {{ rechargeItem.balance }})
                </div>
              </div>
            </div>
            <v-btn icon="mdi-close" variant="text" size="small" density="comfortable"
              @click="dialogRecharge = false"></v-btn>
          </v-card-title>

          <v-divider></v-divider>

          <v-card-text class="pa-4 pt-6">
            <v-text-field :model-value="rechargeType" label="操作类型 (自动识别)" variant="outlined" density="comfortable"
              prepend-inner-icon="mdi-swap-horizontal" readonly class="mb-3">
              <template v-slot:append-inner>
                <v-chip :color="rechargeColor" size="x-small" label variant="flat" class="font-weight-bold">
                  {{ rechargeType }}
                </v-chip>
              </template>
            </v-text-field>

            <v-text-field v-model.number="rechargeForm.amount" label="变动数量 (正数充值，负数扣减)" type="number" variant="outlined"
              density="comfortable" prepend-inner-icon="mdi-counter" suffix="课时" class="mb-3" hint="输入正数增加课时，输入负数减少课时"
              persistent-hint></v-text-field>

            <v-text-field v-model="rechargeForm.note" label="操作备注" placeholder="例如：微信续费 / 活动赠送" variant="outlined"
              density="comfortable" prepend-inner-icon="mdi-comment-outline"
              :hint="inferredTags.length > 0 ? '' : '输入关键词如“微信”、“赠送”可自动分类'" persistent-hint class="mb-1">
            </v-text-field>

            <v-expand-transition>
              <div v-if="inferredTags.length > 0" class="d-flex align-center mt-2 px-1">
                <v-icon size="small" color="primary" class="mr-2">mdi-auto-fix</v-icon>
                <div class="text-caption text-medium-emphasis mr-2">智能识别:</div>
                <div class="d-flex align-center gap-1"
                  style="gap: 4px; overflow-x: auto; white-space: nowrap; -ms-overflow-style: none; scrollbar-width: none;">
                  <v-chip v-for="(tag, index) in inferredTags" :key="index" :color="tag.color" size="small" label
                    variant="tonal" class="font-weight-bold flex-shrink-0">
                    {{ tag.label }}
                  </v-chip>
                </div>
              </div>
            </v-expand-transition>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions class="pa-3">
            <v-spacer></v-spacer>
            <v-btn variant="text" class="mr-2" @click="dialogRecharge = false">取消</v-btn>
            <v-btn :color="rechargeColor === 'error' ? 'error' : 'primary'" variant="elevated" @click="saveRecharge"
              class="px-6" elevation="1">
              确认调整
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

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