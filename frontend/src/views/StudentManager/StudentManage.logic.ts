import { ref, reactive, onMounted, computed } from 'vue'
import type { ResponseWrapper, StudentData, Student, TeacherOption } from '../../types/appModels'
import { useToast } from '../../composables/useToast'
import { useConfirm } from '../../composables/useConfirm'
import { Dispatch } from '../../../wailsjs/go/main/App'
import { StudentDTO } from '../../types/response'
import { CreateStudent, DeleteStudent, ExportStudentsToExcel, GetStudentList, UpdateStudent } from '../../api/student'
import { CreateStudentRequest, DeleteStudentRequest, UpdateStudentRequest } from '../../types/request'

const { success, error, info } = useToast()
const { confirmDelete } = useConfirm()

const search = ref('')
const filterStatus = ref<number | null>(null) // 新增：状态筛选
const toast = useToast()
const confirm = useConfirm()

// 分页状态
const page = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0) // 需要后端支持返回总数
const loading = ref(false)

// 弹窗状态
const dialog = ref(false)

// 数据对象
const defaultItem: StudentData = { id: 0, name: '', phone: '', gender: 'male', remark: '', status: 1 }
const editedItem = reactive<StudentData>({ ...defaultItem })
const editedIndex = ref(-1)

const rechargeItem = reactive<StudentData>({ ...defaultItem })

// 列表数据 (初始化为空)
const students = ref<Student[]>([])

// 计算属性：是否有激活的筛选
const hasActiveFilters = computed(() => {
  return filterStatus.value !== null
})

// 清除所有筛选
const clearAllFilters = () => {
  filterStatus.value = null
  loadData()
}


// --- 初始化加载 ---
const loadData = async () => {
  await loadStudents()
}

async function loadStudents(): Promise<void> {
  loading.value = true
  try {
    // 1. 发起请求
    const data = await GetStudentList({
      Offset: (page.value - 1) * itemsPerPage.value,
      Limit: itemsPerPage.value,
      keyword: search.value,
      Status: filterStatus.value == null ? 0 : filterStatus.value // 传递状态筛选
    })

    console.log('Fetched student list:', data)

    // 2. 处理成功数据 (数据映射/格式化在这里做)
    students.value = (data.students || []).map((item: StudentDTO) => ({
      id: item.id,
      name: item.name,
      phone: item.phone,
      gender: item.gender,
      status: item.status,
      remark: item.remark || '',
      student_number: item.student_number,
      // 格式化日期
      lastModified: item.updated_at ? new Date(item.updated_at).toLocaleString() : '-',
    }))

    totalItems.value = data.total

  } catch (error: any) {
    console.error('获取学生列表失败:', error.message)
    toast.error('获取学生列表失败: ' + error.message, 'top-right')
  } finally {
    loading.value = false
  }
}



function loadItems({ page: newPage, itemsPerPage: newItemsPerPage, sortBy }: { page: number; itemsPerPage: number; sortBy?: string[] | string | undefined }): void {
  console.log('Loading items with params:', { page: newPage, itemsPerPage: newItemsPerPage, sortBy })
  page.value = newPage
  itemsPerPage.value = newItemsPerPage

  loadData()
}

// 组件挂载时加载数据
onMounted(() => {
  loadData()
})

// --- 表头定义 ---
const headers: any = [
  { title: '学号', key: 'student_number', align: 'center', sortable: false },
  { title: '姓名', key: 'name', align: 'center', sortable: false },
  { title: '性别', key: 'gender', sortable: false, align: 'center' },
  { title: '手机号', key: 'phone', sortable: false, align: 'center' },
  { title: '状态', key: 'status', sortable: false, align: 'center' },
  { title: '操作', key: 'actions', sortable: false, align: 'center' },
]

// --- 辅助函数 ---
// status 0 空状态 1 正常上课 2 暂时听课 3 退学
const getStatusColor = (status: number) => {
  const map: Record<number, string> = {
    0: 'grey-lighten-3',
    1: 'green-lighten-4 text-green-darken-2',
    2: 'orange-lighten-4 text-orange-darken-2',
    3: 'red-lighten-4 text-red-darken-2',
  }
  return map[status] || 'grey-lighten-3'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '空状态',
    1: '正常',
    2: '停课',
    3: '退学',
  }
  return map[status] || '未知'
}

