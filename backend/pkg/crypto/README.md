# 密码加密与解密使用指南

## 概述

该项目使用 AES-256-GCM 算法实现了密码的加密与解密功能，确保敏感信息（如 WebDAV 密码）在数据库中安全存储。

## 核心特性

- **算法**: AES-256-GCM (256位密钥、认证加密)
- **密钥生成**: 随机生成或从配置读取
- **Nonce**: 每次加密时生成随机 nonce，确保相同明文的加密结果不同
- **编码**: 使用 Base64 编码密文用于存储和传输

## 使用流程

### 1. 初始化加密器

```go
// 生成或读取加密密钥（必须是16或32字节）
key, err := crypto.GenerateKey(crypto.KeySize256)
if err != nil {
    log.Fatal(err)
}

// 创建加密器
encryptor, err := crypto.NewPasswordEncryptor(key)
if err != nil {
    log.Fatal(err)
}
```

### 2. 加密密码

```go
plainPassword := "user_password"

// 加密（返回 base64 编码的密文）
encryptedPassword, err := encryptor.Encrypt(plainPassword)
if err != nil {
    log.Fatal(err)
}

// 保存 encryptedPassword 到数据库
```

### 3. 解密密码

```go
// 从数据库读取加密的密码
encryptedPassword := "..." // 从数据库获取

// 解密
plainPassword, err := encryptor.Decrypt(encryptedPassword)
if err != nil {
    log.Fatal(err)
}

// 使用明文密码
```

## 在 SettingService 中的应用

### WebDAV 配置加密

```go
// 保存时自动加密
config := entity.WebDavSetting{
    WebDavURL:      "https://dav.example.com",
    WebDavUserName: "user",
    WebDavPassword: "plaintext_password",
}
settingService.UpdateWebDavConfig(config) // 自动加密密码

// 读取时自动解密
webdavConfig, _ := settingService.GetWebDavConfig()
// webdavConfig.WebDavPassword 是解密后的明文
```

## 密钥管理最佳实践

### 1. 密钥存储

**DO:**
- ✅ 将密钥存储在环境变量中
- ✅ 将密钥存储在 .env 文件中（不提交到版本控制）
- ✅ 将密钥存储在密钥管理服务中（如 AWS KMS、Vault）

**DON'T:**
- ❌ 将密钥硬编码到代码中（开发环境除外）
- ❌ 将密钥提交到版本控制系统
- ❌ 在日志中打印密钥

### 2. 生产环境密钥读取示例

```go
// 从环境变量读取密钥
func getEncryptionKey() ([]byte, error) {
    keyStr := os.Getenv("ENCRYPTION_KEY")
    if keyStr == "" {
        return nil, errors.New("ENCRYPTION_KEY 环境变量未设置")
    }
    return crypto.Base64ToKey(keyStr)
}

// 使用
key, err := getEncryptionKey()
encryptor, _ := crypto.NewPasswordEncryptor(key)
```

### 3. 密钥轮换

当需要更换密钥时：

```go
// 1. 使用旧密钥解密所有密码
oldEncryptor, _ := crypto.NewPasswordEncryptor(oldKey)

// 2. 使用新密钥重新加密
newEncryptor, _ := crypto.NewPasswordEncryptor(newKey)

// 3. 批量更新数据库中的密码
for each password in database {
    plaintext, _ := oldEncryptor.Decrypt(password)
    newEncrypted, _ := newEncryptor.Encrypt(plaintext)
    updatePassword(newEncrypted)
}
```

## 安全考虑

1. **Nonce 随机性**: 每次加密时都生成新的随机 nonce，确保同一密码的多次加密结果不同
2. **认证加密**: 使用 GCM 模式提供了认证功能，防止密文被篡改
3. **密钥强度**: 使用 256 位密钥提供足够的加密强度
4. **错误处理**: 解密失败时应该安全处理，而不是暴露详细错误信息

## API 参考

### PasswordEncryptor

```go
// 创建加密器
encryptor, err := crypto.NewPasswordEncryptor(key []byte) (*PasswordEncryptor, error)

// 加密
ciphertext, err := encryptor.Encrypt(plaintext string) (string, error)

// 解密
plaintext, err := encryptor.Decrypt(ciphertext string) (string, error)
```

### 工具函数

```go
// 生成随机密钥
key, err := crypto.GenerateKey(keySize int) ([]byte, error)

// 密钥与 Base64 转换
keyStr := crypto.KeyToBase64(key []byte) string
key, err := crypto.Base64ToKey(keyStr string) ([]byte, error)
```

## 测试

项目包含了全面的单元测试：

```bash
go test ./backend/pkg/crypto -v
```

测试覆盖：
- ✅ 密钥生成和验证
- ✅ 加密/解密功能
- ✅ 无效密钥处理
- ✅ 无效密文处理
- ✅ 密钥转换
- ✅ 长密码和特殊字符支持

## 常见问题

**Q: 为什么每次加密结果不同?**
A: 这是正常的。使用随机 nonce 确保了即使加密相同的明文，结果也不同，这提高了安全性。

**Q: 如果密钥丢失怎么办?**
A: 如果密钥丢失，加密的数据将无法解密。务必做好密钥备份。

**Q: 为什么解密失败时返回原值?**
A: 这是为了向后兼容。如果数据库中有旧的明文密码，系统仍能继续工作。

**Q: 性能如何?**
A: AES-GCM 加密速度很快，对性能影响微乎其微。
