import { ref, reactive, computed, onMounted } from 'vue'
import type { ResponseWrapper, StudentData, StudentItem, TeacherOption } from '../../types/appModels'
import { useToast } from '../../composables/useToast'
import { useConfirm } from '../../composables/useConfirm'
import { Dispatch } from '../../../wailsjs/go/main/App'
import { GetStudentListResponse, GetTeacherListResponse, StudentDTO } from '../../types/response'
import { LogError } from '../../../wailsjs/runtime/runtime'

const { success, error, info } = useToast()
const { confirmDelete } = useConfirm()

const search = ref('')
const toast = useToast()
const confirm = useConfirm()

// 分页状态
const page = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0) // 需要后端支持返回总数
const loading = ref(false)

// 弹窗状态
const dialog = ref(false)
const dialogRecharge = ref(false)
const dialogDetails = ref(false)

// 数据对象
const defaultItem: StudentData = { id: 0, name: '', phone: '', balance: 0, gender: 'male', teacher_id: null, note: '' }
const editedItem = reactive<StudentData>({ ...defaultItem })
const editedIndex = ref(-1)

const rechargeItem = reactive<StudentData>({ ...defaultItem })

// 列表数据 (初始化为空)
const students = ref<StudentItem[]>([])

// 教师列表 (用于搜索)
const teacherOptions = ref<TeacherOption[]>([])

// --- 初始化加载 ---
const loadData = async () => {
  await loadStudents()
}

async function loadStudents(): Promise<void> {
  loading.value = true
  try {
    const reqData = {
      Key: search.value,
      Offset: (page.value - 1) * itemsPerPage.value,
      Limit: itemsPerPage.value,
    }

    console.log('Request Data:', reqData)


    const result: any = await Dispatch('student_manager:get_student_list', JSON.stringify(reqData))
    // ResponseWrapper<StudentDTO[]> 解析
    const resp = JSON.parse(result) as ResponseWrapper<GetStudentListResponse>
    if (resp.code === 200) {
      students.value = (resp.data.students || []).map((item) => ({
        id: item.id,
        name: item.name,
        phone: item.phone,
        balance: item.hours,
        gender: item.gender,
        teacher_id: item.teacher_id,
        note: item.remark, // 后端暂无 remark 字段
        lastModified: new Date(item.updated_at).toLocaleString(),
      }))
      totalItems.value = resp.data.total // 需要后端支持返回总数
    } else {
      console.error('获取学生列表失败:', resp.message)
      toast.error('获取学生列表失败: ' + resp.message, 'top-right')
    }

    console.log('Fetching data...')
  } catch (e) {
    console.error('加载学生列表时出错:', e)
  } finally {
    loading.value = false
  }
}


function loadTeacherOptions() {
  const requestData = {
    Key: "",
    Offset: 0,
    Limit: -1,
  }

  console.log('Fetching teacher list with request:' + JSON.stringify(requestData))
  Dispatch('teacher_manager:get_teacher_list', JSON.stringify(requestData))
    .then((result: any) => {
      console.log('Received teacher list response:' + result)
      const resp = JSON.parse(result) as ResponseWrapper<GetTeacherListResponse>
      if (resp.code === 200) {
        // 如果teacher为null则赋值为空数组
        if (!resp.data.teachers) {
          resp.data.teachers = []
        }


        teacherOptions.value = resp.data.teachers.map((item) => ({
          id: item.id,
          name: item.name,
          phone: item.phone,
          gender: item.gender,
          remark: item.remark,
          lastModified: new Date(item.updated_at).toLocaleString(),
        }))
        totalItems.value = resp.data.total
      } else {
        LogError('获取教师列表失败:' + resp.message)
        toast.error('获取教师列表失败: ' + resp.message, 'top-right')
      }
    })
    .finally(() => {
      loading.value = false
    })

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
  { title: '姓名', key: 'name', align: 'center', sortable: false },
  { title: '剩余课时', key: 'balance', sortable: false, align: 'center' },
  { title: '状态', key: 'status', sortable: false },
  { title: '操作', key: 'actions', sortable: false, align: 'center' },
]

// --- 辅助函数 ---
const getStatusColor = (bal: number) => bal < 0 ? 'error' : (bal < 5 ? 'warning' : 'success')
const getStatusText = (bal: number) => bal < 0 ? '欠费' : (bal < 5 ? '余额不足' : '正常')


