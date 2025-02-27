package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"pictureloader/app_microservice/models"
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

func (rr *RedisRepo) GetPost(ctx context.Context, postID int) (models.PostUnit, bool, error) {
	redisResult, err := rr.rdb.Get(ctx, strconv.Itoa(postID)).Result()
	if errors.Is(err, redis.Nil) {
		return models.PostUnit{}, false, nil
	}
	if err != nil {
		return models.PostUnit{}, false, err
	}
	var result models.PostUnit
	json.Unmarshal([]byte(redisResult), &result)

	return result, true, nil
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

func (rr *RedisRepo) SetMostLikedPosts(ctx context.Context, posts []models.PostUnit) error {
	jsonData, err := json.Marshal(posts)
	if err != nil {
		return errors.New("json marshal posts error")
	}
	_, err = rr.rdb.Set(ctx, "MostLikedPosts", string(jsonData), time.Minute*1).Result()
	if err != nil {
		return err
	}
	return nil
}

func (rr *RedisRepo) GetMostLikedPosts(ctx context.Context) ([]models.PostUnit, error) {
	postsJson, err := rr.rdb.Get(ctx, "MostLikedPosts").Result()
	if err != nil {
		return nil, err
	}
	var posts []models.PostUnit
	err = json.Unmarshal([]byte(postsJson), &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (rr *RedisRepo) SetPost(ctx context.Context, postID int, resultJSON string) (bool, error) {
	err := rr.rdb.Set(ctx, strconv.Itoa(postID), resultJSON, time.Hour*10).Err()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
