package circuitbreaker

// CircuitState describes the possible states the circuit of a circuitbreaker may be in
type CircuitState int

const (
	// Closed - When the circuit is closed. Execution of actions is allowed.
	Closed CircuitState = iota

	// Open - When the automated controller has opened the circuit (typically due to some failure threshold being exceeded by recent actions). Execution of actions is blocked.
	Open

	// HalfOpen - When the circuit is half-open, it is recovering from an open state.
	// The duration of break of the preceding open state has typically passed.
	// In the half-open state, actions may be executed, but the results of these actions may be treated with criteria different to normal operation,
	// to decide if the circuit has recovered sufficiently to be placed back in to the closed state,
	// or if continuing failures mean the circuit should revert to open perhaps more quickly than in normal operation.
	HalfOpen

	// Isolated - When the circuit has been placed into a fixed open state by the isolate call.
	// This isolates the circuit manually, blocking execution of all actions until a reset call is made.
	Isolated
)

// Breaker is the main entry point into the breaker pipeline
type Breaker[TResult any] struct {
	c *controller[TResult]
}

// New instantiates a new Breaker
func New[TResult any](opts ...Option[TResult]) *Breaker[TResult] {
	c := defaultConfig[TResult]()
	for _, o := range opts {
		o(c)
	}

	b := &behavior{
		metrics:           createHealthMetrics(c.samplingDuration, c.timeProvider),
		failureRatio:      c.failureRatio,
		minimumThroughput: int(c.minimumThroughput),
	}

	con := &controller[TResult]{
		breakDuration: c.breakDuration,
		onOpened:      c.onOpened,
		onClosed:      c.onClosed,
		onHalfOpen:    c.onHalfOpen,
		behavior:      b,
		timeProvider:  c.timeProvider,
	}

	return &Breaker[TResult]{c: con}
}
