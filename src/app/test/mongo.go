package main

import (
	"log"

	mgoutil "github.com/qiniu/db/mgoutil.v3"
	"gopkg.in/mgo.v2/bson"
)

var dbConfig = mgoutil.Config{
	Host:           "localhost:27017",
	DB:             "test666",
	Mode:           "strong",
	SyncTimeoutInS: 1,
}

//Person sdf
type Person struct {
	// ID    bson.ObjectId `bson:"_id"`
	Name  string //`bson:"x"`
	Phone string //`bson:"k"`
}

type Face struct {
	Uid  string //`bson:"x"`
	Gid  string //`bson:"k"`
	Id   string
	Name string
	Desc string
}

func main() {
	// session, err := mgo.Dial("") //传入数据库的地址，可以传入多个，具体请看接口文档

	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close() //用完记得关闭

	coll := session.DB(dbConfig.DB).C("face")

	// c := session.DB("test").C("ning")
	// err = coll.Insert(&Person{"user1", "phone1"},
	// 	&Person{"user2", "phone2"},
	// 	&Person{"user3", "phone3"},
	// 	&Person{"user4", "phone3"},
	// 	&Person{"user5", "phone3"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = coll.EnsureIndex(mgo.Index{Key: []string{"uid", "gid", "id"}, Unique: true})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = coll.EnsureIndex(mgo.Index{Key: []string{"uid", "gid"}})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var faces []interface{}
	// for m := 1; m <= 1; m++ {

	// 	for i := 0; i < 500000; i++ {
	// 		face := &Face{
	// 			Uid:  "1",
	// 			Gid:  "group" + strconv.Itoa(m),
	// 			Id:   strconv.Itoa(i),
	// 			Name: "name" + strconv.Itoa(i),
	// 			Desc: "desc" + strconv.Itoa(i),
	// 		}

	// 		faces = append(faces, face)

	// 		if i%100 == 99 {
	// 			if err := coll.Insert(faces...); err != nil {
	// 				panic(err)
	// 			}
	// 			faces = faces[:0]
	// 		}
	// 	}
	// }

	// var results []Face
	// // err = coll.Find(nil).Skip(499000).Limit(1000).All(&results)                                   //如果查询失败，返回“not found”
	// err = coll.Find(bson.M{"uid": "1", "gid": "group160", "id": bson.M{"$gte": "11238"}}).Limit(10).All(&results) //如果查询失败，返回“not found”
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(len(results))
	// fmt.Println(results)

	// err = coll.Find(nil).Skip(499000).Limit(1000).All(&results)                                   //如果查询失败，返回“not found”
	err = coll.RemoveId(bson.ObjectIdHex("13123411234321423"))
	if err != nil {
		log.Fatal(err)
	}

	// result := Person{}
	// var results []Person
	// err = coll.Find(bson.M{"name": "ning"}).All(&results) //如果查询失败，返回“not found”
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("result:", results)

	// condition := []string{}
	// _, err = coll.UpdateAll(bson.M{"name": bson.M{"$in": condition}}, bson.M{"$set": bson.M{"name": "hahaha"}})
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
