# 二、错误与日志

## （一）错误处理规约

### 【MUST】ERR-001 需要保留错误链时用 %w 包装，不使用 %v 或 %s。

- 归属：错误与日志/错误处理规约
- 起始版本：Go 1.13
- 主题：error-handling

需要让调用方用 `errors.Is`、`errors.As` 判断错误时，用 %w 包装，并附上本层的上下文。只有真正终止错误链的边界（对外输出、写日志）才用 %v。

**为什么**

> In Go 1.13, the fmt.Errorf function supports a new %w verb. When this verb is present, the error returned by fmt.Errorf will have an Unwrap method returning the argument of %w, which must be an error. In all other ways, %w is identical to %v.
>
> —— https://go.dev/blog/go1.13-errors

%v 会切断错误链，调用方无法再做类型判断，只能比较字符串，一改文案就失效。

**正例**

```go
if err := repo.Save(ctx, order); err != nil {
	return fmt.Errorf("save order %s: %w", order.ID, err)
}
```

**反例**

```go
if err := repo.Save(ctx, order); err != nil {
	return fmt.Errorf("save order %s: %v", order.ID, err)
}
```

**依据**

- https://go.dev/blog/go1.13-errors
- https://google.github.io/styleguide/go/decisions

**检测**：golangci-lint errorlint

### 【MUST】ERR-003 判断错误种类用 errors.Is 与 errors.As，不用 == 或类型断言。

- 归属：错误与日志/错误处理规约
- 起始版本：Go 1.13
- 主题：error-handling

裸比较与类型断言只看错误链的最外层：错误被 %w 包装后 == 不再相等，`err.(*MyError)` 直接断言失败。`errors.Is` 按哨兵错误逐层比较，`errors.As` 按错误类型逐层匹配，两者都会展开错误链。可供判断的哨兵错误与错误类型由产生错误的包导出；只有确认错误从未被包装时，才可以用 == 比较。

**为什么**

> In the simplest case, the errors.Is function behaves like a comparison to a sentinel error, and the errors.As function behaves like a type assertion. When operating on wrapped errors, however, these functions consider all the errors in a chain.
>
> —— https://go.dev/blog/go1.13-errors

某一层把直接返回改成 %w 包装后，== 与类型断言会静默失配，错误分支不再触发，问题往往到线上才暴露；errorlint 覆盖这两类写法，跟随它能自动拦住。

**正例**

```go
if errors.Is(err, os.ErrNotExist) {
	// file does not exist
}
```

**反例**

```go
if err == os.ErrNotExist {
	// file does not exist
}
```

**依据**

- https://go.dev/blog/go1.13-errors
- https://pkg.go.dev/errors#Is

**检测**：golangci-lint errorlint（覆盖 == 比较与类型断言两类写法）

### 【MUST】ERR-006 同一个错误只处理一次：包装返回或就地记录，不两者都做。

- 归属：错误与日志/错误处理规约
- 起始版本：Go 1.13
- 主题：error-handling

在能决定处置方式的那一层处理错误：上层还有调用方能应对，就补上上下文直接返回；已经到边界、无人再接，就地记录并让流程降级。不写记一条日志再返回 err 的写法，同一处既记录又返回时，上层往往再记一次，同一失败在日志里出现多条。

**为什么**

> You should only handle errors once. Handling an error means inspecting the error value, and making a single decision. ... But making more than one decision in response to a single error is also problematic.
>
> —— https://dave.cheney.net/practical-go/presentations/qcon-china.html#_only_handle_an_error_once

重复记录把日志条数放大，定位问题时先要从这堆重复里分辨真实失败次数；记录权收到一处后，日志条数与失败发生次数一一对应。

**正例**

```go
if err := save(order); err != nil {
	return fmt.Errorf("save order %s: %w", order.ID, err)
}
```

**反例**

```go
if err := save(order); err != nil {
	slog.Error("save order failed", "error", err)
	return err
}
```

**依据**

- https://dave.cheney.net/practical-go/presentations/qcon-china.html#_only_handle_an_error_once
- https://go.dev/wiki/CodeReviewComments#handle-errors

**检测**：人工核对（逐个 if err != nil 分支核对：同一分支内既写日志又 return err 即违反）

### 【SHOULD】ERR-004 可预期的失败返回 error，不用 panic。

- 归属：错误与日志/错误处理规约
- 起始版本：Go 1.16
- 主题：error-handling

panic 留给调用方无法继续的情况，例如启动时缺少必需配置、按逻辑不可能走到的分支被走到。可预期的失败（参数不合法、资源不存在、下游返回错误）用 error 返回，由调用方决定处置方式。库代码不 panic；recover 只在进程入口兜底。

**为什么**

