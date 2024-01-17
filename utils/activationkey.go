package utils

import (
	"crypto/rand"
)

const (
	activationKeyLength = 16
	charset             = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// 生成随机字符串
func generateRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		randomBytes[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(randomBytes), nil
}

func GenerateActivationKey() string {
	randomString, _ := generateRandomString(activationKeyLength)
	return randomString
}

func GenerateActivationKeyByApp(prefix string, length int) string {
	randomString, _ := generateRandomString(length)
	return prefix + randomString
}
