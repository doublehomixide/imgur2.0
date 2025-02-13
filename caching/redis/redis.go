package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedisClient() *RedisRepo {
	rdb := redis.NewClient(&redis.Options{})
	return &RedisRepo{rdb: rdb}
}

type RedisRepo struct {
	rdb *redis.Client
}

func (rr *RedisRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rr.rdb.Set(ctx, key, value, expiration).Err()
}

func (rr *RedisRepo) Get(ctx context.Context, key string) (string, error) {
	result, err := rr.rdb.Get(ctx, key).Result()
	return result, err
}

func (rr *RedisRepo) Delete(ctx context.Context, key string) error {
	return rr.rdb.Del(ctx, key).Err()
}
