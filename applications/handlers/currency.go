package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MrAndreID/goechoms/applications"
	"github.com/MrAndreID/goechoms/applications/types"
	"github.com/MrAndreID/goechoms/configs"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type CurrencyHandler struct {
	Config      *configs.Config
	Application *applications.Application
}

func NewCurrencyHandler(cfg *configs.Config, app *applications.Application) *CurrencyHandler {
	return &CurrencyHandler{
		Config:      cfg,
		Application: app,
	}
}

func (ch *CurrencyHandler) Index(c echo.Context) error {
	var (
		currencyResponse types.HTTPResponse
		tag              string = "Applications.Handlers.Currency.Index."
	)

	ch.Application.Service.Currency.Index(&currencyResponse)

	if currencyResponse.Error != nil || currencyResponse.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{
			"tag":        tag + "01",
			"error":      currencyResponse.Error,
			"statusCode": currencyResponse.StatusCode,
		}).Error("failed to hit currency url (index)")

		return c.JSON(currencyResponse.StatusCode, types.MainResponse{
			Code:        fmt.Sprintf("%04d", currencyResponse.StatusCode),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(currencyResponse.StatusCode), " ", "_")),
		})
	}

	var responseBody map[string]string

	err := json.Unmarshal([]byte(cast.ToString(currencyResponse.Body)), &responseBody)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to json unmarshal (body from currency response)")

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	var data []interface{}

	for i, v := range responseBody {
		data = append(data, map[string]string{
			"code": i,
			"name": v,
		})
	}

	return c.JSON(http.StatusOK, types.MainResponse{
		Code:        fmt.Sprintf("%04d", http.StatusOK),
		Description: "SUCCESS",
		Data:        data,
	})
}
