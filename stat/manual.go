package stat

import (
	"github.com/zero-miao/curl-go/req"
	"github.com/zero-miao/curl-go/run"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-22 09:48
// @File   : manual.go
// @Project: curl-go/stat
// ==========================

type ManualConfig struct {
	// 并发数, 必须指定
	Concurrency int
	// 总请求数
	Count int
	// 总运行时间, 与 Count 二选一. 作为执行的停止条件.
	TimeLimit time.Duration
}

// 手动执行模式, 请求, 控制等参数均需要人为控制.
type ManualMode struct {
	ModeArgs ManualConfig

	Runner   req.Runner
	Requests []interface{}
	Reporter run.Reporter
}

func (m *ManualMode) Run() {
	var controller run.Controller

	args := m.ModeArgs
	if args.TimeLimit == 0 && args.Count == 0 {
		args.Count = 1
	}
	controller = &run.BasicController{
		Concurrency: args.Concurrency,
		Count:       args.Count,
		Period:      args.TimeLimit,
		Requests:    m.Requests,
	}
	controller.Control(m.Runner, m.Reporter)
	// 渲染数据
	m.Reporter.Result(true)
}
