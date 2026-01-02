<template>
  <v-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)" max-width="500">
    <v-card class="rounded-lg">
      <v-card-title class="d-flex justify-space-between align-center py-3 px-4 bg-surface">
        <div class="d-flex align-center">
          <v-icon :icon="isEdit ? 'mdi-notebook-edit-outline' : 'mdi-notebook-plus-outline'" color="primary"
            class="mr-2"></v-icon>
          <span class="text-subtitle-1 font-weight-bold">{{ isEdit ? '编辑课程信息' : '新课报名' }}</span>
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" @click="$emit('update:modelValue', false)"></v-btn>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text class="pa-6">
        <v-autocomplete v-if="!isEdit" label="选择学员" :items="studentOptions" item-title="title" item-value="value"
          :loading="studentLoading" @update:search="$emit('searchStudent', $event)" variant="outlined"
          density="comfortable" prepend-inner-icon="mdi-account" placeholder="输入姓名搜索" class="mb-2"
          v-model="form.studentId" no-filter :return-object="false"></v-autocomplete>
        <div v-else class="mb-4 text-body-2 text-medium-emphasis d-flex align-center bg-grey-lighten-4 pa-3 rounded">
          <v-icon start size="small">mdi-account</v-icon>学员：<span class="font-weight-bold text-high-emphasis ml-1">{{
            course.student?.name }}</span>
        </div>
        <v-row dense>
          <v-col cols="12" sm="6">
            <v-autocomplete label="选择科目" :items="subjectOptions" item-title="title" item-value="value"
              :loading="subjectLoading" @update:search="$emit('searchSubject', $event)" variant="outlined"
              density="comfortable" prepend-inner-icon="mdi-bookshelf" :disabled="isEdit" v-model="form.subjectId"
              no-filter :return-object="false"></v-autocomplete>
          </v-col>
          <v-col cols="12" sm="6">
            <v-autocomplete label="授课老师" :items="teacherOptions" item-title="title" item-value="value"
              :loading="teacherLoading" @update:search="$emit('searchTeacher', $event)" variant="outlined"
              density="comfortable" prepend-inner-icon="mdi-account-tie" v-model="form.teacherId" no-filter
              :return-object="false"></v-autocomplete>
          </v-col>
        </v-row>
        <v-textarea label="备注" rows="2" variant="outlined" density="comfortable"
          prepend-inner-icon="mdi-comment-outline" class="mt-2" v-model="form.remark"></v-textarea>
      </v-card-text>
      <v-card-actions class="pa-4 pt-0">
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="$emit('update:modelValue', false)">取消</v-btn>
        <v-btn color="primary" variant="elevated" @click="$emit('save')" class="px-6">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { Course } from '../../types/appModels'

interface SelectOption {
  title: string
  value: number
  name?: string
}

interface EnrollForm {
  studentId: number | null
  subjectId: number | null
  teacherId: number | null
  remark: string
}

defineProps<{
  modelValue: boolean
  isEdit: boolean
  course: Partial<Course>
  form: EnrollForm
  studentOptions: SelectOption[]
  subjectOptions: SelectOption[]
  teacherOptions: SelectOption[]
  studentLoading?: boolean
  subjectLoading?: boolean
  teacherLoading?: boolean
}>()

defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'save'): void
  (e: 'searchStudent', value: string): void
  (e: 'searchSubject', value: string): void
  (e: 'searchTeacher', value: string): void
}>()
</script>
