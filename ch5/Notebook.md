# 函数

函数这章主要包含函数声明、递归函数、匿名函数、错误处理和函数其它的很多特性。

## 函数声明

函数的类型：又称为函数的签名。如果两个函数形式参数列表和返回值列表中的`变量类型`一一对应，那么这两个函数被认为有相同的类型或签名。形参和返回值的`变量名不影响函数签名`.

每一次函数调用都必须按照声明顺序为所有参数提供实参（参数值）。在函数调用时，Go语言没有默认参数值，也没有任何方法可以通过参数名指定形参。

函数的形参是实参的拷贝。对形参进行修改不会影响实参。但是，如果实参包括引用类型，如指针，slice(切片)、map、function、channel等类型，实参可能会由于函数的间接引用被修改。

```go
//如果函数返回一个无名变量或者没有返回值，返回值列表的括号是可以省略的
func name(parameter-list) (result-list) {
body
}

//如果一组形参或返回值有相同的类型，我们不必为每个形参都写出参数类型。下面2个声明是等价的：
func f(i, j, k int, s, t string)                 { /* ... */ }
func f(i int, j int, k int, s string, t string) { /* ... */ }

//一些针对于方法特殊trick声明
//接收参数却不使用，于是使用_代替
func first(x int, _ int) int { return x }
//两个参数一个都不使用
func zero(int, int) int      { return 0 }
```

偶尔遇到没有函数体的函数声明，这表示该函数不是以Go实现的。这样的声明定义了函数签名。

```go
package math

func Sin(x float64) float //implemented in assembly language

```

## 递归

> findlinks1程序主要描述如何遍历一个html文件打印其中的链接，主要是想通过visit函数来实现。
> outline程序也是为了说明递归，主要的作用递归的方式遍历整个HTML结点树。

## 多返回值

> findlinks2中的findLinks 主要展示多返回值的基础用法，然后说明即使存在垃圾回收机制，为什么还是需要关闭文件操作？

### 多返回值得接收及传递

```go
//调用多返回值函数时，返回给调用者的是一组值，调用者必须显式的将这些值分配给变量:
links, err := findLinks(url)

//如果某个值不被使用，可以将其分配给blank identifier:
links, _ := findLinks(url) // errors ignored

//多返回值的传递。一个函数内部可以将另一个有多返回值的函数调用作为返回值
func findLinksLog(url string) ([]string, error) {
log.Printf("findLinks %s", url)
return findLinks(url)
}

//多返回值的传递。接受多参数的函数时，可以将一个返回多参数的函数调用作为该函数的参数。
//下面两条语句等价
log.Println(findLinks(url))
links, err := findLinks(url)
log.Println(links, err)

//按照惯例，函数的最后一个bool类型的返回值表示函数是否运行成功，error类型的返回值代表函数的错误信息
```

### bare return

```go
//一个函数所有的返回值都有显式的变量名，那么该函数的return语句可以省略操作数。这称之为bare return。
//每一个return语句等价于 return words, images, err
func CountWordsAndImages(url string) (words, images int, err error) {
resp, err := http.Get(url)
if err != nil {
return
}
doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil {
err = fmt.Errorf("parsing HTML: %s", err)
return
}
words, images = countWordsAndImages(doc)
return
}

```

## 错误

错误和panic异常的区别：

* 一部分函数总是能成功的运行。**除非遇到灾难性的、不可预料的情况**，比如运行时的内存溢出。导致这种错误的原因很复杂，难以处理，从错误中恢复的可能性也很低。
* 这种情况下会引发panic异常。panic是来自被调用函数的信号，表示发生了某个已知的bug。一个良好的程序永远不应该发生panic异常。

对于那些将运行失败看作是预期结果的函数，它们会返回一个额外的返回值，通常是最后一个。通常被命名为ok，对于复杂的操作，尤其是对I/O操作而言，用户需要了解更多的错误信息。因此，额外的返回值不再是简单的布尔类型，而是error类型。注意和panic区分开

Go使用控制流机制（如if和return）处理错误，这使得编码人员能更多的关注错误处理。

### 错误处理策略

通常，当函数返回non-nil的error时，其他的返回值是未定义的（undefined）.然而，有少部分函数在发生错误时，仍然会返回一些有用的返回值。正确的处理方式应该是先处理这些不完整的数据，再处理错误。因此对函数的返回值要有清晰的说明，以便于其他人使用。

* 直接传播错误

```go

resp, err := http.Get(url)
if err != nil{
return nil, err
}
//此外如果方法中可能有多个处理的结点出现错误，我们需要提供准确的信息帮助定位。
//如下所示，把信息提供完全
//当对html.Parse的调用失败时，findLinks不会直接返回html.Parse的错误，因为缺少两条重要信息：1、发生错误时的解析器（html parser）；2、发生错误的url。因此，findLinks构造了一个新的错误信息，既包含了这两项，也包括了底层的解析出错的信息。
doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil {
return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
}

```

