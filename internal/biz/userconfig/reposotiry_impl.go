package userconfig

import (
	"errors"

	"configuration-management/global"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"gorm.io/gorm"
)

/*
表结构如下：
CREATE TABLE user_config (
   id INT AUTO_INCREMENT PRIMARY KEY,
   user_id INT NOT NULL,
   config_key VARCHAR(255) NOT NULL,
   config_value TEXT,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   is_deleted TINYINT DEFAULT 0, -- 默认值为0表示未删除
   UNIQUE KEY (user_id, config_key)
);
*/

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserConfigByUserIDAndConfigKey(userID string, configKey string) (*UserConfig, error) {
	var userConfig UserConfig
	err := r.db.Where("user_id = ? AND config_key = ?", userID, configKey).First(&userConfig).Error
	if err != nil {
		// 识别是否为记录不存在的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"user_id":    userID,
				"config_key": configKey,
			}).Error("查询时记录不存在", err)
			return nil, errcode.NotFound
		}
		// 其他错误
		return nil, err
	}
	return &userConfig, nil
}

func (r *repository) CreateUserConfig(userConfig *UserConfig) error {
	err := r.db.Create(userConfig).Error
	if err != nil {
		// 识别是否为记录已存在的错误
		if r.db.Error != nil && r.db.Error.Error() == "Error 1062: Duplicate entry '"+userConfig.UserID+"-"+userConfig.ConfigKey+"' for key 'user_id'" {
			global.Logger.WithFields(logger.Fields{
				"user_id":    userConfig.UserID,
				"config_key": userConfig.ConfigKey,
			}).Error("创建时记录已存在", err)
			return errcode.DuplicateKey
		}
		// 其他错误
		return err
	}
	return nil
}

func (r *repository) UpdateUserConfig(userConfig *UserConfig) error {
	err := r.db.Model(userConfig).Where("user_id = ? AND config_key = ?", userConfig.UserID, userConfig.ConfigKey).Updates(userConfig).Error
	if err != nil {
		// 记录不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.WithFields(logger.Fields{
				"user_id":    userConfig.UserID,
				"config_key": userConfig.ConfigKey,
			}).Error("更新时记录不存在", err)
			return errcode.NotFound
		}
		// 其他错误
		return err
	}
	return nil
}

func (r *repository) DeleteUserConfig(userConfig *UserConfig) error {
	return r.db.Where("user_id = ? AND config_key = ?", userConfig.UserID, userConfig.ConfigKey).Delete(userConfig).Error
}
