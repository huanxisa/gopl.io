# 反射
能够在运行时更新变量和检查它们的值、调用它们的方法和它们支持的内在操作，而不需要在编译时就知道这些变量的具体类型。

## reflect.Type和reflect.Value
函数 reflect.TypeOf 接受任意的 interface{} 类型，并以 reflect.Type 形式返回其动态类型：
```go
t := reflect.TypeOf(3)  // a reflect.Type
fmt.Println(t.String()) // "int"
fmt.Println(t)          // "int"
```

它总是返回具体的类型。

reflect.Value 可以装载任意类型的值。函数 reflect.ValueOf 接受任意的 interface{} 类型，并返回一个装载着其动态值的 reflect.Value。reflect.Value 也可以持有一个接口值。
```go
v := reflect.ValueOf(3) // a reflect.Value
fmt.Println(v)          // "3"
fmt.Printf("%v\n", v)   // "3"
fmt.Println(v.String()) // NOTE: "<int Value>"
```

reflect.ValueOf 的逆操作是 reflect.Value.Interface 方法。它返回一个 interface{} 类型，装载着与 reflect.Value 相同的具体值：
```go
v := reflect.ValueOf(3) // a reflect.Value
//v.Kind()获取类型
x := v.Interface()      // an interface{}
i := x.(int)            // an int
fmt.Printf("%d\n", i)   // "3"
```
> format展示了go如何对于一个值获取类型，并且走不同分支。display写了一个循环打印对象内部结构的例子。
> sexpr也是一个例子

## 通过reflect.value修改值
> sexpr/decode 说了些关于本节的应用

1. 可取地址的reflect.Value来访问变量需要三个步骤。

2. 第一步是调用Addr()方法，它返回一个Value，里面保存了指向变量的指针。

3. 然后是在Value上调用Interface()方法，也就是返回一个interface{}，里面包含指向变量的指针。 
   
4. 最后，如果我们知道变量的类型，我们可以使用类型的断言机制将得到的interface{}类型的接口强制转为普通的类型指针。这样我们就可以通过这个普通指针来更新变量

5. 利用反射机制并不能修改这些未导出的成员.
```go
//有个缺点是必须提前知道这个类型
x := 2
d := reflect.ValueOf(&x).Elem()   // d refers to the variable x
//第一种 通过指针改变值
px := d.Addr().Interface().(*int) // px := &x
*px = 3                           // x = 3
fmt.Println(x)                    // "3"
//第二种 而是通过调用可取地址的reflect.Value的reflect.Value.Set方法来更新对应的值。
d.Set(reflect.ValueOf(4))
fmt.Println(x) // "4"
```

### 额外补充
* Set方法将在运行时执行和编译时进行类似的可赋值性约束的检查
* 对一个不可取地址的reflect.Value调用Set方法也会导致panic异常，要确保改类型的变量可以接受对应的值
* 对于一个引用interface{}类型的reflect.Value调用SetInt会导致panic异常，即使那个interface{}变量对于整数类型也不行。
```go
x := 1
rx := reflect.ValueOf(&x).Elem()
rx.SetInt(2)                     // OK, x = 2
rx.Set(reflect.ValueOf(3))       // OK, x = 3
rx.SetString("hello")            // panic: string is not assignable to int
rx.Set(reflect.ValueOf("hello")) // panic: string is not assignable to int

var y interface{}
ry := reflect.ValueOf(&y).Elem()
ry.SetInt(2)                     // panic: SetInt called on interface Value
ry.Set(reflect.ValueOf(3))       // OK, y = int(3)
ry.SetString("hello")            // panic: SetString called on interface Value
ry.Set(reflect.ValueOf("hello")) // OK, y = "hello"

* CanAddr方法并不能正确反映一个变量是否是可以被修改的。另一个相关的方法CanSet是用于检查对应的reflect.Value是否是可取地址并可被修改的。
```
fmt.Println(fd.CanAddr(), fd.CanSet()) // "true false"

```go
//CanAddr方法并不能正确反映一个变量是否是可以被修改的。另一个相关的方法CanSet是用于检查对应的reflect.Value是否是可取地址并可被修改的
```
### 是否可寻址
对于reflect.Values也有类似的区别。有一些reflect.Values是可取地址的；其它一些则不可以。所有通过reflect.ValueOf(x)返回的reflect.Value都是不可取地址的
```go
//通过调用reflect.Value的CanAddr方法来判断其是否可以被取地址
x := 2                   // value   type    variable?
a := reflect.ValueOf(2)  // 2       int     no
b := reflect.ValueOf(x)  // 2       int     no
c := reflect.ValueOf(&x) // &x      *int    no
d := c.Elem()            // 2       int     yes (x)

fmt.Println(a.CanAddr()) // "false"
fmt.Println(b.CanAddr()) // "false"
```
## 显示一个类型的方法集
> methods 使用reflect.Type来打印任意值的类型和枚举它的方法，使用reflect.Type来打印任意值的类型和枚举它的方法，但是这个例子中只用到了它的类型。

## 忠告  
反射需要被小心使用。
* 第一个原因是，基于反射的代码是比较脆弱的。而反射则是在真正运行到的时候才会抛出panic异常，可能是写完代码很久之后了，而且程序也可能运行了很长的时间。需要非常小心地检查每个reflect.Value的对应值的类型、是否可取地址，还有是否可以被修改等。
避免这种因反射而导致的脆弱性的问题的最好方法是将所有的反射相关的使用控制在包的内部，如果可能的话避免在包的API中直接暴露reflect.Value类型，这样可以限制一些非法输入。
* 反射同样降低了程序的安全性，还影响了自动化重构和分析工具的准确性，因为它们无法识别运行时才能确认的类型信息。
* 第三个原因，基于反射的代码通常比正常的代码运行速度慢一到两个数量级。