package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	"teaching_manage/backend/pkg/pathutil"
	requestx "teaching_manage/backend/service/request"
	responsex "teaching_manage/backend/service/response"
	"time"

	"github.com/studio-b12/gowebdav"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type BackupManager struct {
	settingSvc *SettingService
	Ctx        context.Context
}

func NewBackupManager(settingSvc *SettingService) *BackupManager {
	return &BackupManager{
		settingSvc: settingSvc,
	}
}

func (bm *BackupManager) SetBackupLocalPath(ctx context.Context, req *requestx.SetBackupLocalPathRequest) (string, error) {
	// 校验路径合法性
	if err := pathutil.ValidateBackupPath(req.LocalPath); err != nil {
		return "", fmt.Errorf("路径校验失败: %w", err)
	}

	// 规范化路径
	normalizedPath := pathutil.NormalizePath(req.LocalPath)

	// 调用 SettingService 保存本地备份路径
	if err := bm.settingSvc.UpdateLocalBackupPath(normalizedPath); err != nil {
		logger.Error("保存备份路径失败", logger.ErrorType(err))
		return "", fmt.Errorf("保存备份路径失败: %w", err)
	}

	return normalizedPath, nil
}

// GetBackupLocalPath 获取本地备份路径
func (bm *BackupManager) GetBackupLocalPath(ctx context.Context) (string, error) {
	path, err := bm.settingSvc.GetBackupLocalPath()
	if err != nil {
		return "", fmt.Errorf("获取系统设置失败: %w", err)
	}

	return path, nil
}

// resolveWebDavPassword 允许占位符或空值沿用已保存的密码
func (bm *BackupManager) resolveWebDavPassword(pwd string) (string, error) {
	if pwd != "" && pwd != DEFAULT_PASS_REPLACE {
		return pwd, nil
	}

	cfg, err := bm.settingSvc.GetWebDavConfig()
	if err != nil {
		return "", fmt.Errorf("读取已保存的 WebDav 配置失败: %w", err)
	}
	return cfg.WebDavPassword, nil
}

// 创建 WebDAV 客户端，设置统一超时
func (bm *BackupManager) newWebDavClient(url, username, password string) *gowebdav.Client {
	client := gowebdav.NewClient(url, username, password)
	client.SetTimeout(8 * time.Second)
	return client
}

func (bm *BackupManager) OpenPathSelectorDialog(ctx context.Context, req *requestx.OpenSelectPathRequest) (string, error) {
	logger.Debug("打开选择路径对话框",
		logger.String("title", req.Title),
		logger.String("default_path", req.DefaultPath),
		logger.Int("filters_count", len(req.Filters)),
		logger.String("is_path", fmt.Sprintf("%v", req.IsPath)),
	)
	var filters []wails.FileFilter
	for _, f := range req.Filters {
		filters = append(filters, wails.FileFilter{
			DisplayName: f.DisplayName,
			Pattern:     f.Pattern,
		})
	}
	var selectedPath string
	var err error
	if req.IsPath {
		// 使用 Wails 的运行时打开目录选择对话框
		selectedPath, err = wails.OpenDirectoryDialog(bm.Ctx, wails.OpenDialogOptions{
			Title:            req.Title,
			DefaultDirectory: req.DefaultPath,
			Filters:          filters,
		})
	} else {
		// 使用 Wails 的运行时打开文件选择对话框
		selectedPath, err = wails.OpenFileDialog(bm.Ctx, wails.OpenDialogOptions{
			Title:            req.Title,
			DefaultDirectory: req.DefaultPath,
			Filters:          filters,
		})
	}

	if err != nil {
		logger.Error("打开选择对话框失败", logger.ErrorType(err))
		return "", fmt.Errorf("打开选择对话框失败")
	}

	if selectedPath == "" {
		return "cancel", nil
	}
	return selectedPath, nil
}

// 测试 WebDav 连接（不落库，仅校验）
func (bm *BackupManager) TestWebDavConnection(ctx context.Context, req *requestx.SetWebDavConfigRequest) (string, error) {
	password, err := bm.resolveWebDavPassword(req.WebDavPassword)
	if err != nil {
		return "", err
	}

	client := bm.newWebDavClient(req.WebDavURL, req.WebDavUserName, password)

	if err := client.Connect(); err != nil {
		logger.Error("WebDav 连接失败", logger.ErrorType(err))
		return "", fmt.Errorf("WebDav 连接失败: %w", err)
	}

	if _, err := client.Stat("/"); err != nil {
		logger.Error("WebDav 认证失败", logger.ErrorType(err))
		return "", fmt.Errorf("WebDav 认证失败: %w", err)
	}

	return "WebDav 连接成功", nil
}

func (bm *BackupManager) SetWebDavConfig(ctx context.Context, req *requestx.SetWebDavConfigRequest) (string, error) {
	password, err := bm.resolveWebDavPassword(req.WebDavPassword)
	if err != nil {
		return "", err
	}

	cfg := entity.WebDavSetting{
		WebDavURL:      req.WebDavURL,
		WebDavUserName: req.WebDavUserName,
		WebDavPassword: password,
	}
	// 调用 SettingService 保存 WebDav 配置
	if err := bm.settingSvc.UpdateWebDavConfig(cfg); err != nil {
		logger.Error("保存 WebDav 配置失败", logger.ErrorType(err))
		return "", fmt.Errorf("保存 WebDav 配置失败: %w", err)
	}
	return "WebDav 配置已更新", nil
}

