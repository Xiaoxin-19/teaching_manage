package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"teaching_manage/backend/dao"
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

const MaxBackupRetentionCount = 7

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

// resolveWebDavPassword
// 传入空值或占位符时，返回已保存的密码，否则返回传入的密码
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

	lastBackupInt, err := strconv.ParseInt(cfg.LastCloudBackup, 10, 64)
	if err != nil {
		lastBackupInt = 0
	}
	passwd := DEFAULT_PASS_REPLACE

	if cfg.WebDavPassword == "" {
		passwd = ""
	}
	return responsex.WebDavConfigResponse{
		WebDavURL:        cfg.WebDavURL,
		WebDavUserName:   cfg.WebDavUserName,
		WebDavPassword:   passwd,
		LastCloudBackup:  lastBackupInt,
		EnableAutoBackup: cfg.EnableAutoBackup,
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

// generateBackupFileName 生成备份文件名
func (bm *BackupManager) generateBackupFileName(currentDBPath string) string {
	timestamp := time.Now().Format("20060102_150405")
	dbFileName := path.Base(currentDBPath)
	return fmt.Sprintf("%s_backup_%s", strings.TrimSuffix(dbFileName, path.Ext(dbFileName)), timestamp+path.Ext(dbFileName))
}

// dbReloadWrapper 封装需要关闭和重新打开数据库连接的操作
// 适用于 backup 和 restore 等需要独占访问数据库文件的操作
func dbReloadWrapper(fn func() (string, error)) (string, error) {
	// 关闭数据库连接
	if err := dao.CloseDB(); err != nil {
		logger.Error("close database connection failed", logger.ErrorType(err))
		return "", fmt.Errorf("close database connection failed: %w", err)
	}

	logger.Debug("close database connection success")
	time.Sleep(500 * time.Millisecond) // 确保连接关闭完成
	// 恢复数据库连接（无论操作成功与否）
	defer func() {
		// 使用重试机制重新打开数据库连接，避免固定长时间 sleep
		const (
			reopenTimeout  = 5 * time.Second
			reopenInterval = 500 * time.Millisecond
		)

		start := time.Now()
		for {
			// 在超时时间内，记录错误并短暂等待后重试
			if err := dao.ReopenDB(); err != nil {
				if time.Since(start) >= reopenTimeout {
					logger.Error("re open database failed after retries in defer", logger.ErrorType(err))
					break
				}

				logger.Error("re open database failed in defer, will retry", logger.ErrorType(err))
				time.Sleep(reopenInterval)
				continue
			}

			logger.Debug("re connect to database susccess")
			break
		}
	}()

	// 执行操作
	result, err := fn()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (bm *BackupManager) RestoreBackup(ctx context.Context, req *requestx.RestoreBackupRequest) (string, error) {
	return dbReloadWrapper(func() (string, error) {
		switch req.Type {
		case "local":
			return bm.restoreBackupLocal(ctx, req.BackupPath)
		case "webdav":
			return bm.restoreBackupWebDav(ctx, req.BackupPath)
		}
		return "restore success", nil
	})
}

func (bm *BackupManager) restoreBackupLocal(_ context.Context, path string) (string, error) {
	// 覆盖数据库文件, 先备份现有数据库以防万一
	currentDBPath, err := bm.settingSvc.GetCurrentDBPath()
	if err != nil {
		logger.Error("get current database path failed", logger.ErrorType(err))
		return "", err
	}
	backupPath := currentDBPath + ".bak"
	if err := os.Rename(currentDBPath, backupPath); err != nil {
		logger.Error("backup current database by rename file failed", logger.ErrorType(err))
		return "", err
	}

	backupSuccess := false
	defer func() {
		if r := recover(); r != nil {
			// 恢复备份文件
			if err := os.Rename(backupPath, currentDBPath); err != nil {
				logger.Error("restore database from backup file failed in defer", logger.ErrorType(err))
			}
		} else if !backupSuccess {
			// 恢复备份文件
			if err := os.Rename(backupPath, currentDBPath); err != nil {
				logger.Error("restore database from backup file failed in defer", logger.ErrorType(err))
			}
		}
	}()

	// 复制备份文件到数据库路径
	input, err := os.Open(path)
	if err != nil {
		logger.Error("open backup file failed", logger.ErrorType(err))
		return "", err
	}
	defer input.Close()

	output, err := os.Create(currentDBPath)
	if err != nil {
		logger.Error("create database file failed", logger.ErrorType(err))
		return "", err
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		logger.Error("copy backup file to database file failed", logger.ErrorType(err))
		return "", err
	}

	// 确保文件写入完成
	if err := output.Sync(); err != nil {
		logger.Error("sync database file failed", logger.ErrorType(err))
		return "", err
	}

	// 删除临时备份文件
	if err := os.Remove(backupPath); err != nil {
		logger.Warn("remove temporary backup database file failed", logger.ErrorType(err))
	}

	backupSuccess = true
	return "已恢复本地备份", nil
}

func (bm *BackupManager) restoreBackupWebDav(_ context.Context, backupPath string) (string, error) {
	// 获取 WebDav 配置
	cfg, err := bm.settingSvc.GetWebDavConfig()
	if err != nil {
		logger.Error("get webdav config failed", logger.ErrorType(err))
		return "", err
	}

	client := bm.newWebDavClient(cfg.WebDavURL, cfg.WebDavUserName, cfg.WebDavPassword)
	if err := client.Connect(); err != nil {
		logger.Error("connect to WebDav failed", logger.ErrorType(err))
		return "", err
	}

	// 覆盖数据库文件, 先备份现有数据库以防万一
	currentDBPath, err := bm.settingSvc.GetCurrentDBPath()
	if err != nil {
		logger.Error("get current database path failed", logger.ErrorType(err))
		return "", err
	}
	backupPathLocal := currentDBPath + ".bak"
	if err := os.Rename(currentDBPath, backupPathLocal); err != nil {
		logger.Error("backup current database by rename file failed", logger.ErrorType(err))
		return "", err
	}

	restoreSuccess := false
	defer func() {
		if !restoreSuccess {
			// 恢复失败时回滚到备份文件
			if err := os.Rename(backupPathLocal, currentDBPath); err != nil {
				logger.Error("restore database from backup file failed in defer", logger.ErrorType(err))
			}
		}
	}()

	// 从 WebDav 下载备份文件到数据库路径
	input, err := client.ReadStream(backupPath)
	if err != nil {
		logger.Error("open backup file from webdav failed", logger.ErrorType(err))
		return "", err
	}
	defer input.Close()

	output, err := os.Create(currentDBPath)
	if err != nil {
		logger.Error("create database file failed", logger.ErrorType(err))
		return "", err
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		logger.Error("copy backup file to database file failed", logger.ErrorType(err))
		return "", err
	}

	// 确保文件写入完成
	if err := output.Sync(); err != nil {
		logger.Error("sync database file failed", logger.ErrorType(err))
		return "", err
	}

	// 删除临时备份文件
	if err := os.Remove(backupPathLocal); err != nil {
		logger.Warn("remove temporary backup file failed", logger.ErrorType(err))
	}

	restoreSuccess = true
	return "备份已经恢复", nil
}

func (bm *BackupManager) CreateBackup(ctx context.Context, req *requestx.CreateBackupRequest) (string, error) {
	result, err := dbReloadWrapper(func() (string, error) {
		switch req.Type {
		case "local":
			return bm.CreateBackupLocal(ctx)
		case "webdav":
			// 仅执行备份逻辑，数据库重连后再写入最后备份时间，避免 DB 已关闭导致空指针
			if _, err := bm.CreateBackupWebDav(ctx); err != nil {
				return "", err
			}
			return "备份成功", nil
		}
		return "备份成功", nil
	})
	if err != nil {
		return "", err
	}

	// 确保数据库已在 dbReloadWrapper 的 defer 中重连完毕后，再记录最后备份时间
	if req.Type == "webdav" {
		if err := bm.settingSvc.UpdateLastWebDavBackupTime(time.Now()); err != nil {
			logger.Warn("update last webdav backup time failed", logger.ErrorType(err))
		}
	}

	return result, nil
}

func (bm *BackupManager) CreateBackupLocal(ctx context.Context) (string, error) {
	// 获取本地备份路径
	localPath, err := bm.settingSvc.GetBackupLocalPath()
	if err != nil {
		logger.Error("get backup local path failed", logger.ErrorType(err))
		return "", err
	}

	if localPath == "" {
		return "", fmt.Errorf("本地备份路径未设置")
	}

	// 获取当前数据库路径
	currentDBPath, err := bm.settingSvc.GetCurrentDBPath()
	if err != nil {
		logger.Error("get current database path failed", logger.ErrorType(err))
		return "", err
	}

	// 构造备份文件名
	backupFileName := bm.generateBackupFileName(currentDBPath)
	backupFilePath := path.Join(localPath, backupFileName)

	backupSuccess := false
	defer func() {
		if !backupSuccess {
			// 备份失败时删除不完整的文件
			if err := os.Remove(backupFilePath); err != nil && !os.IsNotExist(err) {
				logger.Warn("remove incomplete backup file failed", logger.ErrorType(err))
			}
		}
	}()

	// 复制数据库文件到备份路径
	input, err := os.Open(currentDBPath)
	if err != nil {
		logger.Error("open database file failed", logger.ErrorType(err))
		return "", err
	}
	defer input.Close()

	output, err := os.Create(backupFilePath)
	if err != nil {
		logger.Error("create backup file failed", logger.ErrorType(err))
		return "", err
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		logger.Error("copy database file to backup file failed", logger.ErrorType(err))
		return "", err
	}

	// 确保文件写入完成
	if err := output.Sync(); err != nil {
		logger.Error("sync backup file failed", logger.ErrorType(err))
		return "", err
	}

	backupSuccess = true

	// 清理旧备份，保留最近7份
	bm.cleanOldLocalBackups(localPath)

	return fmt.Sprintf("本地备份已创建: %s", backupFilePath), nil
}

// 创建 WebDav 备份， 保留最近7份备份文件，删除多余的旧备份
func (bm *BackupManager) CreateBackupWebDav(ctx context.Context) (string, error) {
	logger.Info("start create webdav backup")
	// 获取 WebDav 配置
	cfg, err := bm.settingSvc.GetWebDavConfig()
	if err != nil {
		logger.Error("get webdav config failed", logger.ErrorType(err))
		return "", err
	}

	logger.Debug("connecting to WebDav server")
	client := bm.newWebDavClient(cfg.WebDavURL, cfg.WebDavUserName, cfg.WebDavPassword)
	if err := client.Connect(); err != nil {
		logger.Error("connect to WebDav failed", logger.ErrorType(err))
		return "", err
	}

	logger.Debug("connected to WebDav server")
	// 获取当前数据库路径
	currentDBPath, err := bm.settingSvc.GetCurrentDBPath()
	if err != nil {
		logger.Error("get current database path failed", logger.ErrorType(err))
		return "", err
	}

	// 构造备份文件名
	backupFileName := bm.generateBackupFileName(currentDBPath)
	backupFilePath := path.Join(cfg.WebDavBaseDir, backupFileName)

	logger.Debug("start upload backup file to WebDav", logger.String("backup_file_path", backupFilePath))
	// 上传数据库文件到 WebDav 备份路径
	input, err := os.Open(currentDBPath)
	if err != nil {
		logger.Error("open database file failed", logger.ErrorType(err))
		return "", err
	}
	defer input.Close()

	logger.Debug("uploading backup file to WebDav", logger.String("backup_file_path", backupFilePath))
	if err := client.WriteStream(backupFilePath, input, 0); err != nil {
		logger.Error("upload backup file to WebDav failed", logger.ErrorType(err))
		return "", err
	}

	logger.Debug("uploaded backup file to WebDav", logger.String("backup_file_path", backupFilePath))

	// 上传成功
	logger.Info("WebDav backup created successfully", logger.String("backup_file_path", backupFilePath))
	// 清理旧备份，保留最近7份
	bm.cleanOldWebDavBackups(client, cfg.WebDavBaseDir)
	logger.Debug("cleaned old WebDav backups")
	return fmt.Sprintf("WebDav 备份已创建: %s", backupFilePath), nil
}

// cleanOldLocalBackups 清理本地旧备份，保留最近7份
func (bm *BackupManager) cleanOldLocalBackups(localPath string) {
	entries, err := os.ReadDir(localPath)
	if err != nil {
		logger.Warn("read local backup directory failed", logger.ErrorType(err))
		return
	}

	var backupFiles []os.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		lower := strings.ToLower(name)
		if !(strings.HasSuffix(lower, ".db") || strings.HasSuffix(lower, ".sqlite")) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		backupFiles = append(backupFiles, info)
	}

	if len(backupFiles) > MaxBackupRetentionCount {
		sort.Slice(backupFiles, func(i, j int) bool {
			return backupFiles[i].ModTime().Before(backupFiles[j].ModTime())
		})
		for i := 0; i < len(backupFiles)-MaxBackupRetentionCount; i++ {
			oldBackupPath := path.Join(localPath, backupFiles[i].Name())
			if err := os.Remove(oldBackupPath); err != nil {
				logger.Warn("remove old local backup file failed", logger.String("file", oldBackupPath), logger.ErrorType(err))
			}
		}
	}
}

// cleanOldWebDavBackups 清理 WebDav 旧备份，保留最近7份
func (bm *BackupManager) cleanOldWebDavBackups(client *gowebdav.Client, baseDir string) {
	entries, err := bm.ensureWebDavBaseDir(client, baseDir)
	if err != nil {
		logger.Warn("ensure webdav base dir failed", logger.ErrorType(err))
		return
	}

	var backupFiles []os.FileInfo
	for _, fi := range entries {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		lower := strings.ToLower(name)
		if !(strings.HasSuffix(lower, ".db") || strings.HasSuffix(lower, ".sqlite")) {
			continue
		}
		backupFiles = append(backupFiles, fi)
	}

	if len(backupFiles) > MaxBackupRetentionCount {
		sort.Slice(backupFiles, func(i, j int) bool {
			return backupFiles[i].ModTime().Before(backupFiles[j].ModTime())
		})
		for i := 0; i < len(backupFiles)-MaxBackupRetentionCount; i++ {
			oldBackupPath := path.Join(baseDir, backupFiles[i].Name())
			if err := client.Remove(oldBackupPath); err != nil {
				logger.Warn("remove old webdav backup file failed", logger.String("file", oldBackupPath), logger.ErrorType(err))
			}
		}
	}
}

func (bm *BackupManager) SetAutoBackupEnabled(ctx context.Context, req *requestx.SetAutoBackupRequest) (string, error) {
	if err := bm.settingSvc.UpdateAutoBackupEnabled(req.Enabled); err != nil {
		logger.Error("update auto backup enabled failed", logger.ErrorType(err))
		return "", fmt.Errorf("更新自动备份设置失败: %w", err)
	}
	return "自动备份设置已更新", nil
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

	// 注册更新自动备份设置的路由
	dispatcher.RegisterTyped(d, "backup_manager/set_auto_backup_enabled", bm.SetAutoBackupEnabled)
}
