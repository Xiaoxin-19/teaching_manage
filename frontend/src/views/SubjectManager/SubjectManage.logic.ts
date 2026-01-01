import { ref, reactive, onMounted } from 'vue'
import { useToast } from '../../composables/useToast'
import { useConfirm } from '../../composables/useConfirm'
import { GetSubjectList, CreateSubject, UpdateSubject, DeleteSubject } from '../../api/subject'
import type { Subject } from '../../types/appModels'
import { an } from 'vue-router/dist/router-CWoNjPRp.mjs'
import { Update } from 'vite/types/hmrPayload'
import { CreateSubjectRequest, DeleteSubjectRequest, GetSubjectListRequest, UpdateSubjectRequest } from '../../types/request'

export function useSubjectManage() {
  const { success, error, warning } = useToast()
  const confirm = useConfirm()

  // --- 状态 ---
  const loading = ref(false)
  const search = ref('')
  const subjects = ref<Subject[]>([])

  // 弹窗控制
  const dialogVisible = ref(false)
  const isEdit = ref(false)
  const formRef = ref<any>(null)

  // 分页控制
  const page = ref(1)
  const itemsPerPage = ref(10)
  const totalItems = ref(0)

  // 表单数据
  const formData = reactive({
    id: 0,
    name: '',
    lastModified: '' // 用于在编辑弹窗底部显示
  })

  // --- 数字格式化 (核心优化点) ---
  const formatCount = (count: number) => {
    // 格式化为 4 位数，不足补 0 (例如: 45 -> 0045)
    return (count || 0).toString().padStart(4, '0')
  }

  // --- 数据加载 ---
  const loadData = async () => {
    loading.value = true
    try {
      let reqData: GetSubjectListRequest = { Keyword: search.value, Offset: (page.value - 1) * itemsPerPage.value, Limit: itemsPerPage.value }
      // 调用真实 API
      const data = await GetSubjectList(reqData)

      // 简单的数据转换，确保字段匹配
      subjects.value = data.map(item => ({
        ...item,
        // 如果后端返回的是 created_at 时间戳，这里可以预处理成可读字符串
        // 假设后端 Subject 结构中有 updated_at 字段
        // lastModified: item.updated_at ? new Date(item.updated_at).toLocaleString() : '' 
      }))
    } catch (e: any) {
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  // --- 交互操作 ---

  const openAdd = () => {
    isEdit.value = false
    formData.id = 0
    formData.name = ''
    formData.lastModified = ''
    dialogVisible.value = true
  }

  const openEdit = (item: Subject) => {
    isEdit.value = true
    formData.id = item.id
    formData.name = item.name
    // 假设 item 中有 created_at 作为时间戳演示，实际应使用 updated_at
    formData.lastModified = item.created_at ? new Date(item.created_at).toLocaleString() : '未知'
    dialogVisible.value = true
  }

  const handleSave = async () => {
    if (!formData.name.trim()) {
      warning('科目名称不能为空')
      return
    }

    loading.value = true
    try {
      if (isEdit.value) {
        let reqData: UpdateSubjectRequest = { ID: formData.id, Name: formData.name }
        await UpdateSubject(reqData)
        success('科目更新成功')
      } else {
        let reqData: CreateSubjectRequest = { Name: formData.name }
        await CreateSubject(reqData)
        success('科目添加成功')
      }
      dialogVisible.value = false
      loadData() // 刷新列表
    } catch (e: any) {
      error(e.message)
    } finally {
      loading.value = false
    }
  }

  const handleDelete = async (item: Subject) => {
    // 1. 业务校验：如果有关联学员，禁止删除 (前端再次拦截)
    if (item.student_count > 0) {
      warning(`无法删除：该科目下仍有 ${item.student_count} 名在读学员。`)
      return
    }

    // 2. 二次确认
    const ok = await confirm.confirmDelete(item.name)
    if (ok) {
      try {
        let reqData: DeleteSubjectRequest = { ID: item.id }
        await DeleteSubject(reqData)
        success(`科目 ${item.name} 已删除`)
        loadData()
      } catch (e: any) {
        error(e.message)
      }
    }
  }

  onMounted(() => {
    loadData()
  })

  // 分页加载数据
  function loadItems({ page: newPage, itemsPerPage: newItemsPerPage, sortBy }: { page: number; itemsPerPage: number; sortBy?: string[] | string | undefined }): void {
    loadData()
  }

  // 表头定义：统一居中，使用百分比宽度
  const headers: any = [
    { title: '科目编号', key: 'subject_number', align: 'center', width: '20%', sortable: false },
    { title: '科目名称', key: 'name', align: 'center', width: '30%', sortable: false },
    { title: '在读学员', key: 'student_count', align: 'center', width: '30%', sortable: false },
    { title: '操作', key: 'actions', align: 'center', sortable: false, width: '20%' },
  ]

  return {
    loading,
    search,
    subjects,
    headers,
    dialogVisible,
    isEdit,
    formData,
    formRef,
    page,
    itemsPerPage,
    totalItems,
    loadData,
    openAdd,
    openEdit,
    handleSave,
    handleDelete,
    formatCount,
    loadItems
  }
}