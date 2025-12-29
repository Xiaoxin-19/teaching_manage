<template>
  <v-sheet class="fill-height pa-6 bg-background">
    <div class="d-flex flex-column h-100">
      <!-- 顶部操作栏 -->
      <div class="d-flex justify-end align-center mb-4">
        <!-- 左侧提示信息 (可选) -->
        <div v-if="pendingCount > 0" class="d-flex align-center mr-auto text-warning font-weight-medium">
          <v-icon icon="mdi-alert-circle-outline" class="mr-2"></v-icon>
          <span>当前有 {{ pendingCount }} 条记录待生效</span>
        </div>
        <v-spacer v-else></v-spacer>

        <div class="d-flex align-center" style="gap: 12px">
          <!-- 批量生效：Warning (橙色) + Tonal -->
          <v-btn v-if="pendingCount > 0" prepend-icon="mdi-check-all" color="warning" variant="tonal"
            @click="processAllPending">
            一键生效 ({{ pendingCount }})
          </v-btn>

          <v-divider vertical class="mx-1" v-if="pendingCount > 0"></v-divider>

          <!-- 导出记录：Success (绿色) + Outlined -->
          <v-btn prepend-icon="mdi-microsoft-excel" color="success" variant="outlined" @click="exportRecords">
            导出记录
          </v-btn>

          <!-- 批量导入：Info (蓝色) + Outlined -->
          <v-btn prepend-icon="mdi-upload" color="info" variant="outlined" @click="dialogImport = true">
            批量导入
          </v-btn>

          <!-- 记一笔：Primary (主色) + Flat -->
          <v-btn prepend-icon="mdi-plus" color="primary" variant="flat" @click="openAdd">
            记一笔
          </v-btn>
        </div>
      </div>

      <!-- 数据表格区域 -->
      <v-card elevation="2" class="flex-grow-1 d-flex flex-column rounded-lg overflow-hidden border">
        <v-data-table-server v-model:items-per-page="itemsPerPage" v-model:page="page" :items="serverItems"
          :items-length="totalItems" :headers="headers" :loading="loading" loading-text="正在加载数据..." hover fixed-header
          density="comfortable" class="h-100" show-expand item-value="id" @update:options="loadItems">
          <!-- 顶部插槽：显示生效的筛选条件 -->
          <template v-slot:top>
            <v-expand-transition>
              <div v-if="hasActiveFilters" class="px-4 py-3 border-b d-flex align-center flex-wrap"
                style="gap: 8px; background-color: rgba(var(--v-theme-surface-variant), 0.1);">
                <div class="text-caption font-weight-bold mr-2 text-medium-emphasis d-flex align-center">
                  <v-icon size="small" class="mr-1">mdi-filter-variant</v-icon>
                  当前筛选:
                </div>

                <v-chip v-if="searchStudent" closable size="small" color="primary" variant="flat"
                  @click:close="searchStudent = ''">
                  学生: {{ searchStudent }}
                </v-chip>

                <v-chip v-if="searchTeacher" closable size="small" color="primary" variant="flat"
                  @click:close="searchTeacher = ''">
                  老师: {{ searchTeacher }}
                </v-chip>

                <v-chip v-if="filterDateType !== '全部时间'" closable size="small" color="primary" variant="flat"
                  @click:close="clearDateFilter">
                  时间: {{ dateRangeText }}
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

          <!-- 1. 学生列头筛选 -->
          <template v-slot:header.student_name="{ column }">
            <div class="header-filter-container">
              <span class="font-weight-bold mr-2">{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom start" offset="5">
                <template v-slot:activator="{ props }">
                  <v-icon v-bind="props" :icon="searchStudent ? 'mdi-filter' : 'mdi-filter-outline'" size="small"
                    class="filter-icon" :class="{ active: searchStudent }"></v-icon>
                </template>
                <v-card min-width="260" class="pa-4 rounded-lg elevation-4">
                  <div class="text-subtitle-2 mb-3 font-weight-bold d-flex align-center">
                    <v-icon size="small" class="mr-2" color="primary">
                      mdi-account-school
                    </v-icon>
                    筛选学生姓名
                  </div>
                  <v-text-field v-model="searchStudent" placeholder="输入关键字 (如: 张)" density="compact" variant="outlined"
                    hide-details autofocus prepend-inner-icon="mdi-magnify" clearable
                    @click:clear="searchStudent = ''"></v-text-field>
                  <div class="text-caption text-medium-emphasis mt-2">
                    * 支持模糊搜索
                  </div>
                </v-card>
              </v-menu>
            </div>
          </template>

          <!-- 2. 老师列头筛选 -->
          <template v-slot:header.teacher_name="{ column }">
            <div class="header-filter-container">
              <span class="font-weight-bold mr-2">{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom start" offset="5">
                <template v-slot:activator="{ props }">
                  <v-icon v-bind="props" :icon="searchTeacher ? 'mdi-filter' : 'mdi-filter-outline'" size="small"
                    class="filter-icon" :class="{ active: searchTeacher }"></v-icon>
                </template>
                <v-card min-width="260" class="pa-4 rounded-lg elevation-4">
                  <div class="text-subtitle-2 mb-3 font-weight-bold d-flex align-center">
                    <v-icon size="small" class="mr-2" color="primary">
                      mdi-account-tie
                    </v-icon>
                    筛选老师姓名
                  </div>
                  <v-text-field v-model="searchTeacher" placeholder="输入关键字 (如: 王)" density="compact" variant="outlined"
                    hide-details autofocus prepend-inner-icon="mdi-magnify" clearable
                    @click:clear="searchTeacher = ''"></v-text-field>
                  <div class="text-caption text-medium-emphasis mt-2">
                    * 支持模糊搜索
                  </div>
                </v-card>
              </v-menu>
            </div>
          </template>

          <!-- 3. 日期列头筛选 -->
          <template v-slot:header.date="{ column }">
            <div class="header-filter-container">
              <span class="font-weight-bold mr-2">{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom start" offset="5">
                <template v-slot:activator="{ props }">
                  <v-icon v-bind="props" :icon="filterDateType !== '全部时间'
                    ? 'mdi-filter'
                    : 'mdi-filter-outline'
                    " size="small" class="filter-icon" :class="{ active: filterDateType !== '全部时间' }"></v-icon>
                </template>
                <v-list density="compact" nav class="rounded-lg elevation-4" width="160">
                  <v-list-subheader class="font-weight-bold text-caption">
                    时间范围
                  </v-list-subheader>
                  <v-list-item v-for="item in dateOptions" :key="item" :value="item" :active="filterDateType === item"
                    color="primary" @click="selectDateFilter(item)" class="rounded mb-1">
                    <v-list-item-title class="text-body-2">
                      {{ item }}
                    </v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </div>
          </template>

          <!-- 表格内容插槽 -->
          <template v-slot:item.student_name="{ item }">
            <div class="d-flex align-center py-2">
              <v-avatar color="primary" variant="tonal" size="32" class="mr-3">
                <span class="text-subtitle-2 font-weight-bold">
                  {{ item.studentName.charAt(0) }}
                </span>
              </v-avatar>
              <div>
                <div class="font-weight-medium text-body-2">
                  {{ item.studentName }}
                </div>
              </div>
            </div>
          </template>

          <template v-slot:item.teacher_name="{ item }">
            <div class="d-flex align-center">
              <v-icon icon="mdi-account-tie-outline" size="small" class="mr-1 text-medium-emphasis"></v-icon>
              <span class="text-body-2">{{ item.teacherName }}</span>
            </div>
          </template>

          <!-- 状态列插槽 -->
          <template v-slot:item.status="{ item }">
            <v-chip :color="item.status === 'active' ? 'success' : 'warning'" size="small" variant="flat" label
              class="font-weight-medium">
              <v-icon start size="x-small">
                {{
                  item.status === 'active'
                    ? 'mdi-check-circle'
                    : 'mdi-clock-outline'
                }}
              </v-icon>
              {{ item.status === 'active' ? '已生效' : '未生效' }}
            </v-chip>
          </template>

          <template v-slot:item.date="{ item }">
            <span class="font-weight-medium">{{ item.date }}</span>
          </template>

          <template v-slot:item.time="{ item }">
            <v-sheet class="d-inline-block px-2 py-1 rounded text-caption font-weight-bold" color="grey-lighten-3">
              {{ item.time }}
            </v-sheet>
          </template>

          <template v-slot:item.actions="{ item }">
            <div class="d-flex justify-end align-center">
              <v-tooltip v-if="item.status !== 'active'" location="top" text="立即生效">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="success" v-bind="props"
                    @click.stop="activateRecord(item)">
                    <v-icon>mdi-check</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>

              <v-tooltip location="top" text="撤销/删除记录">
                <template v-slot:activator="{ props }">
                  <v-btn icon size="small" variant="text" color="error" v-bind="props" @click.stop="deleteItem(item)">
                    <v-icon>mdi-delete-outline</v-icon>
                  </v-btn>
                </template>
              </v-tooltip>
            </div>
          </template>

          <template v-slot:expanded-row="{ columns, item }">
            <tr>
              <td :colspan="columns.length" class="pa-0">
                <div class="pa-4 d-flex align-start"
                  style="background-color: rgba(var(--v-theme-surface-variant), 0.3); border-left: 4px solid #1976D2;">
                  <div class="mr-4 mt-1">
                    <v-icon color="primary" size="small">
                      mdi-comment-quote-outline
                    </v-icon>
                  </div>
                  <div class="flex-grow-1">
                    <div class="text-caption font-weight-bold text-primary mb-1">
                      课程备注
                    </div>
                    <div class="text-body-2" style="white-space: pre-wrap; line-height: 1.6; opacity: 0.9;">
                      {{ item.remark || '暂无详细备注信息' }}
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </template>

          <template v-slot:no-data>
            <div class="pa-8 text-center text-medium-emphasis">
              <v-icon size="64" class="mb-2 text-disabled">
                mdi-notebook-off-outline
              </v-icon>
              <div class="text-body-1">暂无教学记录</div>
            </div>
          </template>
        </v-data-table-server>
      </v-card>

      <!-- 弹窗组件 -->
      <ModifyRecord v-model="dialogForm" @save="saveRecord" />
      <ImportRecord v-model="dialogImport" />
      <DateRangeDialog v-model="dialogDateRange" :start-date="customStartDate" :end-date="customEndDate"
        @confirm="handleCustomDateConfirm" @cancel="handleCustomDateCancel" />
    </div>
  </v-sheet>
</template>

<script setup lang="ts">
import { useRecordManage } from './RecordManage.logic';
import ModifyRecord from './ModifyRecord.vue';
import ImportRecord from './ImportRecord.vue';
import DateRangeDialog from './DateRangeDialog.vue';

const {
  searchStudent,
  searchTeacher,
  filterDateType,
  dateOptions,
  page,
  itemsPerPage,
  totalItems,
  loading,
  serverItems,
  headers,
  dialogForm,
  dialogImport,
  dialogDateRange,
  mockStudents,
  pendingCount,
  hasActiveFilters,
  dateRangeText,
  customStartDate,
  customEndDate,
  loadItems,
  selectDateFilter,
  handleCustomDateConfirm,
  handleCustomDateCancel,
  clearAllFilters,
  clearDateFilter,
  openAdd,
  saveRecord,
  activateRecord,
  processAllPending,
  deleteItem,
  exportRecords,
} = useRecordManage();
</script>

<style scoped>
.header-filter-container {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
  border-radius: 4px;
  padding: 2px 4px;
  transition: background-color 0.2s;
}

.header-filter-container:hover {
  background-color: rgba(var(--v-theme-on-surface), 0.05);
}

.filter-icon {
  opacity: 0.6;
  transition: all 0.2s;
}

.header-filter-container:hover .filter-icon {
  opacity: 1;
  color: rgba(var(--v-theme-on-surface), 0.8);
}

.filter-icon.active {
  opacity: 1 !important;
  color: rgb(var(--v-theme-primary)) !important;
  transform: scale(1.1);
}
</style>