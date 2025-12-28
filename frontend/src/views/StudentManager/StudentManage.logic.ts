import { ref, reactive, computed, onMounted } from 'vue'
import type { StudentData, StudentItem, TeacherOption } from '../../types/appModels'
import { useToast } from '../../composables/useToast'
import { useConfirm } from '../../composables/useConfirm'
import { categorizeOrderTags } from '../../utils/classification'
import type { OrderTag } from '../../types/appModels'

export function useStudentManage() {
  const { success, error, info } = useToast()
  const { confirmDelete } = useConfirm()

  const search = ref('')

  // 弹窗状态
  const dialog = ref(false)
  const dialogRecharge = ref(false)
  const dialogDetails = ref(false)

  // 数据对象
  const defaultItem: StudentData = { id: 0, name: '', phone: '', balance: 0, gender: '男', teacher_id: null, note: '' }
  const editedItem = reactive<StudentData>({ ...defaultItem })
  const editedIndex = ref(-1)

  const rechargeItem = reactive<StudentData>({ ...defaultItem })
  const rechargeForm = reactive({ amount: 10, note: '' })

  // 列表数据 (初始化为空)
  const students = ref<StudentItem[]>([])

  // 教师列表 (用于搜索)
  const teacherOptions = ref<TeacherOption[]>([])

  // --- 初始化加载 ---
  const loadData = async () => {
    try {
      // TODO: API - 调用后端获取学生列表
      // const res = await window.go.main.App.GetAllStudents()
      // students.value = res || []

      // TODO: API - 调用后端获取教师列表
      // const teachers = await window.go.main.App.GetAllTeachers()
      // teacherOptions.value = teachers || []

      console.log('Fetching student data...')
    } catch (e) {
      error('加载数据失败')
    }
  }

  onMounted(() => {
    loadData()
  })

  // --- 表头定义 ---
  const headers: any = [
    { title: '姓名', key: 'name', align: 'center', sortable: false, width: '120px' },
    { title: '剩余课时', key: 'balance', sortable: false, width: '120px' },
    { title: '状态', key: 'status', sortable: false, width: '100px' },
    { title: '操作', key: 'actions', sortable: false, align: 'end', width: '180px' },
  ]

  // --- 辅助函数 ---
  const getStatusColor = (bal: number) => bal < 0 ? 'error' : (bal < 5 ? 'warning' : 'success')
  const getStatusText = (bal: number) => bal < 0 ? '欠费' : (bal < 5 ? '余额不足' : '正常')

  // --- 充值逻辑计算 ---
  const rechargeType = computed(() => {
    if (!rechargeForm) return '无变动'
    const amt = Number(rechargeForm.amount)
    if (amt > 0) return '充值/赠送'
    if (amt < 0) return '退费/扣减'
    return '无变动'
  })

  const rechargeColor = computed(() => {
    if (!rechargeForm) return 'grey'
    const amt = Number(rechargeForm.amount)
    if (amt > 0) return 'success'
    if (amt < 0) return 'error'
    return 'grey'
  })

  const inferredTags = computed<OrderTag[]>(() => {
    if (!rechargeForm) return []
    return categorizeOrderTags(rechargeForm.note);
  });

  // --- 操作方法 ---
  const openAdd = () => {
    editedIndex.value = -1
    Object.assign(editedItem, defaultItem)
    dialog.value = true
  }

  const openEdit = (item: StudentItem) => {
    editedIndex.value = students.value.findIndex(s => s.id === item.id)
    Object.assign(editedItem, item)
    dialog.value = true
  }

  const openRecharge = (item: StudentItem) => {
    Object.assign(rechargeItem, item)
    rechargeForm.amount = 10
    rechargeForm.note = ''
    dialogRecharge.value = true
  }

  const openDetails = (item: StudentItem) => {
    Object.assign(editedItem, item)
    dialogDetails.value = true
  }

  const openDelete = async (item: StudentItem) => {
    const ok = await confirmDelete(item.name)
    if (ok) {
      try {
        // TODO: API - 调用后端删除学生
        // await window.go.main.App.DeleteStudent(item.id)

        // 前端移除 (实际应重新加载列表)
        const idx = students.value.findIndex(s => s.id === item.id)
        if (idx > -1) students.value.splice(idx, 1)

        success('删除成功')
      } catch (e) {
        error('删除失败')
      }
    }
  }

  const saveStudent = async (data: StudentData) => {
    try {
      if (editedIndex.value > -1) {
        // TODO: API - 更新学生
        // await window.go.main.App.UpdateStudent(data)

        // 前端模拟更新
        if (students.value[editedIndex.value]) {
          Object.assign(students.value[editedIndex.value], data)
        }
        success('更新成功')
      } else {
        // TODO: API - 新增学生
        // const newId = await window.go.main.App.CreateStudent(data)
        // data.id = newId

        // 前端模拟新增
        students.value.push({ ...data } as StudentItem)
        success('添加成功')
      }
      dialog.value = false
    } catch (e) {
      error('操作失败')
    }
  }

  const saveRecharge = async () => {
    try {
      // TODO: API - 提交充值记录
      // await window.go.main.App.AddOrder({
      //   student_id: rechargeItem.id,
      //   amount: Number(rechargeForm.amount),
      //   comment: rechargeForm.note
      // })

      const target = students.value.find(s => s.id === rechargeItem.id)
      if (target) {
        target.balance += Number(rechargeForm.amount)
        success('课时调整成功')
      }
      dialogRecharge.value = false
    } catch (e) {
      error('充值失败')
    }
  }

  const exportStudents = async () => {
    try {
      info('正在导出学生数据...')
      // TODO: API - 调用后端导出 Excel
      // await window.go.main.App.ExportStudentList()
      success('导出成功')
    } catch (e) {
      error('导出失败')
    }
  }

  return {
    search,
    dialog, dialogRecharge, dialogDetails,
    editedIndex, editedItem, rechargeItem, rechargeForm,
    students, teacherOptions,
    headers,
    rechargeType, rechargeColor, inferredTags,
    getStatusColor, getStatusText,
    openAdd, openEdit, openRecharge, openDetails, openDelete,
    saveStudent, saveRecharge, exportStudents
  }
}