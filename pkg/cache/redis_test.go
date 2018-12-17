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
	if err != nil {
		t.Fatalf("failed to iniate redis client: %v", err)
	}
	st := alerts.State{Firing: []string{"a"}}
	err = rcache.PutState(&st)
	if err != nil {
		t.Fatalf("failed to put cache state: %v", err)
	}
	retrievedState, err := rcache.GetState()
	if err != nil {
		t.Fatalf("failed to get cache state: %v", err)
	}
	assert.Equal(t, retrievedState.Firing, st.Firing, "should be equal")
}
