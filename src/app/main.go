package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("hello ning1111111!")
}

func main1() {
	file, err := os.OpenFile("/Users/zhaoning/Desktop/imageList", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	err = filepath.Walk("/Users/zhaoning/Documents/tomcatimage", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		//skip directory & windows thumb file & linux max hidden file
		if f.IsDir() || strings.ToLower(f.Name()) == "thumb.db" || substring(f.Name(), 0, 1) == "." {
			return nil
		}

		content := fmt.Sprintf("http://localhost:8080/test/%s\n", f.Name())
		buf := []byte(content)
		file.Write(buf)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

}

func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}
