package gopool

import (
	"fmt"
	"sync"
)

type Pool struct {
	capacity int // the size of pool

	// 当pool满的情况下，新的Schedule调用是否阻塞当前goroutine。默认值：true
	block bool

	preAlloc bool // 是否在创建pool的时候，就预创建workers，默认值为：false

	active chan struct{} // 当前pool可用worker available worker in pool
	tasks  chan Task     // 待执行task task to execute

	wg   sync.WaitGroup // 用于pool销毁时等待所有worker退出
	quit chan struct{}  // 用于通知所有worker退出
}

type Task func()

const (
	defaultCapacity = 10
	maxCapacity     = 30
)

func New(capacity int, opts ...Option) *Pool {
	// 防御性校验
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		block:    true,
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
		active:   make(chan struct{}, capacity),
	}

	for _, opt := range opts {
		opt(p)
	}

	// create all workers
	if p.preAlloc {
		for i := 0; i < p.capacity; i++ {
			p.newWorker(i + 1)
			p.active <- struct{}{}
		}
	}
	fmt.Println("gopool start")

	go p.run()

	return p
}

func (p *Pool) Cap() int {
	return p.capacity
}

func (p *Pool) returnTask(t Task) {
	go func() {
		p.tasks <- t
	}()
}

func (p *Pool) run() {
	// 当前pool中active worker
	idx := len(p.active)

	// 根据实际tasks channel中的情况创建worker
	// 如果直接执行下面的for{}，满足的是preAlloc设置，即active pool 一直为满
	if !p.preAlloc {
	loop:
		for t := range p.tasks {
			// 调度循环，目的为了下面的worker创建
			p.returnTask(t)
			select {
			case <-p.quit:
				return
			case p.active <- struct{}{}:
				idx++
				p.newWorker(idx)
			default:
				// active pool满，结束loop
				break loop
			}
		}
	}

	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}: // 添加计数，新增一个G(worker)
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(i int) {
	p.wg.Add(1)
	// 每个worker单独一个G
	go func() {
		// 管理当前的G(worker)
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
				// 释放信号量
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d] starting\n", i)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d] exiting\n", i)
				<-p.active // 释放信号量
				return
			case t := <-p.tasks: // recv from tasks, and execute
				fmt.Printf("worker[%03d]: recv a task\n", i)
				t()
			}
		}
	}()
}

// Schedule adds a task to the pool. Return the ErrPoolFreed if the pool is freed
// The tasks is non-buffer channel,
// it will be blocked when no more idle worker in pool
func (p *Pool) Schedule(task Task) error {
	for {
		select {
		case <-p.quit:
			return ErrPoolFreed
		case p.tasks <- task:
			return nil
		default: // 防止阻塞
			if p.block {
				p.tasks <- task
				return nil
			}
			return ErrNoIdleWokerInPool
		}
	}
}

func (p *Pool) Free() {
	close(p.quit) // make sure all workers are closed
	p.wg.Wait()
}
