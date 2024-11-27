package circuitbreaker

type node[T any] struct {
	Value T
}

type queue[T any] struct {
	nodes []*node[T]
}

func (q *queue[T]) push(v T) {
	n := &node[T]{Value: v}
	q.nodes = append(q.nodes, n)
}

func (q *queue[T]) pop() T {
	if len(q.nodes) > 0 {
		n := q.nodes[0]
		q.nodes = q.nodes[1:]
		return n.Value
	}
	var empty T
	return empty
}

func (q *queue[T]) peek() T {
	if len(q.nodes) > 0 {
		return q.nodes[0].Value
	}
	var empty T
	return empty
}

func (q *queue[T]) clear() {
	q.nodes = make([]*node[T], 0)
}
