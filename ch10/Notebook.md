# 包简介

1. 所有导入的包必须在每个文件的开头显式声明，这样的话编译器就没有必要读取和分析整个源文件来判断包的依赖关系。
2. 禁止包的环状依赖，因为没有循环依赖，包的依赖关系形成一个有向无环图，每个包可以被独立编译，而且很可能是被并发编译。
3. 编译后包的目标文件不仅仅记录包本身的导出信息，目标文件同时还记录了包的依赖关系。因此，在编译一个包的时候，编译器只需要读取每个直接导入包的目标文件，而不需要遍历所有依赖的的文件.

## 包路径
每个包是由一个全局唯一的字符串所标识的导入路径定位。Go语言的规范并没有指明包的导入路径字符串的具体含义，导入路径的具体含义是由构建工具来解释的。
Go语言的规范并没有指明包的导入路径字符串的具体含义，导入路径的具体含义是由构建工具来解释的。

## 包声明

包声明语句的主要目的是确定当前包被其它包导入时默认的标识符。
第一个例外，包对应一个可执行程序，也就是main包，这时候main包本身的导入路径是无关紧要的。名字为main的包是给go build构建命令一个信息，这个包编译完之后必须调用连接器生成一个可执行程序。
所有以_test为后缀包名的测试外部扩展包都由go test命令独立编译，普通包和测试的外部扩展包是相互独立的。
一些依赖版本号的管理工具会在导入路径后追加版本号信息，例如“gopkg.in/yaml.v2”。这种情况下包的名字并不包含版本号后缀，而是yaml

## 导入声明

- 导入的包之间可以通过添加空行来分组；通常将来自不同组织的包独自分组。包的导入顺序无关紧要，但是在每个分组中一般会根据字符串顺序排列。
- 如果我们想同时导入两个有着名字相同的包，例如math/rand包和crypto/rand包，那么导入声明必须至少为一个同名包指定一个新的包名以避免冲突。这叫做导入包的重命名。
```go
import (
    "crypto/rand"
    mrand "math/rand" // alternative name mrand avoids conflict
)
```

## 包的匿名导入
如果只是导入一个包而并不使用导入的包将会导致一个编译错误。但是有时候我们只是想利用导入包而产生的副作用：它会计算包级变量的初始化表达式和执行导入包的init初始化函数.
```go
import _ "image/png" // register PNG decoder
```

## 包的命名

给出了些包的命名意见，但是我todo

## 工具

- go get可以下载一个单一的包或者用...下载整个子目录里面的每个包。
- go get命令获取的代码是真实的本地存储仓库。
- 可以让包用一个自定义的导入路径，但是真实的代码却是由更通用的服务提供。
```shell
$ go build gopl.io/ch1/fetch
# 这个go-import才是关键
$ ./fetch https://golang.org/x/net/html | grep go-import
<meta name="go-import"
      content="golang.org/x/net git https://go.googlesource.com/net">

```
- 如果指定-u命令行标志参数，go get命令将确保所有的包和依赖的包的版本都是最新的，然后重新编译和安装它们。

### go build
* 如果包是一个库，则忽略输出结果；这可以用于检测包是可以正确编译的。如果包的名字是main，go build将调用链接器在当前目录创建一个可执行程序；以导入路径的最后一段作为可执行程序的名字。
* 每个对应可执行程序或者叫Unix术语中的命令的包，会要求放到一个独立的目录中。这些目录有时候会放在名叫cmd目录的子目录下面.
* 用一个相对目录的路径名指定，相对路径必须以.或..开头。如果没有指定参数，那么默认指定为当前目录对应的包。
* 可以指定包的源文件列表，这一般只用于构建一些小程序或做一些临时性的实验。如果是main包，将会以第一个Go源文件的基础文件名作为最终的可执行程序的名字。或者使用go run *.go代替。
* go install会保存每个包的编译成果，而不是将它们都丢弃。被编译的包会被保存到$GOPATH/pkg目录下，目录路径和 src目录路径对应，可执行程序被保存到$GOPATH/bin目录。
* 针对不同操作系统或CPU的交叉构建也是很简单的。只需要设置好目标对应的GOOS和GOARCH，然后运行构建命令即可。
```go
//下面以64位和32位环境分别编译和执行
$ go build gopl.io/ch10/cross
$ ./cross
darwin amd64
$ GOARCH=386 go build gopl.io/ch10/cross
$ ./cross
darwin 386
```

有些包可能需要针对不同平台和处理器类型使用不同版本的代码文件，以便于处理底层的可移植性问题或为一些特定代码提供优化。如果一个文件名包含了一个操作系统或处理器类型名字，例如net_linux.go或asm_amd64.s，Go语言的构建工具将只在对应的平台编译这些文件。还有一个特别的构建注释参数可以提供更多的构建过程控制
```go
//只在go build只在编译程序对应的目标操作系统是Linux或Mac OS X时才编译这个文件
// +build linux darwin
//下面的构建注释则表示不编译这个文件
// +build ignore
```
更多参考：go doc go/build

### 包文档
* go doc命令，该命令打印其后所指定的实体的声明与文档注释，该实体可能是一个包：
```shell
#包级别
$ go doc time
package time // import "time"

Package time provides functionality for measuring and displaying time.

const Nanosecond Duration = 1 ...
func After(d Duration) <-chan Time
func Sleep(d Duration)
func Since(t Time) Duration
func Now() Time
type Duration int64
type Time struct { ... }
...many more...

# 某个具体的包成员
$ go doc time.Since
func Since(t Time) Duration

    Since returns the time elapsed since t.
    It is shorthand for time.Now().Sub(t).

#对于方法的文档
$ go doc time.Duration.Seconds
func (d Duration) Seconds() float64

    Seconds returns the duration as a floating-point number of seconds.

```
* 第二个工具，名字也叫godoc，它提供可以相互交叉引用的HTML页面，[godoc的在线服务](https://godoc.org) ，包含了成千上万的开源包的检索工具。
```shell
#在浏览器查看 http://localhost:8000/pkg
$ godoc -http :8000
```
### 内部包
Go语言的构建工具对包含internal名字的路径段的包导入路径做了特殊处理。这种包叫internal包，一个internal包只能被和internal目录有同一个父目录的包所导入。
例如，net/http/internal/chunked内部包只能被net/http/httputil或net/http包导入

规则：导出路径包含internal关键字的包，只允许internal的父级目录及父级目录的子包导入，其它包无法导入。

### 查询包
```shell
#go list命令可以查询可用包的信息
go list github.com/go-sql-driver/mysql

#go list命令的参数还可以用"..."表示匹配任意的包的导入路径。
go list ...

#或者是特定子目录下的所有包
go list gopl.io/ch3/...

# 和某个主题相关的所有包
go list ...xml...
```