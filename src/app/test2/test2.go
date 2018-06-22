package test2

import (
	"fmt"
)

type ServeMux struct {
	S string
	I int
}

func (h *ServeMux) HandleFunc(s int) {
	fmt.Println(h.I)
}

func (h *ServeMux) ServeHTTP(s string) {
	fmt.Println(h.S)
}
