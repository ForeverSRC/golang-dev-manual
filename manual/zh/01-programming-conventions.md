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

### 【SHOULD】NAMING-007 布尔谓词方法用 Is、Has、Can 等前缀命名。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

返回 bool 的方法按语义选前缀，让名字读成真假问句。

- 状态判定用 Is：IsHealthy、IsEmpty
- 持有判定用 Has：HasExpired、HasPrefix
- 能力判定用 Can：CanRetry、CanClose
- 标准库既有的非前缀谓词（如 `strings.Contains`）按原惯例保留

**为什么**

> it's neither idiomatic nor necessary to put Get into the getter's name
>
> —— https://go.dev/doc/effective_go#Getters

标准库对命名方法的立场是让名字直接表达语义，谓词类方法同理：IsConnected 把调用处读成一个真假问句，而裸形容词 Connected 分不清「已连接」与「连接」，读者必须回查类型才能确认它是否为布尔值。

**正例**

```go
func (s *Server) IsHealthy() bool { return s.err == nil }
```

**反例**

```go
func (s *Server) Healthy() bool { return s.err == nil }
```

**依据**

- https://go.dev/doc/effective_go#Getters

**检测**：人工核对（核对返回 bool 的方法是否带谓词前缀）

### 【SHOULD】NAMING-008 布尔变量与结构体字段用 Is、Has、Can 前缀命名。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

bool 类型的变量与字段用系动词前缀表达判断语义。

- IsConnected、HasCache、CanRetry
- 不用裸形容词或名词：Connected、Cache、Retry
- 临时布尔变量同样适用：isValid、hasMore

**为什么**

> it's neither idiomatic nor necessary to put Get into the getter's name
>
> —— https://go.dev/doc/effective_go#Getters

裸形容词与对象名混淆：cfg.Cache 既可能是「缓存配置」也可能是「启用缓存」的布尔值；加 Is/Has/Can 后类型语义在名字里自明，读 if 条件时不再需要跳转到定义确认类型。

**正例**

```go
type Config struct { IsVerbose bool }
```

**反例**

```go
type Config struct { Verbose bool }
```

**依据**

- https://go.dev/doc/effective_go#Getters

**检测**：人工核对（核对 bool 类型字段与变量的命名前缀）

### 【SHOULD】NAMING-009 单方法接口名用方法名加 er 后缀。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

单方法接口按「方法名 + er」构成行为者名词。

- Reader、Writer、Closer、Stringer、CloseNotifier
- 多方法接口才用具象领域名词：Conn、Pool、Store
- 方法名不便加 er 时允许相近变体，如 Recognize 对应 Recognizer

**为什么**

> By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier.
>
> —— https://go.dev/wiki/CodeReviewComments#interface-names

er 后缀让接口名与它唯一的方法形成直接对应，读者看到名字即知契约；标准库全部遵循该惯例，偏离它会让同族接口出现两套命名风格。

**正例**

```go
type Reader interface { Read(p []byte) (n int, err error) }
```

**反例**

```go
type DataFetcher interface { Read(p []byte) (n int, err error) }
```

**依据**

- https://go.dev/wiki/CodeReviewComments#interface-names

**检测**：人工核对（核对单方法接口类型名与方法名的对应）

### 【SHOULD】NAMING-010 变体函数用 Must、f、WithContext 等固定标记命名。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

同一功能的变体函数靠固定前后缀区分行为差异。

- Must 前缀：失败即 panic，如 `regexp.MustCompile`
- f 后缀：接受格式串与可变参数，如 `errors.Errorf`
- WithContext 后缀：携带 context 的版本，如 `http.NewRequestWithContext`
- 命名标记与标准库保持一致，不自造变体

**为什么**

> MustCompile is like Compile but panics if the expression cannot be parsed. It simplifies safe initialization of global variables holding compiled regular expressions.
>
> —— https://pkg.go.dev/regexp#MustCompile

调用方在 code review 与 IDE 提示里往往只看函数名不看签名，Must 前缀告知「这里会 panic」、f 后缀告知「参数是格式串」、WithContext 告知「取消随 ctx 传播」；变体行为若不进名字，误用只能靠读实现发现。

**正例**

```go
var re = regexp.MustCompile(`^v\d+$`)
```

**反例**

```go
func MustParseTime(s string) time.Time { t, _ := time.Parse(time.RFC3339, s); return t }
```

**依据**

- https://pkg.go.dev/regexp#MustCompile

**检测**：人工核对（核对同族变体函数的命名标记是否与标准库一致）

### 【SHOULD】NAMING-011 标识符不遮蔽 len、cap、error、min 等预声明名。

- 归属：编程规约/命名规约
- 起始版本：Go 1.0
- 主题：naming

局部变量、参数与字段不使用预声明标识符的名字。

- 预声明名包括 len、cap、new、make、min、max、clear、error、any 等
- 误用被遮蔽的名字不产生编译错误，语义静默改变
- 确需同名时（如代码生成）需在评审中说明

**为什么**

> The following identifiers are implicitly declared in the universe block:
>
> —— https://go.dev/ref/spec#Predeclared_identifiers

预声明标识符属于所有作用域的兜底名字，局部遮蔽后 min(x, y) 变成调用用户函数、error 变成局部变量，这类覆盖没有编译提示，只表现为行为异常；predeclared 检查器把它变成可机检项。

**正例**

```go
func dedup(items []int) []int { maxSeen := 0; _ = maxSeen; return items }
```

**反例**

```go
func dedup(items []int) []int { max := 0; _ = max; return items }
```

**依据**

- https://go.dev/ref/spec#Predeclared_identifiers

**检测**：golangci-lint predeclared

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

### 【SHOULD】CONST-004 iota 枚举零值留 Unknown 哨兵或从 iota+1 起步。

