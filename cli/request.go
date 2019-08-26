package cli

import (
	"bytes"
	"fmt"
	"github.com/zero-miao/curl-go/mode"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-26 14:27
// @File   : http_request.go
// @Project: curl-go/cli
// ==========================

type RequestCLIConfig struct {
	Protocol string `json:"protocol" yaml:"protocol"`

	// HTTP
	Method  string   `json:"method" yaml:"method"`
	Headers []string `json:"headers" yaml:"headers"`
	Body    string   `json:"body" yaml:"body"`
	Url     string   `json:"url" yaml:"url"`
}

func (c *RequestCLIConfig) MakeRequest() (interface{}, error) {
	switch mode.Protocol(c.Protocol) {
	case mode.ProtocolHTTP:
		if c.Url == "" {
			return nil, fmt.Errorf("invalid url: %s", c.Url)
		}
		u, err := url.Parse(c.Url)
		if err != nil {
			return nil, fmt.Errorf("invalid url: %s", c.Url)
		}
		var httpHeaders = make(http.Header)
		for _, item := range c.Headers {
			temp := strings.SplitN(item, "=", 2)
			if len(temp) != 2 {
				return nil, fmt.Errorf("headers invalid: %s", item)
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
		return nil, fmt.Errorf("protocol not support: %s", c.Protocol)
	}
}
