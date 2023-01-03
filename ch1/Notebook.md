# 第一章 入门

## hello world
主要介绍go的环境搭建，go的程序的编译，运行指令，以及一些常用的工具命令，比如goget,gofmt。

### go环境搭建。
* 先去官网下载go安装包。然后安装。
* 设置go中的环境变量,如 GOPROXY(依赖下载的镜像)，GOPATH(go的工作目录)，可以参考这个[链接](https://www.xampp.cc/archives/22465) 
* 设置完成之后可以使用 `go env` 命令打印查看设置的环境变量

### go程序的编译，运行指令
1. 编写代码并保存为helloworld.go。如下
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, 世界")
}
```
2.运行命令
```shell
go run helloworld.go
```
3.编译程序,这会在当前文件夹下生成一个各执行文件。
```shell
go build helloworld.go
```

4.获取网络上的第三方代码。

执行如下代码，需要有git等项目管理工具
```shell
go get gopl.io/ch1/helloworld
```
执行成功之后会放在`$GOPATH/src/gopl.io/ch1/helloworld`目录下

### 额外的知识点

1. go程序的方法声明和别的程序不太一样。go的方法签名按照这个顺序来： func 关键字、函数名、参数列表、`返回值列表` 注意这个返回值列表的顺序.

2. Go 语言不需要在语句或者声明的末尾添加分号，除非一行上有多条语句。

3. Go 语言在代码格式上采取了很强硬的态度。gofmt工具把代码格式化为标准格式.

4.goimports，可以根据代码需要，自动地添加或删除 import 声明。这个工具并没有包含在标准的分发包中，可以用下面的命令安装：
```shell
 go get golang.org/x/tools/cmd/goimports
```

## 命令行参数 

> 讲了一些小trick 比如连接字符串，循环，获取命令行参数，以及变量声明，具体见echo系列代码

通常来说，输入来自于程序外部：文件、网络连接、其它程序的输出、敲键盘的用户、命令行参数或其它类似输入源。

如何获取命令行参数：程序的命令行参数可从 os 包的 Args 变量获取；os 包外部使用 os.Args 访问该变量。os.Args[0]，是命令本身的名字；其它的元素则是程序启动时传给它的参数。

参考echo1获取关于echo直观感受，同时你也可以看到一些有关注释,`局部变量`,`for循环的两种用法`，以及`变量声明的小技巧`等内容.

有关变量声明的总结如下所示：

```go
//第一个是一条短变量声明，最简洁，但只能用在函数内部，而不能用于包变量
s := ""
//第二种形式依赖于字符串的默认初始化零值机制
var s string
var s = ""
//第四种形式显式地标明变量的类型，当变量类型与初值类型相同时，类型冗余，但如果两者类型不同，变量类型就必须了
var s string = ""
```

字符串拼接消耗巨大，使用strings.Join来减少开销。

## 查找重复的行
> 主要的内容包含：文件的读取以及从输入端获取数据、fmt.Printf的格式化、向错误流输入信息，以及make创建了一个引用变量使得可以在函数之间传递引用、ioutil.ReadFile整个读入文件然后一口气统计。
> 详情见dup系列代码

Printf常用的格式占位符,默认不会换行。

| ---------------- | ---------------- |  
| %d | 十进制整数 |
| %x, %o, %b | 十六进制，八进制，二进制整数。 |
| %f, %g, %e | 浮点数： 3.141593 3.141592653589793 3.141593e+00 |
| %t | 布尔：true或false |
| %c | 字符（rune） (Unicode码点) |
| %d | 十进制整数 |
| %s | 字符串 |
| %q | 带双引号的字符串"abc"或带单引号的字符'c' |
| %v | 变量的自然形式（natural format） |
| %T | 变量的类型 |
| %% | 字面上的百分号标志（无操作数） |

## GIF动画
> 用了一些新的结构，包括const声明，struct结构体类型，复合声明。以及一个如何画gif图片的演示。详情见lissajous代码

当我们import了一个包路径包含有多个单词的package时，通常我们只需要用最后那个单词表示这个包就可以。

常量有包型和函数型，目前常量声明的值必须是一个数字值、字符串或者一个固定的boolean值。

struct内部的变量可以以一个点（.）来进行访问

## 获取URL

> 主要介绍了如何通过net/http包获取一个网址，并且读取内容。代码见fetch代码。

## 并发获取多个URL

> 简单的体验Go中的goroutine和channel，详情见代码并发获取多个URL。

go 方法名 启动一个goroutine  

ch相当于一个缓冲区

## Web服务
> 主要简单的展示了一个简单的Web服务如何写。代码见server系列。server2 展示了互斥量的访问。server3展示了如何打印http中的信息，以及if语句和其他语言的不同

## 小结以及额外补充：

- 多行注释：和其他语言一样采用 /* ... */ 来包裹。

- 包：在源文件的开头写的注释是这个源文件的文档。在每一个函数之前写一个说明函数行为的注释也是一个好习惯，因为这些内容会被像godoc这样的工具检测到

- 如何查看注释：
```shell
go doc http.ListenAndServe
```

- switch 控制流：Go语言并不需要显式地在每一个case后写break，语言默认执行完case后的逻辑语句会自动退出。当然了，如果你想要相邻的几个case都执行同一逻辑的话，需要自己显式地写上一个fallthrough语句来覆盖这种默认行为。
```go
switch coinflip() {
case "heads":
    heads++
case "tails":
    tails++
default:
    fmt.Println("landed on edge!")
}
```
- 此外，Go语言里的switch还可以不带操作对象，然后将每个case的表达式和true值进行比较
```go
func Signum(x int) int {
    switch {
    case x > 0:
        return +1
    default:
        return 0
    case x < 0:
        return -1
    }
}
```
- 想跳过的是更外层的循环的话，我们可以在相应的位置加上label，这样break和continue就可以根据我们的想法来continue和break任意循环
- type 关键字
```go
type Point struct {
    X, Y int
}
var p Point
```


- 指针：Go语言在这两种范围中取了一种平衡。指针是可见的内存地址，&操作符可以返回一个变量的内存地址，并且*操作符可以获取指针指向的变量内容，但是在Go语言里没有指针运算，也就是`不能像c语言里可以对指针进行加或减操作`