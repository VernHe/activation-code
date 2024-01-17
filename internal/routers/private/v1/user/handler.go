package user

import (
	"configuration-management/internal/biz/user"
)

type Handler struct {
	UserService user.Service
}

func NewHandler() *Handler {
	return &Handler{
		UserService: user.NewService(),
	}
}
