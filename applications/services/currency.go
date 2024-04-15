package services

import (
	"net/http"
	"time"

	"github.com/MrAndreID/goechoms/applications/types"
	"github.com/MrAndreID/goechoms/configs"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type CurrencyService struct {
	Config *configs.Config
}

func NewCurrencyService(cfg *configs.Config) *CurrencyService {
	return &CurrencyService{
		Config: cfg,
	}
}

func (cs *CurrencyService) Index(httpResponse *types.HTTPResponse) {
	var (
		restyClient *resty.Client = resty.New()
		tag         string        = "Applications.Services.Currency.Index."
	)

	restyClient.SetTimeout(time.Second * time.Duration(cs.Config.DefaultTimeout))

	url := cs.Config.CurrencyURL + "/currencies"

	restyResponse, err := restyClient.R().Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":           tag + "01",
			"url":           url,
			"requestHeader": restyResponse.Request.Header,
			"requestBody":   nil,
			"error":         err.Error(),
		}).Error("failed hit to currency url")

		httpResponse.StatusCode = http.StatusServiceUnavailable
		httpResponse.Error = err

		return
	}

	httpResponse.Headers = restyResponse.Header()
	httpResponse.Body = string(restyResponse.Body())
	httpResponse.StatusCode = restyResponse.StatusCode()

	logrus.WithFields(logrus.Fields{
		"tag":                tag + "02",
		"url":                url,
		"requestHeader":      restyResponse.Request.Header,
		"requestBody":        nil,
		"responseHeader":     restyResponse.Header(),
		"responseBody":       string(restyResponse.Body()),
		"responseStatusCode": restyResponse.StatusCode(),
	}).Info("result from hit to currency url")
}
