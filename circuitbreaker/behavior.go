package circuitbreaker

type behavior struct {
	metrics           metrics
	failureRatio      float64
	minimumThroughput int
}

func (b *behavior) onActionSuccess(CircuitState) {
	b.metrics.incrementSuccess()
}

func (b *behavior) onActionFailure(cs CircuitState) bool {
	shouldBreak := false
	switch cs {
	case Closed:
		b.metrics.incrementFailure()
		info := b.metrics.getHealthInfo()
		shouldBreak = info.throughput >= b.minimumThroughput && info.failureRate >= b.failureRatio
	case Open:
		// A failure call result may arrive when the circuit is open, if it was placed before the circuit broke.
		// We take no action beyond tracking the metric
		b.metrics.incrementFailure()
		shouldBreak = false
	case Isolated:
		b.metrics.incrementFailure()
		shouldBreak = false
	default:
		shouldBreak = false
	}

	return shouldBreak
}

func (b *behavior) onCircuitClosed() {
	b.metrics.reset()
}
