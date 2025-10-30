package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"log"
	"bufio"
	"golang.org/x/net/http2"
	"os"
	"io")

func main() {
	// Configure TLS to skip certificate verification for demo purposes
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	
	// Create HTTP/2 Transport with the TLS config
	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{
		Transport: transport,
	}

	pr, pw := io.Pipe()

	req, err := http.NewRequest("POST", "https://localhost:8443/", pr)
	if err != nil {
		log.Printf("❌ Error fetching stream: %v\n", err)
		return
	}

	go func() {
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("❌ Error establishing stream: %v\n", err)
			return
		}
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("📨 Stream message: %s\n", line)
		}
		
	}()
	// Read from user input and send to server
	userScanner := bufio.NewScanner(os.Stdin)
	fmt.Println("💬 Type messages and press Enter to send (Ctrl+C to quit)")
	for userScanner.Scan() {
		line := userScanner.Text()
		io.WriteString(pw, line+"\n")
	}
}