package alerts

import (
	"testing"
	"time"
)

// TestProcess tests the overall process of how related alerts are calculated.
func TestProcess(t *testing.T) {
	data := []Alert{
		Alert{Labels: map[string]string{"d": "e", "q": "z"},
			StartsAt: time.Now().Format(TimeFormat),
			Status:   "firing",
			Related:  make(map[string]uint),
		},
		Alert{Labels: map[string]string{"a": "b", "c": "d"},
			StartsAt: time.Now().Format(TimeFormat),
			EndsAt:   time.Now().Format(TimeFormat),
			Status:   "resolved",
			Related:  make(map[string]uint),
		},
	}

	s := NewState()
	for _, v := range data {
		s.AddAlert(&v)
		s.UpdateRelated(&v)
		s.ParseAlertStatus(&v)
	}

	if data[1].Related[data[0].Hash()] != 1 {
		t.Fatalf("failed to parse related data: 2nd alert has to be related to the 1st one")
	}

	if data[0].Related[data[1].Hash()] == 1 {
		t.Fatalf("failed to parse related data: 1st alert must not be related to the 2nd one")
	}
}
