# Unsafe

## unsafe.Sizeof, Alignof 和 Offsetof
这三个都是为了让程序员理解内存布局而存在。

unsafe.Sizeof函数返回操作数在内存中的字节大小.
```go
import "unsafe"
fmt.Println(unsafe.Sizeof(float64(0))) // "8"
```
unsafe.Alignof 函数返回对应参数的类型需要对齐的倍数。

unsafe.Offsetof 函数的参数必须是一个字段 x.f，然后返回 f 字段相对于 x 起始地址的偏移量，包括可能的空洞。

## unsafe.Pointer
unsafe.Pointer是特别定义的一种指针类型（译注：类似C语言中的void*类型的指针），它可以包含任意类型变量的地址。当然，我们不可以直接通过*p来获取unsafe.Pointer指针指向的真实变量的值.

一个普通的*T类型指针可以被转化为unsafe.Pointer类型指针，并且一个unsafe.Pointer类型指针也可以被转回普通的指针，被转回普通的指针类型并不需要和原始的*T类型相同。\
```go
package math

func Float64bits(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }

fmt.Printf("%#016x\n", Float64bits(1.0)) // "0x3ff0000000000000"

```

todo

## 通过cgo调用C代码

Go语言自带的叫cgo的用于支援C语言函数调用的工具。这类工具一般被称为 foreign-function interfaces （简称ffi），并且在类似工具中cgo也不是唯一的。
SWIG（http://swig.org）是另一个类似的且被广泛使用的工具，SWIG提供了很多复杂特性以支援C++的特性。
> bzip提供了相关示例演示了如何将一个C语言库链接到Go语言程序。相反，将Go编译为静态库然后链接到C程序，或者将Go程序编译为动态库然后在C程序中动态加载也都是可行的。

如果要进一步阅读，可以从 https://golang.org/cmd/cgo 开始。