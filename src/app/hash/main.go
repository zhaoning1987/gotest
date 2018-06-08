package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
	mgoutil "github.com/qiniu/db/mgoutil.v3"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"
)

const (
	BASE64_PREFIX = "data:application/octet-stream;base64,"
)

var (
	dirPath = "/Users/zhaoning/Desktop/testFace/"
	reqURL  = "http://ava-serving-gate.cs.cg.dora-internal.qiniu.io:5001"
	imgURL  = "http://oayjpradp.bkt.clouddn.com/age_gender_test.png"
	// imgURL   = "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1528194279170&di=8d7e47958792fa6719e179b7de9f2cdf&imgtype=0&src=http%3A%2F%2Fpic.58pic.com%2F58pic%2F14%2F79%2F64%2F04I58PICefM_1024.jpg"
	timout   = time.Duration(0 * time.Second)
	dbConfig = mgoutil.Config{
		Host:           "localhost:27017",
		DB:             "testFace",
		Mode:           "strong",
		SyncTimeoutInS: 1,
	}
)

type Face struct {
	Path string
	Spec []byte
}

func getFaceFeatureByURI(ctx context.Context, face feature.FaceFeature, uri string) (fv proto.FeatureValue, err error) {
	// t1 := time.Now()
	runtime.Gosched()
	boudingBox, err := face.FaceBoxes(ctx, proto.ImageURI(uri))
	// fmt.Println("time begin facebox: ", t1, "Time used: ", time.Since(t1))
	if err != nil {
		// TODO
		return nil, errors.Wrapf(err, "error when calling func FaceBoxes")
	}

	if len(boudingBox) == 0 {
		// TODO
		return nil, errors.Errorf("do not contain any face for the image: %s", uri)
	} else if len(boudingBox) > 1 {
		// TODO
		return nil, errors.Errorf("contain more than one face for the image: %s " + uri)
	} else {
		for _, box := range boudingBox {
			// t2 := time.Now()
			fv, err = face.Face(ctx, proto.ImageURI(uri), box.Pts)
			// fmt.Println("time begin facefeature: ", t2, "Time used: ", time.Since(t2))
			if err != nil {
				// TODO
				return nil, errors.Wrapf(err, "error when calling func Face")
			}
		}
	}
	return
}

func getFaceFeatureByData(ctx context.Context, face feature.FaceFeature, data []byte) (fv proto.FeatureValue, err error) {
	imgBase64 := BASE64_PREFIX + base64.StdEncoding.EncodeToString(data)
	return getFaceFeatureByURI(ctx, face, imgBase64)
}

func main() {
	// runtime.GOMAXPROCS(2)
	//init param
	faceFeature := feature.NewFaceFeature(reqURL, timout, 2048)
	ctx := context.Background()
	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	coll := session.DB(dbConfig.DB).C("all")

	// //get data
	// count := 1000
	// // block := 1
	// t1 := time.Now() // get current time
	// // var faceList []interface{}
	// ch := make(chan int, 100)
	// for i := 0; i < count; i++ {
	// 	path := fmt.Sprintf("/usr/local/images/image%d", i)
	// 	ch <- i
	// 	go func(ch chan int, index int) {
	// 		fv, err := getFaceFeatureByURI(ctx, faceFeature, imgURL)
	// 		if err != nil {
	// 			fmt.Println("index", index, err)
	// 			<-ch
	// 			return
	// 		}
	// 		face := &Face{path, fv}
	// 		err = coll.Insert(face)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		<-ch
	// 	}(ch, i)
	// 	// faceList = append(faceList, face)
	// 	// if i%block == block-1 {
	// 	// 	err = coll.Insert(faceList...)
	// 	// 	if err != nil {
	// 	// 		log.Fatal(err)
	// 	// 	}
	// 	// 	faceList = faceList[:0]
	// 	// }
	// }
	// // for i := 0; i < count; i++ {
	// // 	<-ch
	// // }
	// elapsed := time.Since(t1)
	// fmt.Println("Time elapsed: ", elapsed)
	//=========== read folder image
	t2 := time.Now() // get current time
	err = filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		fmt.Println(path)
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			// TODO
			fmt.Println(err)
			return nil
		}
		fv, err := getFaceFeatureByData(ctx, faceFeature, dat)
		if err != nil {
			// TODO
			fmt.Println(err)
			return nil
		}

		face := &Face{path, fv}
		err = coll.Insert(face)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	elapsed2 := time.Since(t2)
	fmt.Println("Time elapsed: ", elapsed2)
}
