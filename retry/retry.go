// Package retry is a simple retry loop
package retry

import (
	"errors"
	"math"
	"time"
)

func Retry(f func() error, opts ...RetryOption) error {
	m := func() (any, error) { return nil, f() }
	_, err := RetryResult(m, opts...)
	return err
}

func RetryResult[TResult any](f func() (TResult, error), opts ...RetryOption) (TResult, error) {
	var emptyResult TResult
	var attempt uint = 0
	var errs error

	c := defaultConfig()

	for _, opt := range opts {
		opt(c)
	}

	for {
		select {
		case <-c.ctxt.Done():
			// handle the context being canceled
			return emptyResult, errors.Join(errs, c.ctxt.Err())
		default:
		}
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

		delay := getDelay(c)
		if delay > 0 {
			<-c.timeProvider.After(delay)
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

func getDelay(c *retryConfig) time.Duration {
	delay := c.delay

	if delay > c.maxDelay {
		delay = c.maxDelay
	}

	return delay
}
