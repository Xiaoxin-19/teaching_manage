import { ref } from 'vue'
import type { RecordItem } from '../../types/appModels'
import { categorizeOrderTags } from '../../utils/classification'
import { useToast } from '../../composables/useToast'
import type { OrderTag } from '../../types/appModels'

export const headers: any = [
  { title: '日期', key: 'date', width: '150px' },
  { title: '类型', key: 'type', width: '100px' },
  { title: '标签', key: 'tags', width: '200px', sortable: false },
  { title: '变动', key: 'amount', align: 'end', width: '100px' },
  { title: '结余', key: 'balanceAfter', align: 'end', width: '100px' },
  { title: '备注', key: 'remark' },
]

export function useDetailsDialog(emit: any) {
  const { success, error } = useToast()
  const studentName = ref('')
  const records = ref<(RecordItem & { tags?: OrderTag[] })[]>([])

  const close = () => {
    emit('update:modelValue', false)
  }

  const updateModelValue = (val: boolean) => {
    emit('update:modelValue', val)
  }

  const onExport = async (studentId?: number) => {
    try {
      if (!studentId) return
      // TODO: API - 导出明细
      // await window.go.main.App.ExportStudentDetails(studentId)
      success('导出成功')
    } catch (e) {
      error('导出失败')
    }
  }

  const getTypeColor = (type: string) => {
    if (type === '充值') return 'success'
    if (type === '消课') return 'error'
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
          tags: categorizeOrderTags(item.remark || '')
        }))
      } catch (e) {
        error('加载明细失败')
      }
    }
  }

  return {
    studentName,
    records,
    close,
    updateModelValue,
    onExport,
    getTypeColor,
    load
  }
}