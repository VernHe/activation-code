package card

import (
	"configuration-management/global"
	"configuration-management/internal/biz/permissions"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

// DeleteCard 创建激活码
func (handler *Handler) DeleteCard(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)
	value := c.Param("value")
	if value == "" {
		global.Logger.WithFields(logger.Fields{
			"value": value,
		}).Error("删除空的激活码")
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails("value", value))
		return
	}
	var userId string
	if !userInfo.IsRoot() {
		currentUser, err := handler.UserService.GetUserByID(userInfo.UserId)
		if err != nil {
			global.Logger.WithFields(logger.Fields{
				"userInfo": userInfo,
			}).Error("用户不存在", err)
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission.WithDetails(err.Error()))
			return
		}
		// 检查用户是否有查询权限
		if !currentUser.HasPermission(permissions.DELETE) {
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission)
			return
		}

		userId = userInfo.UserId
	}
	err := handler.CardService.DeleteCardByValue(value, userId)
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"error": err,
		}).Error("create card failed")
		app.NewResponse(c).ToErrorResponse(errcode.ServerError)
		return
	}

	app.NewResponse(c).ResponseOK()
	return
}
