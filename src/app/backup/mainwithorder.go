package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	mgoutil "github.com/qiniu/db/mgoutil.v3"
	mgo "gopkg.in/mgo.v2"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"

	multithread "app/multithreadwithorder"
)

const (
	BASE64_PREFIX = "data:application/octet-stream;base64,"
)

var (
	dirPath = "/Users/zhaoning/Desktop/testFace/"
	reqURL  = "http://ava-serving-gate.cs.cg.dora-internal.qiniu.io:5001"
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
	jobPool []chan multithread.Job
)

type FaceJob struct {
	index       int
	ctx         *context.Context
	faceFeature feature.FaceFeature
	uri         string
	path        string
	collection  *mgo.Collection
	db          *mgo.Database
	wg          *sync.WaitGroup
}

type Face struct {
	Path string
	Spec []byte
}

func getFaceFeatureByURI(ctx *context.Context, face feature.FaceFeature, uri string) (fv proto.FeatureValue, err error) {
	boudingBox, err := face.FaceBoxes(*ctx, proto.ImageURI(uri))
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
			fv, err = face.Face(*ctx, proto.ImageURI(uri), box.Pts)
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

func getFeatureAndSave(param ...interface{}) {
	fj := param[0].(*FaceJob)
	defer fj.wg.Done()
	fv, err := getFaceFeatureByURI(fj.ctx, fj.faceFeature, fj.uri)
	if err != nil {
		fmt.Println(err)
		// fj.wg.Add(1)
		// retry := multithread.Job{}
		// retry.Param = []interface{}{fj}
		// retry.Fn = getFeatureAndSave
		// jobPool[fj.index] <- retry
		return
	}
	face := &Face{fj.path, fv}
	err = fj.db.C(fmt.Sprintf("coll%d", fj.index)).Insert(face)
	if err != nil {
		fmt.Println(err)
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
	coll := session.DB(dbConfig.DB).C("all")
	db := session.DB(dbConfig.DB)

	poolSize := 1000
	workerSize := 10

	dispatcher := multithread.NewDispatcher(workerSize)
	jobPool = dispatcher.JobPool
	dispatcher.Start()
	t1 := time.Now()
	var wg sync.WaitGroup
	wg.Add(workerSize)
	for thread := 0; thread < workerSize; thread++ {
		go func(index int) {
			defer wg.Done()
			for i := 0; i < poolSize; i++ {
				if i%workerSize != index {
					continue
				}
				path := fmt.Sprintf("/usr/local/images/image%d", i)
				wg.Add(1)
				job := multithread.Job{}
				job.Param = []interface{}{&FaceJob{index, &ctx, &faceFeature, imgURL, path, coll, db, &wg}}
				job.Fn = getFeatureAndSave
				jobPool[index] <- job
			}
			close(jobPool[index])
		}(thread)
	}

	wg.Wait()
	elapsed := time.Since(t1)
	fmt.Println("Time elapsed: ", elapsed)
}
