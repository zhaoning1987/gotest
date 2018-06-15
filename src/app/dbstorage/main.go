package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	mgoutil "github.com/qiniu/db/mgoutil.v3"
	threadpool "qiniu.com/ava/argus/feature_group_private/dbstorage/threadpool"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"qbox.us/cc/config"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"
)

var (
	conf          Config
	jobPool       chan threadpool.Job
	lastMinIndex  int
	lastMaxIndex  int
	wg            sync.WaitGroup
	logPath       string
	systemLogFile string
	errorListFile string
)

func getFeatureAndSave(workerIndex int, param ...interface{}) {
	fj := param[0].(*FaceJob)

	//when start to handle a face, we log its index into current thread's log file
	uploadLog(fmt.Sprintf(logPath+"%d.log", workerIndex), strconv.Itoa(fj.index))
	defer fj.wg.Done()

	var (
		fv  proto.FeatureValue
		err error
		md5 string
	)

	//if image content not passed, get it by url
	if fj.imageContent == nil {
		t1 := time.Now()
		fj.imageContent, err = downloadFile(fj.imageURI)
		fmt.Println("download time for", workerIndex, "is", time.Since(t1))
		if err != nil {
			retry(err.Error(), workerIndex, fj)
			return
		}
	}

	//check if this image exist in db
	md5 = getMd5(fj.imageContent)
	if valueExistInDB(fj.collection, "md5", md5) {
		errorListLog(fmt.Sprintf("%s is duplicated, therefore not imported in database", fj.imageURI))
		return
	}

	//get feature
	fv, err = getFaceFeatureByURI(fj.ctx, fj.faceFeature, fj.imageURI)
	if err != nil {
		//when error, retry
		retry(err.Error(), workerIndex, fj)
		return
	}

	face := &Face{fj.imageURI, fv, md5}
	err = fj.collection.Insert(face)
	if err != nil {
		//if the err is caused by duplicated key, log the error
		//else retry
		if strings.Contains(err.Error(), "duplicate") {
			errorListLog(fmt.Sprintf("%s is duplicated, therefore not imported in database", fj.imageURI))
		} else {
			retry(err.Error(), workerIndex, fj)
		}
	}
	return
}

func retry(errMsg string, workerIndex int, fj *FaceJob) {
	//log error to system log
	systemLog(errMsg)
	if fj.tryTime == conf.MaxTryTime {
		//reach the max retry time, store this file in errorlist log
		errorListLog(errMsg)
	} else {
		//retry
		fj.tryTime++
		fj.wg.Add(1)
		getFeatureAndSave(workerIndex, fj)
		// fj.tryTime++
		// fj.wg.Add(1)
		// retry := threadpool.Job{}
		// retry.Param = []interface{}{fj}
		// retry.Fn = getFeatureAndSave
		// jobPool <- retry
	}
	return
}

