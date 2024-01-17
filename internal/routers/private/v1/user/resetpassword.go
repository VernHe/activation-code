package user

import (
	"configuration-management/global"
	"configuration-management/internal/biz/user"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ResetPasswordRequest struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler *Handler) ResetPassword(c *gin.Context) {
	var request ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		global.Logger.WithFields(logger.Fields{
			"request": request,
		}).Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if err := handler.UserService.ResetPassword(user.ResetPasswordArgs{
		ID:       request.ID,
		Password: request.Password,
	}); err != nil {
		global.Logger.WithFields(logger.Fields{
			"request": request,
			"error":   err.Error(),
		}).Error("重置密码失败", err)
		app.NewResponse(c).ToErrorResponse(errcode.ServerError.WithDetails("重置密码失败"))
		return
	}

	app.NewResponse(c).ResponseOK()
	return
}
