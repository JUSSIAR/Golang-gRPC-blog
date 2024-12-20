package storages

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis(ctx context.Context) redis.Client {
	redisURL, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		redisURL = "localhost:6379"
	}

	redisDB := redis.NewClient(&redis.Options{Addr: redisURL})

	_, err := redisDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return *redisDB
}
