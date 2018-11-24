// Package alerts has functions and data types for storing alert information and
// calculating similar alerts.
package alerts

import (
	"crypto/sha256"
	"time"
)

// Alert stores the necessary data of one alert.
type Alert struct {
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	Status       string            `json:"status"`
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
	h := sha256.New()
	for k, v := range a.Labels {
		h.Write([]byte(k))
		h.Write([]byte(v))
	}
	return string(h.Sum(nil))
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

// AddAlert adds alert to the state if it does not exist already.
func (s *State) AddAlert(a *Alert) {
	if _, ok := s.Alerts[a.Hash()]; ok != true {
		s.Alerts[a.Hash()] = a
	}
}

// UpdateRelated updates the relatedness of an alert with the firing alerts.
func (s *State) UpdateRelated(alert *Alert) {
	for _, f := range s.Firing {
		if f == alert.Hash() {
			continue
		}
		alert.Related[f]++
	}
}

// ParseAlertStatus parses the alert status and either adds it or removes it from firing.
func (s *State) ParseAlertStatus(alert *Alert) {
	newFiring := s.Firing[:0]
	switch alert.Status {
	case "firing":
		s.Firing = append(s.Firing, alert.Hash())
	case "resolved":
		for _, v := range s.Firing {
			if v != alert.Hash() {
				newFiring = append(newFiring, v)
			}
		}
		s.Firing = newFiring
	}
}
