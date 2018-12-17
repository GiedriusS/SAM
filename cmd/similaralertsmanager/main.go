package main

import (
	"log"
	"net/http"

	"github.com/GiedriusS/SAM/pkg/alerts"
	"github.com/GiedriusS/SAM/pkg/api"
	"github.com/GiedriusS/SAM/pkg/cache"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/olivere/elastic.v5"
)

func main() {
	var (
		esinstances   = kingpin.Flag("elasticsearch", "ElasticSearch address").Required().Short('s').Strings()
		addr          = kingpin.Flag("addr", "API listen address").Default(":9888").Short('l').String()
		redisinstance = kingpin.Flag("redis", "Redis address").Required().Short('r').String()
	)
	kingpin.Parse()

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize Zap logger: %v", err)
	}
	defer l.Sync()

	rclient := redis.NewClient(&redis.Options{
		Addr: *redisinstance,
	})

	esclient, err := elastic.NewClient(elastic.SetSniff(false),
		elastic.SetURL(*esinstances...))
	if err != nil {
		l.Fatal("failed to initialize ElasticSearch client",
			zap.Error(err))
	}

	runSAM(l, rclient, esclient, addr)
}

func runSAM(l *zap.Logger, r *redis.Client, e *elastic.Client, addr *string) {
	rcache, err := cache.NewRedisCache(r)
	if err != nil {
		l.Fatal("failed to initialize new Redis cache", zap.Error(err))
	}
	state := alerts.NewState()
	newState, err := rcache.GetState()
	if err != nil {
		l.Info("failed to get cache from Redis", zap.Error(err))
	} else {
		state = newState
	}

	api := api.NewAPI(&state)
	srv := &http.Server{
		Handler: api.R,
		Addr:    *addr,
	}
	l.Fatal("failed to listen", zap.Error(srv.ListenAndServe()))
}
