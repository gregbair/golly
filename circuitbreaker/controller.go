package circuitbreaker

import (
	"context"
	"time"

	"github.com/gregbair/golly"
)

type controller[TResult any] struct {
	breakDuration time.Duration
	onOpened      func(context.Context, time.Duration, TResult, error)
	onClosed      func(context.Context, TResult, error)
	onHalfOpen    func(context.Context)
	behavior      *behavior
	timeProvider  golly.TimeProvider
}