- 归属：编程规约/常量与枚举规约
- 起始版本：Go 1.0
- 主题：constants, naming

枚举类型的零值不承载业务含义。

- 首个成员定义 StatusUnknown，或首个业务成员从 iota+1 起步
- 零值作为「未初始化/无效」哨兵，序列化与校验时把它与业务状态区分开
- 例外：零值有明确业务语义的类型（如 GenderMale 为 0）可保留，需在评审中确认

**为什么**

> Within a const declaration a name not explicitly initialized is implicitly initialized to iota, which is the line number of the corresponding ConstSpec.
>
> —— https://go.dev/ref/spec#Iota

var s Status 的零值静默等于首个 const 成员，未初始化的字段被当成第一个业务状态通过校验，排障时无法区分「默认值」与「真实状态」；把零值留给 Unknown 后，未初始化的状态在校验处显式暴露。

**正例**

```go
type Status int
const (
  StatusUnknown Status = iota
  StatusActive
  StatusClosed
)
```

**反例**

```go
type Status int
const (
  StatusActive Status = iota
  StatusClosed
)
```

**依据**

- https://go.dev/ref/spec#Iota

**检测**：人工核对（核对 iota 枚举首个成员是否占用零值）

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

### 【MUST】FUNC-004 返回接口类型的 nil 分支显式 return nil。

- 归属：编程规约/函数与方法规约
- 起始版本：Go 1.0
- 主题：error-handling

返回值类型是接口（含 error）时，不得返回有类型的 nil 指针变量。

- 确无错误直接 return nil
- 需要返回具体错误类型时，构造后判空再决定返回值
- 适用一切接口返回值，error 是最高频场景

**为什么**

> If we store a nil pointer of type *int inside an interface value, the inner type will be *int regardless of the value of the pointer: (*int)(nil). Thus, the interface value will be non-nil even when the pointer value is inside.
>
> —— https://go.dev/doc/faq#nil_error

接口值携带 (类型, 值) 二元组，typed nil 指针装入后接口非 nil，调用方 err != nil 判定成立、进入错误分支，解引用 nil 字段时 panic；显式 return nil 从根上避免这个判定陷阱。

**正例**

```go
func validate(c Config) error {
  if c.Host == "" {
    return ErrMissingHost
  }
  return nil
}
```

**反例**

```go
func validate(c Config) error {
  var vErr *ValidationError
  if c.Host == "" {
    vErr = &ValidationError{Field: "host"}
  }
  return vErr
}
```

**依据**

- https://go.dev/doc/faq#nil_error

**检测**：人工核对（核对返回接口类型的函数中 nil 分支的写法）

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

### 【SHOULD】FUNC-003 入参超过 4 个或返回值超过 3 个时，改用结构体收束。

- 归属：编程规约/函数与方法规约
- 起始版本：Go 1.0
- 主题：api-design

入参 5 个及以上、或返回值 4 个及以上时，把同向的数据收进一个结构体，调用方按字段名传值。计数口径与收束方式如下：

- 计数口径：入参按参数名计数，含 `ctx`，不含方法接收者；返回值按结果名计数，含 `error`。
- 入参收束：业务参数收进一个入参结构体，`ctx` 仍留在签名第一位，不进结构体，与 `CTX-001` 一致。
- 返回值收束：结果收进一个结果结构体，`error` 留在最后一个返回值位置，不进结构体。
- 例外：签名由外部约定固定的位置不改，例如接口实现、`go:generate` 生成物、mock 生成物。
- 与 `FUNC-002` 的分工：返回值含义不自明时用命名返回值，条数超限才收进结构体。

**为什么**

> Warns when a function receives more parameters than the maximum set by the rule's configuration. ... Specifies the maximum number of results a function can return. Functions returning too many results can be hard to understand/use.
>
> —— https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md

`quote` 里的两条规则给的是可配置上限，Go 官方只说要控制参数数量，没有给数字。带数字的口径都来自工具与社区：revive 的 `argument-limit` 默认 8、`function-result-limit` 默认 3，go-critic 的 `tooManyResultsChecker` 默认 5 且默认关闭。本手册取入参 4、返回值 3 作为自定上限。依据是调用点的可读性：位置参数到第 5 个，读的人已经分不清哪个值对应哪个含义，只能回翻签名逐个对位；返回值到第 4 个同样要靠数位置。收进结构体后每个值都带字段名，调用点自解释，日后增减字段也不会静默错位。

**正例**

```go
func (s *Service) List(ctx context.Context, q ListQuery) ([]ClauseView, error) {
	return s.repo.Find(ctx, q)
}
```

**反例**

```go
func (s *Service) List(ctx context.Context, lang string, levels []string, category string, tags []string) ([]ClauseView, error) {
	return s.repo.Find(ctx, lang, levels, category, tags)
}
```

**依据**

- https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md
- https://go-critic.com/overview.html#toomanyresultschecker

**检测**：golangci-lint revive(argument-limit, function-result-limit)

## （五）数据结构规约

### 【MUST】DATA-003 类型断言用 comma-ok 形式，不用单值形式。

- 归属：编程规约/数据结构规约
- 起始版本：Go 1.0
- 主题：data-structures

断言结果可能不是期望类型时必须取 ok 值分支处理。

- 用 `v, ok := x.(T)`，按 ok 分支处理
- 单值形式 `v := x.(T)` 仅限断言失败即属程序错误的场景，且需注释说明
- 类型 switch 不受本条约束

**为什么**

> If it does not hold, a run-time panic occurs.
>
> —— https://go.dev/ref/spec#Type_assertions

单值断言在类型不匹配时直接 panic，把「值不是这个类型」这一可预期分支变成进程崩溃；comma-ok 把判定权留给代码，调用方可以走降级或报错路径。

