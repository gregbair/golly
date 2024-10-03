package retry

import (
	"context"
	"math/rand"
	"time"
)

type onRetryFunc func(uint, error)

type randomizerFunc func() float64

type delayFunc func(uint, context.Context) time.Duration

// DelayBackoffStrategy represents the different algorithms for introducing a delay between attempts
type DelayBackoffStrategy int

const (
	// Constant provides the same delay duration between attempts
	Constant DelayBackoffStrategy = iota
	// Linear is an algorithm that progresses the duration in a linear fashion
	Linear
)

type retryConfig struct {
	attempts        uint
	onRetry         onRetryFunc
	ctxt            context.Context
	backoffStrategy DelayBackoffStrategy
	jitter          bool
	randomizer      randomizerFunc
	delay           time.Duration
	maxDelay        time.Duration
	delayGenerator  delayFunc
	timeProvider    TimeProvider
}

// RetryOption is a functional configuration function
type RetryOption func(*retryConfig)

func defaultConfig() *retryConfig {
	return &retryConfig{
		attempts:        3,
		onRetry:         func(u uint, err error) {},
		ctxt:            context.Background(),
		backoffStrategy: Constant,
		jitter:          false,
		randomizer:      rand.Float64,
		delay:           0,
		maxDelay:        5 * time.Second,
		timeProvider:    systemTimeProvider{},
	}
}

// Attempts sets the number of attempts to allow before failure with a default of 3
func Attempts(attempts uint) RetryOption {
	return func(c *retryConfig) {
		c.attempts = attempts
	}
}

// OnRetry is a callback that is called after each attempt
func OnRetry(f onRetryFunc) RetryOption {
	return func(c *retryConfig) {
		c.onRetry = f
	}
}

// Context gives a context that is checked for cancellation on every attempt
func Context(ctxt context.Context) RetryOption {
	return func(c *retryConfig) {
		c.ctxt = ctxt
	}
}

// BackoffStrategy provides a choice of backoff algorithms
func BackoffStrategy(d DelayBackoffStrategy) RetryOption {
	return func(c *retryConfig) {
		c.backoffStrategy = d
	}
}

// Jitter decides whether a bit of randomization is added to the delay with a default of false
func Jitter(j bool) RetryOption {
	return func(c *retryConfig) {
		c.jitter = j
	}
}

// MaxDelay sets a maximum delay that is enforced after linear backoff and any jitter is applied
func MaxDelay(t time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.maxDelay = t
	}
}

// Delay sets a baseline delay duration between each attempt
func Delay(t time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.delay = t
	}
}

// TimeProviderImpl provides a TimeProvider that is used for counting the delay with a default of the stdlib time functions
func TimeProviderImpl(t TimeProvider) RetryOption {
	return func(c *retryConfig) {
		c.timeProvider = t
	}
}

// Randomizer provides a randomizing function used when jitter is applied with a default of the rand.Float64 function
func Randomizer(f randomizerFunc) RetryOption {
	return func(c *retryConfig) {
		c.randomizer = f
	}
}

// DelayGenerator provides a function that can be used to do custom calculations of the delay duration
func DelayGenerator(f delayFunc) RetryOption {
	return func(c *retryConfig) {
		c.delayGenerator = f
	}
}
