package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	xlog "github.com/qiniu/xlog.v1"
	"qbox.us/cc/config"
	"qiniu.com/argus/feature_group_private/dbstorage"
)

var (
	jobPool               chan dbstorage.FaceJob
	wg                    sync.WaitGroup
	imageSourceFolderPath = "./source/"
	imageSourceFile       = "./urlSource"
)

func loadFromFile(ctx context.Context, faceGroup *dbstorage.FaceGroup, lastIndex int) {
	xl := xlog.FromContextSafe(ctx)

	file, err := os.Open(imageSourceFile)
	if err != nil {
		xl.Fatalf("error when opening image uri list file <%s>: %s\n", imageSourceFile, err)
	}
	defer file.Close()

	//currently only support csv file
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		xl.Fatalf("error when reading image uri list file <%s>: %s\n", imageSourceFile, err)
	}

	startTime := time.Now()
	for i, data := range lines {
		if i < lastIndex {
			//this line is already handled in the last time, so skip it
			dbstorage.IncrementCount()
			continue
		}

		var uri, tag, desc string
		uri = strings.TrimSpace(data[0])
		if len(data) > 1 {
			tag = strings.TrimSpace(data[1])
		}
		if len(data) > 2 {
			desc = strings.TrimSpace(data[2])
		}

		wg.Add(1)
		job := dbstorage.NewFaceJob(ctx, i, faceGroup, nil, uri, tag, desc, &wg)
		jobPool <- *job
	}

	wg.Wait()
	close(jobPool)
	elapsed := time.Since(startTime)
	xl.Info("Time elapsed: ", elapsed)
	xl.Info("finished!!!")
}

func loadFromFolder(ctx context.Context, faceGroup *dbstorage.FaceGroup, lastIndex int) {
	xl := xlog.FromContextSafe(ctx)

	startTime := time.Now()
	i := -1
	err := filepath.Walk(imageSourceFolderPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		//skip directory & windows thumb file & linux max hidden file
		if f.IsDir() || strings.ToLower(f.Name()) == "thumb.db" || dbstorage.Substring(f.Name(), 0, 1) == "." {
			return nil
		}

		imageContent, err := ioutil.ReadFile(path)
		if err != nil {
			dbstorage.LogError(fmt.Sprintf("%s : error when reading file : %s", path, err.Error()))
			return nil
		}

		//start to handle normal file
		i++
		if i < lastIndex {
			//this file is already handled in the last time, so skip it
			dbstorage.IncrementCount()
			return nil
		}

		tag, desc := dbstorage.GetTagAndDesc(f.Name())
		wg.Add(1)
		job := dbstorage.NewFaceJob(ctx, i, faceGroup, imageContent, path, tag, desc, &wg)
		jobPool <- *job
		return nil
	})
	if err != nil {
		xl.Fatalf("error when filepath.Walk() on [%s]: %v\n", imageSourceFolderPath, err)
	}

	wg.Wait()
	close(jobPool)
	elapsed := time.Since(startTime)
	xl.Info("Time elapsed: ", elapsed)
	xl.Info("finished!!!")
}

func main() {
	var (
		ctx  = xlog.NewContext(context.Background(), xlog.NewDummy())
		xl   = xlog.FromContextSafe(ctx)
		conf dbstorage.Config
	)

	//load config file
	config.Init("f", "dbstorage", "dbstorage.conf")
	if err := config.Load(&conf); err != nil {
		xl.Fatalln("Failed to load configure file")
	}

	xl.Infof("configuration: %+v\n", conf)

	if conf.ImageFolderPath != "" {
		imageSourceFolderPath = conf.ImageFolderPath
	}
	if conf.ImageListFile != "" {
		imageSourceFile = conf.ImageListFile
	}

	faceGroup := dbstorage.NewFaceGroup(conf.ServiceHost, time.Duration(conf.Timeout)*time.Millisecond)

	//setup before import to db
	lastIndex := dbstorage.Setup(ctx, conf)

	//initial thread pool
	poolSize := conf.PoolSize
	workerSize := conf.ThreadNumber
	jobPool = make(chan dbstorage.FaceJob, poolSize)
	dispatcher := dbstorage.NewDispatcher(jobPool, workerSize)
	dispatcher.Start()

	//create group if not exist
	err := faceGroup.CreateGroup(ctx, conf.GroupName)
	if err == nil {
		xl.Infof("create group <%s> successful", conf.GroupName)
	}

	if lastIndex == 0 {
		xl.Infof("start")
	} else {
		xl.Infof("continue from previous stop point: %d", lastIndex)
	}

	if conf.LoadImageFromFolder {
		//load image from folder
		loadFromFolder(ctx, faceGroup, lastIndex)
	} else {
		//load image from uri list file
		loadFromFile(ctx, faceGroup, lastIndex)
	}

	//clear up
	dbstorage.Teardown()
}
