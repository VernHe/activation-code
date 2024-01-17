package user

import (
	"configuration-management/global"
	"configuration-management/internal/biz/user"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type QueryUserListRequest struct {
	Username string      `form:"username"`
	Status   user.Status `form:"status"`
	Page     int         `form:"page"`
	Limit    int         `form:"limit"`
}

func (handler *Handler) QueryUserList(c *gin.Context) {

	var req QueryUserListRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.WithFields(logger.Fields{
			"req": req,
		}).Error("create user failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	result, err := handler.UserService.QueryUserList(user.QueryUserListArgs{
		Username: req.Username,
		Status:   req.Status,
		Page:     req.Page,
		Limit:    req.Limit,
	})
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"req": req,
		}).Error("query user failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	app.NewResponse(c).ToResponseList(user.BatchToView(result.List), result.Total)
}
