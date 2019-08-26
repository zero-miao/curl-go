package cli

import (
	"fmt"
	"github.com/zero-miao/curl-go/mode"
	"github.com/zero-miao/curl-go/req"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-26 14:34
// @File   : runner.go
// @Project: curl-go/cli
// ==========================

type RunnerCLIConfig struct {
	Protocol string `json:"protocol" yaml:"protocol"`

	// HTTP
	RecordRemoteIp bool   `json:"record_remote_ip" yaml:"record_remote_ip"`
	Timeout        string `json:"timeout" yaml:"timeout"`
	CorrectSMax    int    `json:"correct_status_max" yaml:"correct_status_max"`
	CorrectSMin    int    `json:"correct_status_min" yaml:"correct_status_min"`
	CorrectSEnum   []int  `json:"correct_status_enum" yaml:"correct_status_enum"`
	WrongSMax      int    `json:"wrong_status_max" yaml:"wrong_status_max"`
	WrongSMin      int    `json:"wrong_status_min" yaml:"wrong_status_min"`
	WrongSEnum     []int  `json:"wrong_status_enum" yaml:"wrong_status_enum"`
}

func (c *RunnerCLIConfig) MakeRunner() (req.Runner, error) {
	switch mode.Protocol(c.Protocol) {
	case mode.ProtocolHTTP:
		timeout, err := time.ParseDuration(c.Timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout: %s", c.Timeout)
		}

		correctStatusSet := map[int]byte{}
		if c.CorrectSEnum != nil {
			for _, item := range c.CorrectSEnum {
				correctStatusSet[item] = 0
			}
		}
		wrongStatusSet := map[int]byte{}
		if c.WrongSEnum != nil {
			for _, item := range c.WrongSEnum {
				wrongStatusSet[item] = 0
			}
		}
		sChecker := func(code int) bool {
			if _, ok := correctStatusSet[code]; ok {
				// 白名单, 命中
				return true
			} else if len(correctStatusSet) > 0 {
				// 白名单, 未命中
				return false
			} else if _, ok := wrongStatusSet[code]; ok {
				// 黑名单, 命中
				return false
			} else if len(wrongStatusSet) > 0 {
				// 黑名单, 未命中
				return true
			}
			// 名单均为空. 顺序: 正确最大, 正确最小, 错误最大, 错误最小
			if c.CorrectSMax > 0 && code > c.CorrectSMax {
				return false
			}
			if c.CorrectSMin > 0 && code < c.CorrectSMin {
				return false
			}
			if c.WrongSMax > 0 && code > c.WrongSMax {
				return true
			}
			if c.WrongSMin > 0 && code < c.WrongSMin {
				return true
			}
			return true
		}

		return &req.HTTPReqRunner{
			RecordServerIp:    c.RecordRemoteIp,
			Timeout:           timeout,
			StatusCodeChecker: sChecker,
			RetrieveBody:      false,
			BodyChecker:       nil,
		}, nil
	default:
		return nil, fmt.Errorf("protocol not support: %s", c.Protocol)
	}
}
