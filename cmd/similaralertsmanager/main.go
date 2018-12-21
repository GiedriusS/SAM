package main

import (
	"log"
	"net/http"
	"sync"
	"time"

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
		cacheinterval = kingpin.Flag("cacheinterval", "Interval in seconds between updates of cache").Default("5").Int()
		esInterval    = kingpin.Flag("esinterval", "Interval in seconds between parsing new alerts").Default("10").Int()
		esIndexName   = kingpin.Flag("esindex", "ElasticSearch index name").Default("alertmanager").Short('i').String()
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

	runSAM(l, rclient, esclient, addr, cacheinterval, esInterval, esIndexName)
}

func runSAM(l *zap.Logger, r *redis.Client, e *elastic.Client, addr *string, cacheint *int, esint *int, esIndex *string) {
	rcache, err := cache.NewRedisCache(r)
	if err != nil {
		l.Fatal("failed to initialize new Redis cache", zap.Error(err))
	}
	esSource, err := alerts.NewElasticSearchSource(*esIndex, e, l)
	if err != nil {
		l.Fatal("failed to initialize new ElasticSearch source", zap.Error(err))
	}

	state := alerts.NewState()
	newState, err := rcache.GetState()
	if err != nil {
		l.Info("failed to get cache from Redis", zap.Error(err))
	} else {
		l.Info("got cache from Redis", zap.Time("last updated", newState.GetLastUpdated()))
		state = newState
	}

	api := api.NewAPI(&state, l)
	srv := &http.Server{
		Handler: api.R,
		Addr:    *addr,
	}

	stateLock := sync.Mutex{}

	go func() {
		for {
			select {
			case <-time.After(time.Duration(*esint) * time.Second):
			}
			stateLock.Lock()

			l.Info("getting alerts", zap.Time("from", state.GetLastUpdated()), zap.Time("to", time.Now()))
			alerts, err := esSource.GetAlertsFromTo(state.GetLastUpdated(), time.Now())
			if err != nil {
				l.Error("failed to get alerts", zap.Error(err))
				stateLock.Unlock()
				continue
			}
			for _, a := range alerts {
				err := state.AddAlert(a)
				if err != nil {
					l.Error("failed to add alert", zap.Error(err))
				}
			}
			l.Info("finished getting alerts")
			stateLock.Unlock()
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(time.Duration(*cacheint) * time.Second):
			}
			stateLock.Lock()

			l.Info("putting state")
			err := rcache.PutState(&state)
			if err != nil {
				l.Error("failed to put cache", zap.Error(err))
				stateLock.Unlock()
				continue
			}
			l.Info("finished putting cache")
			stateLock.Unlock()
		}
	}()

	l.Fatal("failed to listen", zap.Error(srv.ListenAndServe()))
}
