package cache

import (
	"github.com/GiedriusS/SAM/pkg/alerts"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

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

// PutState saves the state into Redis.
func (r *RedisCache) PutState(s *alerts.State) error {
	return nil
}

// GetState gets the state from Redis.
func (r *RedisCache) GetState() (alerts.State, error) {
	return alerts.State{}, nil
}
