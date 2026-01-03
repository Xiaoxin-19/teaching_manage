<template>
  <v-sheet class="fill-height pa-6 bg-background">
    <div class="d-flex flex-column h-100">

      <!-- 1. 操作栏 (右对齐，去除标题，仅保留导出) -->
      <div class="d-flex justify-end align-center mb-4 flex-shrink-0">
        <div class="d-flex align-center gap-4">
          <v-btn prepend-icon="mdi-microsoft-excel" variant="elevated" color="success" elevation="1"
            @click="handleExport">
            导出报表
          </v-btn>
        </div>
      </div>

      <!-- 2. 筛选条件展示区 Chips -->
      <div v-if="activeFilters.length > 0" class="d-flex align-center flex-wrap gap-2 mb-4 flex-shrink-0">
        <span class="text-caption text-medium-emphasis mr-1">当前筛选:</span>
        <v-chip v-for="filter in activeFilters" :key="filter.key" closable size="small" color="primary" variant="tonal"
          @click:close="clearFilter(filter.key)">
          {{ filter.label }}
        </v-chip>
        <v-btn size="x-small" variant="text" color="grey"
          @click="() => { for (const f of activeFilters) clearFilter(f.key) }">
          清除全部
        </v-btn>
      </div>

      <!-- 3. 数据表格 -->
      <v-card elevation="2" class="flex-grow-1 rounded-lg overflow-hidden border">
        <v-data-table-server v-model:expanded="expanded" :headers="headers" :items="orders" :loading="loading" hover
          :page="page" :items-per-page="itemPerPage" :items-length="totalItems" @update:options="loadItems"
          density="comfortable" class="h-100" show-expand item-value="id" fixed-header>

          <!-- === 表头筛选插槽 === -->

          <!-- 学员筛选 (v-autocomplete) -->
          <template v-slot:header.studentName="{ column }">
            <div class="d-flex align-center justify-start">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': filters.studentId }"
                    :color="filters.studentId ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
                </template>
                <v-sheet min-width="280" class="pa-4 rounded-lg" elevation="4">
                  <v-autocomplete v-model="filters.studentId" :items="studentOptions" :loading="isStudentLoading"
                    @update:search="onStudentSearch" item-title="title" item-value="value" label="搜索学员姓名/学号"
                    placeholder="输入关键词" density="compact" variant="outlined" hide-details
                    prepend-inner-icon="mdi-account-search" clearable no-filter :return-object="false"></v-autocomplete>
                  <div class="text-caption text-medium-emphasis mt-2 pl-1">* 输入文字自动从服务器搜索</div>
                </v-sheet>
              </v-menu>
            </div>
          </template>

          <!-- 科目筛选 (v-autocomplete) -->
          <template v-slot:header._subjectName="{ column }">
            <div class="d-flex align-center justify-center">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': filters.subjectIds.length > 0 }"
                    :color="filters.subjectIds.length > 0 ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
                </template>
                <v-sheet min-width="250" class="pa-4 rounded-lg" elevation="4">
                  <v-autocomplete v-model="filters.subjectIds" :items="subjectOptions" :loading="isSubjectLoading"
                    @update:search="onSubjectSearch" item-title="title" item-value="value" label="选择科目" multiple chips
                    closable-chips density="compact" variant="outlined" hide-details clearable
                    no-filter></v-autocomplete>
                </v-sheet>
              </v-menu>
            </div>
          </template>

          <!-- 类型筛选 -->
          <template v-slot:header.type="{ column }">
            <div class="d-flex align-center justify-center">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': !!filters.type }"
                    :color="filters.type ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
                </template>
                <v-sheet min-width="200" class="pa-2 rounded-lg" elevation="4">
                  <v-radio-group v-model="filters.type" hide-details density="compact" color="primary">
                    <v-radio label="全部" value=""></v-radio>
                    <v-radio v-for="t in typeOptions" :key="t.value" :label="t.title" :value="t.value"></v-radio>
                  </v-radio-group>
                </v-sheet>
              </v-menu>
            </div>
          </template>

          <!-- 时间筛选 -->
          <template v-slot:header.created_at="{ column }">
            <div class="d-flex align-center justify-end">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" location="bottom">
                <template v-slot:activator="{ props }">
                  <v-btn icon variant="text" density="compact" size="small" v-bind="props"
                    class="ml-1 header-filter-icon" :class="{ 'active': filters.dateStart || filters.dateEnd }"
                    :color="(filters.dateStart || filters.dateEnd) ? 'primary' : ''">
                    <v-icon size="16">mdi-filter-variant</v-icon>
                  </v-btn>
                </template>
                <v-sheet min-width="300" class="pa-4 rounded-lg" elevation="4">
                  <div class="text-caption mb-2 font-weight-bold">交易日期范围</div>
                  <div class="d-flex align-center gap-2">
                    <v-text-field v-model="filters.dateStart" type="date" density="compact" variant="outlined"
                      hide-details :max="filters.dateEnd"></v-text-field>
                    <span class="text-medium-emphasis">-</span>
                    <v-text-field v-model="filters.dateEnd" type="date" density="compact" variant="outlined"
                      hide-details :min="filters.dateStart"></v-text-field>
                  </div>
                  <div class="mt-3 d-flex justify-end">
                    <v-btn size="small" variant="text" @click="filters.dateStart = ''; filters.dateEnd = ''">重置</v-btn>
                  </div>
                </v-sheet>
              </v-menu>
            </div>
          </template>

          <!-- === 内容插槽 === -->

          <!-- 学员列 -->
          <template v-slot:item.studentName="{ item }">
            <span class="font-weight-medium text-body-2">{{ item.student.name }}</span>
          </template>

          <!-- 科目列 -->
          <template v-slot:item.subjectName="{ item }">
            <span class="text-body-2">{{ item.subject.name }}</span>
          </template>

          <!-- 类型列 (去色) -->
          <template v-slot:item.type="{ item }">
            <div class="d-flex justify-center">
              <v-chip size="small" :color="getTypeColor(item.type)" label variant="outlined"
                class="font-weight-regular text-caption border-0">
                {{ getTypeText(item.type) }}
              </v-chip>
            </div>
          </template>

          <!-- 标签列 (去色，后置) -->
          <template v-slot:item.tags="{ item }">
            <div class="d-flex gap-2 flex-wrap align-center justify-end">
              <template v-if="item.tags && item.tags.length > 0">
                <v-chip v-for="(tag, index) in item.tags" :key="index" :color="tag.color" size="x-small" label
                  variant="tonal" class="font-weight-medium">
                  {{ tag.label }}
                </v-chip>
              </template>
              <span v-else class="text-caption text-disabled">-</span>
            </div>
          </template>

          <!-- 课时列 -->
          <template v-slot:item.hours="{ item }">
            <div class="d-flex justify-center">
              <v-tooltip location="top">
                <template v-slot:activator="{ props }">
                  <span v-bind="props" class="font-weight-bold tabular-nums text-body-1"
                    :class="item.hours > 0 ? 'text-success' : 'text-deep-orange'">
                    {{ formatHours(item.hours) }}
                  </span>
                </template>
                <span>{{ item.hours }}</span>
              </v-tooltip>
            </div>
          </template>

          <!-- 金额列 -->
          <template v-slot:item.amount="{ item }">
            <span v-if="Math.abs(item.amount) > 0" class="font-weight-bold tabular-nums"
              :class="item.amount > 0 ? 'text-success' : 'text-deep-orange'">
              {{ formatCurrency(item.amount) }}
            </span>
            <span v-else class="text-medium-emphasis">-</span>
          </template>

          <!-- 创建时间列 -->
          <template v-slot:item.created_at="{ item }">
            <span class="text-body-2">
              <!-- 转换UnixMill为本地时间显示 -->
              {{ formatDate(item.created_at) }}
            </span>
          </template>

          <!-- 空状态 -->
          <template v-slot:no-data>
            <div class="pa-8 text-center text-medium-emphasis">
              <v-icon size="64" class="mb-2 text-disabled">mdi-file-document-off-outline</v-icon>
              <div class="text-body-1">暂无财务流水记录</div>
            </div>
          </template>

          <!-- === 展开详情 (增加学员号) === -->
          <template v-slot:expanded-row="{ columns, item }">
            <tr>
              <td :colspan="columns.length" class="pa-0">
                <div class="pa-4 border-b expanded-row-bg">
                  <v-row dense>
                    <!-- 1. 订单编号 -->
                    <v-col cols="12" sm="3">
                      <div class="text-caption text-medium-emphasis mb-1">订单编号</div>
                      <div class="text-body-2 font-family-monospace select-all d-flex align-center">
                        <v-icon size="small" class="mr-1" color="grey">mdi-barcode</v-icon>
                        {{ item.order_number }}
                      </div>
                    </v-col>

                    <!-- 2. 学员号 (新增) -->
                    <v-col cols="12" sm="3">
                      <div class="text-caption text-medium-emphasis mb-1">学员号</div>
                      <div class="text-body-2 font-family-monospace select-all d-flex align-center">
                        <v-icon size="small" class="mr-1" color="grey">mdi-card-account-details-outline</v-icon>
                        {{ item.student.student_number }}
                      </div>
                    </v-col>

                    <!-- 3. 备注说明 -->
                    <v-col cols="12" sm="6">
                      <div class="text-caption text-medium-emphasis mb-1">备注说明</div>
                      <div
                        class="text-body-2 text-pre-wrap select-all bg-surface rounded pa-2 border border-dashed text-medium-emphasis">
                        {{ item.remark || '无备注' }}
                      </div>
                    </v-col>
                  </v-row>
                </div>
              </td>
            </tr>
          </template>

        </v-data-table-server>
      </v-card>
    </div>
  </v-sheet>
