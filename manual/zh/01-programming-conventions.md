# 一、编程规约

## （一）命名规约

### 【MUST】NAMING-001 包名使用小写单词连写，不使用下划线、驼峰或复数。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

包名是调用方引用标识符的前缀，全小写、无分隔符才能保持一致。多个单词直接连写，不使用下划线或驼峰，也不加复数。

**为什么**

> In Go, package names must be concise and use only lowercase letters and numbers (e.g., k8s, oauth2). Multi-word package names should remain unbroken and in all lowercase (e.g., tabwriter instead of tabWriter, TabWriter, or tab_writer).
>
> —— https://google.github.io/styleguide/go/decisions#package-names

包名出现在每一个外部引用点，命名不一致会持续抬高阅读成本；stylecheck 等工具对包名有既定约定，跟随约定可自动检出偏差。

**正例**

```go
package orderbook
```

**反例**

```go
package order_book
```

**依据**

- https://go.dev/wiki/CodeReviewComments#package-names
- https://google.github.io/styleguide/go/decisions

**检测**：golangci-lint stylecheck(ST1003)（检查包名标识符）

### 【MUST】NAMING-002 标识符不使用下划线分隔单词。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

变量、常量、类型、函数与包名一律用 MixedCaps，不用下划线分隔单词。例外只有三类：

- 仅供生成代码使用的包名
- 测试文件里 Test、Benchmark、Example 函数名，可用下划线分组用例
- 与操作系统或 cgo 打交道的底层库复用既有标识符

**为什么**

> Names in Go should in general not contain underscores. There are three exceptions to this principle:
>
> —— https://google.github.io/styleguide/go/decisions#underscores

下划线命名来自其它语言习惯，Go 的约定是 MixedCaps；两套写法混在同一仓库，同类标识符出现两种形态，检索与阅读都要多一步判断。

**正例**

```go
func readConfig() {}
```

**反例**

```go
func read_config() {}
```

**依据**

- https://google.github.io/styleguide/go/decisions#underscores
- https://go.dev/wiki/CodeReviewComments#mixed-caps

**检测**：golangci-lint stylecheck(ST1003)（检查标识符命名风格）

### 【MUST】NAMING-003 首字母缩略词在标识符中大小写统一。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

URL、ID、HTTP 一类缩略词要么全大写、要么全小写，不写成 Url、Id。按缩略词在英文里的写法决定大小写（如 iOS、gRPC），需要导出时把首字母改为大写。

**为什么**

> Words in names that are initialisms or acronyms (e.g., URL and NATO) should have the same case. URL should appear as URL or url (as in urlPony, or URLPony), never as Url.
>
> —— https://google.github.io/styleguide/go/decisions#initialisms

大小写混写会让同一概念出现 URL、Url、url 多种拼写，按名字搜索与替换都会漏；工具对这类命名有既定判定，跟随约定即可自动检出。

**正例**

```go
func ServeHTTP(w http.ResponseWriter, r *http.Request) {}
```

**反例**

