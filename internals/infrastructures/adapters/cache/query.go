package cache

import (
	"context"
	"errors"
	"time"
)

func (repo *RedisRepository) Get(key string) (string, error) {
	ctx := context.Background()
	value, err := repo.redis.Get(ctx, key).Result()
	if err != nil {
		return value, err
	} else {
		return value, nil
	}
}

func (repo *RedisRepository) TTL(key string) (time.Duration, error) {
	ctx := context.Background()
	ttl, err := repo.redis.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if ttl == 0 {
		return 0, errors.New("the key has expired")
	} else if ttl == -1 {
		return 0, errors.New("the key has no expiration")
	} else if ttl == -2 {
		return 0, errors.New("the key does not exist")
	} else if ttl < -2 {
		return 0, errors.New("error eje")
	}
	return ttl, nil
}
