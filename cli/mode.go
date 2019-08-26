package cli

import (
	"fmt"
	"github.com/zero-miao/curl-go/mode"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-26 14:44
// @File   : mode.go
// @Project: curl-go/cli
// ==========================

type ModeCLIConfig struct {
	Mode string `json:"mode_class" yaml:"mode_class"`

	// Manual Mode
	// 并发数, 必须指定
	Concurrency int `json:"concurrency" yaml:"concurrency"`
	// 总请求数
	Count int `json:"count" yaml:"count"`
	// 总运行时间, 与 Count 二选一. 作为执行的停止条件.
	TimeLimit string `json:"time_limit" yaml:"time_limit"`
}

func (c *ModeCLIConfig) MakeMode(requests []interface{}, runner Runner, reporter Reporter) (mode.I, error) {
	switch c.Mode {
	case "manual":
		timeLimit, err := time.ParseDuration(c.TimeLimit)
		if err != nil {
			return nil, fmt.Errorf("invalid time_limit: %s", c.TimeLimit)
		}
		if runner == nil {
			return nil, fmt.Errorf("no runner")
		}
		runner_, err := runner.MakeRunner()
		if err != nil {
			return nil, err
		}
		if reporter == nil {
			return nil, fmt.Errorf("no reporter")
		}
		reporter_, err := reporter.MakeReporter()
		if err != nil {
			return nil, err
		}
		args := mode.ManualConfig{
			Concurrency: c.Concurrency,
			Count:       c.Count,
			TimeLimit:   timeLimit,
		}

		return &mode.ManualMode{
			ModeArgs: args,
			Requests: requests,
			Runner:   runner_,
			Reporter: reporter_,
		}, nil
	default:
		return nil, fmt.Errorf("mode not support: %s", c.Mode)
	}
}
