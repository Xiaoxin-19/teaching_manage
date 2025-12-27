import { ref } from 'vue'
import type { RecordItem } from '../../types/appModels'

// 表头配置（导出以便 SFC 引用）
export const headers: any = [
  { title: '日期', key: 'date', align: 'start', width: '180px' },
  { title: '类型', key: 'type', align: 'center', width: '100px' },
  { title: '变动', key: 'amount', align: 'end', width: '100px' },
  { title: '结余', key: 'balanceAfter', align: 'end', width: '100px' },
  { title: '备注', key: 'remark', align: 'start' },
]

export type FetchDetailsFn = (studentId: number) => Promise<{ studentName: string, records: RecordItem[] }>

// 组合函数：接收 emit 和可选 fetcher，返回组件需要的方法和响应式数据
export function useDetailsDialog(emit: (event: string, ...args: any[]) => void, fetcher?: FetchDetailsFn) {
  const studentName = ref('')
  const records = ref<RecordItem[]>([])

  const close = () => {
    emit('update:modelValue', false)
  }

  const updateModelValue = (val: boolean) => {
    emit('update:modelValue', val)
  }

  const onExport = () => {
    emit('export')
  }

  const getTypeColor = (type: string) => {
    switch (type) {
      case '充值': return 'success'
      case '赠送': return 'info'
      case '消课': return 'error'
      case '退费': return 'warning'
      default: return 'grey'
    }
  }

  // 加载数据：优先使用传入的 fetcher，否则使用 fallback
  const load = async (studentId?: number, fallbackName?: string, fallbackRecords?: RecordItem[]) => {
    if (!studentId) {
      studentName.value = fallbackName ?? ''
      records.value = fallbackRecords ?? []
      return
    }

    if (fetcher) {
      try {
        const res = await fetcher(studentId)
        studentName.value = res.studentName
        records.value = res.records ?? []
      } catch (e) {
        console.error('fetcher failed:', e)
        studentName.value = fallbackName ?? ''
        records.value = fallbackRecords ?? []
      }
    } else {
      // no fetcher provided, use fallback data
      studentName.value = fallbackName ?? ''
      records.value = fallbackRecords ?? []
    }
  }

  return {
    studentName,
    records,
    close,
    updateModelValue,
    onExport,
    getTypeColor,
    load,
  }
}
