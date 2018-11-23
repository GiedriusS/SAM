package alerts

import (
	"context"
	"sort"
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
	alert1.EndsAt = startTs.Add(5 * time.Second).Format(TimeFormat)
	alert1.StartsAt = startTs.Format(TimeFormat)

	alert2 := NewAlert()
	alert2.EndsAt = startTs.Add(3 * time.Second).Format(TimeFormat)
	alert2.StartsAt = startTs.Add(1 * time.Second).Format(TimeFormat)

	notifications := []notification{
		notification{Alerts: []Alert{alert1}, Status: "resolved",
			Timestamp: startTs.Format(TimeFormat)},
		notification{Alerts: []Alert{alert2}, Status: "resolved",
			Timestamp: startTs.Format(TimeFormat)},
	}

	esclient, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		t.Fatalf("failed to setup elastic client: %v", err)
	}
	alertsource, err := NewElasticSearchSource(indexName, esclient, nil)
	if err != nil {
		t.Fatalf("failed to setup elastic alert source: %v", err)
	}
	for _, n := range notifications {
		err = addNotification(esclient, &n, indexName)
		if err != nil {
			t.Fatalf("failed to add alert: %v", err)
		}
	}
	// TODO(g-statkevicius): the ranges look bad here
	alerts, err := alertsource.GetAlertsFromTo("resolved",
		startTs.Add(-15*time.Second), startTs.Add(15*time.Second))
	if err != nil {
		t.Fatalf("failed to get related alerts: %v", err)
	}
	if alerts.Len() != 2 {
		t.Fatalf("retrieved %v alerts instead of %v", alerts.Len(), 2)
	}
	sort.Sort(alerts)
	if alerts.Alerts[0].Starts().After(alerts.Alerts[1].Starts()) {
		t.Fatalf("failed to sort: first alert starts at %v which is after the 2nd at %v",
			alerts.Alerts[0].Starts(), alerts.Alerts[1].Starts())
	}
}