> Don't use panic for normal error handling. Use error and multiple return values.
>
> —— https://go.dev/wiki/CodeReviewComments#dont-panic

panic 越过调用方的错误分支，把本可就地降级的失败升级为整个调用栈中断；recover 只能拦住同一 goroutine 里的 panic，其它 goroutine 的 panic 直接终止进程，兜不住。

**正例**

```go
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parse(data)
}
```

**反例**

```go
func Load(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return parse(data)
}
```

**依据**

- https://go.dev/wiki/CodeReviewComments#dont-panic
- https://go.dev/doc/effective_go#errors

**检测**：人工核对（核对 panic 是否用在可预期的失败上）

### 【SHOULD】ERR-005 失败用 error 或 ok 返回，不用特殊值表示。

- 归属：错误与日志/错误处理规约
- 起始版本：Go 1.0
- 主题：error-handling

不用 -1、空字符串、零值表达「没找到」或「失败」，改用 error，或在不需要说明原因时加一个 ok 返回值。调用方必须先判断 ok 或 err 才能使用结果，漏判当场编译失败。返回值里的 nil、0 本身是合法结果时不受本条限制。

**为什么**

> Go's support for multiple return values provides a better solution. Instead of requiring clients to check for an in-band error value, a function should return an additional value to indicate whether its other return values are valid.
>
> —— https://go.dev/wiki/CodeReviewComments#in-band-errors

特殊值与被测函数的合法结果同在返回值里，调用方漏判时错误顺着计算往下传，直到别处才炸开，排查指向的是错误的位置；多返回值让漏判变成编译错误。

**正例**

```go
func Lookup(id string) (Order, bool) {
	// ...
}

order, ok := Lookup(id)
if !ok {
	return ErrNotFound
}
```

**反例**

```go
func Lookup(id string) Order {
	return Order{} // return zero value when not found
}

order := Lookup(id) // the caller cannot tell not-found from a zero value
```

**依据**

- https://go.dev/wiki/CodeReviewComments#in-band-errors

**检测**：人工核对（核对是否用 -1、空值或零值表达失败）

## （二）错误码规约

### 【SHOULD】CODE-001 对外错误码沿用既有标准分类，不另建编号体系。

- 归属：错误与日志/错误码规约
- 起始版本：Go 1.0
- 主题：api-design

服务对外返回错误时，用所处协议的标准分类表达大类（HTTP 状态码、gRPC 状态码），细分原因放进错误详情或自定义错误类型。业务确有标准分类覆盖不到的大类时，在标准码之上追加原因字段，不新造与标准码并列的编号。

**为什么**

> Google APIs must use the canonical error codes defined by google.rpc.Code. Individual APIs must avoid defining additional error codes, since developers are very unlikely to write logic to handle a large number of error codes.
>
> —— https://cloud.google.com/apis/design/errors

自造编号要调用方为每套码写一遍映射与判断，码一多，调用方实际只按大类处理，细分部分白做；标准分类自带各语言客户端的映射与既定的重试语义，跟着用可以直接复用。

**正例**

```go
http.Error(w, "order not found", http.StatusNotFound)
```

**反例**

```go
http.Error(w, "1001", http.StatusOK)
```

**依据**

- https://cloud.google.com/apis/design/errors

**检测**：人工核对（核对对外错误码是否与标准分类并存了自造编号）

### 【SHOULD】CODE-002 错误消息不作为契约，判断错误只看码与类型。

- 归属：错误与日志/错误码规约
- 起始版本：Go 1.13
- 主题：api-design, error-handling

错误消息面向人，措辞随时可改：调用方、测试与告警规则都按错误码或错误类型判断，不用消息文本。需要给用户看本地化文案时，把文案放进单独的字段或错误详情，与给开发者看的消息分开。

**为什么**

> Error messages are not part of the API surface. They are subject to changes without notice. Application code must not have a hard dependency on error messages.
>
> —— https://cloud.google.com/apis/design/errors

按消息文本判断的分支会在文案调整后静默失配，例如补一个参数名就走进默认分支；错误码与错误类型稳定，判断错误种类要看它们。

**正例**

```go
if errors.Is(err, ErrNotFound) {
	// branch by error type
}
```

**反例**

```go
if strings.Contains(err.Error(), "not found") {
	// branch by message text
}
```

**依据**

- https://cloud.google.com/apis/design/errors

**检测**：人工核对（核对是否有按错误消息文本判断的分支）

### 【SHOULD】CODE-003 跨服务传播错误时转换错误码，不透传下游原始错误。

- 归属：错误与日志/错误码规约
- 起始版本：Go 1.0
- 主题：api-design

