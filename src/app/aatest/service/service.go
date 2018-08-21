package service

import (
	"context"

	"github.com/qiniu/http/restrpc.v1"
)

type Parameter struct {
	Name   string      `json:"name" bson:"name"`
	Param2 interface{} `json:"param2" bson:"meta"`
}

type ITestService interface {
	// PostSet_(context.Context,
	// 	*struct {
	// 		CmdArgs []string
	// 		Param1  Parameter `json:"param1"`
	// 	},
	// 	*restrpc.Env,
	// )
	// PostGet_(context.Context,
	// 	*struct {
	// 		CmdArgs []string
	// 		// Param1  Parameter `json:"param1"`
	// 	},
	// 	*restrpc.Env,
	// ) (Resp, error)

	GetTest(context.Context,
		*struct {
			CmdArgs []string
			Param   string `json:"param"`
		},
		*restrpc.Env)
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
