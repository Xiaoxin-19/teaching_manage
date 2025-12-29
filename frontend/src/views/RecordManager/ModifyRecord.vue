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
                class="mb-3"></v-text-field>
            </v-col>
            <v-col cols="3">
              <v-text-field v-model="formData.endTime" label="结束" type="time" variant="outlined" density="comfortable"
                class="mb-3"></v-text-field>
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
        <v-btn color="primary" variant="flat" @click="save">确认添加</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import { debounce } from 'lodash'; // 需要安装 lodash: npm install lodash @types/lodash
// 假设有从后端获取学生列表的方法，这里模拟引入
// import { GetStudentList } from '../../../wailsjs/go/main/App'; 

// 定义 StudentOption 类型，实际项目中应从 appModels 导入
interface StudentOption {
  id: number;
  name: string;
}

const props = defineProps<{
  modelValue: boolean;
  // studentOptions: StudentOption[]; // 不再通过 props 传入所有学生
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'save', data: any): void;
}>();

const formRef = ref<any>(null);
const studentSearch = ref('');
const studentOptions = ref<StudentOption[]>([]);
const loadingStudents = ref(false);

const defaultFormData = {
  studentId: null,
  date: new Date().toISOString().substring(0, 10),
  startTime: '09:00',
  endTime: '11:00',
  remark: '',
};

const formData = reactive({ ...defaultFormData });

// 每次打开弹窗重置表单
watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      Object.assign(formData, defaultFormData, {
        date: new Date().toISOString().substring(0, 10),
      });
      studentSearch.value = '';
      studentOptions.value = []; // 打开时可以加载默认列表或者清空
      // fetchStudents(''); // 可选：打开时加载初始列表
    }
  }
);

// 防抖搜索函数
const fetchStudents = debounce(async (keyword: string) => {
  if (!keyword) {
    studentOptions.value = [];
    return;
  }

  loadingStudents.value = true;
  try {
    // 模拟后端 API 调用
    // const res = await GetStudentList({ name: keyword, page: 1, pageSize: 10 });
    // studentOptions.value = res.data.list.map(s => ({ id: s.id, name: s.name }));

    // 模拟数据返回
    console.log(`Searching for student: ${keyword}`);
    await new Promise(resolve => setTimeout(resolve, 500)); // 模拟延迟

    // 简单的模拟匹配逻辑 (实际应由后端完成)
    const mockAllStudents = [
      { id: 1, name: '张三' },
      { id: 2, name: '李四' },
      { id: 3, name: '王五' },
      { id: 4, name: '张三丰' },
      { id: 5, name: '赵六' }
    ];
    studentOptions.value = mockAllStudents.filter(s => s.name.includes(keyword));

  } catch (error) {
    console.error("Failed to fetch students", error);
  } finally {
    loadingStudents.value = false;
  }
}, 300); // 300ms 防抖

const onStudentSearch = (val: string) => {
  fetchStudents(val);
};

const close = () => {
  emit('update:modelValue', false);
};

const save = async () => {
  if (formRef.value) {
    const { valid } = await formRef.value.validate();
    if (valid) {
      // 传递完整的学生对象信息可能更有用
      const selectedStudent = studentOptions.value.find(s => s.id === formData.studentId);
      emit('save', { ...formData, studentName: selectedStudent?.name });
      close();
    }
  }
};
</script>