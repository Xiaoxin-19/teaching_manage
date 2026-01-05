package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidKey        = errors.New("密钥长度无效")
	ErrDecryptionFailed  = errors.New("解密失败")
	ErrInvalidCiphertext = errors.New("密文格式无效")
)

// AES密钥长度常量
const (
	KeySize128 = 16 // 128-bit
	KeySize256 = 32 // 256-bit
)

// PasswordEncryptor 密码加密器
type PasswordEncryptor struct {
	key []byte // 加密密钥
}

// NewPasswordEncryptor 创建密码加密器
// key应该是16(AES-128)或32(AES-256)字节长
func NewPasswordEncryptor(key []byte) (*PasswordEncryptor, error) {
	if len(key) != KeySize128 && len(key) != KeySize256 {
		return nil, fmt.Errorf("%w: 密钥长度应为16或32字节，当前为%d", ErrInvalidKey, len(key))
	}

	return &PasswordEncryptor{
		key: key,
	}, nil
}

// Encrypt 加密密码
// 使用 AES-GCM 模式进行加密，返回 base64 编码的密文
func (pe *PasswordEncryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(pe.key)
	if err != nil {
		return "", fmt.Errorf("创建加密块失败: %w", err)
	}

	// 使用 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成nonce失败: %w", err)
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// 返回 base64 编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密密码
// 接收 base64 编码的密文，返回明文
func (pe *PasswordEncryptor) Decrypt(ciphertext string) (string, error) {
	// 解码 base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("%w: base64解码失败", ErrInvalidCiphertext)
	}

	block, err := aes.NewCipher(pe.key)
	if err != nil {
		return "", fmt.Errorf("创建解密块失败: %w", err)
	}

	// 使用 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("%w: 密文长度过短", ErrInvalidCiphertext)
	}

	// 提取 nonce 和实际的密文
	nonce, ciphertextData := data[:nonceSize], data[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
	}

	return string(plaintext), nil
}

// GenerateKey 生成随机密钥
// keySize应为16(AES-128)或32(AES-256)
func GenerateKey(keySize int) ([]byte, error) {
	if keySize != KeySize128 && keySize != KeySize256 {
		return nil, fmt.Errorf("%w: 密钥长度应为16或32字节", ErrInvalidKey)
	}

	key := make([]byte, keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("生成密钥失败: %w", err)
	}

	return key, nil
}

// KeyToBase64 将密钥转换为 base64 字符串
func KeyToBase64(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

// Base64ToKey 将 base64 字符串转换为密钥
func Base64ToKey(keyStr string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("密钥base64解码失败: %w", err)
	}

	if len(key) != KeySize128 && len(key) != KeySize256 {
		return nil, fmt.Errorf("%w: 解码后的密钥长度为%d", ErrInvalidKey, len(key))
	}

	return key, nil
}

// 全局加密密钥和加密器
var (
	globalEncryptionKey []byte
	globalEncryptor     *PasswordEncryptor
)

// InitGlobalEncryptor 初始化全局加密器
// 此函数应在应用启动时调用
func InitGlobalEncryptor(key []byte) error {
	var err error
	globalEncryptor, err = NewPasswordEncryptor(key)
	if err != nil {
		return err
	}
	globalEncryptionKey = key
	return nil
}

// EncryptPassword 加密密码（使用全局加密器）
func EncryptPassword(plaintext string) (string, error) {
	if globalEncryptor == nil {
		return "", fmt.Errorf("全局加密器未初始化，请先调用 InitGlobalEncryptor")
	}
	return globalEncryptor.Encrypt(plaintext)
}

// DecryptPassword 解密密码（使用全局加密器）
func DecryptPassword(ciphertext string) (string, error) {
	if globalEncryptor == nil {
		return "", fmt.Errorf("全局加密器未初始化，请先调用 InitGlobalEncryptor")
	}
	return globalEncryptor.Decrypt(ciphertext)
}
