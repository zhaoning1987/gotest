package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	multithread "app/multithread"

	mgoutil "github.com/qiniu/db/mgoutil.v3"
	"gopkg.in/mgo.v2/bson"
	"qbox.us/cc/config"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"
)

var (
	conf          Config
	imgURL        = "http://oayjpradp.bkt.clouddn.com/age_gender_test.png"
	jobPool       chan multithread.Job
	lastMinIndex  int
	lastMaxIndex  int
	wg            sync.WaitGroup
	logPath       string
	systemLogFile string
	errorListFile string
)

func getFeatureAndSave(workerIndex int, param ...interface{}) {
	fj := param[0].(*FaceJob)

	//when start to handle a face, we log its index
	uploadLog(fmt.Sprintf(logPath+"%d.log", workerIndex), strconv.Itoa(fj.index))
	defer fj.wg.Done()

	var (
		fv  proto.FeatureValue
		err error
	)
	if fj.imageContent != nil {
		fv, err = getFaceFeatureByData(fj.ctx, fj.faceFeature, fj.imageContent, fj.imageURI)
	} else {
		fv, err = getFaceFeatureByURI(fj.ctx, fj.faceFeature, fj.imageURI)
	}

	if err != nil {
		//when error, retry
		retry(err.Error(), fj)
		return
	}

	// wait := rand.Intn(5) + 1
	// time.Sleep(time.Duration(wait) * 100 * time.Millisecond)
	// fv := []byte("for test")
	face := &Face{fj.imageURI, fv}
	err = fj.collection.Insert(face)
	if err != nil {
		//when error, retry
		retry(err.Error(), fj)
	}
}

func retry(errMsg string, fj *FaceJob) {
	//log tye error to system log
	systemLog(errMsg)
	if fj.tryTime == conf.MaxTryTime {
		//reach the max retry time, store this file in errorlist log
		errorListLog(errMsg)
	} else {
		//retry
		fj.tryTime++
		fj.wg.Add(1)
		retry := multithread.Job{}
		retry.Param = []interface{}{fj}
		retry.Fn = getFeatureAndSave
		jobPool <- retry
	}
	return
}

func prepare() {
	// 1.create log dir if not exist
	exist, err := pathExists(logPath)
	if err != nil {
		fmt.Printf("get directory [%s] error: %v\n", logPath, err)
		panic(err)
	}
	if !exist {
		fmt.Printf(fmt.Sprintf("no directory [%s]\n", logPath))
		// create folder
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			fmt.Printf(fmt.Sprintf("create directory [%s] failed: %v\n", logPath, err))
			panic(err)
		} else {
			fmt.Printf(fmt.Sprintf("create directory [%s] success\n", logPath))
		}
	}

	// 2. read log file to get the last min&max stop point
	var indexes []int
	err = filepath.Walk(logPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() || filepath.Ext(path) != ".log" {
			return nil
		}

		dat, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		index, err := strconv.Atoi(string(dat))
		if err != nil {
			return err
		}
		indexes = append(indexes, index)
		return nil
	})
	if err != nil {
		systemLog(fmt.Sprintf("error when filepath.Walk() on [%s]: %v", logPath, err))
		panic(err)
	}
	if indexes == nil {
		lastMinIndex = 0
		lastMaxIndex = 0
	} else {
		sort.Ints(indexes)
		lastMinIndex = indexes[0]
		lastMaxIndex = indexes[len(indexes)-1]
	}
	systemLogFile = logPath + "output"
	errorListFile = logPath + "errorlist"
}

func main() {
	config.Init("f", "/upload", "upload.conf")
	if err := config.Load(&conf); err != nil {
		log.Fatal("Failed to load configure file")
	}
	fmt.Printf("configuration: %+v\n", conf)
	logPath = conf.LogPath
	faceFeature := feature.NewFaceFeature(conf.HTTPHost, time.Duration(conf.Timeout)*time.Millisecond, 2048)
	ctx := context.Background()
	session, err := mgoutil.Dail(conf.DBConfig.Host, conf.DBConfig.Mode, conf.DBConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	prepare()

	coll := session.DB(conf.DBConfig.DB).C(conf.DBCollection)
	poolSize := conf.PoolSize
	workerSize := conf.ThreadNumber
	jobPool = make(chan multithread.Job, poolSize)
	dispatcher := multithread.NewDispatcher(jobPool, workerSize)
	dispatcher.Start()
	t1 := time.Now()

	if lastMinIndex == 0 && lastMaxIndex == 0 {
		systemLog("start")
	} else {
		systemLog(fmt.Sprintf("continue, the previous stop points are from %d to %d", lastMinIndex, lastMaxIndex))
	}

	//====== image from url
	// for i := 0; i < poolSize; i++ {
	// 	//handle the file
	// 	path := fmt.Sprintf("/usr/local/images/image%d", i)

	// 	if i < lastMinIndex {
	// 		// this file is already handled in the last time, so skip it
	// 		continue
	// 	}

	// 	if i <= lastMaxIndex && lastMaxIndex > 0 {
	// 		// not sure this file is handled or not, check it
	// 		// if exist, skip it
	// 		result := Face{}
	// 		err = coll.Find(bson.M{"path": path}).One(&result)
	// 		if err == nil {
	// 			// no error => item found => file handled => we skip this file
	// 			systemLog(fmt.Sprintf("file %s is already handled during the previous run", path))
	// 			continue
	// 		} else {
	// 			systemLog(fmt.Sprintf("file %s is not handled during the previous run, do it this time", path))
	// 		}
	// 	}

	// 	wg.Add(1)
	// 	job := multithread.Job{}
	// 	job.Param = []interface{}{&FaceJob{i, &ctx, &faceFeature, nil, imgURL, coll, &wg, 1}}
	// 	job.Fn = getFeatureAndSave
	// 	jobPool <- job
	// }

	//====== image from folder
	i := 0
	err = filepath.Walk(conf.ImageDirPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			// TODO need to add only handle image
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			// TODO problem file
			fmt.Println(err)
			return nil
		}

		//start to handle normal file
		i++
		if i < lastMinIndex {
			// this file is already handled in the last time, so skip it
			return nil
		}

		if i <= lastMaxIndex && lastMaxIndex > 0 {
			// not sure this file is handled or not, check it. if exist, skip it
			result := Face{}
			if err = coll.Find(bson.M{"path": path}).One(&result); err == nil {
				// no error => item found => file handled => we skip this file
				systemLog(fmt.Sprintf("file %s is already handled during the previous run", path))
				return nil
			}
			systemLog(fmt.Sprintf("file %s is not handled during the previous run, do it this time", path))
		}

		wg.Add(1)
		job := multithread.Job{}
		job.Param = []interface{}{&FaceJob{i, &ctx, &faceFeature, data, path, coll, &wg, 1}}
		job.Fn = getFeatureAndSave
		jobPool <- job

		return nil
	})
	if err != nil {
		systemLog(fmt.Sprintf("error when filepath.Walk() on [%s]: %v", conf.ImageDirPath, err))
		panic(err)
	}

	wg.Wait()
	close(jobPool)
	elapsed := time.Since(t1)
	fmt.Println("Time elapsed: ", elapsed)
	systemLog("finished!!!")
}
