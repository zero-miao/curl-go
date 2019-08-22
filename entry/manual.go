package entry

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zero-miao/curl-go/req"
	"github.com/zero-miao/curl-go/run"
	"github.com/zero-miao/curl-go/stat"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-22 13:54
// @File   : manual.go
// @Project: curl-go/entry
// ==========================

func manual() *cobra.Command {
	// 请求参数
	rc := RequestCLIConfig{}
	// 并发控制
	mc := stat.ManualConfig{}

	reportConfig := ReportCLIConfig{}

	var request interface{}
	var runner req.Runner
	var reporter run.Reporter
	manualCmd := &cobra.Command{
		Use: "curl_ab <url>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("url required")
			}
			// 构造请求.
			rc.Url = args[0]
			if temp, err := rc.Make(); err != nil {
				return err
			} else {
				request = temp
			}
			if temp, err := rc.MakeRunner(); err != nil {
				return err
			} else {
				runner = temp
			}
			if temp, err := reportConfig.Make(); err != nil {
				return err
			} else {
				reporter = temp
			}
			return nil
		},
		Short: "curl_ab just like curl+ab",
		Long:  "curl_ab just like curl+ab",
		Run: func(cmd *cobra.Command, args []string) {
			mode := &stat.ManualMode{
				ModeArgs: mc,

				Runner: runner,
				Requests: []interface{}{
					request,
				},
				Reporter: reporter,
			}
			mode.Run()
		},
	}

	// 请求相关
	manualCmd.Flags().DurationVarP(&rc.Timeout, "timeout", "s", 30*time.Second, "单次请求超时时间")
	manualCmd.Flags().StringVar(&rc.Protocol, "protocol", "HTTP", "协议")
	manualCmd.Flags().StringVarP(&rc.Method, "method", "m", "GET", "请求方法[-p=HTTP]")
	manualCmd.Flags().StringArrayVarP(&rc.Headers, "header", "H", []string{}, "请求头[-p=HTTP], ex: '-H a=b -H c=d'")
	manualCmd.Flags().StringVarP(&rc.Body, "data", "d", "", "请求体[-p=HTTP]")
	manualCmd.Flags().IntVar(&rc.CorrectSMax, "c_max", 499, "大于 c_max 均视为不正确请求")
	manualCmd.Flags().IntVar(&rc.WrongSMin, "w_min", 0, "小于 w_min 均视为正确请求, 晚于 c_max 判断")
	manualCmd.Flags().BoolVar(&rc.RecordRemoteIp, "record_ip", false, "是否记录服务器 ip")

	manualCmd.Flags().IntVarP(&mc.Concurrency, "concurrency", "c", 1, "并发数")
	manualCmd.Flags().IntVarP(&mc.Count, "requests", "n", 0, "总请求数")
	manualCmd.Flags().DurationVarP(&mc.TimeLimit, "timelimit", "t", 0, "总执行时间")

	manualCmd.Flags().StringVarP(&reportConfig.Reporter, "reporter", "r", "default", "接收请求结果(选项: default)")
	manualCmd.Flags().StringVarP(&reportConfig.ReportFile, "output", "o", "-", "输出结果到文件, '-' 表示标准输出")
	manualCmd.Flags().StringArrayVarP(&reportConfig.Filters, "filter", "f", []string{}, "过滤用于计算 TP 的结果集(选项: correct, error)")
	manualCmd.Flags().IntVar(&reportConfig.Reservoir, "reservoir", 0, "蓄水池算法抽样长度, 只计算 filter 之后的结果")
	manualCmd.Flags().IntSliceVar(&reportConfig.TP, "tp", []int{90}, "对最后的结果集计算: Top Percentile (0, 100), 开区间")
	return manualCmd
}
