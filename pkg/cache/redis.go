package cache

import (
	"encoding/json"

	"github.com/GiedriusS/SAM/pkg/alerts"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// RedisCache represents a cache that is based on a underlying Redis.
type RedisCache struct {
	r *redis.Client
	Cache
	key string
}

// NewRedisCache constructs and returns a new RedisCache.
func NewRedisCache(r *redis.Client, key string) (RedisCache, error) {
	_, err := r.Ping().Result()
	if err != nil {
		return RedisCache{}, errors.Wrapf(err, "failed to ping Redis")
	}

	return RedisCache{r: r, key: key}, nil
}

// PutState saves the state into Redis.
func (r *RedisCache) PutState(s *alerts.State) error {
	// TODO(g-statkevicius): make this smarter so we wouldn't have to marshal the whole thing
	b, err := json.Marshal(*s)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal state")
	}

	err = r.r.Do("SET", r.key, string(b)).Err()
	return err
}

// GetState gets the state from Redis.
func (r *RedisCache) GetState() (alerts.State, error) {
	val, err := r.r.Get(r.key).Result()
	if err != nil {
		return alerts.State{}, errors.Wrapf(err, "failed to get key")
	}

	ret := alerts.NewState()
	err = json.Unmarshal([]byte(val), &ret)
	if err != nil {
		return alerts.State{}, errors.Wrapf(err, "failed to unmarshal")
	}
	return ret, nil
}
