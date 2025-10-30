package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // local testing
	}
	client := &http.Client{Transport: tr}

	urls := []string{
		"https://localhost:8443/a",
		"https://localhost:8443/b",
		"https://localhost:8443/c",
		"https://localhost:8443/d",
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			resp, err := client.Get(u)
			if err != nil {
				fmt.Println("❌ Error:", err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("✅ %s -> %s", u, string(body))
		}(url)
	}

	wg.Wait()
}
