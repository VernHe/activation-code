package configuration

import (
	"configuration-management/internal/biz/userconfig"
)

type Handler struct {
	UserConfigService userconfig.Service
}

func NewHandler() *Handler {
	return &Handler{
		UserConfigService: userconfig.NewService(),
	}
}
