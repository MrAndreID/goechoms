package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MrAndreID/goechoms/applications"
	"github.com/MrAndreID/goechoms/applications/handlers"
	"github.com/MrAndreID/goechoms/applications/routes/middlewares"
	"github.com/MrAndreID/goechoms/applications/types"
	"github.com/MrAndreID/goechoms/configs"

	loggerUtil "github.com/hlmn/senyum-go-utils/logger/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"go.elastic.co/apm/module/apmechov4"
)

func New(cfg *configs.Config, app *applications.Application) *echo.Echo {
	var tag string = "Applications.Routes.Main.New."

	echo.NotFoundHandler = func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"tag": tag + "01",
		}).Error("route not found")

		return c.JSON(http.StatusNotFound, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusNotFound),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusNotFound), " ", "_")),
		})
	}

	echo.MethodNotAllowedHandler = func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"tag": tag + "02",
		}).Error("method not allowed")

		return c.JSON(http.StatusMethodNotAllowed, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusMethodNotAllowed),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusMethodNotAllowed), " ", "_")),
		})
	}

	e := echo.New()

	e.Validator = app.NewCustomValidator()

	e.Logger.SetLevel(log.DEBUG)

	e.Debug = true

	e.HTTPErrorHandler = app.NewCustomHTTPErrorHandler

	e.JSONSerializer = app.NewCustomJSON()

	e.Pre(middleware.RemoveTrailingSlash())

	middlewares := middlewares.NewCustomMiddleware(cfg, app)

	e.Pre(middlewares.SetRequestID)

	e.Use(apmechov4.Middleware())

	e.Use(middleware.Recover())

	e.Use(middleware.BodyDump(func(c echo.Context, requestBody, responseBody []byte) {
		request := struct {
			Header interface{} `json:"header"`
			Body   string      `json:"body"`
		}{
			Header: c.Request().Header,
			Body:   string(requestBody),
		}

		response := struct {
			Header interface{} `json:"header"`
			Body   string      `json:"body"`
		}{
			Header: c.Response().Header(),
			Body:   string(responseBody),
		}

		loggerUtil.Info(c, logrus.Fields{
			"request":   request,
			"requestId": c.Get("RequestID"),
			"response":  response,
			"url":       c.Request().Host + c.Request().URL.String(),
		}, "body dump")
	}))

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	secureMiddleware := secure.Options{
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           63072000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		ForceSTSHeader:       true,
		IsDevelopment:        true,
	}

	e.Use(echo.WrapMiddleware(secure.New(secureMiddleware).Handler))

	e.Use(middleware.Logger())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	if cfg.UseSignature {
		e.Use(middlewares.SignatureCheck)
	}

	e.Use(middlewares.NoCache)

	handler := handlers.New(cfg, app)

	v1 := e.Group("/api/v1")

	userRoute := v1.Group("/user")
	userRoute.GET("", handler.User.Index, middlewares.ServiceKeyCheck).Name = "user.index"
	userRoute.POST("", handler.User.Create, middlewares.ServiceKeyCheck).Name = "user.create"
	userRoute.PUT("/:id", handler.User.Edit, middlewares.ServiceKeyCheck).Name = "user.edit"
	userRoute.DELETE("/:id", handler.User.Delete, middlewares.ServiceKeyCheck).Name = "user.delete"

	routes := e.Routes()

	middlewares.RouteList = middlewares.SetRouteList(routes)

	return e
}
