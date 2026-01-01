<template>
  <v-dialog :model-value="modelValue" @update:model-value="close" max-width="500px" scrollable
    transition="dialog-bottom-transition">
    <v-card class="rounded-lg elevation-4">
      <!-- 1. 标题栏 -->
      <!-- 移除 bg-white, text-grey-darken-3，让其自动跟随主题 -->
      <v-card-title class="d-flex justify-space-between align-center py-2 px-4">
        <div class="d-flex align-center text-subtitle-2 font-weight-bold">
          <v-icon :icon="isEdit ? 'mdi-account-edit-outline' : 'mdi-account-plus-outline'" color="primary"
            class="mr-2"></v-icon>
          {{ isEdit ? '编辑教师档案' : '新增教师' }}
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" @click="close"></v-btn>
      </v-card-title>

      <v-divider></v-divider>

      <!-- 2. 表单内容区 -->
      <v-card-text class="pa-4">
        <v-form ref="formRef">
          <v-row dense>

            <!-- 第一行：姓名 + 性别 -->
            <v-col cols="12" sm="6">
              <v-text-field v-model="formData.name" label="教师姓名" variant="outlined" density="compact"
                prepend-inner-icon="mdi-account" hide-details="auto" class="mb-3"></v-text-field>
            </v-col>

            <!-- 性别选择 -->
            <v-col cols="12" sm="6" class="d-flex align-center mb-3">
              <span class="text-caption text-medium-emphasis mr-3 ml-2">性别:</span>
              <v-radio-group v-model="formData.gender" inline hide-details density="compact" class="mt-0">
                <v-radio label="男" value="male" color="primary"></v-radio>
                <v-radio label="女" value="female" color="pink"></v-radio>
              </v-radio-group>
            </v-col>

            <!-- 第二行：手机号 -->
            <v-col cols="12">
              <v-text-field v-model="formData.phone" label="手机号" placeholder="11位手机号" variant="outlined"
                density="compact" prepend-inner-icon="mdi-phone" hide-details="auto" class="mb-3"></v-text-field>
            </v-col>

            <!-- 第三行：备注 -->
            <v-col cols="12">
              <v-textarea v-model="formData.remark" label="备注信息" variant="outlined" density="compact" rows="3"
                hide-details="auto" no-resize></v-textarea>
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>

      <v-divider></v-divider>

      <!-- 3. 底部操作栏 -->
      <!-- 移除 bg-grey-lighten-4，使用 pa-2 保持间距 -->
      <v-card-actions class="pa-3 d-flex align-center">
        <!-- 上次编辑时间 -->
        <div v-if="isEdit && formData.lastModified" class="text-caption text-medium-emphasis ml-2 d-flex align-center">
          <v-icon icon="mdi-clock-outline" size="x-small" class="mr-1"></v-icon>
          上次编辑: {{ formData.lastModified }}
        </div>

        <v-spacer></v-spacer>
        <v-btn variant="text" class="mr-2" @click="close">
          取消
        </v-btn>
        <v-btn prepend-icon="mdi-check" color="primary" variant="elevated" elevation="1" @click="save">
          保存
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

import type { Teacher } from '../../types/appModels'


// Props
const props = defineProps<{
  modelValue: boolean
  isEdit?: boolean
  initialData?: Teacher
}>()

// Emits
const emit = defineEmits(['update:modelValue', 'save'])

// 本地表单数据
const formData = ref<Teacher>({
  name: '',
  phone: '',
  remark: '',
  gender: ''
} as Teacher)

// 监听弹窗开启，重置或填充数据
watch(() => props.modelValue, (val) => {
  if (val) {
    if (props.isEdit && props.initialData) {
      // 深度拷贝以断开引用
      formData.value = JSON.parse(JSON.stringify(props.initialData))
    } else {
      // 重置为默认空状态
      formData.value = {
        name: '',
        phone: '',
        remark: '',
        gender: 'male'
      } as Teacher
    }
  }
})

const close = () => {
  emit('update:modelValue', false)
}

const save = () => {
  // 可以在此处添加表单校验逻辑，例如检查姓名是否为空
  emit('save', formData.value)
  close()
}
</script>