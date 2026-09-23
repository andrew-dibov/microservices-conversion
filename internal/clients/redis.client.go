package clients

import (
	"microservices-conversion/internal/configs"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(appConfig *configs.AppConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: appConfig.RedisClient.Address,

		DB:       appConfig.RedisClient.Database,
		Password: appConfig.RedisClient.Password,

		MaxRetries:      appConfig.RedisClient.Retries.MaxRetries,
		MinRetryBackoff: appConfig.RedisClient.Retries.MinRetryBackoff,
		MaxRetryBackoff: appConfig.RedisClient.Retries.MaxRetryBackoff,
	})

	return &RedisClient{
		client: client,
		ttl:    appConfig.RedisClient.TTL,
	}
}
