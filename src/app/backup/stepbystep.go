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
	feature_group "qiniu.com/ava/argus/feature_group_private"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"

	multithread "app/multithread"
)

const (
	BASE64_PREFIX = "data:application/octet-stream;base64,"
)

var (
	wg      sync.WaitGroup
	dirPath = "/Users/zhaoning/Desktop/testFace/"
	// reqURL  = "http://ava-serving-gate.cs.cg.dora-internal.qiniu.io:5001"
	reqURL = "http://100.100.58.85:9000"
	imgURL = "http://oayjpradp.bkt.clouddn.com/age_gender_test.png"
	// imgURL   = "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1528194279170&di=8d7e47958792fa6719e179b7de9f2cdf&imgtype=0&src=http%3A%2F%2Fpic.58pic.com%2F58pic%2F14%2F79%2F64%2F04I58PICefM_1024.jpg"
	timout   = time.Duration(0)
	dbConfig = mgoutil.Config{
		Host:           "localhost:27017",
		DB:             "testFace",
		Mode:           "strong",
		SyncTimeoutInS: 1,
	}
	dbColl  *mgo.Collection
	jobPool chan multithread.Job
)

type FaceBoxJob struct {
	ctx         *context.Context
	faceFeature feature.FaceFeature
	uri         string
	wg          *sync.WaitGroup
}

type FaceFeatureJob struct {
	ctx         *context.Context
	faceFeature feature.FaceFeature
	uri         string
	pts         [][2]int
	wg          *sync.WaitGroup
}

type DBJob struct {
	fv         proto.FeatureValue
	path       string
	collection *mgo.Collection
	wg         *sync.WaitGroup
}

type Face struct {
	Path string
	Spec []byte
}

func getFaceBox(ctx *context.Context, face feature.FaceFeature, uri string) (box *feature_group.FaceBoudingBox, err error) {
	boudingBox, err := face.FaceBoxes(*ctx, proto.ImageURI(uri))
	if err != nil {
		// TODO
		return nil, errors.Wrapf(err, "error when calling func FaceBoxes")
		return
	}

	if len(boudingBox) == 0 {
		// TODO
		return nil, errors.Errorf("do not contain any face for the image: %s", uri)
	} else if len(boudingBox) > 1 {
		// TODO
		return nil, errors.Errorf("contain more than one face for the image: %s " + uri)
	} else {
		return &boudingBox[0], nil
	}
}

func getFaceFeatureByURI(ctx *context.Context, face feature.FaceFeature, uri string, pts [][2]int) (fv proto.FeatureValue, err error) {
	fv, err = face.Face(*ctx, proto.ImageURI(uri), pts)
	if err != nil {
		// TODO
		return nil, errors.Wrapf(err, "error when calling func Face")
	}
	return
}

func getFaceFeatureByData(ctx *context.Context, face feature.FaceFeature, data []byte, pts [][2]int) (fv proto.FeatureValue, err error) {
	imgBase64 := BASE64_PREFIX + base64.StdEncoding.EncodeToString(data)
	return getFaceFeatureByURI(ctx, face, imgBase64, pts)
}

func insertDB(coll *mgo.Collection, path string, fv proto.FeatureValue) error {
	face := &Face{path, fv}
	err := coll.Insert(face)
	return err
}

func faceBoxJob(param ...interface{}) {
	job := param[0].(*FaceBoxJob)
	defer job.wg.Done()
	box, err := getFaceBox(job.ctx, job.faceFeature, job.uri)
	if err != nil {
		fmt.Println(err)
		//retry
		job.wg.Add(1)
		nextJob := multithread.Job{}
		nextJob.Param = []interface{}{job}
		nextJob.Fn = faceBoxJob
		jobPool <- nextJob
		return
	}
	job.wg.Add(1)
	nextJob := multithread.Job{}
	nextJob.Param = []interface{}{&FaceFeatureJob{job.ctx, job.faceFeature, imgURL, box.Pts, job.wg}}
	nextJob.Fn = faceFeatureJob
	jobPool <- nextJob
}

func faceFeatureJob(param ...interface{}) {
	job := param[0].(*FaceFeatureJob)
	defer job.wg.Done()
	fv, err := getFaceFeatureByURI(job.ctx, job.faceFeature, job.uri, job.pts)
	if err != nil {
		fmt.Println(err)
		//retry
		job.wg.Add(1)
		nextJob := multithread.Job{}
		nextJob.Param = []interface{}{job}
		nextJob.Fn = faceFeatureJob
		jobPool <- nextJob
		return
	}
	job.wg.Add(1)
	nextJob := multithread.Job{}
	nextJob.Param = []interface{}{&DBJob{fv, "path for test", dbColl, job.wg}}
	nextJob.Fn = dbJob
	jobPool <- nextJob
}

func dbJob(param ...interface{}) {
	job := param[0].(*DBJob)
	defer job.wg.Done()
	err := insertDB(job.collection, job.path, job.fv)
	if err != nil {
		fmt.Println(err)
		//retry
		job.wg.Add(1)
		nextJob := multithread.Job{}
		nextJob.Param = []interface{}{job}
		nextJob.Fn = dbJob
		jobPool <- nextJob
		return
	}
}

func main() {
	// runtime.GOMAXPROCS(2)
	// init param
	faceFeature := feature.NewFaceFeature(reqURL, timout, 2048)
	ctx := context.Background()
	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	dbColl = session.DB(dbConfig.DB).C("all")

	poolSize := 100
	workerSize := 10
	jobPool = make(chan multithread.Job, poolSize)
	dispatcher := multithread.NewDispatcher(jobPool, workerSize)
	dispatcher.Start()
	t1 := time.Now()

	for i := 0; i < poolSize; i++ {
		// path := fmt.Sprintf("/usr/local/images/image%d", i)
		wg.Add(1)
		job := multithread.Job{}
		job.Param = []interface{}{&FaceBoxJob{&ctx, &faceFeature, imgURL, &wg}}
		job.Fn = faceBoxJob
		jobPool <- job
	}

	wg.Wait()
	close(jobPool)

	elapsed := time.Since(t1)
	fmt.Println("Time elapsed: ", elapsed)

}