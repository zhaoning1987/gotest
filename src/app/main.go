package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type JsonLine struct {
	Image struct {
		Id   string `json:"id"`
		Uri  string `json:"uri"`
		Tag  string `json:"tag"`
		Desc string `json:"desc"`
	} `json:"image"`
}

type bb struct {
	aa string
}

func (t bb) String() string {
	return "abc"
}

func main() {
	file, _ := os.Open("/Users/zhaoning/Desktop/json2.json")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var obj JsonLine
		_ = json.Unmarshal([]byte(scanner.Text()), &obj)
		fmt.Println(obj)
	}

	fmt.Println("commit11111")
}
