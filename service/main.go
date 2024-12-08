package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/JUSSIAR/Golang-gRPC-blog/service/internal/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectPostgres() *gorm.DB {
	creds := []string{
		"host=localhost",
		"dbname=grpc-blog-storage",
		"user=jussiar",
		"password=admin",
		"port=5432",
		"sslmode=disable",
		"TimeZone=UTC",
	}

	postgresDB, err := gorm.Open(
		postgres.Open(strings.Join(creds, " ")),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	err = postgresDB.AutoMigrate(&models.Post{}, &models.Comment{})
	if err != nil {
		panic(err)
	}

	return postgresDB
}

func connectRedis(ctx context.Context) redis.Client {
	redisDB := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	_, err := redisDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return *redisDB
}

func main() {
	ctx := context.Background()
	postgresDB := connectPostgres()
	redisDB := connectRedis(ctx)

	fmt.Println("OK")
	fmt.Println(postgresDB)
	fmt.Println(redisDB)

	
}
