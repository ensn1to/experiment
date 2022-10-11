package main

import (
	"fmt"
	"net/http"
)

func main() {
	s := http.Server{
		Addr: ":18081",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("receive a request from:", r.RemoteAddr, r.Header)
			w.Write([]byte("ok"))
		}),
	}
	s.ListenAndServe()
}
