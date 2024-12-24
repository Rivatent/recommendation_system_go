package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"recommendation-service/internal/model"
	"time"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	log.Print("NewRedisCache created")

	ctx := context.Background()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Printf("Redis connected: %v", pong)

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisCache) GetRecommendationByID(id string) (model.Recommendation, error) {
	var recommendation model.Recommendation
	data, err := r.client.Get(r.ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return recommendation, redis.Nil
	} else if err != nil {
		return recommendation, err
	}

	err = json.Unmarshal([]byte(data), &recommendation)
	if err != nil {
		return recommendation, err
	}

	return recommendation, nil
}

func (r *RedisCache) GetRecommendationsByUserID(id string) ([]model.Recommendation, error) {
	var recommendations []model.Recommendation
	data, err := r.client.Get(r.ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return recommendations, redis.Nil
	} else if err != nil {
		return recommendations, err
	}

	err = json.Unmarshal([]byte(data), &recommendations)
	if err != nil {
		return recommendations, err
	}

	return recommendations, nil
}

func (r *RedisCache) SetRecommendationByID(id string, recommendation model.Recommendation, expiration time.Duration) error {
	data, err := json.Marshal(recommendation)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, id, data, expiration).Err()
}

func (r *RedisCache) SetRecommendationsByUserID(id string, recommendations []model.Recommendation, expiration time.Duration) error {
	data, err := json.Marshal(recommendations)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, id, data, expiration).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
