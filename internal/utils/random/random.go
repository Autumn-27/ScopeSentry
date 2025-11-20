package random

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	mRand "math/rand"
)

// GenerateString 生成指定长度的随机字符串
func GenerateString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// GenerateSecretKey 生成16位密钥
func GenerateSecretKey() string {
	return GenerateString(16)
}

// GeneratePluginKey 生成6位插件密钥
func GeneratePluginKey() string {
	return GenerateString(32)
}

func GenerateRandomString(length int) string {
	// 定义字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 构建随机字符串
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[mRand.Intn(len(charset))]
	}
	return string(result)
}

func GeneratePassword(length int) (string, error) {
	// 安全字符集：大小写字母 + 数字 + 特殊字符
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"!@#$%^&*()-_=+[]{}<>?/"

	pass := make([]byte, length)

	for i := 0; i < length; i++ {
		// 从 charset 中安全随机取一个字符
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		pass[i] = charset[n.Int64()]
	}

	return string(pass), nil
}
