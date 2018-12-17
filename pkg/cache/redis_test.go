package cache

import (
	"testing"

	"github.com/GiedriusS/SAM/pkg/alerts"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedisSmoke(t *testing.T) {
	rclient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	rcache, err := NewRedisCache(rclient)
	assert.Nil(t, err, "should not fail to create a Redis client")

	st := alerts.State{Firing: []string{"a"}}
	assert.Nil(t, rcache.PutState(&st), "should not fail to put state")

	retrievedState, err := rcache.GetState()
	assert.Nil(t, err, "should not fail to retrieve state")
	assert.Equal(t, retrievedState.Firing, st.Firing, "should be equal")
}
