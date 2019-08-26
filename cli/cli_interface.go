package cli

import (
	"github.com/zero-miao/curl-go/mode"
	"github.com/zero-miao/curl-go/req"
	"github.com/zero-miao/curl-go/run"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-26 14:22
// @File   : cli_interface.go
// @Project: curl-go/entry
// ==========================

type Request interface {
	MakeRequest() (interface{}, error)
}

type Runner interface {
	MakeRunner() (req.Runner, error)
}

type Reporter interface {
	MakeReporter() (run.Reporter, error)
}

type Mode interface {
	MakeMode(requests []interface{}, runner Runner, reporter Reporter) (mode.I, error)
}
