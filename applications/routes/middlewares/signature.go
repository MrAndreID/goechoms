package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MrAndreID/goechoms/applications/types"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func (cm *CustomMiddleware) SignatureCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			tag                    string = "Applications.Routes.Middlewares.Signature.SignatureCheck."
			method                 string = c.Request().Method
			contentType            string = c.Request().Header.Get("Content-Type")
			signature              string = c.Request().Header.Get(cm.Config.SignatureName)
			signatureTransactionID string = c.Request().Header.Get(cm.Config.SignatureTransactionIDName)
			url, stringToSign      string
		)

		methodToSkip := map[string]struct{}{
			"CONNECT": {},
			"HEAD":    {},
			"OPTIONS": {},
			"TRACE":   {},
		}

		if _, ok := methodToSkip[c.Request().Method]; ok {
			return next(c)
		}

		c.Response().Header().Add(cm.Config.SignatureValidationName, "FALSE")

		if signature == "" || signatureTransactionID == "" {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "01",
				"error": "Failed to Check Signature & Signature Transaction ID",
			}).Error("failed to check signature & signature transaction id")

			return c.JSON(http.StatusUnauthorized, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusUnauthorized),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusUnauthorized), " ", "_")),
			})
		}

		if c.Request().URL.RawQuery != "" {
			url = "?" + c.Request().URL.RawQuery
		}

		bodyBytes, err := io.ReadAll(c.Request().Body)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "02",
				"error": err.Error(),
			}).Error("failed to read all from request body")

			return c.JSON(http.StatusUnauthorized, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusUnauthorized),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusUnauthorized), " ", "_")),
			})
		}

		c.Request().Body.Close()

		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		bodyString := string(bodyBytes)

		if bodyString == "" {
			bodyString = "{}"
		}

		transactionId, err := cm.Application.DecryptSignature(cm.Config, signatureTransactionID)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "03",
				"error": err.Error(),
			}).Error("failed to decrypt signature for transaction id")

			return c.JSON(http.StatusUnauthorized, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusUnauthorized),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusUnauthorized), " ", "_")),
			})
		}

		if transactionId == nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "04",
				"error": "Transaction ID is Null",
			}).Error("transaction id is null")

			return c.JSON(http.StatusUnauthorized, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusUnauthorized),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusUnauthorized), " ", "_")),
			})
		}

		if method == "GET" || contentType == "multipart/form-data" {
			requestBody := "{}"

			stringToSign = url + "|" + requestBody + "|" + cast.ToString(transactionId)
		} else {
			stringToSign = url + "|" + bodyString + "|" + cast.ToString(transactionId)
		}

		verifySignature, err := cm.Application.VerifySignature(cm.Config, stringToSign)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "05",
				"error": err.Error(),
			}).Error("failed to verify signature")

			return c.JSON(http.StatusUnauthorized, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusUnauthorized),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusUnauthorized), " ", "_")),
			})
		}

		if verifySignature == nil || cast.ToString(verifySignature) != signature {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "06",
				"error": "Verify Signature is Null or Verify Signature is Not Equal with Signature",
			}).Error("verify signature is null or verify signature is not equal with signature")

			return c.JSON(http.StatusUnauthorized, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusUnauthorized),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusUnauthorized), " ", "_")),
			})
		}

		c.Response().Header().Set(cm.Config.SignatureValidationName, "TRUE")

		return next(c)
	}
}
