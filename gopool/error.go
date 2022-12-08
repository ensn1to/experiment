package gopool

import "errors"

// 哨兵错误
var (
	ErrNoIdleWokerInPool = errors.New("no idle worker in pool")
	ErrPoolFreed         = errors.New("pool freed") // pool终止运行
)
