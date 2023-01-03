# 接口

接口类型是对其它类型行为的抽象和概括；因为接口类型不会和特定的实现细节绑定在一起，通过这种抽象的方式我们可以让我们的函数更加灵活和更具有适应能力。

Go语言中接口类型的独特之处在于它是满足隐式实现的。也就是说，我们没有必要对于给定的具体类型定义所有满足的接口类型；简单地拥有一些必需的方法就足够了。

## 接口是合约

目前为止，我们看到的类型都是具体的类型，当你拿到一个具体的类型时你就知道它的本身是什么和你可以用它来做什么。

接口类型。接口类型是一种抽象的类型。它不会暴露出它所代表的对象的内部值的结构和这个对象支持的基础操作的集合；它们只会表现出它们自己的方法。也就是说当你有看到一个接口类型的值时，你不知道它是什么，唯一知道的就是可以通过它的方法来做什么。

接口举例：
> bytecounter

```go
package fmt

//这个Fprintf第一个参数是io.Writer类型，而它是一个接口类型，规定了有一个Write方法

func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
func Printf(format string, args ...interface{}) (int, error) {
	return Fprintf(os.Stdout, format, args...)
}
func Sprintf(format string, args ...interface{}) string {
	var buf bytes.Buffer
	Fprintf(&buf, format, args...)
	return buf.String()
}

```

## 接口类型
接口声明：方法顺序的变化也没有影响，唯一重要的就是这个集合里面的方法。

```go
//
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Closer interface {
    Close() error
}

//接口内嵌
type ReadWriteCloser interface {
	Reader
	Closer
}
//混合模式
type ReadWriter interface {
	Read(p []byte) (n int, err error)
	Closer
}

```

## 实现接口的条件
一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口。例如*os.File类型实现了io.Reader，Writer，Closer，和ReadWriter接口

### 接口申明
```go
//可以申明接口类型
var w io.Writer
w = os.Stdout           // OK: *os.File has Write method
w = new(bytes.Buffer)   // OK: *bytes.Buffer has Write method
w = time.Second         // compile error: time.Duration lacks Write method

var rwc io.ReadWriteCloser
rwc = os.Stdout         // OK: *os.File has Read, Write, Close methods
rwc = new(bytes.Buffer) // compile error: *bytes.Buffer lacks Close method
//不同的接口类型的变量具有兼容性
w = rwc                 // OK: io.ReadWriteCloser has Write method
rwc = w                 // compile error: io.Writer lacks Close method

```

注意一个细节有关于方法上绑定是指针还是本身的方法签名区别
```go
type IntSet struct { /* ... */ }
//这个方法代表*IntSet类型实现了fmt.String接口,并不意味着IntSet实现了fmt.String接口
func (*IntSet) String() string
var _ = IntSet{}.String() // compile error: String requires *IntSet receiver

var _ fmt.Stringer = &s // OK
var _ fmt.Stringer = s  // compile error: IntSet lacks String method

```
如果声明为接口类型，不能调用除接口中声明的方法以外的方法
```go
os.Stdout.Write([]byte("hello")) // OK: *os.File has Write method
os.Stdout.Close()                // OK: *os.File has Close method

var w io.Writer
// 即使os.Stdout底层存在Close方法，也不能调用
w = os.Stdout
w.Write([]byte("hello")) // OK: io.Writer has Write method
w.Close()                // compile error: io.Writer lacks Close method
```

### inteface{}调用

一个interface{}值持有一个boolean，float，string，map，pointer，或者任意其它的类型；我们当然不能直接对它持有的值做操作，因为interface{}没有任何方法。，除了使用类型断言。

## flag.Value接口

> tempflag 阐述了flag.Value的用法

## 接口值
一个接口的零值就是它的类型和值的部分都是nil，通过使用w==nil或者w!=nil来判断接口值是否为空。调用一个空接口值上的任意方法都会产生panic

