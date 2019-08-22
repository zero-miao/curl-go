package stat

import (
	"github.com/zero-miao/curl-go/req"
	"github.com/zero-miao/curl-go/run"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-21 12:30
// @File   : strategy.go
// @Project: curl-go/stat
// ==========================

type AutoConfig struct {
	// 正确率
	CorrectRate float32
	// p90
	P90 time.Duration
	// p99
	P99 time.Duration
	// 最大并发
	MaxConcurrency int
	// 最小并发
	MinConcurrency int
	// 总执行时间
	MaxPeriod time.Duration
	// 单个请求的超时时间
	Timeout time.Duration
}

type AutoMode struct {
	ModeArgs AutoConfig

	Runner   req.Runner
	Requests []interface{}
	Reporter run.Reporter
}

// 开始调度
func (m *AutoMode) Run() {

	var controller run.Controller

	args := m.ModeArgs
	controller = &run.BasicController{
		Concurrency: args.MinConcurrency,
		Count:       10,
		Period:      args.MaxPeriod,
		Requests:    m.Requests,
	}
	controller.Control(m.Runner, m.Reporter)
	// 渲染数据
	m.Reporter.Result(true)
}
