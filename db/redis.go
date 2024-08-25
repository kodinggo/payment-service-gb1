package db

import (
	"github.com/tubagusmf/payment-service-gb1/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisHost(),
		Password: "",
		DB:       config.GetRedisDB(),
	})
	return rd
}