**正例**

```go
v, ok := x.(string)
if !ok {
  return fmt.Errorf("unexpected type %T", x)
}
```

**反例**

```go
v := x.(string)
```

**依据**

- https://go.dev/ref/spec#Type_assertions

**检测**：golangci-lint forcetypeassert

### 【MUST】DATA-004 nil map 初始化后再写入。

- 归属：编程规约/数据结构规约
- 起始版本：Go 1.0
- 主题：data-structures

写入 map 前确认它已通过 make 或字面量初始化。

- 读取 nil map 合法，返回零值
- 写入 nil map 直接 panic
- 结构体字段的 map 可在构造函数或方法内惰性初始化，使零值可用

**为什么**

> A nil map is equivalent to an empty map except no elements may be added.
>
> —— https://go.dev/blog/maps

`var m map[string]int` 的零值可读不可写，写入路径一旦到达就 panic，且编译期无任何提示；先初始化或让类型的零值可用，二选一写清楚。

**正例**

```go
m := make(map[string]int)
m["k"] = 1
```

**反例**

```go
var m map[string]int
m["k"] = 1
```

**依据**

- https://go.dev/blog/maps

**检测**：人工核对（核对 map 的写入路径是否发生在已初始化的 map 上）

### 【MUST】DATA-005 不依赖 map 遍历顺序，需确定序先排序 key。

- 归属：编程规约/数据结构规约
- 起始版本：Go 1.0
- 主题：data-structures

range map 的输出顺序不构成任何契约。

- 需要稳定输出的序列化、比对、展示场景先对 key 排序
- 依赖顺序的逻辑在 CI 与生产都可能间歇性失败
- Go 运行时对 map 遍历做随机化，属语言保证而非实现细节

**为什么**

> The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next.
>
> —— https://go.dev/ref/spec#For_range

语言规范明确不保证遍历顺序，依赖顺序的代码在同一进程两次运行结果都可能不同，测试通过只是随机性恰好 favor 它；需要确定性时显式排序才构成契约。

**正例**

```go
keys := slices.Sorted(maps.Keys(m))
for _, k := range keys {
  fmt.Println(k, m[k])
}
```

**反例**

```go
for k, v := range m {
  fmt.Println(k, v)
}
```

**依据**

- https://go.dev/ref/spec#For_range

**检测**：人工核对（核对 range map 之后的输出与比对逻辑是否依赖顺序）

### 【MUST】DATA-007 依赖 JSON 序列化的结构体字段必须导出。

- 归属：编程规约/数据结构规约
- 起始版本：Go 1.0
- 主题：data-structures

需要出现在 JSON 中的字段用导出标识符，未导出字段一律静默忽略。

- 未导出字段无法被 `encoding/json` 编码与解码
- 携带 key 映射时配 field tag（见 STRUCT-004）
- 解析到 `interface{}` 的数字一律是 `float64`，大整数经 `json.Number` 或类型化结构体承接

**为什么**

> Unmarshal matches incoming object keys to the keys used by Marshal (either the struct field name or its tag), preferring an exact match but also accepting a case-insensitive match.
>
> —— https://pkg.go.dev/encoding/json#Unmarshal

未导出字段对 encoding/json 不可访问，线上 JSON 悄悄缺字段而本地测试若从未校验契约就发现不了；导出并打 tag 让契约显式落在结构体定义处。

**正例**

```go
type User struct {
  Name string `json:"name"`
}
```

**反例**

```go
type User struct {
  name string `json:"name"`
}
```

**依据**

- https://pkg.go.dev/encoding/json#Unmarshal

**检测**：人工核对（对照结构体字段与序列化契约，核对未导出字段）

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

### 【SHOULD】DATA-006 窄化整型转换前检查数值范围。

- 归属：编程规约/数据结构规约
- 起始版本：Go 1.0
- 主题：data-structures, security

int64 到 int32、uint 到 int 等缩窄转换需显式校验。

- 校验值落在目标类型可表示范围内后再转换
- 转换值来自外部输入时必查（`gosec G115` 覆盖的场景）
- 高位信息明确无意义的位运算场景可豁免并注释

**为什么**

> When converting between integer types, if the value is a signed integer, it is sign extended to infinite precision and then truncated to fit in the result type's size; if the value is an unsigned integer, it is zero-extended and then truncated to fit in the result type's size.
>
> —— https://go.dev/ref/spec#Conversions_to_and_from_integer_string_types

缩窄转换静默截断回绕，30 亿的 int64 转 int32 变成负数且无任何报错；外部输入驱动这类转换时，截断结果直接变成业务错误或安全漏洞。

**正例**

```go
if v > math.MaxInt32 {
  return fmt.Errorf("value %d overflows int32", v)
}
n := int32(v)
```

**反例**

```go
n := int32(v)
```

**依据**

- https://go.dev/ref/spec#Conversions_to_and_from_integer_string_types

**检测**：golangci-lint gosec(G115)

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

## （六）标准库使用规约

### 【MUST】STD-001 time.Duration 不得与 time.Duration 相乘。

- 归属：编程规约/标准库使用规约
- 起始版本：Go 1.0
- 主题：data-structures

Duration 参与运算时只允许与标量相乘。

- 用 `d * n`（n 为整数）表达倍数
- `d * d` 或 `d * (x * time.Second)` 都是量纲错误
- 由变量构造 Duration 用 `time.Duration(n) * time.Second` 这类标量乘法

**为什么**

> A Duration represents the elapsed time between two instants as an int64 nanosecond count.
>
> —— https://pkg.go.dev/time#Duration

`time.Duration` 的底层是纳秒数，两个 Duration 相乘得到「纳秒的平方」，与直觉相差多个数量级：`5 * time.Second * time.Second` 不是 25 秒而是 2.5 万亿纳秒；durationcheck 检查器把它变成可机检项。

