import { ref, computed, reactive, watch } from 'vue';
import { debounce } from 'lodash';
import { useToast } from '../../composables/useToast';
import { useConfirm } from '../../composables/useConfirm';
import type { ResponseWrapper, TeachingRecord } from '../../types/appModels';
import { SaveRecordData } from './types';
import { Dispatch } from '../../../wailsjs/go/main/App';
import { GetRecordListResponse, RecordDTO } from '../../types/response';
import { ActivateRecordRequest, CreateRecordRequest, DeleteRecordRequest, ExportRecordsRequest, GetRecordListRequest, GetStudentListRequest } from '../../types/request';
import { ActivateAllPendingRecords, ActivateRecord, CreateRecord, DeleteRecordByID, ExportRecordToExcel, GetRecordList } from '../../api/record';
import { GetStudentList } from '../../api/student';
import { GetTeacherList } from '../../api/teacher';
import { GetSubjectList } from '../../api/subject';

export function useRecordManage() {
  const { success, error, info, warning } = useToast();
  const confirm = useConfirm();

  // --- 状态定义 ---
  const selectedStudents = ref<number[]>([]);
  const selectedTeachers = ref<number[]>([]);
  const selectedSubjects = ref<number[]>([]);

  const studentOptions = ref<{ title: string, value: number }[]>([]);
  const teacherOptions = ref<{ title: string, value: number }[]>([]);
  const subjectOptions = ref<{ title: string, value: number }[]>([]);

  const loadingStudents = ref(false);
  const loadingTeachers = ref(false);
  const loadingSubjects = ref(false);

  const filterDateType = ref('全部时间');
  const customStartDate = ref('');
  const customEndDate = ref('');
  const activeFilter = ref<boolean | null>(null); // null: 全部, true: 激活, false: 未激活
  const pendingCount = ref(0);

  // 弹窗控制
  const dialogForm = ref(false);
  const dialogImport = ref(false);
  const dialogDateRange = ref(false);

  // 表格分页与数据
  const page = ref(1);
  const itemsPerPage = ref(10);
  const totalItems = ref(0);
  const loading = ref(false);
  const serverItems = ref<TeachingRecord[]>([]);

  // 表头定义 (key 使用下划线风格以兼容 DOM 模板)
  const headers: any = [
    { title: '学生', key: 'student_name', width: '130px', sortable: false },
    { title: '科目', key: 'subject_name', width: '120px', sortable: false },
    { title: '授课老师', key: 'teacher_name', width: '130px', sortable: false },
    { title: '上课日期', key: 'date', width: '140px', sortable: false },
    { title: '时间段', key: 'time', width: '130px', sortable: false },
    { title: '状态', key: 'status', width: '100px', sortable: false, align: 'center' },
    { title: '操作', key: 'actions', align: 'end', width: '120px', sortable: false },
    { title: '详情', key: 'data-table-expand', align: 'end', width: '60px' },
  ];

  const dateOptions = ['全部时间', '本周', '上周', '本月', '上月', '自定义'];
  const activeOptions = [
    { title: '全部状态', value: null },
    { title: '已激活', value: true },
    { title: '未激活', value: false },
  ];

  // --- 辅助函数 ---
  const formatDate = (date: Date) => {
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  const getWeekRange = (offset = 0) => {
    const now = new Date();
    const day = now.getDay() || 7; // 周日为7
    const m = new Date(now);
    m.setDate(now.getDate() - day + 1 + offset * 7); // 本周一
    const s = new Date(m);
    s.setDate(m.getDate() + 6); // 本周日
    return { start: formatDate(m), end: formatDate(s) };
  };

  const getMonthRange = (offset = 0) => {
    const now = new Date();
    const year = now.getFullYear();
    const month = now.getMonth() + offset;
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    return { start: formatDate(firstDay), end: formatDate(lastDay) };
  };

  // --- Computed 属性 ---

  // 生效的日期范围对象
  const effectiveDateRange = computed(() => {
    const type = filterDateType.value;
    if (type === '全部时间') return null;
    if (type === '自定义') {
      return customStartDate.value || customEndDate.value
        ? { start: customStartDate.value || '不限', end: customEndDate.value || '不限' }
        : null;
    }
    if (type === '本周') return getWeekRange(0);
    if (type === '上周') return getWeekRange(-1);
    if (type === '本月') return getMonthRange(0);
    if (type === '上月') return getMonthRange(-1);
    return null;
  });

  // 是否有激活的筛选条件
  const hasActiveFilters = computed(() => {
    return (
      selectedStudents.value.length > 0 ||
      selectedTeachers.value.length > 0 ||
      selectedSubjects.value.length > 0 ||
      filterDateType.value !== '全部时间' ||
      activeFilter.value !== null
    );
  });

  // 日期范围文本显示
  const dateRangeText = computed(() => {
    if (filterDateType.value === '全部时间') return '';
    let text = filterDateType.value;
    if (effectiveDateRange.value) {
      text += ` (${effectiveDateRange.value.start} 至 ${effectiveDateRange.value.end})`;
    }
    return text;
  });

  const getSelectedNames = (ids: number[], options: { title: string; value: number }[], prefix: string) => {
    if (ids.length === 0) return '';
    const names = ids.map(id => {
      const option = options.find(o => o.value === id);
      return option ? option.title.split(' (')[0] : String(id);
    });

    const text = names.join(', ');
    if (text.length > 20 && names.length > 1) {
      return `${prefix}: ${names[0]} 等 ${ids.length} 人`;
    }
    return `${prefix}: ${text}`;
  };

  const selectedStudentText = computed(() => getSelectedNames(selectedStudents.value, studentOptions.value, '学生'));
  const selectedTeacherText = computed(() => getSelectedNames(selectedTeachers.value, teacherOptions.value, '老师'));
  const selectedSubjectText = computed(() => getSelectedNames(selectedSubjects.value, subjectOptions.value, '科目').replace('人', '个')); // 科目单位是个

  // --- 搜索方法 (防抖) ---
  const onStudentSearch = debounce(async (keyword: string) => {
    if (!keyword) return;
    loadingStudents.value = true;
    try {
      const res = await GetStudentList({ Offset: 0, Limit: 25, Keyword: keyword, Status_Level: 3, Status_Target: 0 } as any);
      const newOptions = res.students.map(s => ({ title: `${s.name} (${s.student_number})`, value: s.id }));

      // 保留已选中的项
      if (selectedStudents.value.length > 0) {
        const selected = studentOptions.value.filter(o => selectedStudents.value.includes(o.value));
        selected.forEach(s => {
          if (!newOptions.find(n => n.value === s.value)) {
            newOptions.push(s);
          }
        });
      }
      studentOptions.value = newOptions;
    } catch (e) {
      console.error(e);
    } finally {
      loadingStudents.value = false;
    }
  }, 300);

  const onTeacherSearch = debounce(async (keyword: string) => {
    if (!keyword) return;
    loadingTeachers.value = true;
    try {
      const res = await GetTeacherList({ Offset: 0, Limit: 25, Keyword: keyword } as any);
      const newOptions = res.teachers.map(t => ({ title: `${t.name} (${t.teacher_number})`, value: t.id }));

      if (selectedTeachers.value.length > 0) {
        const selected = teacherOptions.value.filter(o => selectedTeachers.value.includes(o.value));
        selected.forEach(s => {
          if (!newOptions.find(n => n.value === s.value)) {
            newOptions.push(s);
          }
        });
      }
      teacherOptions.value = newOptions;
    } catch (e) {
      console.error(e);
    } finally {
      loadingTeachers.value = false;
    }
  }, 300);

  const onSubjectSearch = debounce(async (keyword: string) => {
    if (!keyword) return;
    loadingSubjects.value = true;
    try {
      const res = await GetSubjectList({ Offset: 0, Limit: 25, Keyword: keyword } as any);
      const newOptions = res.subjects.map(s => ({ title: s.name, value: s.id }));

      if (selectedSubjects.value.length > 0) {
        const selected = subjectOptions.value.filter(o => selectedSubjects.value.includes(o.value));
        selected.forEach(s => {
          if (!newOptions.find(n => n.value === s.value)) {
            newOptions.push(s);
          }
        });
      }
      subjectOptions.value = newOptions;
    } catch (e) {
      console.error(e);
    } finally {
      loadingSubjects.value = false;
    }
  }, 300);

  // --- 核心方法 ---

  // 加载数据 
  const loadItems = async ({ page: p, itemsPerPage: ipp, sortBy }: any) => {
    loading.value = true;

    let reqData: GetRecordListRequest = {
      student_ids: selectedStudents.value.length > 0 ? selectedStudents.value : undefined,
      teacher_ids: selectedTeachers.value.length > 0 ? selectedTeachers.value : undefined,
      subject_ids: selectedSubjects.value.length > 0 ? selectedSubjects.value : undefined,
      start_date: effectiveDateRange.value?.start || '',
      end_date: effectiveDateRange.value?.end || '',
      offset: (p - 1) * ipp,
      limit: ipp,
      active: activeFilter.value,
    };

    console.log('加载记录列表，参数:', reqData);

    try {
      const data = await GetRecordList(reqData);
      const items: TeachingRecord[] = data.records.map((item: RecordDTO) => ({
        id: item.id,
        status: item.active ? 'active' : 'pending',
        date: item.teaching_date,
        time: `${item.start_time}-${item.end_time}`,
        startTime: item.start_time,
        endTime: item.end_time,
        studentId: item.student.id,
        studentName: item.student.name,
        teacherId: item.teacher.id,
        teacherName: item.teacher.name,
        subjectName: item.subject.name,
        remark: item.remark,
      }));
      serverItems.value = items;
      totalItems.value = data.total;
      pendingCount.value = data.total_pending;
    } catch (e: any) {
      console.error('获取记录列表失败:', e);
      error('获取记录列表失败: ' + e.message);
    } finally {
      loading.value = false;
    }
  };

  // 监听筛选条件变化
  watch([selectedStudents, selectedTeachers, selectedSubjects, effectiveDateRange, activeFilter], () => {
    page.value = 1;
    loadItems({ page: 1, itemsPerPage: itemsPerPage.value, sortBy: [] });
  });

  // --- 筛选操作 ---

  const selectDateFilter = (type: string) => {
    if (type === '自定义') {
      dialogDateRange.value = true;
    } else {
      customStartDate.value = '';
      customEndDate.value = '';
      filterDateType.value = type;
    }
  };

  const handleCustomDateConfirm = (start: string, end: string) => {
    customStartDate.value = start;
    customEndDate.value = end;
    filterDateType.value = '自定义';
  };

  const handleCustomDateCancel = () => {
    // 如果是自定义模式且没值，回滚到全部时间
    if (filterDateType.value === '自定义' && !customStartDate.value) {
      filterDateType.value = '全部时间';
    }
  };

  const clearAllFilters = () => {
    selectedStudents.value = [];
    selectedTeachers.value = [];
    selectedSubjects.value = [];
    filterDateType.value = '全部时间';
    customStartDate.value = '';
    customEndDate.value = '';
    activeFilter.value = null;
    info('已清空所有筛选条件');
  };

  const clearDateFilter = () => {
    filterDateType.value = '全部时间';
    customStartDate.value = '';
    customEndDate.value = '';
  };

  // --- CRUD 操作 ---

  const openAdd = () => {
    dialogForm.value = true;
  };

  const saveRecord = async (data: SaveRecordData) => {
    console.log('保存记录数据:', data);
    let reqData: CreateRecordRequest = {
      student_id: data.student_id,
      subject_id: data.subject_id,
      teaching_date: data.teaching_date,
      start_time: data.start_time,
      end_time: data.end_time,
      remark: data.remark || '',
    }

    try {
      console.log('Creating record with data:', reqData);
      let result = await CreateRecord(reqData);
      console.log('Record created successfully, result:', result);
      success('记录添加成功 (待生效)');
      loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
    } catch (e) {
      console.error('Failed to create record:', e);
      if (e instanceof Error) {
        error(`添加记录失败: ${e.message}`);
      }
    }
  };

  const activateRecord = async (item: TeachingRecord) => {
    let reqData: ActivateRecordRequest = {
      id: item.id,
    };

    console.log('Activating record:', reqData);
    try {
      await ActivateRecord(reqData);
      success(`记录已激活`);
      loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
    } catch (e) {
      console.error('Failed to activate record:', e);
      if (e instanceof Error) {
        error(`激活记录失败: ${e.message}`);
      }
    }
  };

  const processAllPending = async () => {
    const confirmed = await confirm.confirm(
      "批量激活确认",
      "确定要激活所有待处理记录吗？",
      {
        type: "warning",
        confirmText: '激活',
        cancelText: '取消'
      }
    );
    if (confirmed) {
      try {
        await ActivateAllPendingRecords();
        success('所有待处理记录已激活');
        loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
      } catch (e) {
        console.error('Failed to activate all pending records:', e);
        if (e instanceof Error) {
          error(`批量激活记录失败: ${e.message}`);
        }
      }
    }
  }


  const deleteItem = async (item: TeachingRecord) => {
    const confirmed = await confirm.confirm(
      "确认删除",
      `学生 ${item.studentName} 于 ${item.date} 的记录, 如果已生效，删除则会返还学生的课时。确定删除？`,
      {
        type: "warning",
        confirmText: '删除',
        cancelText: '取消'
      }
    );
    if (confirmed) {
      let reqData: DeleteRecordRequest = {
        id: item.id,
      };
      console.log('Deleting record:', reqData);

      try {
        await DeleteRecordByID(reqData);
        success('记录已删除, 学生课时已返还');
        loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
      } catch (e) {
        console.error('Failed to delete record:', e);
        if (e instanceof Error) {
          error(`删除记录失败: ${e.message}`);
        }
      }
    }
  };

  const exportRecords = async () => {
    const params: ExportRecordsRequest = {
      student_ids: selectedStudents.value.length > 0 ? selectedStudents.value : undefined,
      teacher_ids: selectedTeachers.value.length > 0 ? selectedTeachers.value : undefined,
      subject_ids: selectedSubjects.value.length > 0 ? selectedSubjects.value : undefined,
      start_date: effectiveDateRange.value?.start,
      end_date: effectiveDateRange.value?.end,
      active: activeFilter.value,
    };
    console.log('导出参数:', params);

    try {
      let result = await ExportRecordToExcel(params);
      if (result.includes('cancel')) {
        info('已取消导出操作');
        return;
      }
      success('记录导出成功');
    } catch (e) {
      if (e instanceof Error) {
        error('导出记录失败: ' + e.message, 'top-right');
      }

      console.error('Failed to export records:', e);
    }
  };

  // 错误弹窗控制
  const dialogError = ref(false);
  const importErrorInfos = ref<string[][]>([]);

  const onImportSuccess = () => {
    loadItems({ page: page.value, itemsPerPage: itemsPerPage.value, sortBy: [] });
  };

  const onImportFailed = (errorInfo: any) => {
    console.log("Import failed", errorInfo);
    if (errorInfo && Array.isArray(errorInfo) && errorInfo.length > 0) {
      importErrorInfos.value = errorInfo;
      dialogError.value = true;
    }
  };

  return {
    // State
    selectedStudents,
    selectedTeachers,
    selectedSubjects,
    studentOptions,
    teacherOptions,
    subjectOptions,
    loadingStudents,
    loadingTeachers,
    loadingSubjects,
    filterDateType,
    dateOptions,
    page,
    itemsPerPage,
    totalItems,
    loading,
    serverItems,
    headers,
    dialogForm,
    dialogImport,
    dialogDateRange,
    dialogError,
    importErrorInfos,
    customStartDate,
    customEndDate,
    mockStudents: ref([]),
    pendingCount,
    hasActiveFilters,
    dateRangeText,
    selectedStudentText,
    selectedTeacherText,
    selectedSubjectText,
    activeFilter,
    activeOptions,

    // Methods
    loadItems,
    onStudentSearch,
    onTeacherSearch,
    onSubjectSearch,
    selectDateFilter,
    handleCustomDateConfirm,
    handleCustomDateCancel,
    clearAllFilters,
    clearDateFilter,
    openAdd,
    saveRecord,
    activateRecord,
    processAllPending,
    deleteItem,
    exportRecords,
    onImportSuccess,
    onImportFailed,
  };
}