```go
func ServeHttp(w http.ResponseWriter, r *http.Request) {}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#initialisms
- https://google.github.io/styleguide/go/decisions#initialisms

**检测**：golangci-lint stylecheck(ST1003)

### 【SHOULD】NAMING-004 名称长度与作用域范围成正比。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

只在一两行内使用的循环下标、方法接收者用单字母即可；跨整个函数或包级作用域的变量用完整词。同一作用域里有多个相近概念时，加限定词区分。

**为什么**

> The general rule of thumb is that the length of a name should be proportional to the size of its scope and inversely proportional to the number of times that it is used within that scope.
>
> —— https://google.github.io/styleguide/go/decisions#variable-names

短名字在短作用域里省去噪音，放到长作用域里会让人读到后面忘了它指什么；按作用域调长度，读者不必回翻声明。

**正例**

```go
orderCount := 0
for _, order := range orders {
	if order.Paid {
		orderCount++
	}
}
```

**反例**

```go
c := 0
for _, o := range orders {
	if o.Paid {
		c++
	}
}
```

**依据**

- https://google.github.io/styleguide/go/decisions#variable-names
- https://go.dev/wiki/CodeReviewComments#variable-names

**检测**：人工核对（人工判断名称长度是否与作用域匹配）

### 【SHOULD】NAMING-005 名称不复刻类型与外围上下文信息。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

类型由声明给出，包名、方法名、文件名已经提供上下文，命名里不再重复，只在类型确实需要区分时加限定词。

- 用 users，不用 userSlice、numUsers
- 在 report 包里用 Report，不用 ReportReport
- 在 UserCount 方法里用 count，不用 userCount

**为什么**

> Names that include information from their surrounding context often create extra noise without benefit. The package name, method name, type name, function name, import path, and even filename can all provide context that automatically qualifies all names within.
>
> —— https://google.github.io/styleguide/go/decisions#external-context-vs-local-names

重复的类型与上下文不增加分辨力，只加长名字与代码行，跨行断行也跟着变多；去掉冗余后，名字仍能唯一指向对应值。

**正例**

```go
var users []User
```

**反例**

```go
var userSlice []User
```

**依据**

- https://google.github.io/styleguide/go/decisions#variable-name-vs-type
- https://google.github.io/styleguide/go/decisions#external-context-vs-local-names

**检测**：人工核对（核对名称是否重复了类型或外围上下文）

### 【SHOULD】NAMING-006 接收者名用类型名的短缩写，同类型内保持一致。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

接收者名通常一到两个字母，取自类型名缩写；同一个类型的所有方法用同一个接收者名。不用 this、self，也不用不到接收者时留一个占位名。

**为什么**

> The name of a method's receiver should be a reflection of its identity; often a one or two letter abbreviation of its type suffices (such as "c" or "cl" for "Client"). Don't use generic names such as "me", "this" or "self". Be consistent, too: if you call the receiver "c" in one method, don't call it "cl" in another.
>
> —— https://go.dev/wiki/CodeReviewComments#receiver-names

同一类型的方法散在多个文件，各用各的接收者名会让人以为指向不同对象；this、self 是其它语言的惯性，在 Go 里不是约定。

**正例**

```go
func (w *ReportWriter) Write(p []byte) (int, error) {}
```

**反例**

```go
func (this *ReportWriter) Write(p []byte) (int, error) {}
```

**依据**

- https://google.github.io/styleguide/go/decisions#receiver-names
- https://go.dev/wiki/CodeReviewComments#receiver-names

**检测**：人工核对（核对同一类型的接收者名是否一致、是否用了 this / self）

## （二）常量与枚举规约

### 【MUST】CONST-001 常量名用 MixedCaps，不用全大写或 k 前缀。

- 归属：编程规约/常量与枚举规约
- 起始版本：Go 1.0
- 主题：naming

导出常量首字母大写，未导出小写，单词之间首字母大写。不使用 MAX_LENGTH 一类全大写，也不使用 kMaxLength 的 k 前缀。

**为什么**

> Constant names must use MixedCaps like all other names in Go. ... Do not use non-MixedCaps constant names or constants with a K prefix.
>
> —— https://google.github.io/styleguide/go/decisions#constant-names

全大写与 k 前缀来自其它语言，Go 里所有标识符共用一套 MixedCaps 约定；常量单独换一套写法，同一仓库就会出现两种命名风格。

**正例**

```go
const MaxPacketSize = 512
```

**反例**

```go
const MAX_PACKET_SIZE = 512
```

**依据**

- https://go.dev/wiki/CodeReviewComments#mixed-caps
- https://google.github.io/styleguide/go/decisions#constant-names

**检测**：golangci-lint stylecheck(ST1003)

### 【SHOULD】CONST-002 常量按用途命名，不复刻它的取值。

- 归属：编程规约/常量与枚举规约
- 起始版本：Go 1.0
- 主题：naming

常量名说明它代表什么，而不是等于多少；除了取值之外没有别的含义时，不必定义成常量。同一取值承担不同用途时，各自命名。

**为什么**

> Name constants based on their role, not their values. If a constant does not have a role apart from its value, then it is unnecessary to define it as a constant.
>
> —— https://google.github.io/styleguide/go/decisions#constant-names

按取值命名（Twelve = 12、ThreeSeconds = 3）在取值调整后名字就变成错的，调用处也没表达出用途；按用途命名才能承载调用处的意图。

**正例**

```go
const maxRetries = 3
```

**反例**

```go
const three = 3
```

**依据**

- https://google.github.io/styleguide/go/decisions#constant-names

**检测**：人工核对（核对常量名是否只是在复述取值）

### 【MAY】CONST-003 同源递增的枚举值用 iota 生成。

- 归属：编程规约/常量与枚举规约
- 起始版本：Go 1.0
- 主题：constants

一组同源、按序递增的枚举用 iota 生成，避免手写数字后中途插值忘改。iota 的取值随它在 const 块内的行数变化，块里插入行后要复核后续取值。

**为什么**

> In Go, enumerated constants are created using the iota enumerator. Since iota can be part of an expression and expressions can be implicitly repeated, it is easy to build intricate sets of values.
>
> —— https://go.dev/doc/effective_go#constants

手写数字时在中间插入一个枚举值，后面的编号要逐个改，漏改一个就产生重复值；iota 按行生成，插入行不会与既有值冲突。

**正例**

```go
const (
	StatusPending = iota
	StatusRunning
	StatusDone
)
```

**反例**

```go
const (
	StatusPending = 0
	StatusRunning = 1
	StatusDone    = 1
)
```

**依据**

- https://go.dev/doc/effective_go#constants
- https://go.dev/ref/spec#Constant_declarations

**检测**：人工核对（核对同源递增的枚举是否手写了数字）

## （三）格式化与代码风格规约

### 【SHOULD】STYLE-001 用 any 代替 interface{}。

- 归属：编程规约/格式化与代码风格规约
- 起始版本：Go 1.18
- 主题：modernization

Go 1.18 引入 any 作为 interface{} 的别名，两者完全等价。新代码统一写 any；改动到旧代码里的 interface{} 时顺带替换。

**为什么**

> The new predeclared identifier any is an alias for the empty interface. It may be used instead of interface{}.
>
> —— https://go.dev/doc/go1.18

any 是官方推荐写法，语义直接，与泛型约束里的写法一致。

**正例**

```go
func printValue(v any) {
	fmt.Println(v)
}
```

**反例**

```go
func printValue(v interface{}) {
	fmt.Println(v)
}
```

**依据**

- https://go.dev/doc/go1.18
- https://google.github.io/styleguide/go/decisions

**检测**：golangci-lint modernize(any)（同一分析器 go fix 可直接改写）

## （四）函数与方法规约

### 【MUST】FUNC-001 导出的函数与方法放在文件前部，未导出的放在后部。

- 归属：编程规约/函数与方法规约
- 起始版本：Go 1.0
- 主题：code-organization

同一个文件里，导出的函数排在未导出的函数前面，某个类型导出的方法排在它未导出的方法前面。这条只管导出与未导出的相对次序，类型与构造函数放在哪个位置不在此列。

**为什么**

> Functions should be sorted in rough call order. Functions in a file should be grouped by receiver. Therefore, exported functions should appear first in a file, after struct, const, var definitions.
>
> —— https://github.com/uber-go/guide/blob/master/style.md#function-grouping-and-ordering

文件从上往下读，先看到的是这个包对外提供什么，再是需要时才展开的实现细节；对外 API 夹在私有辅助之间，读代码的人得先跳过不关心的实现才能找到入口。次序固定后，评审和工具都能直接判定。

**正例**

```go
func Load(path string) (*Manual, error) {
	return readJSON(path)
}