**正例**

```go
timeout := 5 * time.Second
result := timeout * attempts
```

**反例**

```go
timeout := 5 * time.Second * time.Second
```

**依据**

- https://pkg.go.dev/time#Duration

**检测**：golangci-lint durationcheck

### 【SHOULD】STD-002 time.Time 比较用 Equal，不用 ==。

- 归属：编程规约/标准库使用规约
- 起始版本：Go 1.0
- 主题：data-structures

判断同一时刻用 Equal 或 Compare。

- `==` 同时比较墙钟、时区与单调时钟读数
- 解析自字符串、经网络传输的 time.Time 不带单调时钟，但混用来源时 `==` 不可靠
- 零值判定用 `IsZero`，先后比较用 `Before`/`After`

**为什么**

> Equal reports whether t and u represent the same time instant. Two times can be equal even if they are in different locations.
>
> —— https://pkg.go.dev/time#Time.Equal

time.Time 内部携带单调时钟读数，== 把读数差异也算进相等判定，同一墙钟时刻可能判为不等；Equal 只比较时间瞬间，语义与业务预期一致。

**正例**

```go
if t1.Equal(t2) {
  return true
}
```

**反例**

```go
if t1 == t2 {
  return true
}
```

**依据**

- https://pkg.go.dev/time#Time.Equal

**检测**：人工核对（核对 time.Time 的比较表达式）

### 【SHOULD】STD-003 去前后缀用 strings.TrimPrefix 与 TrimSuffix。

- 归属：编程规约/标准库使用规约
- 起始版本：Go 1.0
- 主题：data-structures

按子串裁剪时不用 strings.Trim 系列函数。

- `strings.TrimPrefix`/`strings.TrimSuffix` 精确去掉一次子串：`strings.TrimPrefix(s, "api/")`
- `strings.Trim`/`TrimLeft`/`TrimRight` 的参数是字符集合，会反复去掉两端所有命中字符
- 混用两者是常见过度裁剪来源

**为什么**

> Trim returns a slice of the string s with all leading and trailing Unicode code points contained in cutset removed.
>
> —— https://pkg.go.dev/strings#Trim

Trim 的第二个参数是字符集合而非子串：strings.Trim(s, "ab") 会去掉两端所有 a 和 b 字符，与「去掉前缀 ab」的意图完全不同；字符集语义藏在函数名背后，误用只表现为偶发裁剪过头。

**正例**

```go
path := strings.TrimPrefix(raw, "/api/")
```

**反例**

```go
path := strings.Trim(raw, "/api/")
```

**依据**

- https://pkg.go.dev/strings#Trim

**检测**：人工核对（核对字符串裁剪意图与所用 API 是否匹配）

## （七）数据库访问规约

### 【MUST】DB-001 sql.Rows 迭代结束检查 rows.Err 并 Close。

- 归属：编程规约/数据库访问规约
- 起始版本：Go 1.0
- 主题：data-structures, error-handling

rows 的生命周期跟随整个迭代。

- `for rows.Scan` 循环结束后检查 `rows.Err()`
- rows 用完即 `Close`，或 `defer rows.Close()` 紧跟查询
- 不返回行的语句（INSERT/UPDATE/DELETE）用 `Exec`，不用 `Query`

**为什么**

> Err returns the error, if any, that was encountered during iteration. Err should be checked explicitly after the iteration is done.
>
> —— https://pkg.go.dev/database/sql#Rows.Err

漏 Close 会让连接一直被占用，不归还连接池；漏 Err 会吞掉迭代中途的 I/O 中断（网络抖动时循环悄悄少读行）；Err 的错误不会从 Scan 冒出来，只能显式检查。

**正例**

```go
rows, err := db.QueryContext(ctx, q)
if err != nil {
  return err
}
defer rows.Close()
for rows.Next() {
  if err := rows.Scan(&id); err != nil {
    return err
  }
}
return rows.Err()
```

**反例**

```go
rows, err := db.QueryContext(ctx, q)
if err != nil {
  return err
}
for rows.Next() {
  if err := rows.Scan(&id); err != nil {
    return err
  }
}
return nil
```

**依据**

- https://pkg.go.dev/database/sql#Rows.Err

**检测**：golangci-lint rowserrcheck, sqlclosecheck

### 【MUST】DB-002 多语句写操作包事务，Begin 后立即 defer Rollback。

- 归属：编程规约/数据库访问规约
- 起始版本：Go 1.0
- 主题：data-structures, error-handling

事务必须有统一的收尾路径。

- `tx := db.BeginTx(ctx, nil)` 后紧跟 `defer tx.Rollback()`
- 提交成功后 Rollback 返回 `ErrTxDone`，无害
- panic 或提前 return 的路径不会悬挂事务

**为什么**

> After calling Commit or Rollback, all operations on the transaction fail with ErrTxDone.
>
> —— https://pkg.go.dev/database/sql#Tx

没有 defer Rollback 的事务在 panic 或提前 return 时悬挂到连接回收，中间的行锁持续占用；Rollback 在 Commit 之后调用只返回 ErrTxDone，作为兜底总是安全。

**正例**

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
  return err
}
defer tx.Rollback()
if _, err := tx.ExecContext(ctx, q1); err != nil {
  return err
}
return tx.Commit()
```

**反例**

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
  return err
}
if _, err := tx.ExecContext(ctx, q1); err != nil {
  return err
}
return tx.Commit()
```

**依据**

- https://pkg.go.dev/database/sql#Tx

**检测**：人工核对（核对事务开启后是否紧跟 defer Rollback）

### 【MUST】DB-004 sql.ErrNoRows 用 errors.Is 判别。

