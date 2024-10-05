package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default address for Redis
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		redisDB = "0"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis at", redisAddr)

	return client
}

func Set(client *redis.Client, key string, value string, expiration time.Duration) error {
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Failed to set key: %s in Redis: %v", key, err)
		return err
	}
	return nil
}

func Get(client *redis.Client, key string) (string, error) {
	value, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key not found: %s", key)
		return "", nil
	} else if err != nil {
		log.Printf("Failed to get key: %s from Redis: %v", key, err)
		return "", err
	}
	return value, nil
}

func Del(client *redis.Client, key string) error {
	err := client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Failed to delete key: %s from Redis: %v", key, err)
		return err
	}
	return nil
}
