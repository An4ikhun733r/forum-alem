package handlers

import (
	"forum/internal/app"
	"forum/internal/service"
)

type HandlerApp struct {
	service service.ServiceI
	*app.Application
}

func New(s service.ServiceI, a *app.Application) *HandlerApp {
	return &HandlerApp{
		s,
		a,
	}
}