- 归属：编程规约/数据库访问规约
- 起始版本：Go 1.0
- 主题：error-handling

「没查到」是领域结果不是普通错误。

- 用 `errors.Is(err, sql.ErrNoRows)` 判别后转成领域语义（如返回不存在标记）
- 直接 `==` 比较在错误被包装后失效
- `QueryRow` 的单行查询是 ErrNoRows 的主要来源

**为什么**

> ErrNoRows is returned by Scan when QueryRow doesn't return a row. In this case, QueryRow will return this error.
>
> —— https://pkg.go.dev/database/sql#ErrNoRows

把 ErrNoRows 当普通错误上抛，上游无法区分「查无此人」与「数据库故障」，往往表现为一次不存在的请求触发 5xx；errors.Is 判别后转成「不存在」结果，错误链语义才完整。

**正例**

```go
err := db.QueryRowContext(ctx, q, id).Scan(&name)
if errors.Is(err, sql.ErrNoRows) {
  return ErrNotFound
}
if err != nil {
  return err
}
```

**反例**

```go
err := db.QueryRowContext(ctx, q, id).Scan(&name)
if err == sql.ErrNoRows {
  return ErrNotFound
}
if err != nil {
  return err
}
```

**依据**

- https://pkg.go.dev/database/sql#ErrNoRows

**检测**：人工核对（核对 ErrNoRows 的判定方式）

### 【SHOULD】DB-003 可空列用指针或 sql.Null 承接。

- 归属：编程规约/数据库访问规约
- 起始版本：Go 1.22
- 主题：data-structures

NULL 与 Go 零值是不同语义。

- 列可空时 Scan 目标用 `*T` 或 `sql.Null[T]`（Go 1.22+ 泛型）
- NULL 扫进裸基础类型直接报错
- 表结构变化（新增可空列）时同步检查 Scan 目标

**为什么**

> NullString represents a string that may be null. NullString provides a convenient way of dealing with NULL string values.
>
> —— https://pkg.go.dev/database/sql#NullString

裸类型 Scan 到 NULL 直接报 conversion 错误，且报错位置在运行期数据边界，靠下游 DDL 变更触发；指针与 sql.Null[T] 把「可空」写进类型，承接语义与数据库一致。

**正例**

```go
var nick *string
if err := rows.Scan(&id, &nick); err != nil {
  return err
}
```

**反例**

```go
var nick string
if err := rows.Scan(&id, &nick); err != nil {
  return err
}
```

**依据**

- https://pkg.go.dev/database/sql#NullString

**检测**：人工核对（对照 DDL 的可空列与 Scan 目标类型）

## （八）并发规约

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

### 【MUST】CONC-010 wg 并发等待用 wg.Go，手工 Add 必须在 go 语句之前。

- 归属：编程规约/并发规约
- 起始版本：Go 1.25
- 主题：concurrency

Go 1.25+ 优先用 wg.Go(f) 替代手工 Add/go/Done 配对。

- `wg.Go(f)` 把 Add、goroutine 启动与 Done 封装成一步，消除配对错误
- 手工使用时 `wg.Add(1)` 必须写在 go 语句之前，写在协程内会与 Wait 竞争
- 需要收集每个任务的 error 时用 `errgroup`；Go 1.25 以下仍手工配对

**为什么**

> Go calls the given function in a new goroutine. It blocks until all goroutines have exited.
>
> —— https://pkg.go.dev/sync#WaitGroup.Go

Add 写在协程内时，主 goroutine 的 Wait 可能在 Add 执行前就通过计数为零的检查提前返回，等待形同虚设；wg.Go 由标准库保证配对时序，从根上消除这类间歇性竞态。

**正例**

```go
var wg sync.WaitGroup
for _, t := range tasks {
  wg.Go(func() { handle(t) })
}
wg.Wait()
```

**反例**

```go
var wg sync.WaitGroup
for _, t := range tasks {
  go func() {
    wg.Add(1)
    defer wg.Done()
    handle(t)
  }()
}
wg.Wait()
```

**依据**

- https://pkg.go.dev/sync#WaitGroup.Go
- https://go.dev/doc/go1.25

**检测**：人工核对（核对 wg.Add 与 go 语句的相对位置）

### 【MUST】CONC-012 持有 RLock 时不请求同一 RWMutex 的 Lock。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

RWMutex 不可升级，读写锁请求次序要单向。

- 持 RLock 期间需要写操作时，先 `RUnlock` 再 `Lock`
- 需要读改写原子性时直接用 `Lock`
- 递归读锁同样禁止

**为什么**

> If a goroutine holds a RWMutex for reading and another goroutine might call Lock, no goroutine should expect to be able to acquire a read lock until the initial read lock is released.
>
> —— https://pkg.go.dev/sync#RWMutex

RWMutex 的写锁请求会阻塞后续读锁，持读锁再请求写锁形成自等待死锁：自己等自己释放读锁；文档明确把递归读锁一并禁止，升级语义必须拆成两步或改用写锁。

**正例**

```go
mu.RLock()
v := m[key]
mu.RUnlock()
if needUpdate(v) {
  mu.Lock()
  m[key] = compute()
  mu.Unlock()
}
```

**反例**

```go
mu.RLock()
if needUpdate(m[key]) {
  mu.Lock()
  m[key] = compute()
  mu.Unlock()
}
mu.RUnlock()
```

**依据**

- https://pkg.go.dev/sync#RWMutex

**检测**：人工核对（核对 RLock 与 Lock 的嵌套持有）

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

### 【SHOULD】CONC-005 函数签名中的 channel 参数标注单向类型。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

按参数实际用途选用 chan<- 或 <-chan。

- 只发送的参数用 `chan<- T`，只接收的用 `<-chan T`
- 单向约束在编译期拦截误发送与误关闭
- 返回值与结构体字段按同样口径

