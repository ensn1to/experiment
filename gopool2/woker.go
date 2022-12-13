package gopool

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var workerpool sync.Pool

var globalCnt uint32

type worker struct {
	pool *pool
}

func init() {
	workerpool.New = newWorker
}

func newWorker() any { return &worker{} }

func (w *worker) run() {
	go func() {
		for {
			var t *task
			w.pool.taskLock.Lock()
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCount, -1)
			}
			// 没有任何任务，退出
			if t == nil {
				w.close()
				w.pool.taskLock.Unlock()
				w.Recycle()
				return
			}
			w.pool.taskLock.Unlock()
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("GOPOOL: panic in pool: %s: %v: %s", w.pool.name, r, debug.Stack())
						if w.pool.panicHandler != nil {
							w.pool.panicHandler(t.ctx, r)
						}
					}
				}()

				t.f()
			}()

			t.Recycle()
		}
	}()
}

func (w *worker) close() { w.pool.deWokerCount() }

func (w *worker) zero() { w.pool = nil }

// 资源回收
func (w *worker) Recycle() {
	w.zero()
	workerpool.Put(w)
}
