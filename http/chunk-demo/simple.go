package chunkdemo

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/chunk", chunkHandle)
	http.HandleFunc("/", connHandler)

	fmt.Println(http.ListenAndServe(":18080", nil))
}

func connHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Connection", "keep-alive")
	w.Write([]byte("connect back"))
}

func chunkHandle(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "Chunk #%d\n", i)
		flusher.Flush() // Trigger "chunked" encoding and send a chunk...
		time.Sleep(500 * time.Millisecond)
	}
}
