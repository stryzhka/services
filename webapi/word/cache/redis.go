package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"webapi/models"

	"github.com/redis/go-redis/v9"
)

type RedisWordCache struct {
	r *redis.Client
}

func NewRedisWordCache(r *redis.Client) *RedisWordCache {
	return &RedisWordCache{r}
}

func (c *RedisWordCache) GetOne(ctx context.Context, key string) (interface{}, error) {
	val, err := c.r.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result *models.Word
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RedisWordCache) GetMany(ctx context.Context, key string) ([]interface{}, error) {
	val, err := c.r.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result []interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RedisWordCache) Set(ctx context.Context, key string, exp time.Duration, value interface{}) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	if err := c.r.Set(ctx, key, data, exp).Err(); err != nil {
		return false, err
	}

	return true, nil
}

func (c *RedisWordCache) Delete(ctx context.Context, key string) (bool, error) {
	deleted, err := c.r.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return deleted > 0, nil
}
