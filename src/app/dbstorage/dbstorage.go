package dbstorage

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	xlog "github.com/qiniu/xlog.v1"
	"qiniu.com/argus/argus/facec/client"
	"qiniu.com/argus/feature_group_private/manager"
	"qiniu.com/argus/feature_group_private/proto"
)

var (
	conf                             Config
	sha1Dict                         *SafeMap
	count                            int64
	errorFile, sha1File, processFile *SafeFile
	threadFiles                      []*SafeFile
	outputLogPath                    = "./log/"
	internalLogPath                  = "./processlog/"
)

//====== FaceJob ======
type FaceJob struct {
	index        int
	ctx          context.Context
	faceGroup    *FaceGroup
	imageContent []byte
	imageURI     string
	tag          proto.FeatureTag
	desc         proto.FeatureDesc
	wg           *sync.WaitGroup
}

func NewFaceJob(ctx context.Context, index int, faceGroup *FaceGroup, imageContent []byte, imageURI, tag, desc string, wg *sync.WaitGroup) *FaceJob {
	return &FaceJob{
		index:        index,
		ctx:          ctx,
		faceGroup:    faceGroup,
		imageContent: imageContent,
		imageURI:     imageURI,
		tag:          proto.FeatureTag(tag),
		desc:         proto.FeatureDesc(desc),
		wg:           wg,
	}
}

func (fj *FaceJob) execute(workerIndex int) {
	defer fj.wg.Done()
	xl := xlog.FromContextSafe(fj.ctx)

	//when start to handle a face, we save its index in file
	threadFiles[workerIndex].Write(strconv.Itoa(fj.index))

	var err error

	//if image content not passed, try to get it by url
	if fj.imageContent == nil {
		for i := 0; i < conf.MaxTryDownloadTime; i++ {
			fj.imageContent, err = downloadFile(fj.imageURI)
			if err == nil {
				break
			}
		}
		if err != nil {
			xl.Infof("%s : %s\n", fj.imageURI, err.Error())
			LogError(fmt.Sprintf("%s : %s", fj.imageURI, err.Error()))
		}
	}

	if fj.imageContent != nil {
		//check if this image exist
		existed := false
		sha1 := getSha1(fj.imageContent)
		sha1Dict.Mutex.Lock()
		if _, ok := sha1Dict.Map[sha1]; ok {
			existed = true
		} else {
			sha1Dict.Map[sha1] = struct{}{}
		}
		sha1Dict.Mutex.Unlock()

		if !existed {
			//call group_add service to store the image
			imgBase64 := BASE64_PREFIX + base64.StdEncoding.EncodeToString(fj.imageContent)

			for i := 0; i < conf.MaxTryServiceTime; i++ {
				_, err = fj.faceGroup.Add(fj.ctx, conf.GroupName, proto.FeatureID(fj.imageURI), proto.ImageURI(imgBase64), fj.tag, fj.desc)
				if err == nil {
					break
				}
			}

			if err != nil {
				errMsg := err.Error()
				if strings.Contains(errMsg, manager.ErrGroupNotExist.Err) {
					xl.Fatalf("group <%s> not exist\n", conf.GroupName)
				} else if !strings.Contains(errMsg, manager.ErrFeatureExist.Err) {
					LogError(fmt.Sprintf("%s : %s", fj.imageURI, err.Error()))
				}
				xl.Infof("%s : %s\n", fj.imageURI, err.Error())
			}

			//no matter call service success or not, always write to sha1 file
			sha1File.Write(sha1)
		}
	}

	//after handling an image, update the count & save to file
	IncrementCount()
	processFile.Write(fmt.Sprintf("%d", count))
}

//====== FaceGroup ======
type FaceGroup struct {
	host    string
	timeout time.Duration
}

func NewFaceGroup(host string, timeout time.Duration) *FaceGroup {
	return &FaceGroup{
		host:    host,
		timeout: timeout,
	}
}

func (fg *FaceGroup) Add(ctx context.Context, groupName string, id proto.FeatureID, uri proto.ImageURI, tag proto.FeatureTag, desc proto.FeatureDesc) (respID proto.FeatureID, err error) {
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

func (fg *FaceGroup) CreateGroup(ctx context.Context, groupName string) error {
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

//====== Setup & Teardown action======
func Setup(ctx context.Context, config Config) int {
	lastMinIndex := 0
	conf = config
	errorFilePath := outputLogPath + "errorlist"
	processFilePath := outputLogPath + "process"
	sha1FilePath := internalLogPath + "sha1"
	sha1Dict = NewSafeMap()
	xl := xlog.FromContextSafe(ctx)

	//1.create log dir if not exist
	CreatePath(xl, outputLogPath)
	CreatePath(xl, internalLogPath)

	//2.read process log file to get the last min stop point
	var indexes []int
	err := filepath.Walk(internalLogPath, func(path string, f os.FileInfo, err error) error {
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

		content := string(dat)
		if content != "" {
			index, err := strconv.Atoi(string(dat))
			if err != nil {
				return err
			}
			indexes = append(indexes, index)
		}
		return nil
	})
	if err != nil {
		xl.Fatalf("error when filepath.Walk() on [%s]: %v\n", internalLogPath, err)
	}
	if indexes != nil {
		sort.Ints(indexes)
		lastMinIndex = indexes[0]
	}

	//3. get the uploaded files' sha1
	if file, err := os.Open(sha1FilePath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			sha1Dict.Map[scanner.Text()] = struct{}{}
		}
	}

	//4. open file for logging
	for i := 0; i < conf.ThreadNumber; i++ {
		var threadFile *SafeFile
		threadFilePath := fmt.Sprintf(internalLogPath+"%d.log", i)
		if threadFile, err = NewSafeFile(threadFilePath, REPLACE, false); err != nil {
			xl.Fatalf("error when opening file [%s]: %v\n", threadFilePath, err)
		}
		threadFiles = append(threadFiles, threadFile)
	}
	if processFile, err = NewSafeFile(processFilePath, REPLACE, true); err != nil {
		xl.Fatalf("error when opening file [%s]: %v\n", processFilePath, err)
	}
	if errorFile, err = NewSafeFile(errorFilePath, APPEND, true); err != nil {
		xl.Fatalf("error when opening file [%s]: %v\n", errorFilePath, err)
	}
	if sha1File, err = NewSafeFile(sha1FilePath, APPEND, true); err != nil {
		xl.Fatalf("error when opening file [%s]: %v\n", sha1FilePath, err)
	}

	return lastMinIndex
}

func Teardown() {
	for _, file := range threadFiles {
		file.Close()
	}
	processFile.Close()
	errorFile.Close()
	sha1File.Close()
}

func IncrementCount() {
	atomic.AddInt64(&count, 1)
}

func LogError(msg string) {
	errorFile.Write(msg)
}
