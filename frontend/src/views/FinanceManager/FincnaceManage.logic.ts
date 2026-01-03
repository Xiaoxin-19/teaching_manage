import { ref, reactive, computed, onMounted, watch } from 'vue'
import { debounce } from 'lodash'
import { useToast } from '../../composables/useToast'
import { GetStudentList } from '../../api/student'
import { GetSubjectList } from '../../api/subject'
import type { GetOrderListRequest, GetStudentListRequest, GetSubjectListRequest } from '../../types/request'
import { categorizeOrderTags } from '../../utils/classification' // 使用项目现有的工具
import type { Order, OrderTag } from '../../types/appModels'
import { ExportOrdersToExcel, GetOrderList } from '../../api/order'
import { GetOrderListResponse } from '../../types/response'

interface SelectOption {
  title: string
  value: string | number
}


export function useFinanceManage() {
  const { success, info, error } = useToast()

  const loading = ref(false)
  const orders = ref<Order[]>([])
  const expanded = ref<string[]>([]) // 控制展开行


  // 分页控制
  const page = ref(1)
  const itemPerPage = ref(10)
  const totalItems = ref(0)

  // --- 筛选状态 ---
  const filters = reactive({
    studentId: null as number | null, // 动态搜索选中 ID
    studentNameLabel: '', // 用于 Chip 显示名字
    type: '' as string,
    subjectIds: [] as number[], // 动态搜索选中 IDs
    dateStart: "" as string,
    dateEnd: "" as string,
    keyword: '' // 仅用于前端显示的搜索词绑定
  })

  // --- 动态搜索状态 ---
  const studentOptions = ref<SelectOption[]>([])
  const isStudentLoading = ref(false)

  const subjectOptions = ref<SelectOption[]>([])
  const isSubjectLoading = ref(false)

  // 字典配置
  const typeOptions = [
    { title: '充值', value: 'increase' },
    { title: '退费', value: 'decrease' }
  ]

  // 表头定义 (标签列后置)
  const headers: any = [
    { title: '学员', key: 'studentName', align: 'start', sortable: false },
    { title: '科目', key: '_subjectName', align: 'center', sortable: false },
    { title: '类型', key: 'type', align: 'center', sortable: false },
    { title: '课时变动', key: 'hours', align: 'center', sortable: false },
    { title: '发生金额', key: 'amount', align: 'end', sortable: false },
    { title: '交易时间', key: 'created_at', align: 'end', sortable: false },
    { title: '标签', key: 'tags', align: 'start', sortable: false },
    { title: '', key: 'data-table-expand', align: 'end', sortable: false },
  ]

  // --- 核心方法: 服务器动态搜索 ---

  // 1. 搜索学员 (防抖)
  const searchStudents = async (keyword: string) => {
    isStudentLoading.value = true
    try {
      const req: GetStudentListRequest = {
        keyword,
        Offset: 0,
        Limit: 25,
        Status_Level: 3, // 仅搜索在读/潜在
        Status_Target: 0
      }
      const res = await GetStudentList(req)

      const newOptions: SelectOption[] = (res.students || []).map(s => ({
        title: `${s.name} (${s.student_number || '无学号'})`,
        value: s.id
      }))

      // 保留当前已选中的项
      if (filters.studentId) {
        const currentSelected = studentOptions.value.find(o => o.value === filters.studentId)
        if (currentSelected && !newOptions.find(o => o.value === filters.studentId)) {
          newOptions.unshift(currentSelected)
        }
      }
      studentOptions.value = newOptions
    } catch (e) {
      console.error('Search students failed', e)
    } finally {
      isStudentLoading.value = false
    }
  }

  // 2. 搜索科目 (防抖)
  const searchSubjects = async (keyword: string) => {
    isSubjectLoading.value = true
    try {
      const req: GetSubjectListRequest = {
        Keyword: keyword,
        Offset: 0,
        Limit: 25
      }
      const res = await GetSubjectList(req)

      const newOptions: SelectOption[] = (res.subjects || []).map(s => ({
        title: s.name,
        value: s.id
      }))

      // 保留已选项
      if (filters.subjectIds.length > 0) {
        const selected = subjectOptions.value.filter(o => filters.subjectIds.includes(o.value as number))
        selected.forEach(s => {
          if (!newOptions.find(o => o.value === s.value)) {
            newOptions.push(s)
          }
        })
      }
      subjectOptions.value = newOptions
    } catch (e) {
      console.error('Search subjects failed', e)
    } finally {
      isSubjectLoading.value = false
    }
  }

  // 创建防抖函数
  const onStudentSearch = debounce((v: string) => searchStudents(v), 500)
  const onSubjectSearch = debounce((v: string) => searchSubjects(v), 500)

  // --- 筛选逻辑 ---

  const activeFilters = computed(() => {
    const list = []
    if (filters.studentId) {
      const option = studentOptions.value.find(o => o.value === filters.studentId)
      const label = option ? option.title : (filters.studentNameLabel || '指定学员')
      list.push({ key: 'studentId', label: `学员: ${label}` })
    }

    if (filters.subjectIds.length > 0) {
      const names = filters.subjectIds.map(id => {
        const opt = subjectOptions.value.find(o => o.value === id)
        return opt ? opt.title : '未知科目'
      })
      list.push({ key: 'subjects', label: `科目: ${names.join(', ')}` })
    }

    if (filters.type) {
      const label = filters.type === 'increase' ? '充值' : '退费'
      list.push({ key: 'type', label: `类型: ${label}` })
    }

    if (filters.dateStart || filters.dateEnd) {
      const start = filters.dateStart || '不限'
      const end = filters.dateEnd || '不限'
      list.push({ key: 'date', label: `时间: ${start} 至 ${end}` })
    }
    return list
  })

  const clearFilter = (key: string) => {
    if (key === 'studentId') {
      filters.studentId = null
      filters.studentNameLabel = ''
    }
    else if (key === 'type') filters.type = ''
    else if (key === 'subjects') filters.subjectIds = []
    else if (key === 'date') {
      filters.dateStart = ""
      filters.dateEnd = ""
    }
    // 自动触发 watch 加载
  }

  function loadItems({ page: newPage, itemsPerPage: newItemsPerPage, sortBy }: { page: number; itemsPerPage: number; sortBy?: string[] | string | undefined }): void {
    page.value = newPage
    itemPerPage.value = newItemsPerPage
    loadData()
  }

  // --- 模拟数据加载 ---
  const loadData = async () => {
    loading.value = true
    try {
      let reqData: GetOrderListRequest = {
        student_id: filters.studentId || 0,
        subject_ids: filters.subjectIds.length > 0 ? filters.subjectIds : [],
        type: filters.type ? [filters.type] : [],
        date_start: filters.dateStart || "",
        date_end: filters.dateEnd || "",
        Offset: (page.value - 1) * itemPerPage.value,
        Limit: itemPerPage.value
      }
      console.log('Request data:', reqData)

      let data: GetOrderListResponse = await GetOrderList(reqData)
      console.log('Response data:', data)
      data.orders = data.orders || []
      orders.value = data.orders.map((o: Order) => {

        let data: Order = {
          ...o,
          tags: categorizeOrderTags(o.remark || '', o.hours) as OrderTag[],
          _subjectName: o.subject.name // 新增字段用于显示
        }
        return data
      })
      totalItems.value = data.total
    } catch (e) {
      console.error('Load orders failed', e)
    } finally {
      loading.value = false
    }
  }

  // 监听筛选变化
  const onFilterChange = debounce(() => {
    page.value = 1
    loadData()
  }, 500)

  watch(filters, () => {
    onFilterChange()
  }, { deep: true })

  // 监听学员选择变化，保存名字
  watch(() => filters.studentId, (newId) => {
    if (newId) {
      const option = studentOptions.value.find(o => o.value === newId)
      if (option) filters.studentNameLabel = option.title
    }
  })

  // --- 导出逻辑 (实现需求：访问筛选信息) ---
  const handleExport = async () => {
    let reqData: GetOrderListRequest = {
      student_id: filters.studentId || 0,
      subject_ids: filters.subjectIds.length > 0 ? filters.subjectIds : [],
      type: filters.type ? [filters.type] : [],
      date_start: filters.dateStart || "",
      date_end: filters.dateEnd || "",
      Offset: 0,
      Limit: -1 // 导出所有符合条件的数据
    }
    console.log('Export request data:', reqData)
    try {
      info('正在导出订单数据，请稍候...')
      let result = await ExportOrdersToExcel(reqData)
      console.log('Export result:', result)
      if (result == "cancel") {
        info('已取消导出操作')
        return
      }
      success('订单数据已成功导出至 Excel 文件')

    } catch (e) {
      if (e instanceof Error) {
        error('导出学生数据异常: ' + e.message, 'top-right')
      }
      console.error('导出提示异常:', e)
    }
  }

  // --- 辅助函数 ---

  const getTypeColor = (type: string) => 'grey-darken-1' // 统一灰色
  const getTypeText = (type: string) => type === 'increase' ? '充值' : '退费'

  const formatCurrency = (val: number) => {
    if (val === 0) return '¥0'
    const sign = val > 0 ? '+' : (val < 0 ? '-' : '')
    return `${sign}¥${Math.abs(val).toLocaleString()}`
  }

  const formatHours = (val: number) => {
    if (val > 99) return '99+'
    if (val < -99) return '99+'
    return Math.abs(val)
  }

  const formatDate = (timestamp: number) => {
    if (!timestamp) return '-'
    return new Date(timestamp).toLocaleString()
  }

  onMounted(() => {
    // 预加载
    searchStudents('')
    searchSubjects('')
    loadData()
  })

  return {

    page,
    itemPerPage,
    totalItems,
    loading,
    filters,
    orders,
    headers,
    expanded,

    studentOptions,
    isStudentLoading,
    onStudentSearch,
    subjectOptions,
    isSubjectLoading,
    onSubjectSearch,
    typeOptions,

    activeFilters,
    loadItems,
    clearFilter,
    getTypeColor,
    getTypeText,
    formatCurrency,
    formatHours,
    formatDate,
    handleExport
  }
}