**为什么**

> The optional <- operator specifies the channel direction, send or receive. If a direction is given, the channel is directional, otherwise it is bidirectional.
>
> —— https://go.dev/ref/spec#Channel_types

双向 chan T 对误发送与误关闭没有任何编译期约束，误关闭生产者的 channel 会让消费方 panic；单向类型把使用边界写进签名，误用直接编译失败。

**正例**

```go
func consume(in <-chan int) { for v := range in { process(v) } }
```

**反例**

```go
func consume(in chan int) { for v := range in { process(v) } }
```

**依据**

- https://go.dev/ref/spec#Channel_types

**检测**：人工核对（核对公开函数的 channel 参数方向）

### 【SHOULD】CONC-006 缓冲 channel 容量大于 1 时注释容量依据与写满行为。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

缓冲容量是并发设计决策，不是随意数字。

- 默认无缓冲；缓冲大于 1 须注释容量依据与写满时的行为（阻塞、丢弃还是降级）
- 大缓冲掩盖背压问题，堆积的任务无界增长
- 生产者消费者速率对齐的设计说明写在 make 语句旁

**为什么**

> Channels should usually be size one or give an unbuffered implementation. A channel size of zero or one is sufficient in most scenarios.
>
> —— https://github.com/uber-go/guide/blob/master/style.md#channel-size

无缓冲 channel 天然提供背压，写满的缓冲 channel 把「生产过快」的信号推迟到堆积之后；没有注释的容量值无法区分「有意为之」与「随手写的」，评审无从判断写满后的行为是否安全。

**正例**

```go
// 容量 64 对齐下游批处理批次；写满时生产端阻塞，保证背压
jobs := make(chan Job, 64)
```

**反例**

```go
jobs := make(chan Job, 1000)
```

**依据**

- https://github.com/uber-go/guide/blob/master/style.md#channel-size

**检测**：grep 正则 make\(chan[^)]*,\s*([^)1][^)]*)\)

### 【SHOULD】CONC-007 可取消路径上的 select 必须包含 ctx.Done() 分支。

- 归属：编程规约/并发规约
- 起始版本：Go 1.7
- 主题：concurrency, context

长驻 goroutine 的 select 循环要有取消出口。

- select 的 case 集合里放 `ctx.Done()`，收到取消后退出
- 只有一次性、确定会结束的 select 可省略
- 缺失取消分支的 goroutine 在服务停止时无法退出

**为什么**

> Done returns a channel that's closed when work done on behalf of this context should be canceled. Done is provided for use in select statements
>
> —— https://pkg.go.dev/context#Context.Done

select 阻塞在多个 channel 上，取消信号若不在 case 集合里，调用方 cancel 后该 goroutine 仍然阻塞，服务优雅停机时这些 goroutine 全部泄漏。

**正例**

```go
for {
  select {
  case <-ctx.Done():
    return
  case v := <-in:
    process(v)
  }
}
```

**反例**

```go
for {
  select {
  case v := <-in:
    process(v)
  }
}
```

**依据**

- https://pkg.go.dev/context#Context.Done

**检测**：人工核对（核对长驻 goroutine 的 select 分支集合）

### 【SHOULD】CONC-008 for 循环内的 select 不得用 default 制造忙等。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

循环内 select 的 default 分支会让循环空转。

- 需要等待数据时去掉 `default`，让 select 阻塞
- 需要非阻塞轮询时在 `default` 分支 sleep 或退避
- 接收方发现 channel 关闭后将其置 nil，退出多路复用

**为什么**

> Execution of a "select" statement proceeds in several steps: ... If multiple cases can proceed, a uniform pseudo-random choice is made
>
> —— https://go.dev/ref/spec#Select_statements

`for { select { default: } }` 在没有 case 就绪时立即重试，这个空转循环会一直占用一个 CPU 核；被关闭的 channel 永远就绪，混进 select 后同样占用 CPU，把 channel 置 nil 是标准的退出方法。

**正例**

```go
for {
  select {
  case <-ctx.Done():
    return
  case v, ok := <-in:
    if !ok {
      in = nil
      continue
    }
    process(v)
  }
}
```

**反例**

```go
for {
  select {
  case v := <-in:
    process(v)
  default:
  }
}
```

**依据**

- https://go.dev/ref/spec#Select_statements

**检测**：人工核对（核对循环内 select 的 default 分支）

### 【SHOULD】CONC-009 循环内周期触发用 time.NewTicker 加 Stop。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

需要停止的周期任务不用 time.After 与 time.Tick。

- 循环内每次 `time.After` 新建定时器，旧计时器触发前不可回收
- `time.Tick` 的 ticker 无法停止，永久泄漏
- 用 `ticker := time.NewTicker(d)` 加 `defer ticker.Stop()`

**为什么**

> While Tick is useful for clients that have no need to shut down the Ticker, be aware that without a way to shut it down the underlying Ticker can never be recovered by the garbage collector; it "leaks".
>
> —— https://pkg.go.dev/time#Tick

time.After 在循环里每轮分配一个新定时器，触发前 GC 无法回收，长循环表现为内存与定时器堆积；官方文档对 time.Tick 直接标注 leak，需要停止语义时只有 NewTicker + Stop 一条路。

**正例**

```go
ticker := time.NewTicker(interval)
defer ticker.Stop()
for {
  select {
  case <-ticker.C:
    tick()
  case <-ctx.Done():
    return
  }
}
```

**反例**

```go
for {
  select {
  case <-time.After(interval):
    tick()
  case <-ctx.Done():
    return
  }
}
```

**依据**

- https://pkg.go.dev/time#Tick

**检测**：人工核对（核对循环体内的 time.After 与 time.Tick）

