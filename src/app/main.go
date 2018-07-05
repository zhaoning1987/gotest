package main

import (
	"fmt"
)

type bb struct {
	s   string
	b   bool
	i   int
	Map map[string]struct{}
}

func (b *bb) print() {
	fmt.Println(b.s, b.b, b.i, b.Map)
}
func main() {
	var test bb
	test.print()
}
