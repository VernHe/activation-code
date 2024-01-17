package userconfig

import "time"

// UserConfig 结构体用于映射数据库中的 user_config 表
type UserConfig struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ConfigKey   string    `json:"config_key"`
	ConfigValue string    `json:"config_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsDeleted   bool      `json:"is_deleted"`
}

// TableName 指定 UserConfig 结构体对应的表名
func (UserConfig) TableName() string {
	return "user_config"
}
