package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zero-miao/curl-go/cli"
	"github.com/zero-miao/curl-go/mode"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-22 13:54
// @File   : manual.go
// @Project: curl-go/entry
// ==========================

func manual() *cobra.Command {
	var fileConfig, protocol string

	// 请求参数
	reqConfig := &cli.RequestCLIConfig{}
	// 运行逻辑,
	runConfig := &cli.RunnerCLIConfig{}
	// 结果处理
	reportConfig := &cli.ReportCLIConfig{}
	// 并发控制
	mc := &cli.ModeCLIConfig{}

	var modeInterface mode.I
	manualCmd := &cobra.Command{
		Use: "curl_ab [url]",
		Args: func(cmd *cobra.Command, args []string) error {
			if fileConfig != "" {
				data, err := ioutil.ReadFile(fileConfig)
				if err != nil {
					return fmt.Errorf("读取文件失败: %v", err)
				}
				var file = new(cli.FileCLIConfig)
				if err := yaml.Unmarshal(data, &file); err != nil {
					return fmt.Errorf("文件解析异常: %v", err)
				}
				temp, err := file.MakeMode()
				if err != nil {
					return err
				}
				modeInterface = temp
			} else {
				if protocol == "" {
					return fmt.Errorf("no protocol")
				}
				reqConfig.Protocol = protocol
				runConfig.Protocol = protocol
				if reqConfig.Url == "" {
					if len(args) != 1 {
						return fmt.Errorf("url required")
					}
					reqConfig.Url = args[0]
				}
				if reportConfig.Reporter == "" {
					return fmt.Errorf("no reporter")
				}
				mc.Mode = "manual"
				r, err := reqConfig.MakeRequest()
				if err != nil {
					return err
				}
				temp, err := mc.MakeMode([]interface{}{r}, runConfig, reportConfig)
				if err != nil {
					return err
				}
				modeInterface = temp
			}
			return nil
		},
		Short: "curl_ab just like curl+ab",
		Long:  "curl_ab just like curl+ab",
		Run: func(cmd *cobra.Command, args []string) {
			modeInterface.Run()
		},
	}

	manualCmd.Flags().StringVarP(&fileConfig, "file", "f", "", "配置文件")
	// 请求相关
	manualCmd.Flags().StringVarP(&protocol, "protocol", "p", "HTTP", "协议")
	// RequestCLIConfig
	manualCmd.Flags().StringVarP(&reqConfig.Method, "method", "m", "GET", "请求方法[-p=HTTP]")
	manualCmd.Flags().StringArrayVarP(&reqConfig.Headers, "header", "H", []string{}, "请求头[-p=HTTP], ex: '-H a=b -H c=d'")
	manualCmd.Flags().StringVarP(&reqConfig.Body, "data", "d", "", "请求体[-p=HTTP]")
	manualCmd.Flags().StringVar(&reqConfig.Url, "url", "", "请求url[-p=HTTP]")

	// RunnerCLIConfig
	manualCmd.Flags().BoolVar(&runConfig.RecordRemoteIp, "record_ip", false, "是否记录服务器 ip")
	manualCmd.Flags().StringVarP(&runConfig.Timeout, "timeout", "s", "30s", "单次请求超时时间 `duration`")
	manualCmd.Flags().IntVar(&runConfig.CorrectSMax, "c_max", 499, "大于 c_max 均视为不正确请求")
	manualCmd.Flags().IntVar(&runConfig.WrongSMin, "w_min", 0, "小于 w_min 均视为正确请求, 晚于 c_max 判断")

	// manual 控制相关
	manualCmd.Flags().IntVarP(&mc.Concurrency, "concurrency", "c", 1, "并发数")
	manualCmd.Flags().IntVarP(&mc.Count, "requests", "n", 0, "总请求数")
	manualCmd.Flags().StringVarP(&mc.TimeLimit, "timelimit", "t", "0s", "总执行时间 `duration`")

	// 输出相关
	manualCmd.Flags().StringVarP(&reportConfig.Reporter, "reporter", "r", "default", "接收请求结果(选项: default)")
	manualCmd.Flags().StringVar(&reportConfig.ResultCSV, "csv", "-", "输出结果到文件, '-' 表示标准输出")
	manualCmd.Flags().StringArrayVarP(&reportConfig.Filters, "filter", "F", []string{}, "过滤用于计算 TP 的结果集(选项: correct, error)")
	manualCmd.Flags().IntVarP(&reportConfig.Reservoir, "reservoir", "R", 1000, "蓄水池算法抽样长度, 只计算 filter 之后的结果")
	manualCmd.Flags().IntSliceVarP(&reportConfig.TP, "tp", "T", []int{90}, "对最后的结果集计算: Top Percentile (0, 100), 开区间")
	return manualCmd
}
