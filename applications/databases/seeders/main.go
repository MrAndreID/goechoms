package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MrAndreID/goechoms/applications"
	"github.com/MrAndreID/goechoms/applications/databases/models"
	"github.com/MrAndreID/goechoms/configs"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var seeder map[string]map[string]interface{} = map[string]map[string]interface{}{
	"user1": {
		"model": &models.User{},
		"data": &models.User{
			ID:        "aaaaaaaa-1111-aaaa-1111-aaaaaaaaaaaa",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "Andrea Adam",
		},
	},
	"email1": {
		"model": &models.Email{},
		"data": &models.Email{
			ID:        "aaaaaaaa-1111-aaaa-1111-aaaaaaaaaaaa",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    "aaaaaaaa-1111-aaaa-1111-aaaaaaaaaaaa",
			Email:     "mrandreid@email.com",
		},
	},
	"email2": {
		"model": &models.Email{},
		"data": &models.Email{
			ID:        "bbbbbbbb-2222-bbbb-2222-bbbbbbbbbbbb",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    "aaaaaaaa-1111-aaaa-1111-aaaaaaaaaaaa",
			Email:     "mrandreid@email.co.id",
		},
	},
}

func main() {
	var tag string = "Applications.Databases.Seeders.Main.Main."

	cfg, err := configs.New()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to initiate configuration")

		return
	}

	if !cfg.UseDatabase {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": "The Database is not yet used",
		}).Error("failed to migrate")

		return
	}

	app, err := applications.New(cfg)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to initiate application")

		return
	}

	fmt.Println("Start Seeder")

	seedFlag := flag.String("seed", "default", "For Seed")

	flag.Parse()

	if cast.ToString(seedFlag) == "default" {
		fmt.Println("Start Seed")

		for i, v := range seeder {
			fmt.Println("Seeding: " + i + " Data")

			for key, data := range v {

				if key == "model" {
					if !app.Database.Migrator().HasTable(data) {
						logrus.WithFields(logrus.Fields{
							"tag":   tag + "04",
							"error": "Failed to Initiate Table",
						}).Error("failed to initiate table")

						return
					}
				}

				if key == "data" {
					result := app.Database.Create(data)

					if result.Error != nil {
						logrus.WithFields(logrus.Fields{
							"tag":   tag + "05",
							"error": result.Error.Error(),
						}).Error("failed to create data")

						return
					}

					if result.RowsAffected == 0 {
						logrus.WithFields(logrus.Fields{
							"tag":   tag + "06",
							"error": "failed to create data",
						}).Error("failed to create data")

						return
					}
				}
			}

			fmt.Println("Seeded: " + i + " Data")
		}

		fmt.Println("End Seed")
	}

	fmt.Println("End Seeder")
}
