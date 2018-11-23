package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gopkg.in/olivere/elastic.v5"
)

// This is the type that alertmanager2es uses.
type notification struct {
	Alerts            []Alert           `json:"alerts"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	CommonLabels      map[string]string `json:"commonLabels"`
	ExternalURL       string            `json:"externalURL"`
	GroupLabels       map[string]string `json:"groupLabels"`
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`

	// Timestamp records when the alert notification was received
	Timestamp string `json:"@timestamp"`
}

// ElasticSearchSource represents ElasticSearch as a source for alerts.
type ElasticSearchSource struct {
	AlertSource

	client *elastic.Client
	logger *zap.Logger
	index  string
}

// NewElasticSearchSource returns a new ElasticSearchSource.
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

// GetAlertsFromTo retrieves the alerts with specified status between specified boundaries.
func (es ElasticSearchSource) GetAlertsFromTo(status string, from, to time.Time) (AugmentedAlerts, error) {
	t := elastic.NewTermQuery("status", status)
	q := elastic.NewBoolQuery().Must(t).Filter(
		elastic.NewRangeQuery("@timestamp").From(from).
			To(to))

	searchResult, err := es.client.Search().
		Index(es.index).
		Query(q).Sort("alerts.startsAt", true).
		Do(context.Background())

	if err != nil {
		return AugmentedAlerts{}, err
	}

	if searchResult.Hits.TotalHits == 0 {
		return AugmentedAlerts{}, nil
	}

	ret := AugmentedAlerts{}
	for _, hit := range searchResult.Hits.Hits {
		var n notification
		if err := json.Unmarshal(*hit.Source, &n); err != nil {
			return AugmentedAlerts{}, fmt.Errorf("failed to unmarshal notification: %v", err)
		}

		for _, a := range n.Alerts {
			ret.Alerts = append(ret.Alerts, a)
		}
	}

	return ret, nil
}
