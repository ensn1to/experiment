package queue

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestQueue(t *testing.T) {
	queues := map[string]Queue{
		"lock-free queue": NewLKQueue(),
		"slice queue":     NewSliceQueue(0),
	}

	for n, q := range queues {
		t.Run(n, func(t *testing.T) {
			count := 100
			for i := 0; i < count; i++ {
				q.Enqueue(i)
			}

			for i := 0; i < count; i++ {
				v := q.Dequeue()
				if v == nil {
					t.Fatalf("got a nil value")
				}

				if v.(int) != i {
					t.Fatalf("expect %d but got %d", i, v)
				}
			}
		})
	}

	assert.Equal(t, false, true)
}
