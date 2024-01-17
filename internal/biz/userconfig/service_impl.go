package userconfig

import (
	"configuration-management/global"
	"configuration-management/utils"
	"time"
)

type service struct {
	repo Repository
}

func NewService() Service {
	return &service{repo: NewRepository(global.DBEngine)}
}

//type UserConfig struct {
//	ID          int       `json:"id"`
//	UserID      int       `json:"user_id"`
//	ConfigKey   string    `json:"config_key"`
//	ConfigValue string    `json:"config_value"`
//	CreatedAt   time.Time `json:"created_at"`
//	UpdatedAt   time.Time `json:"updated_at"`
//	IsDeleted   bool      `json:"is_deleted"`
//}

//CREATE TABLE user_config (
//    id VARCHAR(255) PRIMARY KEY,
//    user_id VARCHAR(255) NOT NULL,
//    config_key VARCHAR(255) NOT NULL,
//    config_value TEXT,
//    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//    is_deleted TINYINT DEFAULT 0,
//    UNIQUE KEY (user_id, config_key)
//);

func (s *service) GetUserConfigByUserIDAndConfigKey(userID string, configKey string) (*UserConfig, error) {
	userConfig, err := s.repo.GetUserConfigByUserIDAndConfigKey(userID, configKey)
	if err != nil {
		return nil, err
	}
	return userConfig, nil
}

func (s *service) CreateUserConfig(args CreateUserConfigArgs) error {
	userConfig := &UserConfig{
		ID:          utils.GenerateUUID(),
		UserID:      args.UserID,
		ConfigKey:   args.ConfigKey,
		ConfigValue: args.ConfigValue,
		CreatedAt:   time.Now(),
		IsDeleted:   false,
	}

	err := s.repo.CreateUserConfig(userConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateUserConfig(userConfig *UserConfig) error {
	err := s.repo.UpdateUserConfig(userConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteUserConfig(userConfig *UserConfig) error {
	err := s.repo.DeleteUserConfig(userConfig)
	if err != nil {
		return err
	}
	return nil
}
