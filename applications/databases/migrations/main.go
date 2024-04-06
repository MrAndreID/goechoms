package main

import (
	"flag"
	"fmt"

	"github.com/MrAndreID/goechoms/applications"
	"github.com/MrAndreID/goechoms/applications/databases/models"
	"github.com/MrAndreID/goechoms/configs"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var tables map[string]interface{} = map[string]interface{}{
	"users":  &models.User{},
	"emails": &models.Email{},
}

func main() {
	var tag string = "Applications.Databases.Migrations.Main.Main."

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

	migrateFlag := flag.String("migrate", "default", "For Migrate")

	flag.Parse()

	if cast.ToString(migrateFlag) == "fresh" {
		fmt.Println("Start Drop All Tables")

		existingTables, err := app.Database.Migrator().GetTables()

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "04",
				"error": err.Error(),
			}).Error("failed to get tables from database")

			return
		}

		for _, v := range existingTables {
			fmt.Println("Dropping: " + v + " Table")

			err := app.Database.Migrator().DropTable(v)

			if err != nil {
				logrus.WithFields(logrus.Fields{
					"tag":   tag + "05",
					"error": err.Error(),
				}).Error("failed to drop table")

				return
			}

			fmt.Println("Dropped: " + v + " Table")
		}

		fmt.Println("End Drop All Tables")
	}

	fmt.Println()

	fmt.Println("Start Migration")

	err = Migrate(app)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "06",
			"error": err.Error(),
		}).Error("failed to migrate")

		return
	}

	fmt.Println("End Migration")
}

func Migrate(app *applications.Application) error {
	for i, v := range tables {
		fmt.Println("Migrating: " + i + " Table")

		err := app.Database.Migrator().CreateTable(v)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   "Applications.Databases.Migrations.Main.Migrate.01",
				"error": err.Error(),
			}).Error("failed to create table")

			return err
		}

		fmt.Println("Migrated: " + i + " Table")
	}

	return nil
}
