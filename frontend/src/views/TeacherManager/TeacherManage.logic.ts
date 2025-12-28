import { onMounted, ref } from 'vue'
import type { ResponseWrapper, TeacherData } from '../../types/appModels'
import type { GetTeacherListResponse } from '../../types/response'

import { useToast } from '../../composables/useToast'
import { LogDebug, LogError, LogInfo } from '../../../wailsjs/runtime/runtime'
import { useConfirm } from '../../composables/useConfirm'
import { Dispatch } from '../../../wailsjs/go/main/App'


const toast = useToast()
const confirm = useConfirm()

// 分页状态
const page = ref(1)
const itemsPerPage = ref(10)
const totalItems = ref(0) // 需要后端支持返回总数
const loading = ref(false)

// 表头定义
const headers: any = [
  { title: '姓名', key: 'name', align: 'start', sortable: false },
  { title: '性别', key: 'gender', align: 'center', sortable: false },
  { title: '备注', key: 'remark', align: 'center', sortable: false },
  { title: '操作', key: 'actions', align: 'center', sortable: false, width: '120px' },
]

// 教师数据列表
const teachers = ref<TeacherData[]>([])
const search = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
var currentData = ref<TeacherData | undefined>(undefined)


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
function fetchTeachers() {
  loading.value = true
  const reqData = {
    key: search.value,
    offset: (page.value - 1) * itemsPerPage.value,
    limit: itemsPerPage.value,
  }

  console.log('Fetching teacher list with request:' + JSON.stringify(reqData))
  LogDebug('Fetching teacher list with request:' + JSON.stringify(reqData))
  Dispatch('teacher_manager:get_teacher_list', JSON.stringify(reqData))
    .then((result: any) => {
      console.log('Received teacher list response:' + result)
      const resp = JSON.parse(result) as ResponseWrapper<GetTeacherListResponse>
      if (resp.code === 200) {
        // 如果teacher为null则赋值为空数组
        if (!resp.data.teachers) {
          resp.data.teachers = []
        }


        teachers.value = resp.data.teachers.map((item) => ({
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


// 创建教师
function createTeacher(data: TeacherData) {
  // 创建匿名结构体以匹配后端要求
  var reqData = {
    name: data.name,
    phone: data.phone,
    remark: data.remark,
    gender: data.gender,
  }

  // json 转换
  var reqStr: string = JSON.stringify(reqData)
  Dispatch('teacher_manager:create_teacher', reqStr).then((resp) => {
    // 使用ResponseWrapper解析响应
    const response = JSON.parse(resp) as ResponseWrapper<null>
    if (response.code === 200) {
      toast.success('教师创建成功', 'top-right')
      fetchTeachers()
    } else {
      console.error('创建教师失败:', response.message)
      toast.error('教师创建失败: ' + response.message, 'top-right')
    }
  })
}


// 修改教师
function updateTeacher(data: TeacherData) {
  // 创建匿名结构体以匹配后端要求
  var reqData = {
    id: data.id,
    name: data.name,
    phone: data.phone,
    remark: data.remark,
    gender: data.gender,
  }
  // json 转换
  var reqStr: string = JSON.stringify(reqData)
  console.log('Updating teacher with data:' + reqStr)
  // 调用后端方法
  Dispatch('teacher_manager:update_teacher', reqStr).then((resp) => {
    // 使用ResponseWrapper解析响应
    const response = JSON.parse(resp) as ResponseWrapper<string>
    if (response.code === 200) {
      toast.success('教师更新成功', 'top-right')
      fetchTeachers()
    } else {
      console.error('更新教师失败:', response.message)
      toast.error('教师更新失败: ' + response.message, 'top-right')
    }
  })
}

// 删除教师
function deleteTeacher(item: TeacherData) {
  confirm.confirmDelete(item.name).then((confirmed) => {
    if (confirmed) {
      Dispatch('teacher_manager:delete_teacher', JSON.stringify({ id: item.id })).then((resp) => {
        const response = JSON.parse(resp) as ResponseWrapper<string>
        if (response.code === 200) {
          toast.success('教师删除成功', 'top-right')
          fetchTeachers()
        } else {
          console.error('删除教师失败:', response.message)
          toast.error('教师删除失败: ' + response.message, 'top-right')
        }
      })
    }
  })
}


// 导出教师数据
function exportTeacher2Excel() {
  Dispatch('teacher_manager:export_teacher_to_excel', '').then((resp) => {
    const response = JSON.parse(resp) as ResponseWrapper<string>
    if (response.code === 200) {
      toast.success('教师数据已导出到: ' + response.data, 'top-right')
    } else {
      console.error('导出教师数据失败:', response.message)
      toast.error('导出教师数据失败: ' + response.message, 'top-right')
    }
  })
}

function loadItems({ page: newPage, itemsPerPage: newItemsPerPage, sortBy }: { page: number; itemsPerPage: number; sortBy?: string[] | string | undefined }): void {
  loading.value = true
  page.value = newPage
  itemsPerPage.value = newItemsPerPage
  fetchTeachers()
  loading.value = false
}

const handleSave = (data: TeacherData) => {
  console.log('Save teacher data:', data)
  if (isEdit.value) {
    updateTeacher(data)
  } else {
    createTeacher(data)
  }
  // 操作成功后重新获取列表...
}
// 初始化

export function useTeacherManage() {
  const openAdd = () => {
    isEdit.value = false
    currentData.value = undefined
    dialogVisible.value = true
  }

  const openEdit = (item: TeacherData) => {
    isEdit.value = true
    currentData.value = item
    dialogVisible.value = true
  }

  const deleteItem = (item: TeacherData) => {
    deleteTeacher(item)
  }

  const exportData = () => {
    // TODO: 调用后端 API 导出 Excel
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
