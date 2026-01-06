<template>
  <v-sheet class="fill-height pa-6 bg-background overflow-y-auto">
    <div class="d-flex flex-column h-100">

      <!-- 本地备份路径设置 -->
      <v-card elevation="1" border class="rounded-lg mb-4">
        <v-card-item class="py-3 px-6 cursor-pointer" @click="openLocalBackupDirSetting" v-ripple>
          <div class="d-flex align-center justify-space-between">
            <div class="d-flex align-center">
              <v-icon icon="mdi-folder-cog" color="primary" class="mr-3"></v-icon>
              <div>
                <div class="text-body-2 font-weight-bold">本地备份路径</div>
                <div class="text-caption text-medium-emphasis">
                  {{ localBackupDir || '点击设置本地备份保存位置' }}
                </div>
              </div>
            </div>
            <v-icon>mdi-chevron-right</v-icon>
          </div>
        </v-card-item>
      </v-card>

      <!-- 自动备份设置 -->
      <v-card elevation="1" border class="rounded-lg mb-4">
        <v-card-item class="py-3 px-6">
          <div class="d-flex align-center justify-space-between">
            <div class="d-flex align-center">
              <v-icon icon="mdi-autorenew" color="primary" class="mr-3"></v-icon>
              <div>
                <div class="text-body-2 font-weight-bold">自动备份</div>
                <div class="text-caption text-medium-emphasis">
                  退出程序时自动创建备份
                </div>
              </div>
            </div>
            <v-switch v-model="autoBackupEnabled" :disabled="isAllowEnableAutoBackup" @change="updateAutoBackupEnabled"
              color="primary" density="compact" class="mr-1"></v-switch>
          </div>
        </v-card-item>
      </v-card>

      <!-- 模块：数据管理 -->
      <v-card elevation="2" border class="rounded-lg overflow-hidden">

        <!-- 1. 云服务连接状态 (顶部通栏) -->
        <v-card-item class="py-4 px-6 bg-surface cursor-pointer" @click="openConfig" v-ripple>
          <div class="d-flex align-center justify-space-between">
            <div class="d-flex align-center">
              <v-avatar :color="isConfigured ? 'success' : 'grey-lighten-2'" size="48" variant="tonal" class="mr-4">
                <v-icon :icon="isConfigured ? 'mdi-cloud-check' : 'mdi-cloud-off-outline'"
                  :color="isConfigured ? 'success' : 'grey'"></v-icon>
              </v-avatar>
              <div>
                <div class="text-subtitle-1 font-weight-bold">WebDAV 云同步</div>
                <div class="text-caption" :class="isConfigured ? 'text-success' : 'text-medium-emphasis'">
                  {{ isConfigured ? '已连接至 ' + config.url : '未配置云端存储，仅支持本地操作' }}
                </div>
              </div>
            </div>
            <v-btn icon variant="text" color="grey">
              <v-icon>mdi-cog-outline</v-icon>
            </v-btn>
          </div>
        </v-card-item>

        <v-divider></v-divider>

        <v-card-text class="pa-6">
          <div class="text-caption text-medium-emphasis mb-4">如果设置了备份方式，则每次关闭程序都会自动备份，且最多保留最近的7次备份</div>
          <v-row>
            <!-- 左侧：数据备份 -->
            <v-col cols="12" md="6">
              <v-sheet class="rounded-lg  pa-5 h-100 d-flex flex-column border">
                <div class="d-flex align-center mb-4">
                  <v-icon icon="mdi-backup-restore" color="primary" class="mr-2"></v-icon>
                  <span class="text-subtitle-2 font-weight-bold">创建备份</span>
                </div>

                <div class="flex-grow-1 d-flex flex-column justify-center align-center py-4">
                  <div class="text-h4 font-weight-bold mb-1">{{ lastBackupDate || '--' }}</div>
                  <div class="text-caption text-medium-emphasis mb-6">上次云端备份时间</div>
                </div>

                <!-- 操作按钮组 -->
                <div class="mt-auto">
                  <!-- 主按钮：根据配置变更为同步或导出 -->
                  <v-btn block color="primary" height="44" variant="flat" class="rounded-lg mb-2" :loading="backingUp"
                    @click="handleMainBackupAction"
                    :prepend-icon="isConfigured ? 'mdi-cloud-upload' : 'mdi-folder-download-outline'">
                    <template v-if="backingUp">
                      {{ backupStatusText || (isConfigured ? '正在同步到云端...' : '正在导出到本地...') }}
                    </template>
                    <template v-else>
                      {{ isConfigured ? '立即同步到云端' : '导出到本地' }}
                    </template>
                  </v-btn>

                  <!-- 次级按钮：当配置了云端时，提供仅导出本地的选项 -->
                  <v-btn v-if="isConfigured" block variant="text" size="small" color="grey-darken-1"
                    @click="exportLocalOnly" :disabled="backingUp">
                    <v-icon start size="small">mdi-folder-download</v-icon>
                    <span v-if="backingUp">{{ backupStatusText || '正在导出到本地...' }}</span>
                    <span v-else>仅导出到本地文件</span>
                  </v-btn>
                </div>
              </v-sheet>
            </v-col>

            <!-- 右侧：数据恢复 -->
            <v-col cols="12" md="6">
              <v-sheet class="rounded-lg border pa-0 h-100 d-flex flex-column overflow-hidden">
                <v-tabs border v-model="restoreTab" density="compact" color="primary" grow>
                  <v-tab value="cloud" :disabled="!isConfigured" class="text-caption">
                    <v-icon start size="small">mdi-cloud-download</v-icon> 云端恢复
                  </v-tab>
                  <v-tab value="local" class="text-caption">
                    <v-icon start size="small">mdi-laptop</v-icon> 本地导入
                  </v-tab>
                </v-tabs>

                <v-window v-model="restoreTab" class="flex-grow-1">
                  <!-- 云端面板 -->
                  <v-window-item value="cloud" class="pa-5 h-100">
                    <div v-if="!isConfigured"
                      class="d-flex flex-column align-center justify-center h-100 text-center py-4">
                      <v-icon icon="mdi-cloud-off-outline" size="40" color="grey-lighten-2" class="mb-2"></v-icon>
                      <div class="text-caption text-disabled">请先配置 WebDAV</div>
                    </div>
                    <div v-else class="d-flex flex-column h-100">
                      <div class="text-caption text-medium-emphasis mb-2">选择云端备份点:</div>
                      <v-select v-model="selectedBackup" :items="cloudBackups" item-title="name" item-value="path"
                        return-object variant="outlined" density="compact" placeholder="点击加载列表"
                        prepend-inner-icon="mdi-history" hide-details :loading="loadingBackups"
                        @click="fetchCloudBackups" class="mb-4">
                        <template v-slot:item="{ props, item }">
                          <v-list-item v-bind="props">
                            <template v-slot:subtitle>
                              {{ formatSize(item.raw.size) }} · {{ new Date(item.raw.mod_time *
                                1000).toLocaleString('zh-CN') }}
                            </template>
                          </v-list-item>
                        </template>
                        <template v-slot:selection="{ item }">
                          {{ item.raw.name }} ({{ formatSize(item.raw.size) }})
                        </template>
                      </v-select>

                      <v-spacer></v-spacer>
                      <!-- 绑定 btn-breathe 类 -->
                      <v-btn block variant="tonal" color="warning" :disabled="!selectedBackup" @click="confirmRestore"
                        :class="{ 'btn-breathe': selectedBackup }">
                        开始恢复
                      </v-btn>
                    </div>
                  </v-window-item>

                  <!-- 本地面板 -->
                  <v-window-item value="local" class="pa-5 h-100">
                    <div class="d-flex flex-column h-100">
                      <div
                        class="upload-zone d-flex flex-column align-center justify-center py-6 cursor-pointer mb-4 flex-grow-1"
                        style="--wails-drop-target: drop" @click="selectLocalFile">
                        <v-icon :icon="localFile ? 'mdi-file-check' : 'mdi-cloud-upload'"
                          :color="localFile ? 'success' : 'grey'" size="32" class="mb-2"></v-icon>
                        <div class="text-caption font-weight-bold"
                          :class="localFile ? 'text-high-emphasis' : 'text-medium-emphasis'">
                          {{ localFile ? '已选择文件' : '点击或拖拽 .db 文件至此' }}
                        </div>
                        <div class="text-caption text-medium-emphasis mt-1">
                          {{ localFile || '支持 .db 格式备份文件' }}
                        </div>
                      </div>
                      <v-spacer></v-spacer>
                      <!-- 绑定 btn-breathe 类 -->
                      <v-btn block variant="tonal" color="warning" :disabled="!localFile" @click="confirmRestore"
                        :class="{ 'btn-breathe': localFile }">
                        开始恢复
                      </v-btn>
                    </div>
                  </v-window-item>
                </v-window>
              </v-sheet>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>

      <!-- 配置弹窗 -->
      <v-dialog v-model="configDialog" max-width="480">
        <v-card class="rounded-lg elevation-4">
          <v-card-title class="text-subtitle-1 font-weight-bold px-6 pt-6">WebDAV 连接配置</v-card-title>
          <v-card-text class="px-6 py-4">
            <v-form ref="formRef" v-model="valid" @submit.prevent="saveConfig">
              <v-text-field v-model="config.url" label="服务器地址" placeholder="https://dav.jianguoyun.com/dav/"
                variant="outlined" density="comfortable" color="primary" class="mb-1"
                :rules="[rules.required, rules.url]">
                <template v-slot:append-inner>
                  <v-menu>
                    <template v-slot:activator="{ props }">
                      <v-btn v-bind="props" variant="text" size="small" color="primary" class="px-0">快速填充</v-btn>
                    </template>
                    <v-list density="compact">
                      <v-list-item v-for="(p, i) in providers" :key="i" @click="config.url = p.value"
                        :title="p.name"></v-list-item>
                    </v-list>
                  </v-menu>
                </template>
              </v-text-field>

              <v-text-field v-model="config.username" label="账号" variant="outlined" density="comfortable"
                color="primary" class="mb-1" :rules="[rules.required]"></v-text-field>

              <v-text-field v-model="config.password" label="应用密码/密码" type="password" variant="outlined"
                density="comfortable" color="primary" :rules="[rules.required]"></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions class="px-6 pb-6 pt-0">
            <v-btn variant="text" color="grey-darken-1" @click="configDialog = false">取消</v-btn>
            <v-spacer></v-spacer>
            <v-btn variant="outlined" color="primary" class="mr-2" :loading="testing" @click="testConnection"
              :disabled="!valid">测试连接</v-btn>
            <v-btn color="primary" variant="flat" :loading="saving" @click="saveConfig" :disabled="!valid">保存配置</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <!-- 本地备份路径设置弹窗 -->
      <v-dialog v-model="localBackupDirDialog" max-width="480">
        <v-card class="rounded-lg elevation-4">
          <v-card-title class="text-subtitle-1 font-weight-bold px-6 pt-6">本地备份路径设置</v-card-title>
          <v-card-text class="px-6 py-4">
            <div class="upload-zone d-flex flex-column align-center justify-center rounded-lg py-8 mb-3 cursor-pointer"
              style="--wails-drop-target: drop" @click="selectBackupDir">
              <v-icon :icon="localBackupDir ? 'mdi-folder-check' : 'mdi-folder-open'"
                :color="localBackupDir ? 'success' : 'primary'" size="48" class="mb-2">
              </v-icon>
              <div class="text-body-2 font-weight-bold mb-1">
                {{ localBackupDir ? '已选择路径' : '点击选择或拖拽文件夹至此' }}
              </div>
              <div class="text-caption text-medium-emphasis">
                {{ localBackupDir || '用于保存本地备份文件' }}
              </div>
            </div>
            <v-text-field v-model="localBackupDir" label="备份保存路径" placeholder="例如：D:/Backups/TeachingManage"
              variant="outlined" density="compact" color="primary" hide-details readonly>
            </v-text-field>
          </v-card-text>
          <v-card-actions class="px-6 pb-6 pt-3">
            <v-btn variant="text" color="grey-darken-1" @click="localBackupDirDialog = false">取消</v-btn>
            <v-spacer></v-spacer>
            <v-btn color="primary" variant="flat" :loading="savingLocalPath" @click="saveLocalBackupDir"
              :disabled="!localBackupDir">保存</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <!-- 恢复确认弹窗 -->
      <v-dialog v-model="restoreDialog" max-width="320">
        <v-card class="rounded-lg text-center pa-4 elevation-4">
          <v-icon icon="mdi-alert-circle" color="warning" size="48" class="mb-2 mx-auto"></v-icon>
          <div class="text-subtitle-1 font-weight-bold mb-2">确认恢复数据？</div>
          <div class="text-body-2 text-medium-emphasis mb-4">
            此操作将用备份文件覆盖当前所有数据，且<strong class="text-error">不可撤销</strong>。建议先进行一次备份。
          </div>
          <div class="d-flex gap-2 justify-center">
            <v-btn variant="text" @click="restoreDialog = false" class="flex-grow-1">取消</v-btn>
            <v-btn color="warning" variant="flat" @click="executeRestore" :loading="restoring"
              class="flex-grow-1">确认覆盖</v-btn>
          </div>
        </v-card>
      </v-dialog>

    </div>
  </v-sheet>
