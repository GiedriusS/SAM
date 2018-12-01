package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

// addNotification sends the notification to specified alertmanager2es.
func addNotification(url string, n notification) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(n)
	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func TestGetAlertsFromTo(t *testing.T) {
	indexName := fmt.Sprintf("alertmanager-%s", time.Now().Format("2006.01"))

	startTs := time.Now()

	alert1 := NewAlert()
	alert1.Labels["fgbfgb"] = "xcvcxv"
	alert1.Status = "firing"
	alert1.StartsAt = startTs.Format(TimeFormat)

	alert2 := NewAlert()
	alert2.Labels["gbgb"] = "dfd"
	alert2.Status = "firing"
	alert2.StartsAt = startTs.Add(2 * time.Second).Format(TimeFormat)

	notification1 := notification{Alerts: []Alert{alert1},
		Timestamp: startTs.Format(TimeFormat), Version: "4"}

	notification2 := notification{Alerts: []Alert{alert2},
		Timestamp: startTs.Format(TimeFormat), Version: "4"}

	esclient, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		t.Fatalf("failed to setup elastic client: %v", err)
	}
	alertsource, err := NewElasticSearchSource(indexName, esclient, nil)
	if err != nil {
		t.Fatalf("failed to setup elastic alert source: %v", err)
	}
	err = addNotification("http://localhost:9097/webhook", notification1)
	if err != nil {
		t.Fatalf("failed to add notification #1: %v", err)
	}
	err = addNotification("http://localhost:9097/webhook", notification2)
	if err != nil {
		t.Fatalf("failed to add notification #2: %v", err)
	}

	t.Logf("waiting 10 seconds until the alerts appear in ES")
	select {
	case <-time.After(10 * time.Second):
		t.Logf("10 seconds passed, trying to query")
	}

	from := startTs
	to := startTs.Add(2 * time.Minute)
	alerts, err := alertsource.GetAlertsFromTo(from, to)
	if err != nil {
		t.Fatalf("failed to get alerts: %v", err)
	}
	if len(alerts) != 2 {
		t.Fatalf("retrieved %v alerts instead of %v", len(alerts), 2)
	}
}
