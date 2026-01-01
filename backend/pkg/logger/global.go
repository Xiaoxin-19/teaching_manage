package logger

import "sync"

// globalLogger 默认为 NopLogger，确保未初始化时调用不会 Panic
var (
	globalLogger Logger = NewNopLogger()
	mu           sync.RWMutex
)

// SetGlobalLogger 设置全局 Logger (并发安全)
func SetGlobalLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = l
}

// GetLogger 获取全局 Logger
func GetLogger() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return globalLogger
}

// --- 静态代理方法，方便直接调用 ---

func Debug(msg string, args ...Field) {
	GetLogger().Debug(msg, args...)
}

func Info(msg string, args ...Field) {
	GetLogger().Info(msg, args...)
}

func Warn(msg string, args ...Field) {
	GetLogger().Warn(msg, args...)
}

func Error(msg string, args ...Field) {
	GetLogger().Error(msg, args...)
}
