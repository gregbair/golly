package circuitbreaker

import (
	"context"
	"errors"
	"time"

	"github.com/gregbair/golly"
)

type config[TResult any] struct {
	failureRatio      float64
	minimumThroughput uint
	samplingDuration  time.Duration
	breakDuration     time.Duration
	timeProvider      golly.TimeProvider
	onOpened          func(context.Context, time.Duration, TResult, error)
	onClosed          func(context.Context, TResult, error)
	onHalfOpen        func(context.Context)
}

func defaultConfig[TResult any]() *config[TResult] {
	return &config[TResult]{
		failureRatio:      0.1,
		minimumThroughput: 100,
		samplingDuration:  30 * time.Second,
		breakDuration:     5 * time.Second,
		timeProvider:      golly.SystemTimeProvider{},
		onOpened:          nil,
		onClosed:          nil,
	}
}

// Option is a function configuration function
type Option[TResult any] func(*config[TResult])

// FailureRatio sets the failure percentage at which the circuit will open.
// The default value is 0.1. Value must be between 0 and 1.0.
func FailureRatio(r float64) Option[any] {
	return func(c *config[any]) {
		c.failureRatio = r
	}
}

// MinimumThroughput sets the minimum sets how many actions or more must pass through the circuit in the time-slice,
// for statistics to be considered significant and the circuit-breaker to come into action.
// Default value is 100, minimum is 2.
func MinimumThroughput(u uint) Option[any] {
	return func(c *config[any]) {
		c.minimumThroughput = u
	}
}

// SamplingDuration is the duration of samppling for failures.
// Default is 30 seconds. Value must be at least 0.5 seconds and at most an hour.
func SamplingDuration(d time.Duration) Option[any] {
	return func(c *config[any]) {
		c.samplingDuration = d
	}
}

// BreakDuration is the duration of break the circuit will stay open before resetting.
// Default value is 5 seconds. Value must be at least 0.5 seconds and at most an hour.
func BreakDuration(d time.Duration) Option[any] {
	return func(c *config[any]) {
		c.breakDuration = d
	}
}

// OnOpened is the callback that is called when the circuit breaker enters the Open state
func OnOpened[TResult any](f func(context.Context, time.Duration, TResult, error)) Option[TResult] {
	return func(c *config[TResult]) {
		c.onOpened = f
	}
}

// OnClosed is the callback that is called when the circuit breaker enters the Closed state
func OnClosed[TResult any](f func(context.Context, TResult, error)) Option[TResult] {
	return func(c *config[TResult]) {
		c.onClosed = f
	}
}

func onHalfOpened(f func(context.Context)) Option[any] {
	return func(c *config[any]) {
		c.onHalfOpen = f
	}
}

// timeProvider is a way to inject time handling into the pipeline.
func timeProvider(t golly.TimeProvider) Option[any] {
	return func(c *config[any]) {
		c.timeProvider = t
	}
}

func validateConfig(c *config[any]) error {
	var err error
	if c.failureRatio < 0 || c.failureRatio > 1.0 {
		err = errors.Join(err, errors.New("failure ration out of range, must be between 0 and 1.0"))
	}

	if c.minimumThroughput < 2 {
		err = errors.Join(err, errors.New("minimum throughput is too low, must be at least 2"))
	}

	if c.samplingDuration.Milliseconds() < 500 || c.samplingDuration.Hours() > 1 {
		err = errors.Join(err, errors.New("sampling duration is out of range, must be at between 500 milliseconds and an hour"))
	}

	if c.breakDuration.Milliseconds() < 500 || c.breakDuration.Hours() > 1 {
		err = errors.Join(err, errors.New("break duration is out of range, must be at between 500 milliseconds and an hour"))
	}

	return err
}
