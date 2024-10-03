# Golly

Golly is a suite of resiliency functions that allow you to build safer code.

It is somewhat inspired by [Polly](https://www.thepollyproject.org/)

# Retry

Retry allows you to continuously attempt an operation until either success or the associated number of retries (or a timeout) occurs.

## Retry example

You can either use retry with or without a return value.

### Without return value

```go
import "github.com/gregbair/golly/retry"

func postData() error {
    // do something
}

err := retry.Retry(
    postData,
    Attempts(5),
    )

if err != nil {
    // handle errors
}
```

### With Return value
```go
import "github.com/gregbair/golly/retry"

func fetchData() (string, error) {
    // do something
}

result, err := retry.RetryResult(
    fetchData,
    Attempts(3),
)
```

Configuration details are provided in comments in the [package](retry/options.go)