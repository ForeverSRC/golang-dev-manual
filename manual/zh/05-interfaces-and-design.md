# 五、接口与设计规约

## （一）接口定义与落位

### 【SHOULD】ABSTRACT-001 抽象按实际需要引入，不提前铺开。

- 归属：接口与设计规约/接口定义与落位
- 起始版本：Go 1.0

接口、工厂、注册表在出现第二个实现、需要测试替身或确有变化点要隔离时才引入；只有一个实现且没有替换与测试需求时，先用具体类型。分层与端口在确有隔离需求时属于正当抽象，不受本条限制。

**为什么**

> Do not define interfaces before they are used: without a realistic example of usage, it is too difficult to see whether an interface is even necessary, let alone what methods it ought to contain.
>
> —— https://go.dev/wiki/CodeReviewComments#interfaces

提前铺的抽象多一层跳转与维护成本，真实需求到来时方向常常与预设不同，接口还要回头改；等真实用例出现再抽，接口由使用方定义，形状才对得上。

**正例**

```go
type Loader struct {
	path string
}

func (l *Loader) Load() (*domain.Manual, error) {
	return readJSON(l.path)
}
```

**反例**

```go
// only JSON exists today, yet an interface, a factory, and a format registry are defined up front for a future YAML source
type Source interface {
	Load() (*domain.Manual, error)
}

type SourceFactory interface {
	New(format string) (Source, error)
}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#interfaces

**检测**：人工核对（核对接口或工厂是否已有第二个实现、替换或测试需求）

### 【SHOULD】IFACE-001 接口按实现数量落位：单实现与实现同包，多实现由使用方定义。

- 归属：接口与设计规约/接口定义与落位
- 起始版本：Go 1.0

落位看预期有几个实现：

- 只会有单实现：接口与实现放在同一个包，接口紧邻实现，读代码时不必跨包对照；接口与构造函数导出，实现结构体不导出。
- 预计有多个实现或需要替换（存储、外部服务、通知一类）：接口放在使用方所在包，按使用方真正用得着的方法选取；实现包导出具体类型与构造函数，并按 IFACE-002 写编译期断言，把承诺写进代码。
- 接口签名已由标准库或第三方固定（`io.Reader` 一类）：本条不适用，实现方按原样满足即可。

**为什么**

> Go interfaces generally belong in the package that uses values of the interface type, not the package that implements those values. The implementing package should return concrete (usually pointer or struct) types: that way, new methods can be added to implementations without requiring extensive refactoring.
>
> —— https://go.dev/wiki/CodeReviewComments#interfaces

单实现时把接口另立一个包，读实现的上下文要在两个包之间来回对照，改一处方法签名也要两处对齐；接口与实现同包，一个文件就能看全契约与行为。上游 Go wiki 的通行说法是「接口属于使用接口类型值的包」（下面引文），针对的是会被替换或由多方使用的依赖，取向与本条的单实现分支不同——单实现分支按本地查阅习惯定，是本手册自己的约定；多实现分支沿用上游说法，那种场景下接口只含使用方用得着的方法，实现方加方法不必回头改接口。

**正例**

```go
// single implementation: interface and implementation live in the same package, contract and behavior are readable in one place
package service

type ManualService interface {
	List(ctx context.Context) ([]Clause, error)
}

type manualService struct{ repo ClauseRepository }
```

**反例**

```go
// single implementation, yet the interface lives in the consumer package, so reading the implementation means cross-package comparison
package clihandler

type ManualService interface {
	List(ctx context.Context) ([]Clause, error)
}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#interfaces

**检测**：人工核对（核对接口所在包与实现数量是否对应：单实现看是否与实现同包，多实现看是否在使用方包）

### 【SHOULD】IFACE-002 用编译期断言声明类型对接口的实现。

- 归属：接口与设计规约/接口定义与落位
- 起始版本：Go 1.0

在类型定义附近写 `var _ 接口 = 该类型零值`，指针、slice、map 用 nil，结构体用空结构体字面量。接口由使用方定义时，本条与 IFACE-001 成对使用，使用方声明接口、实现包写断言，缺一项该做法就不成立。适用场景：

- 使用方在多实现场景下按 IFACE-001 定义的接口
- 导出类型按 API 契约必须实现标准库或第三方定义的接口
- 一组类型共同实现同一接口
- 不满足接口会让调用方直接编译失败

**为什么**

> Verify interface compliance at compile time where appropriate. ... The statement `var _ http.Handler = (*Handler)(nil)` will fail to compile if `*Handler` ever stops matching the `http.Handler` interface.
>
> —— https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance

断言把「类型是否满足接口」这一判断交给编译器：接口改了、方法签名对不上，构建当场失败，改的人立刻知道；没有断言时，不匹配只在某处赋值给接口类型时才暴露，位置往往离改动很远。

