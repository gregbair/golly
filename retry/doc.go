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
