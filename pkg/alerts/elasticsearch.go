package alerts

import (
	"time"

	"go.uber.org/zap"
	"gopkg.in/olivere/elastic.v5"
)

// ElasticSearchSource represents ElasticSearch as a source for alerts
type ElasticSearchSource struct {
	AlertSource

	client *elastic.Client
	logger *zap.Logger
}

// NewElasticSearchSource returns a new ElasticSearchSource
func NewElasticSearchSource(client *elastic.Client, logger *zap.Logger) (ElasticSearchSource, error) {
	if logger == nil {
		l, err := zap.NewProduction()
		if err != nil {
			return ElasticSearchSource{}, err
		}
		logger = l
	}
	return ElasticSearchSource{client: client, logger: logger}, nil
}

// GetAlertsFromTo retrieves the alerts with specified status where EndsAt <= UntilEndsAt
func (es ElasticSearchSource) GetAlertsFromTo(status string, UntilEndsAt time.Time) (RetrievedAlerts, error) {
	return RetrievedAlerts{}, nil
}
