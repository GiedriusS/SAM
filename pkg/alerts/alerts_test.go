package alerts

import (
	"testing"
	"time"
)

// TestProcess tests the overall process of how related alerts are calculated.
func TestProcess(t *testing.T) {
	start := time.Now()
	data := []Alert{
		Alert{Labels: map[string]string{"d": "e", "q": "z"},
			StartsAt: start.Format(TimeFormat),
			Status:   "firing",
			Related:  make(map[string]uint),
		},
		Alert{Labels: map[string]string{"a": "b", "c": "d"},
			StartsAt: start.Format(TimeFormat),
			EndsAt:   start.Format(TimeFormat),
			Status:   "firing",
			Related:  make(map[string]uint),
		},
		Alert{Labels: map[string]string{"d": "e", "q": "z"},
			StartsAt: start.Format(TimeFormat),
			EndsAt:   start.Add(1 * time.Second).Format(TimeFormat),
			Status:   "resolved",
			Related:  make(map[string]uint),
		},
	}

	s := NewState()
	for _, v := range data {
		err := s.AddAlert(v)
		if err != nil {
			t.Fatalf("failed to add alert: %v", err)
		}
	}

	if data[1].Related[data[0].Hash()] != 1 {
		t.Fatalf("failed to parse related data: 2nd alert has to be related to the 1st one (got %v)",
			data[1].Related[data[0].Hash()])
	}

	if data[0].Related[data[1].Hash()] == 1 {
		t.Fatalf("failed to parse related data: 1st alert must not be related to the 2nd one (got %v)",
			data[0].Related[data[1].Hash()])
	}

	if len(data[2].Related) != 0 {
		t.Fatalf("2nd alert is related to other %v alerts", len(data[2].Related))
	}

	if len(s.Firing) != 1 {
		t.Fatalf("%v alerts are still firing even though it was supposed to be 1", len(s.Firing))
	}
}

// TestCollision tests the case when two exact alerts are firing which is invalid.
func TestCollision(t *testing.T) {
	start := time.Now()
	data := []Alert{
		Alert{Labels: map[string]string{"foo": "bar"},
			StartsAt: start.Format(TimeFormat),
			Status:   "firing",
			Related:  make(map[string]uint),
		},
		Alert{Labels: map[string]string{"foo": "bar"},
			StartsAt: start.Format(TimeFormat),
			Status:   "firing",
			Related:  make(map[string]uint),
		},
	}

	s := NewState()
	s.AddAlert(data[0])
	err := s.AddAlert(data[1])
	if err == nil {
		t.Fatalf("adding the same alert did not return an error")
	}

	if len(data[0].Related) != 0 {
		t.Fatalf("1st alert is related to something (got %v, expected %v)",
			len(data[0].Related), 0)
	}

	if len(data[1].Related) != 0 {
		t.Fatalf("2nd alert is related to something (got %v, expected %v)",
			len(data[1].Related), 0)
	}

	if len(s.Firing) != 1 {
		t.Fatalf("%v alerts are firing when it should be %v", len(s.Firing),
			1)
	}
}
