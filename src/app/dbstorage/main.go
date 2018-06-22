package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	xlog "github.com/qiniu/xlog.v1"
	"qbox.us/cc/config"
	"qiniu.com/ava/argus/feature_group_private/manager"
	"qiniu.com/ava/argus/feature_group_private/proto"
)

var (
	conf           Config
	jobPool        chan FaceJob
	lastMinIndex   int
	wg             sync.WaitGroup
	outputLogPath  string
	systemLogFile  string
	errorFile      string
	processFile    string
	md5Dict        map[string]struct{}
	xl             *xlog.Logger
	mutex          = &sync.Mutex{}
	processLogPath = "./processlog/"
)

func (fj *FaceJob) execute(workerIndex int) {

	//when start to handle a face, we save its index in file
	writeToFile(xl, fmt.Sprintf(processLogPath+"%d.log", workerIndex), strconv.Itoa(fj.index))
	defer fj.wg.Done()

	var err error

	//if image content not passed, get it by url
	if fj.imageContent == nil {
		for i := 0; i < conf.MaxTryDownloadTime; i++ {
			fj.imageContent, err = downloadFile(fj.imageURI)
			if err == nil {
				break
			}
		}
		if err != nil {
			xl.Infof("%s : %s\n", fj.imageURI, err.Error())
			appendToFile(xl, errorFile, fmt.Sprintf("%s : %s", fj.imageURI, err.Error()))
			return
		}
	}

	//check if this image exist
	existed := false
	md5 := getMd5(fj.imageContent)
	mutex.Lock()
	if _, ok := md5Dict[md5]; ok {
		existed = true
	} else {
		md5Dict[md5] = struct{}{}
	}
	mutex.Unlock()

	if existed {
		//image exist, skip it
		return
	}

	//call group_add service to store the image
	imgBase64 := BASE64_PREFIX + base64.StdEncoding.EncodeToString(fj.imageContent)

	for i := 0; i < conf.MaxTryServiceTime; i++ {
		_, err = fj.faceGroup.Add(*fj.ctx, conf.GroupName, proto.FeatureID(fj.imageURI), proto.ImageURI(imgBase64))
		if err == nil {
			break
		}
	}

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, manager.ErrGroupNotExist.Err) {
			xl.Fatalf("group <%s> not exist\n", conf.GroupName)
		} else if !strings.Contains(errMsg, manager.ErrFeatureExist.Err) {
			appendToFile(xl, errorFile, fmt.Sprintf("%s : %s", fj.imageURI, err.Error()))
		}
		xl.Infof("%s : %s\n", fj.imageURI, err.Error())
	}

	//no matter call service success or not, always write to process file
	appendToFile(xl, processFile, md5)
	return
}

func prepare() {
	//1.create log dir if not exist
	createPath(xl, outputLogPath)
	createPath(xl, processLogPath)

	//2.read process log file to get the last min stop point
	var indexes []int
	err := filepath.Walk(processLogPath, func(path string, f os.FileInfo, err error) error {
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
		xl.Fatalf("error when filepath.Walk() on [%s]: %v\n", processLogPath, err)
	}
	if indexes == nil {
		lastMinIndex = 0
	} else {
		sort.Ints(indexes)
		lastMinIndex = indexes[0]
	}

	systemLogFile = outputLogPath + "output"
	errorFile = outputLogPath + "errorlist"
	processFile = processLogPath + "process"

	//3. get the uploaded files' md5
	md5Dict = map[string]struct{}{}
	if file, err := os.Open(processFile); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			md5Dict[scanner.Text()] = struct{}{}
		}
	}
}

func loadFromURIList(ctx *context.Context, faceGroup *faceGroup) {
	file, err := os.Open(conf.ImageListFile)
	if err != nil {
		xl.Fatalf("error when reading image uri list file <%s>: %s\n", conf.ImageListFile, err)
	}
	defer file.Close()

	startTime := time.Now()
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		imageURL := strings.TrimSpace(scanner.Text())

		if imageURL == "" {
			//skip empty line
			continue
		}
		if i < lastMinIndex {
			//this file is already handled in the last time, so skip it
			continue
		}

		wg.Add(1)
		job := FaceJob{i, ctx, faceGroup, nil, imageURL, &wg}
		jobPool <- job
		i++
	}

	wg.Wait()
	close(jobPool)
	elapsed := time.Since(startTime)
	xl.Info("Time elapsed: ", elapsed)
	xl.Info("finished!!!")
}

func loadFromFolder(ctx *context.Context, faceGroup *faceGroup) {
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
			appendToFile(xl, errorFile, fmt.Sprintf("error when reading file %s: %s", path, err.Error()))
			return nil
		}

		//start to handle normal file
		i++
		if i < lastMinIndex {
			//this file is already handled in the last time, so skip it
			return nil
		}

		wg.Add(1)
		job := FaceJob{i, ctx, faceGroup, imageContent, path, &wg}
		jobPool <- job
		return nil
	})
	if err != nil {
		xl.Fatalf("error when filepath.Walk() on [%s]: %v\n", conf.ImageDirPath, err)
	}

	wg.Wait()
	close(jobPool)
	elapsed := time.Since(startTime)
	xl.Info("Time elapsed: ", elapsed)
	xl.Info("finished!!!")
}

func main() {
	ctx := context.Background()
	xl = xlog.FromContextSafe(ctx)

	//load config file
	config.Init("f", "dbstorage", "dbstorage.conf")
	if err := config.Load(&conf); err != nil {
		xl.Fatalln("Failed to load configure file")
	}

	xl.Infof("configuration: %+v\n", conf)
	outputLogPath = conf.LogPath
	faceGroup := NewFaceGroup(conf.HTTPHost, time.Duration(conf.Timeout)*time.Millisecond)

	//do preparation before import db
	prepare()

	//initial thread pool
	poolSize := conf.PoolSize
	workerSize := conf.ThreadNumber
	jobPool = make(chan FaceJob, poolSize)
	dispatcher := NewDispatcher(jobPool, workerSize)
	dispatcher.Start()

	//create group if not exist
	err := faceGroup.CreateGroup(ctx, conf.GroupName)
	if err == nil {
		xl.Infof("create group <%s> successful", conf.GroupName)
	}

	if lastMinIndex == 0 {
		xl.Infof("start")
	} else {
		xl.Infof("continue from previous stop point: %d", lastMinIndex)
	}

	if conf.UseImageDirPath {
		//load image from folder
		loadFromFolder(&ctx, faceGroup)
	} else {
		//load image from uri list
		loadFromURIList(&ctx, faceGroup)
	}

}
