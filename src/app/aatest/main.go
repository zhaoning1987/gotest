package main

import (
	"net/http"
	"os"
	"runtime"
	"runtime/debug"

	"app/aatest/service"

	"github.com/qiniu/http/restrpc.v1"
	"github.com/qiniu/http/servestk.v1"
	"github.com/qiniu/log.v1"
)

var (
	PORT_HTTP string
)

func init() {
	PORT_HTTP = os.Getenv("PORT_HTTP")
}

func main() {
	prefix := "v1/test"
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetOutputLevel(0)

	alMux := servestk.New(restrpc.NewServeMux(), func(
		w http.ResponseWriter, req *http.Request, f func(http.ResponseWriter, *http.Request)) {
		// req.Header.Set("Authorization", "QiniuStub uid=1&ut=0")
		f(w, req)
	})

	alMux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok feature_group"))
	})

	mux := http.NewServeMux()
	alMux.SetDefault(mux)

	service, err := service.NewTestService()
	if nil != err {
		log.Error("NewTestService failed", err)
		return
	}

	router := restrpc.Router{
		PatternPrefix: prefix,
		Mux:           alMux,
	}
	router.Register(service)

	runtime.GC()
	debug.FreeOSMemory()

	if err := http.ListenAndServe(":9091", alMux); err != nil {
		log.Errorf("feature group private server start error: %v", err)
	}

	log.Error("shutdown...")
}