func readJSON(path string) (*Manual, error) {
	// ...
}
```

**反例**

```go
func readJSON(path string) (*Manual, error) {
	// ...
}

func Load(path string) (*Manual, error) {
	return readJSON(path)
}
```

**依据**

- https://github.com/uber-go/guide/blob/master/style.md#function-grouping-and-ordering
- https://github.com/manuelarte/funcorder

**检测**：golangci-lint funcorder

### 【SHOULD】FUNC-002 返回值含义不自明时用命名返回值。

- 归属：编程规约/函数与方法规约
- 起始版本：Go 1.0
- 主题：naming

匿名返回值已经能说明含义时不额外命名；返回多个同类型值、或返回值的含义从签名看不出来时，给结果参数命名，让签名与 godoc 自解释。判断依据是含义是否清楚，与函数长短无关。

- 返回两个 float64 时命名：`(lat, long float64, err error)`
- 类型已经说明含义时不命名：`(*Manual, error)`

**为什么**

> On the other hand, if a function returns two or three parameters of the same type, or if the meaning of a result isn't clear from context, adding names may be useful in some contexts. ... Clarity of docs is always more important than saving a line or two in your function.
>
> —— https://go.dev/wiki/CodeReviewComments#named-result-parameters

返回两个同类型值时，光看签名分不清哪个是哪个，只能翻实现或文档；命名让签名本身说明含义。含义已经清楚的返回值再命名，只会给 godoc 增加噪音。

**正例**

```go
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return x, y
}
```

**反例**

```go
func split(sum int) (int, int) {
	x := sum * 4 / 9
	y := sum - x
	return x, y
}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#named-result-parameters

