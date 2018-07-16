package main

import (
	"bufio"
	"fmt"

	mgoutil "github.com/qiniu/db/mgoutil.v3"
)

var dbConfig2 = mgoutil.Config{
	Host:           "localhost:27017",
	DB:             "test1",
	Mode:           "strong",
	SyncTimeoutInS: 1,
}

func read() {
	// session, err := mgo.Dial("") //传入数据库的地址，可以传入多个，具体请看接口文档

	session, err := mgoutil.Dail(dbConfig2.Host, dbConfig2.Mode, dbConfig2.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close() //用完记得关闭

	db := session.DB(dbConfig2.DB)

	file, err := db.GridFS("fs").Open("hello232.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// _, err = io.Copy(os.Stdout, file)
	// if err != nil {
	// 	panic(err)
	// }
	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func write() {
	session, err := mgoutil.Dail(dbConfig2.Host, dbConfig2.Mode, dbConfig2.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close() //用完记得关闭

	db := session.DB(dbConfig2.DB)

	file, err := db.GridFS("fs").Create("hello232.txt")

	if err != nil {
		panic(err)
	}

	file.Write([]byte("abcd 1234\nssdfsdf\nwf"))
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func delete() {
	session, err := mgoutil.Dail(dbConfig2.Host, dbConfig2.Mode, dbConfig2.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close() //用完记得关闭

	db := session.DB(dbConfig2.DB)

	err = db.GridFS("fs").Remove("hello232.txt")

	if err != nil {
		panic(err)
	}

}

type fileinfo struct {
	//文件大小
	LENGTH int32
	//md5
	MD5 string
	//文件名
	FILENAME string
}

func main() {
	read()
	// write()
	// delete()
}
