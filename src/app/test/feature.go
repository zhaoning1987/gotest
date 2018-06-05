package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	xlog "github.com/qiniu/xlog.v1"
	"qiniu.com/ava/argus/feature_group_private/feature"
	"qiniu.com/ava/argus/feature_group_private/proto"
)

func httpGet() {
	resp, err := http.Get("http://odum9helk.qnssl.com/resource/gogopher.jpg?imageInfo")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func httpPost() {
	resp, err := http.Post("http://www.01happy.com/demo/accept.php",
		"application/x-www-form-urlencoded",
		strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func httpPostForm() {
	resp, err := http.PostForm("http://www.01happy.com/demo/accept.php",
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))

}

func httpDo() {
	client := &http.Client{}

	reqURL := "http://ava-serving-gate.cs.cg.dora-internal.qiniu.io:5001/v1/eval/image-feature"

	reqData := map[string]interface{}{"data": map[string]string{"uri": "http://oayjpradp.bkt.clouddn.com/age_gender_test.png"}}

	msg, err := json.Marshal(reqData)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(msg))
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()
	xl := xlog.FromContextSafe(ctx)
	req.Header.Set("X-Reqid", xl.ReqId())

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "QiniuStub uid=1&ut=2")
	req.Header.Set("User-Agent", "Golang qiniu/rpc package")
	req.Header.Set("X-Reqid", xl.ReqId())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func main() {
	// httpDo()

	//========= image-feature
	// reqURL := "http://ava-serving-gate.cs.cg.dora-internal.qiniu.io:5001"
	// imgURL := "http://oayjpradp.bkt.clouddn.com/age_gender_test.png"
	// timout := time.Duration(0 * time.Second)
	// image := feature.NewImageFeature(reqURL, timout, 16384)
	// ctx := context.Background()
	// fv, err := image.Image(ctx, proto.ImageURI(imgURL))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(fv)

	//========== face-feature
	// FaceBoxes 方法返回指定图片中的人脸坐标，返回值为array,包含0到多个检测到的人脸坐标
	// 人脸坐标范围为正方形4个点的x，y坐标。 照片左上角坐标为（0，0），x表示离该坐标的水平像素差值，y为垂直像素差值
	// 例子：包含一张人脸的图片返回值为 [{[[526 643] [1332 643] [1332 1652] [526 1652]] 0.9999474}]， 最后的是score准确率
	// 例子：如果包含两张人脸，则返回值为如下形式 {[[762 0] [1238 0] [1238 621] [762 621]] 0.9999913} {[[73 16] [537 16] [537 618] [73 618]] 0.9999831}]
	// 例子：如果不包含人脸，则返回值为[]

	// Face 方法返回人脸图片的特征值，需要指定人脸的范围，为4个x，y坐标值，可由FaceBoxes方法获得

	reqURL := "http://ava-serving-gate.cs.cg.dora-internal.qiniu.io:5001"
	imgURL := "http://oayjpradp.bkt.clouddn.com/age_gender_test.png1"
	// imgURL := "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1528783194&di=1f0587e1eb038d24102cc41254be8e93&imgtype=jpg&er=1&src=http%3A%2F%2Fimg.zcool.cn%2Fcommunity%2F0176565a1e2322a80120908d93209c.png%401280w_1l_2o_100sh.png"
	// imgURL := "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1528180701580&di=58f7ca5c2b97cc9e871f9f62c88c212c&imgtype=0&src=http%3A%2F%2Fe.hiphotos.baidu.com%2Fzhidao%2Fpic%2Fitem%2Fa08b87d6277f9e2faba04ea51a30e924b999f382.jpg"
	// imgURL := "/Users/zhaoning/Desktop/download.jpeg"
	// input := []byte(imgURL)
	// imgBase64 := base64.StdEncoding.EncodeToString(input)
	timout := time.Duration(0 * time.Second)
	face := feature.NewFaceFeature(reqURL, timout, 2048)
	ctx := context.Background()

	faceBoudingBox, err := face.FaceBoxes(ctx, proto.ImageURI(imgURL))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(faceBoudingBox)

	for _, box := range faceBoudingBox {
		fv, err := face.Face(ctx, proto.ImageURI(imgURL), box.Pts)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(len(fv))
	}

}
