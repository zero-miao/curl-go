package run

import (
	"fmt"
	"github.com/zero-miao/curl-go/req"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-21 11:33
// @File   : control_interface.go
// @Project: curl-go/run
// ==========================

type Controller interface {
	fmt.Stringer
	// 调用逻辑控制
	Control(runner req.Runner, reporter Reporter)
}
