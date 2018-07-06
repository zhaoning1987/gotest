package service

import (
	"context"
	"io"

	"github.com/qiniu/http/restrpc.v1"
)

type parameter struct {
	Name string `json:"name"`
}

type ITestService interface {
	GetTest_(context.Context,
		*struct {
			CmdArgs []string
			Par1    string        "par1"
			File    io.ReadCloser "abc"
		},
		*restrpc.Env,
	) (*ret, error)
}

// type ITestService interface {
// 	PostTest_(context.Context,
// 		*struct {
// 			CmdArgs []string
// 			Config  parameter `json:"config"`
// 		},
// 		*restrpc.Env,
// 	) error
// }

// type ITestService interface {
// 	PostTest_(context.Context,
// 		*BaseReq,
// 		*restrpc.Env,
// 	) error
// }
