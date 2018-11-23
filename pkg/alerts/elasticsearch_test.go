package alerts

import (
	"context"
	"sort"
	"testing"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

// addAlert adds a new alert to a specified index using the specified client.
func addAlert(c *elastic.Client, a *Alert, index string) error {
	_, err := c.Index().Index(index).Type("alert").BodyJson(*a).Do(context.Background())
	return err
}

func TestGetAlertsFromTo(t *testing.T) {
	alert1 := NewAlert()
	alert1.Labels["a"] = "b"
	alert1.EndsAt = time.Now().Add(5 * time.Second).Format(TimeFormat)
	alert1.StartsAt = time.Now().Format(TimeFormat)

	alert2 := NewAlert()
	alert2.Labels["b"] = "a"
	alert2.EndsAt = time.Now().Add(3 * time.Second).Format(TimeFormat)
	alert2.StartsAt = time.Now().Add(1 * time.Second).Format(TimeFormat)

	const indexName = "alertmanager-2018.11"

	esclient, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		t.Fatalf("failed to setup elastic client: %v", err)
	}
	alertsource, err := NewElasticSearchSource(indexName, esclient, nil)
	if err != nil {
		t.Fatalf("failed to setup elastic alert source: %v", err)
	}
	err = addAlert(esclient, &alert1, indexName)
	if err != nil {
		t.Fatalf("failed to add alert #1: %v", err)
	}
	err = addAlert(esclient, &alert2, indexName)
	if err != nil {
		t.Fatalf("failed to add alert #2: %v", err)
	}
	alerts, err := alertsource.GetAlertsFromTo("resolved", time.Now(),
		time.Now().Add(99999*time.Second))
	if err != nil {
		t.Fatalf("failed to get related alerts: %v", err)
	}
	if alerts.Len() != 2 {
		t.Fatalf("retrieved %v alerts instead of %v", alerts.Len(), 2)
	}
	sort.Sort(alerts)
	if alerts.Alerts[0].Starts().After(alerts.Alerts[1].Starts()) {
		t.Fatalf("failed to sort: first alert starts at before %v which is before the 2nd at %v",
			alerts.Alerts[0].Starts(), alerts.Alerts[1].Starts())
	}
}