* 如果错误的发生是偶然性的，或由不可预知的问题导致的。一个明智的选择是重新尝试失败的操作。在重试时，我们需要限制重试的时间间隔或重试的次数，防止无限制的重试。

> wait中的WaitForServer函数展示如果发生错误重试，注意重试策略的设置

* 输出错误信息并结束程序。需要注意的是，这种策略只应在main中执行。

```go
// (In function main.)
if err := WaitForServer(url); err != nil {
//log.Fatalf("Site is down: %v\n", err) 都默认会在错误信息之前输出时间信息。
fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
os.Exit(1)
}

```

* 我们只需要输出错误信息就足够了，不需要中断程序的运行。我们可以通过log包提供

```go
//通过log包提供函数
if err := Ping(); err != nil {
log.Printf("ping failed: %v; networking disabled", err)
}
//标准错误流输出错误信息
if err := Ping(); err != nil {
fmt.Fprintf(os.Stderr, "ping failed: %v; networking disabled\n", err)
}

```

* 我们可以直接忽略掉错误

```go
dir, err := ioutil.TempDir("", "scratch")
if err != nil {
return fmt.Errorf("failed to create temp dir: %v", err)
}
// ...use temp dir…
//虽然程序没有处理错误，但程序的逻辑不会因此受到影响。我们应该在每次函数调用后，都养成考虑错误处理的习惯，当你决定忽略某个错误时，你应该清晰地写下你的意图。
os.RemoveAll(dir) // ignore errors; $TMPDIR is cleaned periodically

```

### 文件结尾错误（EOF）

文件有一个独特的定义错误的类型：

```go
 if err == io.EOF {
break // finished reading
}
```

## 函数值

函数向其他值一样，拥有类型，可以赋值给其他变量传递给函数

```go
    func square(n int) int { return n * n }
func negative(n int) int { return -n }
func product(m, n int) int { return m * n }

f := square
fmt.Println(f(3)) // "9"

f = negative
fmt.Println(f(3)) // "-3"
fmt.Printf("%T\n", f) // "func(int) int"
//注意这里的函数签名，这一句无法通过编译
f = product // compile error: can't assign func(int, int) int to func(int) int

//函数值的零值为nil，调用值为nil的函数值会引起panic错误
var f func (int) int
f(3)

```

### 函数值得布尔比较

* 只能和nil比较，不能作为map的key

### 实现类似于lambda表达式的效果

> outline2 主要是通过forEachNode将许多逻辑抽离分开

```go
func add1(r rune) rune { return r + 1 }
//Map/reduce
fmt.Println(strings.Map(add1, "HAL-9000")) // "IBM.:111"
```

## 匿名函数

函数字面量的语法和函数声明相似，区别在于func关键字后没有函数名。函数值字面量是一种表达式，它的值被称为匿名函数。

### 匿名函数作为闭包

```go
// squares返回一个匿名函数。
// 该匿名函数每次被调用时都会返回下一个数的平方。
func squares() func () int {
var x int
return func () int {
x++
return x * x
}
}
func main() {
f := squares()
fmt.Println(f()) // "1"
fmt.Println(f()) // "4"
fmt.Println(f()) // "9"
fmt.Println(f()) // "16"
}

```

> toposort 主要展示了一个拓扑排序问题，额外的补充了，如果想要递归调用匿名函数不是一件容易的事
> links 主要将findlinks中的visit改写了，具体见代码注释
> findlinks3 展示如何广度优先的遍历树。

### 局部变量的陷阱

```go
// 这个是正确的的
var rmdirs []func ()
for _, d := range tempDirs() {
dir := d // NOTE: necessary!
os.MkdirAll(dir, 0755) // creates parent directories too
rmdirs = append(rmdirs, func() {
os.RemoveAll(dir)
})
}
// ...do some work…
for _, rmdir := range rmdirs {
rmdir() // clean up
}

// 这个是错误的，循环变量dir在这个词法块中被声明。在该循环中生成的所有函数值都共享相同的循环变量。需要注意，函数值中记录的是循环变量的内存地址，而不是循环变量某一时刻的值。
//for循环已完成，dir中存储的值等于最后一次迭代的值。这意味着，每次对os.RemoveAll的调用删除的都是相同的目录
var rmdirs []func ()
for _, dir := range tempDirs() {
os.MkdirAll(dir, 0755)
rmdirs = append(rmdirs, func () {
os.RemoveAll(dir) // NOTE: incorrect!
})
}


```

## 可变参数

首先接收一个必备的参数，之后接收任意个数的后续参数。

### 声明可变参数

在声明可变参数函数时，需要在参数列表的最后一个参数类型之前加上省略符号“...”，这表示该函数会接收任意数量的该类型参数。

