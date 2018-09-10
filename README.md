# tutorial
Golang程序性能调优教程,性能调优是一个很大的课题,
本文主要介绍golang的逃逸分析机制和一些辅助调优的工具

## 逃逸分析

golang在编译阶段对代码进行分析,将可能被超出函数范围引用的变量对象,分配到堆(heap)内存中,
否则分配到栈(stack)内存中

### 值传递和引用传递

golang中只有slice,map,chan是引用类型,其余都是值类型;
值类型变量在传递时,会产生数据拷贝,而引用类型,可以理解为指针类型,变量传递时仅是指针值拷贝

#### Benchmark

测试相关[`demo`](/heap/value/value_test.go)

```bash
cd heap/value/
go test -bench .
```

得到benchmark结果

```bash
# go test -bench .
...
BenchmarkValue-8        50000000                25.7 ns/op            24 B/op          0 allocs/op
BenchmarkPointer-8      30000000                60.9 ns/op            40 B/op          1 allocs/op
PASS
ok      github.com/little-cui/tutorial/heap/value       6.269s
```

#### GC trace

golang编译的二进制在运行时支持使用环境变量`GODEBUG="gctrace=1"`打开gc报告,
每一次gc都会打印到控制台,包括cpu占用比率,内存变化情况等

测试相关[`demo`](/heap/value/bin/value.go)

- 使用值传递打印结果
```bash
# GODEBUG="gctrace=1" go run heap/value/bin/value.go
...
run TestValue
```

- 使用引用传递打印结果
```bash
# GODEBUG="gctrace=1" go run heap/value/bin/value.go value
...
run TestPointer
gc 2 @0.007s 11%: 0.003+2.3+0.060 ms clock, 0.026+0.50/4.2/0.23+0.48 ms cpu, 8->8->8 MB, 15 MB goal, 8 P
gc 3 @0.020s 8%: 0.003+3.9+0.044 ms clock, 0.026+0.051/7.4/6.0+0.35 ms cpu, 13->14->14 MB, 17 MB goal, 8 P
gc 4 @0.046s 9%: 0.005+14+0.073 ms clock, 0.041+0.057/28/1.3+0.59 ms cpu, 25->29->29 MB, 29 MB goal, 8 P
```

总结下
- 如果业务需要驻留到堆内存的变量在传递时,建议声明为指针类型
- 如果业务中间过程数据,且数据量不大,建议声明为值类型
- 使用指针是有代价的,如:间接寻址,非空判断,堆内存回收

### 分析规则

逃逸分析中,一条基本判定规则: 变量的引用从声明它的函数中被返回

以下是使用golang提供的`gcflags`编译参数,打印出golang编译器是如何判定变量逃逸的

测试相关[`demo`](/heap/escape/bin/escape.go)

```bash
# go build -gcflags '-m -l' heap/escape/bin/escape.go
...
heap/escape/bin/escape.go:22:22: leaking param: f
heap/escape/bin/escape.go:22:7: (*parent).set s does not escape
heap/escape/bin/escape.go:29:12: leaking param: s to result ~r1 level=0
heap/escape/bin/escape.go:38:15: leaking param: s to result ~r1 level=0
heap/escape/bin/escape.go:44:16: Case2 &s does not escape
heap/escape/bin/escape.go:48:9: &s escapes to heap
heap/escape/bin/escape.go:47:15: moved to heap: s
...
```

失败场景:

- 返回slice和map,会分配到堆内存,slice大小超过64K也会分配到堆内存
- slice,map,chan传递对象引用,会分配套堆内存
- interface函数调用,会导致对象和参数都分配到堆内存中

## 性能调优工具

### go tool pprof

golang原生的性能调优工具,可以查看代码的cpu消耗,堆对象和内存的分配情况等,是一个非常强大的工具

##### 加入调试代码

以调优一个使用原生http server实现的服务端为例
```go
import _ "net/http/pprof"
```

##### 启动测试服务端

```bash
go run pprof/bin/bad.go
```

##### 执行`go tool pprof`命令

```bash
# go tool pprof http://127.0.0.1:8080/debug/pprof/profile?seconds=30
...
(pprof) top -cum
Showing nodes accounting for 0.18s, 3.35% of 5.38s total
Dropped 112 nodes (cum <= 0.03s)
Showing top 10 nodes out of 165
      flat  flat%   sum%        cum   cum%
     0.07s  1.30%  1.30%      3.18s 59.11%  runtime.systemstack
         0     0%  1.30%      2.77s 51.49%  net/http.(*conn).serve
         0     0%  1.30%      2.70s 50.19%  main.(*badHandler).ServeHTTP
         0     0%  1.30%      2.70s 50.19%  net/http.(*ServeMux).ServeHTTP
         0     0%  1.30%      2.70s 50.19%  net/http.serverHandler.ServeHTTP
     0.02s  0.37%  1.67%      2.66s 49.44%  main.dumpRequest
         0     0%  1.67%      1.93s 35.87%  runtime.concatstring2
     0.01s  0.19%  1.86%      1.93s 35.87%  runtime.concatstrings
     0.07s  1.30%  3.16%      1.46s 27.14%  runtime.mallocgc
     0.01s  0.19%  3.35%      1.33s 24.72%  runtime.gcBgMarkWorker
```

### go-torch 火焰图

go-torch是在go tool pprof基础上,增加以图形化的方式展示golang程序的性能报告,
开发人员可以查看程序的函数调用栈占比大小,快速定位到性能瓶颈点,从而高效地完成代码优化

![torch](/docs/torch.svg)

##### Requirements

- 安装`go-torch`
```bash
go get github.com/uber/go-torch
```

- 下载`FlameGraph`项目
```bash
git clone https://github.com/brendangregg/FlameGraph
export PATH=$(pwd)/FlameGraph:$PATH
```

##### 执行命令生成svg图片
```bash
# go-torch -u http://127.0.0.1:8080 -t 30
INFO[18:00:38] Run pprof command: go tool pprof -raw -seconds 10 http://127.0.0.1:8080/debug/pprof/profile
INFO[18:00:48] Writing svg to torch.svg
```