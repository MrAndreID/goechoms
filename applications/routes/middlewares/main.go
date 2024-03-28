package middlewares

import (
	"github.com/MrAndreID/goechoms/applications"
	"github.com/MrAndreID/goechoms/configs"
)

type CustomMiddleware struct {
	Config      *configs.Config
	Application *applications.Application
	RouteList   map[string]map[string]string
}

func NewCustomMiddleware(cfg *configs.Config, app *applications.Application) *CustomMiddleware {
	return &CustomMiddleware{
		Config:      cfg,
		Application: app,
	}
}
