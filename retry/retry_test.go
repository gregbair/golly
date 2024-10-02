package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoRetrySuccess(t *testing.T) {
	expectedResult := 5
	var attempts uint = 0
	result, err := RetryResult(func() (int, error) { return expectedResult, nil }, Attempts(3), OnRetry(func(u uint, err error) { attempts = u }))
	assert.Equal(t, uint(0), attempts)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestRetryOnErr(t *testing.T) {
	var expectedAttempts uint = 5
	var actualAttempts uint = 0
	err := Retry(func() error { return errors.New("foo") }, Attempts(expectedAttempts), OnRetry(func(u uint, err error) { actualAttempts = u }))
	assert.Error(t, err)
	assert.Equal(t, expectedAttempts, actualAttempts)
}

func TestContextCancellation(t *testing.T) {
	t.Run("cancel before execution", func(t *testing.T) {
		ctxt, cancel := context.WithCancel(context.Background())
		executed := false
		cancel()
		err := Retry(
			func() error { return errors.New("foo") },
			Attempts(5),
			OnRetry(func(u uint, e error) { executed = true }),
			Context(ctxt),
		)
		assert.Error(t, err)
		assert.Equal(t, "context canceled", err.Error())
		assert.False(t, executed)
	})

	t.Run("cancel during execution", func(t *testing.T) {
		ctxt, _ := context.WithTimeout(context.Background(), 1000*time.Millisecond)

		var expectedAttempts uint = 10
		var actualAttempts uint = 0
		err := Retry(
			func() error { return errors.New("foo") },
			Attempts(expectedAttempts),
			Context(ctxt),
			OnRetry(func(u uint, err error) {
				actualAttempts = u
				time.Sleep(500 * time.Millisecond)
			}),
		)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "context deadline exceeded")
		assert.Less(t, actualAttempts, expectedAttempts)

	})
}
