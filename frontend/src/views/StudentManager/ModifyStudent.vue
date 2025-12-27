<template>
  <v-dialog :model-value="modelValue" @update:model-value="close" max-width="600px" scrollable
    transition="dialog-bottom-transition">
    <v-card class="rounded-lg elevation-4">
      <!-- 1. 标题栏 -->
      <v-card-title class="d-flex justify-space-between align-center py-2 px-4 bg-white border-b">
        <div class="d-flex align-center text-subtitle-2 font-weight-bold text-grey-darken-3">
          <v-icon :icon="isEdit ? 'mdi-account-edit-outline' : 'mdi-account-plus-outline'" color="primary"
            class="mr-2"></v-icon>
          {{ isEdit ? '编辑学生档案' : '新增学生' }}
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" color="grey" @click="close"></v-btn>
      </v-card-title>

      <!-- 2. 表单内容区 -->
      <v-card-text class="pa-4">
        <v-form ref="formRef">
          <v-row dense>

            <!-- 第一行：姓名 + 性别 -->
            <v-col cols="12" sm="6">
              <v-text-field v-model="formData.name" label="学生姓名" variant="outlined" density="compact"
                prepend-inner-icon="mdi-account" hide-details="auto" class="mb-2"></v-text-field>
            </v-col>

            <v-col cols="12" sm="6" class="d-flex align-center mb-2">
              <span class="text-caption text-grey mr-3 ml-2">性别:</span>
              <v-radio-group v-model="formData.gender" inline hide-details density="compact" class="mt-0">
                <v-radio label="男" value="男" color="primary"></v-radio>
                <v-radio label="女" value="女" color="pink"></v-radio>
              </v-radio-group>
            </v-col>

            <!-- 第二行：电话 + 绑定教师 -->
            <v-col cols="12" sm="6">
              <v-text-field v-model="formData.phone" label="联系电话" placeholder="11位手机号" variant="outlined"
                density="compact" prepend-inner-icon="mdi-phone" hide-details="auto" class="mb-2"></v-text-field>
            </v-col>

            <v-col cols="12" sm="6">
              <v-autocomplete v-model="formData.teacherId" :items="teacherOptions" item-title="name" item-value="id"
                label="绑定教师" placeholder="请选择教师" variant="outlined" density="compact"
                prepend-inner-icon="mdi-account-tie" menu-icon="mdi-chevron-down" hide-details="auto" clearable
                class="mb-2">
                <template v-slot:item="{ props, item }">
                  <v-list-item v-bind="props" prepend-icon="mdi-account-outline"></v-list-item>
                </template>
              </v-autocomplete>
            </v-col>

            <!-- 第三行：备注 -->
            <v-col cols="12">
              <v-textarea v-model="formData.remark" label="备注信息" variant="outlined" density="compact" rows="3"
                hide-details="auto" no-resize></v-textarea>
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>

      <!-- 3. 底部操作栏 -->
      <v-card-actions class="pa-2 bg-grey-lighten-4 d-flex align-center">
        <!-- 上次编辑时间 (仅编辑模式显示) -->
        <div v-if="isEdit && formData.lastModified" class="text-caption text-medium-emphasis ml-2 d-flex align-center">
          <v-icon icon="mdi-clock-outline" size="x-small" class="mr-1"></v-icon>
          {{ formData.lastModified }}
        </div>

        <v-spacer></v-spacer>
        <v-btn color="grey-darken-1" size="small" variant="text" class="mr-2" @click="close">
          取消
        </v-btn>
        <v-btn prepend-icon="mdi-check" color="primary" size="small" variant="elevated" elevation="1" @click="save">
          保存
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { useModifyStudent } from './ModifyStudent.logic'
import type { StudentData, TeacherOption } from '../../types/appModels'

// Props
const props = defineProps<{
  modelValue: boolean
  isEdit?: boolean
  initialData?: StudentData
  // 教师列表改为从父组件传入
  teacherOptions?: TeacherOption[]
}>()

// Emits
const emit = defineEmits(['update:modelValue', 'save'])

const { formRef, formData, close, save } = useModifyStudent(props as any, emit as any)
</script>

<style scoped>
.bg-grey-lighten-4 {
  background-color: #f5f5f5 !important;
}

.border-b {
  border-bottom: 1px solid rgba(0, 0, 0, 0.12);
}
</style>