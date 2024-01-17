package apps

import (
	"configuration-management/internal/biz/apps"
	"configuration-management/internal/biz/user"
)

type Handler struct {
	AppService  apps.Service
	UserService user.Service
}

func NewHandler() *Handler {
	return &Handler{
		AppService:  apps.NewService(),
		UserService: user.NewService(),
	}
}
