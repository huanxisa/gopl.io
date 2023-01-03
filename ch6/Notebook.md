# 方法 

## 方法声明
> 主要见程序geometry注释

## 基于指针对象的方法

主要好处是原本基于原本对象的方法将会拷贝对象，但是如果希望避免拷贝的话可以基于指针
```go
//方法签名为(*Point).ScaleBy而不是*(Point.ScaleBy)
func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}

```

指针类型方法调用：
```go
r := &Point{1, 2}
r.ScaleBy(2)
fmt.Println(*r) // "{2, 4}"

//语法糖，编译器会隐式地帮我们用&p去调用ScaleBy这个方法。
p.ScaleBy(2)

// 这种简写方法只适用于“变量，下面写法不可以
Point{1, 2}.ScaleBy(2) // compile error: can't take address of Point literal

```
额外注意：如果函数签名不是指针类型方法，但是存在一个指针对象，go也提供了语法糖 如下例：
```go
//方法签名并非指针类型
func (p Point) Distance(q Point) float64 {
	// 所以保持其在方法间传递时的一致性和简短性是不错的主意。这里的建议是可以使用其类型的第一个字母，比如这里使用了Point的首字母p。
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

//变量为指针类型
pptr := &Point{1, 2}
//下面两种写法等价，语法糖
pptr.Distance(q)
(*pptr).Distance(q)

```

### Nil是一个合法的接收器类型
```go
// An IntList is a linked list of integers.
//记得注释说明好
// A nil *IntList represents the empty list.
type IntList struct {
    Value int
    Tail  *IntList
}
// Sum returns the sum of the list elements.
func (list *IntList) Sum() int {
	//重点在这个地方
    if list == nil {
        return 0
    }
    return list.Value + list.Tail.Sum()
}

```

如果一个类型定义非常简单如：type Values map[string][]string，仅仅只有单变量。客户端使用这个变量的时候可以使用map固有的一些操作（make，切片，m[key]等等），也可以使用这里提供的操作方法，或者两者并用，都是可以的
> 见urlvalues

## 通过嵌入结构体来扩展类型
成员变量的延伸访问，已经在第一章讲过，故略

成员函数的延伸访问。
>见coloerpoint，除了表达成员函数的用法，还有匿名指针成员对象的访问
 
### 匿名struct
下面两种写法作用是相同的，但是在Lookup中的写法有细微的区别，可以体会一下
```go
var (
    mu sync.Mutex // guards mapping
    mapping = make(map[string]string)
)

func Lookup(key string) string {
    mu.Lock()
    v := mapping[key]
    mu.Unlock()
    return v
}

//匿名的stuct
var cache = struct {
    sync.Mutex
    mapping map[string]string
}{
    mapping: make(map[string]string),
}


func Lookup(key string) string {
	//更具表达性的名字：cache。因为sync.Mutex字段也被嵌入到了这个struct里，其Lock和Unlock方法也就都被引入到了这个匿名结构中了
    cache.Lock()
    v := cache.mapping[key]
    cache.Unlock()
    return v
}
```
## 方法值和方法表达式

方法值：一个将方法（Point.Distance）绑定到特定接收器变量的函数。这个函数可以不通过指定其接收器即可被调用；即调用时不需要指定接收器。
```go
p := Point{1, 2}
q := Point{4, 6}

distanceFromP := p.Distance        // method value
fmt.Println(distanceFromP(q))      // "5"
var origin Point                   // {0, 0}
fmt.Println(distanceFromP(origin)) // "2.23606797749979", sqrt(5)

scaleP := p.ScaleBy // method value
scaleP(2)           // p becomes (2, 4)
scaleP(3)           //      then (6, 12)
scaleP(10)          //      then (60, 120)

```
```go

//当T是一个类型时，方法表达式可能会写作T.f或者(*T).f，会返回一个函数“值”，这种函数会将其第一个参数用作接收器
p := Point{1, 2}
q := Point{4, 6}

distance := Point.Distance   // method expression
fmt.Println(distance(p, q))  // "5"
fmt.Printf("%T\n", distance) // "func(Point, Point) float64"

scale := (*Point).ScaleBy
scale(&p, 2)
fmt.Println(p)            // "{2 4}"
fmt.Printf("%T\n", scale) // "func(*Point, float64)"

//// 看起来本书中函数和方法的区别是指有没有接收器，而不像其他语言那样是指有没有返回值。
```
### 方法表达式的用法
一个变量来决定调用同一个类型的哪个函数时，方法表达式就显得很有用了。你可以根据选择来调用接收器各不相同的方法。下面的例子，变量op代表Point类型的addition或者subtraction方法

```go
type Point struct{ X, Y float64 }

func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }
func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }

type Path []Point

func (path Path) TranslateBy(offset Point, add bool) {
    var op func(p, q Point) Point
    //通过add决定调用的方法
    if add {
        op = Point.Add
    } else {
        op = Point.Sub
    }
    for i := range path {
        // Call either path[i].Add(offset) or path[i].Sub(offset).
        path[i] = op(path[i], offset)
    }
}

```
在绑定方法时，指针并非万能的，有时会出意外。具体问题具体分析。

## 封装
大写首字母的标识符会从定义它们的包中被导出，小写字母的则不会。

在命名一个getter方法时，我们通常会省略掉前面的Get前缀。
