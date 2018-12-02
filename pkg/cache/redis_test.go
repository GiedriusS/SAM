package cache

import (
	"testing"
	"time"

	"github.com/GiedriusS/SAM/pkg/alerts"
	"github.com/go-redis/redis"
)

func TestRedisSmoke(t *testing.T) {
	rclient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	rcache, err := NewRedisCache(rclient)
	if err != nil {
		t.Fatalf("failed to iniate redis client: %v", err)
	}
	st := alerts.State{LastUpdated: time.Now()}
	err = rcache.PutState(&st)
	if err != nil {
		t.Fatalf("failed to put cache state: %v", err)
	}
	retrievedState, err := rcache.GetState()
	if err != nil {
		t.Fatalf("failed to get cache state: %v", err)
	}
	if !retrievedState.LastUpdated.Equal(st.LastUpdated) {
		t.Fatalf("retrieved and saved states' last updated times are different")
	}
}
