package crypto

import (
	"encoding/base64"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		wantErr bool
	}{
		{
			name:    "生成128位密钥",
			keySize: KeySize128,
			wantErr: false,
		},
		{
			name:    "生成256位密钥",
			keySize: KeySize256,
			wantErr: false,
		},
		{
			name:    "无效的密钥长度",
			keySize: 24,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateKey(tt.keySize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(key) != tt.keySize {
				t.Errorf("GenerateKey() returned key of length %d, want %d", len(key), tt.keySize)
			}
		})
	}
}

func TestPasswordEncryptor(t *testing.T) {
	// 生成测试密钥
	key, err := GenerateKey(KeySize256)
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	// 创建加密器
	encryptor, err := NewPasswordEncryptor(key)
	if err != nil {
		t.Fatalf("创建加密器失败: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "简单密码",
			plaintext: "password123",
		},
		{
			name:      "特殊字符密码",
			plaintext: "P@ssw0rd!#$%",
		},
		{
			name:      "中文密码",
			plaintext: "密码123",
		},
		{
			name:      "长密码",
			plaintext: "this is a very long password with many characters and symbols!@#$%^&*()",
		},
		{
			name:      "空密码",
			plaintext: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 加密
			ciphertext, err := encryptor.Encrypt(tt.plaintext)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
				return
			}

			// 验证密文是有效的 base64
			_, err = base64.StdEncoding.DecodeString(ciphertext)
			if err != nil {
				t.Errorf("密文不是有效的 base64: %v", err)
				return
			}

			// 解密
			decrypted, err := encryptor.Decrypt(ciphertext)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}

			// 验证解密结果
			if decrypted != tt.plaintext {
				t.Errorf("解密结果不匹配: got %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

func TestNewPasswordEncryptor_InvalidKey(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		wantErr bool
	}{
		{
			name:    "密钥长度为8",
			keySize: 8,
			wantErr: true,
		},
		{
			name:    "密钥长度为24",
			keySize: 24,
			wantErr: true,
		},
		{
			name:    "密钥长度为64",
			keySize: 64,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invalidKey := make([]byte, tt.keySize)
			_, err := NewPasswordEncryptor(invalidKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPasswordEncryptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptionDeterminism(t *testing.T) {
	// 虽然 AES-GCM 使用随机 nonce，所以同样的明文加密结果会不同
	// 但解密后应该得到相同的明文
	key, err := GenerateKey(KeySize256)
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	encryptor, err := NewPasswordEncryptor(key)
	if err != nil {
		t.Fatalf("创建加密器失败: %v", err)
	}

	plaintext := "test_password"

	// 加密多次
	ciphertext1, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("第一次加密失败: %v", err)
	}

	ciphertext2, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("第二次加密失败: %v", err)
	}

	// 密文应该不同（因为使用了随机 nonce）
	if ciphertext1 == ciphertext2 {
		t.Errorf("相同明文的两次加密结果不应该相同（理论上概率极低）")
	}

	// 但解密结果应该相同
	decrypted1, err := encryptor.Decrypt(ciphertext1)
	if err != nil {
		t.Fatalf("解密1失败: %v", err)
	}

	decrypted2, err := encryptor.Decrypt(ciphertext2)
	if err != nil {
		t.Fatalf("解密2失败: %v", err)
	}

	if decrypted1 != decrypted2 || decrypted1 != plaintext {
		t.Errorf("解密结果不一致: got %q and %q, want %q", decrypted1, decrypted2, plaintext)
	}
}

func TestKeyConversion(t *testing.T) {
	// 生成原始密钥
	originalKey, err := GenerateKey(KeySize256)
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	// 转换为 base64
	keyStr := KeyToBase64(originalKey)

	// 从 base64 转换回来
	restoredKey, err := Base64ToKey(keyStr)
	if err != nil {
		t.Fatalf("从base64恢复密钥失败: %v", err)
	}

	// 验证密钥是否相同
	if len(originalKey) != len(restoredKey) {
		t.Errorf("密钥长度不匹配: %d vs %d", len(originalKey), len(restoredKey))
	}

	for i := range originalKey {
		if originalKey[i] != restoredKey[i] {
			t.Errorf("密钥内容不匹配")
			break
		}
	}

	// 验证转换后的密钥是否可用
	encryptor, err := NewPasswordEncryptor(restoredKey)
	if err != nil {
		t.Fatalf("使用转换后的密钥创建加密器失败: %v", err)
	}

	plaintext := "test"
	ciphertext, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("加密/解密失败: got %q, want %q", decrypted, plaintext)
	}
}

func TestInvalidCiphertext(t *testing.T) {
	key, err := GenerateKey(KeySize256)
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	encryptor, err := NewPasswordEncryptor(key)
	if err != nil {
		t.Fatalf("创建加密器失败: %v", err)
	}

	tests := []struct {
		name       string
		ciphertext string
		wantErr    bool
	}{
		{
			name:       "无效的base64",
			ciphertext: "invalid!!!",
			wantErr:    true,
		},
		{
			name:       "太短的密文",
			ciphertext: base64.StdEncoding.EncodeToString([]byte("short")),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptor.Decrypt(tt.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
