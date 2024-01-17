package biz

import (
	"fmt"

	"configuration-management/global"
	"configuration-management/pkg/logger"
)

// IsValidTimeType 检查时间类型是否有效
func IsValidTimeType(timeType string, days int, minutes int) bool {
	// 定义不同时间类型的有效时长范围（分钟）
	timeLimits := map[string]struct {
		Min int
		Max int
	}{
		"hourly":  {Min: 1, Max: 24 * 60},
		"daily":   {Min: 24 * 60, Max: 7 * 24 * 60},
		"weekly":  {Min: 7 * 24 * 60, Max: 30 * 24 * 60},
		"monthly": {Min: 30 * 24 * 60, Max: 365 * 24 * 60},
		"yearly":  {Min: 365 * 24 * 60, Max: 999 * 24 * 60}, // 999天的上限，可以根据实际需求调整
	}

	// 检查输入参数范围
	if days < 0 || minutes < 0 || (days == 0 && minutes == 0) {
		global.Logger.WithFields(logger.Fields{
			"days":    days,
			"minutes": minutes,
		}).Error("错误的时间参数")
		return false
	}

	// 计算总有效时间（分钟）
	totalMinutes := days*24*60 + minutes

	// 检查是否在有效时长范围内
	limit, ok := timeLimits[timeType]
	if !ok {
		fmt.Println("Invalid time type")
		return false
	}

	if totalMinutes < limit.Min || totalMinutes >= limit.Max {
		return false
	}

	return true
}
