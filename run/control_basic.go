package run

import (
	"context"
	"fmt"
	"github.com/zero-miao/curl-go/req"
	"math/rand"
	"sync"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-21 10:50
// @File   : control.go
// @Project: curl-go/run
// ==========================

type BasicController struct {
	// 并发数, 必须指定
	Concurrency int
	// 总请求数
	Count int
	// 总运行时间
	Period time.Duration
	// 请求
	Requests []interface{}
}

func (c *BasicController) String() string {
	return fmt.Sprintf("BasicController(并发=%d, QPS=%.3f)", c.Concurrency, float64(c.Count)/c.Period.Seconds())
}

func (c *BasicController) Control(runner req.Runner, reporter Reporter) {
	if c.Concurrency <= 0 {
		panic("concurrency should be specified")
	}
	if c.Count > 0 {
		fmt.Printf("control begin: 并发数=%d, 总运行次数=%d\n", c.Concurrency, c.Count)
	} else if c.Period > 0 {
		fmt.Printf("control begin: 并发数=%d, 总运行时间~=%s\n", c.Concurrency, c.Period.String())
	} else {
		panic("参数错误")
	}
	fmt.Println("check the above information, bench will begin 3s later...")
	time.Sleep(3 * time.Second)

	var wg sync.WaitGroup
	ctx := context.Background()
	st := time.Now()
	if c.Period > 0 {
		ctx, _ = context.WithTimeout(ctx, c.Period)
	}
	total := 0
	for i := 0; i < c.Concurrency; i++ {
		wg.Add(1)
		go func(ctx context.Context, uid int, requests []interface{}) {
			defer wg.Done()
			l := len(requests)
			randStart := rand.Intn(l)

			if c.Count > 0 {
				for total < c.Count {
					request := requests[randStart]
					select {
					case <-ctx.Done():
						return
					default:
						total++
						id := fmt.Sprintf("%d-%d-%d", uid, time.Now().UnixNano(), rand.Int())
						res := runner.Run(id, request)
						reporter.StreamReport(res)
					}
					randStart = (randStart + 1) % l
				}
			} else {
				for {
					request := requests[randStart]
					select {
					case <-ctx.Done():
						return
					default:
						total++
						id := fmt.Sprintf("%d-%d-%d", uid, time.Now().UnixNano(), rand.Int())
						res := runner.Run(id, request)
						reporter.StreamReport(res)
					}
					randStart = (randStart + 1) % l
				}
			}
		}(ctx, i, c.Requests)
	}
	wg.Wait()
	reporter.ResultReport(&BasicController{
		Period:      time.Now().Sub(st),
		Count:       total,
		Concurrency: c.Concurrency,
	})
	return
}
