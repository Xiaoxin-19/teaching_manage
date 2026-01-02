import { ref, reactive, computed, watch, toRef } from 'vue'
import { useToast } from '../../composables/useToast'
import { categorizeOrderTags } from '../../utils/classification'

export function useRechargeDialog(props: any, emit: any) {
  const { success, error } = useToast()

  const loading = ref(false)

  // 双向绑定弹窗显示状态
  const dialogVisible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })

  const currentCourse = toRef(props, 'course')

  // 内部状态
  const mode = ref<'charge' | 'refund'>('charge')
  const form = reactive({
    hours: 10,
    amount: null as number | null,
    remark: ''
  })

  // 监听打开事件，初始化数据
  watch(() => props.modelValue, (val) => {
    if (val) {
      mode.value = props.initialMode || 'charge'
      form.hours = mode.value === 'charge' ? 10 : 1
      form.amount = null
      form.remark = ''
    }
  })

  // 校验逻辑
  const isValid = computed(() => {
    const hoursVal = Number(form.hours)
    const amountVal = form.amount === null ? 0 : Number(form.amount)

    const isHoursValid = !isNaN(hoursVal) && hoursVal > 0
    const isAmountValid = !isNaN(amountVal) && amountVal >= 0

    return isHoursValid && isAmountValid
  })

  // 智能标签配置
  const smartTags = computed(() => {
    return mode.value === 'charge' ? [10, 20, 30, 50, 100] : [1, 2, 4, 8]
  })

  // 实时智能分类
  const inferredTags = computed(() => {
    const val = Number(form.hours) || 0
    const signedAmount = mode.value === 'charge' ? val : -val
    return categorizeOrderTags(form.remark, signedAmount)
  })

  const applyTag = (val: number) => {
    form.hours = val
  }

  const getBalanceColor = (b: number | undefined) => {
    const val = b || 0
    return val < 0 ? 'error' : (val < 5 ? 'warning' : 'success')
  }

  // 提交处理
  const handleSubmit = async () => {
    if (!isValid.value) return

    loading.value = true
    try {
      // 计算实际变动值 (充值为正，退费为负)
      const change = mode.value === 'charge' ? Number(form.hours) : -Number(form.hours)
      const amount = form.amount ? Number(form.amount) : 0

      // 触发父组件的提交事件，将数据传出去
      // 注意：这里不直接调用 API，而是让父组件决定如何处理（如调用 API 后刷新列表）
      emit('submit', {
        courseId: currentCourse.value.id,
        hours: change, // 变动课时
        amount: amount,     // 实际金额
        remark: form.remark
      })

      // 这里我们假设父组件处理成功后会关闭弹窗，或者我们可以先关闭
      // 为了 UI 体验，我们在组件内不等待父组件 API 结束就先关闭吗？
      // 通常最好是父组件处理完后通知关闭，或者抛出 Promise。
      // 这里简化处理：发射事件后由父组件控制 loading 结束或弹窗关闭

    } catch (e: any) {
      error(e.message)
      loading.value = false
    }
  }

  // 提供给父组件调用，用于停止 loading 状态并关闭弹窗
  const onApiSuccess = () => {
    loading.value = false
    dialogVisible.value = false
    const action = mode.value === 'charge' ? '充值' : '扣除'
    success(`${action}成功：${form.hours} 节`)
  }

  const onApiError = () => {
    loading.value = false
  }

  return {
    loading,
    dialogVisible,
    mode,
    form,
    isValid,
    smartTags,
    inferredTags,
    applyTag,
    handleSubmit,
    getBalanceColor,
    onApiSuccess,
    onApiError
  }
}