</template>

<script setup lang="ts">
import { useFinanceManage } from './FincnaceManage.logic'

const {
  page,
  itemPerPage,
  totalItems,

  loading,
  filters,
  orders,
  headers,
  expanded,

  studentOptions,
  isStudentLoading,
  onStudentSearch,
  subjectOptions,
  isSubjectLoading,
  onSubjectSearch,
  typeOptions,

  activeFilters,
  loadItems,
  clearFilter,
  getTypeColor,
  getTypeText,
  formatCurrency,
  formatHours,
  formatDate,
  handleExport
} = useFinanceManage()
</script>

<style scoped>
/* 隐藏滚动条但保留功能 */
div::-webkit-scrollbar {
  display: none;
}

.gap-2 {
  gap: 8px;
}

.gap-4 {
  gap: 16px;
}

/* 数字等宽对齐 */
.tabular-nums {
  font-feature-settings: "tnum";
  font-variant-numeric: tabular-nums;
}

.font-family-monospace {
  font-family: monospace;
}

.select-all {
  user-select: all;
}

.text-pre-wrap {
  white-space: pre-wrap;
}

/* 表头筛选图标样式 */
.header-filter-icon {
  opacity: 0.4;
  transition: opacity 0.2s;
}

.header-filter-icon:hover,
.header-filter-icon.active {
  opacity: 1;
}

/* 展开行背景色适配 */
.expanded-row-bg {
  background-color: #f5f5f5;
  /* Light mode */
}

/* Vuetify 暗黑模式自动适配类名 */
:global(.v-theme--dark) .expanded-row-bg {
  background-color: #1E1E1E;
  /* Dark mode */
}
</style>