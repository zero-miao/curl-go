package entry

import (
	"bytes"
	"fmt"
	"github.com/zero-miao/curl-go/mode"
	"github.com/zero-miao/curl-go/req"
	"github.com/zero-miao/curl-go/run"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-22 13:59
// @File   : req.go
// @Project: curl-go/entry
// ==========================

// 通过命令行的简单参数, 转换成对应协议所需的参数.
type RequestCLIConfig struct {
	Protocol       string
	Method         string
	Body           string
	Headers        []string
	Timeout        time.Duration
	Url            string
	RecordRemoteIp bool

	CorrectSMax  int
	CorrectSMin  int
	CorrectSEnum []int
	WrongSMax    int
	WrongSMin    int
	WrongSEnum   []int
}

func (c *RequestCLIConfig) Make() (interface{}, error) {
	switch mode.Protocol(c.Protocol) {
	case mode.ProtocolHTTP:
		u, err := url.Parse(c.Url)
		if err != nil {
			return nil, fmt.Errorf("url invalid")
		}

		var httpHeaders = make(http.Header)
		for _, item := range c.Headers {
			temp := strings.SplitN(item, "=", 2)
			if len(temp) != 2 {
				return nil, fmt.Errorf("headers invalid")
			}
			httpHeaders.Add(temp[0], temp[1])
		}
		var httpBody io.ReadCloser = nil
		if c.Body != "" {
			httpBody = ioutil.NopCloser(bytes.NewReader([]byte(c.Body)))
		}
		return &http.Request{
			Method: c.Method,
			URL:    u,
			Header: httpHeaders,
			Body:   httpBody,
		}, nil
	default:
		return nil, fmt.Errorf("protocol not support")
	}
}

func (c *RequestCLIConfig) MakeRunner() (req.Runner, error) {
	switch mode.Protocol(c.Protocol) {
	case mode.ProtocolHTTP:
		cscheckmap := map[int]byte{}
		if c.CorrectSEnum != nil {
			for _, item := range c.CorrectSEnum {
				cscheckmap[item] = 0
			}
		}
		wscheckmap := map[int]byte{}
		if c.WrongSEnum != nil {
			for _, item := range c.WrongSEnum {
				wscheckmap[item] = 0
			}
		}
		sChecker := func(code int) bool {
			if _, ok := cscheckmap[code]; ok {
				return true
			} else if len(cscheckmap) > 0 {
				return false
			} else if _, ok := wscheckmap[code]; ok {
				return false
			} else if len(wscheckmap) > 0 {
				return true
			}
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
			Timeout:           c.Timeout,
			StatusCodeChecker: sChecker,
			RetrieveBody:      false,
			BodyChecker:       nil,
		}, nil
	default:
		return nil, fmt.Errorf("protocol not support")
	}
}

type ReportCLIConfig struct {
	Reporter   string
	ReportFile string

	Filters   []string
	Reservoir int
	TP        []int
}

func (c *ReportCLIConfig) Make() (run.Reporter, error) {
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
	if c.ReportFile != "" {
		if c.ReportFile == "-" {
			writer = os.Stdout
		} else if temp, err := os.OpenFile(c.ReportFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
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
		return nil, fmt.Errorf("invalid reporter")
	}
}
