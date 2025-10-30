package main

import (
	"fmt"
	"net/http"
	"log"
	"bufio"
	"crypto/tls"
	"golang.org/x/net/http2"
	
	"os"
)

func serverHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("➡️  Received request for %s\n", path)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}


	go func() {
		scanner := bufio.NewScanner(r.Body)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("📨 Received stream message: %s\n", line)
			
		}
	}()

	writeScanner := bufio.NewScanner(os.Stdin)
	for writeScanner.Scan() {
		line := writeScanner.Text()
		fmt.Fprintf(w, "%s\n", line)	
		flusher.Flush()
	}

}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", serverHandler)
	
	server := &http.Server{
		Addr:    ":8443",
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,		},
	}

	http2.ConfigureServer(server, &http2.Server{})

	fmt.Println("🚀 HTTP/2 server running at https://localhost:8443")

	// Go’s standard lib enables HTTP/2 automatically for HTTPS
	err := server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		panic(err)
	}

}