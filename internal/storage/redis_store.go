package storage

import (
	"context"
	"os"

	"log"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore() *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"), // No password set
		DB:       0,                           // Use default DB
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisStore{Client: rdb}
}

func (store *RedisStore) SetPlan(id string, planData string) error {
	err := store.Client.Set(Ctx, id, planData, 0).Err()
	return err
}

func (store *RedisStore) GetPlan(id string) (string, error) {
	val, err := store.Client.Get(Ctx, id).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (store *RedisStore) DeletePlan(id string) error {
	err := store.Client.Del(Ctx, id).Err()
	return err
}
