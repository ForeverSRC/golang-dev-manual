# 八、静态检查规约

## （一）配置与运行

### 【MUST】TOOL-001 统一用 golangci-lint 做静态代码检查。

- 归属：静态检查规约/配置与运行
- 起始版本：Go 1.0

全仓库只保留一份 .golangci.yml，本地与 CI 跑同一条命令；`golangci-lint` 版本由 go.mod 的 tool 指令固定，不依赖各机器上的全局安装。

**为什么**

> Golangci-lint is a fast linters runner for Go. It runs linters in parallel, uses caching, supports YAML configuration, integrates with all major IDEs, and includes over a hundred linters.
>
> —— https://golangci-lint.run/

各人各机各版本的裸 linter 结果不一致，评审时无法对齐结论；集中一份配置再加固定版本，问题才能复现。

**正例**

```makefile
lint:
	golangci-lint run ./...
```

**反例**

```makefile
lint:
	go vet ./...
	staticcheck ./...
```

**依据**

- https://golangci-lint.run/

**检测**：人工核对（核对仓库根有 .golangci.yml 且 lint 入口调用 golangci-lint）

### 【MUST】TOOL-009 抑制告警要指名检查器并写明理由。

- 归属：静态检查规约/配置与运行
- 起始版本：Go 1.0

确认是误报后，在触发处写 //nolint:<检查器> // 理由，检查器名与理由缺一不可；不写裸 //nolint 与 //nolint:all，它们不表明被抑制的是哪条检查；也不在文件或包的顶部一次性关掉检查。配置里开启 nolintlint，用 require-specific、require-explanation、allow-unused 三项设置强制。

**为什么**

> You may add a comment explaining or justifying why a `nolint` directive is being used on the same line as the flag itself:
>
> —— https://golangci-lint.run/docs/linters/false-positives/#nolint-directive

抑制是给检查器开的口子，写清检查器与理由，评审时才知道被绕过的是哪条规则、绕过的代价是什么；裸 //nolint 与 //nolint:all 把口子开在所有检查上，日后新加的分析器也一并放行，问题不再有人看见。

**正例**

```go
raw, err := os.ReadFile(path) //nolint:gosec // path is specified explicitly by the caller
```

**反例**

```go
raw, err := os.ReadFile(path) //nolint
```

**依据**

- https://golangci-lint.run/docs/linters/false-positives/#nolint-directive
- https://golangci-lint.run/docs/linters/configuration/#nolintlint

**检测**：golangci-lint nolintlint（配置为 require-specific + require-explanation + allow-unused，拦下不指名检查器、缺理由与已失效的抑制）

## （二）必开规则

### 【MUST】TOOL-006 至少开启 golangci-lint 的 standard 集。

- 归属：静态检查规约/必开规则
- 起始版本：Go 1.0

standard 集是 `golangci-lint` 默认启用的五个检查器，配置里不要关掉，也不要只留其中几个：

- errcheck：未处理的错误返回值
- govet：可疑构造，等价于 `go vet`
- ineffassign：赋了值却不再读的赋值
- staticcheck：staticcheck 的诊断规则
- unused：未使用的标识符

**为什么**

> Golangci-lint can be used with zero configuration. By default, the following linters are enabled: errcheck, govet, ineffassign, staticcheck, unused.
>
> —— https://golangci-lint.run/docs/welcome/quick-start/

这五个对准的都是能确定判定的缺陷：漏掉的错误、写错但能编译的可疑构造、算了不用的中间值、永远走不到的死代码。少开哪个就等于放弃对应类别的机器拦截，只能退回评审靠记忆。

**正例**

```yaml
linters:
  default: standard
```

**反例**

```yaml
linters:
  default: none
  enable:
    - gocritic
```

**依据**

- https://golangci-lint.run/docs/welcome/quick-start/

**检测**：人工核对（核对 .golangci.yml 未关闭 standard 集）

## （三）推荐规则

### 【SHOULD】TOOL-002 golangci-lint 中开启 modernize，跟进新版 Go 惯用法。

- 归属：静态检查规约/推荐规则
- 起始版本：Go 1.22

Go 版本升级后，旧写法由 modernize 提示可简化的地方；未启用则这些惯用法只能靠人工记忆。

**为什么**

> modernize: A suite of analyzers that suggest simplifications to Go code, using modern language and library features.
>
> —— https://golangci-lint.run/docs/linters/configuration/#modernize

标准库与语言特性持续补齐，旧写法逐版变成技术债；交给分析器批量提示，成本低于事后重构。

**正例**

```go
for range n {
	handle()
}
```

**反例**

```go
for i := 0; i < n; i++ {
	handle()
}
```

**依据**

- https://golangci-lint.run/docs/linters/configuration/#modernize

**检测**：人工核对（核对 .golangci.yml 的 linters.enable 含 modernize）

### 【SHOULD】TOOL-003 使用 testify 时，同步开启 testifylint。

- 归属：静态检查规约/推荐规则
- 起始版本：Go 1.0

testifylint 拦下 `assert.Nil` 判错误、参数顺序颠倒一类高频误用；未用 `testify` 的项目不必启用。

**为什么**

> testify is the most popular Golang testing framework in recent years. But it has a terrible ambiguous API in places, and the purpose of this linter is to protect you from annoying mistakes.
>
> —— https://github.com/Antonboom/testifylint

testify 部分 API 误用后编译期不报错，如 `assert.Nil(t, err)` 在 err 为接口时判断失效；交给 linter 比靠评审记忆可靠。

**正例**

```go
require.NoError(t, err)
```

**反例**

```go
assert.Nil(t, err)
```

**依据**

- https://github.com/Antonboom/testifylint

**检测**：人工核对（核对 .golangci.yml 的 linters.enable 含 testifylint）

### 【SHOULD】TOOL-007 standard 集之外，推荐开启常用的补充分析器。

- 归属：静态检查规约/推荐规则
- 起始版本：Go 1.0

standard 集之外，推荐开启以下分析器：

- errorlint：错误链相关的写法问题
- nilerr：判断了 err != nil 却返回 nil
- bodyclose：HTTP 响应体未关闭
- noctx：请求未带 `context.Context`
- nolintlint：抑制告警未指名检查器、缺理由或已失效
- gosec：安全隐患
- gocritic：缺陷与性能诊断
- exhaustive：枚举 switch 未穷举
- copyloopvar：Go 1.22 起多余的循环变量拷贝
- testpackage：同包测试
- funcorder：导出与未导出的函数次序
- misspell：拼写错误

modernize 与 testifylint 见同节单列条款。与项目形态不符的可以不启用，理由写进配置注释或评审记录。

**为什么**

> To see a list of supported linters and which linters are enabled/disabled: golangci-lint help linters
>
> —— https://golangci-lint.run/docs/linters/

standard 集之外没有默认开启，不主动打开就等于没有。这些分析器各自对准一类具体缺陷或高频偏差，误报可控；清单只收与项目无关的通用项，项目专属检查在各自项目里另加。

**正例**

```yaml
linters:
  default: standard
  enable:
    - errorlint
    - gosec
    - nolintlint
    - testpackage
```

**反例**

```yaml
linters:
  default: standard
  # none of the supplementary analyzers is enabled
```

**依据**

- https://golangci-lint.run/docs/linters/

**检测**：人工核对（核对 enable 列表是否覆盖清单，未启用项是否有理由）