// 性别文本转换
const getGenderLabel = (gender: string) => {
  const map: Record<string, string> = {
    'male': '男',
    'female': '女'
  }
  return map[gender] || gender || '未知'
}

// 性别颜色区分
const getGenderColor = (gender: string) => {
  if (gender === '1' || gender === 'male') return 'blue-lighten-4 text-blue-darken-2'
  if (gender === '2' || gender === 'female') return 'pink-lighten-4 text-pink-darken-2'
  return 'grey-lighten-3'
}

async function createStudent(data: StudentData) {
  let reqData: CreateStudentRequest = {
    Name: data.name,
    Gender: data.gender,
    Phone: data.phone,
    Remark: data.remark,
  }

  console.log('Creating student with data:', reqData)
  // 调用 API
  try {
    await CreateStudent(reqData)
    loadData()
    success('添加成功', 'top-right')
  } catch (e) {
    console.error('创建学生异常:', e)
    toast.error('创建学生异常: ' + (e as Error).message, 'top-right')
    return
  }
}

async function updateStudent(data: StudentData) {

  if (!data.id) {
    console.error('更新学生失败: 缺少学生ID')
    toast.error('更新学生失败: 缺少学生ID', 'top-right')
    return
  }

  let reqData: UpdateStudentRequest = {
    ID: data.id,
    Name: data.name,
    Gender: data.gender,
    Phone: data.phone,
    Remark: data.remark,
    Status: data.status,
  }

  console.log('Updating student with data:', reqData)

  try {
    await UpdateStudent(reqData)
    loadData()
    success('更新成功', 'top-right')
  } catch (e) {
    console.error('更新学生异常:', e)
    toast.error('更新学生异常: ' + (e as Error).message, 'top-right')
    return
  }
}


async function deleteStudent(item: Student) {
  console.log('Deleting student with ID:', item.id)
  let reqData: DeleteStudentRequest = { ID: item.id }
  try {
    if (!item.id) {
      console.error('删除学生失败: 缺少学生ID')
      toast.error('删除学生失败: 缺少学生ID', 'top-right')
      return
    }

    await DeleteStudent(reqData)
    loadData()
    success('删除成功', 'top-right')
  } catch (e) {
    console.error('删除学生异常:', e)
    toast.error('删除学生异常: ' + (e as Error).message, 'top-right')
    return
  }
}

async function exportStudents2Excel() {
  console.log('Exporting student data...')

  try {
    info('正在导出学生数据，请稍候...', 'top-right')
    let data: string = await ExportStudentsToExcel()
    if (data === 'cancel') {
      info('已取消导出', 'top-right')
      return
    }
    success('学生数据已导出到: ' + data, 'top-right')
  } catch (e) {
    if (e instanceof Error) {
      toast.error('导出学生数据异常: ' + e.message, 'top-right')
    }
    console.error('导出提示异常:', e)
  }
}



export function useStudentManage() {

  // --- 操作方法 ---
  const openAdd = () => {
    editedIndex.value = -1
    Object.assign(editedItem, defaultItem)
    dialog.value = true
  }

  const openEdit = (item: Student) => {
    editedIndex.value = students.value.findIndex(s => s.id === item.id)
    Object.assign(editedItem, item)
    dialog.value = true
  }


  const openDelete = async (item: Student) => {
    const ok = await confirmDelete(item.name)
    if (ok) {
      deleteStudent(item)
    }
  }

  const saveStudent = async (data: StudentData) => {
    try {
      if (editedIndex.value > -1) {
        updateStudent(data)
      } else {
        createStudent(data)
      }
    } catch (e) {
      error('操作失败')
    }
  }

  const exportStudents = async () => {
    exportStudents2Excel()
  }



  return {
    search,
    filterStatus,
    hasActiveFilters,
    clearAllFilters,
    page,
    itemsPerPage,
    totalItems,
    loading,
    dialog,
    editedIndex, editedItem,
    students,
    headers,
    loadItems,
    getStatusColor, getStatusText, getGenderColor, getGenderLabel,
    openAdd, openEdit, openDelete,
    saveStudent, exportStudents
  }
}