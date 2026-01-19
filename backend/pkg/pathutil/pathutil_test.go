package pathutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateBackupPath(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errType error
	}{
		{
			name:    "空路径",
			path:    "",
			wantErr: true,
			errType: ErrPathEmpty,
		},
		{
			name:    "空格路径",
			path:    "   ",
			wantErr: true,
			errType: ErrPathEmpty,
		},
		{
			name:    "有效的临时目录",
			path:    tempDir,
			wantErr: false,
		},
		{
			name:    "不存在的路径",
			path:    filepath.Join(tempDir, "nonexistent"),
			wantErr: true,
			errType: ErrPathNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBackupPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBackupPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnsureDirExists(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "空路径",
			path:    "",
			wantErr: true,
		},
		{
			name:    "已存在的目录",
			path:    tempDir,
			wantErr: false,
		},
		{
			name:    "需要创建的目录",
			path:    filepath.Join(tempDir, "new_dir"),
			wantErr: false,
		},
		{
			name:    "多层嵌套目录",
			path:    filepath.Join(tempDir, "a", "b", "c"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureDirExists(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDirExists() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 如果没有错误，验证目录是否真的存在
			if err == nil && tt.path != "" {
				info, statErr := os.Stat(tt.path)
				if statErr != nil {
					t.Errorf("目录创建后无法访问: %v", statErr)
				}
				if !info.IsDir() {
					t.Errorf("创建的路径不是目录")
				}
			}
		})
	}
}

func TestIsValidFilePath(t *testing.T) {
	tempDir := t.TempDir()

	// 创建一个测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "空路径",
			path:    "",
			wantErr: true,
		},
		{
			name:    "有效的文件路径",
			path:    testFile,
			wantErr: false,
		},
		{
			name:    "目录路径（应该失败）",
			path:    tempDir,
			wantErr: true,
		},
		{
			name:    "不存在的文件",
			path:    filepath.Join(tempDir, "nonexistent.txt"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDirFromPath(t *testing.T) {
	tempDir := t.TempDir()

	// 创建一个测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "目录路径",
			path: tempDir,
			want: tempDir,
		},
		{
			name: "文件路径",
			path: testFile,
			want: tempDir,
		},
		{
			name: "不存在的路径（假定为目录）",
			path: filepath.Join(tempDir, "nonexistent"),
			want: filepath.Join(tempDir, "nonexistent"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDirFromPath(tt.path)
			if got != tt.want {
				t.Errorf("GetDirFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "空路径",
			path: "",
			want: "",
		},
		{
			name: "简单路径",
			path: "/home/user/backup",
			want: filepath.Clean("/home/user/backup"),
		},
		{
			name: "带有多余分隔符的路径",
			path: "/home//user///backup",
			want: filepath.Clean("/home/user/backup"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizePath(tt.path)
			if got != tt.want {
				t.Errorf("NormalizePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
