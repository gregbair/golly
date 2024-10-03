package retry

import (
	"context"
	"math/rand"
	"time"
)

type OnRetryFunc func(uint, error)

type RandomizerFunc func() float64

type DelayFunc func(uint, context.Context) time.Duration

type DelayBackOffStrategy int

const (
	Constant DelayBackOffStrategy = iota
	Linear
)

type retryConfig struct {
	attempts        uint
	onRetry         OnRetryFunc
	ctxt            context.Context
	backOffStrategy DelayBackOffStrategy
	jitter          bool
	randomizer      RandomizerFunc
	delay           time.Duration
	maxDelay        time.Duration
	delayGenerator  DelayFunc
	timeProvider    TimeProvider
}

type RetryOption func(*retryConfig)

func defaultConfig() *retryConfig {
	return &retryConfig{
		attempts:        3,
		onRetry:         func(u uint, err error) {},
		ctxt:            context.Background(),
		backOffStrategy: Constant,
		jitter:          false,
		randomizer:      rand.Float64,
		delay:           0,
		maxDelay:        5 * time.Second,
		timeProvider:    systemTimeProvider{},
	}
}

func Attempts(attempts uint) RetryOption {
	return func(c *retryConfig) {
		c.attempts = attempts
	}
}

func OnRetry(f OnRetryFunc) RetryOption {
	return func(c *retryConfig) {
		c.onRetry = f
	}
}

func Context(ctxt context.Context) RetryOption {
	return func(c *retryConfig) {
		c.ctxt = ctxt
	}
}

func BackOffStrategy(d DelayBackOffStrategy) RetryOption {
	return func(c *retryConfig) {
		c.backOffStrategy = d
	}
}

func Jitter(j bool) RetryOption {
	return func(c *retryConfig) {
		c.jitter = j
	}
}

func MaxDelay(t time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.maxDelay = t
	}
}

func Delay(t time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.delay = t
	}
}

func TimeProviderImpl(t TimeProvider) RetryOption {
	return func(c *retryConfig) {
		c.timeProvider = t
	}
}

func Randomizer(f RandomizerFunc) RetryOption {
	return func(c *retryConfig) {
		c.randomizer = f
	}
}

func DelayGenerator(f DelayFunc) RetryOption {
	return func(c *retryConfig) {
		c.delayGenerator = f
	}
}
