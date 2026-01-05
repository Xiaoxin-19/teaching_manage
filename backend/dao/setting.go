package dao

import (
	"fmt"
	"reflect"
	"teaching_manage/backend/model"
	"teaching_manage/backend/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SystemSetting struct {
	// 若需新增配置，只需在此处添加字段，并在 json tag 中定义 key 即可
	WebDavURL          string `json:"web_dav_url" default:""`
	WebDavUserName     string `json:"web_dav_user_name" default:""`
	WebDavPassword     string `json:"web_dav_password" default:""`
	LocalBackupDirPath string `json:"local_backup_dir_path" default:""`
	WebDavBaseDir      string `json:"web_dav_base_dir" default:"/TeachingManageBackups"`
	LastCloudBackup    string `json:"last_cloud_backup" default:""`
}

type SettingDAO interface {
	// GetSystemSetting 获取配置。
	// setting: 用于接收数据的结构体指针。
	// targets: 可选。指定要获取的字段指针 (如 &s.WebDavURL)。若不传，则获取所有字段。
	GetSystemSetting(setting *SystemSetting, targets ...interface{}) error

	// UpdateSystemSetting 更新配置。
	// setting: 包含最新数据的结构体指针。
	// targets: 可选。指定要更新的字段指针 (如 &s.WebDavURL)。若不传，则更新所有字段。
	UpdateSystemSetting(setting *SystemSetting, targets ...interface{}) error
}

type SettingGORMDAO struct {
	db *gorm.DB
}

func NewSettingDAO(db *gorm.DB) SettingDAO {
	return &SettingGORMDAO{
		db: db,
	}
}

// GetSystemSetting 实现
func (dao *SettingGORMDAO) GetSystemSetting(setting *SystemSetting, targets ...interface{}) error {
	if setting == nil {
		return fmt.Errorf("setting cannot be nil")
	}

	// 1. 解析传入的指针，获取对应的 JSON Tag (数据库 Key)
	targetKeys := dao.resolveJsonTags(setting, targets...)

	// 2. 准备查询
	var dbSettings []model.Setting
	query := dao.db.Model(&model.Setting{})

	// 如果指定了特定字段，则添加过滤条件
	if len(targetKeys) > 0 {
		query = query.Where("`key` IN ?", targetKeys)
	}

	if err := query.Find(&dbSettings).Error; err != nil {
		logger.Error("failed to find settings", logger.ErrorType(err))
		return fmt.Errorf("failed to find settings: %w", err)
	}

	// 3. 构建 Map: Key -> Value
	keyMap := make(map[string]string)
	for _, s := range dbSettings {
		keyMap[s.Key] = s.Value
	}

	// 4. 反射回填结构体
	v := reflect.ValueOf(setting).Elem()
	t := v.Type()

	// 制作过滤 Set 用于快速判断
	checkFilter := len(targetKeys) > 0
	filterSet := make(map[string]struct{})
	for _, k := range targetKeys {
		filterSet[k] = struct{}{}
	}

	for i := 0; i < v.NumField(); i++ {
		fieldInfo := t.Field(i)
		jsonTag := fieldInfo.Tag.Get("json")

		if jsonTag == "" {
			continue
		}

		// 过滤逻辑：如果指定了 targets 且当前字段不在 targets 中，跳过
		if checkFilter {
			if _, ok := filterSet[jsonTag]; !ok {
				continue
			}
		}

		// 赋值逻辑
		fieldVal := v.Field(i)
		if fieldVal.CanSet() && fieldVal.Kind() == reflect.String {
			if dbVal, exist := keyMap[jsonTag]; exist {
				// A. 数据库有值 -> 用数据库的
				fieldVal.SetString(dbVal)
			} else {
				// B. 数据库无值 -> 用 default tag
				fieldVal.SetString(fieldInfo.Tag.Get("default"))
			}
		}
	}

	return nil
}

// UpdateSystemSetting 实现
func (dao *SettingGORMDAO) UpdateSystemSetting(setting *SystemSetting, targets ...interface{}) error {
	if setting == nil {
		return fmt.Errorf("setting cannot be nil")
	}

	// 1. 解析目标 Keys
	targetKeys := dao.resolveJsonTags(setting, targets...)

	// 2. 准备反射数据
	v := reflect.ValueOf(setting).Elem()
	t := v.Type()

	checkFilter := len(targetKeys) > 0
	filterSet := make(map[string]struct{})
	for _, k := range targetKeys {
		filterSet[k] = struct{}{}
	}

	var settingsToSave []model.Setting

	// 3. 遍历结构体，收集需要保存的数据
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := t.Field(i)
		jsonTag := fieldInfo.Tag.Get("json")

		if jsonTag == "" {
			continue
		}

		// 过滤逻辑
		if checkFilter {
			if _, ok := filterSet[jsonTag]; !ok {
				continue
			}
		}

		fieldVal := v.Field(i)
		// 仅支持字符串存储
		if fieldVal.Kind() == reflect.String {
			settingsToSave = append(settingsToSave, model.Setting{
				Key:   jsonTag,
				Value: fieldVal.String(),
			})
		}
	}

	if len(settingsToSave) == 0 {
		return nil
	}

	// 4. 批量 Upsert (插入或更新)
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},              // 冲突检测列
			DoUpdates: clause.AssignmentColumns([]string{"value"}), // 冲突时更新列
		}).Create(&settingsToSave).Error
	})

	if err != nil {
		logger.Error("failed to update settings", logger.ErrorType(err))
		return fmt.Errorf("failed to update settings: %w", err)
	}

	return nil
}

// resolveJsonTags 将传入的结构体字段指针转换为对应的 json tag 字符串列表
func (dao *SettingGORMDAO) resolveJsonTags(s *SystemSetting, targets ...interface{}) []string {
	var keys []string

	// 如果没有传 targets，返回空切片，代表“全部”
	if len(targets) == 0 {
		return keys
	}

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	// 建立内存地址到 JSON Tag 的映射
	addrToTag := make(map[uintptr]string)

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldInfo := t.Field(i)
		tag := fieldInfo.Tag.Get("json")

		if tag != "" {
			// 获取字段的内存地址
			addrToTag[fieldVal.Addr().Pointer()] = tag
		}
	}

	// 匹配传入的指针
	for _, target := range targets {
		val := reflect.ValueOf(target)
		// 确保传入的是指针
		if val.Kind() != reflect.Ptr {
			logger.Warn("passed non-pointer to targets, ignoring", logger.Field{Key: "type", Val: val.Type().String()})
			continue
		}

		ptr := val.Pointer()
		if tag, ok := addrToTag[ptr]; ok {
			keys = append(keys, tag)
		}
	}

	return keys
}
