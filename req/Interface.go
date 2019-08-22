package req

import (
	"fmt"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-21 10:11
// @File   : Interface.go
// @Project: curl-go/req
// ==========================

type ReportV1 struct {
	Id         string
	Latency    time.Duration
	IsCorrect  bool
	StatusCode int
	Err        error
	Msg        interface{}
}

func (r *ReportV1) CSV() string {
	return fmt.Sprintf("\"%s\",\"%f\",\"%v\",\"%d\",\"%v\",\"%v\"\n", r.Id, r.Latency.Seconds(), r.IsCorrect, r.StatusCode, r.Err, r.Msg)
}

// 执行 http 请求. 依据不同的要求, 实现 Report 方法, 输出不同的内容. (如对端 ip)
type Runner interface {
	// 执行请求, 并报告, 不同的请求式, 通过不同的结构体实现.
	Run(id string, r interface{}) *ReportV1
}