**检测**：人工核对（核对返回值含义不自明时是否命名）

## （五）数据结构规约

### 【SHOULD】DATA-001 声明空切片时用 var 得到 nil 切片。

- 归属：编程规约/数据结构规约
- 起始版本：Go 1.0
- 主题：data-structures

空切片统一用 var 声明，非 nil 的空切片只在确有需要时使用：JSON 序列化时 nil 切片输出 null、空切片输出 []。接口设计上不把 nil 切片与非 nil 空切片当作两种状态。

**为什么**

> When declaring an empty slice, prefer var t []string over t := []string{}. The former declares a nil slice value, while the latter is non-nil but zero-length.
>
> —— https://go.dev/wiki/CodeReviewComments#declaring-empty-slices

两者 len 与 cap 都是 0，功能等价但写法有两种；统一成 nil 后，读者不必停下来判断这里为什么用了另一种。

**正例**

```go
var ids []string
```

**反例**

```go
ids := []string{}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#declaring-empty-slices

**检测**：人工核对（需要输出 JSON 数组 [] 时属例外）

### 【SHOULD】DATA-002 结构体字面量用字段名初始化，不按位置罗列。

- 归属：编程规约/数据结构规约
- 起始版本：Go 1.0
- 主题：data-structures

字面量里逐字段写「字段名: 值」。字段少、顺序稳定的私有小结构体与测试表格数据可以例外；跨包结构体一律用字段名，因为字段增删或换序时位置写法会静默错位。

**为什么**

> You should almost always specify field names when initializing structs. This is now enforced by go vet.
>
> —— https://github.com/uber-go/guide/blob/master/style.md#use-field-names-to-initialize-structs

位置写法把字段顺序变成隐式契约，对方调整结构体后，类型相同的位置会静默错配且编译不报错；字段名写法让每个值对应到字段，改动当场暴露。`go vet` 的 composites 检查覆盖跨包位置字面量。

**正例**

```go
user := User{
	Name: "alice",
	Age:  30,
}
```

**反例**

```go
user := User{"alice", 30}
```

**依据**

- https://pkg.go.dev/cmd/vet
- https://github.com/uber-go/guide/blob/master/style.md#use-field-names-to-initialize-structs

**检测**：golangci-lint govet(composites)（govet 只覆盖跨包位置字面量，需启用 composites 检查）

### 【MAY】MODERN-001 切片与映射的常见操作优先使用 slices、maps 标准库，替代手写循环。

- 归属：编程规约/数据结构规约
- 起始版本：Go 1.21
- 主题：modernization

`slices` 与 `maps` 提供 `Contains`、`Sort`、`Clone`、`Keys` 等高频操作。能用标准库表达时不再手写循环，减少边界错误并统一写法。

**为什么**

> The new slices package provides many common operations on slices, using generic functions that work with slices of any element type.
>
> —— https://go.dev/doc/go1.21

手写循环容易在空切片、重复元素、排序稳定性上出错，标准库已覆盖这些边界。

**正例**

```go
if slices.Contains(ids, target) {
	return true
}
```

**反例**

```go
for _, id := range ids {
	if id == target {
		return true
	}
}
```

**依据**

- https://pkg.go.dev/slices
- https://go.dev/doc/go1.21

**检测**：golangci-lint modernize（go fix 也会提示等价的现代写法）

## （六）并发规约

### 【MUST】CONC-001 启动 goroutine 前明确退出时机，并由启动方等待结束。

- 归属：编程规约/并发规约
- 起始版本：Go 1.25
- 主题：concurrency

写 go 之前先想清三件事：它何时结束、谁等它结束、出错谁处理。启动方在同一个函数里用 `sync.WaitGroup` 或 `errgroup` 收束；子任务会返回错误时用 `errgroup`。不写没有等待方的 fire-and-forget。

**为什么**

> When you spawn goroutines, make it clear when - or whether - they exit. Goroutines can leak by blocking on channel sends or receives: the garbage collector will not terminate a goroutine even if the channels it is blocked on are unreachable.
>
> —— https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes

goroutine 阻塞时垃圾回收不会回收它，它连同栈上的引用一直占内存；没有等待方时，后台任务的错误与 panic 无人接手，只能在线上发作。

**正例**

```go
var wg sync.WaitGroup
for _, item := range items {
	wg.Go(func() {
		process(item)
	})
}
wg.Wait()
```

**反例**

```go
for _, item := range items {
	go process(item)
}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes
- https://pkg.go.dev/golang.org/x/sync/errgroup

