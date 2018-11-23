package alerts

import (
	"testing"
	"time"
)

func TestCalculateRelated(t *testing.T) {
	alert1 := NewAlert()
	alert1.Labels["a"] = "b"
	alert1.StartsAt = time.Now().Format(TimeFormat)
	alert1.EndsAt = time.Now().Add(10 * time.Second).Format(TimeFormat)

	alert2 := NewAlert()
	alert2.Labels["b"] = "a"
	alert2.StartsAt = time.Now().Add(1 * time.Second).Format(TimeFormat)
	alert2.EndsAt = time.Now().Add(5 * time.Second).Format(TimeFormat)

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
