import { onMounted, ref } from 'vue'
import type { RecordItem } from '../../types/appModels'
import { categorizeOrderTags } from '../../utils/classification'
import { useToast } from '../../composables/useToast'
import type { OrderTag } from '../../types/appModels'
import { Dispatch } from '../../../wailsjs/go/main/App'
import type { ResponseWrapper } from '../../types/appModels'
import { GetOrdersByStudentIdResponse } from '../../types/response'

export const headers: any = [
  { title: '日期', key: 'date' },
  { title: '类型', key: 'type' },
  { title: '标签', key: 'tags' },
  { title: '变动', key: 'amount' },
  { title: '备注', key: 'remark' },
]

// studentId moved into the useDetailsDialog composable to avoid shared state across instances

const { success, error } = useToast()
const studentName = ref('')

// loadItems moved into the useDetailsDialog composable so it can reference the local studentId ref

function exportToExcel(studentId?: number) {
  let reqData = {
    Student_Id: studentId,
  }

  Dispatch("order_manager:export_orders_by_student_id", JSON.stringify(reqData)).then((result: any) => {
    const res = JSON.parse(result) as ResponseWrapper<string>
    if (res.code === 200) {
      success('导出成功，文件已保存至: ' + res.data, "top-right")
    } else {
      error('导出失败: ' + res.message, "top-right")
    }
  }).catch(() => {
    error('导出失败', "top-right")
  })
}


export function useDetailsDialog(emit: any) {


  // pagination state is now local to the composable to avoid shared state across instances
  const page = ref(1)
  const itemsPerPage = ref(10)
  const totalItems = ref(0) // 需要后端支持返回总数
  const loading = ref(false)

  const close = () => {
    emit('update:modelValue', false)
  }

  const updateModelValue = (val: boolean) => {
    emit('update:modelValue', val)
  }

  // keep studentId local to the composable to avoid shared state across instances
  const studentId = ref<number>(0)
  const records = ref<(RecordItem & { tags?: OrderTag[] })[]>([])

  // move loadOrders into composable so it can access pagination state
  function loadOrders(studentId: number) {
    let reqData = {
      Student_Id: studentId,
      Offset: (itemsPerPage.value * (page.value - 1)),
      Limit: itemsPerPage.value,
    }

    Dispatch("order_manager:get_orders_by_student_id", JSON.stringify(reqData)).then((result: any) => {
      const res = JSON.parse(result) as ResponseWrapper<GetOrdersByStudentIdResponse>
      const rawRecords = res.data.orders || []
      totalItems.value = res.data.total || 0
      records.value = rawRecords.map(item => ({
        id: item.id,
        // unix 毫秒时间戳转换为日期字符串
        date: new Date(item.created_at).toLocaleString(),
        type: item.type,
        amount: item.hours,
        // 如果备注为空，显示一个默认文本
        remark: item.comment || '-',
        tags: categorizeOrderTags(item.comment || '', item.hours)
      }))
    }).catch(() => {
      error('获取明细失败', "top-right")
    }).finally(() => {
      loading.value = false
    })
  }

  const onExport = async (studentId?: number) => {
    if (!studentId) {
      error('导出失败，学生ID无效')
      return
    }
    exportToExcel(studentId)
  }

  const getTypeLabel = (type: string) => {
    if (type === 'increase') return '充值'
    if (type === 'decrease') return '扣减'
    return '其他'
  }

  const getTypeColor = (type: string) => {
    if (type === 'increase') return 'success'
    if (type === 'decrease') return 'error'
    return 'default'
  }

  const load = async (studentId?: number, name?: string) => {
    studentName.value = name || ''
    records.value = [] // Reset

    if (studentId) {
      try {
        // TODO: API - 获取学生明细
        // const res = await window.go.main.App.GetStudentRecords(studentId)
        // const rawRecords = res || []

        // 模拟数据填充，实际应使用 rawRecords
        const rawRecords: RecordItem[] = []

        records.value = rawRecords.map(item => ({
          ...item,
          tags: categorizeOrderTags(item.remark || '', item.amount)
        }))
      } catch (e) {
        error('加载明细失败')
      }
    }
  }

  // loadItems now uses the composable-local studentId ref
  const loadItems = ({ page: newPage, itemsPerPage: newItemsPerPage, sortBy }: { page: number; itemsPerPage: number; sortBy?: string[] | string | undefined }): void => {
    console.log('Loading items with params:', { page: newPage, itemsPerPage: newItemsPerPage, sortBy })
    page.value = newPage
    itemsPerPage.value = newItemsPerPage
    if (studentId.value > 0) {
      loadOrders(studentId.value)
    } else {
      error('获取明细失败，无效的学生ID', "top-right")
    }
  }

  return {
    studentId,
    studentName,
    records,
    headers,
    page,
    itemsPerPage,
    totalItems,
    loading,
    close,
    updateModelValue,
    onExport,
    getTypeColor,
    getTypeLabel,
    load,
    loadItems,
  }
}