const DEFAULT_PASS_REPLACE = "(●'◡'●)"

func (bm *BackupManager) GetWebDavConfig(ctx context.Context) (responsex.WebDavConfigResponse, error) {
	// 调用 SettingService 获取 WebDav 配置
	cfg, err := bm.settingSvc.GetWebDavConfig()
	if err != nil {
		logger.Error("获取 WebDav 配置失败", logger.ErrorType(err))
		return responsex.WebDavConfigResponse{}, fmt.Errorf("获取 WebDav 配置失败: %w", err)
	}

	return responsex.WebDavConfigResponse{
		WebDavURL:      cfg.WebDavURL,
		WebDavUserName: cfg.WebDavUserName,
		WebDavPassword: DEFAULT_PASS_REPLACE,
	}, nil
}

// ListWebDavBackups 获取 WebDav 上的备份列表（仅返回 .db / .sqlite 文件）
func (bm *BackupManager) ListWebDavBackups(ctx context.Context) ([]responsex.WebDavBackupItem, error) {
	cfg, err := bm.settingSvc.GetWebDavConfig()
	if err != nil {
		logger.Error("读取 WebDav 配置失败", logger.ErrorType(err))
		return nil, fmt.Errorf("读取 WebDav 配置失败: %w", err)
	}

	if cfg.WebDavURL == "" {
		return nil, fmt.Errorf("未配置 WebDav")
	}

	client := bm.newWebDavClient(cfg.WebDavURL, cfg.WebDavUserName, cfg.WebDavPassword)
	if err := client.Connect(); err != nil {
		logger.Error("连接 WebDav 失败", logger.ErrorType(err))
		return nil, fmt.Errorf("连接 WebDav 失败: %w", err)
	}

	entries, err := bm.ensureWebDavBaseDir(client, cfg.WebDavBaseDir)
	if err != nil {
		return nil, err
	}
	logger.Debug("WebDav 目录文件列表获取成功", logger.Int("count", len(entries)))
	items := make([]responsex.WebDavBackupItem, 0)
	for _, fi := range entries {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		lower := strings.ToLower(name)
		if !(strings.HasSuffix(lower, ".db") || strings.HasSuffix(lower, ".sqlite")) {
			continue
		}
		items = append(items, responsex.WebDavBackupItem{
			Name:    name,
			Size:    fi.Size(),
			ModTime: fi.ModTime().Unix(),
			Path:    path.Join(cfg.WebDavBaseDir, name),
		})
	}

	return items, nil
}

// ensureWebDavBaseDir 确保基础目录存在并返回目录下的文件列表
func (bm *BackupManager) ensureWebDavBaseDir(client *gowebdav.Client, baseDir string) ([]os.FileInfo, error) {
	if _, err := client.Stat(baseDir); err != nil {
		// gowebdav 的 Stat 404 不总是转换为 os.ErrNotExist，这里直接尝试创建
		if mkErr := client.MkdirAll(baseDir, 0o755); mkErr != nil && !os.IsExist(mkErr) {
			logger.Error("创建 WebDav 目录失败", logger.ErrorType(mkErr))
			return nil, fmt.Errorf("创建 WebDav 目录失败: %w", mkErr)
		}
	}
	entries, err := client.ReadDir(baseDir)
	if err != nil {
		logger.Error("读取 WebDav 目录失败", logger.ErrorType(err))
		return nil, fmt.Errorf("读取 WebDav 目录失败: %w", err)
	}
	return entries, nil
}

func (bm *BackupManager) RestoreBackup(ctx context.Context, req *requestx.RestoreBackupRequest) (string, error) {
	return "restore success", nil
}

func (bm *BackupManager) CreateBackup(ctx context.Context, req *requestx.CreateBackupRequest) (string, error) {
	return "backup success", nil
}

func (bm *BackupManager) RegisterRoute(d *dispatcher.Dispatcher) {
	// 注册设置本地备份路径的路由
	dispatcher.RegisterTyped(d, "backup_manager/set_backup_local_path", bm.SetBackupLocalPath)

	// 注册获取本地备份路径的路由
	dispatcher.RegisterNoReq(d, "backup_manager/get_backup_local_path", bm.GetBackupLocalPath)

	// 注册打开目录选择对话框的路由
	dispatcher.RegisterTyped(d, "backup_manager/open_path_selector_dialog", bm.OpenPathSelectorDialog)

	// 注册设置 WebDav 配置的路由
	dispatcher.RegisterTyped(d, "backup_manager/set_webdav_config", bm.SetWebDavConfig)

	// 注册获取 WebDav 配置的路由
	dispatcher.RegisterNoReq(d, "backup_manager/get_webdav_config", bm.GetWebDavConfig)

	// 注册测试 WebDav 连接的路由
	dispatcher.RegisterTyped(d, "backup_manager/test_webdav_connection", bm.TestWebDavConnection)

	// 注册获取 WebDav 备份列表的路由
	dispatcher.RegisterNoReq(d, "backup_manager/list_webdav_backups", bm.ListWebDavBackups)

	// 注册恢复备份的路由
	dispatcher.RegisterTyped(d, "backup_manager/restore_backup", bm.RestoreBackup)

	// 注册手动创建备份的路由
	dispatcher.RegisterTyped(d, "backup_manager/create_backup", bm.CreateBackup)
}
