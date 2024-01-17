package apps

import (
	"configuration-management/global"
	"configuration-management/internal/biz/apps"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type QueryAppListRequest struct {
	ID    string `form:"id"`
	Name  string `form:"name"`
	Page  int    `form:"page" binding:"min=1"`
	Limit int    `form:"limit" binding:"min=1,max=50"`
}

// QueryAppList 批量获取app信息
func (handler *Handler) QueryAppList(c *gin.Context) {
	userInfo := app.GetUserInfoFromContext(c)

	var req QueryAppListRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	if !userInfo.IsRoot() {
		currentUser, err := handler.UserService.GetUserByID(userInfo.UserId)
		if err != nil {
			global.Logger.WithFields(logger.Fields{
				"userInfo": userInfo,
			}).Error("用户不存在", err)
			app.NewResponse(c).ToErrorResponse(errcode.NoPermission.WithDetails(err.Error()))
			return
		}
		app.NewResponse(c).ToResponseList(currentUser.Apps, len(currentUser.Apps))
		return
	}

	result, err := handler.AppService.QueryAppList(apps.QueryAppListArgs{
		ID:    req.ID,
		Name:  req.Name,
		Page:  req.Page,
		Limit: req.Limit,
	})
	if err != nil {
		global.Logger.Error("get apps failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ToResponseList(result.List, result.Total)
	return
}
