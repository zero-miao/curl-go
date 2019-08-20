package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"os"
	"sync"
	"time"
)

var rootCmd *cobra.Command

func init() {
	var concurrent int
	var count int
	var period time.Duration

	rootCmd = &cobra.Command{
		Use: "curl_go <url>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("url required")
			}
			return nil
		},
		Short: "curl_go just like curl get",
		Long:  "curl_go just like curl get",
		Run: func(cmd *cobra.Command, args []string) {
			if period == 0 && count == 0 {
				count = 1
			}
			fmt.Printf("run with: concurrent=%d, count=%d, period=%s\nurl=%s\nbench will begin 3s later...\n\n", concurrent, count, period.String(), args[0])
			time.Sleep(3 * time.Second)
			var wg sync.WaitGroup
			ctx := context.Background()
			st := time.Now()
			if period > 0 {
				ctx, _ = context.WithTimeout(ctx, period)
			}
			total := 0
			for i := 0; i < concurrent; i++ {
				wg.Add(1)
				go func(ctx context.Context, uid int) {
					defer wg.Done()
					if count > 0 {
						for total < count {
							select {
							case <-ctx.Done():
								return
							default:
								total++
								DoReq(args[0], fmt.Sprintf("%d-%d-%d", uid, time.Now().UnixNano(), rand.Int()))
							}
						}
					} else {
						for {
							select {
							case <-ctx.Done():
								return
							default:
								DoReq(args[0], fmt.Sprintf("%d-%d-%d", uid, time.Now().UnixNano(), rand.Int()))
							}
						}
					}
				}(ctx, i)
			}
			wg.Wait()
			fmt.Println("total cost:", time.Now().Sub(st).String())
		},
	}
	rootCmd.Flags().IntVarP(&concurrent, "concurrent", "c", 1, "并发数")
	rootCmd.Flags().IntVarP(&count, "count", "n", 0, "总请求数")
	rootCmd.Flags().DurationVarP(&period, "period", "t", 0, "总执行时间")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
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
