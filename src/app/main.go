package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup
var f *os.File

func main() {
	t1 := time.Now()
	f, err := os.OpenFile("/Users/zhaoning/Desktop/multi", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go saveErrorToFile(f, i)
	}
	wg.Wait()

	fmt.Println(time.Since(t1))
}

func saveErrorToFile(file *os.File, i int) {

	time.Sleep(time.Duration(rand.Intn(5)) * 100 * time.Millisecond)
	content := fmt.Sprintf("%d\n", i)
	buf := []byte(content)
	file.Write(buf)
	wg.Done()

}
