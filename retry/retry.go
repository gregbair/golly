// Package retry is a simple retry loop
package retry

import (
	"errors"
	"math"
)

func Retry(f func() error, opts ...RetryOption) error {
	m := func() (any, error) { return nil, f() }
	_, err := RetryResult(m, opts...)
	return err
}

func RetryResult[TResult any](f func() (TResult, error), opts ...RetryOption) (TResult, error) {
	var attempt uint = 0

	c := defaultConfig()

	for _, opt := range opts {
		opt(c)
	}

	for {
		var errs error
		result, err := f()

		if err == nil {
			return result, nil
		}

		c.onRetry(attempt, err)

		errs = errors.Join(err)

		isLastAttempt, increment := isLastAttempt(attempt, c)
		if isLastAttempt {
			return result, errs
		}

		if increment {
			attempt++
		}
	}
}

func isLastAttempt(attempt uint, c *retryConfig) (isLastAttempt bool, increment bool) {
	if attempt == math.MaxUint {
		return false, false
	}
	return attempt >= c.attempts, true
}
