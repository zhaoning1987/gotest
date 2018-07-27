package main

import (
	"fmt"
)

type stru struct {
	str string
}

type stru2 struct {
	stru
	str string
}

func main() {
	a := " "
	fmt.Println(len(a))
}

func test1() (res stru) {
	res = *test()
	return
}

func test() (res *stru) {
	res = &stru{}
	res.str = "aa"
	return
}
