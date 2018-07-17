package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func change(i int) {
	i = 3
}
func main() {
	var count int64
	count = 0
	var wg sync.WaitGroup
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			atomic.AddInt64(&count, 1)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(count)
}
