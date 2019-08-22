# curl-go
压测小工具

## just use it
    curl_go 是 macos 版本的二进制
    curl_go_linux 是 linux 版本的二进制. 
    
使用方法

    curl_go -h
    
开始尝试吧. 

### example
示例 1: 当成 curl 用(一般可以用于压测前的请求验证)

    curl_ab -m GET -H Authorization=zero https://www.baidu.com
    
示例 2: 并发 100, 请求 1000 次, 记录对端 ip:port, 输出请求结果到文件: /tmp/zero96.cc, 并且统计 p90, p99. 

    curl_ab -c 100 -n 1000 --record_ip --tp=90 --tp=99 -o /tmp/zero96.cc https://www.baidu.com
    
示例 3: 结果太多, 使用蓄水池算法来抽样(只保留 1000 个结果, 来计算 p90, p99 的值)

    curl_ab -c 1000 -n 100000 --reservoir=1000 --record_ip --tp=90 --tp=99 -o /tmp/baidu.com.csv https://www.baidu.com
    
