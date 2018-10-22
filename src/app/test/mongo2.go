package main

// import (
// 	"fmt"
// 	"log"
// 	"sync"
// 	"sync/atomic"

// 	mgoutil "github.com/qiniu/db/mgoutil.v3"
// 	"gopkg.in/mgo.v2/bson"
// )

// var dbConfig3 = mgoutil.Config{
// 	Host:           "localhost:27017",
// 	DB:             "test333",
// 	Mode:           "strong",
// 	SyncTimeoutInS: 1,
// }

// //Person sdf
// type Person struct {
// 	// ID    bson.ObjectId `bson:"_id"`
// 	Name  string   //`bson:"x"`
// 	Phone []string //`bson:"k"`'
// 	Count int
// }

// func main() {
// 	// session, err := mgo.Dial("") //传入数据库的地址，可以传入多个，具体请看接口文档

// 	session, err := mgoutil.Dail(dbConfig3.Host, dbConfig3.Mode, dbConfig3.SyncTimeoutInS)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close() //用完记得关闭

// 	coll := session.DB(dbConfig3.DB).C("Person")

// 	//Optional. Switch the session to a monotonic behavior.
// 	// session.SetMode(mgo.Monotonic, true) //读模式，与副本集有关，详情参考https://docs.mongodb.com/manual/reference/read-preference/ & https://docs.mongodb.com/manual/replication/

// 	// c := session.DB("test").C("ning")
// 	err = coll.Insert(&Person{"ning", []string{"+55 53 8116 9639", "sdfsdf"}, 0},
// 		&Person{"zhao", []string{"+55 53 8116 9639", "wwwwwww"}, 0})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	coll.Update(bson.M{"name": "ning"},
// 		bson.M{"$push": bson.M{"phone": "a new phone"}})

// 	// result := Person{}
// 	var results Person
// 	err = coll.Find(bson.M{"name": "zhao"}).One(&results) //如果查询失败，返回“not found”
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("result:", results)

// 	err = coll.Find(bson.M{"name": "ning"}).One(&results) //如果查询失败，返回“not found”
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("result:", results)

// 	var projection struct {
// 		Phone []string `bson:"phone"`
// 	}
// 	err = coll.Find(bson.M{"name": "ning"}).Select(bson.M{"phone": 1}).One(&projection)
// 	// for _, v := range projection {
// 	fmt.Println(projection.Phone)
// 	// }
// 	// err = c.Find(bson.M{}).All(&result)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// } else {
// 	// 	fmt.Println("Phone:", result)
// 	// }

// 	err = coll.Insert(&Person{"concurrTest", []string{}, 0})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var count int64
// 	count = 0
// 	// var lock sync.Mutex
// 	var wg sync.WaitGroup
// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go func(i int) {
// 			// count++
// 			// lock.Lock()
// 			count = atomic.AddInt64(&count, 1)
// 			// count++
// 			// count++
// 			coll.Update(bson.M{"name": "concurrTest"},
// 				bson.M{"$set": bson.M{"count": count}})
// 			// lock.Unlock()
// 			wg.Done()

// 		}(i)
// 	}

// 	wg.Wait()
// 	fmt.Println(count)

// }
