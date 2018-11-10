/*
Package alerts has needed things for storing alert information and calculating
similar alerts.

Honestly, maybe one can make a smarter algorithm but this algorithm at least is simple and understandable.
It is O(n^2) where n is the number of alerts. Sort all data by StartsAt.
Go through each alert and go forward while StartAt is equal, and backwards while StartsAt is lower or equal
to the current one. Only retrieve resolved alerts to reduce the noise. Retrieve alerts with the EndsAt bound.

Save all data in a hash -> Alert structure.
*/
package alerts

import (
	"crypto/sha256"
	"sort"
	"time"
)

// Alert stores the necessary data of one alert
type Alert struct {
	Labels   map[string]string // Labels identify the alert
	StartsAt time.Time
	EndsAt   time.Time
	Related  map[string]uint
}

// NewAlert constructs a new Alert object
func NewAlert() Alert {
	return Alert{Labels: make(map[string]string), Related: make(map[string]uint)}
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
	return ra.Alerts[i].StartsAt.Before(ra.Alerts[j].StartsAt)
}

// AlertSource is an interface for all alerts sources
type AlertSource interface {
	GetAlertsFromTo(status string, UntilEndsAt time.Time) (RetrievedAlerts, error)
}
