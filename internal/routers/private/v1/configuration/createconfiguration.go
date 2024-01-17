package configuration

import (
	"configuration-management/global"
	"configuration-management/internal/biz/userconfig"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"github.com/gin-gonic/gin"
)

// CreateConfigurationRequest gin binding for body
type CreateConfigurationRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value" binding:"required"`
}

// CreateConfiguration
// @Summary 获取用户的配置信息
// @Produce  json
// @Param user_id query string true "用户id"
// @Param config_key query int true "用户配置key"
// @Failure 10000001 {object} errcode.Error "请求错误"
// @Failure 10000002 {object} errcode.Error "找不到"
// @Router /private/v1/configuration [post]
func (handler *Handler) CreateConfiguration(c *gin.Context) {
	var request CreateConfigurationRequest
	if err := c.ShouldBind(&request); err != nil {
		global.Logger.Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	err := handler.UserConfigService.CreateUserConfig(userconfig.CreateUserConfigArgs{
		UserID:      request.UserID,
		ConfigKey:   request.ConfigKey,
		ConfigValue: request.ConfigValue,
	})
	if err != nil {
		global.Logger.Error("create user config failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ResponseOK()
}
