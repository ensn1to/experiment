package log

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ensn1to/experiment/tree/master/log/elog/pkg/httpc"
	"go.uber.org/zap"
)

func init() {
	// 这里注册http sink
	err := zap.RegisterSink("http", httpSink)
	if err != nil {
		fmt.Println("Register http sink fail", err)
	}
}

func httpSink(url *url.URL) (zap.Sink, error) {
	return &Http{
		// httpc是我封装的httpClient，没什么别的逻辑，直接当成http.Client就好
		httpc: httpc.New(context.Background()),
		url:   url,
	}, nil
}

type Http struct {
	httpc *httpc.HttpC
	url   *url.URL
}

// 主要逻辑就是Write
func (h *Http) Write(p []byte) (n int, err error) {
	// 初始化request
	req, err := http.NewRequest("POST", h.url.String(), bytes.NewReader(p))
	if err != nil {
		return 0, err
	}

	// 执行http请求
	resp, err := h.httpc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 获取response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return 0, errors.New(hex.EncodeToString(respBody))
	}

	return len(p), nil
}

// 可以搞个内置的queue或者[]log，在Sync函数里用来做批量发送提升性能，这里只是简单的实现，所以Sync没什么逻辑
func (h *Http) Sync() error {
	return nil
}

func (h *Http) Close() error {
	return h.httpc.Close()
}
