package dao

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"teaching_manage/backend/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	globalDB      *gorm.DB
	currentDBPath string
	providerMu    sync.RWMutex
	provider      DBProvider = &refreshableProvider{}
)

// DBProvider 提供可刷新 DB 指针，用于备份/恢复后的重连
type DBProvider interface {
	DB() *gorm.DB
	Set(db *gorm.DB)
}

type refreshableProvider struct {
	mu sync.RWMutex
	db *gorm.DB
}

func (p *refreshableProvider) DB() *gorm.DB {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.db
}

func (p *refreshableProvider) Set(db *gorm.DB) {
	p.mu.Lock()
	p.db = db
	p.mu.Unlock()
}

var ErrDuplicatedKey = gorm.ErrDuplicatedKey
var ErrRecordNotFound = gorm.ErrRecordNotFound

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
	if err := db.AutoMigrate(
		&model.Student{},
		&model.Teacher{},
		&model.Record{},
		&model.Subject{},
		&model.StudentSubject{},
		&model.RechargeOrder{},
		&model.Setting{},
	); err != nil {
		return err
	}
	providerMu.Lock()
	currentDBPath = path
	globalDB = db
	providerMu.Unlock()
	provider.Set(db)
	return nil
}

// DBGetter 是用于获取 *gorm.DB 的函数类型。
//
// 设计意图：
//   - 调用方不直接持有 *gorm.DB 指针，而是持有一个“获取 DB 的函数”，
//     在每次需要访问数据库时通过调用该函数拿到当前应当使用的 *gorm.DB；
//   - 在数据库重连、备份/恢复或在全局 DB 与事务 DB（tx）之间切换时，
//     能够透明地获得最新、正确的连接，避免长期持有陈旧 DB 指针导致的并发问题或操作已关闭连接。
//
// 典型使用场景：
//   - DAO / 仓储层依赖注入时，优先接收 DBGetter 而不是裸 *gorm.DB 指针，
//     便于在普通查询与事务之间无感切换；
//   - 长生命周期对象（如 service、dispatcher）避免在初始化时缓存全局 DB，
//     通过 DBGetter 始终按需获取当前可用连接；
//   - 配合 DBProvider，在备份/恢复后刷新底层 DB 时，使上层代码无需感知重建过程。
//
// 调用者通过函数指针获取 DB，始终能获取最新的 DB 实例，
// 避免持有陈旧的 DB 指针导致的并发竞争问题。
type DBGetter func() *gorm.DB

// NewDBGetter 返回一个 DBGetter，始终返回传入的 db 实例。
// 通常用于事务场景下，将事务对象 db（例如 tx）传入，
// 可能原因：
// - 尚未初始化数据库（未调用 dao.InitDB(path)）
// - 数据库连接已通过 dao.CloseDB() 关闭，且尚未调用 dao.ReopenDB()
// 建议：
// - 在应用启动阶段调用 dao.InitDB(path) 完成数据库初始化
// - 若执行了数据库备份/恢复，请按顺序调用 dao.CloseDB() 和 dao.ReopenDB() 重新建立连接`)
func NewDBGetter(db *gorm.DB) DBGetter {
	return func() *gorm.DB {
		return db
	}
}

// GetDBTarget 是 NewDBGetter 的别名，用于兼容性（事务处理等场景）
var GetDBTarget = NewDBGetter

func GetDB() *gorm.DB {
	db := provider.DB()
	if db == nil {
		panic("database connection is not initialized or has been closed")
	}
	return db
}

// GetDBProvider 返回一个 DBGetter 函数，该函数始终通过 provider 获取最新的 DB 实例。
// 这种设计避免了调用者长期持有陈旧 DB 指针导致的并发竞争问题。
//
// 并发安全保证：
// - DBGetter 函数可以被安全地并发调用
// - 即使在 CloseDB() 或 ReopenDB() 期间调用，也能正确刷新获取新的 DB 实例
// - 不建议调用者缓存返回的 *gorm.DB；每次需要时应调用 DBGetter() 重新获取
//
// 注意：DBGetter 返回的 *gorm.DB 可能为 nil（当数据库已关闭时），调用者应检查或使用 GetDB()。
//
// 用法示例：
//
//	getDB := dao.GetDBProvider()
//	db := getDB()  // 获取当前 DB，总是最新的
//	db.Create(&model)
func GetDBProvider() DBGetter {
	return func() *gorm.DB {
		return provider.DB()
	}
}

// CloseDB 关闭底层连接，备份/恢复前调用。
// 并发安全：内部通过互斥锁保护全局变量的更新。
func CloseDB() error {
	providerMu.Lock()
	defer providerMu.Unlock()
	if globalDB == nil {
		// 若在此之前从未调用 InitDB，则会返回错误 "数据库路径未设置，请先调用 InitDB 初始化数据库"。
		return fmt.Errorf("数据库未初始化或已关闭")
	}
	sqlDB, err := globalDB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	globalDB = nil
	currentDBPath = ""
	provider.Set(nil)
	return nil
}

// ReopenDB 重新打开数据库连接。
//
// 用途：根据最近一次成功调用 InitDB 时记录的数据库路径，重新初始化全局数据库连接
// 并刷新底层 DBProvider。
//
// 调用场景：通常用于执行数据库文件备份或恢复操作之后，先调用 CloseDB 关闭当前连接，
// 再调用 ReopenDB 使用原有配置重新建立连接。
//
// 前置条件：必须至少成功调用过一次 InitDB，使 currentDBPath 记录了有效的数据库路径。
// 若在此之前从未调用 InitDB，则会返回错误 "db path not set"。
func ReopenDB() error {
	providerMu.RLock()
	path := currentDBPath
	providerMu.RUnlock()
	if path == "" {
		return fmt.Errorf("db path not set")
	}
	err := InitDB(path)
	if err != nil {
		panic("re set up database fail")
	}
	return nil
}
