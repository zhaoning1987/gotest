package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	mgoutil "github.com/qiniu/db/mgoutil.v3"
	mgo "gopkg.in/mgo.v2"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"
)

const (
	BASE64_PREFIX = "data:application/octet-stream;base64,"
)

type Config struct {
	LogPath      string         `json:"log_path"`
	ImageDirPath string         `json:"image_dir_path"`
	HTTPHost     string         `json:"qiniu_host_url"`
	Timeout      int            `json:"http_timeout_in_millisecond"`
	MaxTryTime   int            `json:"max_try_time"`
	ThreadNumber int            `json:"thread_number"`
	PoolSize     int            `json:"job_pool_size"`
	DBConfig     mgoutil.Config `json:"db_config"`
	DBCollection string         `json:"db_collection_name"`
}

type FaceJob struct {
	index        int
	ctx          *context.Context
	faceFeature  feature.FaceFeature
	imageContent []byte
	imageURI     string
	collection   *mgo.Collection
	wg           *sync.WaitGroup
	tryTime      int
}

type Face struct {
	Path string
	Spec []byte
}

func getFaceFeature(ctx *context.Context, face feature.FaceFeature, imageContent string, imageSource string) (fv proto.FeatureValue, err error) {
	boudingBox, err := face.FaceBoxes(*ctx, proto.ImageURI(imageContent))
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("error when calling func FaceBoxes for the image: %s", imageSource))
	}

	if len(boudingBox) == 0 {
		return nil, errors.Errorf("do not contain any face for the image: %s", imageSource)
	} else if len(boudingBox) > 1 {
		return nil, errors.Errorf("contain more than one face for the image: %s " + imageSource)
	} else {
		for _, box := range boudingBox {
			fv, err = face.Face(*ctx, proto.ImageURI(imageContent), box.Pts)
			if err != nil {
				return nil, errors.Wrapf(err, fmt.Sprintf("error when calling func Face for the image: %s", imageSource))
			}
		}
	}
	return
}

func getFaceFeatureByURI(ctx *context.Context, face feature.FaceFeature, uri string) (fv proto.FeatureValue, err error) {
	return getFaceFeature(ctx, face, uri, uri)
}

func getFaceFeatureByData(ctx *context.Context, face feature.FaceFeature, imageContent []byte, imageSource string) (fv proto.FeatureValue, err error) {
	imgBase64 := BASE64_PREFIX + base64.StdEncoding.EncodeToString(imageContent)
	return getFaceFeature(ctx, face, imgBase64, imageSource)
}

func uploadLog(path string, msg string) {
	err := ioutil.WriteFile(path, []byte(msg), 0644)
	if err != nil {
		panic(err)
	}
}

func systemLog(msg string) {
	fmt.Println(msg)
	f, err := os.OpenFile(systemLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	curTime := time.Now().Format("2006-01-02 15:04:05")
	content := fmt.Sprintf("%s %s\n", curTime, msg)
	buf := []byte(content)
	f.Write(buf)
}

func errorListLog(msg string) {
	f, err := os.OpenFile(errorListFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	content := fmt.Sprintf("%s\n", msg)
	buf := []byte(content)
	f.Write(buf)
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
