package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 180; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d\n", index)
		}(i)
	}

	wg.Wait()
}
