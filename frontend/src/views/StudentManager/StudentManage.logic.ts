import { ref, computed, nextTick } from 'vue'
import type { ResponseWrapper } from '../../types/appModels'
import type { StudentDTO } from '../../types/response'
import { useToast } from '../../composables/useToast'
import { LogError } from '../../../wailsjs/runtime/runtime'
import { Dispatch } from '../../../wailsjs/go/main/App'
// --- 类型定义 (TypeScript Interface) ---

const toast = useToast()

interface Student {
  id: number
  name: string
  phone: string
  balance: number
  gender: string
  teacher_id: number
  note: string
}

export function useStudentManage() {
  // --- 页面状态 ---
  const search = ref('')
  const dialog = ref(false)         // 新增/编辑弹窗状态
  const dialogRecharge = ref(false) // 充值弹窗状态
  const dialogDelete = ref(false)   // 删除确认弹窗状态
  const dialogDetails = ref(false) // 详情弹窗状态

  // --- 表单状态 ---
  const editedIndex = ref(-1) // -1 表示新增模式，>=0 表示编辑模式
  const defaultItem: Student = {
    id: 0,
    name: '',
    phone: '',
    balance: 0,
    gender: '',
    teacher_id: 0,
    note: ''
  }

  // 当前操作的对象副本
  const editedItem = ref<Student>({ ...defaultItem })
  const rechargeItem = ref<Student>({ ...defaultItem })

  // 充值表单数据
  const rechargeForm = ref({
    type: '充值',
    amount: 10,
    note: ''
  })

  // 实际开发中，这里会替换为从 Go 后端获取的数据
  const students = ref<Student[]>([])

  // 从后端中获取学生列表（使用 Promise .then() 风格）
  const fetchStudents = () => {

    const reqData = {
      Key: search.value,
      Offset: 0,
      Limit: 100,
    }


    Dispatch('student_manager:get_student_list', JSON.stringify(reqData))
      .then((result: any) => {
        // ResponseWrapper<StudentDTO[]> 解析
        const resp = JSON.parse(result) as ResponseWrapper<StudentDTO[]>
        if (resp.code === 200) {
          students.value = resp.data.map((item) => ({
            id: item.id,
            name: item.name,
            phone: item.phone,
            balance: item.hours,
            gender: item.gender,
            teacher_id: item.teacher_id,
            note: '', // 后端暂无 remark 字段
          }))
        } else {
          LogError('获取学生列表失败:' + resp.message)
          toast.error('获取学生列表失败: ' + resp.message, 'top-right')
        }

      })
  }

  // 初始化获取学生列表
  fetchStudents()

  // --- 表格配置 ---
  const headers: any = [
    { title: '姓名', key: 'name', align: 'start', sortable: false, width: '120px' },
    { title: '剩余课时', key: 'balance', sortable: false, width: '120px' },
    { title: '状态', key: 'status', sortable: false, width: '100px' }, // 虚拟列
    { title: '操作', key: 'actions', sortable: false, align: 'end', width: '150px' },
  ]

  // --- 计算属性 ---
  const formTitle = computed(() => {
    return editedIndex.value === -1 ? '新增学生' : '编辑学生'
  })

  // --- 核心业务逻辑 ---

  // 1. 状态显示逻辑
  const getStatusColor = (balance: number) => {
    if (balance < 0) return 'error'   // 红色
    if (balance < 5) return 'warning' // 橙色
    return 'success'                  // 绿色
  }

  const getStatusText = (balance: number) => {
    if (balance < 0) return '欠费'
    if (balance < 5) return '余额不足'
    return '正常'
  }

  // 2. CRUD 逻辑
  const openEdit = (item: Student) => {
    editedIndex.value = students.value.indexOf(item)
    editedItem.value = { ...item } // 复制对象，防止直接修改原数据
    dialog.value = true
  }

  const openDelete = (item: Student) => {
    editedIndex.value = students.value.indexOf(item)
    editedItem.value = { ...item }
    dialogDelete.value = true
  }

  const deleteItemConfirm = () => {
    students.value.splice(editedIndex.value, 1) // 模拟删除
    closeDelete()
  }

  const closeDialog = () => {
    dialog.value = false
    nextTick(() => {
      editedItem.value = { ...defaultItem }
      editedIndex.value = -1
    })
  }

  const closeDelete = () => {
    dialogDelete.value = false
    nextTick(() => {
      editedItem.value = { ...defaultItem }
      editedIndex.value = -1
    })
  }

  const saveStudent = () => {
    if (editedIndex.value > -1) {
      // 编辑保存
      Object.assign(students.value[editedIndex.value], editedItem.value)
    } else {
      // 新增保存 (生成一个简单的模拟ID)
      editedItem.value.id = Math.max(...students.value.map(s => s.id), 0) + 1
      students.value.push(editedItem.value)
    }
    closeDialog()
  }

  // 3. 充值逻辑
  const openRecharge = (item: Student) => {
    rechargeItem.value = { ...item }
    dialogRecharge.value = true
  }

  const saveRecharge = () => {
    const index = students.value.findIndex(s => s.id === rechargeItem.value.id)
    if (index !== -1) {
      const amount = Number(rechargeForm.value.amount)

      // 根据类型计算余额
      if (rechargeForm.value.type === '退费') {
        students.value[index].balance -= amount
      } else {
        // 充值 或 赠送
        students.value[index].balance += amount
      }
    }
    dialogRecharge.value = false
    // 重置表单
    rechargeForm.value = { type: '充值', amount: 10, note: '' }
  }

  return {
    search,
    dialog,
    dialogRecharge,
    dialogDelete,
    dialogDetails,
    editedIndex,
    editedItem,
    rechargeItem,
    rechargeForm,
    students,
    headers,
    formTitle,
    getStatusColor,
    getStatusText,
    openEdit,
    openDelete,
    deleteItemConfirm,
    closeDialog,
    closeDelete,
    saveStudent,
    openRecharge,
    saveRecharge,
  }
}
