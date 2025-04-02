package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // TODO: Use config
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	log.Println("✅ Redis connected successfully!")
}

func SetCache(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func GetCache(key string) (string, error) {
	ctx := context.Background()
	return RedisClient.Get(ctx, key).Result()
}

func DeleteCache(key string) error {
	ctx := context.Background()
	return RedisClient.Del(ctx, key).Err()
}
