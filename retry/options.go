package retry

import (
	"context"
	"time"
)

type OnRetryFunc func(uint, error)

type DelayBackOffType int

const (
	Constant DelayBackOffType = iota
	Linear
)

type retryConfig struct {
	attempts     uint
	onRetry      OnRetryFunc
	ctxt         context.Context
	backoffType  DelayBackOffType
	jitter       bool
	delay        time.Duration
	maxDelay     time.Duration
	timeProvider TimeProvider
}

type RetryOption func(*retryConfig)

func defaultConfig() *retryConfig {
	return &retryConfig{
		attempts:     3,
		onRetry:      func(u uint, err error) {},
		ctxt:         context.Background(),
		backoffType:  Constant,
		jitter:       false,
		delay:        0,
		maxDelay:     5 * time.Second,
		timeProvider: systemTimeProvider{},
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

func BackOffType(d DelayBackOffType) RetryOption {
	return func(c *retryConfig) {
		c.backoffType = d
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
