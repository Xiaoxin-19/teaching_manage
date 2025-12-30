package dao

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var global_db *gorm.DB

func InitDB(path string) error {
	// 确保数据库目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	// 在 DSN 中开启外键并启用 WAL 和共享缓存
	dsn := path + "?_foreign_keys=1&_cache=shared&_journal_mode=WAL"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return err
	}

	// 冗余确保外键和性能相关 PRAGMA 已启用
	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")
	db.Exec("PRAGMA synchronous = NORMAL")

	// 调整底层连接池（SQLite 文件 DB 建议限制连接数以避免写冲突）
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0)

	// 生产环境建议使用版本化迁移工具；AutoMigrate 可用于开发/快速原型
	if err := db.AutoMigrate(&Student{}, &Teacher{}, &Order{}, &Record{}); err != nil {
		return err
	}
	global_db = db
	return nil
}

func GetDB() *gorm.DB {
	return global_db
}
