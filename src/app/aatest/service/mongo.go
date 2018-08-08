package service

import (
	"log"

	"gopkg.in/mgo.v2/bson"

	mgoutil "github.com/qiniu/db/mgoutil.v3"
)

var dbConfig = mgoutil.Config{
	Host:           "localhost:27017",
	DB:             "test_raw",
	Mode:           "strong",
	SyncTimeoutInS: 1,
}

//Person sdf
type Person struct {
	// ID    bson.ObjectId `bson:"_id"`
	Name  string //`bson:"x"`
	Phone string //`bson:"k"`
}

func insert(parameter *Parameter) {
	// session, err := mgo.Dial("") //传入数据库的地址，可以传入多个，具体请看接口文档

	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close() //用完记得关闭

	coll := session.DB(dbConfig.DB).C("Person")

	err = coll.Insert(parameter)
	if err != nil {
		log.Fatal(err)
	}
}

func get(id string) *Parameter {
	// session, err := mgo.Dial("") //传入数据库的地址，可以传入多个，具体请看接口文档

	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close() //用完记得关闭

	coll := session.DB(dbConfig.DB).C("Person")

	p := &Parameter{}
	err = coll.Find(bson.M{"name": id}).One(p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
