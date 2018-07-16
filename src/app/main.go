package main

import "os"

type abc struct {
	aa string
	bb string
}

func (a abc) String() string {
	return "ssss"
}

func main() {

	_, err := os.OpenFile("/Users/zhaoning/Desktop/bbbb", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	_, err = os.OpenFile("/Users/zhaoning/Desktop/bbbb", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	_, err = os.OpenFile("/Users/zhaoning/Desktop/bbbb", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
}
