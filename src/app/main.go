package main

import (
	"os"
	"time"
)

func main() {
	fd, _ := os.OpenFile("a.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd_time := time.Now().Format("2006-01-02 15:04:05\n")
	buf := []byte(fd_time)
	fd.Write(buf)
	fd.Close()
}
