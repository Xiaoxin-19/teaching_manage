import { ref, computed, reactive, watch } from 'vue';
import { useToast } from '../../composables/useToast';
import { useConfirm } from '../../composables/useConfirm';
import type { ResponseWrapper, TeachingRecord } from '../../types/appModels';
import { SaveRecordData } from './types';
import { Dispatch } from '../../../wailsjs/go/main/App';
import { GetRecordListResponse, RecordDTO } from '../../types/response';

export function useRecordManage() {
  const { success, error, info, warning } = useToast();
  const confirm = useConfirm();

  // --- 状态定义 ---
  const searchStudent = ref('');
  const searchTeacher = ref('');
  const filterDateType = ref('全部时间');
  const customStartDate = ref('');
  const customEndDate = ref('');
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
    { title: '授课老师', key: 'teacher_name', width: '130px', sortable: false },
    { title: '上课日期', key: 'date', width: '140px', sortable: false },
    { title: '时间段', key: 'time', width: '130px', sortable: false },
    { title: '状态', key: 'status', width: '100px', sortable: false, align: 'center' },
    { title: '操作', key: 'actions', align: 'end', width: '120px', sortable: false },
    { title: '详情', key: 'data-table-expand', align: 'end', width: '60px' },
  ];

  const dateOptions = ['全部时间', '本周', '上周', '本月', '上月', '自定义'];

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
      !!searchStudent.value ||
      !!searchTeacher.value ||
      filterDateType.value !== '全部时间'
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

  // --- 核心方法 ---

  // 加载数据 
  const loadItems = async ({ page: p, itemsPerPage: ipp, sortBy }: any) => {
    loading.value = true;

    let reqData: any = {
      student_key: searchStudent.value || '',
      teacher_key: searchTeacher.value || '',
      start_date: effectiveDateRange.value?.start || '',
      end_date: effectiveDateRange.value?.end || '',
      offset: (p - 1) * ipp,
      limit: ipp,
    };

    console.log('加载记录列表，参数:', reqData);

    // json 转换
    const reqStr: string = JSON.stringify(reqData);
    Dispatch("record_manager:get_record_list", reqStr).then(
      (result: any) => {
        console.log('Received record list response:' + result);
        const resp = JSON.parse(result) as ResponseWrapper<GetRecordListResponse>;
        if (resp.code === 200) {
          const items: TeachingRecord[] = resp.data.records.map((item) => ({
            id: item.id,
            status: item.active ? 'active' : 'pending',
            date: item.teaching_date,
            time: `${item.start_time}-${item.end_time}`,
            startTime: item.start_time,
            endTime: item.end_time,
            studentId: item.student_id,
            studentName: item.student_name,
            teacherId: item.teacher_id,
            teacherName: item.teacher_name,
            remark: item.remark,
          }));
          serverItems.value = items;
          totalItems.value = resp.data.total; // 需要后端支持返回总数
          pendingCount.value = resp.data.total_pending;
        } else {
          console.error('获取记录列表失败:', resp.message);
          error('获取记录列表失败: ' + resp.message);
        }
      }
    );

    loading.value = false;
  };

  // 监听筛选条件变化
  watch([searchStudent, searchTeacher, effectiveDateRange], () => {
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
    searchStudent.value = '';
    searchTeacher.value = '';
    filterDateType.value = '全部时间';
    customStartDate.value = '';
    customEndDate.value = '';
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
    Dispatch('record_manager:create_record', JSON.stringify(data)).then((result: any) => {
      console.log('Record created:', result);
      const resp = JSON.parse(result) as ResponseWrapper<string>;
      if (resp.code === 200) {
        success('记录添加成功 (待生效)');
        loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
      } else {
        error(`添加记录失败: ${resp.message}`);
      }
    });
  };

  const activateRecord = (item: TeachingRecord) => {
    let reqData = {
      record_id: item.id,
    };

    console.log('Activating record:', reqData);
    Dispatch('record_manager:activate_record', JSON.stringify(reqData)).then((result: any) => {
      console.log('Record activated:', result);
      const resp = JSON.parse(result) as ResponseWrapper<string>;
      if (resp.code === 200) {
        success('记录已激活');
        loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
      } else if (resp.message.includes('fail')) {
        console.error('激活记录失败:', resp.message);
        error(`激活记录失败`);
      }
    });
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
      console.log('Activating all pending records');
      Dispatch('record_manager:activate_all_pending_records', "{}").then((result: any) => {
        console.log('All pending records activated:', result);
        const resp = JSON.parse(result) as ResponseWrapper<string>;
        if (resp.code === 200) {
          success('所有待处理记录已激活');
          loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
        } else if (resp.message.includes('fail')) {
          console.error('批量激活记录失败:', resp.message);
          error(`批量激活记录失败`);
        }
      })
    };
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
      let reqData = {
        record_id: item.id,
      };
      console.log('Deleting record:', reqData);
      Dispatch('record_manager:delete_record_by_id', JSON.stringify(reqData)).then((result: any) => {
        console.log('Record deleted:', result);
        const resp = JSON.parse(result) as ResponseWrapper<string>;
        if (resp.code === 200) {
          if (item.status === 'active') {
            info('记录已删除，学生课时已返还');
          } else {
            info('记录已删除');
          }
          loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
        } else if (resp.message.includes('fail')) {
          error(`删除记录失败`);
        }
      });
    }
  };

  const exportRecords = () => {
    const params = {
      student_key: searchStudent.value || '',
      teacher_key: searchTeacher.value || '',
      start_date: effectiveDateRange.value?.start,
      end_date: effectiveDateRange.value?.end,
    };
    console.log('导出参数:', params);
    Dispatch('record_manager:export_record_to_excel', JSON.stringify(params)).then((result: any) => {
      console.log('Export records result:', result);
      const resp = JSON.parse(result) as ResponseWrapper<string>;
      if (resp.code === 200) {
        success('记录导出成功');
      } else {
        if (resp.message.includes('cancel')) {
          info('已取消导出操作');
          return;
        }
        console.error('导出记录失败:', resp.message);
        error('导出记录失败: ' + resp.message, 'top-right');
      }
    });
  };

  return {
    // State
    searchStudent,
    searchTeacher,
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
    customStartDate,
    customEndDate,
    mockStudents: ref([]),
    pendingCount,
    hasActiveFilters,
    dateRangeText,

    // Methods
    loadItems,
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
  };
}