### 【SHOULD】CONC-011 临界区内不执行 I/O、网络调用与 channel 操作。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

锁的持有时间只覆盖内存操作。

- 网络、磁盘、数据库等阻塞调用放在锁外
- 需要锁内数据做 I/O 时，先拷贝出结果再释放锁
- channel 收发可能无限期阻塞，同样移出临界区

**为什么**

> Programs that modify data being simultaneously accessed by multiple goroutines must serialize such access.
>
> —— https://go.dev/ref/mem

锁内做网络调用时持锁时长等于网络往返，毫秒级 I/O 把其他所有等待该锁的 goroutine 串行化；把 I/O 移出临界区后，持锁时长回到纳秒级，并发度不随外部延迟衰减。

**正例**

```go
resp, err := http.Get(url)
if err != nil {
  return err
}
mu.Lock()
cache[key] = resp
mu.Unlock()
```

**反例**

```go
mu.Lock()
resp, err := http.Get(url)
if err != nil {
  mu.Unlock()
  return err
}
cache[key] = resp
mu.Unlock()
```

**依据**

- https://go.dev/ref/mem

**检测**：人工核对（核对 Lock 与 Unlock 之间的阻塞调用）

### 【SHOULD】CONC-013 简单计数与标志位用类型化原子类型。

- 归属：编程规约/并发规约
- 起始版本：Go 1.19
- 主题：concurrency

atomic.Int64、atomic.Bool 等类型化原子替代函数式原子调用。

- 计数、开关、版本号等单值状态用 `atomic.Int64`/`atomic.Int32`/`atomic.Bool`/`atomic.Pointer`
- 不用裸 `atomic.AddInt64(&counter, 1)` 的函数式写法
- 复合状态仍用互斥锁保护

**为什么**

> Int64 is an atomically operated int64 value.
>
> —— https://pkg.go.dev/sync/atomic#Int64

类型化原子把「这个字段是原子访问的」写进类型本身，误传普通 int64 的编译期即可发现；函数式原子调用的目标字段与普通字段在类型上无区别，新增访问点时容易漏用原子操作。

**正例**

```go
var counter atomic.Int64
counter.Add(1)
```

**反例**

```go
var counter int64
atomic.AddInt64(&counter, 1)
```

**依据**

- https://pkg.go.dev/sync/atomic#Int64

**检测**：人工核对（核对原子操作的使用形式）

### 【SHOULD】CONC-014 写入密集场景不用 sync.Map。

- 归属：编程规约/并发规约
- 起始版本：Go 1.9
- 主题：concurrency, data-structures

sync.Map 只服务两类访问模式。

- 键一次写入多次读取（只增的缓存）
- 多个 goroutine 操作不相交的键集合
- 其余场景用 `map` 加 `RWMutex`，写入密集时 `sync.Map` 更慢且类型不安全

**为什么**

> The Map type is optimized for two common use cases: (1) when the entry for a given key is only ever written once but read many times, as in caches that only grow, or (2) when multiple goroutines read, write, and overwrite entries for disjoint keys.
>
> —— https://pkg.go.dev/sync#Map

sync.Map 用读多写少的内部结构换取无锁读，写入密集时维护成本高于 map+RWMutex；官方文档明确只优化两种模式，超出这两种模式的选用没有性能依据且失去编译期类型检查。

**正例**

```go
mu.Lock()
m[key] = value
mu.Unlock()
```

**反例**

```go
sm.Store(key, value)
```

**依据**

- https://pkg.go.dev/sync#Map

**检测**：人工核对（核对 sync.Map 的读写比与键集模式）

### 【SHOULD】CONC-015 批量并发任务设置并发上限。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

按输入规模无界启动 goroutine 会耗尽资源。

- 用 `errgroup.Group` 的 `SetLimit` 或带缓冲 channel 信号量限流
- 上限依据下游容量（连接池、下游 QPS）确定并注释
- 一次性小任务（如有界输入的 fan-out）可豁免并说明

**为什么**

> SetLimit limits the number of active goroutines in this group to at most n.
>
> —— https://pkg.go.dev/golang.org/x/sync/errgroup#Group.SetLimit

每个 goroutine 都占用栈内存与调度成本，输入一万条就启动一万个 goroutine 时，内存和下游连接很快耗尽；设上限后并发度固定，错误也能经 errgroup 统一收集。

