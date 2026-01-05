package pathutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	ErrPathEmpty       = errors.New("路径不能为空")
	ErrPathNotExist    = errors.New("路径不存在")
	ErrPathNotDir      = errors.New("路径不是一个目录")
	ErrPathNotWritable = errors.New("路径没有写入权限")
	ErrPathInvalid     = errors.New("路径格式无效")
	ErrPathTooLong     = errors.New("路径长度超出限制")
)

const (
	// Windows 系统路径最大长度
	MaxPathLengthWindows = 260
	// Unix 系统路径最大长度
	MaxPathLengthUnix = 4096
)

// ValidateBackupPath 验证备份路径的有效性
// 检查路径是否存在、是否为目录、是否有写入权限
func ValidateBackupPath(path string) error {
	// 检查路径是否为空
	if strings.TrimSpace(path) == "" {
		return ErrPathEmpty
	}

	// 检查路径长度
	maxLen := MaxPathLengthUnix
	if runtime.GOOS == "windows" {
		maxLen = MaxPathLengthWindows
	}
	if len(path) > maxLen {
		return fmt.Errorf("%w: 最大长度为 %d", ErrPathTooLong, maxLen)
	}

	// 清理路径
	cleanPath := filepath.Clean(path)

	// 检查路径格式（基本格式检查）
	if !filepath.IsAbs(cleanPath) {
		return fmt.Errorf("%w: 请使用绝对路径", ErrPathInvalid)
	}

	// 检查路径是否存在
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrPathNotExist, cleanPath)
		}
		return fmt.Errorf("无法访问路径: %w", err)
	}

	// 检查是否为目录
	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrPathNotDir, cleanPath)
	}

	// 检查写入权限（尝试创建临时文件）
	testFile := filepath.Join(cleanPath, ".backup_test_write_permission")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrPathNotWritable, cleanPath)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}

// EnsureDirExists 确保目录存在，如果不存在则创建
func EnsureDirExists(path string) error {
	if strings.TrimSpace(path) == "" {
		return ErrPathEmpty
	}

	cleanPath := filepath.Clean(path)

	// 检查路径是否存在
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 目录不存在，尝试创建
			if err := os.MkdirAll(cleanPath, 0755); err != nil {
				return fmt.Errorf("创建目录失败: %w", err)
			}
			return nil
		}
		return fmt.Errorf("无法访问路径: %w", err)
	}

	// 路径存在，检查是否为目录
	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrPathNotDir, cleanPath)
	}

	return nil
}

// IsValidFilePath 验证文件路径的有效性
// 检查文件是否存在、是否为普通文件
func IsValidFilePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return ErrPathEmpty
	}

	cleanPath := filepath.Clean(path)

	// 检查路径格式
	if !filepath.IsAbs(cleanPath) {
		return fmt.Errorf("%w: 请使用绝对路径", ErrPathInvalid)
	}

	// 检查文件是否存在
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrPathNotExist, cleanPath)
		}
		return fmt.Errorf("无法访问文件: %w", err)
	}

	// 检查是否为普通文件
	if info.IsDir() {
		return fmt.Errorf("路径是一个目录，而非文件: %s", cleanPath)
	}

	return nil
}

// GetDirFromPath 从路径中提取目录
// 如果是文件路径，返回其父目录；如果是目录路径，返回该目录
func GetDirFromPath(path string) string {
	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		// 如果路径不存在，假定它是目录路径
		return cleanPath
	}

	if info.IsDir() {
		return cleanPath
	}

	// 是文件，返回父目录
	return filepath.Dir(cleanPath)
}

// NormalizePath 规范化路径，统一路径分隔符
func NormalizePath(path string) string {
	if path == "" {
		return ""
	}

	// 清理路径
	cleanPath := filepath.Clean(path)

	// 在 Windows 上将反斜杠转换为正斜杠（如果需要）
	if runtime.GOOS == "windows" {
		// 保留 Windows 路径格式
		return cleanPath
	}

	return cleanPath
}
