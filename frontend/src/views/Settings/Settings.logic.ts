import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useToast } from '../../composables/useToast'
import { useGlobalOverlay } from '../../composables/useGlobalOverlay'
import { OnFileDrop, OnFileDropOff } from '../../../wailsjs/runtime/runtime'
import { SetBackupLocalPath, GetBackupLocalPath, SetBackupLocalPathRequest, OpenFileDialog, SetWebDavConfig, GetWebDavConfig, TestWebDavConnection, ListWebDavBackups } from '../../api/backup'
import { OpenFileDialogResponse, SetWebDavConfigRequest } from '../../types/request'
import { WebDavBackupItem, WebDavConfigResponse } from '../../types/response'

export function useSettings() {
  // 使用项目中封装的 toast 钩子
  const { success, info, error } = useToast()
  const overlay = useGlobalOverlay()

  // --- UI 状态 ---
  const configDialog = ref(false)
  const restoreDialog = ref(false)
  const isConfigured = ref(false) // 实际项目中应从后端加载配置状态
  const valid = ref(false)
  const formRef = ref<any>(null)

  // --- 备份状态 ---
  const lastBackupDate = ref('') // 实际应从后端获取
  const backingUp = ref(false)
  const backupProgress = ref(0)
  const backupStatusText = ref('')

  // --- 恢复状态 ---
  const restoreTab = ref('local')
  const restoring = ref(false)
  const loadingBackups = ref(false)
  const cloudBackups = ref<string[]>([])
  const selectedBackup = ref<string | null>(null)
  const localFile = ref<string>('')
  const localFilePath = ref<string>('')

  // --- 配置表单 ---
  const config = reactive({
    url: '',
    username: '',
    password: ''
  })

  // --- 本地备份设置 ---
  const localBackupDir = ref('')
  const localBackupDirDialog = ref(false)
  const savingLocalPath = ref(false)

  // 预设的 WebDAV 提供商，方便用户快速填写
  const providers = [
    { name: '坚果云 (Jianguoyun)', value: 'https://dav.jianguoyun.com/dav/' },
    { name: 'Nextcloud', value: 'https://your-domain.com/remote.php/dav/files/user/' },
  ]

  // 表单校验规则
  const rules = {
    required: (v: string) => !!v || '此项必填',
    url: (v: string) => /^(https?:\/\/)/.test(v) || '地址需以 http 或 https 开头'
  }

  // --- 交互反馈状态 ---
  const testing = ref(false)
  const saving = ref(false)

  // --- 方法实现 ---
  const PASSWORD_PLACEHOLDER = "(●'◡'●)"

  const openConfig = async () => {
    try {
      const saved: WebDavConfigResponse = await GetWebDavConfig()
      config.url = saved.url || ''
      config.username = saved.username || ''
      config.password = saved.password
      isConfigured.value = !!saved.url
    } catch (err) {
      console.error('获取 WebDav 配置失败:', err)
    }
    configDialog.value = true
  }

  const openLocalBackupDirSetting = async () => {
    // 加载当前保存的路径
    try {
      const currentPath = await GetBackupLocalPath()
      if (currentPath) {
        localBackupDir.value = currentPath
      }
    } catch (err) {
      console.error('获取备份路径失败:', err)
    }
    localBackupDirDialog.value = true
  }

  const saveLocalBackupDir = async () => {
    if (!localBackupDir.value) return
    savingLocalPath.value = true

    try {
      let req: SetBackupLocalPathRequest = {
        local_path: localBackupDir.value
      }
      const result = await SetBackupLocalPath(req)
      localBackupDirDialog.value = false
      success('本地备份路径已保存: ' + result)
    } catch (err) {
      const errMsg = err instanceof Error ? err.message : String(err)
      error('保存失败: ' + errMsg)
      console.error('保存备份路径失败:', err)
    } finally {
      savingLocalPath.value = false
    }
  }

  // 测试 WebDav 连接
  const testConnection = async () => {
    if (!valid.value) return
    testing.value = true
    const req: SetWebDavConfigRequest = {
      url: config.url,
      username: config.username,
      password: config.password,
    }
    try {
      await TestWebDavConnection(req)
      success('WebDav 连接成功')
    } catch (err) {
      const errMsg = err instanceof Error ? err.message : String(err)
      error('连接失败: ' + errMsg)
    } finally {
      testing.value = false
    }
  }

  // 保存 WebDav 配置
  const saveConfig = async () => {
    if (!formRef.value) return
    const { valid: isValid } = await formRef.value.validate()
    if (!isValid) return

    saving.value = true
    const req: SetWebDavConfigRequest = {
      url: config.url,
      username: config.username,
      password: config.password,
    }
    try {
      await SetWebDavConfig(req)
      isConfigured.value = true
      configDialog.value = false
      success('WebDav 配置已保存')
      restoreTab.value = 'cloud'
      fetchCloudBackups()
    } catch (err) {
      const errMsg = err instanceof Error ? err.message : String(err)
      error('保存失败: ' + errMsg)
      console.error('保存 WebDav 配置失败:', err)
    } finally {
      saving.value = false
    }
  }

  // 云端备份逻辑 (模拟)
  const startCloudBackup = () => {
    overlay.show('正在备份到云端，请稍候...')
    backingUp.value = true
    backupProgress.value = 0
    backupStatusText.value = '正在加密并上传至云端...'

    const interval = setInterval(() => {
      if (backupProgress.value >= 100) {
        clearInterval(interval)
        backingUp.value = false;
        overlay.hide()
        // 获取当前时间字符串，兼容处理
        const now = new Date();
        lastBackupDate.value = `${now.getFullYear()}/${now.getMonth() + 1}/${now.getDate()} ${now.getHours()}:${now.getMinutes()}`;
        success('备份成功！已上传至云端')
      } else {
        backupProgress.value += Math.floor(Math.random() * 10)
        if (backupProgress.value > 100) backupProgress.value = 100
      }
    }, 300)
  }

  // 仅本地导出逻辑 (模拟)
  const exportLocalOnly = () => {
    overlay.show('正在导出本地备份，请稍候...')
    backingUp.value = true
    backupProgress.value = 0
    backupStatusText.value = '正在打包数据...'

    setTimeout(() => {
      backupProgress.value = 50
    }, 300)

    setTimeout(() => {
      backupProgress.value = 100
      backingUp.value = false
      overlay.hide()
      const now = new Date();
      lastBackupDate.value = `${now.getFullYear()}/${now.getMonth() + 1}/${now.getDate()} ${now.getHours()}:${now.getMinutes()}`;
      info('已导出到 Downloads 文件夹')
    }, 800)
  }

  // 统一的备份入口
  const handleMainBackupAction = () => {
    if (isConfigured.value) {
      startCloudBackup()
    } else {
      exportLocalOnly()
    }
  }

  const formatSize = (bytes: number) => {
    if (!bytes) return '0B'
    const units = ['B', 'KB', 'MB', 'GB']
    let i = 0
    let val = bytes
    while (val >= 1024 && i < units.length - 1) {
      val = val / 1024
      i++
    }
    return `${val.toFixed(1)}${units[i]}`
  }

  // 获取云端备份列表（实时）
  const fetchCloudBackups = async () => {
    loadingBackups.value = true
    try {
      const items: WebDavBackupItem[] = await ListWebDavBackups()
      let list = items || []
      cloudBackups.value = list.map(item => `${item.name} (${formatSize(item.size)})`)
    } catch (err) {
      const errMsg = err instanceof Error ? err.message : String(err)
      error('获取云端备份失败: ' + errMsg)
      cloudBackups.value = []
    } finally {
      loadingBackups.value = false
    }
  }

  const selectLocalFile = async () => {
    console.log(localBackupDir.value)
    let reqData: OpenFileDialogResponse = {
      title: '选择本地备份文件',
      default_path: localBackupDir.value || '',
      filters: [{ display_name: '数据库文件', pattern: '*.db' }]
    }

    try {
      let result = await OpenFileDialog(reqData)
      if (result.includes('cancel')) {
        info('已取消选择文件')
        return
      }
      localFilePath.value = result
      // 提取文件名
      const fileName = result.substring(result.lastIndexOf('\\') + 1)
      localFile.value = fileName
      info(`已选择备份文件: ${fileName}`)
    }
    catch (err) {
      const errMsg = err instanceof Error ? err.message : String(err)
      error('选择文件失败: ' + errMsg)
      console.error('打开文件对话框失败:', err)
    }
  }

  const selectBackupDir = async () => {
    let req: OpenFileDialogResponse = {
      title: '选择本地备份目录',
      default_path: localBackupDir.value || '',
      filters: [],
      is_path: true,
    }
    try {
      let result = await OpenFileDialog(req)
      if (result.includes('cancel')) {
        info('已取消选择目录')
        return
      }
      localBackupDir.value = result
      success(`已选择目录: ${result}`)
    } catch (err) {
      const errMsg = err instanceof Error ? err.message : String(err)
      error('选择目录失败: ' + errMsg)
      console.error('打开文件对话框失败:', err)
    }
  }

  const confirmRestore = () => {
    restoreDialog.value = true
  }

  const executeRestore = () => {
    overlay.show('正在恢复数据，请稍候...')
    restoring.value = true
    // 模拟恢复过程
    setTimeout(() => {
      restoring.value = false
      restoreDialog.value = false
      overlay.hide()
      success('数据恢复成功！系统将重新加载')
    }, 2000)
  }

  // 注册文件拖拽监听
  onMounted(async () => {
    // 初始化时加载本地备份路径
    try {
      const currentPath = await GetBackupLocalPath()
      if (currentPath) {
        localBackupDir.value = currentPath
      }
    } catch (err) {
      console.error('初始化备份路径失败:', err)
    }

    // 尝试读取已保存的 WebDav 配置以设置状态
    try {
      const saved = await GetWebDavConfig()
      console.log('已加载 WebDav 配置:', saved)
      if (saved.url) {
        config.url = saved.url
        config.username = saved.username
        config.password = saved.password
        isConfigured.value = true
      }
    } catch (err) {
      console.error('初始化 WebDav 配置失败:', err)
    }

    OnFileDrop((x: number, y: number, paths: string[]) => {
      if (paths.length === 0) return

      const filePath = paths[0]

      // 如果本地备份路径设置弹窗打开，处理文件夹拖拽
      if (localBackupDirDialog.value) {
        // 提取目录路径（如果是文件，则取其父目录）
        const dirPath = filePath.includes('.')
          ? filePath.substring(0, filePath.lastIndexOf('\\'))
          : filePath
        localBackupDir.value = dirPath
        info(`已设置备份路径: ${dirPath}`)
        return
      }

      // 如果是本地恢复标签页，处理 .db 文件拖拽
      if (restoreTab.value === 'local') {
        if (filePath.endsWith('.db')) {
          localFilePath.value = filePath
          // 提取文件名
          const fileName = filePath.substring(filePath.lastIndexOf('\\') + 1)
          localFile.value = fileName
          info(`已选择备份文件: ${fileName}`)
        } else {
          info('仅支持 .db 备份文件')
        }
      }
    }, true)
  })

  // 组件卸载时移除监听
  onUnmounted(() => {
    OnFileDropOff()
  })

  return {
    // State
    configDialog,
    restoreDialog,
    isConfigured,
    valid,
    formRef,
    lastBackupDate,
    backingUp,
    backupProgress,
    backupStatusText,
    restoreTab,
    restoring,
    loadingBackups,
    cloudBackups,
    selectedBackup,
    localFile,
    localFilePath,
    config,
    providers,
    rules,
    testing,
    saving,
    localBackupDir,
    localBackupDirDialog,
    savingLocalPath,

    // Methods
    openConfig,
    testConnection,
    saveConfig,
    handleMainBackupAction,
    exportLocalOnly,
    fetchCloudBackups,
    selectLocalFile,
    confirmRestore,
    executeRestore,
    openLocalBackupDirSetting,
    saveLocalBackupDir,
    selectBackupDir
  }
}