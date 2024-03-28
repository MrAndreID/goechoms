package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MrAndreID/goechoms/applications/types"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (cm *CustomMiddleware) ServiceKeyCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var tag string = "Applications.Routes.Middlewares.ServiceKey.ServiceKeyCheck."

		if c.Request().Header["Service-Key"] == nil || c.Request().Header["Service-Key"][0] == "" {
			logrus.WithFields(logrus.Fields{
				"tag": tag + "01",
			}).Error("failed to checking service key from header")

			return c.JSON(http.StatusBadRequest, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusBadRequest),
				Description: "INVALID_SERVICE_KEY",
			})
		}

		err := validation.Validate(c.Request().Header["Service-Key"][0], validation.Required, validation.By(types.BlacklistValidation("Service-Key")))

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "02",
				"error": err.Error(),
			}).Error("failed to checking service key from header")

			return c.JSON(http.StatusBadRequest, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusBadRequest),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
				Data: map[string]string{
					"Service-Key": err.Error(),
				},
			})
		}

		if c.Request().Header["Service-Key"][0] != cm.Config.ServiceKey {
			logrus.WithFields(logrus.Fields{
				"tag": tag + "03",
			}).Error("failed to checking service key from header")

			return c.JSON(http.StatusBadRequest, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusBadRequest),
				Description: "INVALID_SERVICE_KEY",
			})
		}

		return next(c)
	}
}
