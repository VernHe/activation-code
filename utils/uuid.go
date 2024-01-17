package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	// 生成一个新的 UUID
	newUUID := uuid.New()
	// 将 UUID 转换成字符串，去掉 -
	return strings.Replace(newUUID.String(), "-", "", -1)
}

// 输入 string 返回 MD5
func MD5(row string) string {
	hash := md5.Sum([]byte(row))
	return hex.EncodeToString(hash[:])
}
