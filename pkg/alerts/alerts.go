/*
Package alerts has needed things for storing alert information and calculating
similar alerts.

Honestly, maybe one can make a smarter algorithm but this algorithm at least is simple and understandable.
It is O(n^2) where n is the number of alerts. Sort all data by StartsAt.
Go through each alert and go forward while StartAt is equal, and backwards while StartsAt is lower or equal
to the current one. Only retrieve resolved alerts to reduce the noise.
*/
package alerts

import (
	"crypto/sha256"
	"sort"
	"time"
)

// Alert stores the necessary data of one alert
type Alert struct {
	Annotations  map[string]string `json:"annotations"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	StartsAt     string            `json:"startsAt"`
	Status       string            `json:"status"`
}

// NewAlert constructs a new Alert object
func NewAlert() Alert {
	return Alert{Labels: make(map[string]string), Annotations: make(map[string]string)}
}

// Starts parses StartsAt and retrieves time.Time
func (a *Alert) Starts() time.Time {
	starts, _ := time.Parse("0001-01-01T00:00:00Z", a.StartsAt)
	return starts
}

// Ends parses EndsAt and retrieves time.Time
func (a *Alert) Ends() time.Time {
	ends, _ := time.Parse("0001-01-01T00:00:00Z", a.EndsAt)
	return ends
}

// Hash calculates the alert's hash. Used to identify identical alerts.
func (a *Alert) Hash() string {
	h := sha256.New()
	for k, v := range a.Labels {
		h.Write([]byte(k))
		h.Write([]byte(v))
	}
	return string(h.Sum(nil))
}

// RetrievedAlerts is a wrapper around retrieved data which implements sort.Interface
type RetrievedAlerts struct {
	Alerts []Alert
	sort.Interface
	Related map[string]uint
}

// Len is part of sort.Interface for RetrievedAlerts
func (ra RetrievedAlerts) Len() int {
	return len(ra.Alerts)
}

// Swap is part of sort.Interface for RetrievedAlerts
func (ra RetrievedAlerts) Swap(i, j int) {
	ra.Alerts[i], ra.Alerts[j] = ra.Alerts[j], ra.Alerts[i]
}

// Less is part of sort.Interface. Sorted by StartsAt.
func (ra RetrievedAlerts) Less(i, j int) bool {
	return ra.Alerts[i].Starts().Before(ra.Alerts[j].Ends())
}

// AlertSource is an interface for all alerts sources
type AlertSource interface {
	GetAlertsFromTo(status string, StartsAt, EndsAt time.Time) (RetrievedAlerts, error)
}
