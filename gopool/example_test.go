package gopool_test

import (
	"time"

	"github.com/ensn1to/experiment/tree/master/gopool"
)

func ExampleNew() {
	p := gopool.New(10)

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

func ExampleNewWithOpts() {
	p := gopool.New(10, gopool.WithBlock(false), gopool.WithPreAlloc(false))
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
