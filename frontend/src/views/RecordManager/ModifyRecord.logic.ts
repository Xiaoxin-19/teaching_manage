import { ref, reactive, watch, computed, nextTick } from 'vue';
import { debounce } from 'lodash';
import { Dispatch } from '../../../wailsjs/go/main/App';
import { ResponseWrapper } from '../../types/appModels';
import { GetStudentListResponse } from '../../types/response';
import { useToast } from '../../composables/useToast';
import { SaveRecordData, StudentOption } from './types';

export function useModifyRecord(props: { modelValue: boolean }, emit: any) {
  const toast = useToast();
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

  // 时间解析与校验
  const parseTimeToMinutes = (t: string) => {
    if (!t) return null;
    const m = t.match(/^(\d{2}):(\d{2})$/);
    if (!m) return null;
    return parseInt(m[1], 10) * 60 + parseInt(m[2], 10);
  };

  const MINUTES_PER_HOUR = 60;
  const MAX_DURATION_HOURS = 8;
  const MAX_DURATION_MINUTES = MAX_DURATION_HOURS * MINUTES_PER_HOUR;
  const MIN_DURATION_MINUTES = 5;

  const DURATION_ERROR_TOO_LONG = `上课时长不能超过${MAX_DURATION_HOURS}小时`;
  const DURATION_ERROR_TOO_SHORT = '上课时长太短';

  const timeValidationMessage = (start: string, end: string) => {
    if (!start || !end) return '';
    const s = parseTimeToMinutes(start);
    const e = parseTimeToMinutes(end);
    if (s === null || e === null) return '时间格式不正确';
    if (e <= s) return '结束时间必须晚于开始时间';
    const duration = e - s;
    if (duration > MAX_DURATION_MINUTES) return DURATION_ERROR_TOO_LONG;
    if (duration < MIN_DURATION_MINUTES) return DURATION_ERROR_TOO_SHORT;
    return '';
  };

  const startTimeRules = [
    (v: string) => !!v || '请选择开始时间',
    (v: string) => {
      const msg = timeValidationMessage(v, formData.endTime);
      return msg || true;
    },
  ];

  const endTimeRules = [
    (v: string) => !!v || '请选择结束时间',
    (v: string) => {
      const msg = timeValidationMessage(formData.startTime, v);
      return msg || true;
    },
  ];

  const isFormValid = computed(() => {
    if (!formData.studentId) return false;
    if (!formData.date) return false;
    if (!formData.startTime || !formData.endTime) return false;
    const s = parseTimeToMinutes(formData.startTime);
    const e = parseTimeToMinutes(formData.endTime);
    if (s === null || e === null) return false;
    if (e <= s) return false;
    if (e - s > 8 * 60) return false;
    if (e - s < 5) return false;
    return true;
  });

  const isSubmitDisabled = computed(() => !isFormValid.value);

  // 当开始/结束时间任一变化时，防抖触发表单校验，确保两个字段的错误能同时更新
  const debouncedValidate = debounce(() => {
    if (formRef.value && typeof formRef.value.validate === 'function') {
      formRef.value.validate();
    }
  }, 150);

  // 监听时间变化，触发整体校验
  watch(() => formData.startTime, () => {
    debouncedValidate();
  });
  watch(() => formData.endTime, () => {
    debouncedValidate();
  });

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
        // 在下一个 tick 重置字段校验状态，避免旧错误残留
        nextTick(() => {
          if (formRef.value && typeof formRef.value.resetValidation === 'function') {
            formRef.value.resetValidation();
          }
        });
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
      let reqData = {
        Key: keyword,
        Offset: 0,
        Limit: -1,
      };

      console.log('Fetching students with keyword:', keyword, 'Request data:', reqData);
      Dispatch('student_manager:get_student_list', JSON.stringify(reqData))
        .then((result: any) => {
          const resp = JSON.parse(result) as ResponseWrapper<GetStudentListResponse>
          if (resp.code === 200) {
            studentOptions.value = (resp.data.students || []).map((item) => ({
              id: item.id,
              name: item.name,
            }));
          } else {
            console.error('获取学生列表失败:', resp.message)
            toast.error('获取学生列表失败: ' + resp.message, 'top-right')
          }
        })
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
    if (!isFormValid.value) {
      toast.error('请检查必填项和时间范围是否正确', 'top-right');
      if (formRef.value) {
        await formRef.value.validate();
      }
      return;
    }

    if (formRef.value) {
      const { valid } = await formRef.value.validate();
      if (valid) {
        // 传递完整的学生对象信息可能更有用
        const selectedStudent = studentOptions.value.find(s => s.id === formData.studentId);

        // 格式化日期和时间
        const formattedDate = new Date(formData.date).toISOString().substring(0, 10);

        const saveData: SaveRecordData = {
          student_id: formData.studentId!,
          student_name: selectedStudent ? selectedStudent.name : '',
          teaching_date: formattedDate,
          start_time: formData.startTime,
          end_time: formData.endTime,
          remark: formData.remark,
        };

        emit('save', saveData);
        close();
      }
    }
  };

  return {
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
  };
}
