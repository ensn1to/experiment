package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	s := http.Server{
		Addr: ":18081",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("receive a request from:", r.RemoteAddr, r.Header)
			w.Write([]byte("ok"))
		}),
		// 设置空闲连接时间
		IdleTimeout: 5 * time.Second,
	}

	s.ListenAndServe()
}
