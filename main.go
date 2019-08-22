package main

import (
	"fmt"
	"github.com/zero-miao/curl-go/entry"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func main() {
	if err := entry.Entry(); err != nil {
		os.Exit(1)
	}
}

func DoReq(u string, uid interface{}) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Printf("%v: code=None, err=%v\n", uid, err)
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("%v: %s\n", uid, connInfo.Conn.RemoteAddr())
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{}
	st := time.Now()
	resp, err := client.Do(req)
	cost := time.Now().Sub(st)
	if err != nil {
		fmt.Printf("%v: code=None, err=%v, cost=%s\n", uid, err, cost.String())
	} else {
		fmt.Printf("%v: code=%d, cost=%s\n", uid, resp.StatusCode, cost.String())
	}
}
