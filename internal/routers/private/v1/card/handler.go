package card

import (
	"configuration-management/internal/biz/activationattempt"
	"configuration-management/internal/biz/apps"
	"configuration-management/internal/biz/card"
	"configuration-management/internal/biz/user"
)

type Handler struct {
	CardService       card.Service
	AppService        apps.Service
	UserService       user.Service
	ActivationAttempt activationattempt.Service
}

func NewHandler() *Handler {
	return &Handler{
		CardService:       card.NewService(),
		AppService:        apps.NewService(),
		UserService:       user.NewService(),
		ActivationAttempt: activationattempt.NewService(),
	}
}
