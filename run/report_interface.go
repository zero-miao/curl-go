package run

import (
	"github.com/zero-miao/curl-go/req"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-21 11:33
// @File   : report_interface.go
// @Project: curl-go/run
// ==========================

// 通过不同的策略, 去实现不同的抽样.
type Reporter interface {
	// 结果流式处理: 存盘, 抽样统计
	StreamReport(report *req.ReportV1)
	//
	ResultReport(interface{})
	// 渲染自己, 最后自己就是结果.
	Result(print bool)
}
