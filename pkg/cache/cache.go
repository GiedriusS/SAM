// Package cache is used to store the alerts parser state into Redis.
package cache

import (
	"github.com/GiedriusS/SAM/pkg/alerts"
)

// Cache is the general interface all cache providers must implement.
type Cache interface {
	PutState(s *alerts.State) error
	GetState() (alerts.State, error)
}