func prepare() {
	//1.create log dir if not exist
	exist, err := pathExists(logPath)
	if err != nil {
		fmt.Printf("get directory [%s] error: %v\n", logPath, err)
		panic(err)
	}
	if !exist {
		fmt.Printf(fmt.Sprintf("no directory [%s]\n", logPath))
		//create folder
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			fmt.Printf(fmt.Sprintf("create directory [%s] failed: %v\n", logPath, err))
			panic(err)
		} else {
			fmt.Printf(fmt.Sprintf("create directory [%s] success\n", logPath))
		}
	}

	//2.read log file to get the last min&max stop point
	var indexes []int
	err = filepath.Walk(logPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() || filepath.Ext(f.Name()) != ".log" {
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

func loadFromURIList(ctx *context.Context, coll *mgo.Collection, faceFeature feature.FaceFeature) {
	startTime := time.Now()

	file, err := os.Open(conf.ImageListFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		imageURL := scanner.Text()

		if i < lastMinIndex {
			//this file is already handled in the last time, so skip it
			continue
		}

		if i <= lastMaxIndex && lastMaxIndex > 0 {
			//not sure this file is handled or not, check it. if exist, skip it
			result := Face{}
			if err = coll.Find(bson.M{"uri": imageURL}).One(&result); err == nil {
				//no error => item found => file handled => we skip this file
				systemLog(fmt.Sprintf("file %s is already handled during the previous run", imageURL))
				continue
			}
			systemLog(fmt.Sprintf("file %s is not handled during the previous run, do it this time", imageURL))
		}

		wg.Add(1)
		job := threadpool.Job{}
		job.Param = []interface{}{&FaceJob{i, ctx, faceFeature, nil, imageURL, coll, &wg, 1}}
		job.Fn = getFeatureAndSave
		jobPool <- job
		i++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	wg.Wait()
	close(jobPool)
	elapsed := time.Since(startTime)
	fmt.Println("Time elapsed: ", elapsed)
	systemLog("finished!!!")
}

func loadFromFolder(ctx *context.Context, coll *mgo.Collection, faceFeature feature.FaceFeature) {
	startTime := time.Now()
	i := -1
	err := filepath.Walk(conf.ImageDirPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		//skip directory & windows thumb file & linux max hidden file
		if f.IsDir() || strings.ToLower(f.Name()) == "thumb.db" || substring(f.Name(), 0, 1) == "." {
			return nil
		}

		imageContent, err := ioutil.ReadFile(path)
		if err != nil {
			errorListLog(fmt.Sprintf("error when reading file %s: %s", path, err.Error()))
			return nil
		}

		//start to handle normal file
		i++
		if i < lastMinIndex {
			//this file is already handled in the last time, so skip it
			return nil
		}

		if i <= lastMaxIndex && lastMaxIndex > 0 {
			//not sure this file is handled or not, check it. if exist, skip it
			result := Face{}
			if err = coll.Find(bson.M{"uri": path}).One(&result); err == nil {
				//no error => item found => file handled => we skip this file
				systemLog(fmt.Sprintf("file %s is already handled during the previous run", path))
				return nil
			}
			systemLog(fmt.Sprintf("file %s is not handled during the previous run, do it this time", path))
		}

		wg.Add(1)
		job := threadpool.Job{}
		job.Param = []interface{}{&FaceJob{i, ctx, faceFeature, imageContent, path, coll, &wg, 1}}
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
	elapsed := time.Since(startTime)
	fmt.Println("Time elapsed: ", elapsed)
	systemLog("finished!!!")
}

func main() {
	//load config file
	config.Init("f", "/dbstorage", "dbstorage.conf")
	if err := config.Load(&conf); err != nil {
		log.Fatal("Failed to load configure file")
	}
	fmt.Printf("configuration: %+v\n", conf)
	logPath = conf.LogPath
	faceFeature := feature.NewFaceFeature(conf.HTTPHost, time.Duration(conf.Timeout)*time.Millisecond, FEATURE_SIZE)
	ctx := context.Background()

	//connnect db
	session, err := mgoutil.Dail(conf.DBConfig.Host, conf.DBConfig.Mode, conf.DBConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//ensure index
	coll := session.DB(conf.DBConfig.DB).C(conf.DBCollection)
	index := mgo.Index{
		Key:    []string{"md5"},
		Unique: true,
	}
	if err := coll.EnsureIndex(index); err != nil {
		panic(err)
	}

	//do preparation before import db
	prepare()

	//initial thread pool
	poolSize := conf.PoolSize
	workerSize := conf.ThreadNumber
	jobPool = make(chan threadpool.Job, poolSize)
	dispatcher := threadpool.NewDispatcher(jobPool, workerSize)
	dispatcher.Start()

	if lastMinIndex == 0 && lastMaxIndex == 0 {
		systemLog("start")
	} else {
		systemLog(fmt.Sprintf("continue, the previous stop points are from %d to %d", lastMinIndex, lastMaxIndex))
	}

	//load image from folder
	//loadFromFolder(&ctx, coll, &faceFeature)

	//load image from urlList
	loadFromURIList(&ctx, coll, &faceFeature)
}
