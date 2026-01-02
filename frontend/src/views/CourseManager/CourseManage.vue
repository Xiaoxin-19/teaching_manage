<template>
  <v-sheet class="fill-height pa-6 bg-background">
    <div class="d-flex flex-column h-100">

      <!-- 1. 顶部工具栏 -->
      <div class="d-flex justify-space-between align-center mb-4">
        <div>
          <h2 class="text-h5 font-weight-bold text-high-emphasis">教务中心</h2>
        </div>
        <div class="d-flex align-center">
          <v-btn color="primary" prepend-icon="mdi-book-plus" elevation="2" @click="openEnroll">新课报名</v-btn>
        </div>
      </div>

      <!-- 筛选条件展示区 -->
      <div v-if="activeFilters.length > 0" class="d-flex align-center flex-wrap gap-2 mb-4">
        <span class="text-caption text-medium-emphasis mr-1">筛选条件:</span>
        <v-chip v-for="filter in activeFilters" :key="filter.key" closable size="small" color="primary" variant="tonal"
          @click:close="clearFilter(filter.key)">{{ filter.label }}</v-chip>
        <v-btn size="x-small" variant="text" color="grey"
          @click="() => { for (const f of activeFilters) clearFilter(f.key) }">清除全部</v-btn>
      </div>

      <!-- 2. 数据表格 -->
      <v-card elevation="2" class="flex-grow-1 rounded-lg overflow-hidden border">
        <v-data-table-server :headers="headers" :items="courses" :items-length="totalItems" :loading="loading"
          v-model:page="page" v-model:items-per-page="itemsPerPage" @update:options="loadData" hover
          density="comfortable" class="h-100">
          <!-- 表头筛选插槽 -->
          <template v-slot:header.studentName="{ column }">
            <div class="d-flex align-center justify-start"><span>{{ column.title }}</span><v-menu
                :close-on-content-click="false" location="bottom"><template v-slot:activator="{ props }"><v-btn icon
                    variant="text" density="compact" size="small" v-bind="props" class="ml-1 header-filter-icon"
                    :class="{ 'active': filters.studentName }" :color="filters.studentName ? 'primary' : ''"><v-icon
                      size="16">mdi-filter-variant</v-icon></v-btn></template><v-sheet min-width="250"
                  class="pa-4 rounded-lg" elevation="4"><v-text-field v-model="filters.studentName" label="搜索姓名或学号"
                    density="compact" variant="outlined" hide-details prepend-inner-icon="mdi-magnify"
                    clearable></v-text-field></v-sheet></v-menu></div>
          </template>
          <template v-slot:header.subjectName="{ column }">
            <div class="d-flex align-center justify-center"><span>{{ column.title }}</span><v-menu
                :close-on-content-click="false" location="bottom"><template v-slot:activator="{ props }"><v-btn icon
                    variant="text" density="compact" size="small" v-bind="props" class="ml-1 header-filter-icon"
                    :class="{ 'active': filters.subjects.length > 0 }"
                    :color="filters.subjects.length > 0 ? 'primary' : ''"><v-icon
                      size="16">mdi-filter-variant</v-icon></v-btn></template><v-sheet min-width="250"
                  class="pa-4 rounded-lg" elevation="4"><v-autocomplete v-model="filters.subjects"
                    :items="subjectOptions" item-title="title" item-value="title" :loading="isSubjectLoading"
                    @update:search="onSubjectSearch" label="搜索或选择科目" multiple chips closable-chips density="compact"
                    variant="outlined" hide-details clearable no-filter
                    :return-object="false"></v-autocomplete></v-sheet></v-menu></div>
          </template>
          <template v-slot:header.teacherName="{ column }">
            <div class="d-flex align-center justify-center"><span>{{ column.title }}</span><v-menu
                :close-on-content-click="false" location="bottom"><template v-slot:activator="{ props }"><v-btn icon
                    variant="text" density="compact" size="small" v-bind="props" class="ml-1 header-filter-icon"
                    :class="{ 'active': filters.teachers.length > 0 }"
                    :color="filters.teachers.length > 0 ? 'primary' : ''"><v-icon
                      size="16">mdi-filter-variant</v-icon></v-btn></template><v-sheet min-width="250"
                  class="pa-4 rounded-lg" elevation="4"><v-autocomplete v-model="filters.teachers"
                    :items="teacherOptions" item-title="title" item-value="title" :loading="isTeacherLoading"
                    @update:search="onTeacherSearch" label="搜索或选择老师" multiple chips closable-chips density="compact"
                    variant="outlined" hide-details clearable no-filter
                    :return-object="false"></v-autocomplete></v-sheet></v-menu></div>
          </template>
          <template v-slot:header.balance="{ column }">
            <div class="d-flex align-center justify-center"><span>{{ column.title }}</span><v-menu
                :close-on-content-click="false" location="bottom"><template v-slot:activator="{ props }"><v-btn icon
                    variant="text" density="compact" size="small" v-bind="props" class="ml-1 header-filter-icon"
                    :class="{ 'active': filters.balanceMin || filters.balanceMax }"
                    :color="(filters.balanceMin || filters.balanceMax) ? 'primary' : ''"><v-icon
                      size="16">mdi-filter-variant</v-icon></v-btn></template><v-sheet min-width="280"
                  class="pa-4 rounded-lg" elevation="4">
                  <div class="text-caption mb-2">课时范围筛选</div>
                  <div class="d-flex align-center gap-2"><v-text-field v-model="filters.balanceMin" type="number"
                      label="最小" density="compact" variant="outlined" hide-details></v-text-field><span
                      class="text-medium-emphasis">-</span><v-text-field v-model="filters.balanceMax" type="number"
                      label="最大" density="compact" variant="outlined" hide-details></v-text-field></div>
                  <div class="mt-2 d-flex justify-end"><v-btn size="x-small" variant="text"
                      @click="filters.balanceMin = null; filters.balanceMax = null">重置</v-btn></div>
                </v-sheet></v-menu></div>
          </template>
          <template v-slot:header.status="{ column }">
            <div class="d-flex align-center justify-center"><span>{{ column.title }}</span><v-menu
                :close-on-content-click="false" location="bottom"><template v-slot:activator="{ props }"><v-btn icon
                    variant="text" density="compact" size="small" v-bind="props" class="ml-1 header-filter-icon"
                    :class="{ 'active': filters.status.length > 0 }"
                    :color="filters.status.length > 0 ? 'primary' : ''"><v-icon
                      size="16">mdi-filter-variant</v-icon></v-btn></template><v-sheet min-width="220"
                  class="pa-2 rounded-lg" elevation="4"><v-list density="compact"
                    select-strategy="leaf"><v-list-subheader class="text-caption">选择状态</v-list-subheader><template
                      v-for="s in statusOptions" :key="s.value"><v-checkbox v-model="filters.status" :label="s.title"
                        :value="s.value" density="compact" hide-details
                        color="primary"></v-checkbox></template></v-list></v-sheet></v-menu>
            </div>
          </template>

          <!-- 表格内容 -->
          <template v-slot:item.studentName="{ item }">
            <span class="font-weight-bold text-body-2">{{ item.studentName }}</span>
            <span class="text-caption text-medium-emphasis ml-1">- {{ item.studentCode }}</span>
          </template>
          <template v-slot:item.subjectName="{ item }">
            <v-chip size="small" variant="outlined" color="indigo" label class="font-weight-medium fixed-chip"
              :class="{ 'text-shrink': (item.subjectName || '').length > 4 }">{{ item.subjectName }}</v-chip>
          </template>
          <template v-slot:item.teacherName="{ item }">
            <div class="d-flex align-center justify-center"><span class="text-medium-emphasis">{{ item.teacherName
                }}</span><span class="text-caption text-disabled ml-1">- {{ item.teacherCode }}</span></div>
          </template>
          <template v-slot:item.balance="{ item }">
            <div class="d-flex align-center justify-center"><v-tooltip location="top"><template
                  v-slot:activator="{ props }"><span v-bind="props"
                    class="font-weight-bold text-subtitle-1 tabular-nums cursor-help"
                    :class="'text-' + getBalanceColor(item.balance)">{{ formatBalance(item.balance)
                    }}</span></template><span>{{
                      item.balance }} 节 - {{ getBalanceLabel(item.balance) }}</span></v-tooltip></div>
          </template>
          <template v-slot:item.status="{ item }">
            <v-tooltip location="top" :disabled="!getEffectiveStatus(item).disabled && item.courseStatus !== 2">
              <template v-slot:activator="{ props }">
                <div v-bind="props" class="d-flex align-center justify-center"><v-chip size="small"
                    :color="getEffectiveStatus(item).color" variant="tonal" label class="font-weight-medium fixed-chip"
                    :class="{ 'text-shrink': (getEffectiveStatus(item).label || '').length > 4 }">{{
                      getEffectiveStatus(item).label }}</v-chip></div>
              </template>
              <span>{{ getEffectiveStatus(item).desc }}</span>
            </v-tooltip>
          </template>
          <template v-slot:item.actions="{ item }">
            <div class="d-flex justify-end align-center gap-2">
              <v-btn size="small" color="primary" variant="text" prepend-icon="mdi-wallet-plus"
                @click="openRecharge(item)" :disabled="getEffectiveStatus(item).label === '学员退学'">续费</v-btn>
              <v-menu location="bottom end"><template v-slot:activator="{ props }"><v-btn icon="mdi-dots-vertical"
                    variant="text" size="small" v-bind="props" color="medium-emphasis"></v-btn></template><v-list
                  density="compact" elevation="3" class="rounded-lg">
                  <v-list-item @click="openEdit(item)" color="primary"
                    :disabled="getEffectiveStatus(item).disabled && item.courseStatus !== 2"><template
                      v-slot:prepend><v-icon
                        color="primary">mdi-account-switch-outline</v-icon></template><v-list-item-title
                      class="text-body-2">更换老师</v-list-item-title></v-list-item>
                  <v-list-item @click="openDeduction(item)" color="warning"
                    :disabled="getEffectiveStatus(item).label === '学员退学'"><template v-slot:prepend><v-icon
                        color="warning">mdi-cash-minus</v-icon></template><v-list-item-title
                      class="text-body-2 text-warning">退费/扣课时</v-list-item-title></v-list-item>
                  <v-list-item @click="toggleStatus(item)" :disabled="item.studentStatus !== 1"><template
                      v-slot:prepend><v-icon :color="item.courseStatus === 1 ? 'warning' : 'success'">{{
                        item.courseStatus === 1 ? 'mdi-pause-circle-outline' : 'mdi-play-circle-outline'
                      }}</v-icon></template><v-list-item-title class="text-body-2"
                      :class="item.courseStatus === 1 ? 'text-warning' : 'text-success'">{{ item.courseStatus === 1 ?
                        '办理停课' : '恢复上课' }}</v-list-item-title><template v-if="item.studentStatus !== 1"
                      v-slot:append><v-icon size="x-small" color="grey">mdi-lock</v-icon></template></v-list-item>
                  <v-divider class="my-1"></v-divider>
                  <v-list-item @click="openDelete(item)" color="warning"><template v-slot:prepend><v-icon
                        color="warning">mdi-archive-remove-outline</v-icon></template><v-list-item-title
                      class="text-body-2 text-warning">办理退课</v-list-item-title></v-list-item>
                  <v-list-item @click="openForceDelete(item)" color="error"><template v-slot:prepend><v-icon
                        color="error">mdi-delete-forever</v-icon></template><v-list-item-title
                      class="text-body-2 text-error">彻底删除</v-list-item-title></v-list-item>
                </v-list></v-menu>
            </div>
          </template>
        </v-data-table-server>
      </v-card>

      <!-- 1. RechargeDialog (独立组件) -->
      <RechargeDialog ref="rechargeDialogRef" v-model="rechargeDialogVisible" :course="currentItem"
        :initial-mode="rechargeMode" @submit="handleRechargeSubmit" />

      <!-- 2. Delete Dialog -->
      <DeleteCourseDialog v-model="deleteDialogVisible" :course="currentItem" :remark="deleteForm.remark"
        @update:remark="deleteForm.remark = $event" :is-valid="isDeleteValid === true" @confirm="handleDeleteConfirm" />

      <!-- 3. Force Delete Dialog -->
      <ForceDeleteCourseDialog v-model="forceDeleteDialogVisible" :course="currentItem"
        @confirm="handleForceDeleteConfirm" />

      <!-- 4. Edit Dialog -->
      <CourseEditDialog v-model="dialogVisible" :is-edit="isEdit" :course="currentItem" :form="enrollForm"
        :student-options="studentOptions" :subject-options="subjectOptions" :teacher-options="teacherOptions"
        :student-loading="isStudentLoading" :subject-loading="isSubjectLoading" :teacher-loading="isTeacherLoading"
        @search-student="onStudentSearch" @search-subject="onSubjectSearch" @search-teacher="onTeacherSearch"
        @save="handleSave" />
    </div>
  </v-sheet>
