package log

import (
	"errors"

	"github.com/go-redis/redis"
)

// 为logger提供redis的writer io接口
type redisWriter struct {
	c       *redis.Client
	listKey string
}

// 实现writer方法
func (w *redisWriter) Write(p []byte) (int, error) {
	n, err := w.c.RPush(w.listKey, p).Result()
	return int(n), err
}

func NewRedisWriter(key string, addr string) (*redisWriter, error) {
	c := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if c == nil {
		return nil, errors.New("redis client init failed")
	}
	return &redisWriter{
		c:       c,
		listKey: key,
	}, nil
}
