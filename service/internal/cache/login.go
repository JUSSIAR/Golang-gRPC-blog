package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	LOGIN_PREFIX               string        = "LOGIN_"
	INVALIDATION_SESSION_DELAY time.Duration = time.Hour * 2
)

func Login(ctx context.Context, redisDB redis.Client, user string) {
	key := LOGIN_PREFIX + user
	panicIfErr(redisDB.Set(ctx, key, true, INVALIDATION_SESSION_DELAY).Err())
}

func Logout(ctx context.Context, redisDB redis.Client, user string) {
	key := LOGIN_PREFIX + user
	_, err := redisDB.Del(ctx, key).Result()
	panicIfErr(err)
}

func CheckLogin(ctx context.Context, redisDB redis.Client, user string) bool {
	key := LOGIN_PREFIX + user
	_, err := redisDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	panicIfErr(err)

	return true
}
