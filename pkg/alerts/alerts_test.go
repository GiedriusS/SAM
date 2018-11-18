package alerts

import (
	"testing"
	"time"
)

func TestCalculateRelated(t *testing.T) {
	alert1 := NewAlert()
	alert1.Labels["a"] = "b"
	alert1.EndsAt = "0001-01-01T00:00:10Z"
	alert1.StartsAt = "0001-01-01T00:00:00Z"

	alert2 := NewAlert()
	alert2.Labels["b"] = "a"
	alert2.EndsAt = "0001-01-01T00:00:05Z"
	alert2.StartsAt = "0001-01-01T00:00:01Z"

	data := AugmentedAlerts{
		Alerts:        []Alert{alert1, alert2},
		LastTimestamp: time.Time(time.Now()),
	}

	data.CalculateRelated()
	if data.Alerts[1].Related[alert1.Hash()] != 1 {
		t.Fatalf("expected 2nd alert to be related to the 1st alert")
	}

	if data.Alerts[0].Related[alert2.Hash()] == 1 {
		t.Fatalf("expected 1st alert not to be related to the 1st alert")
	}
}
