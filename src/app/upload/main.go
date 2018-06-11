package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	mgoutil "github.com/qiniu/db/mgoutil.v3"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"

	multithread "app/multithread"
)

const (
	BASE64_PREFIX = "data:application/octet-stream;base64,"
)

var (
	logPath      = "./log/"
	lastMinIndex = 0
	lastMaxIndex = 0
	wg           sync.WaitGroup
	dirPath      = "/Users/zhaoning/Desktop/testFace/"
	reqURL       = "http://ava-serving-gate.cs.cg.dora-internal.qiniu.io:5001"
	// reqURL = "http://100.100.58.85:9000"
	imgURL = "http://oayjpradp.bkt.clouddn.com/age_gender_test.png"
	// imgURL   = "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1528194279170&di=8d7e47958792fa6719e179b7de9f2cdf&imgtype=0&src=http%3A%2F%2Fpic.58pic.com%2F58pic%2F14%2F79%2F64%2F04I58PICefM_1024.jpg"
	timout   = time.Duration(0 * time.Second)
	dbConfig = mgoutil.Config{
		Host:           "localhost:27017",
		DB:             "testFace",
		Mode:           "strong",
		SyncTimeoutInS: 1,
	}
	jobPool chan multithread.Job
)

type FaceJob struct {
	index       int
	ctx         *context.Context
	faceFeature feature.FaceFeature
	uri         string
	path        string
	collection  *mgo.Collection
	wg          *sync.WaitGroup
}

type Face struct {
	Path string
	Spec []byte
}

func getFaceFeatureByURI(ctx *context.Context, face feature.FaceFeature, uri string) (fv proto.FeatureValue, err error) {
	// t1 := time.Now()
	boudingBox, err := face.FaceBoxes(*ctx, proto.ImageURI(uri))
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
			fv, err = face.Face(*ctx, proto.ImageURI(uri), box.Pts)
			// fmt.Println("time begin facefeature: ", t2, "Time used: ", time.Since(t2))
			if err != nil {
				// TODO
				return nil, errors.Wrapf(err, "error when calling func Face")
			}
		}
	}
	return
}

func getFaceFeatureByData(ctx *context.Context, face feature.FaceFeature, data []byte) (fv proto.FeatureValue, err error) {
	imgBase64 := BASE64_PREFIX + base64.StdEncoding.EncodeToString(data)
	return getFaceFeatureByURI(ctx, face, imgBase64)
}

func getFeatureAndSave(workerIndex int, param ...interface{}) {
	//when start to handle a face, we log it
	fj := param[0].(*FaceJob)
	log(workerIndex, strconv.Itoa(fj.index))
	defer fj.wg.Done()
	fv, err := getFaceFeatureByURI(fj.ctx, fj.faceFeature, fj.uri)
	if err != nil {
		fmt.Println(err)
		fj.wg.Add(1)
		retry := multithread.Job{}
		retry.Param = []interface{}{fj}
		retry.Fn = getFeatureAndSave
		jobPool <- retry
		fmt.Println("retry", fj.index)
		return
	}

	// wait := rand.Intn(5) + 1
	// time.Sleep(time.Duration(wait) * 100 * time.Millisecond)
	// fv := []byte("for test")
	face := &Face{fj.path, fv}
	err = fj.collection.Insert(face)
	if err != nil {
		fmt.Println(err)
	}
}

func log(index int, msg string) {
	err := ioutil.WriteFile(fmt.Sprintf(logPath+"%d.log", index), []byte(msg), 0644)
	if err != nil {
		panic(err)
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func prepare() {
	// 1.create log dir if not exist
	exist, err := pathExists(logPath)
	if err != nil {
		fmt.Printf("get directory error![%v]\n", err)
		panic(err)
	}
	if !exist {
		fmt.Printf("no directory: %v\n", logPath)
		// create folder
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			fmt.Printf("create directory failed: %v\n", logPath)
			panic(err)
		} else {
			fmt.Printf("create directory success: %v\n", logPath)
		}
	}

	// 2. read log file to get the last min&max stop point
	var indexes []int
	err = filepath.Walk(logPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		dat, err := ioutil.ReadFile(path)
		if err != nil {
			// TODO
			fmt.Println(err)
			return nil
		}
		index, err := strconv.Atoi(string(dat))
		if err != nil {
			// TODO
			fmt.Println(err)
			return nil
		}
		indexes = append(indexes, index)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	if indexes == nil {
		lastMinIndex = 0
		lastMaxIndex = 0
	} else {
		sort.Ints(indexes)
		lastMinIndex = indexes[0]
		lastMaxIndex = indexes[len(indexes)-1]
	}

}

func main() {
	faceFeature := feature.NewFaceFeature(reqURL, timout, 2048)
	ctx := context.Background()
	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	prepare()

	coll := session.DB(dbConfig.DB).C("all")
	poolSize := 10
	workerSize := 10
	jobPool = make(chan multithread.Job, poolSize)
	dispatcher := multithread.NewDispatcher(jobPool, workerSize)
	dispatcher.Start()
	t1 := time.Now()

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < poolSize; i++ {
		//handle the file
		path := fmt.Sprintf("/usr/local/images/image%d", i)

		if i < lastMinIndex {
			// this file is already handled in the last time
			continue
		}

		if i <= lastMaxIndex && lastMaxIndex > 0 {
			// not sure this file is handled or not, check it
			// if exist, skip it
			result := Face{}
			err = coll.Find(bson.M{"path": path}).One(&result)
			if err == nil {
				// if no error, then item found, so we skip this file
				fmt.Println("found item", i)
				continue
			} else {
				fmt.Println("not found item", i)
			}
		}

		wg.Add(1)
		job := multithread.Job{}
		job.Param = []interface{}{&FaceJob{i, &ctx, &faceFeature, imgURL, path, coll, &wg}}
		job.Fn = getFeatureAndSave
		jobPool <- job
	}

	wg.Wait()
	close(jobPool)
	elapsed := time.Since(t1)
	fmt.Println("Time elapsed: ", elapsed)

}
