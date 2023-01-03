#　复合数据类型

主要讲 数组、slice、map和结构体 这四种复合体。

我们将演示如何使用结构体来解码和编码到对应JSON格式的数据，并且通过结合使用模板来生成HTML页面。

数组和结构体都是有固定内存大小的数据结构。相比之下，slice和map则是动态的数据结构，它们将根据需要动态增长。

## 数组

因为数组的长度是固定的，因此在Go语言中很少直接使用数组。和数组对应的类型是Slice（切片），它是可以增长和收缩的动态序列，slice功能也更灵活

### 数组的初始化 

* 数组常规初始化
```go

var q [3]int = [3]int{1, 2, 3}
var r [3]int = [3]int{1, 2}
fmt.Println(r[2]) // "0"

//可以使用特殊的省略号来初始化
q := [...]int{1, 2, 3}
fmt.Printf("%T\n", q) // "[3]int"

//注意 数组的类型是会由长度确认的，所以下面的赋值会报错
q := [3]int{1, 2, 3}
q = [4]int{1, 2, 3, 4} // compile error: cannot assign [4]int to [3]int

// 注意这个地方和python不一样,
r := [...]int{99: -1}  //定义了一个含有100个元素的数组r，最后一个元素被初始化为-1，其它元素都是用0初始化。

```

数组的逻辑比较：

> sha256 程序简单展示了crypto/sha256包如何使用，以及数组比较
```go
a := [2]int{1, 2}
b := [...]int{1, 2}
c := [2]int{1, 3}
fmt.Println(a == b, a == c, b == c) // "true false false"
d := [3]int{1, 2}
fmt.Println(a == d) // 注意这个，由于这两个数组的长度不一致，直接认为不能比较
```

数组作为方法参数时程序的行为：

>函数的每个调用参数将会被赋值给函数内部的参数变量，所以函数参数变量接收的是一个复制的副本，并不是原始调用的变量,我们可以显式地传入一个数组指针

```go

//数组的类型包含了僵化的长度信息。上面的zero函数并不能接收指向[16]byte类型数组的指针，而且也没有任何添加或删除数组元素的方法
func zero(ptr *[32]byte) {
    for i := range ptr {
        ptr[i] = 0
    }
}

```

## Slice

一个slice类型一般写作[]T，其中T代表slice中元素的类型；slice的语法和数组很像，只是没有固定长度而已。

一个slice由三个部分构成：指针、长度和容量。指针指向第一个slice元素对应的底层数组元素的地址，要注意的是slice的第一个元素并不一定就是数组的第一个元素。

长度对应slice中元素的数目；长度不能超过容量，容量一般是从slice的开始位置到底层数据的结尾位置。内置的len和cap函数分别返回slice的长度和容量。

```go
// go 中声明slice可以指定index.
months := [...]string{1: "January", /* ... */, 12: "December"}

//采用 make声明
make([]T, len)
make([]T, len, cap) // same as make([]T, cap)[:len]

// 申明 slice 为nil
var s []int    // len(s) == 0, s == nil
s = nil        // len(s) == 0, s == nil
s = []int(nil) // len(s) == 0, s == nil
s = []int{}    // len(s) == 0, s != nil

```

Slice可以从另一个切片中获取：
```go

Q2 := months[4:7]
summer := months[6:9]
fmt.Println(Q2)     // ["April" "May" "June"]
fmt.Println(summer) // ["June" "July" "August"]

fmt.Println(summer[:20]) // panic: out of range

//注意这种声明可以超过底层slice长度，但是不能超过slice cap
endlessSummer := summer[:5] // extend a slice (within capacity)
fmt.Println(endlessSummer)  // "[June July August September October]"


```

### Slice 的比较：

唯一合法的是和nil比较。

slice之间不能比较，因此我们不能使用==操作符来判断两个slice是否含有全部相等元素。不过标准库提供了高度优化的bytes.Equal函数来判断两个字节型slice是否相等（[]byte），但是对于其他类型的slice，我们必须自己展开每个元素进行比较：

