package retry

type OnRetryFunc func(uint, error)

type retryConfig struct {
	attempts uint
	onRetry  OnRetryFunc
}

type RetryOption func(*retryConfig)

func defaultConfig() *retryConfig {
	return &retryConfig{
		attempts: 3,
		onRetry:  func(u uint, err error) {},
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
