package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	mgoutil "github.com/qiniu/db/mgoutil.v3"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"
)

var (
	reqURL   = "http://ava-serving-gate.cs.cg.dora-internal.qiniu.io:5001"
	imgURL   = "http://oayjpradp.bkt.clouddn.com/age_gender_test.png"
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

func getFaceFeature(ctx context.Context, face feature.FaceFeature, imageURI string) (fv proto.FeatureValue, err error) {

	boudingBox, err := face.FaceBoxes(ctx, proto.ImageURI(imgURL))
	if err != nil {
		return nil, errors.Wrap(err, "error when calling func FaceBoxes")
	}

	if len(boudingBox) == 0 {
		return nil, errors.Wrap(err, "do not contain any face for the image: "+imageURI)
	} else if len(boudingBox) > 1 {
		return nil, errors.Wrap(err, "contain more than one face for the image: "+imageURI)
	} else {
		for _, box := range boudingBox {
			fv, err = face.Face(ctx, proto.ImageURI(imgURL), box.Pts)
			if err != nil {
				return nil, errors.Wrap(err, "error when calling func Face")
			}
		}
	}
	return
}

func main() {
	faceFeature := feature.NewFaceFeature(reqURL, timout, 2048)
	ctx := context.Background()
	session, err := mgoutil.Dail(dbConfig.Host, dbConfig.Mode, dbConfig.SyncTimeoutInS)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count := 100
	block := 100

	t1 := time.Now() // get current time
	coll := session.DB(dbConfig.DB).C("all")
	var faceList []interface{}
	for i := 0; i < count; i++ {
		path := fmt.Sprintf("usr/local/images/image%d", i)
		fv, err := getFaceFeature(ctx, faceFeature, imgURL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		face := &Face{path, fv}
		faceList = append(faceList, face)

		if i%block == block-1 {
			err = coll.Insert(faceList...)
			if err != nil {
				log.Fatal(err)
			}
			faceList = faceList[:0]
		}
	}

	elapsed := time.Since(t1)
	fmt.Println("Time elapsed: ", elapsed)

}
