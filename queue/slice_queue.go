package queue

import "sync"

type SliceQueue struct {
	data []interface{}
	mu   sync.Mutex
}

func NewSliceQueue(n int) *SliceQueue {
	return &SliceQueue{data: make([]interface{}, n)}
}

func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

func (q *SliceQueue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.data) == 0 {
		return nil
	}

	v := q.data[0]
	q.data = q.data[1:]
	return v
}
