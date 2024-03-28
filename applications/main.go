package applications

import (
	"time"

	"github.com/MrAndreID/goechoms/applications/databases"
	"github.com/MrAndreID/goechoms/applications/redis"
	"github.com/MrAndreID/goechoms/applications/services"
	"github.com/MrAndreID/goechoms/configs"

	redisPackage "github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Application struct {
	TimeLocation *time.Location
	Database     *gorm.DB
	Redis        *redisPackage.Client
	Service      *services.Service
}

func New(cfg *configs.Config) (*Application, error) {
	var (
		tag          string = "Applications.Main.New."
		timeLocation *time.Location
	)

	timeLocation, err := time.LoadLocation(cfg.TimeZone)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to load location for time")

		return nil, err
	}

	var databaseConnection *gorm.DB

	if cfg.UseDatabase {
		databaseConnection, err = databases.New(&databases.Database{
			Connection: cfg.DatabaseConnection,
			Host:       cfg.DatabaseHost,
			Port:       cfg.DatabasePort,
			Username:   cfg.DatabaseUsername,
			Password:   cfg.DatabasePassword,
			Name:       cfg.DatabaseName,
			SSLMode:    cfg.DatabaseSSLMode,
			ParseTime:  cfg.DatabaseParseTime,
			Charset:    cfg.DatabaseCharset,
			Timezone:   cfg.DatabaseTimezone,
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "02",
				"error": err.Error(),
			}).Error("failed to connect database")

			return nil, err
		}
	}

	var redisConnection *redisPackage.Client

	if cfg.UseRedis {
		redisConnection, err = redis.New(&redis.Redis{
			Host:     cfg.RedisHost,
			Port:     cfg.RedisPort,
			Username: cfg.RedisUsername,
			Password: cfg.RedisPassword,
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "03",
				"error": err.Error(),
			}).Error("failed to connect redis")

			return nil, err
		}
	}

	return &Application{
		TimeLocation: timeLocation,
		Database:     databaseConnection,
		Redis:        redisConnection,
		Service:      services.New(cfg, redisConnection, databaseConnection),
	}, nil
}

func (app *Application) Start(cfg *configs.Config, e *echo.Echo) error {
	return e.Start(":" + cfg.Port)
}
