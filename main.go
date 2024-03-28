package main

import (
	"github.com/MrAndreID/goechoms/applications"
	"github.com/MrAndreID/goechoms/applications/routes"
	"github.com/MrAndreID/goechoms/configs"

	"github.com/sirupsen/logrus"
)

func main() {
	var tag string = "Main."

	cfg, err := configs.New()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to initiate configuration")

		return
	}

	app, err := applications.New(cfg)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to initiate application")

		return
	}

	e := routes.New(cfg, app)

	if app.Start(cfg, e); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to run application")

		return
	}
}
