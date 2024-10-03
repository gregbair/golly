package retry

import "time"

type TimeProvider interface {
	After(d time.Duration) <-chan time.Time
}

type systemTimeProvider struct{}

func (systemTimeProvider) After(d time.Duration) <-chan time.Time { return time.After(d) }
