import { ref, onMounted, onUnmounted } from 'vue'
import { OnFileDrop, OnFileDropOff } from '../../../wailsjs/runtime/runtime'
import { useToast } from '../../composables/useToast'
import { Dispatch } from '../../../wailsjs/go/main/App'
import { ResponseWrapper } from '../../types/appModels'

export function useImportRecord(props: { modelValue: boolean }, emit: any) {
  const { info, success, error } = useToast()
  const selectedFile = ref<string>('')

  const close = () => {
    emit('update:modelValue', false)
    // 关闭弹窗时清除选择的文件
    selectedFile.value = ''
  }

  const triggerFileInput = () => {
    // TODO: 调用系统文件选择对话框
    info('请直接拖拽文件到此处')
  }

  const downloadTemplate = () => {
    Dispatch('record_manager:download_import_template', "").then((result: any) => {
      const resp = JSON.parse(result) as ResponseWrapper<string>;
      console.log('下载模板响应:', resp);
      if (resp.code === 200) {
        if (resp.message.includes('cancel') || resp.data.includes('cancel')) {
          info('已取消操作');
          return;
        }
        success('模板下载成功，文件路径: ' + resp.data, 'top-right');
      } else {

        console.error('下载模板失败:', resp.message);
        error('下载模板失败: ' + resp.message, 'top-right');
      }
    })
  }

  const startImport = () => {
    if (!selectedFile.value) {
      error('请先选择或拖拽文件')
      return
    }
    // TODO: 这里需要对接真实的后端导入接口
    success(`开始导入文件: ${selectedFile.value}`)
    // 模拟导入过程
    setTimeout(() => {
      success('导入成功 (模拟)')
      close()
    }, 1000)
  }

  // 注册文件拖拽监听
  onMounted(() => {
    // useDropTarget = true: 启用拖拽目标检测
    // Wails 会自动检测带有 --wails-drop-target 样式的元素
    // 当拖拽经过这些元素时，会自动添加 wails-drop-target-active 类
    OnFileDrop((x: number, y: number, paths: string[]) => {
      // 只有当弹窗打开时才处理拖拽
      if (props.modelValue && paths.length > 0) {
        if (paths.length > 1) {
          info('一次仅支持导入一个文件，已使用第一个文件进行处理')
        }
        const filePath = paths[0]
        if (filePath.endsWith('.xlsx') || filePath.endsWith('.xls')) {
          selectedFile.value = filePath
          info(`已选择文件: ${filePath}`)
        } else {
          error('仅支持 Excel 文件 (.xlsx, .xls)')
        }
      }
    }, true)
  })

  // 组件卸载时移除监听
  onUnmounted(() => {
    OnFileDropOff()
  })

  return {
    selectedFile,
    close,
    triggerFileInput,
    downloadTemplate,
    startImport
  }
}
