import { ref, computed, reactive, watch } from 'vue';
import { useToast } from '../../composables/useToast';
import { useConfirm } from '../../composables/useConfirm';
import type { TeachingRecord } from '../../types/appModels';

export function useRecordManage() {
  const { success, error, info, warning } = useToast();
  const { confirmDelete, confirmInfo } = useConfirm();

  // --- 状态定义 ---
  const searchStudent = ref('');
  const searchTeacher = ref('');
  const filterDateType = ref('全部时间');
  const customStartDate = ref('');
  const customEndDate = ref('');

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

  // 选项数据 (Mock) - 实际应通过 API 获取
  const mockStudents = ref([
    { id: 1, name: '张三' },
    { id: 2, name: '李四' },
    { id: 3, name: '王五' },
  ]);

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

  // --- Mock 数据 (模拟后端数据库) ---
  const allMockData = ref<TeachingRecord[]>([
    { id: 1, status: 'active', date: '2024-01-15', time: '10:00-12:00', startTime: '10:00', endTime: '12:00', studentId: 1, studentName: '张三', teacherId: 1, teacherName: '王老师', remark: '本节课重点讲解了二次函数的图像与性质。' },
    { id: 2, status: 'active', date: '2024-01-14', time: '14:00-15:00', startTime: '14:00', endTime: '15:00', studentId: 2, studentName: '李四', teacherId: 2, teacherName: '李老师', remark: '英语口语练习，重点纠正了发音问题。' },
    { id: 3, status: 'pending', date: '2024-01-14', time: '16:00-18:00', startTime: '16:00', endTime: '18:00', studentId: 3, studentName: '王五', teacherId: 1, teacherName: '王老师', remark: '物理力学复习，稍微有点跟不上进度。' },
    { id: 4, status: 'active', date: '2024-01-12', time: '09:00-11:00', startTime: '09:00', endTime: '11:00', studentId: 1, studentName: '张三', teacherId: 1, teacherName: '王老师', remark: '化学实验基础，学生表现很积极。' },
    { id: 5, status: 'pending', date: '2024-01-12', time: '13:00-15:00', startTime: '13:00', endTime: '15:00', studentId: 2, studentName: '李四', teacherId: 1, teacherName: '王老师', remark: '模拟模糊搜索匹配测试。' },
  ]);

  // --- 辅助函数 ---
  const formatDate = (date: Date) => date.toISOString().split('T')[0];

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

  // 计算待生效记录数
  const pendingCount = computed(() => {
    // 实际项目中可能需要单独的 API 或者在 list 接口返回
    return allMockData.value.filter((item) => item.status === 'pending').length;
  });

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

  // 加载数据 (模拟后端 API)
  const loadItems = async ({ page: p, itemsPerPage: ipp, sortBy }: any) => {
    loading.value = true;

    // 模拟网络延迟
    await new Promise((resolve) => setTimeout(resolve, 300));

    let items = [...allMockData.value];

    // 1. 过滤
    if (searchStudent.value) {
      items = items.filter((item) =>
        item.studentName.toLowerCase().includes(searchStudent.value.toLowerCase())
      );
    }
    if (searchTeacher.value) {
      items = items.filter((item) =>
        item.teacherName.toLowerCase().includes(searchTeacher.value.toLowerCase())
      );
    }
    const range = effectiveDateRange.value;
    if (range) {
      if (range.start !== '不限') {
        items = items.filter((item) => item.date >= range.start);
      }
      if (range.end !== '不限') {
        items = items.filter((item) => item.date <= range.end);
      }
    }

    // 2. 总数
    totalItems.value = items.length;

    // 3. 分页
    if (ipp > 0) {
      const start = (p - 1) * ipp;
      const end = start + ipp;
      items = items.slice(start, end);
    }

    serverItems.value = items;
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

  const saveRecord = (data: any) => {
    // TODO: 调用后端 API 创建
    const student = mockStudents.value.find((s) => s.id === data.studentId);
    allMockData.value.unshift({
      ...data,
      id: Date.now(),
      status: 'pending',
      studentName: student?.name || '未知',
      teacherId: 1,
      teacherName: '默认教师',
      time: `${data.startTime}-${data.endTime}`,
    });
    success('已添加记录 (待生效)');
    loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
  };

  const activateRecord = (item: TeachingRecord) => {
    // TODO: 调用后端 API 激活
    const record = allMockData.value.find((r) => r.id === item.id);
    if (record) {
      record.status = 'active';
      success('记录已生效，课时已扣除');
      loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
    }
  };

  const processAllPending = async () => {
    const count = allMockData.value.filter((item) => item.status === 'pending').length;
    if (count === 0) return;

    const confirmed = await confirmInfo(
      `确定要将所有 ${count} 条待生效记录转为已生效吗？\n这将扣除对应学生的课时。`
    );

    if (confirmed) {
      // TODO: 调用后端批量激活 API
      allMockData.value.forEach((item) => {
        if (item.status === 'pending') item.status = 'active';
      });
      success(`已成功处理 ${count} 条记录`);
      loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
    }
  };

  const deleteItem = async (item: TeachingRecord) => {
    const confirmed = await confirmDelete(
      `记录 ${item.date} ${item.studentName}`
    );
    if (confirmed) {
      // TODO: 调用后端删除 API
      allMockData.value = allMockData.value.filter((r) => r.id !== item.id);
      success('记录已撤销');
      loadItems({ page: page.value, itemsPerPage: itemsPerPage.value });
    }
  };

  const exportRecords = () => {
    // TODO: 调用后端导出接口
    const params = {
      studentName: searchStudent.value || '',
      teacherName: searchTeacher.value || '',
      startDate: effectiveDateRange.value?.start,
      endDate: effectiveDateRange.value?.end,
    };
    console.log('导出参数:', params);
    success('正在导出记录...');
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
    mockStudents,
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