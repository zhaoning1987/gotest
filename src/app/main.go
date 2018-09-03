package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

type ss struct {
	a string
}
type ss2 struct {
	param ss
}

func main() {
	str := "/sdf/sd"
	bb := strings.TrimSuffix(str, "/") + "/"
	fmt.Println(bb)
}
func main2() {
	count := 10000
	bar := pb.StartNew(count)
	// bar.Postfix("test")
	bar.ShowSpeed = true
	for i := 0; i < count; i++ {
		bar.Increment()

		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("The End!")
}
func main1() {
	var s [3]ss
	s[0] = ss{"1"}
	s[1] = ss{"2"}
	s[2] = ss{"3"}
	fmt.Println(s)

	b := ss2{param: s[0]}
	s[0].a = "aaa"
	fmt.Println(b)

	path2 := "/Users/zhaoning/Desktop/"
	path2 = strings.Replace(path2, "/", ":", -1)
	path := "./%s"
	err := os.MkdirAll(fmt.Sprintf(path, path2), os.ModePerm)
	if err != nil {
		fmt.Printf("create directory [%s] failed: %v", path, err)
	}
}
