package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ss interface {
	fun()
}

type ssStruct struct {
	p1 string
}

func (s *ssStruct) fun() {
	s.p1 = "abc"
}

func generateFace() {
	dir := "/Users/zhaoning/Desktop/face_img"
	var list []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			list = append(list, "http://pfraaslhy.bkt.clouddn.com/"+f.Name())
		}
		return nil
	})
	if err != nil {
		panic("error when filepath.Walk")
		return
	}

	msg := strings.Join(list, "\n")
	err = ioutil.WriteFile("./face_url_list", []byte(msg), 0644)
}

func generateImage() {
	dir := "/Users/zhaoning/Desktop/image_img"
	var list []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			list = append(list, "http://pfradxv0k.bkt.clouddn.com/"+f.Name())
		}
		return nil
	})
	if err != nil {
		panic("error when filepath.Walk")
		return
	}

	msg := strings.Join(list, "\n")
	err = ioutil.WriteFile("./image_url_list", []byte(msg), 0644)
}

const (
	PENDING_READ  byte = 0x01
	PENDING_WRITE byte = 0x02
)

func test123(src [2]int) {
	fmt.Println(src)
	src[1] = 60
	// src["www"] = 1
}

type interface1 interface {
	todo()
}
type str1 struct {
	AA int    `json:"aa"`
	BB string `json:"bb"`
}

func (t str1) todo() {

}

type str2 struct {
	str1
	CC string `json:"cc"`
}

func main() {
	var i interface1
	i = str1{1, "b"}
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	fmt.Println(t)
	fmt.Println(v)
	fmt.Println(v.Type())
	fmt.Println(v.Kind())
	fmt.Println(v.CanSet())

	if !v.CanSet() {
		v = reflect.New(t)
	}
	fmt.Println("======")
	fmt.Println(v)
	fmt.Println(v.Type())
	fmt.Println(v.Kind())
	fmt.Println(v.CanSet())

	vv := v.Elem()
	vv.Field(0).Set(reflect.ValueOf(4))
	fmt.Println(v)

	v2 := v.Interface().(*str1)
	fmt.Println(v2)
}
func main1() {
	// generateFace()
	// generateImage()

	// http.Handle("/metrics", promhttp.Handler())
	// log.Fatal(http.ListenAndServe(":8081", nil))

	// cpuTemp := prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Namespace:   "namespace",
	// 	Subsystem:   "subnamespace",
	// 	Name:        "cpu_temperature_celsius",
	// 	Help:        "Current temperature of the CPU.",
	// 	ConstLabels: prometheus.Labels{"key1": "val1", "key2": "val2"},
	// })

	cpuTemp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "namespace",
			Subsystem:   "subnamespace",
			Name:        "cpu_temperature_celsius",
			Help:        "Current temperature of the CPU.",
			ConstLabels: prometheus.Labels{"key1": "val1", "key2": "val2"},
		},
		[]string{"device1", "device2"},
	)

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(cpuTemp)

	// cpuTemp.Set(65.3)

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	// http.Handle("/metrics", promhttp.Handler())
	// http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		cpuTemp.WithLabelValues("devicsadfe1", "deviqqqqqce2").Inc()
		cpuTemp.With(prometheus.Labels{"device1": "devsda2", "device2": "XXX"}).Inc()
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, req)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func test(a ss) {
	a.fun()
}
