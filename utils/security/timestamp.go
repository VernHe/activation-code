package security

import (
	"time"

	"configuration-management/global"
	"configuration-management/pkg/logger"
)

// IsValidTimestamp 验证时间戳是否有效
func IsValidTimestamp(timestamp string) bool {
	clientTimestamp, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return false
	}

	// 假设有效期为 30 秒
	validityDuration := 5 * time.Minute
	currentTime := time.Now().UTC()

	global.Logger.WithFields(logger.Fields{
		"client_timestamp": clientTimestamp,
		"timestamp":        timestamp,
	}).Debug("校验时间")

	return clientTimestamp.After(currentTime.Add(-validityDuration)) && clientTimestamp.Before(currentTime.Add(validityDuration))
}
