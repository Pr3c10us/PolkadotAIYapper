package cache

import "time"

type Repository interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	TTL(key string) (time.Duration, error)
}
