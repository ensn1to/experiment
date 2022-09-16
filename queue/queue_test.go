package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queues := map[string]Queue[int]{
		"lock-free queue": NewLKQueue[int](),
		"slice queue":     NewSliceQueue[int](0),
	}

	for n, q := range queues {
		t.Run(n, func(t *testing.T) {
			count := 100
			for i := 1; i <= count; i++ {
				q.Enqueue(i)
			}

			for i := 1; i <= count; i++ {
				v := q.Dequeue()
				if v == 0 {
					t.Fatalf("got a 0 value")
				}

				if v != i {
					t.Fatalf("expect %d but got %d", i, v)
				}
			}
		})
	}
}