**正例**

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(8)
for _, t := range tasks {
  g.Go(func() error { return handle(ctx, t) })
}
err := g.Wait()
```

**反例**

```go
for _, t := range tasks {
  go handle(ctx, t)
}
```

**依据**

- https://pkg.go.dev/golang.org/x/sync/errgroup#Group.SetLimit

**检测**：人工核对（核对批量任务的并发创建是否设上限）

### 【SHOULD】CONC-016 sync.Pool 对象放回前 Reset。

- 归属：编程规约/并发规约
- 起始版本：Go 1.3
- 主题：concurrency, performance

放回池中的对象必须清掉内部状态。

- `Reset` 清空 slice、map、指针字段等引用，避免脏数据带入下一次 `Get`
- 不假设 `Get` 返回零值对象
- 池化对象文档需注明 `Reset` 义务

**为什么**

> Pool's purpose is to cache allocated but unused items for later reuse, relieving pressure on the garbage collector.
>
> —— https://pkg.go.dev/sync#Pool

Pool 只缓存内存不复位状态，脏对象把上一次使用的数据带进下一次，跨请求泄漏用户数据或造成逻辑错乱；Get 的语义是「拿到可用内存」而非「拿到零值」，复位义务在使用方。

**正例**

```go
buf := pool.Get().(*bytes.Buffer)
defer func() {
  buf.Reset()
  pool.Put(buf)
}()
```

**反例**

```go
buf := pool.Get().(*bytes.Buffer)
defer pool.Put(buf)
```

**依据**

- https://pkg.go.dev/sync#Pool

**检测**：人工核对（核对 Pool.Put 前是否 Reset）

### 【SHOULD】CONC-017 常驻 goroutine 入口捕获 panic 并记日志。

- 归属：编程规约/并发规约
- 起始版本：Go 1.0
- 主题：concurrency

进程级长期存活的 goroutine 在入口处兜住 panic。

- 适用对象是后台 worker、消费循环、定时任务这类常驻 goroutine
- 捕获后记录 panic 值与堆栈（`debug.Stack()`），只吞不记等于把故障藏起来
- 有意 fail-fast（任其崩溃、交由编排系统重启）时可不捕获，需在代码或评审里写明理由
- 一次性 goroutine 与 `errgroup` 管理的 goroutine 不在本条范围

**为什么**

> If ServeHTTP panics, the server (the caller of ServeHTTP) assumes that the effect of the panic was isolated to the active request. It recovers the panic, logs a stack trace to the server error log, and either closes the network connection or sends an HTTP/2 RST_STREAM, depending on the HTTP protocol.
>
> —— https://pkg.go.dev/net/http#Server

panic 只能被同一 goroutine 内的 recover 捕获。`net/http` 对 handler 的 panic 自己兜底，把影响限制在当次请求；自建常驻 goroutine 没有框架兜底，一个 nil 解引用或切片越界就让整个进程退出，在途任务一起丢失。入口处 recover 并记录堆栈，把进程级故障降级为单任务失败，堆栈又保留了定位线索。

**正例**

```go
go func() {
  defer func() {
    if r := recover(); r != nil {
      slog.Error("worker panic", "panic", r, "stack", string(debug.Stack()))
    }
  }()
  handle(req)
}()
```

**反例**

```go
go handle(req)
```

**依据**

- https://pkg.go.dev/net/http#Server

**检测**：人工核对（核对常驻 goroutine 的入口是否 recover 并记录堆栈）

## （九）控制语句规约

### 【MUST】CTRL-003 循环体内不得 defer 资源释放。

- 归属：编程规约/控制语句规约
- 起始版本：Go 1.0
- 主题：control-flow, data-structures

defer 绑定函数而非循环迭代。

- 循环内打开的资源在当次迭代内用完即关，或把循环体抽成函数
- 长循环内 defer 会在函数返回时一次性释放，文件描述符与连接累积耗尽
- 需要「循环结束统一清理」时用显式收集再统一关闭

**为什么**

> The deferred call's arguments are evaluated immediately, but the function call is not executed until the surrounding function returns.
>
> —— https://go.dev/doc/effective_go#defer

defer 的执行时机是所在函数返回，循环一千次就推迟一千次释放，句柄在函数存活期内只增不减；服务器进程在长循环场景下表现为 fd 用尽、连接池占死。

**正例**

```go
for _, name := range names {
  func() {
    f, err := os.Open(name)
    if err != nil {
      return
    }
    defer f.Close()
    process(f)
  }()
}
```

**反例**

```go
for _, name := range names {
  f, err := os.Open(name)
  if err != nil {
    return
  }
  defer f.Close()
  process(f)
}
```

**依据**

- https://go.dev/doc/effective_go#defer

**检测**：人工核对（grep for 循环体内的 defer 语句）

### 【MUST】CTRL-004 HTTP handler 写出错误响应后立即 return。

- 归属：编程规约/控制语句规约
- 起始版本：Go 1.0
- 主题：control-flow, error-handling

http.Error 或 w.Write 之后不得继续执行业务逻辑。

- 错误分支写响应后必须 return，避免二次写响应
- 继续执行的代码可能再次 `WriteHeader`，触发 superfluous WriteHeader 警告
- 中间件改写响应的场景同样适用

**为什么**

> Changing the header map after a call to WriteHeader (or Write) has no effect.
>
> —— https://pkg.go.dev/net/http#ResponseWriter

写出响应后再执行的代码或者二次 WriteHeader（响应头已发出，改动无效并打日志），或者把本该保密的错误细节写进响应体；立即 return 让错误分支封闭。

**正例**

```go
if err := parse(r); err != nil {
  http.Error(w, "bad request", http.StatusBadRequest)
  return
}
```

**反例**

```go
if err := parse(r); err != nil {
  http.Error(w, "bad request", http.StatusBadRequest)
}
process(r)
```

**依据**

- https://pkg.go.dev/net/http#ResponseWriter

**检测**：人工核对（核对 handler 内 http.Error/w.Write 之后是否 return）

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

### 【SHOULD】CTRL-005 内层作用域不用 := 遮蔽外层 err。

- 归属：编程规约/控制语句规约
- 起始版本：Go 1.0
- 主题：control-flow, error-handling

错误变量跨作用域时保持同一个标识符。

- `if v, err := f(); err != nil` 的 err 与外层 err 无关，赋值不会带到外层
- 需要把内层结果带出作用域时先声明再赋值
- shadow 检查器需显式开启，误报场景用 nolint 注明

**为什么**

> check for shadowed variables
>
> —— https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow

:= 在新作用域声明新变量，if 语句里的 err 是全新变量，检查它之后外层 err 仍是 nil，错误处理被静默架空；这类缺陷在编译期与运行期都不报错，只有评审或 shadow 检查能拦住。

**正例**

```go
var err error
for _, p := range paths {
  data, err := os.ReadFile(p)
  if err != nil {
    return err
  }
  _ = data
}
return err
```

**反例**

```go
var err error
for _, p := range paths {
  data, err := os.ReadFile(p)
  if err != nil {
    return err
  }
  _ = data
}
return nil
```

**依据**

- https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow

**检测**：golangci-lint govet(shadow)

## （十）注释规约

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

