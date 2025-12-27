import { ref, watch, Ref } from 'vue'
import type { StudentData, TeacherOption } from '../../types/appModels'

export type ModifyProps = {
  modelValue: boolean
  isEdit?: boolean
  initialData?: StudentData
  teacherOptions?: TeacherOption[]
}

export type ModifyEmit = (event: 'update:modelValue' | 'save', ...args: any[]) => void

export function useModifyStudent(props: ModifyProps, emit: ModifyEmit) {
  const formRef: Ref<any> = ref(null)

  // 本地表单数据
  const formData = ref<StudentData>({
    name: '',
    gender: '男',
    phone: '',
    teacherId: null,
    remark: ''
  })

  // 监听弹窗开启，重置或填充数据
  watch(() => props.modelValue, (val) => {
    if (val) {
      if (props.isEdit && props.initialData) {
        // 深度拷贝以断开引用
        formData.value = JSON.parse(JSON.stringify(props.initialData))
      } else {
        // 重置为默认空状态
        formData.value = {
          name: '',
          gender: '男',
          phone: '',
          teacherId: null,
          remark: ''
        }
      }
    }
  })

  const close = () => {
    emit('update:modelValue', false)
  }

  const save = () => {
    // 可以在此处添加表单校验逻辑
    emit('save', formData.value)
    close()
  }

  return {
    formRef,
    formData,
    close,
    save,
  }
}
