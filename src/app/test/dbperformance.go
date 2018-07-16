package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type person struct {
	Val string
}

func main3() {
	bb := "中国"
	aa := substring(bb, 1, 2)
	fmt.Println(aa)
	fmt.Println(bb[0:1])
	// session, err := mgoutil.Dail("localhost:27017", "strong", 1)
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()

	// coll := session.DB("md5").C("md5")
	// index := mgo.Index{
	// 	Key:        []string{"val"}, // 索引字段， 默认升序,若需降序在字段前加-
	// 	Unique:     true,            // 唯一索引 同mysql唯一索引
	// 	Background: false,           // 后台创建索引
	// }
	// if err := coll.EnsureIndex(index); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// data := gen()

	// block := 1000
	// var list []interface{}
	// start := time.Now()
	// for i := 0; i < 10000000; i++ {
	// 	p := &person{data[i]}
	// 	list = append(list, p)
	// 	if i%block == block-1 {
	// 		err = coll.Insert(list...)
	// 		if err != nil {
	// 			fmt.Println(i)
	// 			panic(err)
	// 		}
	// 		list = list[:0]
	// 	}
	// }
	// elapsed := time.Since(start)
	// fmt.Println("Time elapsed: ", elapsed)

	// h := md5.New()
	// h.Write([]byte(fmt.Sprintf("%d", 1000)))
	// cipherStr := h.Sum(nil)
	// source := hex.EncodeToString(cipherStr)

	// start2 := time.Now()
	//
	// elapsed2 := time.Since(start2)
	// fmt.Println("Time elapsed: ", elapsed2)

	//=======
	// h := md5.New()
	// h.Write([]byte(fmt.Sprintf("%d", 9999)))
	// cipherStr := h.Sum(nil)
	// source := hex.EncodeToString(cipherStr)
	// p1 := &person{source}
	// start3 := time.Now()
	// err = coll.Insert(p1)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "duplicate") {
	// 		fmt.Println("duplicate key")
	// 	}
	// }
	// elapsed3 := time.Since(start3)
	// fmt.Println("Time elapsed: ", elapsed3)
}

// 模拟生成一百万随机数
func gen() []string {
	start := time.Now()
	arr := make([]string, 10000000)
	for i := 0; i < 10000000; i++ {
		h := md5.New()
		h.Write([]byte(fmt.Sprintf("%d", i)))
		cipherStr := h.Sum(nil)
		arr[i] = hex.EncodeToString(cipherStr)
	}

	elapsed := time.Since(start)
	fmt.Println("Time elapsed: ", elapsed)
	return arr
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
