package applications

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}

	CustomHTTPErrorResponse struct {
		Code        string      `json:"code"`
		Description string      `json:"description"`
		Data        interface{} `json:"data"`
		Internal    error       `json:"-"`
	}
)

func (app *Application) NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func (app *Application) NewCustomHTTPErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)

	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")))
	}

	logrus.WithFields(logrus.Fields{
		"tag": "Applications.Validator.NewCustomHTTPErrorHandler.01",
	}).Error(strings.ToLower(cast.ToString(report.Message)))

	c.Logger().Error(report)

	c.JSON(report.Code, CustomHTTPErrorResponse{
		Internal:    report.Internal,
		Code:        fmt.Sprintf("%04d", report.Code),
		Description: strings.ToUpper(strings.ReplaceAll(cast.ToString(report.Message), " ", "_")),
	})
}
