package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (cm *CustomMiddleware) SetRequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		UUID, err := uuid.NewRandom()

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   "Applications.Routes.Middlewares.RequestID.SetRequestID.01",
				"error": err.Error(),
			}).Error("failed to generate uuid for request id")

			return err
		}

		c.Set("RequestID", &UUID)

		return next(c)
	}
}
