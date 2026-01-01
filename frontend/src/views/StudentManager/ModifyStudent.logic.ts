import { ref, watch, computed } from 'vue'
import type { StudentData } from '../../types/appModels'

export function useModifyStudent(props: { modelValue?: boolean, isEdit?: boolean, initialData?: StudentData }, emit: any) {
  const formRef = ref<any>(null)
  const isFormValid = ref(false)

  const defaultData = (): StudentData => ({
    name: '',
    phone: '',
    gender: 'male',
    status: 1,
    remark: ''
  })

  const formData = ref<StudentData>(defaultData())

  const resetFormFromInitial = () => {
    if (props.initialData) {
      console.log('Editing student data:', props.initialData)
      // 深度拷贝，防止直接修改父组件数据
      formData.value = JSON.parse(JSON.stringify(props.initialData))
    } else {
      console.log('Creating new student')
      // 重置
      formData.value = defaultData()
    }
    // 重置表单校验状态，防止旧错误残留
    if (formRef.value && typeof formRef.value.resetValidation === 'function') {
      formRef.value.resetValidation()
    }
  }

  // 监听弹窗打开（每次打开都重置表单，从 initialData 或默认值）
  watch(() => props.modelValue, (val) => {
    if (val) {
      resetFormFromInitial()
    }
  }, { immediate: true })

  // 如果 initialData 发生变化（例如外部更新），并且弹窗已打开，也同步更新表单
  watch(() => props.initialData, (val) => {
    if (props.modelValue && val) {
      resetFormFromInitial()
    }
  }, { deep: true })

  // 根据表单的校验状态计算是否禁用提交
  const isSubmitDisabled = computed(() => {
    return !isFormValid.value
  })


  const close = () => {
    emit('update:modelValue', false)
  }

  const save = async () => {
    if (!formRef.value) return false
    const { valid } = await formRef.value.validate()
    if (valid) {
      emit('save', formData.value)
      close()
      return true
    }
    return false
  }

  return {
    formRef,
    isFormValid,
    formData,
    close,
    save,
    isSubmitDisabled
  }

}