function createStudent(data: StudentData) {
  let reqData = {
    Name: data.name,
    Gender: data.gender,
    Hours: data.balance,
    Phone: data.phone,
    Teacher_Id: data.teacher_id,
    Remark: data.note,
  }

  console.log('Creating student with data:', reqData)
  Dispatch('student_manager:create_student', JSON.stringify(reqData))
    .then((result: any) => {
      const resp = JSON.parse(result) as ResponseWrapper<string>
      if (resp.code === 200) {
        loadData()
        success('添加成功', 'top-right')
      } else {
        console.error('添加学生失败:', resp.message)
        toast.error('添加学生失败: ' + resp.message, 'top-right')
      }
    })
}

function updateStudent(data: StudentData) {
  let reqData = {
    Id: data.id,
    Name: data.name,
    Gender: data.gender,
    Hours: data.balance,
    Phone: data.phone,
    Teacher_Id: data.teacher_id,
    Remark: data.note,
  }
  console.log('Updating student with data:', reqData)
  Dispatch('student_manager:update_student', JSON.stringify(reqData))
    .then((result: any) => {
      const resp = JSON.parse(result) as ResponseWrapper<string>
      if (resp.code === 200) {
        loadData()
        success('更新成功', 'top-right')
      } else {
        console.error('更新学生失败:', resp.message)
        toast.error('更新学生失败: ' + resp.message, 'top-right')
      }
    })
}

function deleteStudent(item: StudentItem) {
  console.log('Deleting student with ID:', item.id)
  Dispatch('student_manager:delete_student', JSON.stringify({ Id: item.id }))
    .then((result: any) => {
      const resp = JSON.parse(result) as ResponseWrapper<string>
      if (resp.code === 200) {
        loadData()
        success('删除成功', 'top-right')
      } else {
        console.error('删除学生失败:', resp.message)
        toast.error('删除学生失败: ' + resp.message, 'top-right')
      }
    })
}

function exportStudents2Excel() {
  console.log('Exporting student data...')
  Dispatch('student_manager:export_students', '').then((result: any) => {
    console.log('Received export response:' + result)
    const resp = JSON.parse(result) as ResponseWrapper<string>
    if (resp.code === 200) {
      if (resp.data === 'cancel') {
        info('已取消导出', 'top-right')
        return
      }
      success('导出成功', 'top-right')
    } else {
      console.error('导出学生数据失败:', resp.message)
      toast.error('导出学生数据失败: ' + resp.message, 'top-right')
    }
  })
}

function createOrder(data: { studentId: number, amount: number, note: string }) {
  console.log('Processing recharge with data:', data)

  let reqData = {
    Student_Id: data.studentId,
    Hours: data.amount,
    Comment: data.note,
  };

  Dispatch('order_manager:create_order', JSON.stringify(reqData)).then((result: any) => {
    console.log('Received recharge response:' + result)
    const resp = JSON.parse(result) as ResponseWrapper<string>
    if (resp.code === 200) {
      const target = students.value.find(s => s.id === data.studentId)
      if (target) {
        target.balance += data.amount
        success('课时调整成功')
      }
      dialogRecharge.value = false
      loadData()
    } else {
      console.error('课时调整失败:', resp.message)
      toast.error('课时调整失败: ' + resp.message, 'top-right')
    }
  })
}


export function useStudentManage() {

  // --- 操作方法 ---
  const openAdd = () => {
    editedIndex.value = -1
    loadTeacherOptions()
    Object.assign(editedItem, defaultItem)
    dialog.value = true
  }

  const openEdit = (item: StudentItem) => {
    editedIndex.value = students.value.findIndex(s => s.id === item.id)
    loadTeacherOptions()
    Object.assign(editedItem, item)
    dialog.value = true
  }

  const openRecharge = (item: StudentItem) => {
    Object.assign(rechargeItem, item)
    dialogRecharge.value = true
  }

  const openDetails = (item: StudentItem) => {
    Object.assign(editedItem, item)
    dialogDetails.value = true
  }

  const openDelete = async (item: StudentItem) => {
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

  const saveRecharge = async (data: { studentId: number, amount: number, note: string }) => {
    createOrder(data)
  }

  const exportStudents = async () => {
    exportStudents2Excel()
  }

  return {
    search,
    page,
    itemsPerPage,
    totalItems,
    loading,
    dialog, dialogRecharge, dialogDetails,
    editedIndex, editedItem, rechargeItem,
    students, teacherOptions,
    headers,
    loadItems,
    getStatusColor, getStatusText,
    openAdd, openEdit, openRecharge, openDetails, openDelete,
    saveStudent, saveRecharge, exportStudents
  }
}