package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gopkg.in/olivere/elastic.v5"
)

// ElasticSearchSource represents ElasticSearch as a source for alerts
type ElasticSearchSource struct {
	AlertSource

	client *elastic.Client
	logger *zap.Logger
	index  string
}

// NewElasticSearchSource returns a new ElasticSearchSource
func NewElasticSearchSource(index string, client *elastic.Client, logger *zap.Logger) (ElasticSearchSource, error) {
	if logger == nil {
		l, err := zap.NewProduction()
		if err != nil {
			return ElasticSearchSource{}, err
		}
		logger = l
	}
	return ElasticSearchSource{index: index, client: client, logger: logger}, nil
}

// GetAlertsFromTo retrieves the alerts with specified status where EndsAt <= UntilEndsAt
func (es ElasticSearchSource) GetAlertsFromTo(status string, UntilEndsAt time.Time) (RetrievedAlerts, error) {
	t := elastic.NewTermQuery("status", status)
	r := elastic.NewRangeQuery("alerts.endsAt").Lt(UntilEndsAt)
	q := elastic.NewBoolQuery().Must(t).Must(r)

	searchResult, err := es.client.Search().
		Index(es.index).
		Query(q).Sort("alerts.startsAt", true).
		Do(context.Background())

	if err != nil {
		return RetrievedAlerts{}, err
	}

	if searchResult.Hits.TotalHits == 0 {
		return RetrievedAlerts{}, nil
	}

	ret := RetrievedAlerts{}
	for _, hit := range searchResult.Hits.Hits {
		var a Alert
		if err := json.Unmarshal(*hit.Source, &a); err != nil {
			return RetrievedAlerts{}, fmt.Errorf("failed to unmarshal alert: %v", err)
		}

		ret.Alerts = append(ret.Alerts, a)
	}

	return ret, nil
}
