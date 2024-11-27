package golly

import "time"

// TimeProvider defines a contract for providing time-related functionality, primarily for creating timers that fire after a specified duration.
type TimeProvider interface {
	UtcNow() time.Time
	After(d time.Duration) <-chan time.Time
}

type SystemTimeProvider struct{}

func (SystemTimeProvider) UtcNow() time.Time { return time.Now().UTC() }

func (SystemTimeProvider) After(d time.Duration) <-chan time.Time { return time.After(d) }
