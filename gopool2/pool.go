package gopool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool interface {
	Name() string
	SetCap(cap int32)

	Go(func())
	CtxGo(ctx context.Context, f func())

	SetPanicHandler(f func(context.Context, any))

	WorkerCount() int32
}

var taskPool sync.Pool

func init() {
	taskPool.New = newTask
}

// ctx 主要是为了打日志的时候用，这样如果有 logid 的话调用链追踪可以查找到
type task struct {
	ctx context.Context

	f func()

	// 指向下一个任务的指针
	next *task
}

func (t *task) zero() {
	t.ctx = nil
	t.f = nil
	t.next = nil
}

func (t *task) Recycle() {
	t.zero()
	taskPool.Put(t)
}
func newTask() any { return &task{} }

type pool struct {
	name string

	cnt uint32 // create index
	cap int32  // capacity of the pool

	// 多链表改为单链表
	taskHead  *task
	taskTail  *task
	taskLock  sync.Mutex
	taskCount int32

	workerCount int32 // number of workers running

	cfg *Config

	panicHandler func(context.Context, any)
}

func NewPool(name string, cap int32, cfg *Config) Pool {
	p := &pool{
		name: name,
		cap:  cap,
		cfg:  cfg,
	}

	return p
}

func (p *pool) Name() string {
	return p.name
}

func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

func (p *pool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f

	p.taskLock.Lock()
	// add task to the tasklist's head
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()
	atomic.AddInt32(&p.taskCount, 1)

	if (atomic.LoadInt32(&p.taskCount) >= p.cfg.ScaleThreshold) &&
		(p.WorkerCount() < atomic.LoadInt32(&p.cap)) ||
		p.WorkerCount() == 0 {
		p.inWokerCount()
		w := workerpool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

func (p *pool) SetPanicHandler(f func(context.Context, any)) {
	p.panicHandler = f
}

func (p *pool) WorkerCount() int32 { return atomic.LoadInt32(&p.workerCount) }

func (p *pool) inWokerCount() { atomic.AddInt32(&p.workerCount, 1) }

func (p *pool) deWokerCount() { atomic.AddInt32(&p.workerCount, -1) }
