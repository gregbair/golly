/*
Package retry is a retry provider providing constant and linear backoff strategies, configurable delay, and more.

# Key Features

  - Context-aware: The function respects a provided context and can handle cancellation.
  - Customizable: The RetryOption functions allow for flexible configuration of retry behavior.
  - Error handling: The function accumulates errors and returns them if all retries fail.
  - Delay strategies: Supports various delay strategies, including fixed delays and backoff algorithms.
  - Callbacks: Provides callbacks for handling retries and customizing behavior.

# Example usage

The following is a basic example of using retry with no return value. In this
example, the operation will be retried 5 times and then all errors will be returned.

		func doSomething() error {return errors.New("foo")}

	    err:= Retry(
			doSomething,
			Attempts(5)
		)
*/
package retry

import (
	"errors"
	"math"
	"time"
)

const jitterFactor float64 = 0.5

/*
Retry retries a given function f a specified number of times or until a timeout occurs, handling errors that may arise during execution.

# Parameters
  - f: A function that returns an error.
  - opts...: A variable number of RetryOption functions that can be used to customize the retry behavior, such as setting the maximum number of
    attempts, specifying a delay strategy, and providing callbacks for handling retries.

# Return Values

  - error: An error if all retries fail or the context is canceled.

# Example Usage

	func doSomething() error {return errors.New("foo")}

	err:= Retry(
	    doSomething,
	    Attempts(5)
	)
*/
func Retry(operation func() error, opts ...RetryOption) error {
	_, err := RetryResult(func() (any, error) { return nil, f() }, opts...)
	return err
}

/*
RetryResult retries a given function f a specified number of times
or until a timeout occurs, handling errors that may arise during execution.

# Parameters

  - f: A function that returns a result of type TResult and an error.
  - opts...: A variable number of RetryOption functions that can be used to
    customize the retry behavior, such as setting the maximum number of
    attempts, specifying a delay strategy, and providing callbacks for
    handling retries.

# Return Values

  - TResult: The result of the function f if it succeeds.
  - error: An error if all retries fail or the context is canceled.

# Example Usage

	func fetchData() (string, error) {
	    // ...
	}

	result, err := RetryResult(fetchData, RetryMaxAttempts(3), RetryDelay(time.Second))
	if err != nil {
	    // Handle error
	}
*/
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
