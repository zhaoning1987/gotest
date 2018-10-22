package main

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	mgoutil "github.com/qiniu/db/mgoutil.v3"
)

func BenchmarkHello(b *testing.B) {
	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close() //用完记得关闭

	coll := session.DB(dbConfig.DB).C("face")

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		// for j := 0; j < 500; j++ {
		// n := rand.Intn(500)
		// m := rand.Intn(200) + 1
		// n := 199
		m := 123
		k := "120"
		var results []Face
		// err := coll.Find(bson.M{"uid": "1", "gid": "group" + strconv.Itoa(m)}).Skip(n * 1000).Limit(1000).All(&results) //如果查询失败，返回“not found”
		err := coll.Find(bson.M{"uid": "1", "gid": "group" + strconv.Itoa(m), "id": bson.M{"$gt": k}}).Sort("id").Limit(1000).All(&results) //如果查询失败，返回“not found”
		if err != nil {
			log.Fatal(err)
		}
		// }
	}
}
