package cache

import (
	"context"
	"github.com/Pr3c10us/boilerplate/internals/domains/cache"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository struct {
	redis                *redis.Client
	environmentVariables *configs.EnvironmentVariables
}

func NewRedisRepository(redis *redis.Client, environmentVariables *configs.EnvironmentVariables) cache.Repository {
	return &RedisRepository{
		redis, environmentVariables,
	}
}

func (repo *RedisRepository) Set(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	err := repo.redis.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
