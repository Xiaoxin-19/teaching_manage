<script setup lang="ts">
import { useRechargeDialog } from './RechargeDialog.logic'
import type { StudentData } from '../../types/appModels'

const props = defineProps<{
  modelValue: boolean
  student: StudentData
}>()

const emit = defineEmits(['update:modelValue', 'save'])

const {
  rechargeForm,
  rechargeType,
  rechargeColor,
  inferredTags,
  isValid,
  save,
  close
} = useRechargeDialog(props, emit)
</script>

<template>
  <v-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)" max-width="450px"
    transition="dialog-bottom-transition">
    <v-card class="rounded-lg elevation-4">
      <v-card-title class="d-flex justify-space-between align-center py-3 px-4">
        <div class="d-flex align-center">
          <v-icon color="primary" class="mr-3">mdi-wallet-plus</v-icon>
          <div>
            <div class="text-subtitle-1 font-weight-bold" style="line-height: 1.2;">课时调整</div>
            <div class="text-caption text-medium-emphasis mt-1 font-weight-normal">
              操作对象: {{ student.name }} (当前: {{ student.balance }})
            </div>
          </div>
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" density="comfortable" @click="close"></v-btn>
      </v-card-title>

      <v-divider></v-divider>

      <v-card-text class="pa-4 pt-6">
        <v-text-field :model-value="rechargeType" label="操作类型 (自动识别)" variant="outlined" density="comfortable"
          prepend-inner-icon="mdi-swap-horizontal" readonly class="mb-3">
          <template v-slot:append-inner>
            <v-chip :color="rechargeColor" size="x-small" label variant="flat" class="font-weight-bold">
              {{ rechargeType }}
            </v-chip>
          </template>
        </v-text-field>

        <v-text-field v-model.number="rechargeForm.amount" label="变动数量 (正数充值，负数扣减)" type="number" variant="outlined"
          density="comfortable" prepend-inner-icon="mdi-counter" suffix="课时" class="mb-3" hint="输入正数增加课时，输入负数减少课时"
          :rules="[v => v != 0 || '变动数量不能为0']" persistent-hint></v-text-field>

        <v-text-field v-model="rechargeForm.note" label="操作备注" placeholder="例如：微信续费 / 活动赠送" variant="outlined"
          density="comfortable" prepend-inner-icon="mdi-comment-outline"
          :hint="inferredTags.length > 0 ? '' : '输入关键词如“微信”、“赠送”可自动分类'" persistent-hint class="mb-1">
        </v-text-field>

        <v-expand-transition>
          <div v-if="inferredTags.length > 0" class="d-flex align-center mt-2 px-1">
            <v-icon size="small" color="primary" class="mr-2">mdi-auto-fix</v-icon>
            <div class="text-caption text-medium-emphasis mr-2">智能识别:</div>
            <div class="d-flex align-center gap-1"
              style="gap: 4px; overflow-x: auto; white-space: nowrap; -ms-overflow-style: none; scrollbar-width: none;">
              <v-chip v-for="(tag, index) in inferredTags" :key="index" :color="tag.color" size="small" label
                variant="tonal" class="font-weight-bold flex-shrink-0">
                {{ tag.label }}
              </v-chip>
            </div>
          </div>
        </v-expand-transition>
      </v-card-text>

      <v-divider></v-divider>

      <v-card-actions class="pa-3">
        <v-spacer></v-spacer>
        <v-btn variant="text" class="mr-2" @click="close">取消</v-btn>
        <v-btn :color="rechargeColor === 'error' ? 'error' : 'primary'" variant="elevated" @click="save" class="px-6"
          elevation="1" :disabled="!isValid">
          确认调整
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
