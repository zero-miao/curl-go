package cli

import (
	"fmt"
	"github.com/zero-miao/curl-go/req"
	"github.com/zero-miao/curl-go/run"
	"io"
	"os"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-26 14:42
// @File   : reporter.go
// @Project: curl-go/cli
// ==========================

type ReportCLIConfig struct {
	Reporter  string `json:"reporter_class" yaml:"reporter_class"`
	ResultCSV string `json:"result_csv" yaml:"result_csv"`
	// 先过滤, 后抽样, 再统计.
	Filters   []string `json:"filters" yaml:"filters"`
	Reservoir int      `json:"reservoir" yaml:"reservoir"`
	TP        []int    `json:"tp" yaml:"tp"`
}

func (c *ReportCLIConfig) MakeReporter() (run.Reporter, error) {
	for _, item := range c.TP {
		if item >= 100 || item <= 0 {
			return nil, fmt.Errorf("TP should in (0, 100)")
		}
	}
	var filters []func(v1 *req.ReportV1) bool
	for _, item := range c.Filters {
		switch item {
		case "correct":
			filters = append(filters, func(v1 *req.ReportV1) bool {
				return v1.IsCorrect
			})
		case "error":
			filters = append(filters, func(v1 *req.ReportV1) bool {
				return v1.Err != nil
			})
		}
	}
	var writer io.WriteCloser
	if c.ResultCSV != "" {
		if c.ResultCSV == "-" {
			writer = os.Stdout
		} else if temp, err := os.OpenFile(c.ResultCSV, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return nil, err
		} else {
			writer = temp
		}
	}
	switch c.Reporter {
	case "default":
		return &run.Statistic{
			WriterCSV: writer,
			Filter: func(v1 *req.ReportV1) bool {
				for _, f := range filters {
					if !f(v1) {
						return false
					}
				}
				return true
			},
			SampleReservoir: c.Reservoir,
			ReservoirIndex:  0,
			TP:              c.TP,
		}, nil
	default:
		return nil, fmt.Errorf("invalid reporter: %s", c.Reporter)
	}
}
