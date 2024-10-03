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
	var attempts uint
	result, err := RetryResult(func() (int, error) { return expectedResult, nil }, Attempts(3), OnRetry(func(u uint, _ error) { attempts = u }))
	assert.Equal(t, uint(0), attempts)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestRetryOnErr(t *testing.T) {
	var expectedAttempts uint = 5
	var actualAttempts uint
	err := Retry(func() error { return errors.New("foo") }, Attempts(expectedAttempts), OnRetry(func(u uint, _ error) { actualAttempts = u }))
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
			OnRetry(func(uint, error) { executed = true }),
			Context(ctxt),
		)
		assert.Error(t, err)
		assert.Equal(t, "context canceled", err.Error())
		assert.False(t, executed)
	})

	t.Run("cancel during execution", func(t *testing.T) {
		ctxt, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)

		var expectedAttempts uint = 10
		var actualAttempts uint
		err := Retry(
			func() error { return errors.New("foo") },
			Attempts(expectedAttempts),
			Context(ctxt),
			OnRetry(func(u uint, _ error) {
				actualAttempts = u
				time.Sleep(500 * time.Millisecond)
			}),
		)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "context deadline exceeded")
		assert.Less(t, actualAttempts, expectedAttempts)
		cancel()
	})
}

func TestDelay(t *testing.T) {
	t.Run("test base delay", func(t *testing.T) {
		p := fakeProvider{}
		expectedDelay := 100 * time.Millisecond

		err := Retry(
			func() error { return errors.New("foo") },
			Attempts(2),
			Delay(expectedDelay),
			TimeProviderImpl(&p),
		)
		assert.Error(t, err)
		assert.Equal(t, expectedDelay, p.delay)
	})

	t.Run("test max delay", func(t *testing.T) {
		p := &fakeProvider{}
		expectedDelay := 100 * time.Millisecond

		err := Retry(
			func() error { return errors.New("foo") },
			Attempts(2),
			Delay(1*time.Second),
			MaxDelay(expectedDelay),
			TimeProviderImpl(p),
		)
		assert.Error(t, err)
		assert.Equal(t, expectedDelay, p.delay)
	})

	t.Run("no delay, provider not called", func(t *testing.T) {
		p := &failingProvider{t: t}

		_ = Retry(
			func() error { return errors.New("foo") },
			Attempts(2),
			Delay(0),
			TimeProviderImpl(p),
		)
	})

	t.Run("constant delay, does not change", func(t *testing.T) {
		p := &fakeProvider{}

		_ = Retry(
			func() error { return errors.New("foo") },
			Attempts(20),
			Delay(17*time.Millisecond),
			BackoffStrategy(Constant),
			TimeProviderImpl(p),
		)

		assert.Equal(t, 17*time.Millisecond, p.delay)
	})

	t.Run("Linear delay, grows linearly", func(t *testing.T) {
		p := &fakeProvider{}

		_ = Retry(
			func() error { return errors.New("foo") },
			Attempts(2),
			Delay(1*time.Millisecond),
			BackoffStrategy(Linear),
			TimeProviderImpl(p),
		)

		assert.Equal(t, 2*time.Millisecond, p.delay)
	})

	t.Run("delay generator respected", func(t *testing.T) {
		p := &fakeProvider{}
		expectedDelay := 12 * time.Millisecond

		_ = Retry(
			func() error { return errors.New("foo") },
			Attempts(2),
			Delay(72*time.Millisecond),
			TimeProviderImpl(p),
			DelayGenerator(func(uint, context.Context) time.Duration { return expectedDelay }),
		)

		assert.Equal(t, expectedDelay, p.delay)
	})
}

type fakeProvider struct {
	delay time.Duration
}

func (f *fakeProvider) After(d time.Duration) <-chan time.Time {
	f.delay = d
	return time.After(0)
}

type failingProvider struct {
	t *testing.T
}

func (f *failingProvider) After(time.Duration) <-chan time.Time {
	assert.FailNow(f.t, "AFter should not have been called")
	return nil
}