```go
func sum(vals ...int) int {
total := 0
for _, val := range vals {
total += val
}
return total
}

// 如果原始参数已经是切片类型，我们该如何传递给sum？只需在最后一个参数后加上省略符
values := []int{1, 2, 3, 4}
fmt.Println(sum(values...)) // "10"

```

注意：int 型参数的行为看起来很像切片类型，但实际上，可变参数函数和以切片作为参数的函数是不同的。

```go
func f(...int) {}
func g([]int) {}
fmt.Printf("%T\n", f) // "func(...int)"
fmt.Printf("%T\n", g) // "func([]int)"

```

## Deferred函数

> title1 获取HTML页面并输出页面的标题。title函数会检查服务器返回的Content-Type字段，如果发现页面不是HTML，将终止函数运行，返回错误。
> title2 是对于title1的改进，主要表现是释放资源的defer操作。

释放资源的defer应该直接跟在获取资源的语句后。或是处理互斥锁 调试复杂程序时，defer机制也常被用于记录何时进入和退出函数

### defer访问return之后的值

```go
//被延迟执行的匿名函数甚至可以修改函数返回给调用者的返回值：
func triple(x int) (result int){
//	注意这里访问了result中的值
defer func () { result += x }()
return double(x)
}
fmt.Println(triple(4)) // "12"
```

### 循环中的defer

在循环体中的defer语句需要特别注意，因为只有在函数执行完毕后，这些被延迟的函数才会执行。下面的代码会导致系统的文件描述符耗尽，因为在所有文件都被处理之前，没有文件会被关闭。

一种解决方法是将循环体中的defer语句移至另外一个函数。在每次循环时，调用这个函数。

```go
for _, filename := range filenames {
f, err := os.Open(filename)
if err != nil {
return err
}
defer f.Close() // NOTE: 可能导致文件描述符好景
// ...process f…
}
```

### 文件系统的不推荐采用defer

> 见fetch中的注释

## Panic异常

有些错误只能在运行时检查，如数组访问越界、空指针引用等。这些运行时错误会引起panic异常。

当某些不应该发生的场景发生时，我们就应该调用panic。比如，当程序到达了某条逻辑上不可能到达的路径

```go
switch s := suit(drawCard()); s {
case "Spades":   // ...
case "Hearts":   // ...
case "Diamonds": // ...
case "Clubs": // ...
default:
panic(fmt.Sprintf("invalid suit %q", s)) // Joker?
}
```

除非你能提供更多的错误信息，或者能更快速的发现错误，否则不需要断言那些运行时会检查的条件.

我们应该使用Go提供的错误机制，而不是panic，尽量避免程序的崩溃。在健壮的程序中，任何可以预料到的错误，如不正确的输入、错误的配置或是失败的I/O操作都应该被优雅的处理，最好的处理方式，就是使用Go的错误机制

### panic的牵扯到的简单性

```go
package regexp

func Compile(expr string) (*Regexp, error) { /* ... */ }

// 如果调用者明确的知道正确的输入不会引起函数错误时，要求调用者检查这个错误是不必要和累赘的。可以采用下面这种
//Must前缀是一种针对此类函数的命名约定
func MustCompile(expr string) *Regexp {
	re, err := Compile(expr)
	if err != nil {
		panic(err)
	}
	return re
}

```

### panic异常和程序调用之间的关系
> 见defer1 延迟函数的调用在释放堆栈信息之前。

## Recover捕获 

通常来说，不应该对panic异常做任何处理，但有时，也许我们可以从异常中恢复，至少我们可以在程序崩溃前，做一些操作。举个例子，当web服务器遇到不可预料的严重问题时，在崩溃前应该将所有的连接关闭；如果不做任何处理，会使得客户端一直处于等待状态。如果web服务器还在开发阶段，服务器甚至可以将异常信息反馈到客户端，帮助调试。

### 具体实现形式
在panic之前捕获异常
```go
func Parse(input string) (s *Syntax, err error) {
    defer func() {
        if p := recover(); p != nil {
            //panic value被附加到错误信息中；并用err变量接收错误信息，返回给调用者
            err = fmt.Errorf("internal error: %v", p)
        }
    }()
    // ...parser...
}

```
### 代码编写规范
公有的API应该将函数的运行失败作为error返回，而不是panic。同样的，你也不应该恢复一个由他人开发的函数引起的panic，比如说调用者传入的回调函数，因为你无法确保这样做是安全的。

有时我们很难完全遵循规范,安全的做法是有选择性的recover。为了标识某个panic是否应该被恢复，我们可以将`panic value设置成特殊类型`。在recover时对panic value进行检查，如果发现panic value是特殊类型，就将这个panic作为error处理，如果不是，则按照正常的panic进行处理