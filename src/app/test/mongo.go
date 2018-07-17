package main

import (
	"fmt"
	"log"

	mgoutil "github.com/qiniu/db/mgoutil.v3"
	"gopkg.in/mgo.v2/bson"
)

var dbConfig = mgoutil.Config{
	Host:           "localhost:27017",
	DB:             "test134",
	Mode:           "strong",
	SyncTimeoutInS: 1,
}

//Person sdf
type Person struct {
	// ID    bson.ObjectId `bson:"_id"`
	Name  string //`bson:"x"`
	Phone string //`bson:"k"`
}

func main() {
	// session, err := mgo.Dial("") //传入数据库的地址，可以传入多个，具体请看接口文档

	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close() //用完记得关闭

	coll := session.DB(dbConfig.DB).C("Person")

	//Optional. Switch the session to a monotonic behavior.
	// session.SetMode(mgo.Monotonic, true) //读模式，与副本集有关，详情参考https://docs.mongodb.com/manual/reference/read-preference/ & https://docs.mongodb.com/manual/replication/

	// c := session.DB("test").C("ning")
	err = coll.Insert(&Person{"ning", "+55 53 8116 9639"},
		&Person{"zhao", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	// result := Person{}
	var results []Person
	err = coll.Find(bson.M{"name": "ning"}).All(&results) //如果查询失败，返回“not found”
	if err != nil {
		log.Fatal(err)
	}
	// err = c.Find(bson.M{}).All(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("Phone:", result)
	// }

	fmt.Println("result:", results)

	condition := []string{}
	_, err = coll.UpdateAll(bson.M{"name": bson.M{"$in": condition}}, bson.M{"$set": bson.M{"name": "hahaha"}})
	if err != nil {
		log.Fatal(err)
	}
}
