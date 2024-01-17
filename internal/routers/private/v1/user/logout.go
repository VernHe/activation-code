package user

import (
	"time"

	"configuration-management/global"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type LogoutRequest struct {
	Token string `json:"token"`
}

func (handler *Handler) Logout(c *gin.Context) {
	var request LogoutRequest
	if err := c.ShouldBind(&request); err != nil {
		global.Logger.WithFields(logger.Fields{
			"request": request,
		}).Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 标记 token 为已失效
	global.InvalidTokenCache.Set(request.Token, nil, 5*time.Minute)

	app.NewResponse(c).ResponseOK()
	return
}
