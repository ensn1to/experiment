package main

import (
	"fmt"
	"net/http"
)

func main() {
	s := http.Server{
		Addr: ":18082",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("receive a request from:", r.RemoteAddr, r.Header)
			w.Write([]byte("resp from server-2"))
		}),
	}
	s.ListenAndServe()
}