**检测**：人工核对（核对每个 go 是否有对应的等待方与退出路径）

### 【MUST】CONC-002 channel 只由发送方关闭，接收方不关闭。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

关闭动作绑定在生产者的退出路径上，通常紧跟在所有发送完成之后。存在多个生产者时不直接关闭数据 channel，用单独的停止信号配合 context 收束。取消协作统一走 context，不用关闭 channel 传递取消。

**为什么**

> Sends on a closed channel panic, so it's important to ensure all sends are done before calling close. ... stages close their outbound channels when all the send operations are done.
>
> —— https://go.dev/blog/pipelines

向已关闭的 channel 发送会 panic，接收方关闭或重复关闭都会在运行期炸掉；把关闭权收到发送方，谁写谁负责收尾，边界才唯一。

**正例**

```go
ch := make(chan int)
go func() {
	defer close(ch)
	for _, v := range vals {
		ch <- v
	}
}()
```

**反例**

```go
ch := make(chan int)
go func() {
	for _, v := range vals {
		ch <- v
	}
}()
<-ch
close(ch)
```

**依据**

- https://go.dev/blog/pipelines

**检测**：人工核对（核对 close 是否都在发送方一侧）

### 【SHOULD】CONC-003 对外 API 保持同步语义，由调用方决定是否并发。

- 归属：编程规约/并发规约
- 起始版本：Go 1.7
- 主题：api-design, concurrency

内部可以起多个 goroutine 并发取数，但在返回前收束，对外返回普通结果与 error，不返回 channel 让调用方自己等待。

**为什么**

> Prefer synchronous functions - functions which return their results directly or finish any callbacks or channel ops before returning - over asynchronous ones. ... If callers need more concurrency, they can add it easily by calling the function from a separate goroutine. But it is quite difficult - sometimes impossible - to remove unnecessary concurrency at the caller side.
>
> —— https://go.dev/wiki/CodeReviewComments#synchronous-functions

返回 channel 的 API 把调度与生命周期管理推给每个调用方，调用方各写一遍等待与取消，容易漏掉；同步接口把并发收在实现内部，调用方按普通函数使用。

**正例**

```go
func Fetch(ctx context.Context) (Data, error) {
	// concurrency inside, joined before returning
}
```

**反例**

```go
func FetchAsync() <-chan Data {
	// the caller waits on its own
}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#synchronous-functions

**检测**：人工核对（核对对外函数是否返回 channel 让调用方等待）

### 【SHOULD】CONC-004 共享状态与保护它的锁放在同一结构体，仅经方法读写。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

被保护的字段不导出，读写都走加锁的方法，锁不暴露到结构体之外。选锁还是 channel 按场景：保护缓存、状态等共享数据用 `sync.Mutex`；传递数据所有权、分发任务、传递异步结果用 channel。

**为什么**

> Use whichever is most expressive and/or most simple. A common Go newbie mistake is to over-use channels and goroutines just because it's possible, and/or because it's fun. Don't be afraid to use a sync.Mutex if that fits your problem best.
>
> —— https://go.dev/wiki/MutexOrChannel

状态与锁分离时，调用方可以绕开方法直接改字段，锁形同虚设；放在一起后唯一的读写入口就是方法，保护范围可以逐行核对。

**正例**

```go
type counter struct {
	mu sync.Mutex
	n  int
}

func (c *counter) inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
```

**反例**

```go
type counter struct {
	mu sync.Mutex
	N  int // exported field, callers can bypass the lock and mutate it directly
}
```

**依据**

- https://go.dev/wiki/MutexOrChannel
- https://pkg.go.dev/sync

**检测**：人工核对（核对受锁保护的字段是否导出、是否只能经方法访问）

## （七）控制语句规约

### 【SHOULD】CTRL-001 先处理错误与边界并提前返回，正常路径保持零缩进。

- 归属：编程规约/控制语句规约
- 起始版本：Go 1.0
- 主题：control-flow

把 if 成功分支套 else 的写法改成先判断错误返回，正常逻辑跟在后面；带短变量声明的 if 在后续还要用变量时，把声明移到 if 之前。

**为什么**

> Try to keep the normal code path at a minimal indentation, and indent the error handling, dealing with it first. This improves the readability of the code by permitting visually scanning the normal path quickly.
>
> —— https://go.dev/wiki/CodeReviewComments#indent-error-flow

正常路径是读代码的人最常看的部分，套在 else 里要多缩进一层、多一次条件取反；先返回后，正常路径贴在最左侧，从上往下读一遍就到。

**正例**

```go
if err := save(); err != nil {
	return err
}
process()
```

**反例**

```go
if err := save(); err == nil {
	process()
} else {
	return err
}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#indent-error-flow

