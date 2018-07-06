package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	xlog "github.com/qiniu/xlog.v1"
	"qiniu.com/argus/argus/facec/client"
	"qiniu.com/argus/feature_group_private/proto"
)

const (
	BASE64_PREFIX = "data:application/octet-stream;base64,"
)

type Config struct {
	ImageFolderPath     string `json:"image_folder_path"`
	ImageListFile       string `json:"image_list_file"`
	LoadImageFromFolder bool   `json:"load_image_from_folder"`
	ServiceHost         string `json:"service_host_url"`
	Timeout             int    `json:"http_timeout_in_millisecond"`
	MaxTryServiceTime   int    `json:"max_try_service_time"`
	MaxTryDownloadTime  int    `json:"max_try_download_time"`
	ThreadNumber        int    `json:"thread_number"`
	PoolSize            int    `json:"job_pool_size"`
	GroupName           string `json:"group_name"`
}

type FaceJob struct {
	index        int
	ctx          *context.Context
	faceGroup    *faceGroup
	imageContent []byte
	imageURI     string
	tag          proto.FeatureTag
	desc         proto.FeatureDesc
	wg           *sync.WaitGroup
}

type faceGroup struct {
	host    string
	timeout time.Duration
}

func NewFaceGroup(host string, timeout time.Duration) *faceGroup {
	return &faceGroup{
		host:    host,
		timeout: timeout,
	}
}

func (fg *faceGroup) Add(ctx context.Context, groupName string, id proto.FeatureID, uri proto.ImageURI, tag proto.FeatureTag, desc proto.FeatureDesc) (respID proto.FeatureID, err error) {
	cli := client.NewRPCClient(client.EvalEnv{Uid: 1, Utype: 0}, fg.timeout)

	if len(uri) == 0 {
		return "", errors.New("image do not contain any data")
	}

	req := map[string]interface{}{"image": map[string]string{"id": string(id), "uri": string(uri), "tag": string(tag), "desc": string(desc)}}

	resp, err := cli.DoRequestWithJson(ctx, "POST", fg.host+"/v1/face/groups/"+groupName+"/add", req)
	if err != nil {
		return "", errors.Wrap(err, "request to PostGroup_Add failed")
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read PostGroup_Add response failed")
	}

	if resp.StatusCode/100 != 2 {
		return "", errors.Errorf("PostGroup_Add return error: %s", string(content))
	}

	result := struct {
		ID proto.FeatureID `json:"id"`
	}{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", errors.Wrapf(err, "parse PostGroup_Add response failed, response body is: %s", string(content))
	}

	return result.ID, nil
}

func (fg *faceGroup) CreateGroup(ctx context.Context, groupName string) error {
	cli := client.NewRPCClient(client.EvalEnv{Uid: 1, Utype: 0}, fg.timeout)

	req := map[string]interface{}{"config": map[string]int{"capacity": 100000000}}

	resp, err := cli.DoRequestWithJson(ctx, "POST", fg.host+"/v1/face/groups/"+groupName, req)
	if err != nil {
		return errors.Wrap(err, "request to PostGroup_ failed")
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read PostGroup_ response failed")
	}

	if resp.StatusCode/100 != 2 {
		return errors.Errorf("PostGroup_ return error: %s", string(content))
	}

	return nil
}

func (fg *faceGroup) Delete(ctx context.Context, groupName string, id ...proto.FeatureID) error {
	cli := client.NewRPCClient(client.EvalEnv{Uid: 1, Utype: 0}, fg.timeout)

	req := map[string]interface{}{"ids": id}

	resp, err := cli.DoRequestWithJson(ctx, "POST", fg.host+"/v1/face/groups/"+groupName+"/delete", req)
	if err != nil {
		return errors.Wrap(err, "request to PostGroup_Delete failed")
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read PostGroup_Delete response failed")
	}

	if resp.StatusCode/100 != 2 {
		return errors.Errorf("PostGroup_Delete return error: %s", string(content))
	}

	return nil
}

func writeToFile(log *xlog.Logger, path string, msg string) {
	err := ioutil.WriteFile(path, []byte(msg), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func appendToFile(log *xlog.Logger, path string, msg string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
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

func createPath(log *xlog.Logger, path string) {
	exist, err := pathExists(path)
	if err != nil {
		xl.Fatalf("get directory [%s] error: %v\n", path, err)
	}
	if !exist {
		xl.Infof("no directory [%s]\n", path)
		//create folder
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			xl.Fatalf("create directory [%s] failed: %v\n", path, err)
		} else {
			xl.Infof("create directory [%s] success\n", path)
		}
	}
}

func getSha1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func downloadFile(url string) (content []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("error when trying to get image: %s", url))
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, errors.Errorf("failed to get image, response status code : %d", resp.StatusCode)
	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("error when trying to read the content of image : %s", url))
	}
	return content, nil
}

func getTagAndDesc(name string) (tag, desc string) {
	if name != "" {
		if i := strings.LastIndex(name, "."); i >= 0 {
			name = name[0:i]
		}
		blocks := strings.SplitN(name, "_", 2)
		if len(blocks) == 1 {
			return blocks[0], ""
		}
		return blocks[0], blocks[1]
	}
	return "", ""
}
