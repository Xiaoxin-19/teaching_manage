package main

import (
	"context"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	"teaching_manage/backend/service"
)

// App struct
type App struct {
	ctx        context.Context
	dispatcher *dispatcher.Dispatcher
	setting    *service.SettingService
	backSvc    *service.BackupManager
}

// NewApp creates a new App application struct
func NewApp(dis *dispatcher.Dispatcher, setting *service.SettingService, backSvc *service.BackupManager) *App {
	return &App{
		dispatcher: dis,
		setting:    setting,
		backSvc:    backSvc,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Dispatch dispatches a method with payload to the registered handlers
func (a *App) Dispatch(router string, payload string) string {
	resp, err := a.dispatcher.Dispatch(a.ctx, router, []byte(payload))
	if err != nil {
		// Log the error
	}
	return resp
}

// OnShutdown 应用关闭时尝试备份数据
func (a *App) OnShutdown() {
	// 检查是否开启了自动备份
	autoBackupEnabled, err := a.setting.GetAutoBackupEnabled()
	if err != nil {
		logger.Error("读取自动备份设置失败", logger.ErrorType(err))
		return
	}

	if !autoBackupEnabled {
		logger.Info("自动备份已关闭，跳过备份")
		return
	}

	webDavCfg, err := a.setting.GetWebDavConfig()
	if err != nil {
		logger.Error("读取系统配置失败，无法进行 WebDAV 备份", logger.ErrorType(err))
		return
	}

	localPath, err := a.setting.GetBackupLocalPath()
	if err != nil {
		logger.Error("读取本地备份路径失败，无法进行本地备份", logger.ErrorType(err))
		return
	}

	// 检测是否设置了 WebDAV 备份配置，如果设置了则进行 WebDAV 备份
	if webDavCfg.WebDavBaseDir != "" && webDavCfg.WebDavURL != "" && webDavCfg.WebDavUserName != "" && webDavCfg.WebDavPassword != "" {
		_, err := a.backSvc.CreateBackupWebDav(a.ctx)
		if err != nil {
			logger.Error("WebDAV 备份失败", logger.ErrorType(err))
		} else {
			logger.Info("WebDAV 备份成功")
		}
	}

	// 检测是否设置了本地备份路径，如果设置了则进行本地备份
	if localPath != "" {
		_, err := a.backSvc.CreateBackupLocal(a.ctx)
		if err != nil {
			logger.Error("本地备份失败", logger.ErrorType(err))
		} else {
			logger.Info("本地备份成功")
		}
	}
}
