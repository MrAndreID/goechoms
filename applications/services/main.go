package services

import (
	"github.com/MrAndreID/goechoms/configs"

	redisPackage "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Service struct {
	//
}

func New(cfg *configs.Config, redisConnection *redisPackage.Client, databaseConnection *gorm.DB) *Service {
	return &Service{
		//
	}
}
