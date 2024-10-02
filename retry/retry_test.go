package retry

import (
	"errors"
	"testing"

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
