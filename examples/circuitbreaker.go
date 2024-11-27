package examples

import "github.com/gregbair/golly/circuitbreaker"

func doThings() {
	_ = circuitbreaker.New[int]()
}
