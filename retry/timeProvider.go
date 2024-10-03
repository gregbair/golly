package retry

import "time"

// TimeProvider defines a contract for providing time-related functionality, primarily for creating timers that fire after a specified duration.
type TimeProvider interface {
	After(d time.Duration) <-chan time.Time
}

type systemTimeProvider struct{}

func (systemTimeProvider) After(d time.Duration) <-chan time.Time { return time.After(d) }