### Slice的元素操作

* append函数 
```go
var runes []rune
for _, r := range "Hello, 世界" {
	//必须先检测slice底层数组是否有足够的容量来保存新添加的元素。如果有足够空间的话，直接扩展slice（依然在原有的底层数组之上），将新添加的y元素复制到新扩展的空间，并返回slice。
	//没有足够的增长空间的话，appendInt函数则会先分配一个足够大的slice用于保存新的结果，先将输入的x复制到新的空间
	//注意这个前面的runes = 不能省略
    runes = append(runes, r)
}
fmt.Printf("%q\n", runes) // "['H' 'e' 'l' 'l' 'o' ',' ' ' '世' '界']"

```

## Map

map的创建：

```go
ages := make(map[string]int) // mapping from strings to ints

ages := map[string]int{
"alice":   31,
"charlie": 34,
}

```

Map的操作：

删除元素：
```go
//元素不在map中也没关系，不在map中会以默认初始值返回
delete(ages, "alice") // remove element ages["alice"]

```

Map中的元素不是一个变量，我们无法对于value进行取址操作：
```go
_ = &ages["bob"] // compile error: cannot take address of map element
```

遍历元素：
Map的迭代顺序是不确定的，并且不同的哈希函数实现可能导致不同的遍历顺序。在实践中，遍历的顺序是随机的，每一次遍历的顺序都不相同。
```go
for name, age := range ages {
    fmt.Printf("%s\t%d\n", name, age)
}
```

判断是否存在于Map中：
```go
age, ok := ages["bob"]
if !ok { /* "bob" is not a key in this map; age == 0. */ }

```

Map之间的比较操作：

* map之间也不能进行相等比较；唯一的例外是和nil进行比较。要判断两个map是否包含相同的key和value，我们必须通过一个循环实现：

采用Map实现Set功能：

* 见代码dedup

如何将无法比较的的类型作为key,例如Slice：

```go
//定义辅助函数，将无法比较类型转化为key

```

> graph 展示了一个图如何使用map来存储有向边的信息，主要是为了描述Map的value可以是聚合类型。

## 结构体

> treesort 程序展示结构体如何嵌套自身指针。 

### 结构体声明
```go
type Employee struct {
	//如果结构体成员名字是以大写字母开头的，那么该成员就是导出的
    ID        int
    //通常，我们只是将相关的成员写到一起。
    Name,Address     string //同种类型可以同时声明。
    DoB       time.Time
    Position  string
    Salary    int
    ManagerID int
}

var dilbert Employee

//一个命名为S的结构体类型将不能再包含S类型的成员：因为一个聚合的值不能包含它自身。（该限制同样适用于数组。）但是S类型的结构体可以包含*S指针类型的成员
//具体见treesort

//如果结构体没有任何成员的话就是空结构体，写作struct{}。它的大小为0，也不包含任何信息，但是有时候依然是有价值的。有些Go语言程序员用map来模拟set数据结构时，用它来代替map中布尔类型的value，只是强调key的重要性，但是因为节约的空间有限，而且语法比较复杂，所以我们通常会避免这样的用法。
seen := make(map[string]struct{}) // set of strings
```
### 结构体变量成员的访问
```go
dilbert.Salary -= 5000 // demoted, for writing too few lines of code

//对于里面的成员变量取指针，访问成员变量的话就需要取内容操作符
position := &dilbert.Position
*position = "Senior " + *position // promoted, for outsourcing to Elbonia

//
var employeeOfTheMonth *Employee = &dilbert
employeeOfTheMonth.Position += " (proactive team player)"

```
### 结构体赋值

