package retry

import "context"

type OnRetryFunc func(uint, error)

type retryConfig struct {
	attempts uint
	onRetry  OnRetryFunc
	ctxt     context.Context
}

type RetryOption func(*retryConfig)

func defaultConfig() *retryConfig {
	return &retryConfig{
		attempts: 3,
		onRetry:  func(u uint, err error) {},
		ctxt:     context.Background(),
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