</template>

<script setup lang="ts">
import { useSettings } from './Settings.logic'

const {
  configDialog, restoreDialog, isConfigured, valid, formRef,
  lastBackupDate, backingUp, backupProgress, backupStatusText,
  restoreTab, restoring, loadingBackups, cloudBackups, selectedBackup, localFile, localFilePath,
  config, providers, rules, testing, saving,
  localBackupDir, localBackupDirDialog, savingLocalPath,
  autoBackupEnabled, isAllowEnableAutoBackup,
  openConfig, testConnection, saveConfig,
  handleMainBackupAction, exportLocalOnly,
  fetchCloudBackups, selectLocalFile, confirmRestore, executeRestore,
  openLocalBackupDirSetting, saveLocalBackupDir, selectBackupDir, updateAutoBackupEnabled, formatSize
} = useSettings()
</script>

<style scoped>
/* 虚线边框区域 */
.dashed-zone {
  border: 2px dashed rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 8px;
  transition: all 0.3s ease;
}

.dashed-zone:hover {
  border-color: rgb(var(--v-theme-primary));
  background-color: rgba(var(--v-theme-primary), 0.04);
}

/* 上传拖拽区域 */
.upload-zone {
  border: 2px dashed rgba(var(--v-border-color), var(--v-border-opacity));
  transition: all 0.3s ease;
  background-color: rgba(var(--v-theme-on-surface), 0.02);
}

.upload-zone:hover,
.upload-zone.wails-drop-target-active {
  border-color: rgb(var(--v-theme-primary));
  background-color: rgba(var(--v-theme-primary), 0.08);
  transform: scale(1.01);
}

/* 呼吸动画效果 */
@keyframes breathe {
  0% {
    box-shadow: 0 0 0 0 rgba(var(--v-theme-warning), 0.7);
    transform: scale(1);
  }

  70% {
    box-shadow: 0 0 0 8px rgba(var(--v-theme-warning), 0);
    transform: scale(1.02);
  }

  100% {
    box-shadow: 0 0 0 0 rgba(var(--v-theme-warning), 0);
    transform: scale(1);
  }
}

.btn-breathe {
  animation: breathe 2s infinite ease-in-out;
  z-index: 1;
}
</style>