```go
type Point struct{ X, Y int }

//第一种按照顺序赋值
p := Point{1, 2}
// 第二种按照名称赋值
anim := gif.GIF{LoopCount: nframes}

//两种不同形式的写法不能混合使用。而且，你不能企图在外部包中用第一种顺序赋值的技巧来偷偷地初始化结构体中未导出的成员。
package p
type T struct{ a, b int } // a and b are not exported

package q
import "p"
var _ = p.T{a: 1, b: 2} // compile error: can't reference a, b
//并没有显式提到未导出的成员，但是这样企图隐式使用未导出成员的行为也是不允许的。
var _ = p.T{1, 2}       // compile error: can't reference a, b

//结构体的指针，第一种
pp := &Point{1, 2}
//第二种
pp := new(Point)
*pp = Point{1, 2}

//但是下面的不能通过编译
w = Wheel{8, 8, 5, 20}                       // compile error: unknown fields
w = Wheel{X: 8, Y: 8, Radius: 5, Spokes: 20} // compile error: unknown fields
//对于复杂对象构造只能用下面两种。
w = Wheel{Circle{Point{8, 8}, 5}, 20}
w = Wheel{
Circle: Circle{
Point:  Point{X: 8, Y: 8},
Radius: 5,
},
Spokes: 20, // NOTE: trailing comma necessary here (and at Radius)
}
```
### 结构体比较
```go
//如果结构体的全部成员都是可以比较的，那么结构体也是可以比较的，那样的话两个结构体将可以使用==或!=运算符进行比较

type Point struct{ X, Y int }

p := Point{1, 2}
q := Point{2, 1}
// 下面两种等价
fmt.Println(p.X == q.X && p.Y == q.Y) // "false"
fmt.Println(p == q)                   // "false"

```
### 结构体的匿名访问

```go
type Point struct {
X, Y int
}

type Circle struct {
Center Point
Radius int
}

type Wheel struct {
Circle Circle
Spokes int
}

type Circle struct {
Point
Radius int
}

type Wheel struct {
Circle
Spokes int
}

// 下面写法访问内部的成员变量变得简单了
var w Wheel
w.X = 8            // equivalent to w.Circle.Point.X = 8
w.Y = 8            // equivalent to w.Circle.Point.Y = 8
w.Radius = 5       // equivalent to w.Circle.Radius = 5
w.Spokes = 20

// 注意如果匿名成员是私有的，也可以这样访问匿名成员
w.X = 8 // equivalent to w.circle.point.X = 8

```

> 一个很奇怪的点  
```go
//todo： 下面这两种写法会对于第三句写法有影响
// 返回的是指针，第三句代码能通过编译
func EmployeeByID(id int) *Employee { /* ... */ }
// 返回的是实体，无法通过编译，因为在赋值语句的左边并不确定是一个变量（译注：调用函数返回的是值，并不是一个可取地址的变量）。无法理解最后一句，反正无脑返回指针就行。
func EmployeeByID(id int) Employee { /* ... */ }

EmployeeByID(id).Salary = 0 // fired for... no real reason

```

## JSON

不过JSON使用的是\Uhhhh转义数字来表示一个UTF-16编码（译注：UTF-16和UTF-8一样是一种变长的编码，有些Unicode码点较大的字符需要用4个字节表示；而且UTF-16还有大端和小端的问题），而不是Go语言的rune类型。

> movie程序展示了如何将go结构体序列化json和反序列化json，github程序展示如何从流对象序列化JSON，issues主要用来调用github包。


## 文本和HTML模板

有时候会需要复杂的打印格式，这时候一般需要将格式化代码分离出来以便更安全地修改。这些功能是由text/template和html/template等模板包提供的，它们提供了一个将变量值填充到一个文本或HTML格式的模板的机制。

> issuesreport程序展示了文本模板怎么使用，并且如何选择输入和输出源。issueshtml程序展示如何。autoescape展示如何设置受信任的字符串和不受信任的字符串

关于更多有关于模板的用法见:

```shell
$ go doc text/template
$ go doc html/template
```
