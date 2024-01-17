package security

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// GenerateCipherText 生成暗号，如果一小时内尝试激活次数超过 3 次，或者累计激活超过 5 次，返回假暗号，否则返回真暗号
// 真暗号: 与激活码的 MD5 中的数字部分之和相同的数字字符串，比如激活码的 MD5 是 1eds23s4w567g89f01，数字部分之和是 1+2+3+4+5+6+7+8+9+0+1=46，那么生成的真暗号中的数字之和页是 46
// 假暗号: 与激活码的 MD5 中的数字部分之和不同的数字字符串，比如激活码的 MD5 是 1eds23s4w567g89f01，数字部分之和是 1+2+3+4+5+6+7+8+9+0+1=46，那么假暗号的数字之和就不是 46
func GenerateCipherText(cardValue string, isTrueCipher bool) string {
	md5Value := GenerateMD5(cardValue)
	sum := getNumberSum(md5Value)
	// 一小时内尝试激活次数超过 3 次，或者累计激活超过 5 次，返回假暗号
	if isTrueCipher {
		// 生成真暗号，告诉客户端进行相应的操作
		return generateString(sum)
	}

	return generateString(rand.Intn(sum))
}

// 生成一个包含数字的字符串，数字之和等于目标值
func generateString(targetSum int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成一个包含数字的字符串
	var result string
	currentSum := 0

	// 随机生成数字并追加到字符串中，直到和等于目标值
	for currentSum != targetSum {
		// 生成一个随机数字
		num := r.Intn(10)
		// 如果加上这个数字不超过目标值，将数字部分加到结果字符串
		if currentSum+num <= targetSum {
			result += fmt.Sprintf("%d", num)
			// 更新当前和
			currentSum += num
		}
	}

	return result
}

// getNumberSum 值计算字符串中数字的和
func getNumberSum(str string) int {
	sum := 0
	for _, c := range str {
		if c >= '0' && c <= '9' {
			sum += int(c - '0')
		}
	}

	return sum
}

// GenerateMD5 生成字符串的 MD5 哈希
func GenerateMD5(input string) string {
	// 创建 MD5 哈希对象
	hasher := md5.New()
	// 将字符串转换为字节数组并写入哈希对象
	hasher.Write([]byte(input))
	// 计算哈希值并将其转换为十六进制字符串
	hashInBytes := hasher.Sum(nil)
	md5String := hex.EncodeToString(hashInBytes)

	return md5String
}