参考下面的程序：
```go
//此时 w 变量的 type是nil，value是nil
var w io.Writer
//此时 w 变量的 type是*os.File，value是指向变量的地址
w = os.Stdout
//此时 w 变量的 type是*bytes.Buffer，value是指向变量的地址
w = new(bytes.Buffer)
//此时 w 变量的 type是nil，value是nil
w = nil
```

### 接口值得布尔操作

两个接口值相等仅当它们都是nil值，或者它们的动态类型相同并且动态值（如果动态类型不能比较将会引发panic）也根据这个动态类型的==操作相等。所以它们可以用在map的键或者作为switch语句的操作数。

使用%T获取底层类型。

### 注意点。
注意，在调用接口类型时需要小心。如下面的方法。
```go

// If out is non-nil, output will be written to it.
//这是一个参数接收接口类型的方法。
func f(out io.Writer) {
    // ...do something...
    if out != nil {
        out.Write([]byte("done!\n"))
    }
}

//如果调用方法为具体的类如下，例子中为*bytes.Buffer.类型为*bytes.Buffer,值为nil，传入f中。`out != nil`这个判断语句将会为false
var buf *bytes.Buffer
if debug {
buf = new(bytes.Buffer) // enable collection of output
}
f(buf)

//正确的是方法是，如果是接口类型的参数，就提前声明接口类型，然后中途给他赋值
var buf io.Writer
if debug {
buf = new(bytes.Buffer) // enable collection of output
}
f(buf) // OK
```

## sort.Interface接口如何使用

> 见sorting代码，其中可以看到如果想要根据不同的方式排序的做法（封装为一个包装器），
> sort.Reverse的代码的实现思路是将Sort.Interface接口包装，然后定义自己的reverse方法，然后暴露Reverse方法
> 比起这章展示的sort接口如何调用，更重要的是如何将变量和不变量分离的代码组织做法。这是一种能力

sort.Interface的定义
```go
package sort

type Interface interface {
    Len() int
    Less(i, j int) bool // i, j are indices of sequence elements
    Swap(i, j int)
}

//然后使用sort.Sort()方法来进行排序
```

## http.Handle接口

> http系列代码讲了如何写一个http程序，然后http4是最简练版本。第3版本讲了一些语法知识。（主要是func值的运用）

主要了解net/http包去实现网络客户端和服务器

注意：web服务器在一个新的协程中调用每一个handler（handle是在协程中处理），所以当handler获取其它协程或者这个handler本身的其它请求也可以访问到变量时，一定要使用预防措施，比如锁机制。

## error接口

### error定义
```go
type error interface {
    Error() string
}

```

### error的创建
```go

// 第一种创建error最简单方法
//创建一个error最简单的方法就是调用errors包的New函数 errors.New。内部存储了一个string类型字符串
package errors

//每个New函数的调用都分配了一个独特的和其他错误不相同的实例，因为指针的缘故
func New(text string) error { return &errorString{text} }
//创建了一个字符串
type errorString struct { text string }
// 注意是*errorString类型实现了error借口，而不是errorString类型实现了error借口
//实现接口定义
func (e *errorString) Error() string { return e.text }

//这使得下面的代码返回false
fmt.Println(errors.New("EOF") == errors.New("EOF")) // "false"

//第二种创建error的方法，常用
//使用fmt的Errorf方法，它会处理字符串格式化。
package fmt
func Errorf(format string, args ...interface{}) error {
	return errors.New(Sprintf(format, args...))
}
```

### 其他error类型
syscall包提供了Go语言底层系统调用API。在多个平台上，它定义一个实现error接口的数字类型Errno，并且在Unix平台上，Errno的Error方法会从一个字符串表中查找错误消息
````go
//syscall中的定义
package syscall

type Errno uintptr // operating system error code
var errors = [...]string{
    1:   "operation not permitted",   // EPERM
    2:   "no such file or directory", // ENOENT
    3:   "no such process",           // ESRCH
    // ...
}

