<template>
  <v-dialog :model-value="modelValue" @update:model-value="close" max-width="400px">
    <v-card class="rounded-lg">
      <v-card-title class="d-flex align-center py-3 px-4 text-subtitle-1 font-weight-bold">
        <v-icon color="primary" icon="mdi-calendar-edit" class="mr-2"></v-icon>
        设置自定义时间段
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text class="pa-5">
        <v-row dense>
          <v-col cols="12">
            <v-text-field v-model="tempStartDate" label="开始日期" type="date" variant="outlined" density="comfortable"
              class="mb-4" hide-details="auto"></v-text-field>
          </v-col>
          <v-col cols="12">
            <v-text-field v-model="tempEndDate" label="结束日期" type="date" variant="outlined" density="comfortable"
              :error-messages="dateErrorMsg" hide-details="auto"></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="pa-3 pt-0">
        <v-spacer></v-spacer>
        <v-btn color="grey-darken-1" variant="text" @click="close">取消</v-btn>
        <v-btn color="primary" variant="flat" @click="confirm" :disabled="!isDateRangeValid">
          确认筛选
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

const props = defineProps<{
  modelValue: boolean;
  startDate: string;
  endDate: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'confirm', start: string, end: string): void;
  (e: 'cancel'): void;
}>();

const tempStartDate = ref('');
const tempEndDate = ref('');

// 同步 props 到本地状态
watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      tempStartDate.value = props.startDate;
      tempEndDate.value = props.endDate;
    }
  }
);

// 校验日期有效性
const isDateRangeValid = computed(() => {
  if (!tempStartDate.value || !tempEndDate.value) return false;
  if (tempStartDate.value > tempEndDate.value) return false;
  return true;
});

// 错误提示信息
const dateErrorMsg = computed(() => {
  if (
    tempStartDate.value &&
    tempEndDate.value &&
    tempStartDate.value > tempEndDate.value
  ) {
    return '开始日期不能晚于结束日期';
  }
  return '';
});

const close = () => {
  emit('update:modelValue', false);
  emit('cancel');
};

const confirm = () => {
  if (isDateRangeValid.value) {
    emit('confirm', tempStartDate.value, tempEndDate.value);
    emit('update:modelValue', false);
  }
};
</script>