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

                <v-chip v-if="selectedStudents.length > 0" closable size="small" color="primary" variant="flat"
                  @click:close="selectedStudents = []">
                  {{ selectedStudentText }}
                </v-chip>

                <v-chip v-if="selectedTeachers.length > 0" closable size="small" color="primary" variant="flat"
                  @click:close="selectedTeachers = []">
                  {{ selectedTeacherText }}
                </v-chip>

                <v-chip v-if="selectedSubjects.length > 0" closable size="small" color="primary" variant="flat"
                  @click:close="selectedSubjects = []">
                  {{ selectedSubjectText }}
                </v-chip>

                <v-chip v-if="filterDateType !== '全部时间'" closable size="small" color="primary" variant="flat"
                  @click:close="clearDateFilter">
                  时间: {{ dateRangeText }}
                </v-chip>

                <v-chip v-if="activeFilter !== null" closable size="small" color="primary" variant="flat"
                  @click:close="activeFilter = null">
                  状态: {{ activeFilter ? '已激活' : '未激活' }}
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
            <div class="d-flex align-center justify-start">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom start" offset="5">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': selectedStudents.length > 0 }"
                    :color="selectedStudents.length > 0 ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
                </template>
                <v-card min-width="260" class="pa-4 rounded-lg elevation-4">
                  <v-autocomplete v-model="selectedStudents" :items="studentOptions" item-title="title"
                    item-value="value" :loading="loadingStudents" @update:search="onStudentSearch" label="搜索学生姓名"
                    multiple chips closable-chips density="compact" variant="outlined" hide-details clearable no-filter
                    :return-object="false" placeholder="输入关键字搜索" autocomplete="off"></v-autocomplete>
                </v-card>
              </v-menu>
            </div>
          </template>

          <!-- 2. 科目列头筛选 -->
          <template v-slot:header.subject_name="{ column }">
            <div class="d-flex align-center justify-start">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom start" offset="5">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': selectedSubjects.length > 0 }"
                    :color="selectedSubjects.length > 0 ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
                </template>
                <v-card min-width="260" class="pa-4 rounded-lg elevation-4">
                  <v-autocomplete v-model="selectedSubjects" :items="subjectOptions" item-title="title"
                    item-value="value" :loading="loadingSubjects" @update:search="onSubjectSearch" label="搜索科目名称"
                    multiple chips closable-chips density="compact" variant="outlined" hide-details clearable no-filter
                    :return-object="false" placeholder="输入关键字搜索" autocomplete="off"></v-autocomplete>
                </v-card>
              </v-menu>
            </div>
          </template>

          <!-- 3. 老师列头筛选 -->
          <template v-slot:header.teacher_name="{ column }">
            <div class="d-flex align-center justify-start">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom start" offset="5">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': selectedTeachers.length > 0 }"
                    :color="selectedTeachers.length > 0 ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
                </template>
                <v-card min-width="260" class="pa-4 rounded-lg elevation-4">
                  <v-autocomplete v-model="selectedTeachers" :items="teacherOptions" item-title="title"
                    item-value="value" :loading="loadingTeachers" @update:search="onTeacherSearch" label="搜索老师姓名"
                    multiple chips closable-chips density="compact" variant="outlined" hide-details clearable no-filter
                    :return-object="false" placeholder="输入关键字搜索" autocomplete="off"></v-autocomplete>
                </v-card>
              </v-menu>
            </div>
          </template>

          <!-- 3. 日期列头筛选 -->
          <template v-slot:header.date="{ column }">
            <div class="d-flex align-center justify-start">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom start" offset="5">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': filterDateType !== '全部时间' }"
                    :color="filterDateType !== '全部时间' ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
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

          <!-- 4. 状态列头筛选 -->
          <template v-slot:header.status="{ column }">
            <div class="d-flex align-center justify-center">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom start" offset="5">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': activeFilter !== null }"
                    :color="activeFilter !== null ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
                </template>
                <v-list density="compact" nav class="rounded-lg elevation-4" width="160">
                  <v-list-subheader class="font-weight-bold text-caption">
                    激活状态
                  </v-list-subheader>
                  <v-list-item v-for="item in activeOptions" :key="String(item.value)" :value="item.value"
                    :active="activeFilter === item.value" color="primary" @click="activeFilter = item.value"
                    class="rounded mb-1">
                    <v-list-item-title class="text-body-2">
                      {{ item.title }}
                    </v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </div>
          </template>

          <!-- 表格内容插槽 -->
          <template v-slot:item.student_name="{ item }">
            <div class="d-flex align-center py-2">
              <v-icon icon="mdi-account-outline" size="small" class="mr-1 text-medium-emphasis"></v-icon>
              <span class="text-body-2">{{ item.studentName }}</span>
            </div>
          </template>

          <template v-slot:item.subject_name="{ item }">
            <v-chip size="small" variant="outlined" color="indigo" label class="font-weight-medium">
              {{ item.subjectName }}
            </v-chip>
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
      <ImportRecord v-model="dialogImport" @import-success="onImportSuccess" @import-failed="onImportFailed" />
      <DateRangeDialog v-model="dialogDateRange" :start-date="customStartDate" :end-date="customEndDate"
        @confirm="handleCustomDateConfirm" @cancel="handleCustomDateCancel" />
      <ImportErrorDialog v-model="dialogError" :errors="importErrorInfos" />
    </div>
  </v-sheet>
</template>

<script setup lang="ts">
import { onActivated } from 'vue';
import { useRecordManage } from './RecordManage.logic';
import ModifyRecord from './ModifyRecord.vue';
import ImportRecord from './ImportRecord.vue';
import DateRangeDialog from './DateRangeDialog.vue';
import ImportErrorDialog from './ImportErrorDialog.vue';


const {
  selectedStudents,
  selectedTeachers,
  selectedSubjects,
  studentOptions,
  teacherOptions,
  subjectOptions,
  loadingStudents,
  loadingTeachers,
  loadingSubjects,
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
  dialogError,
  importErrorInfos,
  pendingCount,
  hasActiveFilters,
  dateRangeText,
  selectedStudentText,
  selectedTeacherText,
  selectedSubjectText,
  customStartDate,
  customEndDate,
  activeFilter,
  activeOptions,
  loadItems,
  onStudentSearch,
  onTeacherSearch,
  onSubjectSearch,
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
  onImportSuccess,
  onImportFailed,
} = useRecordManage();

onActivated(() => {
  loadItems({ page: page.value, itemsPerPage: itemsPerPage.value, sortBy: [] });
});
</script>

<style scoped>
.header-filter-icon {
  opacity: 0.4;
  transition: opacity 0.2s;
}

.header-filter-icon:hover,
.header-filter-icon.active {
  opacity: 1;
}
</style>