**正例**

```go
type Handler struct{}

var _ http.Handler = (*Handler)(nil)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
```

**反例**

```go
type Handler struct{}

// without the assertion, the moment *Handler stops satisfying http.Handler only surfaces at the call site
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
```

**依据**

- https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance

**检测**：人工核对（核对按 API 契约实现接口的导出类型处是否有编译期断言）

### 【SHOULD】IFACE-003 不定义指向接口的指针。

- 归属：接口与设计规约/接口定义与落位
- 起始版本：Go 1.0

参数、返回值与字段都用接口值本身，不用 `*Interface`。接口值内部已含类型信息与数据指针，底层数据是值还是指针由具体类型决定；需要方法改到底层数据时，由实现方用指针接收者表达。

**为什么**

> You almost never need a pointer to an interface. You should be passing interfaces as values—the underlying data can still be a pointer.
>
> —— https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces

`*Interface` 让调用处多一层解引用，nil 判断也变复杂：接口值为 nil 与接口内部数据指针为 nil 是两种状态，混在一起后判空条件要么写错要么靠猜。传接口值本身表达同一意图，判空只有一种含义。

**正例**

```go
func Save(w io.Writer) error
```

**反例**

```go
func Save(w *io.Writer) error
```

**依据**

- https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces

**检测**：人工核对（核对函数签名与结构体字段是否出现 *接口类型）

## （二）组合与复用

### 【SHOULD】COMPOSE-001 复用能力通过组合表达，不搭继承式层次。

- 归属：接口与设计规约/组合与复用
- 起始版本：Go 1.0

复用一个类型的能力时，用字段显式持有它；需要把它的方法一并暴露出去时才用匿名内嵌。给类型起 Base、Abstract 一类名字再内嵌，是在模拟继承层次，应避免。接口之间同样用内嵌组合小接口，不造大而全的接口。

**为什么**

> Go does not provide the typical, type-driven notion of subclassing, but it does have the ability to "borrow" pieces of an implementation by embedding types within a struct or interface. ... There's an important way in which embedding differs from subclassing. When we embed a type, the methods of that type become methods of the outer type, but when they are invoked the receiver of the method is the inner type, not the outer one.
>
> —— https://go.dev/doc/effective_go#embedding

组合关系写在字段上，替换、测试与阅读都能直接定位；Go 的嵌入没有 Java 那样的覆写规则，方法落到哪一层靠名字解析推断，模拟继承会让阅读者判断不出实际调用的是哪个实现。

**正例**

```go
type Server struct {
	logger *slog.Logger
	store  Store
}
```

**反例**

```go
type Server struct {
	BaseServer
}
```

**依据**

- https://go.dev/doc/faq#inheritance
- https://go.dev/doc/effective_go#embedding

**检测**：人工核对（核对内嵌类型是否承担基类角色（Base / Abstract 命名、多层嵌入））

### 【SHOULD】COMPOSE-002 嵌入不改变外层的零值、拷贝语义与公开接口。

- 归属：接口与设计规约/组合与复用
- 起始版本：Go 1.0

只有当内层类型的方法应当成为外层方法时才用匿名内嵌，其余情况用具名字段。嵌入后外层零值仍可用、拷贝语义不变、不额外暴露无关方法。`sync.Mutex` 一类同步原语一律用具名字段，匿名嵌入会把 `Lock`、`Unlock` 挂到外层 API 上。

**为什么**

> Embedding should provide tangible benefit, like adding or augmenting functionality in a semantically-appropriate way. It should do this with zero adverse user-facing effects.
>
> —— https://github.com/uber-go/guide/blob/master/style.md#embedding-in-structs

嵌入选了公开类型，内层以后新增的方法就都成为外层的公开方法，外层被替掉或改换内层都是破坏性变更；嵌入指针类型还会让零值的外层一调用就 panic。具名字段只暴露想暴露的那一组方法，外层形状自己说了算。

**正例**

```go
type Server struct {
	mu     sync.Mutex
	logger *slog.Logger
}
```

**反例**

```go
type Server struct {
	sync.Mutex
	*bytes.Buffer
}
```

**依据**

- https://github.com/uber-go/guide/blob/master/style.md#embedding-in-structs
- https://github.com/uber-go/guide/blob/master/style.md#avoid-embedding-types-in-public-structs

**检测**：人工核对（核对嵌入是否改变外层零值、拷贝语义或公开无关方法（同步原语尤其））

## （三）依赖注入与装配

### 【SHOULD】DI-001 依赖装配用 wire 在编译期生成。

- 归属：接口与设计规约/依赖注入与装配
- 起始版本：Go 1.0

各层构造函数与接口绑定集中在一个 ProviderSet，接口到实现的绑定写 `wire.Bind`；入口只调用生成好的 `Initialize` 函数。`wire` 版本经 go.mod 的 tool 指令固定，改动装配后用 `go generate` 重新生成 wire_gen.go。

