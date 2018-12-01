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

// addNotification sends the alerts to url.
func addNotification(url string, alerts ...Alert) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(alerts)
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
	alert1.Labels["a"] = "b"

	alert2 := NewAlert()
	alert2.Labels["a"] = "bbb"

	esclient, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		t.Fatalf("failed to setup elastic client: %v", err)
	}
	alertsource, err := NewElasticSearchSource(indexName, esclient, nil)
	if err != nil {
		t.Fatalf("failed to setup elastic alert source: %v", err)
	}
	err = addNotification("http://localhost:9093/api/v1/alerts", alert1, alert2)
	if err != nil {
		t.Fatalf("failed to add alert: %v", err)
	}
	t.Logf("waiting 30 seconds for alerts to appear in ES")
	select {
	case <-time.After(15 * time.Second):
		t.Logf("30 seconds have passed, trying to query")
	}
	alerts, err := alertsource.GetAlertsFromTo(startTs, startTs.Add(5*time.Second))
	if err != nil {
		t.Fatalf("failed to get alerts: %v", err)
	}
	if len(alerts) != 2 {
		t.Fatalf("retrieved %v alerts instead of %v", len(alerts), 2)
	}
}
