<template>
  <v-dialog :model-value="modelValue" @update:model-value="close" max-width="500px">
    <v-card class="rounded-lg elevation-4">
      <v-card-title class="d-flex justify-space-between align-center py-3 px-4">
        <div class="d-flex align-center text-subtitle-1 font-weight-bold">
          <v-icon icon="mdi-plus-circle" color="primary" class="mr-2"></v-icon>
          新增教学记录
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" @click="close"></v-btn>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text class="pa-4 pt-6">
        <div
          class="text-caption text-warning mb-4 d-flex align-center bg-yellow-lighten-5 pa-2 rounded border border-warning"
          style="border-style: dashed !important;">
          <v-icon start size="small" color="warning">mdi-information</v-icon>
          添加的记录默认为“未生效”状态，请稍后确认生效。
        </div>
        <v-form ref="formRef">
          <!-- 学生选择 (支持后端模糊搜索) -->
          <v-autocomplete v-model="formData.studentId" v-model:search="studentSearch" :items="studentOptions"
            :loading="loadingStudents" item-title="name" item-value="id" label="上课学生" placeholder="输入姓名搜索..."
            variant="outlined" density="comfortable" prepend-inner-icon="mdi-account-school" class="mb-3"
            :rules="[(v) => !!v || '请选择学生']" no-filter hide-no-data @update:search="onStudentSearch"></v-autocomplete>

          <!-- 日期与时间 -->
          <v-row dense>
            <v-col cols="6">
              <v-text-field v-model="formData.date" label="上课日期" type="date" variant="outlined" density="comfortable"
                class="mb-3" :rules="[(v) => !!v || '请选择日期']"></v-text-field>
            </v-col>
            <v-col cols="3">
              <v-text-field v-model="formData.startTime" label="开始" type="time" variant="outlined" density="comfortable"
                class="mb-3" :rules="startTimeRules"></v-text-field>
            </v-col>
            <v-col cols="3">
              <v-text-field v-model="formData.endTime" label="结束" type="time" variant="outlined" density="comfortable"
                class="mb-3" :rules="endTimeRules"></v-text-field>
            </v-col>
          </v-row>
          <v-textarea v-model="formData.remark" label="备注" variant="outlined" density="comfortable" rows="2" no-resize
            prepend-inner-icon="mdi-comment-text-outline"></v-textarea>
        </v-form>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions class="pa-3">
        <v-spacer></v-spacer>
        <v-btn variant="text" class="mr-2" @click="close">取消</v-btn>
        <v-btn color="primary" variant="flat" :disabled="isSubmitDisabled" @click="save">确认添加</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { useModifyRecord } from './ModifyRecord.logic';
import { SaveRecordData } from './types';

const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'save', data: SaveRecordData): void;
}>();

const {
  formRef,
  studentSearch,
  studentOptions,
  loadingStudents,
  formData,
  startTimeRules,
  endTimeRules,
  isSubmitDisabled,
  onStudentSearch,
  close,
  save,
} = useModifyRecord(props, emit);
</script>