# Goroutines和Channels
* 本章讲解goroutine和channel，其支持“顺序通信进程”（communicating sequential processes）或被简称为CSP。CSP是一种现代的并发编程模型.
* 第9章覆盖更为传统的并发模型：多线程共享内存

## Goroutines

> spinner简单的展示了Goroutines的用法
> Clock系列代码展示了一个http链接怎么写，并且展示协程在网络中的应用。 
> Echo系列好像没说明语法用

## Channels
channels则是它们之间的通信机制。一个channel是一个通信机制，它可以让一个goroutine通过它给另一个goroutine发送值信息。

```go
//声明
ch := make(chan int) // ch has type 'chan int'
ch = make(chan int, 3) // 带缓存的 3
//使用
ch <- x  // 发送数据给通道
x = <-ch // 接收值并分配给一个变量
<-ch     // 接收语句但是抛弃了了

//关闭
close(ch)

```

当一个channel被关闭后，再向该channel发送数据将导致panic异常。

### 检测通道关闭
没有办法直接测试一个channel是否被关闭，但是接收操作有一个变体形式：它多接收一个结果，多接收的第二个结果是一个布尔值ok，ture表示成功从channels接收到值，false表示channels已经被关闭并且里面没有值可接收
```go
// Squarer
go func() {
    for {
        x, ok := <-naturals
        if !ok {
            break // channel was closed and drained
        }
        squares <- x * x
    }
    close(squares)
}()

```
Go语言的range循环可直接在channels上面迭代。使用range循环是上面处理模式的简洁语法，它依次从channel接收数据，当channel被关闭并且没有值可接收时跳出循环

```go
    go func() {
        for x := range naturals {
            squares <- x * x
        }
        close(squares)
    }()
```

### 单方向channel
> pipeline3提供了单方面channel的接收

### 带缓存的chan
介绍带缓存的chan以及查询它的容量.以及如何获取channel内部缓存队列中有效元素的个数。
```go
ch = make(chan string, 3)
fmt.Println(cap(ch)) // "3"
//我觉得这个用的少，因为这个有效元素个数随着多进程的进行，时时刻刻而变化
fmt.Println(len(ch)) // "2"

```

## 并发循环
> thumbnail_test提供了许多有趣的例子，这些例子说明了，一个易于并发的任务如何执行，然后还有如果其中任务组出错时如何做，以及WaitGroup的用法
> 注意阻塞引起的协程泄露问题。

## 基于select的多路循环
> 具体见countdown系列

一个没有任何case的select语句写作select{}，会永远地等待下去。

额外补充：nil的channel有时候也是有一些用处的，对一个nil的channel发送和接收操作会永远阻塞，
在select语句中会使case子句永远禁用，可以用nil来激活或者禁用case，来达成处理其它输入或输出事件时超时和取消的逻辑。


* du3将channel当做信号量来控制并发，巧妙
## 并发的退出
* 退出单个协程：可以构造一个channel,让想要退出的协程检查这个channel，然后得到一个信号量之后主动退出。
* 退出多个协程：不考虑构造带buffer的channel,除非你确保有足够的协程消费channel中的信息，否则发送消息的channel协程将会阻塞。
* 退出协程：通过关闭通道来作为一个信号，具体见du4，如何具体实现需要分析自己的程序开启的协程，然后根据协程逻辑设计关闭方案。

## 聊天案例
>具体见chat


