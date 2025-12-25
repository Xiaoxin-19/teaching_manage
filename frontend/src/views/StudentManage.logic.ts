import { ref, computed, nextTick } from 'vue'

// --- 类型定义 (TypeScript Interface) ---
interface Student {
  id: number
  name: string
  phone: string
  balance: number
  note: string
}

export function useStudentManage() {
  // --- 页面状态 ---
  const search = ref('')
  const dialog = ref(false)         // 新增/编辑弹窗状态
  const dialogRecharge = ref(false) // 充值弹窗状态
  const dialogDelete = ref(false)   // 删除确认弹窗状态

  // --- 表单状态 ---
  const editedIndex = ref(-1) // -1 表示新增模式，>=0 表示编辑模式
  const defaultItem: Student = {
    id: 0,
    name: '',
    phone: '',
    balance: 0,
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

  // --- 模拟数据 (Mock Data) ---
  // 实际开发中，这里会替换为从 Go 后端获取的数据
  const students = ref<Student[]>([
    { id: 1, name: '张子轩', phone: '13812345678', balance: 24, note: '钢琴基础' },
    { id: 2, name: '李梓涵', phone: '13987654321', balance: 3, note: '周六上午班' },
    { id: 3, name: '王浩宇', phone: '15011112222', balance: -2, note: '需催费' },
    { id: 4, name: '陈思睿', phone: '13666668888', balance: 45, note: '' },
    { id: 5, name: '刘一诺', phone: '18999990000', balance: 10, note: '' },
  ])

  // --- 表格配置 ---
  const headers: any = [
    { title: '姓名', key: 'name', align: 'start', width: '120px' },
    { title: '联系电话', key: 'phone', width: '150px' },
    { title: '剩余课时', key: 'balance', sortable: true, width: '120px' },
    { title: '状态', key: 'status', sortable: false, width: '100px' }, // 虚拟列
    { title: '备注', key: 'note' },
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
