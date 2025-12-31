<template>
  <v-dialog :model-value="modelValue" @update:model-value="close" max-width="800px" scrollable>
    <v-card class="rounded-lg elevation-4">
      <v-card-title class="d-flex justify-space-between align-center py-3 px-4"
        style="background-color: rgba(var(--v-theme-error), 0.1)">
        <div class="d-flex align-center text-subtitle-1 font-weight-bold text-error">
          <v-icon icon="mdi-alert-circle" color="error" class="mr-2"></v-icon>
          导入失败详情
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" @click="close"></v-btn>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text class="pa-0" style="max-height: 65vh;">
        <div class="pa-4 text-body-2 text-medium-emphasis">
          以下记录存在问题导致导入失败，请根据提示修正 Excel 文件后重新导入。
        </div>
        <v-table density="compact" hover>
          <thead>
            <tr>
              <th class="text-left font-weight-bold" style="width: 100px">行号</th>
              <th class="text-left font-weight-bold">错误原因</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in errorList" :key="index">
              <td class="font-weight-medium">第 {{ item.row }} 行</td>
              <td class="text-error text-body-2 py-2">
                <div v-for="(msg, i) in item.messages" :key="i" class="d-flex align-center">
                  <v-icon size="x-small" color="error" class="mr-1">mdi-circle-small</v-icon>
                  {{ msg }}
                </div>
              </td>
            </tr>
            <tr v-if="errorList.length === 0">
              <td colspan="2" class="text-center text-medium-emphasis py-4">
                无详细错误信息
              </td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions class="pa-3" style="background-color: rgba(var(--v-theme-on-surface), 0.04)">
        <v-spacer></v-spacer>
        <v-btn color="primary" variant="flat" @click="close" min-width="100">关闭</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  modelValue: boolean;
  errors: string[][];
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
}>();

const close = () => {
  emit('update:modelValue', false);
};

const errorList = computed(() => {
  if (!props.errors || !Array.isArray(props.errors)) return [];
  const list = [];
  for (let i = 0; i < props.errors.length; i++) {
    const rowErrors = props.errors[i];
    if (rowErrors && rowErrors.length > 0) {
      // 这里的 rowErrors 应该是一个字符串数组，但为了保险起见，做一下处理
      // 后端返回的是 [][]string
      // 假设后端逻辑：errInfo[i] 对应的是 rows[i+1] (因为 rows[0] 是表头)
      // 所以索引 i 对应的是 Excel 的第 i + 2 行

      // 有时候后端返回的错误信息可能已经包含了"第 x 行"的前缀，这里我们只提取纯粹的错误信息如果可能的话
      // 不过看后端代码：errInfo[i] = append(errInfo[i], fmt.Sprintf("第 %d 行: ...", i+2))
      // 后端已经格式化了字符串。
      // 如果我们想在前端分列显示行号，可能需要解析字符串，或者直接显示。
      // 既然我们有索引 i，我们可以自己生成行号，但是后端返回的信息里已经有了。
      // 让我们看看后端代码：
      // errInfo[i] = append(errInfo[i], fmt.Sprintf("第 %d 行: ...", i+2))
      // 所以 rowErrors 是 ["第 2 行: 错误1", "第 2 行: 错误2"]

      // 为了界面美观，我们可以尝试去除 "第 x 行: " 前缀，因为我们在第一列显示了行号。
      const cleanMessages = rowErrors.map(msg => {
        return msg.replace(/^第 \d+ 行: /, '');
      });

      list.push({
        row: i + 2,
        messages: cleanMessages
      });
    }
  }
  return list;
});
</script>
