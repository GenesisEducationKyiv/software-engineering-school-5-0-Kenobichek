package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"internal/services/weather-service/domain"

	"github.com/redis/go-redis/v9"
)

const redisConnectionTimeout = 5 * time.Second

type redisClientManager interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
	Ping(ctx context.Context) *redis.StatusCmd
}

type RedisCache struct {
	client redisClientManager
	ttl    time.Duration
}

func NewRedisCache(addr, password string, db int, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), redisConnectionTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ttl:    ttl,
	}, nil
}

func (r *RedisCache) Get(ctx context.Context, city string) (*domain.Metrics, error) {
	key := r.buildKey(city)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var metrics domain.Metrics
	if err := json.Unmarshal([]byte(data), &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return &metrics, nil
}

func (r *RedisCache) Set(ctx context.Context, city string, metrics domain.Metrics) error {
	return r.SetWithTTL(ctx, city, metrics, r.ttl)
}

func (r *RedisCache) SetWithTTL(ctx context.Context, city string, metrics domain.Metrics, ttl time.Duration) error {
	key := r.buildKey(city)
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, city string) error {
	key := r.buildKey(city)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) buildKey(city string) string {
	return fmt.Sprintf("weather:%s", city)
}

func (r *RedisCache) GetDefaultTTL() time.Duration {
	return r.ttl
}
