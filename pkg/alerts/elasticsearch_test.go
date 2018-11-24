package alerts

import (
	"context"
	"testing"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

var alert1 = `
{
	  "alerts": [
		{
		  "annotations": {
			"info": "The disk sda1 is running full",
			"summary": "please check the instance example1"
		  },
		  "endsAt": "0001-01-01T00:00:00Z",
		  "generatorURL": "",
		  "labels": {
			"alertname": "DiskRunningFull",
			"dev": "sda1",
			"instance": "example1"
		  },
		  "startsAt": "2018-11-24T21:19:34.3730271Z",
		  "status": "firing"
		}
	  ],
	  "commonAnnotations": {
		"info": "The disk sda1 is running full",
		"summary": "please check the instance example1"
	  },
	  "commonLabels": {
		"alertname": "DiskRunningFull",
		"dev": "sda1",
		"instance": "example1"
	  },
	  "externalURL": "http://ce1fa40d0cb5:9093",
	  "groupLabels": {
		"alertname": "DiskRunningFull"
	  },
	  "receiver": "web\\.hook",
	  "status": "firing",
	  "version": "4",
	  "groupKey": "{}:{alertname=\"DiskRunningFull\"}"
}
`

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
	alert1.EndsAt = startTs.Add(1 * time.Second).Format(TimeFormat)
	alert1.Labels["a"] = "b"

	n := notification{Alerts: []Alert{alert1}, Status: "resolved", Timestamp: startTs}

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
	alerts, err := alertsource.GetAlertsFromTo(startTs.Add(-90*time.Second), startTs.Add(90*time.Second))
	if err != nil {
		t.Fatalf("failed to get alerts: %v", err)
	}
	if len(alerts) != 1 {
		t.Fatalf("retrieved %v alerts instead of %v", len(alerts), 1)
	}
}
