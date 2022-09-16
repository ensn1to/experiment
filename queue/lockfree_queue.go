package queue

import (
	"sync/atomic"
	"unsafe"
)

type LKQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node[T any] struct {
	value T
	next  unsafe.Pointer
}

// NewLKQueue an empty queue
func NewLKQueue[T any]() *LKQueue[T] {
	// empty node
	n := unsafe.Pointer(&node[T]{})
	return &LKQueue[T]{head: n, tail: n}
}

func (q *LKQueue[T]) Enqueue(v T) {
	n := &node[T]{value: v}
	for {
		tail := load[T](&q.tail)
		next := load[T](&tail.next)

		// tail and next is consistent?
		if tail == load[T](&q.tail) {
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n)
					return
				} else {
					cas(&q.tail, tail, next)
				}
			}
		}
	}
}

func (q *LKQueue[T]) Dequeue() T {
	var t T
	for {
		head := load[T](&q.head)
		tail := load[T](&q.tail)
		next := load[T](&head.next)

		if head == load[T](&q.head) {
			if head == tail {
				if next == nil {
					return t
				}
				cas(&q.tail, tail, next)
			} else {
				value := next.value
				if cas(&q.head, head, next) {
					return value
				}
			}
		}
	}
}

func load[T any](p *unsafe.Pointer) (n *node[T]) {
	return (*node[T])(atomic.LoadPointer(p))
}

func cas[T any](p *unsafe.Pointer, old, new *node[T]) (ok bool) {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
