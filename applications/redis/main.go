package redis

import (
	"context"

	redisPackage "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	Host     string
	Port     string
	Username string
	Password string
}

func New(redis *Redis) (*redisPackage.Client, error) {
	redisClient := redisPackage.NewClient(&redisPackage.Options{
		Addr:     redis.Host + ":" + redis.Port,
		Username: redis.Username,
		Password: redis.Password,
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   "Applications.Redis.Main.New.01",
			"error": err.Error(),
		}).Error("failed to connect redis")

		return nil, err
	}

	return redisClient, nil
}
