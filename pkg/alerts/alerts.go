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

// Starts parses StartsAt and retrieves time.Time.
func (a *Alert) Starts() time.Time {
	starts, _ := time.Parse(TimeFormat, a.StartsAt)
	return starts
}

// Ends parses EndsAt and retrieves time.Time.
func (a *Alert) Ends() time.Time {
	ends, _ := time.Parse(TimeFormat, a.EndsAt)
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

// AugmentedAlerts is a wrapper around retrieved data which implements sort.Interface and
// stores all information about the current state of alerts.
type AugmentedAlerts struct {
	Alerts []Alert
	sort.Interface
	LastTimestamp time.Time
}

// Len is part of sort.Interface for AugmentedAlerts.
func (aa AugmentedAlerts) Len() int {
	return len(aa.Alerts)
}

// Swap is part of sort.Interface for AugmentedAlerts.
func (aa AugmentedAlerts) Swap(i, j int) {
	aa.Alerts[i], aa.Alerts[j] = aa.Alerts[j], aa.Alerts[i]
}

// Less is part of sort.Interface. Sorted by StartsAt.
func (aa AugmentedAlerts) Less(i, j int) bool {
	return aa.Alerts[i].Starts().Before(aa.Alerts[j].Starts())
}

// Merge merges the specified AugmentedAlerts into the current one.
func (aa *AugmentedAlerts) Merge(src *AugmentedAlerts) error {
	for _, v := range src.Alerts {
		aa.Alerts = append(aa.Alerts, v)
	}
	aa.LastTimestamp = src.LastTimestamp
	return nil
}

// CalculateRelated calculates related alerts in augmented alerts.
// Invoke it before merging.
func (aa *AugmentedAlerts) CalculateRelated() error {
	sort.Sort(*aa)
	for k, v := range aa.Alerts {
		now := k
		for now < len(aa.Alerts) && aa.Alerts[now].Starts() == aa.Alerts[k].Starts() {
			v.Related[aa.Alerts[now].Hash()]++
			now++
		}

		now = k
		for now >= 0 {
			v.Related[aa.Alerts[now].Hash()]++
			now--
		}
	}
	return nil
}

// AlertSource is an interface for all alerts sources.
type AlertSource interface {
	GetAlertsFromTo(status string, StartsAt, EndsAt time.Time) (AugmentedAlerts, error)
}
