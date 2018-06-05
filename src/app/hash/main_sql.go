package main

import (
	"fmt"
	"log"
	"time"

	chash "app/hash"

	mgoutil "github.com/qiniu/db/mgoutil.v3"
)

var dbConfig = mgoutil.Config{
	Host:           "localhost:27017",
	DB:             "testImage",
	Mode:           "strong",
	SyncTimeoutInS: 5,
}

type Image struct {
	// ID    bson.ObjectId `bson:"_id"`
	Spec string
	Path string
	Data []byte
}

type person struct {
	name string
	age  int
}

func main() {
	var data []byte
	for i := 0; i < 2000; i++ {
		data = append(data, 147)
	}
	cHashRing := chash.NewConsistent(1000)

	for i := 0; i < 10; i++ {
		si := fmt.Sprintf("%d", i)
		cHashRing.Add(chash.NewNode(i, "172.18.1."+si, 8080, "host_"+si, 1, "image"+si))
	}

	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count := 1000000
	block := 10000

	t1 := time.Now() // get current time
	coll := session.DB(dbConfig.DB).C("all")
	var imageList []interface{}
	for i := 0; i < count; i++ {
		spec := fmt.Sprintf("spec%d", i)
		path := fmt.Sprintf("usr/local/images/image%d", i)
		// node := cHashRing.Get(spec)
		image := &Image{spec, path, data}
		imageList = append(imageList, image)

		// coll := session.DB(dbConfig.DB).C(node.Collection)
		if i%block == block-1 {
			err = coll.Insert(imageList...)
			if err != nil {
				log.Fatal(err)
			}
			imageList = imageList[:0]
		}
	}

	elapsed := time.Since(t1)
	fmt.Println("Time elapsed: ", elapsed)

}
