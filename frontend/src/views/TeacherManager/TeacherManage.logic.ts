import { onMounted, ref } from 'vue'
import type { ResponseWrapper, Teacher } from '../../types/appModels'
import type { GetTeacherListResponse } from '../../types/response'

import { useToast } from '../../composables/useToast'
import { LogDebug, LogError, LogInfo } from '../../../wailsjs/runtime/runtime'
import { useConfirm } from '../../composables/useConfirm'
import { Dispatch } from '../../../wailsjs/go/main/App'
import { CreateTeacherRequest, DeleteTeacherRequest, GetTeacherListRequest, UpdateTeacherRequest } from '../../types/request'
import { CreateTeacher, DeleteTeacher, ExportTeachersToExcel, GetTeacherList, UpdateTeacher } from '../../api/teacher'


const toast = useToast()
const confirm = useConfirm()

// 分页状态
const page = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0) // 需要后端支持返回总数
const loading = ref(false)

// 表头定义
const headers: any = [
  { title: '编号', key: 'teacher_number', align: 'center', sortable: false },
  { title: '姓名', key: 'name', align: 'start', sortable: false },
  { title: '性别', key: 'gender', align: 'center', sortable: false },
  { title: '电话', key: 'phone', align: 'center', sortable: false },
  { title: '备注', key: 'remark', align: 'center', sortable: false },
  { title: '操作', key: 'actions', align: 'center', sortable: false, width: '120px' },
]

// 教师数据列表
const teachers = ref<Teacher[]>([])
const search = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
var currentData = ref<Teacher | undefined>(undefined)


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


// 从后端中获取教师列表
async function fetchTeachers() {
  loading.value = true
  const reqData: GetTeacherListRequest = {
    Keyword: search.value,
    Offset: (page.value - 1) * itemsPerPage.value,
    Limit: itemsPerPage.value,
  }



  console.log('Fetching teacher list with request:' + JSON.stringify(reqData))

  try {
    let result = await GetTeacherList(reqData)
    teachers.value = result.teachers.map((item: Teacher) => ({
      id: item.id,
      name: item.name,
      phone: item.phone,
      gender: item.gender,
      remark: item.remark,
      teacher_number: item.teacher_number,
      updated_at: item.updated_at,
      created_at: item.created_at,
      lastModified: new Date(item.updated_at).toLocaleString(),
    }))

    totalItems.value = result.total
    console.log('Fetched teacher list successfully:', {
      teachers: teachers.value, total: totalItems.value
    })

  } catch (error) {
    console.error('获取教师列表失败:', error)
    toast.error('获取教师列表失败: ' + error, 'top-right')
    loading.value = false
    return
  } finally {
    loading.value = false
  }
}


// 创建教师
async function createTeacher(data: Teacher) {
  // 创建匿名结构体以匹配后端要求
  let reqData: CreateTeacherRequest = {
    Name: data.name,
    Phone: data.phone,
    Remark: data.remark,
    Gender: data.gender,
  }

  try {
    await CreateTeacher(reqData)
    toast.success('教师创建成功', 'top-right')
    fetchTeachers()
  } catch (error: any) {
    console.error('创建教师失败:', error.message)
    toast.error('教师创建失败: ' + error.message, 'top-right')

  }
}


// 修改教师
async function updateTeacher(data: Teacher) {
  // 创建匿名结构体以匹配后端要求
  var reqData: UpdateTeacherRequest = {
    ID: data.id,
    Name: data.name,
    Phone: data.phone,
    Remark: data.remark,
    Gender: data.gender,
  }

  try {
    await UpdateTeacher(reqData)
    toast.success('教师更新成功', 'top-right')
    fetchTeachers()
  }
  catch (error: any) {
    console.error('更新教师失败:', error.message)
    toast.error('教师更新失败: ' + error.message, 'top-right')
  }
}

// 删除教师
function deleteTeacher(item: Teacher) {
  confirm.confirmDelete(item.name).then(async (confirmed) => {
    if (confirmed) {
      let reqData: DeleteTeacherRequest = { ID: item.id }
      try {
        await DeleteTeacher(reqData)
        toast.success('教师删除成功', 'top-right')
        fetchTeachers()
      }
      catch (error: any) {
        console.error('删除教师失败:', error.message)
        toast.error('教师删除失败: ' + error.message, 'top-right')
      }
    }
  })
}


// 导出教师数据
function exportTeacher2Excel() {
  ExportTeachersToExcel().then((filePath) => {
    if (filePath === 'cancel') {
      toast.info('已取消导出', 'top-right')
      return
    }
    toast.success('教师数据已导出到: ' + filePath, 'top-right')
  }).catch((error: any) => {
    console.error('导出教师数据失败:', error.message)
    toast.error('导出教师数据失败: ' + error.message, 'top-right')
  })
}

function loadItems({ page: newPage, itemsPerPage: newItemsPerPage, sortBy }: { page: number; itemsPerPage: number; sortBy?: string[] | string | undefined }): void {
  loading.value = true
  page.value = newPage
  itemsPerPage.value = newItemsPerPage
  fetchTeachers()
  loading.value = false
}

const handleSave = (data: Teacher) => {
  console.log('Save teacher data:', data)
  if (isEdit.value) {
    updateTeacher(data)
  } else {
    createTeacher(data)
  }
}
// 初始化

export function useTeacherManage() {
  const openAdd = () => {
    isEdit.value = false
    currentData.value = undefined
    dialogVisible.value = true
  }

  const openEdit = (item: Teacher) => {
    isEdit.value = true
    currentData.value = item
    dialogVisible.value = true
  }

  const deleteItem = (item: Teacher) => {
    deleteTeacher(item)
  }

  const exportData = () => {
    console.log('Trigger Export Excel')
    exportTeacher2Excel()
  }

  return {
    search,
    page,
    itemsPerPage,
    totalItems,
    loading,
    dialogVisible,
    isEdit,
    currentData,
    headers,
    teachers,
    fetchTeachers,
    loadItems,
    openAdd,
    openEdit,
    deleteItem,
    exportData,
    handleSave,
    getGenderColor,
    getGenderLabel,
  }
}
