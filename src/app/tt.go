package main

import (
	"fmt"
	"sync"
)

func main() {
	maxGoroutines := 5
	maxBusi := 10
	ch := make(chan bool, maxGoroutines)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < maxBusi; i++ {
		ch <- true // maxGoroutines allowed by this buffered channel
		go func(n int) {
			worker(n)
			<-ch // release the channel to let the rest start
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func worker(i int) {
	fmt.Println("doing work on", i)
}
