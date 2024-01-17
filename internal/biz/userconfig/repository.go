package userconfig

type Repository interface {
	GetUserConfigByUserIDAndConfigKey(userID string, configKey string) (*UserConfig, error)
	CreateUserConfig(userConfig *UserConfig) error
	UpdateUserConfig(userConfig *UserConfig) error
	DeleteUserConfig(userConfig *UserConfig) error
}
