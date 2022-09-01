package httpc

import (
	"context"
	"net/http"
	"time"
)

// 简单http client封装
type (
	HttpC struct {
		c *http.Client
	}

	httpConfig struct {
		host string
		port int
	}
)

func New(ctx context.Context) *HttpC {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &HttpC{
		c: &http.Client{
			Transport: tr,
		},
	}
}

func (h *HttpC) Do(req *http.Request) (*http.Response, error) {
	return h.c.Do(req)
}

func (h *HttpC) Close() error {
	return nil
}
