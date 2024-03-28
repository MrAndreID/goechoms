package handlers

import (
	"github.com/MrAndreID/goechoms/applications"
	"github.com/MrAndreID/goechoms/configs"
)

type Handler struct {
	User *UserHandler
}

func New(cfg *configs.Config, app *applications.Application) *Handler {
	return &Handler{
		User: NewUserHandler(cfg, app),
	}
}
