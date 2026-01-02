import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { debounce } from 'lodash'
import { useToast } from '../../composables/useToast'
import {
  GetCourseList, EnrollCourse, RechargeCourse, ToggleCourseStatus, DeleteCourse, UpdateCourse,
  type Course
} from '../../api/course'
import { GetStudentList } from '../../api/student'
import { GetSubjectList } from '../../api/subject'
import { GetTeacherList } from '../../api/teacher'
import { GetStudentListResponse, GetSubjectListResponse, GetTeacherListResponse } from '../../types/response'
import { GetStudentListRequest } from '../../types/request'

export function useCourseManage() {
  const { success, error } = useToast()

  // --- 状态定义 ---
  const loading = ref(false)
  const page = ref(1)
  const itemsPerPage = ref(10)
  const totalItems = ref(0)
  const courses = ref<Course[]>([])

  // 弹窗状态
  const dialogVisible = ref(false) // 编辑/报名
  const rechargeDialogVisible = ref(false) // 充值/退费 (新组件控制)
  const deleteDialogVisible = ref(false) // 办理退课
  const forceDeleteDialogVisible = ref(false) // 彻底删除

  const isEdit = ref(false)
  const currentItem = ref<Partial<Course>>({})

  // 充值模式状态 (传递给子组件)
  const rechargeMode = ref<'charge' | 'refund'>('charge')

  // 引用子组件实例 (用于调用 success/error 方法)
  const rechargeDialogRef = ref<any>(null)

  // 报名/编辑表单
  const enrollForm = reactive({
    studentId: null as number | null,
    subjectId: null as number | null,
    teacherId: null as number | null,
    remark: ''
  })

  // 退课表单
  const deleteForm = reactive({
    clearBalance: true,
    remark: ''
  })

  // 筛选状态
  const filters = reactive({
    studentName: '',
    subjects: [] as string[],
    teachers: [] as string[],
    balanceMin: null as number | null,
    balanceMax: null as number | null,
    status: [] as string[]
  })

  // 字典数据
  const studentOptions = ref<{ title: string, value: number }[]>([])
  const subjectOptions = ref<{ title: string, value: number }[]>([])
  const teacherOptions = ref<{ title: string, value: number }[]>([])

  const isStudentLoading = ref(false)
  const isSubjectLoading = ref(false)
  const isTeacherLoading = ref(false)

  const statusOptions = [
    { title: '正常上课', value: 'normal' },
    { title: '课程暂停', value: 'paused' },
    { title: '已结课', value: 'finished' },
    { title: '学员停课', value: 'student_suspended' },
    { title: '学员退学', value: 'student_withdrawn' }
  ]

  // --- 搜索方法 ---
  const searchStudents = async (keyword: string) => {
    if (!keyword) return
    isStudentLoading.value = true
    try {
      let reqData: GetStudentListRequest = { Offset: 0, Limit: 25, keyword, Status_Level: 1, Status_Target: 1 }
      const res: GetStudentListResponse = await GetStudentList(reqData)
      let students = res.students
      studentOptions.value = students.map(s => ({ title: `${s.name} (${s.student_number})`, value: s.id }))
    } catch (e) {
      console.error(e)
    } finally {
      isStudentLoading.value = false
    }
  }

  const searchSubjects = async (keyword: string) => {
    isSubjectLoading.value = true
    try {
      const res: GetSubjectListResponse = await GetSubjectList({ Offset: 0, Limit: 25, Keyword: keyword })
      let subjectsList = res.subjects || []
      subjectOptions.value = subjectsList.map(s => ({ title: s.name, value: s.id }))
    } catch (e) {
      console.error(e)
    } finally {
      isSubjectLoading.value = false
    }
  }

  const searchTeachers = async (keyword: string) => {
    isTeacherLoading.value = true
    try {
      const res: GetTeacherListResponse = await GetTeacherList({ Offset: 0, Limit: 25, Keyword: keyword })
      let teachersList = res.teachers || []
      teacherOptions.value = teachersList.map(t => ({ title: `${t.name} (${t.teacher_number})`, value: t.id }))
    } catch (e) {
      console.error(e)
    } finally {
      isTeacherLoading.value = false
    }
  }

  const onStudentSearch = debounce(searchStudents, 500)
  const onSubjectSearch = debounce(searchSubjects, 500)
  const onTeacherSearch = debounce(searchTeachers, 500)

  // --- 辅助计算属性 ---

  const activeFilters = computed(() => {
    const list = []
    if (filters.studentName) list.push({ key: 'studentName', label: `学员: ${filters.studentName}` })
    if (filters.subjects.length > 0) list.push({ key: 'subjects', label: `科目: ${filters.subjects.join(', ')}` })
    if (filters.teachers.length > 0) list.push({ key: 'teachers', label: `老师: ${filters.teachers.join(', ')}` })
    if (filters.balanceMin !== null) list.push({ key: 'balanceMin', label: `课时 >= ${filters.balanceMin}` })
    if (filters.balanceMax !== null) list.push({ key: 'balanceMax', label: `课时 <= ${filters.balanceMax}` })
    if (filters.status.length > 0) {
      const titles = filters.status.map(v => statusOptions.find(o => o.value === v)?.title).filter(Boolean)
      list.push({ key: 'status', label: `状态: ${titles.join(', ')}` })
    }
    return list
  })

  const isDeleteValid = computed(() => {
    if ((currentItem.value.balance || 0) > 0) {
      return deleteForm.remark && deleteForm.remark.trim().length > 0;
    }
    return true;
  })

  // --- 核心方法 ---

  const loadData = async () => {
    loading.value = true
    try {
      const req = {
        search: filters.studentName,
        subjects: filters.subjects.length ? filters.subjects : undefined,
        teachers: filters.teachers.length ? filters.teachers : undefined,
        balanceMin: filters.balanceMin,
        balanceMax: filters.balanceMax,
        status: filters.status.length ? filters.status : undefined,
        page: page.value,
        pageSize: itemsPerPage.value
      }

      const data = await GetCourseList(req)
      courses.value = data.courses
      totalItems.value = data.total
    } catch (e: any) {
      console.error(e)
      courses.value = []
      totalItems.value = 0
    } finally {
      loading.value = false
    }
  }

  const clearFilter = (key: string) => {
    if (key === 'studentName') filters.studentName = ''
    else if (key === 'subjects') filters.subjects = []
    else if (key === 'teachers') filters.teachers = []
    else if (key === 'balanceMin') filters.balanceMin = null
    else if (key === 'balanceMax') filters.balanceMax = null
    else if (key === 'status') filters.status = []
    loadData()
  }

  // --- 业务操作 ---

  const openEnroll = () => {
    isEdit.value = false
    currentItem.value = {}
    enrollForm.studentId = null
    enrollForm.subjectId = null
    enrollForm.teacherId = null
    enrollForm.remark = ''

    // 重置选项
    studentOptions.value = []
    subjectOptions.value = []
    teacherOptions.value = []

    // 预加载部分数据 (可选)
    searchSubjects('')
    searchTeachers('')

    dialogVisible.value = true
  }

  const openEdit = (item: Course) => {
    isEdit.value = true
    currentItem.value = { ...item }
    enrollForm.teacherId = item.teacherId as any
    enrollForm.remark = ''

    // 初始化选项，确保当前值能正确显示
    teacherOptions.value = [{ title: item.teacherName, value: item.teacherId }]
    // 预加载更多老师以便更换
    searchTeachers('')

    dialogVisible.value = true
  }

  const handleSave = async () => {
    loading.value = true
    try {
      if (isEdit.value) {
        await UpdateCourse({
          id: currentItem.value.id!,
          teacherId: enrollForm.teacherId!,
          remark: enrollForm.remark
        })
        success('课程信息更新成功')
      } else {
        await EnrollCourse(enrollForm as any)
        success('新课报名成功')
      }
      dialogVisible.value = false
      loadData()
    } catch (e: any) {
      error(e.message)
    } finally {
      loading.value = false
    }
  }

  // --- 充值/退费 (重构后) ---

  const openRecharge = (item: Course) => {
    const status = getEffectiveStatus(item)
    if (status.label === '学员退学') {
      error('操作失败：学员已退学')
      return
    }
    currentItem.value = { ...item }
    rechargeMode.value = 'charge'
    rechargeDialogVisible.value = true
  }

  const openDeduction = (item: Course) => {
    const status = getEffectiveStatus(item)
    if (status.disabled && status.label !== '课程暂停') {
      error(`操作失败：${status.desc}`)
      return
    }
    currentItem.value = { ...item }
    rechargeMode.value = 'refund'
    rechargeDialogVisible.value = true
  }

  // 处理子组件提交的数据
  const handleRechargeSubmit = async (data: any) => {
    // data: { courseId, hours, amount, remark }
    try {
      await RechargeCourse({
        courseId: data.courseId,
        hours: data.hours,
        amount: 0,
        remark: data.remark
      })

      // 更新前端数据
      const idx = courses.value.findIndex(c => c.id === data.courseId)
      if (idx !== -1) {
        courses.value[idx].balance += data.hours
        if (data.hours > 0) { // 仅充值增加累计
          courses.value[idx].totalBuy += data.hours
        }
      }

      // 通知子组件成功
      if (rechargeDialogRef.value) {
        rechargeDialogRef.value.onApiSuccess()
      }
      // loadData() // 可选：重新加载列表
    } catch (e: any) {
      error(e.message)
      // 通知子组件失败
      if (rechargeDialogRef.value) {
        rechargeDialogRef.value.onApiError()
      }
    }
  }

  const toggleStatus = async (item: Course) => {
    if (item.studentStatus !== 1) {
      error('请先恢复学员档案状态')
      return
    }
    try {
      await ToggleCourseStatus(item.id)
      success('状态已更新')
      loadData()
    } catch (e: any) {
      error(e.message)
    }
  }

  // --- 退课 / 删除 ---

  const openDelete = (item: Course) => {
    currentItem.value = { ...item }
    deleteForm.clearBalance = (item.balance || 0) > 0
    deleteForm.remark = ''
    deleteDialogVisible.value = true
  }

  const handleDeleteConfirm = async () => {
    loading.value = true
    try {
      await DeleteCourse({
        courseId: currentItem.value.id!,
        isHardDelete: false,
        remark: deleteForm.remark
      })
      success('已办理退课')
      deleteDialogVisible.value = false
      loadData()
    } catch (e: any) {
      error(e.message)
    } finally {
      loading.value = false
    }
  }

  const openForceDelete = (item: Course) => {
    currentItem.value = { ...item }
    forceDeleteDialogVisible.value = true
  }

  const handleForceDeleteConfirm = async () => {
    loading.value = true
    try {
      await DeleteCourse({
        courseId: currentItem.value.id!,
        isHardDelete: true
      })
      success('记录已彻底删除')
      forceDeleteDialogVisible.value = false
      loadData()
    } catch (e: any) {
      error(e.message)
    } finally {
      loading.value = false
    }
  }

  // --- UI Helpers ---

  const getEffectiveStatus = (item: Partial<Course>) => {
    if (item.studentStatus === 3) return { label: '学员退学', color: 'error', icon: 'mdi-account-off', disabled: true, desc: '该学员已退学，课程终止' }
    if (item.studentStatus === 2) return { label: '学员停课', color: 'blue-grey', icon: 'mdi-account-clock', disabled: true, desc: '因学员档案处于停课状态，该课程被冻结' }
    if (item.courseStatus === 3) return { label: '已结课', color: 'grey', icon: 'mdi-flag-checkered', disabled: true, desc: '课程已结束' }
    if (item.courseStatus === 2) return { label: '课程暂停', color: 'warning', icon: 'mdi-pause-circle', disabled: false, desc: '该课程已暂停，可恢复' }
    return { label: '正常上课', color: 'success', icon: 'mdi-check-circle', disabled: false, desc: '状态正常' }
  }

  const getBalanceColor = (b: number | undefined) => {
    const val = b || 0
    return val < 0 ? 'error' : (val < 5 ? 'warning' : 'success')
  }

  const getBalanceLabel = (b: number | undefined) => {
    const val = b || 0
    return val < 0 ? '已欠费' : (val < 5 ? '余额不足' : '余额充足')
  }

  const formatBalance = (b: number | undefined) => {
    const val = b || 0
    if (val > 99) return '99+'
    if (val < -99) return '-99+'
    const sign = val < 0 ? '-' : ''
    return `${sign}${Math.abs(val).toString().padStart(2, '0')}`
  }

  const headers: any = [
    { title: '学员', key: 'studentName', align: 'start', width: '20%' },
    { title: '科目', key: 'subjectName', align: 'center', width: '15%' },
    { title: '授课老师', key: 'teacherName', align: 'center', width: '20%' },
    { title: '剩余课时', key: 'balance', align: 'center', width: '15%' },
    { title: '有效状态', key: 'status', align: 'center', width: '20%' },
    { title: '操作', key: 'actions', align: 'end', sortable: false, width: '10%' },
  ]

  // 监听筛选变化
  watch(filters, () => {
    page.value = 1
    loadData()
  }, { deep: true })

  onMounted(() => {
    loadData()
  })

  return {
    loading, page, itemsPerPage, totalItems, courses, headers,
    dialogVisible, rechargeDialogVisible, deleteDialogVisible, forceDeleteDialogVisible,
    isEdit, currentItem,
    rechargeMode, rechargeDialogRef, // 导出 ref
    deleteForm, isDeleteValid, enrollForm,
    filters, activeFilters,
    studentOptions, subjectOptions, teacherOptions, statusOptions,
    isStudentLoading, isSubjectLoading, isTeacherLoading,
    onStudentSearch, onSubjectSearch, onTeacherSearch,
    loadData, clearFilter,
    openEnroll, openEdit, handleSave,
    openRecharge, openDeduction, handleRechargeSubmit,
    toggleStatus,
    openDelete, handleDeleteConfirm,
    openForceDelete, handleForceDeleteConfirm,
    getEffectiveStatus, getBalanceColor, getBalanceLabel, formatBalance
  }
}