package alerts

import (
	"context"
	"testing"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

// addNotification adds a new notification to a index using the client.
func addNotification(c *elastic.Client, n *notification, index string) error {
	_, err := c.Index().Index(index).Type("alert_group").BodyJson(*n).Do(context.Background())
	return err
}

func TestGetAlertsFromTo(t *testing.T) {
	const indexName = "alertmanager-2018.11"

	startTs := time.Now()

	alert1 := NewAlert()
	alert1.StartsAt = startTs.Format(TimeFormat)
	alert1.EndsAt = startTs.Format(TimeFormat)
	alert1.Status = "resolved"
	alert1.Labels["a"] = "b"

	n := notification{Alerts: []Alert{alert1}, Timestamp: startTs.Format(time.RFC3339)}

	esclient, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		t.Fatalf("failed to setup elastic client: %v", err)
	}
	alertsource, err := NewElasticSearchSource(indexName, esclient, nil)
	if err != nil {
		t.Fatalf("failed to setup elastic alert source: %v", err)
	}
	err = addNotification(esclient, &n, indexName)
	if err != nil {
		t.Fatalf("failed to add notification: %v", err)
	}
	alerts, err := alertsource.GetAlertsFromTo(time.Now().Add(-15*time.Second),
		time.Now().Add(15*time.Second))
	if err != nil {
		t.Fatalf("failed to get alerts: %v", err)
	}
	if len(alerts) != 1 {
		t.Fatalf("retrieved %v alerts instead of %v", len(alerts), 1)
	}
}
