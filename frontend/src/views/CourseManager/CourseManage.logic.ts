import { ref, reactive, computed, onMounted, onActivated, watch, nextTick } from 'vue'
import { debounce } from 'lodash'
import { useToast } from '../../composables/useToast'
import {
  GetCourseList, CreateCourse, RechargeCourse, ToggleCourseStatus, DeleteCourse, UpdateCourse,
} from '../../api/course'
import { GetStudentList } from '../../api/student'
import { GetSubjectList } from '../../api/subject'
import { GetTeacherList } from '../../api/teacher'
import { GetStudentListResponse, GetSubjectListResponse, GetTeacherListResponse } from '../../types/response'
import { CreateCourseRequest, GetCourseListRequest, GetStudentListRequest } from '../../types/request'
import { Course } from '../../types/appModels'

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
    studentId: null as number | null,
    studentNameLabel: '',
    subjects: [] as number[],
    teachers: [] as number[],
    balanceMin: null as number | null,
    balanceMax: null as number | null,
    status: [] as number[]
  })

  // 字典数据
  const studentOptions = ref<{ title: string, value: number, name: string }[]>([])
  const enrollStudentOptions = ref<{ title: string, value: number, name: string }[]>([])
  const subjectOptions = ref<{ title: string, value: number }[]>([])
  const teacherOptions = ref<{ title: string, value: number, name: string }[]>([])

  const isStudentLoading = ref(false)
  const isEnrollStudentLoading = ref(false)
  const isSubjectLoading = ref(false)
  const isTeacherLoading = ref(false)

  const statusOptions = [
    { title: '正常上课', value: 1, color: 'success' },
    { title: '课程暂停', value: 2, color: 'warning' },
    { title: '已结课', value: 3, color: 'grey' },
    { title: '学员停课', value: 4, color: 'blue-grey' },
    { title: '学员退学', value: 5, color: 'error' }
  ]

  // --- 搜索方法 ---
  const searchStudents = async (keyword: string) => {
    if (!keyword) return
    isStudentLoading.value = true
    try {
      let reqData: GetStudentListRequest = { Offset: 0, Limit: 25, keyword, Status_Level: 3, Status_Target: 0 }
      const res: GetStudentListResponse = await GetStudentList(reqData)
      let students = res.students
      const newOptions = students.map(s => ({ title: `${s.name} (${s.student_number})`, value: s.id, name: s.name }))

      // 保留当前选中的项
      if (filters.studentId) {
        const selected = studentOptions.value.find(o => o.value === filters.studentId)
        if (selected && !newOptions.find(o => o.value === filters.studentId)) {
          newOptions.push(selected)
        }
      }
      studentOptions.value = newOptions
    } catch (e) {
      console.error(e)
    } finally {
      isStudentLoading.value = false
    }
  }

  const searchEnrollStudents = async (keyword: string) => {
    if (!keyword) return
    isEnrollStudentLoading.value = true
    try {
      let reqData: GetStudentListRequest = { Offset: 0, Limit: 25, keyword, Status_Level: 1, Status_Target: 1 }
      const res: GetStudentListResponse = await GetStudentList(reqData)
      let students = res.students
      const newOptions = students.map(s => ({ title: `${s.name} (${s.student_number})`, value: s.id, name: s.name }))

      // 保留当前选中的项
      if (enrollForm.studentId) {
        const selected = enrollStudentOptions.value.find(o => o.value === enrollForm.studentId)
        if (selected && !newOptions.find(o => o.value === enrollForm.studentId)) {
          newOptions.push(selected)
        }
      }
      enrollStudentOptions.value = newOptions
    } catch (e) {
      console.error(e)
    } finally {
      isEnrollStudentLoading.value = false
    }
  }

  const searchSubjects = async (keyword: string) => {
    isSubjectLoading.value = true
    try {
      const res: GetSubjectListResponse = await GetSubjectList({ Offset: 0, Limit: 25, Keyword: keyword })
      let subjectsList = res.subjects || []
      const newOptions = subjectsList.map(s => ({ title: s.name, value: s.id }))

      // 保留筛选中选中的项
      filters.subjects.forEach(id => {
        const selected = subjectOptions.value.find(o => o.value === id)
        if (selected && !newOptions.find(o => o.value === id)) {
          newOptions.push(selected)
        }
      })
      // 保留表单中选中的项
      if (enrollForm.subjectId) {
        const selected = subjectOptions.value.find(o => o.value === enrollForm.subjectId)
        if (selected && !newOptions.find(o => o.value === enrollForm.subjectId)) {
          newOptions.push(selected)
        }
      }

      subjectOptions.value = newOptions
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
      const newOptions = teachersList.map(t => ({ title: `${t.name} (${t.teacher_number})`, value: t.id, name: t.name }))

      // 保留筛选中选中的项
      filters.teachers.forEach(id => {
        const selected = teacherOptions.value.find(o => o.value === id)
        if (selected && !newOptions.find(o => o.value === id)) {
          newOptions.push(selected)
        }
      })
      // 保留表单中选中的项
      if (enrollForm.teacherId) {
        const selected = teacherOptions.value.find(o => o.value === enrollForm.teacherId)
        if (selected && !newOptions.find(o => o.value === enrollForm.teacherId)) {
          newOptions.push(selected)
        }
      }

      teacherOptions.value = newOptions
    } catch (e) {
      console.error(e)
    } finally {
      isTeacherLoading.value = false
    }
  }

  const onStudentSearch = debounce(searchStudents, 500)
  const onEnrollStudentSearch = debounce(searchEnrollStudents, 500)
  const onSubjectSearch = debounce(searchSubjects, 500)
  const onTeacherSearch = debounce(searchTeachers, 500)

  // --- 辅助计算属性 ---

  const activeFilters = computed(() => {
    const list = []
    if (filters.studentId) {
      list.push({ key: 'studentId', label: `学员: ${filters.studentNameLabel || filters.studentId}` })
    }
    if (filters.subjects.length > 0) {
      const names = filters.subjects.map(id => subjectOptions.value.find(o => o.value === id)?.title || id)
      list.push({ key: 'subjects', label: `科目: ${names.join(', ')}` })
    }
    if (filters.teachers.length > 0) {
      const names = filters.teachers.map(id => teacherOptions.value.find(o => o.value === id)?.name || id)
      list.push({ key: 'teachers', label: `老师: ${names.join(', ')}` })
    }
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
      // 构建请求参数
      const req: GetCourseListRequest = {
        Students: filters.studentId ? [filters.studentId] : undefined,
        Subjects: filters.subjects.length ? filters.subjects : undefined,
        Teachers: filters.teachers.length ? filters.teachers : undefined,
        Balance_Min: filters.balanceMin !== null && filters.balanceMin !== '' as any ? Number(filters.balanceMin) : undefined,
        Balance_Max: filters.balanceMax !== null && filters.balanceMax !== '' as any ? Number(filters.balanceMax) : undefined,
        Status: filters.status.length ? filters.status : undefined,
        Offset: (page.value - 1) * itemsPerPage.value,
        Limit: itemsPerPage.value
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
    if (key === 'studentId') {
      filters.studentId = null
      filters.studentNameLabel = ''
    }
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
    enrollStudentOptions.value = []
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
    enrollForm.teacherId = item.teacher.id as any
    enrollForm.remark = ''

    // 初始化选项，确保当前值能正确显示
    teacherOptions.value = [{ title: item.teacher.name, value: item.teacher.id, name: item.teacher.name }]
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
          teacher_id: enrollForm.teacherId!,
          remark: enrollForm.remark
        })
        success('课程信息更新成功')
      } else {
        let reqData: CreateCourseRequest = {
          student_Id: enrollForm.studentId!,
          subject_Id: enrollForm.subjectId!,
          teacher_Id: enrollForm.teacherId!,
          remark: enrollForm.remark
        }
        await CreateCourse(reqData)
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
        course_id: data.courseId,
        hours: data.hours,
        amount: data.amount,
        remark: data.remark
      })

      // 更新前端数据
      const idx = courses.value.findIndex(c => c.id === data.courseId)
      if (idx !== -1) {
        courses.value[idx].balance += data.hours
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
    if (item.student.status !== 1) {
      error('请先恢复学员档案状态')
      return
    }
    try {
      await ToggleCourseStatus(item.id)
      success('状态已更新')
      // 本地更新状态，避免重新加载整个列表
      item.status = item.status === 1 ? 2 : 1
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
        course_id: currentItem.value.id!,
        is_hard_delete: false,
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
        course_id: currentItem.value.id!,
        is_hard_delete: true
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
    if (!item.student) {
      return { label: '未知状态', color: 'grey', icon: 'mdi-help-circle', disabled: true, desc: '学员信息缺失' }
    }
    if (item.student.status === 3) return { label: '学员退学', color: 'error', icon: 'mdi-account-off', disabled: true, desc: '该学员已退学，课程终止' }
    if (item.student.status === 2) return { label: '学员停课', color: 'blue-grey', icon: 'mdi-account-clock', disabled: true, desc: '因学员档案处于停课状态，该课程被冻结' }
    if (item.status === 3) return { label: '已结课', color: 'grey', icon: 'mdi-flag-checkered', disabled: true, desc: '课程已结束' }
    if (item.status === 2) return { label: '课程暂停', color: 'warning', icon: 'mdi-pause-circle', disabled: false, desc: '该课程已暂停，可恢复' }
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

  watch(() => filters.studentId, (val) => {
    if (val) {
      const s = studentOptions.value.find(o => o.value === val)
      if (s) filters.studentNameLabel = s.name
    } else {
      filters.studentNameLabel = ''
    }
  })

  onMounted(() => {
    loadData()
  })
  onActivated(() => {
    loadData()
  })

  return {
    loading, page, itemsPerPage, totalItems, courses, headers,
    dialogVisible, rechargeDialogVisible, deleteDialogVisible, forceDeleteDialogVisible,
    isEdit, currentItem,
    rechargeMode, rechargeDialogRef, // 导出 ref
    deleteForm, isDeleteValid, enrollForm,
    filters, activeFilters,
    studentOptions, enrollStudentOptions, subjectOptions, teacherOptions, statusOptions,
    isStudentLoading, isEnrollStudentLoading, isSubjectLoading, isTeacherLoading,
    onStudentSearch, onEnrollStudentSearch, onSubjectSearch, onTeacherSearch,
    loadData, clearFilter,
    openEnroll, openEdit, handleSave,
    openRecharge, openDeduction, handleRechargeSubmit,
    toggleStatus,
    openDelete, handleDeleteConfirm,
    openForceDelete, handleForceDeleteConfirm,
    getEffectiveStatus, getBalanceColor, getBalanceLabel, formatBalance
  }
}