**为什么**

> Wire is a code generation tool that automates connecting components using dependency injection. ... Because Wire operates without runtime state or reflection, code written to be used with Wire is useful even for hand-written initialization.
>
> —— https://github.com/google/wire

运行时反射注入把字段名写错、绑定缺失推迟到启动或运行才暴露；编译期生成在缺绑定或多绑定时直接报错，装配关系也能在生成代码里逐行核对。

**正例**

```go
var ProviderSet = wire.NewSet(
	jsonfile.New,
	wire.Bind(new(service.ClauseRepository), new(*jsonfile.Repository)),
)
```

**反例**

```go
container := dig.New()
if err := container.Provide(jsonfile.New); err != nil {
	panic(err)
}
```

**依据**

- https://github.com/google/wire

**检测**：人工核对（核对仓库有 wire_gen.go 且端口绑定集中在 ProviderSet）

### 【SHOULD】DI-002 依赖经构造函数显式传入，不读写包级可变变量。

- 归属：接口与设计规约/依赖注入与装配
- 起始版本：Go 1.0

时间源、随机源、客户端、配置这些依赖作为构造参数或字段传入，包级变量只放只读常量。测试要替换行为时传替身，不改全局状态。函数指针与其他类型的值同样适用。

**为什么**

> Avoid mutating global variables, instead opting for dependency injection. This applies to function pointers as well as other kinds of values.
>
> —— https://github.com/uber-go/guide/blob/master/style.md#avoid-mutable-globals

可变全局把依赖藏在调用链之外：读代码看不出这段逻辑用的是哪个时间源或客户端，测试要替换得先存旧值、测完再还原，并行执行的用例还会互相干扰。作为参数传入后，依赖写在构造函数签名上，替换与阅读都直接。

**正例**

```go
type Signer struct {
	now func() time.Time
}

func NewSigner(now func() time.Time) *Signer {
	return &Signer{now: now}
}
```

**反例**

```go
var now = time.Now

func Sign(msg string) string {
	return signWithTime(msg, now())
}
```

**依据**

- https://github.com/uber-go/guide/blob/master/style.md#avoid-mutable-globals

**检测**：人工核对（核对包级变量是否存在赋值，依赖是否从构造函数传入）

## （四）context 传递

### 【MUST】CTX-002 不向函数传 nil context，占位用 context.TODO()。

- 归属：接口与设计规约/context 传递
- 起始版本：Go 1.7

暂时拿不到上游 context 时传 `context.TODO()`，程序入口、初始化与测试用 `context.Background()`。不把 nil 当作「没有 context」的表达。

**为什么**

> Do not pass a nil Context, even if a function permits it. Pass context.TODO if you are unsure about which Context to use.
>
> —— https://pkg.go.dev/context#pkg-overview

nil 是合法的接口值，能在编译期通过；接收方一旦调用 `Done` 或 `Value` 就直接 panic，是否 panic 取决于其内部实现，函数签名上看不出来。传 `context.TODO()` 语义相同且始终安全，静态检查也能据此找出这类遗漏。

**正例**

```go
ctx := context.TODO()
return svc.Load(ctx, path)
```

**反例**

```go
return svc.Load(nil, path)
```

**依据**

- https://pkg.go.dev/context#pkg-overview
- https://staticcheck.dev/docs/checks/#SA1012

**检测**：golangci-lint staticcheck(SA1012)（SA1012 检出向函数传 nil context 的调用）

### 【SHOULD】CTX-001 不把 context.Context 存进结构体，作为函数第一个参数传递。

- 归属：接口与设计规约/context 传递
- 起始版本：Go 1.7

context 承载单次请求的取消、超时与传递值，生命周期与请求一致。跨层传递时把 context 作为第一个参数，命名统一用 ctx。

**为什么**

> Most functions that use a Context should accept it as their first parameter. ... Don't add a Context member to a struct type; instead add a ctx parameter to each method on that type that needs to pass it along.
>
> —— https://go.dev/wiki/CodeReviewComments#contexts

存进结构体后，context 的生命周期脱离请求，框架的取消传播无法覆盖，容易形成泄漏，也无法在接口上表达取消边界。

**正例**

```go
func (s *Service) Get(ctx context.Context, id string) (*Order, error) {
	return s.repo.Find(ctx, id)
}
```

**反例**

