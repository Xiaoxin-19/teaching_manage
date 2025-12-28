import { ref, reactive, computed, onMounted } from 'vue'
import type { ResponseWrapper, StudentData, StudentItem, TeacherOption } from '../../types/appModels'
import { useToast } from '../../composables/useToast'
import { useConfirm } from '../../composables/useConfirm'
import { categorizeOrderTags } from '../../utils/classification'
import type { OrderTag } from '../../types/appModels'
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
const rechargeForm = reactive({ amount: 10, note: '' })

// 列表数据 (初始化为空)
const students = ref<StudentItem[]>([])

// 教师列表 (用于搜索)
const teacherOptions = ref<TeacherOption[]>([])

// --- 初始化加载 ---
const loadData = async () => {
  await loadStudents()
}

function loadStudents() {
  loading.value = true
  try {
    // TODO: API - 调用后端获取学生列表
    const reqData = {
      Key: search.value,
      Offset: (page.value - 1) * itemsPerPage.value,
      Limit: itemsPerPage.value,
    }

    console.log('Request Data:', reqData)


    Dispatch('student_manager:get_student_list', JSON.stringify(reqData))
      .then((result: any) => {
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
            note: '', // 后端暂无 remark 字段
          }))
          totalItems.value = resp.data.total // 需要后端支持返回总数
        } else {
          console.error('获取学生列表失败:', resp.message)
          toast.error('获取学生列表失败: ' + resp.message, 'top-right')
        }

      })

    loading.value = false


    console.log('Fetching data...')
  } catch (e) {
    error('加载数据失败')
  }
}


function loadTeacherOptions() {
  const requestData = {
    key: "",
    offset: 0,
    limit: -1,
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
  { title: '姓名', key: 'name', align: 'center', sortable: false, width: '120px' },
  { title: '剩余课时', key: 'balance', sortable: false, width: '120px' },
  { title: '状态', key: 'status', sortable: false, width: '100px' },
  { title: '操作', key: 'actions', sortable: false, align: 'end', width: '180px' },
]

// --- 辅助函数 ---
const getStatusColor = (bal: number) => bal < 0 ? 'error' : (bal < 5 ? 'warning' : 'success')
const getStatusText = (bal: number) => bal < 0 ? '欠费' : (bal < 5 ? '余额不足' : '正常')
export function useStudentManage() {




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
        let reqData = {
          Name: data.name,
          Gender: data.gender,
          Hours: data.balance,
          Phone: data.phone,
          Teacher_Id: data.teacher_id,
          Note: data.note,
        }

        console.log('Creating student with data:', reqData)
        Dispatch('student_manager:create_student', JSON.stringify(reqData))
          .then((result: any) => {
            const resp = JSON.parse(result) as ResponseWrapper<string>
            if (resp.code === 200) {
              loadData()
              success('添加成功')
            } else {
              console.error('添加学生失败:', resp.message)
              toast.error('添加学生失败: ' + resp.message, 'top-right')
            }
          })
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
    page,
    itemsPerPage,
    totalItems,
    loading,
    dialog, dialogRecharge, dialogDetails,
    editedIndex, editedItem, rechargeItem, rechargeForm,
    students, teacherOptions,
    headers,
    rechargeType, rechargeColor, inferredTags,
    loadItems,
    getStatusColor, getStatusText,
    openAdd, openEdit, openRecharge, openDetails, openDelete,
    saveStudent, saveRecharge, exportStudents
  }
}