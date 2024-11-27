package circuitbreaker

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultNoErrors(t *testing.T) {
	assert.NoError(t, validateConfig(defaultConfig[any]()))
}

func TestValidationFailureRatio(t *testing.T) {
	type cases struct {
		description string
		input       float64
		isError     bool
	}

	c := defaultConfig[any]()

	for _, tt := range []cases{
		{"at min", 0.0, false},
		{"at max", 1.0, false},
		{"mid", 0.5, false},
		{"too low", -0.1, true},
		{"too high", 1.1, true},
	} {
		FailureRatio(tt.input)(c)

		err := validateConfig(c)

		if tt.isError {
			assert.ErrorContains(t, err, "failure ratio")
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidationMinimumThroughput(t *testing.T) {
	type cases struct {
		description string
		input       int
		isError     bool
	}

	c := defaultConfig[any]()

	for _, tt := range []cases{
		{"at min", 2, false},
		{"too low", 1, true},
		{"at max", math.MaxInt, false},
	} {
		MinimumThroughput(uint(tt.input))(c)

		err := validateConfig(c)

		if tt.isError {
			assert.ErrorContains(t, err, "minimum throughput")
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidationSamplingDuration(t *testing.T) {
	type cases struct {
		description string
		input       time.Duration
		isErr       bool
	}

	c := defaultConfig[any]()

	for _, tt := range []cases{
		{"at min", 500 * time.Millisecond, false},
		{"at max", 1 * time.Hour, false},
		{"too low", 150 * time.Millisecond, true},
		{"too high", 2 * time.Hour, true},
	} {
		SamplingDuration(tt.input)(c)
		err := validateConfig(c)

		if tt.isErr {
			assert.ErrorContains(t, err, "sampling duration")
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateBreakDuration(t *testing.T) {
	type cases struct {
		description string
		input       time.Duration
		isErr       bool
	}

	c := defaultConfig[any]()

	for _, tt := range []cases{
		{"at min", 500 * time.Millisecond, false},
		{"at max", 1 * time.Hour, false},
		{"too low", 150 * time.Millisecond, true},
		{"too high", 2 * time.Hour, true},
	} {
		BreakDuration(tt.input)(c)
		err := validateConfig(c)

		if tt.isErr {
			assert.ErrorContains(t, err, "break duration")
		} else {
			assert.NoError(t, err)
		}
	}
}
