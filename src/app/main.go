package main

import (
	"fmt"
)

type fff struct {
	Process []int
}

func main() {
	fmt.Println(len(ggg().Process))
}
func ggg() *fff {
	return &fff{}
}
