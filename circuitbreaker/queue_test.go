package circuitbreaker

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	t.Run("push 2 items, pop gets first", func(t *testing.T) {
		expectedResult := rand.Int()
		q := queue[int]{nodes: make([]*node[int], 0)}
		q.push(expectedResult)
		q.push(math.MaxInt)
		assert.Equal(t, expectedResult, q.pop())
	})
}
