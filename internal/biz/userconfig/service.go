package userconfig

type CreateUserConfigArgs struct {
	UserID      string `json:"user_id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
}

type Service interface {
	GetUserConfigByUserIDAndConfigKey(userID string, configKey string) (*UserConfig, error)
	CreateUserConfig(args CreateUserConfigArgs) error
	UpdateUserConfig(userConfig *UserConfig) error
	DeleteUserConfig(userConfig *UserConfig) error
}
