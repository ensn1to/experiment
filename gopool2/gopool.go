package gopool

import (
	"context"
	"fmt"
	"sync"
)

// defaultPool is the default global pool
var defaultPool Pool

func init() {
	defaultPool = NewPool("gopool.DefaultPool", 10000, NewConfig())
}

// Go is an alternative to the go key.
// gopool.Go(func(){}())
func Go(f func()) {
	defaultPool.Go(f)
}

func CtxGo(ctx context.Context, f func()) {
	defaultPool.CtxGo(ctx, f)
}

// SetCap is not recommended to be called,
// this func changes the global pool's capacity which will affect other callers.
func SetCap(cap int32) {
	defaultPool.ChangeCap(cap)
}

func SetPanicHandler(f func(context.Context, any)) {
	defaultPool.SetPanicHandler(f)
}

func WorkerCount() int32 {
	return defaultPool.WorkerCount()
}

var poolMap sync.Map

// RegisterPool registers a new pool to the global map.
func RegisterPool(p Pool) error {
	_, loaded := poolMap.LoadOrStore(p.Name(), p)
	if !loaded {
		return fmt.Errorf("pool:%s already register", p.Name())
	}

	return nil
}

// GetPool gets the registered pool by name.
func GetPool(name string) Pool {
	p, ok := poolMap.Load(name)
	if !ok {
		return nil
	}
	return p.(Pool)
}
