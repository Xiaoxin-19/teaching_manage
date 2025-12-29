import { reactive, computed, watch } from 'vue'
import { categorizeOrderTags } from '../../utils/classification'
import type { OrderTag, StudentData } from '../../types/appModels'

export function useRechargeDialog(props: { modelValue: boolean, student: StudentData }, emit: any) {

  const rechargeForm = reactive({
    amount: 10,
    note: ''
  })

  // 监听弹窗打开或学生变化，重置表单
  watch(() => props.modelValue, (val) => {
    if (val) {
      rechargeForm.amount = 10
      rechargeForm.note = ''
    }
  })

  // --- 充值逻辑计算 ---
  const rechargeType = computed(() => {
    const amt = Number(rechargeForm.amount)
    if (amt > 0) return '充值/赠送'
    if (amt < 0) return '退费/扣减'
    return '无变动'
  })

  const rechargeColor = computed(() => {
    const amt = Number(rechargeForm.amount)
    if (amt > 0) return 'success'
    if (amt < 0) return 'error'
    return 'grey'
  })

  const inferredTags = computed<OrderTag[]>(() => {
    return categorizeOrderTags(rechargeForm.note, Number(rechargeForm.amount));
  });

  const isValid = computed(() => {
    return Number(rechargeForm.amount) !== 0
  })

  const close = () => {
    emit('update:modelValue', false)
  }

  const save = () => {
    const amt = Number(rechargeForm.amount)
    if (amt === 0) return
    emit('save', {
      studentId: props.student.id,
      amount: amt,
      note: rechargeForm.note
    })
    close()
  }

  return {
    rechargeForm,
    rechargeType,
    rechargeColor,
    inferredTags,
    isValid,
    save,
    close
  }
}