收到下游错误后判断语义，再按本服务的错误码对外返回：下游内部故障归为不可用或内部错误，参数错误只在确由本服务调用方引起时保留。下游的原始错误与错误链写进本地日志，不放进对外响应。

**为什么**

> If your API service depends on other services, you should not blindly propagate errors from those services to your clients. When translating errors, we suggest the following: Hide implementation details and confidential information. Adjust the party responsible for the error.
>
> —— https://cloud.google.com/apis/design/errors

透传会把下游服务名、参数名与内部结构暴露给外部调用方，还把责任归错方——下游报参数非法，外部看到的却是本服务的参数非法；转换后对外只有本服务的语义，细节留在自己的日志里。

**正例**

```go
if err := callOrder(ctx); err != nil {
	// map the downstream failure to unavailable, write the original error to the local log
	return status.Error(codes.Unavailable, "order service unavailable")
}
```

**反例**

```go
if err := callOrder(ctx); err != nil {
	return err // pass through the downstream error directly
}
```

**依据**

- https://cloud.google.com/apis/design/errors

**检测**：人工核对（核对对外返回的错误是否直接来自下游）

## （三）日志规约

### 【MUST】LOG-001 业务日志用 log/slog 记录，不用 fmt 与 log 直接打印。

- 归属：错误与日志/日志规约
- 起始版本：Go 1.21
- 主题：logging

结构化日志统一走 `log/slog`：进程入口用 `slog.New` 装配 `Handler`，生产环境用 `JSONHandler`，其余位置用注入的 `*slog.Logger` 或 `slog.Default()`。`fmt.Print` 与标准 `log` 只用于一次性脚本的输出；接入第三方日志库时，经自定义 `Handler` 收口到同一套输出。

**为什么**

> Structured logs use key-value pairs so they can be parsed, filtered, searched, and analyzed quickly and reliably.
>
> —— https://go.dev/blog/slog

`fmt` 与 `log` 打出的整行文本无法按字段解析，收集端只能整段存下，按 order_id 检索、按耗时聚合都要另写解析规则；`slog` 的键值对与 JSON 输出可直接被收集端识别。

**正例**

```go
slog.Info("order created", "order_id", order.ID)
```

**反例**

```go
fmt.Printf("order created: %s\n", order.ID)
```

**依据**

- https://go.dev/blog/slog
- https://pkg.go.dev/log/slog

**检测**：grep 正则 \b(fmt\.Print|log\.Print)（命中处人工核对是否为业务日志）

### 【SHOULD】LOG-002 日志字段写成键值对，不拼进消息文本。

- 归属：错误与日志/日志规约
- 起始版本：Go 1.21
- 主题：logging

随日志带出的取值（ID、数量、耗时、错误）都作为参数写成 key-value，消息只写发生了什么。同一个含义在仓库内用同一个键名；错误值用 `slog.Any("error", err)` 传入，不写进消息。

**为什么**

> Unlike with the log package, we can easily add key-value pairs to our output by writing them after the message.
>
> —— https://go.dev/blog/slog

取值拼进消息后，收集端只能对整行做模糊匹配，按字段精确检索与聚合都做不到；键值对在收集端成为独立字段，可以直接过滤与统计。

**正例**

```go
slog.Info("order created", "order_id", order.ID, "amount", order.Amount)
```

**反例**

```go
slog.Info(fmt.Sprintf("order %s created, amount %d", order.ID, order.Amount))
```

**依据**

- https://go.dev/blog/slog

**检测**：人工核对（核对日志取值是否拼进了消息字符串）

### 【SHOULD】LOG-003 请求路径上的日志用带 Context 的方法。

- 归属：错误与日志/日志规约
- 起始版本：Go 1.21
- 主题：context, logging

处理请求时用 `slog` 的 `InfoContext`、`ErrorContext` 等方法把 ctx 一并传入，handler 据此补上请求 ID、trace 标识一类贯穿整个请求的字段。启动、退出等没有请求上下文的路径用不带 Context 的方法。

**为什么**

> As the call to LogAttrs shows, you can pass a context.Context to some log functions so a handler can extract context information like trace IDs. (Canceling the context does not prevent the log entry from being written.)
>
> —— https://go.dev/blog/slog

日志不带 ctx 时，同一次请求产生的多条日志在收集端无法自动串起来，排查要先按时间窗捞一批再人工分辨；带上 ctx 后请求级字段由 handler 统一补全，不必每处手写。

**正例**

```go
slog.ErrorContext(ctx, "order save failed", "error", err)
```

**反例**

```go
slog.Error("order save failed", "error", err)
```

**依据**

- https://pkg.go.dev/log/slog#Logger.InfoContext
- https://go.dev/blog/slog

**检测**：人工核对（核对请求路径上的日志是否传入 ctx）

