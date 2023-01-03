# 测试
有三种类型的函数：测试函数、基准测试（benchmark）函数、示例函数。
- 一个测试函数是以Test为函数名前缀的函数，用于测试程序的一些逻辑行为是否正确；go test命令会调用这些测试函数并报告测试结果是PASS或FAIL。
- 基准测试函数是以Benchmark为函数名前缀的函数，它们用于衡量一些函数的性能；
- go test命令会多次运行基准测试函数以计算一个平均的执行时间。示例函数是以Example为函数名前缀的函数，提供一个由编译器保证正确性的示例文档。

## 测试函数
> 见word1
```go
//必须以Test开头，可选的后缀名必须以大写字母开头
func TestSin(t *testing.T) { /* ... */ }
func TestCos(t *testing.T) { /* ... */ }
func TestLog(t *testing.T) { /* ... */ }

```
参数-v可用于打印每个测试函数的名字和运行时间.

参数-run对应一个正则表达式，只有测试函数名被它正确匹配的测试函数才会被go test测试命令运行.

### 随机测试
也就是通过构造更广泛的随机输入来测试探索函数的行为。
- 第一个是编写另一个对照函数，使用简单和清晰的算法，虽然效率较低但是行为和要测试的函数是一致的，然后针对相同的随机输入检查两者的输出结果。
- 第二种是生成的随机输入的数据遵循特定的模式，这样我们就可以知道期望的输出的模式。

如果你使用的是定期运行的自动化测试集成系统，随机测试将特别有价值。

## 测试覆盖率

语句的覆盖率是指在测试中至少被运行一次的代码占总代码数的比例。
> ch7/eval
```shell
go tool cover #显示cover工具的用法

go test -run=Coverage -coverprofile=c.out gopl.io/ch7/eval
```
## benchmark 函数
在Go语言中，基准测试函数和普通测试函数写法类似，但是以Benchmark为前缀名，并且带有一个`*testing.B`类型的参数；*testing.B参数除了提供和*testing.T类似的方法，还有额外一些和性能测量相关的方法。
> 见word2 中BenchmarkIsPalindrome函数
```shell
#默认情况下不运行任何基准测试。我们需要通过-bench命令行标志参数手工指定要运行的基准测试函数。
go test -bench=.
#-benchmem命令行标志参数将在报告中包含内存的分配数据统计。我们可以比较优化前后内存的分配情况
go test -bench=. -benchmem
```

但我们往往想知道的是两个不同的操作的时间对比:如果一个函数需要1ms处理1,000个元素，那么处理10000或1百万将需要多少时间呢?I/O缓存该设置为多大呢?
```go
//todo 这种比较型函数我没有看到具体的调用方式，需要额外资料
func benchmark(b *testing.B, size int) { /* ... */ }
func Benchmark10(b *testing.B)         { benchmark(b, 10) }
func Benchmark100(b *testing.B)        { benchmark(b, 100) }
func Benchmark1000(b *testing.B)       { benchmark(b, 1000) }

```

## 性能剖析
Go语言支持多种类型的剖析性能分析，每一种关注不同的方面，但它们都涉及到每个采样记录的感兴趣的一系列事件消息，每个事件都包含函数调用时函数调用堆栈的信息。内建的go test工具对几种分析方式都提供了支持。

* CPU剖析数据标识了最耗CPU时间的函数。
* 堆剖析则标识了最耗内存的语句。
* 阻塞剖析则记录阻塞goroutine最久的操作。

```shell
# 当同时使用多个标志参数时需要当心，因为一项分析操作可能会影响其他项的分析结果。
$ go test -cpuprofile=cpu.out
$ go test -blockprofile=block.out
$ go test -memprofile=mem.out
```
使用pprof来分析这些数据。它对应go tool pprof命令。该命令有许多特性和选项，但是最基本的是两个参数：生成这个概要文件的可执行程序和对应的剖析数据。
```shell
#-run=NONE 代表禁用哪些简单测试 对于net/http包生成cpu.log
go test -run=NONE -bench=ClientServerParallelTLS64 \
    -cpuprofile=cpu.log net/http
#    -text用于指定输出格式，在这里每行是一个函数，根据使用CPU的时间长短来排序。其中-nodecount=10参数限制了只输出前10行的结果。
go tool pprof -text -nodecount=10 ./http.test cpu.log
```
你可能需要使用pprof的图形显示功能。这个需要安装GraphViz工具，可以从 http://www.graphviz.org 下载。参数-web用于生成函数的有向图，标注有CPU的使用和最热点的函数等信息。

Go官方博客的“Profiling Go Programs”一文。


## Example
godoc这个web文档服务器会将示例函数关联到某个具体函数或包本身，因此ExampleIsPalindrome示例函数将是IsPalindrome函数文档的一部分，Example示例函数将是包文档的一部分。

如果示例函数内含有类似上面例子中的// Output:格式的注释，那么测试工具会执行这个示例函数，然后检查示例函数的标准输出与注释是否匹配。