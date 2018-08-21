package service

import (
	"context"
	"fmt"
	"io"

	"github.com/qiniu/http/restrpc.v1"
)

type TestService struct {
}

type BaseReq struct {
	CmdArgs []string
	ReqBody io.ReadCloser
}

type CustomError struct {
	code string
	msg  string
}

func (e *CustomError) Error() string {
	return e.code + "||" + e.msg
}

// func (e *CustomError) RpcError() string {
// 	return e.code + "||" + e.msg
// }
var _ ITestService = new(TestService)

func NewTestService() (ITestService, error) {
	s := &TestService{}
	return s, nil
}

type ret struct {
	A string
	B string
}

type Resp struct {
	Test string      `json:"test"`
	PPP  interface{} `json:"meta"`
}

func (s *TestService) GetTest(ctx context.Context,
	args *struct {
		CmdArgs []string
		Param   string `json:"param"`
	},
	env *restrpc.Env,
) {
	fmt.Println("===========")
	fmt.Println(args.Param)
}

func (s *TestService) PostGet_(ctx context.Context,
	args *struct {
		CmdArgs []string
		// Param1  Parameter `json:"param1"`
	},
	env *restrpc.Env,
) (res Resp, err error) {
	fmt.Println("sfdd")
	result := get(args.CmdArgs[0])
	res.Test = result.Name
	res.PPP = result.Param2
	return
}

func (s *TestService) PostSet_(ctx context.Context,
	args *struct {
		CmdArgs []string
		Param1  Parameter `json:"param1"`
	},
	env *restrpc.Env,
) {
	fmt.Println(args.CmdArgs[0])
	fmt.Println(args.Param1.Name)
	fmt.Println(args.Param1.Param2)

	p := &Parameter{
		Name:   args.CmdArgs[0],
		Param2: args.Param1.Param2,
	}

	insert(p)

	// env.W.Header().Set("Content-Disposition", "attachment;fileName="+"a.txt")

	// env.W.Write([]byte("zhao ning"))

	// env.W.WriteHeader(999)
	// env.W.Write([]byte("wwwww ww"))
	// body := ioutil.NopCloser(bytes.NewReader([]byte("zhao ning")))
	// return &ret{}, nil
	// return errors.New("invalid arguments")
	// return httputil.NewError(403, "hello")
	// return httputil.NewRpcError(301, 100, "", "www")
	// return &CustomError{"sss", "www"}
	// fmt.Println()
	// fmt.Println(env.Req.Header.Get("Content-Type"))
	// paraName := args.CmdArgs[0]
	// if 0 == len(paraName) {
	// 	return errors.New("invalid arguments")
	// }
	// fmt.Println("paraName", paraName)
	// fmt.Println("Par1", args.Par1)

	// bb, _ := ioutil.ReadAll(args.File)
	// fmt.Println("ppppppp", string(bb))
	// fmt.Println("post data:", env.Req.FormValue("par1"))
	// formFile, header, err := env.Req.FormFile("abc")
	// if err != nil {
	// 	log.Printf("Get form file failed: %s\n", err)
	// 	return nil
	// }
	// defer formFile.Close()
	// fmt.Println("======]]]=====")
	// fmt.Println(header.Filename)
	// bt, _ := ioutil.ReadAll(formFile)
	// fmt.Println(string(bt))
	// return nil
}

// func (s *TestService) PostTest_(ctx context.Context,
// 	args *BaseReq,
// 	env *restrpc.Env,
// ) error {
// 	fmt.Println()
// 	fmt.Println(env.Req.Header.Get("Content-Type"))
// 	paraName := args.CmdArgs[0]
// 	if 0 == len(paraName) {
// 		return errors.New("invalid arguments")
// 	}
// 	fmt.Println("paraName", paraName)

// 	fmt.Println("post data:", env.Req.FormValue("par1"))
// 	formFile, header, err := env.Req.FormFile("abc")
// 	if err != nil {
// 		log.Printf("Get form file failed: %s\n", err)
// 		return nil
// 	}
// 	defer formFile.Close()
// 	fmt.Println("======]]]=====")
// 	fmt.Println(header.Filename)
// 	bt, _ := ioutil.ReadAll(formFile)
// 	fmt.Println(string(bt))
// 	return nil
// }

// func (s *TestService) PostTest_(ctx context.Context,
// 	args *struct {
// 		CmdArgs []string
// 		Config  parameter `json:"config"`
// 	},
// 	env *restrpc.Env,
// ) error {
// 	fmt.Println()
// 	fmt.Println(env.Req.Header.Get("Content-Type"))
// 	paraName := args.CmdArgs[0]
// 	if 0 == len(paraName) {
// 		return errors.New("invalid arguments")
// 	}
// 	fmt.Println("paraName", paraName)
// 	fmt.Println("Config", args.Config)
// 	return nil
// }
