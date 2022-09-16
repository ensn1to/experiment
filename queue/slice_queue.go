package queue

import "sync"

type SliceQueue[T any] struct {
	data []T
	mu   sync.Mutex
}

func NewSliceQueue[T any](n int) *SliceQueue[T] {
	return &SliceQueue[T]{data: make([]T, 0, n)}
}

func (q *SliceQueue[T]) Enqueue(v T) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

func (q *SliceQueue[T]) Dequeue() T {
	var t T
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.data) == 0 {
		return t
	}

	v := q.data[0]
	q.data = q.data[1:]
	return v
}
