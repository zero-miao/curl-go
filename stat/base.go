package stat

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-22 09:57
// @File   : base.go
// @Project: curl-go/stat
// ==========================

type Protocol string

const (
	ProtocolHTTP Protocol = "HTTP"
)

// 通过命令行的简单参数, 转换成对应协议所需的参数.
type RequestConfig struct {
	Protocol   string
	Method     string
	Body       string
	OutputFile string
	Headers    []string
	Timeout    time.Duration
	Url        string
}

func (c *RequestConfig) Make() (interface{}, error) {
	switch Protocol(c.Protocol) {
	case ProtocolHTTP:
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
