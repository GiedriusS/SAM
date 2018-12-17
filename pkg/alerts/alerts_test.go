package alerts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		assert.Nil(t, s.AddAlert(v), "should have succeeded")
	}

	assert.Contains(t, s.Alerts[data[1].Hash()].Related, data[0].Hash(), "2nd alert has to be related to the 1st one")
	assert.NotContains(t, s.Alerts[data[0].Hash()].Related, data[1].Hash(), "1st alert must not be related to the 2nd one")
	assert.Empty(t, s.Alerts[data[2].Hash()].Related, "must not be related to any other alerts")
	assert.Len(t, s.Firing, 1, "one alert must be firing")
}

// TestCollision tests the case when two exact alerts are firing/resolved which is invalid.
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
		Alert{Labels: map[string]string{"foo": "bar"},
			StartsAt: start.Format(TimeFormat),
			Status:   "resolved",
			Related:  make(map[string]uint),
		},
		Alert{Labels: map[string]string{"foo": "bar"},
			StartsAt: start.Format(TimeFormat),
			Status:   "resolved",
			Related:  make(map[string]uint),
		},
	}

	s := NewState()
	assert.Nil(t, s.AddAlert(data[0]), "adding alert must not error")
	assert.NotNil(t, s.AddAlert(data[1]), "adding same alert must return an error")
	assert.Len(t, data[0].Related, 0, "1st alert must not be related to the 2nd one")
	assert.Len(t, data[1].Related, 0, "2nd alert must not be related to the 1nd one")
	assert.Equal(t, len(s.Firing), 1, "one alert should be firing")

	assert.Nil(t, s.AddAlert(data[2]), "returned an error when it should not have")
	assert.Len(t, s.Firing, 0, "no alerts should be firing")

	assert.NotNil(t, s.AddAlert(data[3]), "should have returned an error")
}