func (e Errno) Error() string {
if 0 <= int(e) && int(e) < len(errors) {
return errors[e]
}
return fmt.Sprintf("errno %d", e)
}

//调用
var err error = syscall.Errno(2)
fmt.Println(err.Error()) // "no such file or directory"
fmt.Println(err)         // "no such file or directory"

````

## 类型断言
如果断言的类型T是一个具体类型，具体类型的类型断言从它的操作对象中获得具体的值。如果检查失败，接下来这个操作会抛出panic
```go
var w io.Writer
w = os.Stdout
f := w.(*os.File)      // success: f == os.Stdout
c := w.(*bytes.Buffer) // panic: interface holds *os.File, not *bytes.Buffer

```

如果相反地断言的类型T是一个接口类型，然后类型断言检查是否x的动态类型满足T。如果这个检查成功了,改变了可以获取的方法集合（通常更大），但是它保留了接口值内部的动态类型和值的部分。
```go
var w io.Writer
w = os.Stdout
// 成功则获得括号中说明的类型
rw := w.(io.ReadWriter) 
w = new(ByteCounter)
// 失败则抛出错误
rw = w.(io.ReadWriter)
```
你有时会看见原来的变量名重用而不是声明一个新的本地变量名，这个重用的变量原来的值会被覆盖（理解：其实是声明了一个同名的新的本地变量，外层原来的w不会被改变），如下面这样
```go
if w, ok := w.(*os.File); ok {
    // ...use w...
}
```
### 类型断言的应用
* 使用类型断言识别错误类型
```go
import (
    "errors"
    "syscall"
)

var ErrNotExist = errors.New("file does not exist")

// IsNotExist returns a boolean indicating whether the error is known to
// report that a file or directory does not exist. It is satisfied by
// ErrNotExist as well as some syscall errors.
func IsNotExist(err error) bool {
	//使用类型断言识别错误类型，前面说过错误类型一般是指针类型，不可比较，只能和nil比较
    if pe, ok := err.(*PathError); ok {
        err = pe.Err
    }
    return err == syscall.ENOENT || err == ErrNotExist
}

```
* 通过类型断言询问行为

```go
//假设一个糟糕的程序是这样。猜想w变量持有的动态类型也有一个允许字符串高效写入的WriteString方法；这个方法会避免去分配一个临时的拷贝。

func writeHeader(w io.Writer, contentType string) error {
	//
    if _, err := w.Write([]byte("Content-Type: ")); err != nil {
    	//。这个转换分配内存并且做一个拷贝，但是这个拷贝在转换后几乎立马就被丢弃掉。让我们假装这是一个web服务器的核心部分并且我们的性能分析表示这个内存分配使服务器的速度变慢。
        return err
    }
    if _, err := w.Write([]byte(contentType)); err != nil {
        return err
    }
// ...
}

```

## 类型分支 switch
* 使用字面量type
```go
switch x.(type) {
    case nil:       // ...
    case int, uint: // ...
    case bool:      // ...
    case string:    // ...
    default:        // ...
}

```
* 拓展形式
> xmlselect 展示了具体做法
```go
//类型分支语句有一个扩展的形式，它可以将提取的值绑定到一个在每个case范围内都有效的新变量。
func sqlQuote(x interface{}) string {
    switch x := x.(type) { //这个地方获取值
        case nil:
            return "NULL"
        case int, uint:
            return fmt.Sprintf("%d", x) // x has type interface{} here.
        case bool:
            if x {
                return "TRUE"
            }
            return "FALSE"
        case string:
            return sqlQuoteString(x) // (not shown)
        default:
            panic(fmt.Sprintf("unexpected type %T: %v", x, x))
    }
}

```
## 一些注意
* 新手Go程序员总是先创建一套接口，然后再定义一些满足它们的具体类型。这种方式的结果就是有很多的接口，它们中的每一个仅只有一个实现。不要再这么做了。这种接口是不必要的抽象；
