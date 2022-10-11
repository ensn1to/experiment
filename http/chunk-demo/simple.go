package chunkdemo

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", chunkHandle)

	fmt.Println(http.ListenAndServe(":18080", nil))
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
