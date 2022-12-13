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
	l := len(w.pool.taskLists)
	go func() {
		for {
			var t *task
			for i := 0; i < l; i++ {
				idx := int(atomic.AddUint32(&globalCnt, 1)) % l
				w.pool.taskLists[idx].Lock()
				if w.pool.taskLists[idx].taskHead != nil {
					t = w.pool.taskLists[idx].taskHead
					w.pool.taskLists[idx].taskHead = w.pool.taskLists[idx].taskHead.next
					atomic.AddInt32(&w.pool.taskCount, -1)
					w.pool.taskLists[idx].Unlock()
					break
				} else {
					// 没有任务且是最后一次循环
					if i == l-1 {
						w.close()
						w.pool.taskLists[idx].Unlock()
						w.Recycle()
						return
					}
					w.pool.taskLists[idx].Unlock()
				}
			}
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
