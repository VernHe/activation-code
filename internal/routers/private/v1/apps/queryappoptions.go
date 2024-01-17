package apps

import (
	"net/http"

	"configuration-management/global"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

// QueryAppOptions 批量获取app信息
func (handler *Handler) QueryAppOptions(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	if !userInfo.IsRoot() {
		currentUser, err := handler.UserService.GetUserByID(userInfo.UserId)
		if err != nil {
			global.Logger.WithFields(logger.Fields{
				"userInfo": userInfo,
			}).Error("用户不存在", err)
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission.WithDetails(err.Error()))
			return
		}
		apps, err := handler.AppService.GetAppByIDs(currentUser.Apps)
		if err != nil {
			global.Logger.Error("获取 app 信息失败", err)
			app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
			return
		}
		app.NewResponse(c).ToResponse(app.ResponseContent{
			StatusCode: http.StatusOK,
			Data:       apps,
		})
		return
	}

	appOptions, err := handler.AppService.QueryAppOptions()
	if err != nil {
		global.Logger.Error("获取 app 信息失败", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ToResponse(app.ResponseContent{
		StatusCode: http.StatusOK,
		Data:       appOptions,
	})
	return
}
