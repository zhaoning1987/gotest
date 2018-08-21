package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MockEnv struct {
	param       map[string]string
	fileContent string
}

func (b *MockEnv) generate() *authstub.Env {
	begin := "----------------------------181519181644778910353215\n"
	end := "----------------------------181519181644778910353215--\n"
	var buf bytes.Buffer
	for k, v := range b.param {
		buf.WriteString(begin)
		buf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"\n\n", k))
		buf.WriteString(v)
		buf.WriteString("\n")
	}
	if b.fileContent != "" {
		buf.WriteString(begin)
		buf.WriteString("Content-Disposition: form-data; name=\"file\"; filename=\"csv\"\n")
		buf.WriteString("Content-Type: application/octet-stream\n\n")
		buf.WriteString(b.fileContent)
		buf.WriteString("\n")
	}
	buf.WriteString(end)

	req := &http.Request{Header: http.Header{"Content-Type": []string{"multipart/form-data"}}}
	req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
	// _ = req.ParseMultipartForm(10000)
	resp := &MockResponse{}
	return &authstub.Env{Req: req, W: resp}
}

type MockResponse struct{}

func (r *MockResponse) Header() http.Header         { return http.Header{} }
func (r *MockResponse) Write(b []byte) (int, error) { return 0, nil }
func (r *MockResponse) WriteHeader(statusCode int)  {}
func main() {
	env := MockEnv{param: make(map[string]string)}
	env.param["group_name"] = "g123"
	fmt.Println(env.generate())
}
