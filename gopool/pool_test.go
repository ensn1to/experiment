package gopool

import (
	"testing"
)

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
