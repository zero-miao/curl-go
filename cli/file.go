package cli

import (
	"github.com/zero-miao/curl-go/mode"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-26 15:24
// @File   : file.go
// @Project: curl-go/cli
// ==========================

type FileCLIConfig struct {
	Protocol string              `json:"protocol" yaml:"protocol"`
	Requests []*RequestCLIConfig `json:"requests" yaml:"requests"`
	Runner   *RunnerCLIConfig    `json:"runner" yaml:"runner"`
	Reporter *ReportCLIConfig    `json:"reporter" yaml:"reporter"`
	Mode     *ModeCLIConfig      `json:"mode" yaml:"mode"`
}

func (c *FileCLIConfig) MakeMode() (mode.I, error) {
	reqs := make([]interface{}, 0)
	for _, item := range c.Requests {
		item.Protocol = c.Protocol
		req, err := item.MakeRequest()
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	c.Runner.Protocol = c.Protocol
	return c.Mode.MakeMode(reqs, c.Runner, c.Reporter)
}
