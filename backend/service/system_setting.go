package service

import (
	"sync"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/pkg/crypto"
)

// SettingService 负责管理系统配置
type SettingService struct {
	settingDAO dao.SettingDAO

	// 简单的内存缓存，避免频繁查库
	cache      *dao.SystemSetting
	cacheMutex sync.RWMutex
}

func (s *SettingService) GetBackupLocalPath() (string, error) {
	setting, err := s.getSettingWithCache()
	if err != nil {
		return "", err
	}
	return setting.LocalBackupDirPath, nil
}

func NewSettingService(settingDAO dao.SettingDAO) *SettingService {
	return &SettingService{
		settingDAO: settingDAO,
	}
}

// GetWebDavConfig 供其他 Service (如 BackupService) 获取 WebDAV 配置
// 返回时自动解密密码
func (s *SettingService) GetWebDavConfig() (entity.WebDavSetting, error) {
	setting, err := s.getSettingWithCache()
	if err != nil {
		return entity.WebDavSetting{}, err
	}

	// 解密密码
	decryptedPassword := setting.WebDavPassword
	if setting.WebDavPassword != "" {
		pwd, err := crypto.DecryptPassword(setting.WebDavPassword)
		if err != nil {
			return entity.WebDavSetting{}, err
		}
		decryptedPassword = pwd
	}

	return entity.WebDavSetting{
		WebDavURL:       setting.WebDavURL,
		WebDavUserName:  setting.WebDavUserName,
		WebDavPassword:  decryptedPassword,
		LastCloudBackup: setting.LastCloudBackup,
		WebDavBaseDir:   setting.WebDavBaseDir,
	}, nil
}

// UpdateWebDavConfig 供内部逻辑更新配置
// 自动加密密码后存储
func (s *SettingService) UpdateWebDavConfig(config entity.WebDavSetting) error {
	// 加密密码
	encryptedPassword := config.WebDavPassword
	if config.WebDavPassword != "" {
		pwd, err := crypto.EncryptPassword(config.WebDavPassword)
		if err != nil {
			return err
		}
		encryptedPassword = pwd
	}

	// 构造更新对象
	setting := &dao.SystemSetting{
		WebDavURL:       config.WebDavURL,
		WebDavUserName:  config.WebDavUserName,
		WebDavPassword:  encryptedPassword,
		LastCloudBackup: config.LastCloudBackup,
	}

	// 调用 DAO 更新指定字段
	err := s.settingDAO.UpdateSystemSetting(setting,
		&setting.WebDavURL,
		&setting.WebDavUserName,
		&setting.WebDavPassword,
	)

	if err != nil {
		return err
	}

	// 失效缓存
	s.invalidateCache()
	return nil
}

func (s *SettingService) UpdateLocalBackupPath(path string) error {
	// 1. 构造更新对象
	setting := &dao.SystemSetting{
		LocalBackupDirPath: path,
	}
	// 2. 调用 DAO 更新指定字段
	err := s.settingDAO.UpdateSystemSetting(setting,
		&setting.LocalBackupDirPath,
	)
	if err != nil {
		return err
	}
	// 3. 失效缓存
	s.invalidateCache()
	return nil
}

// getSettingWithCache 内部辅助方法：带缓存读取
func (s *SettingService) getSettingWithCache() (*dao.SystemSetting, error) {
	s.cacheMutex.RLock()
	if s.cache != nil {
		defer s.cacheMutex.RUnlock()
		return s.cache, nil
	}
	s.cacheMutex.RUnlock()

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// 双重检查
	if s.cache != nil {
		return s.cache, nil
	}

	var setting dao.SystemSetting
	if err := s.settingDAO.GetSystemSetting(&setting); err != nil {
		return nil, err
	}

	s.cache = &setting
	return &setting, nil
}

func (s *SettingService) invalidateCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.cache = nil
}
