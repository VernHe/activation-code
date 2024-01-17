package user

import (
	"net/http"

	"configuration-management/global"
	"configuration-management/internal/biz/user"
	"configuration-management/pkg/app"
	"configuration-management/pkg/errcode"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username    string `json:"username"`     // 用户名
	PasswordMD5 string `json:"password_md5"` // 密码的md5值
}

func (handler *Handler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBind(&request); err != nil {
		global.Logger.WithFields(logger.Fields{
			"request": request,
		}).Error("invalid params", err)
		app.NewResponse(c).ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	targetUser, err := handler.UserService.Login(user.LoginArgs{
		Username:    request.Username,
		PasswordMD5: request.PasswordMD5,
	})
	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"request": request,
		}).Error("login failed", err)
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails(err.Error()))
		return
	}

	if targetUser.ID == "" {
		global.Logger.WithFields(logger.Fields{
			"request": request,
		}).Error("login failed, targetUser id is empty")
		app.NewResponse(c).ToErrorResponse(errcode.NotFound.WithDetails("login failed"))
		return
	}

	// 生成 token
	token, err := app.CreateToken(app.UserInfo{
		UserId:   targetUser.ID,
		Username: targetUser.Username,
		MaxCnt:   targetUser.MaxCnt,
		Roles:    targetUser.Roles,
	})
	if err != nil {
		return
	}

	responseContent := app.ResponseContent{
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"token": token},
	}

	app.NewResponse(c).ToResponse(responseContent)
	return
}
