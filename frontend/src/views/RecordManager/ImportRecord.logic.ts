import { ref, onMounted, onUnmounted } from 'vue'
import { OnFileDrop, OnFileDropOff } from '../../../wailsjs/runtime/runtime'
import { useToast } from '../../composables/useToast'
import { Dispatch } from '../../../wailsjs/go/main/App'
import { ResponseWrapper } from '../../types/appModels'
import { ImportExcelResponse, SelectFileResponse } from '../../types/response'
import { DownloadRecordImportTemplate, ImportFromExcel, SelectFilePath } from '../../api/record'
import { ImportFromExcelRequest } from '../../types/request'

export function useImportRecord(props: { modelValue: boolean }, emit: any) {
  const { info, success, error } = useToast()
  const selectedFile = ref<string>('')

  const close = () => {
    emit('update:modelValue', false)
    // 关闭弹窗时清除选择的文件
    selectedFile.value = ''
  }

  const importFile = async (filePath: string) => {
    const reqData: ImportFromExcelRequest = { filepath: filePath }
    const reqDataJson = JSON.stringify(reqData)

    try {
      console.log('导入请求数据:', reqDataJson);
      let result = await ImportFromExcel(reqData);
      success('导入成功，文件路径: ' + result.filepath, 'top-right');
      success(`导入总数: ${result.total_rows}, 错误数: ${result.error_infos.length}`, 'top-right');
      // 重新加载记录列表
      emit('import-success')
      close()
    } catch (e) {
      close();
      const err = e instanceof Error ? e.message : String(e);
      const info = e && (e as any).data ? (e as any).data.error_infos : null;
      if (info) {
        console.error('导入失败详情:', info);
        emit('import-failed', info)
      }
      error('导入失败: ' + err, 'top-right');
      return;
    }
  }

  const triggerFileInput = async () => {
    // 触发文件选择对话框
    try {
      console.log('触发文件选择对话框');
      let result = await SelectFilePath();
      if (result.includes('cancel')) {
        info('文件选择已取消', 'top-right');
        return;
      }
      selectedFile.value = result;
      info(`已选择文件: ${selectedFile.value}`, 'top-right');
    } catch (e) {
      error('文件选择失败: ' + e, 'top-right');
      console.error('Failed to select file path:', e);
    }
  }

  const downloadTemplate = async () => {
    try {
      let result = await DownloadRecordImportTemplate();
      if (result.includes('cancel')) {
        info('已取消操作');
        return;
      }
      success('模板下载成功，文件路径: ' + result, 'top-right');
    } catch (e) {
      if (e instanceof Error) {
        error('下载模板失败: ' + e.message, 'top-right');
        console.error('Failed to download import template:', e);
      }
    }
  }

  const startImport = () => {
    console.log('开始导入文件:', selectedFile.value)
    if (!selectedFile.value) {
      error('请先选择或拖拽文件')
      return
    }

    importFile(selectedFile.value)
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
