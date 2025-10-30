package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Printf("➡️  Received request for %s\n", path)

	// Simulate variable latency
	delay := time.Duration(500+1000*int64(len(path))) * time.Millisecond
	time.Sleep(delay)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Response for %s after %v\n", path, delay)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}

	fmt.Println("🚀 HTTP/2 server running at https://localhost:8443")
	// Go’s standard lib enables HTTP/2 automatically for HTTPS
	err := server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		panic(err)
	}
}
