package run

import (
	"fmt"
	"github.com/zero-miao/curl-go/req"
	"io"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-21 12:28
// @File   : data.go
// @Project: curl-go/stat
// ==========================

func partition(du []time.Duration, start, end int) int {
	if start >= end {
		return -1
	}
	pivot := du[start]
	var j = start // 最后一个小于等于 pivot 的值.
	for i := start + 1; i < end; i++ {
		if du[i] <= pivot {
			j++
			du[i], du[j] = du[j], du[i]
		}
	}
	du[start], du[j] = du[j], du[start]

	partition(du, start, j)
	partition(du, j+1, end)
	return j
}

type Statistic struct {
	WriterCSV       io.WriteCloser
	Filter          func(v1 *req.ReportV1) bool // 抽样方法
	SampleReservoir int
	ReservoirIndex  int
	TP              []int

	// 实际的运行数据, 并发多少, 总请求多少, 总运行时间多少.
	result interface{}

	mu sync.Mutex
	// 成功的请求数.
	success int
	// 失败的请求数.
	failed int
	// 已排序的响应时间列表, 每一次运行都能报告一次
	lats []time.Duration
	// 统计值
	topPercentile map[string]time.Duration
}

func (s *Statistic) String() string {
	tp := make([]string, 0, len(s.topPercentile))
	for k, v := range s.topPercentile {
		tp = append(tp, fmt.Sprintf("%s=%s", k, v.String()))
	}
	sort.Strings(tp)
	return fmt.Sprintf("请求信息: %s; 成功率: %d/%d=%d%%\n\n样本统计(length=%d): %s", s.result, s.success, s.success+s.failed, s.success*100/(s.success+s.failed), len(s.lats), strings.Join(tp, ", "))
}

func (s *Statistic) StreamReport(report *req.ReportV1) {
	if s.WriterCSV != nil {
		s.WriterCSV.Write([]byte(report.CSV()))
	}
	ok := true
	if s.Filter != nil {
		ok = s.Filter(report)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if report.IsCorrect {
		s.success++
	} else {
		s.failed++
	}
	if ok {
		if s.SampleReservoir > 0 {
			// 流式数据中抽样 s.SampleReservoir 个.
			if s.ReservoirIndex < s.SampleReservoir {
				s.lats = append(s.lats, report.Latency)
			} else {
				// 从第 reservoirIndex+1 个数开始, 以 SampleReservoir/reservoirIndex 的概率选择随机放入一个位置.
				if rand.Intn(s.ReservoirIndex) < s.SampleReservoir {
					pos := rand.Intn(s.SampleReservoir)
					s.lats[pos] = report.Latency
				}
			}
			s.ReservoirIndex++
		} else {
			s.lats = append(s.lats, report.Latency)
		}
	}
}

func (s *Statistic) ResultReport(result interface{}) {
	s.result = result
}

func (s *Statistic) Result(print bool) {
	if print {
		fmt.Println("\n正在统计结果...")
	}
	length := len(s.lats)
	s.topPercentile = map[string]time.Duration{}
	partition(s.lats, 0, length)
	for _, tp := range s.TP {
		index := (tp * length) / 100
		s.topPercentile[fmt.Sprintf("tp%d", tp)] = s.lats[index]
	}
	if print {
		fmt.Println(s)
	}
}
