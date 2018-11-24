// Package cache is used to store the alerts parser state into Redis.
package cache

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// Cache is the general interface all cache providers must implement.
type Cache interface {
}

// RedisCache represents a cache that is based on a underlying Redis.
type RedisCache struct {
	r *redis.Client
	Cache
}

// NewRedisCache constructs and returns a new RedisCache.
func NewRedisCache(r *redis.Client) (RedisCache, error) {
	_, err := r.Ping().Result()
	if err != nil {
		return RedisCache{}, errors.Wrapf(err, "failed to ping Redis")
	}

	return RedisCache{r: r}, nil
}
