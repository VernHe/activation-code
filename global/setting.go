package global

import (
	"crypto/rsa"
	"time"

	"configuration-management/pkg/logger"
	"configuration-management/pkg/setting"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

// 读取配置信息，然后保存到全局变量中
var (
	ServerSetting     *setting.ServerSettingS
	AppSetting        *setting.AppSettingS
	DatabaseSetting   *setting.DatabaseSettingS
	Logger            *logger.Logger
	DBEngine          *gorm.DB
	InvalidTokenCache *cache.Cache

	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey

	Key = "cc32digitkey12345678901234567890"
	IV  = "cc16digitIvKey12"
)

func init() {
	InvalidTokenCache = cache.New(5*time.Minute, 5*time.Minute)
}
