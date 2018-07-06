package dbstorage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	xlog "github.com/qiniu/xlog.v1"
)

const (
	BASE64_PREFIX = "data:application/octet-stream;base64,"
)

type Mode int

const (
	REPLACE Mode = iota
	APPEND
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

type SafeMap struct {
	Map   map[string]struct{}
	Mutex sync.Mutex
}

func NewSafeMap() *SafeMap {
	return &SafeMap{Map: make(map[string]struct{}, 0)}
}

type SafeFile struct {
	File     *os.File
	Path     string
	EditMode Mode
	NeedLock bool
	Mutex    sync.Mutex
}

func NewSafeFile(path string, mode Mode, needLock bool) (sf *SafeFile, err error) {
	var f *os.File
	if mode == APPEND {
		f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	}
	if err != nil {
		return nil, err
	}
	return &SafeFile{
		File:     f,
		Path:     path,
		EditMode: mode,
		NeedLock: needLock,
	}, nil
}

func (f *SafeFile) Close() {
	if f.File != nil {
		f.File.Close()
	}
}

func (f *SafeFile) Write(msg string) (err error) {
	if f.NeedLock {
		f.Mutex.Lock()
		defer f.Mutex.Unlock()
	}
	if f.EditMode == APPEND {
		content := fmt.Sprintf("%s\n", msg)
		_, err = f.File.Write([]byte(content))
	} else {
		err = ioutil.WriteFile(f.Path, []byte(msg), 0644)
	}
	return
}

func CreatePath(log *xlog.Logger, path string) {
	exist, err := pathExists(path)
	if err != nil {
		log.Fatalf("get directory [%s] error: %v\n", path, err)
	}
	if !exist {
		log.Infof("no directory [%s]\n", path)
		//create folder
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatalf("create directory [%s] failed: %v\n", path, err)
		} else {
			log.Infof("create directory [%s] success\n", path)
		}
	}
}

func Substring(source string, start int, end int) string {
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

func GetTagAndDesc(name string) (tag, desc string) {
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

func getSha1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
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
