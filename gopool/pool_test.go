package gopool

import (
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p := New(10)

	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			time.Sleep(1 * time.Second)
		})
		if err != nil {
			println("task", i, "failed", err)
		}
	}

	p.Free()
}

func TestPoolWith(t *testing.T) {
	p := New(10, WithBlock(false), WithPreAlloc(false))
	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			time.Sleep(1 * time.Second)
		})
		if err != nil {
			println("task", i, "failed", err)
		}
	}

	p.Free()
}

func TestSchedule(t *testing.T) {
	p := New(10)
	if p.capacity != 10 {
		t.Errorf("want 10, got: %v", p.capacity)
	}

	err := p.Schedule(func() {})
	if err != nil {
		t.Errorf("want nil, got: %v", err)
	}

	p.Free()
	err = p.Schedule(func() {})
	if err == nil {
		t.Errorf("want no nil, got nil\n")
	}
}
