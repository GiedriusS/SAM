package alerts

import (
	"sort"
	"testing"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

func TestGetAlertsFromTo(t *testing.T) {
	esclient, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		t.Fatalf("failed to setup elastic client: %v", err)
	}
	alertsource, err := NewElasticSearchSource("alertmanager-2018.11", esclient, nil)
	if err != nil {
		t.Fatalf("failed to setup elastic alert source: %v", err)
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
