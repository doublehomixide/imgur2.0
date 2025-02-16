package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func NewRedisClient(imageRepo imageRepo) *RedisRepo {
	rdb := redis.NewClient(&redis.Options{})
	return &RedisRepo{rdb: rdb, imageRepo: imageRepo}
}

type imageRepo interface {
	GetImageLinkedPost(ctx context.Context, imageSK string) (int, error)
}

type RedisRepo struct {
	rdb       *redis.Client
	imageRepo imageRepo
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

func (rr *RedisRepo) InvalidatePost(ctx context.Context, postID int) (bool, error) {
	result, err := rr.rdb.Del(ctx, strconv.Itoa(postID)).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
