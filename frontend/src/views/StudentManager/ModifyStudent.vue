<template>
  <v-dialog :model-value="modelValue" @update:model-value="close" max-width="600px" scrollable
    transition="dialog-bottom-transition">
    <v-card class="rounded-lg elevation-4">
      <v-card-title class="d-flex justify-space-between align-center py-3 px-4">
        <div class="d-flex align-center text-subtitle-1 font-weight-bold">
          <v-icon :icon="isEdit ? 'mdi-account-edit-outline' : 'mdi-account-plus-outline'" color="primary"
            class="mr-2"></v-icon>
          {{ isEdit ? '编辑学生档案' : '新增学生' }}
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" density="comfortable" @click="close"></v-btn>
      </v-card-title>

      <v-divider></v-divider>

      <v-card-text class="pa-4">
        <v-form ref="formRef" v-model="isFormValid">
          <v-row dense>
            <v-col cols="12" sm="6">
              <v-text-field v-model="formData.name" label="学生姓名" variant="outlined" density="compact"
                prepend-inner-icon="mdi-account" hide-details="auto" class="mb-3"
                :rules="[v => !!v || '姓名必填']"></v-text-field>
            </v-col>

            <v-col cols="12" sm="6" class="d-flex align-center mb-3">
              <span class="text-caption text-medium-emphasis mr-3 ml-2">性别:</span>
              <v-radio-group v-model="formData.gender" inline hide-details density="compact" class="mt-0">
                <v-radio label="男" value="male" color="primary"></v-radio>
                <v-radio label="女" value="female" color="pink"></v-radio>
              </v-radio-group>
            </v-col>

            <v-col cols="12" sm="6">
              <v-text-field v-model="formData.phone" label="联系电话" placeholder="11位手机号" variant="outlined"
                density="compact" prepend-inner-icon="mdi-phone" hide-details="auto" class="mb-3"></v-text-field>
            </v-col>

            <!-- 状态编辑下拉框 -->
            <v-col v-show="isEdit" cols="12" sm="6">
              <v-select v-model="formData.status" :items="[
                { label: '正常', value: 1 },
                { label: '停课', value: 2 },
                { label: '退学', value: 3 }
              ]" label="学生状态" item-title="label" item-value="value" variant="outlined" density="compact"
                prepend-inner-icon="mdi-account-check" hide-details="auto" class="mb-3"></v-select>
              <!-- 状态提示语 -->
              <div class="text-caption px-1" style="min-height: 20px;">
                <span v-if="formData.status === 3" class="text-error d-flex align-center">
                  <v-icon icon="mdi-alert-circle-outline" size="x-small" class="mr-1"></v-icon>
                  退学状态无法充值和产生教学记录
                </span>
                <span v-else-if="formData.status === 2" class="text-warning d-flex align-center">
                  <v-icon icon="mdi-alert-outline" size="x-small" class="mr-1"></v-icon>
                  停课状态可以充值，但无法产生教学记录
                </span>
              </div>
            </v-col>


            <v-col cols="12">
              <v-textarea v-model="formData.remark" label="备注信息" variant="outlined" density="compact" rows="3"
                hide-details="auto" no-resize prepend-inner-icon="mdi-comment-text-outline"></v-textarea>
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>

      <v-divider></v-divider>

      <v-card-actions class="pa-3 d-flex align-center">
        <div v-if="isEdit && formData.lastModified" class="text-caption text-medium-emphasis ml-2 d-flex align-center">
          <v-icon icon="mdi-clock-outline" size="x-small" class="mr-1"></v-icon>
          上次编辑: {{ formData.lastModified }}
        </div>
        <v-spacer></v-spacer>
        <v-btn variant="text" class="mr-2" @click="close">取消</v-btn>
        <v-btn prepend-icon="mdi-check" color="primary" variant="elevated" elevation="1" :disabled="isSubmitDisabled"
          @click="save">
          保存
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { useModifyStudent } from './ModifyStudent.logic'
import type { StudentData, TeacherOption } from '../../types/appModels'

const props = defineProps<{
  modelValue: boolean
  isEdit?: boolean
  initialData?: StudentData
  teacherOptions?: TeacherOption[]
}>()

const emit = defineEmits(['update:modelValue', 'save'])

const { formData, close, save, isSubmitDisabled, formRef, isFormValid } = useModifyStudent(props as any, emit)
</script>