// Package retry is a simple retry loop
package retry

import (
	"errors"
	"math"
	"time"
)

const jitterFactor float64 = 0.5

func Retry(f func() error, opts ...RetryOption) error {
	_, err := RetryResult(func() (any, error) { return nil, f() }, opts...)
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

		errs = errors.Join(errs, err)

		isLastAttempt, increment := isLastAttempt(attempt, c)
		if isLastAttempt {
			return result, errs
		}

		delay := getDelay(c, attempt)
		if c.delayGenerator != nil {
			newDelay := c.delayGenerator(attempt, c.ctxt)
			if newDelay >= 0 {
				delay = newDelay
			}
		}

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

func getDelay(c *retryConfig, attempt uint) time.Duration {
	var delay time.Duration

	switch c.backOffStrategy {
	case Linear:
		delay = time.Duration(int64(attempt+1)*c.delay.Milliseconds()) * time.Millisecond
	case Constant:
		delay = c.delay
	default:
		return c.delay
	}

	if c.jitter {
		delay = applyJitter(delay, c)
	}

	if delay > c.maxDelay {
		delay = c.maxDelay
	}

	return delay
}

func applyJitter(d time.Duration, c *retryConfig) time.Duration {
	offset := (float64(d.Milliseconds()) * jitterFactor) / 2
	randomDelay := (float64(d.Milliseconds()) * jitterFactor * c.randomizer()) - offset
	newDelay := float64(d.Milliseconds()) + randomDelay

	return time.Duration(newDelay * float64(time.Millisecond))
}
