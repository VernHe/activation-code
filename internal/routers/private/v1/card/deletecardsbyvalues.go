package card

import (
	"configuration-management/global"
	"configuration-management/internal/biz/permissions"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type DeleteCardsByValuesRequest struct {
	Values []string `json:"values" binding:"required" form:"values"`
}

// DeleteCardsByValues 批量删除激活码
func (handler *Handler) DeleteCardsByValues(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)
	var req DeleteCardsByValuesRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.WithFields(logger.Fields{
			"error": err,
		}).Error("参数错误")
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
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
	}

	err := handler.CardService.DeleteCardsByValues(req.Values, userId)
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
