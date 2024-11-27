package circuitbreaker

import (
	"time"

	"github.com/gregbair/golly"
)

const numberOfWindows = 10

type healthInfo struct {
	throughput   int
	failureRate  float64
	failureCount int
}

func createHealthInfo(successes int, failures int) healthInfo {
	total := successes + failures
	if total == 0 {
		return healthInfo{0, 0, failures}
	}
	return healthInfo{total, float64(failures / total), failures}
}

type metrics interface {
	incrementSuccess()
	incrementFailure()
	reset()
	getHealthInfo() healthInfo
}

func createHealthMetrics(samplingDuration time.Duration, timeProvider golly.TimeProvider) metrics {
	if samplingDuration < numberOfWindows {
		return newSingleMetrics(samplingDuration, timeProvider)
	}
	return newRollingMetrics(samplingDuration, timeProvider)
}

type singleMetrics struct {
	samplingDuration time.Duration
	failures         int
	successes        int
	startedAt        time.Time
	timeProvider     golly.TimeProvider
}

func newSingleMetrics(samplingDuration time.Duration, timeProvider golly.TimeProvider) *singleMetrics {
	return &singleMetrics{
		samplingDuration: samplingDuration,
		timeProvider:     timeProvider,
	}
}

func (s *singleMetrics) incrementSuccess() {
	s.tryReset()
	s.successes++
}

func (s *singleMetrics) incrementFailure() {
	s.tryReset()
	s.failures++
}

func (s *singleMetrics) reset() {
	s.startedAt = s.timeProvider.UtcNow()
	s.successes = 0
	s.failures = 0
}

func (s *singleMetrics) getHealthInfo() healthInfo {
	s.tryReset()
	return createHealthInfo(s.successes, s.failures)
}

func (s *singleMetrics) tryReset() {
	if s.timeProvider.UtcNow().Sub(s.startedAt) >= s.samplingDuration {
		s.reset()
	}
}

type rollingMetrics struct {
	timeProvider     golly.TimeProvider
	samplingDuration time.Duration
	windowDuration   time.Duration
	windows          *queue[*healthWindow]
	currentWindow    *healthWindow
}

func newRollingMetrics(samplingDuration time.Duration, timeProvider golly.TimeProvider) *rollingMetrics {
	return &rollingMetrics{
		samplingDuration: samplingDuration,
		timeProvider:     timeProvider,
		windows:          &queue[*healthWindow]{nodes: make([]*node[*healthWindow], 0)},
		windowDuration:   samplingDuration / numberOfWindows,
	}
}

func (r *rollingMetrics) incrementSuccess() {
	r.updateCurrentWindow().successes++
}

func (r *rollingMetrics) incrementFailure() {
	r.updateCurrentWindow().failures++
}

func (r *rollingMetrics) reset() {
	r.currentWindow = nil
	r.windows.clear()
}

func (r *rollingMetrics) getHealthInfo() healthInfo {
	_ = r.updateCurrentWindow()
	successes := 0
	failures := 0

	for _, window := range r.windows.nodes {
		successes += window.Value.successes
		failures += window.Value.failures
	}

	return createHealthInfo(successes, failures)
}

func (r *rollingMetrics) updateCurrentWindow() *healthWindow {
	now := r.timeProvider.UtcNow()
	if r.currentWindow == nil || now.Sub(r.currentWindow.startedAt) >= r.windowDuration {
		r.currentWindow = &healthWindow{startedAt: now}
		r.windows.push(r.currentWindow)
	}

	for now.Sub(r.windows.peek().startedAt) >= r.samplingDuration {
		_ = r.windows.pop()
	}
	return r.currentWindow
}

type healthWindow struct {
	successes int
	failures  int
	startedAt time.Time
}
