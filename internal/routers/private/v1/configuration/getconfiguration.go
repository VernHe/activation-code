package configuration

import (
	"configuration-management/global"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetConfigurationRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	ConfigKey string `json:"config_key" binding:"required"`
}

// GetConfiguration
// @Summary 获取用户的配置信息
// @Produce  json
// @Param user_id query string true "用户id"
// @Param config_key query int true "用户配置key"
// @Failure 10000001 {object} errcode.Error "请求错误"
// @Failure 10000002 {object} errcode.Error "找不到"
// @Router /private/v1/configuration [get]
func (handler *Handler) GetConfiguration(c *gin.Context) {
	var request GetConfigurationRequest
	request.UserID = c.Query("user_id")
	request.ConfigKey = c.Query("config_key")
	if request.UserID == "" || request.ConfigKey == "" {
		global.Logger.Error("invalid params", request.UserID, request.ConfigKey)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(request.UserID, request.ConfigKey))
		return
	}

	userConfig, err := handler.UserConfigService.GetUserConfigByUserIDAndConfigKey(request.UserID, request.ConfigKey)
	if err != nil {
		global.Logger.Error("get user config failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error(), request.UserID, request.ConfigKey))
		return
	}

	app.NewResponse(c).ToResponse(app.ResponseContent{
		StatusCode: http.StatusOK,
		Data:       userConfig,
	})
	return
}
