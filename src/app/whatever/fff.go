package main

import (
	"fmt"
)

func func1() (fv string, num int) {
	return "hello", 1
}

func func2() (fv string, num int) {
	fv, num = func1()
	return
}

func main1() {
	fmt.Println(func2())
}
