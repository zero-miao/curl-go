package req

import (
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-21 12:38
// @File   : http_impl.go
// @Project: curl-go/req
// ==========================

type HTTPReqRunner struct {
	// 最大请求时间
	Timeout           time.Duration
	StatusCodeChecker func(code int) bool
	// 是否获取响应体
	RetrieveBody bool
	BodyChecker  func(data []byte) bool

	RecordServerIp bool
}

func (hr *HTTPReqRunner) Run(id string, request interface{}) (report *ReportV1) {
	report = &ReportV1{
		Id:        id,
		IsCorrect: false,
	}

	r, ok := request.(*http.Request)
	if !ok {
		return nil
	}
	if hr.RecordServerIp {
		trace := &httptrace.ClientTrace{
			GotConn: func(connInfo httptrace.GotConnInfo) {
				report.Msg = connInfo.Conn.RemoteAddr()
			},
		}
		r = r.WithContext(httptrace.WithClientTrace(r.Context(), trace))
	}

	client := http.Client{
		Timeout: hr.Timeout,
	}
	st := time.Now()
	res, err := client.Do(r)
	cost := time.Now().Sub(st)
	report.Latency = cost

	if err != nil {
		report.Err = err
		return
	}
	report.StatusCode = res.StatusCode
	if hr.StatusCodeChecker != nil && !hr.StatusCodeChecker(res.StatusCode) {
		return
	}

	if hr.RetrieveBody {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			report.Err = err
			return
		}
		if hr.BodyChecker != nil && !hr.BodyChecker(data) {
			return
		}
	}
	report.IsCorrect = true
	return
}