```go
type Service struct {
	ctx  context.Context
	repo Repository
}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#contexts
- https://pkg.go.dev/context

**检测**：golangci-lint containedctx（首参约定需人工核对）

### 【SHOULD】CTX-003 context.Value 的 key 用未导出的自定义类型。

- 归属：接口与设计规约/context 传递
- 起始版本：Go 1.7

每个包为自己的取值定义 `type key int` 一类未导出类型与未导出常量，再提供类型安全的存取函数，调用方不直接接触 key。不用 string、int 等内建类型作 key。

**为什么**

> A key can be any type that supports equality; packages should define keys as an unexported type to avoid collisions.
>
> —— https://pkg.go.dev/context#Context.Value

内建类型作 key 时，不同包只要用同一个字符串或整数值就落到同一个槽位，读写互相覆盖且没有编译期提示；未导出类型只有本包能构造，跨包撞不上，取值也只能经本包提供的函数，取回的类型由该函数保证。

**正例**

```go
type key int

const userKey key = 0

func WithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}
```

**反例**

```go
ctx = context.WithValue(ctx, "user", u)
```

**依据**

- https://pkg.go.dev/context#Context.Value
- https://staticcheck.dev/docs/checks/#SA1029

**检测**：golangci-lint staticcheck(SA1029)（SA1029 检出用内建类型作 context key）

### 【SHOULD】CTX-004 context.Value 只放请求作用域数据，不放可选参数。

- 归属：接口与设计规约/context 传递
- 起始版本：Go 1.7

用户身份、追踪信息、截止时间这类跨进程、跨 API 传递的请求数据才放 context；函数自己的可选参数与配置项走显式参数或选项结构体。

**为什么**

> Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.
>
> —— https://pkg.go.dev/context#pkg-overview

塞进 context 的参数从函数签名里消失，读签名看不出函数依赖哪些输入；调用方漏传也只在运行时取到 nil，错误发生点离真正原因很远。走显式参数或选项结构体，依赖写在签名上，编译器与评审都能当场核对。

**正例**

```go
// request-scoped data: user identity goes into context
ctx := context.WithValue(r.Context(), userKey, user)

// optional parameters: template and render options are passed via the function signature
func Render(ctx context.Context, tmpl *Template, opts Options) error {
	u, _ := ctx.Value(userKey).(*User)
	return tmpl.Execute(u, opts)
}
```

**反例**

```go
// optional parameters: the template is stuffed into context and invisible in the function signature
ctx = context.WithValue(ctx, tmplKey, tmpl)

func Render(ctx context.Context) error {
	tmpl, _ := ctx.Value(tmplKey).(*Template)
	return tmpl.Execute(nil, Options{})
}
```

**依据**

- https://pkg.go.dev/context#pkg-overview

**检测**：人工核对（核对 context.Value 里放的是否为请求作用域数据）

## （五）结构体安全形状

### 【MUST】STRUCT-001 含 sync.Mutex 等同步原语的结构体禁止值拷贝，一律用指针传递或接收。

- 归属：接口与设计规约/结构体安全形状
- 起始版本：Go 1.0

同步原语复制后不再共享同一把锁，互斥语义会静默失效。含锁的结构体用指针接收者定义方法，并在文档注释里写明不可拷贝。

**为什么**

> copylocks check for locks erroneously passed by value
>
> —— https://pkg.go.dev/cmd/vet

值拷贝在编译期合法，`go vet` 的 copylocks 检查是主要防线；一旦漏检，问题只在并发压力下暴露，排查成本高。

**正例**

```go
type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
```

**反例**

```go
func (c Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
```

**依据**

- https://pkg.go.dev/cmd/vet
- https://google.github.io/styleguide/go/decisions

**检测**：golangci-lint govet(copylocks)

### 【MAY】STRUCT-002 不应支持比较的结构体，用不可比较字段让编译期拦截 ==。

- 归属：接口与设计规约/结构体安全形状
- 起始版本：Go 1.0

含内部状态、缓存或指针字段、用 == 比不到业务身份的结构体，嵌入一个 `[0]func()` 命名的辅助类型；此后 == 与把它用作 map key 的写法都不再编译。只影响比较写法，不占额外内存。

**为什么**

> DoNotCompare can be embedded in a struct to prevent comparability.
>
> —— https://github.com/protocolbuffers/protobuf-go/blob/master/internal/pragma/pragma.go

字段全部可比较时，== 在编译期合法，比的是逐字段内容；这类比较偶然而非有意，只在特定数据组合下才成立。嵌入不可比较字段后，误用直接编译失败，代价是零长度数组，不改变结构体大小与字段含义。

**正例**

```go
type DoNotCompare [0]func()

type Session struct {
	DoNotCompare
	id string
}
```

**反例**

```go
type Session struct {
	id   string
	seen bool
}

// all fields are comparable, so == compiles, but it compares internal state rather than business identity
if s1 == s2 {
	// ...
}
```

**依据**

- https://github.com/protocolbuffers/protobuf-go/blob/master/internal/pragma/pragma.go
- https://go.dev/ref/spec#Comparison_operators

**检测**：人工核对（核对有意禁止比较的导出结构体是否嵌入了不可比较字段）

