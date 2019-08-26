package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-22 13:55
// @File   : auto.go
// @Project: curl-go/entry
// ==========================

func auto() *cobra.Command {

	// 请求相关参数
	var timeout time.Duration
	var method, protocol, body, outputFile string
	var headers []string

	var correctRate float32
	var p90, p99, timelimitAuto time.Duration
	var maxConcurrency, minConcurrency int
	autoCmd := &cobra.Command{
		Use:   "auto",
		Short: "run curl_go in auto mode",
		Long:  "run curl_go in auto mode",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("还在思考策略哈~")
		},
	}

	// auto 控制
	autoCmd.Flags().DurationVarP(&timelimitAuto, "timelimit", "t", 0, "总执行时间")
	autoCmd.Flags().DurationVar(&p90, "p90", 0, "Top Percentile 90, 即 90% 的请求响应时间小于 p90")
	autoCmd.Flags().DurationVar(&p99, "p99", 0, "Top Percentile 99, 即 99% 的请求响应时间小于 p99")
	autoCmd.Flags().Float32VarP(&correctRate, "correct_rate", "r", 1, "请求的正确率, 有效区间[0-1]")
	autoCmd.Flags().IntVar(&maxConcurrency, "maxc", 0, "最大并发数")
	autoCmd.Flags().IntVar(&minConcurrency, "minc", 1, "最小并发数")

	// auto 请求
	autoCmd.Flags().DurationVarP(&timeout, "timeout", "s", 30*time.Second, "单次请求超时时间")
	autoCmd.Flags().StringVarP(&protocol, "protocol", "p", "HTTP", "协议")
	autoCmd.Flags().StringVarP(&method, "method", "m", "GET", "请求方法[-p=HTTP]")
	autoCmd.Flags().StringArrayVarP(&headers, "headers", "H", []string{}, "请求头[-p=HTTP], ex: '-H a=b -H c=d'")
	autoCmd.Flags().StringVarP(&body, "data", "d", "", "请求体[-p=HTTP]")
	autoCmd.Flags().StringVarP(&outputFile, "output", "o", "-", "输出结果到文件, '-' 表示标准输出[-p=HTTP]")
	return autoCmd
}
