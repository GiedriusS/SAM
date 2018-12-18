package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/olivere/elastic.v5"
)

// ElasticSearchSource represents ElasticSearch as a source for alerts.
type ElasticSearchSource struct {
	AlertSource

	client *elastic.Client
	logger *zap.Logger
	index  string
}

// This is the type that alertmanager2es uses.
type notification struct {
	Alerts            []Alert           `json:"alerts"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	CommonLabels      map[string]string `json:"commonLabels"`
	ExternalURL       string            `json:"externalURL"`
	GroupLabels       map[string]string `json:"groupLabels"`
	Receiver          string            `json:"receiver"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`

	// Timestamp records when the alert notification was received
	Timestamp string `json:"@timestamp"`
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

// GetAlertsFromTo retrieves the alerts between specified boundaries.
func (es ElasticSearchSource) GetAlertsFromTo(from, to time.Time) (ret []Alert, err error) {
	query := elastic.NewRangeQuery("@timestamp").From(from).To(to)

	searchResult, err := es.client.Search(es.index).
		Query(query).
		Do(context.Background())

	if err != nil {
		return
	}

	if searchResult.Hits.TotalHits == 0 {
		return ret, nil
	}

	for _, hit := range searchResult.Hits.Hits {
		var n notification
		if err := json.Unmarshal(*hit.Source, &n); err != nil {
			return ret, fmt.Errorf("failed to unmarshal notification: %v", err)
		}

		for _, alert := range n.Alerts {
			a := NewAlert()
			if err := mergo.Merge(&a, alert); err != nil {
				return ret, errors.Wrapf(err, "failed to merge alert")
			}
			ret = append(ret, a)
		}
	}

	return ret, nil
}
