package main

import (
	"log"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/olivere/elastic.v5"
)

func main() {
	var (
		esinstances = kingpin.Flag("elasticsearch", "IP:PORT pair of an ES instance").Required().Short('s').Strings()
		//listen        = kingpin.Flag("listenaddr", "Listen address").Default(":9888").Short('l').String()
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
	_, err = rclient.Ping().Result()
	if err != nil {
		l.Fatal("failed to ping Redis",
			zap.Error(err))
	}

	esclient, err := elastic.NewSimpleClient(elastic.SetURL(*esinstances...))
	if err != nil {
		l.Fatal("failed to initialize Elasticsearch client",
			zap.Error(err))
	}

	runSAM(l, rclient, esclient)
}

func runSAM(l *zap.Logger, r *redis.Client, e *elastic.Client) {
	l.Info("hello", zap.String("hi", "a"))
}
