package gopool_test

import (
	"sync"
	"testing"

	"github.com/ensn1to/experiment/tree/master/gopool"
)

func TestSchedule(t *testing.T) {
	p := gopool.New(10)
	if p.Cap() != 10 {
		t.Errorf("want 10, got: %v", p.Cap())
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

const benchmarkTimes = 10000

func testFunc() {
	panic("panic test")
}

func BenchmarkSchedule(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	p := gopool.New(10000)

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(benchmarkTimes)
		for j := 0; j < benchmarkTimes; j++ {
			go func() {
				p.Schedule(testFunc)
				wg.Done()
			}()
		}

		wg.Wait()

	}
}
