package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var requestCount uint64
var memoryStore [][]byte

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&requestCount, 1)
		// Allocate and retain memory to simulate load (1MB per request)
		data := make([]byte, 3*1024*1024)
		for i := range data {
			data[i] = byte(i % 256)
		}
		memoryStore = append(memoryStore, data)

		// Keep memory for a while then release oldest
		if len(memoryStore) > 3 {
			memoryStore = memoryStore[1:]
		}

		time.Sleep(100 * time.Millisecond)
		fmt.Fprintf(w, "Pod ID: %s\n", hostname)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		count := atomic.LoadUint64(&requestCount)
		fmt.Fprintf(w, "# HELP http_requests_total The total number of HTTP requests.\n")
		fmt.Fprintf(w, "# TYPE http_requests_total counter\n")
		fmt.Fprintf(w, "http_requests_total %d\n", count)
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
