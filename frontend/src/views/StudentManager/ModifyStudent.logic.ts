import { ref, watch } from 'vue'
import type { StudentData } from '../../types/appModels'

export function useModifyStudent(props: { isEdit: boolean, initialData?: StudentData }, emit: any) {
  const formRef = ref<any>(null)

  const formData = ref<StudentData>({
    name: '',
    phone: '',
    gender: 'male',
    balance: 0,
    teacher_id: null,
    note: ''
  })

  // 监听弹窗打开
  watch(() => props.initialData, (val) => {
    if (val) {
      // 深度拷贝，防止直接修改父组件数据
      formData.value = JSON.parse(JSON.stringify(val))
    } else {
      // 重置
      formData.value = {
        name: '',
        phone: '',
        gender: 'male',
        balance: 0,
        teacher_id: null,
        note: ''
      }
    }
  }, { deep: true, immediate: true })

  const close = () => {
    emit('update:modelValue', false)
  }

  const save = async () => {
    if (!formRef.value) return
    const { valid } = await formRef.value.validate()
    if (valid) {
      emit('save', formData.value)
    }
  }

  return {
    formRef,
    formData,
    close,
    save
  }
}