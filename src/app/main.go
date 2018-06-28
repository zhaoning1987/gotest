package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CsvLine struct {
	Column1 string
	Column2 string
	Column3 string
}

func main() {

	filename := "/Users/zhaoning/Desktop/Book1.csv"

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	// Loop through lines & turn into object
	for _, line := range lines {
		data := CsvLine{
			Column1: line[0],
			Column2: line[1],
			Column3: line[2],
		}
		if data.Column1 == " aaa" {
			fmt.Println("hello")
		}
		fmt.Println(data.Column1, data.Column2, data.Column3)
	}
}