**检测**：人工核对（核对正常路径是否被 else 缩进）

### 【SHOULD】CTRL-002 不再显式拷贝循环变量。

- 归属：编程规约/控制语句规约
- 起始版本：Go 1.22
- 主题：modernization

Go 1.22 起 for 循环的迭代变量每次迭代各自独立，闭包与 goroutine 捕获到的就是当次的值，不需要再写 `v := v` 或把变量当参数传进 goroutine。本条对 go.mod 声明 Go 1.22 及以上的模块生效。

**为什么**

> Previously, the variables declared by a "for" loop were created once and updated by each iteration. In Go 1.22, each iteration of the loop creates new variables, to avoid accidental sharing bugs.
>
> —— https://go.dev/doc/go1.22

显式拷贝是 Go 1.22 之前规避共享变量的写法，语言修好后留着只会让读者以为这里还需要防护；交给 copyloopvar 检出即可统一清理。

**正例**

```go
fns := make([]func(), 0, len(items))
for _, item := range items {
	fns = append(fns, func() {
		process(item)
	})
}
```

**反例**

```go
fns := make([]func(), 0, len(items))
for _, item := range items {
	item := item
	fns = append(fns, func() {
		process(item)
	})
}
```

**依据**

- https://go.dev/doc/go1.22
- https://go.dev/blog/loopvar-preview

**检测**：golangci-lint copyloopvar

## （八）注释规约

### 【SHOULD】COMMENT-001 注释只写代码表达不出的东西，不复述代码在做什么。

- 归属：编程规约/注释规约
- 起始版本：Go 1.21
- 主题：comments

命名与结构能说清的，不写注释。只有这几类情况才值得写：

- 代码这样写的缘由，从代码本身看不出来
- 实现的是某个知名算法
- 套用了某个数学公式

注释解释为什么这么做，不复述代码做了什么。

**为什么**

> It is often better for comments to explain why something is done, not what the code is doing.
>
> —— https://google.github.io/styleguide/go/guide#clarity-rationale

代码本身就能读出它做了什么，再写一遍等于维护两份随时会分叉的说明，读者遇到还要判断哪份才对。把注释留给代码说不清的地方，注释才值得被信任。

**正例**

```go
// Make a copy: the caller may keep reusing the incoming slice.
entries := slices.Clone(in)
```

**反例**

```go
// copy the slice
entries := slices.Clone(in)
```

**依据**

- https://google.github.io/styleguide/go/guide#clarity-rationale

**检测**：人工核对（核对注释是否在复述代码本身能读出的内容）

### 【SHOULD】COMMENT-002 文档注释以被描述的名称开头，写成完整句子。

- 归属：编程规约/注释规约
- 起始版本：Go 1.0
- 主题：comments

需要写文档注释时（包注释、确有解释价值的导出标识符），以「名称 + 一句话」开头，结尾用句号。本条只管文档注释的写法，不要求给每个导出标识符补注释，写不写由 COMMENT-001 的取向决定。

**为什么**

> Comments documenting declarations should be full sentences, even if that seems a little redundant. This approach makes them format well when extracted into godoc documentation. Comments should begin with the name of the thing being described and end in a period.
>
> —— https://go.dev/wiki/CodeReviewComments#comment-sentences

注释以名称开头，godoc 的列表页里摘要是完整句子，名称链接能正确识别；不写名称时摘要读起来像半句话，链接也容易漏。

**正例**

```go
// Load reads and parses the clause data file.
func Load(path string) (*Manual, error) {
```

**反例**

```go
// read the clause data
func Load(path string) (*Manual, error) {
```

**依据**

- https://go.dev/wiki/CodeReviewComments#comment-sentences
- https://go.dev/doc/effective_go#commentary

**检测**：人工核对（核对文档注释是否以名称开头并写成完整句子）

