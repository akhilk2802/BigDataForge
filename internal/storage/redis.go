package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// NewRedisClient initializes and returns a Redis client
func NewRedisClient() *redis.Client {
	// Fetch Redis connection details from environment variables (if available)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default address for Redis
	}

	redisPassword := os.Getenv("REDIS_PASSWORD") // Optional, set in environment for production
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		redisDB = "0" // Default Redis DB
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // Redis address
		Password: redisPassword, // No password by default
		DB:       0,             // Default DB
	})

	// Check the connection by pinging the Redis server
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis at", redisAddr)

	return client
}

// SetKey stores a key-value pair in Redis with an optional expiration time
func Set(client *redis.Client, key string, value string, expiration time.Duration) error {
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Failed to set key: %s in Redis: %v", key, err)
		return err
	}
	return nil
}

// GetKey retrieves a value from Redis for a given key
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

// DeleteKey removes a key-value pair from Redis
func Del(client *redis.Client, key string) error {
	err := client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Failed to delete key: %s from Redis: %v", key, err)
		return err
	}
	return nil
}