</template>

<script setup lang="ts">
import { useCourseManage } from './CourseManage.logic'
import RechargeDialog from './RechargeDialog.vue'
import DeleteCourseDialog from './DeleteCourseDialog.vue'
import ForceDeleteCourseDialog from './ForceDeleteCourseDialog.vue'
import CourseEditDialog from './CourseEditDialog.vue'

const {
  loading, page, itemsPerPage, totalItems, courses, headers,
  dialogVisible, rechargeDialogVisible, deleteDialogVisible, forceDeleteDialogVisible,
  isEdit, currentItem,
  rechargeMode, rechargeDialogRef, // 导出 ref
  deleteForm, isDeleteValid, enrollForm,
  filters, activeFilters,
  studentOptions, subjectOptions, teacherOptions, statusOptions,
  isStudentLoading, isSubjectLoading, isTeacherLoading,
  onStudentSearch, onSubjectSearch, onTeacherSearch,
  loadData, clearFilter,
  openEnroll, openEdit, handleSave,
  openRecharge, openDeduction, handleRechargeSubmit,
  toggleStatus,
  openDelete, handleDeleteConfirm,
  openForceDelete, handleForceDeleteConfirm,
  getEffectiveStatus, getBalanceColor, getBalanceLabel, formatBalance
} = useCourseManage()
</script>

<style scoped>
.tabular-nums {
  font-feature-settings: "tnum";
  font-variant-numeric: tabular-nums;
}

.cursor-help {
  cursor: help;
}

.fixed-chip {
  width: 76px !important;
  justify-content: center !important;
  padding-left: 0 !important;
  padding-right: 0 !important;
}

.text-shrink {
  font-size: 10px !important;
  letter-spacing: -0.5px !important;
}

.header-filter-icon {
  opacity: 0.4;
  transition: opacity 0.2s;
}

.header-filter-icon:hover,
.header-filter-icon.active {
  opacity: 1;
}
</style>