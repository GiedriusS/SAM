// Package alerts has functions and data types for storing alert information and
// calculating similar alerts.
package alerts

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"
)

// Alert stores the necessary data of one alert.
type Alert struct {
	Annotations  map[string]string `json:"annotations,omitempty"`
	StartsAt     string            `json:"startsAt,omitempty"`
	EndsAt       string            `json:"endsAt,omitempty"`
	GeneratorURL string            `json:"generatorURL,omitempty"`
	Labels       map[string]string `json:"labels"`
	Status       string            `json:"status,omitempty"`
	Related      map[string]uint   `json:"-"`
}

// NewAlert constructs a new Alert object.
func NewAlert() Alert {
	return Alert{Labels: make(map[string]string),
		Annotations: make(map[string]string),
		Related:     make(map[string]uint),
	}
}

// TimeFormat is the time format of alert boundaries
const TimeFormat = time.RFC3339

// Hash calculates the alert's hash. Used to identify identical alerts.
func (a *Alert) Hash() string {
	keys := []string{}
	h := sha256.New()

	for k := range a.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(a.Labels[k]))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// AlertSource is an interface for all alerts sources.
type AlertSource interface {
	GetAlertsFromTo(StartsAt, EndsAt time.Time) ([]Alert, error)
}

// State is the current state of the alerts parser.
type State struct {
	Firing      []string
	Alerts      map[string]*Alert
	LastUpdated time.Time
}

// NewState initializes a new State variable.
func NewState() State {
	return State{Firing: []string{},
		Alerts: make(map[string]*Alert),
	}
}

// AddAlert adds alert to the state and parses it.
func (s *State) AddAlert(a *Alert) {
	if _, ok := s.Alerts[a.Hash()]; ok != true {
		s.Alerts[a.Hash()] = a
	}

	s.updateRelated(a)
	s.parseAlertStatus(a)
}

// updateRelated updates the relatedness of an alert with the firing alerts.
func (s *State) updateRelated(alert *Alert) {
	for _, f := range s.Firing {
		if f == alert.Hash() {
			continue
		}
		alert.Related[f]++
	}
}

// parseAlertStatus parses the alert status and either adds it or removes it from firing.
func (s *State) parseAlertStatus(alert *Alert) {
	newFiring := []string{}

	switch alert.Status {
	case "firing":
		newFiring = append(newFiring, s.Firing...)
		newFiring = append(newFiring, alert.Hash())
	case "resolved":
		for _, f := range s.Firing {
			if f != alert.Hash() {
				newFiring = append(newFiring, f)
			}
		}
	}
	s.Firing = newFiring
}
