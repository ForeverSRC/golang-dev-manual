# 九、工具规约

## （一）命令行

### 【SHOULD】TOOL-004 命令行程序用 cobra 组织子命令。

- 归属：工具规约/命令行
- 起始版本：Go 1.0

根命令下挂子命令，业务逻辑写进各命令的 `RunE`；输出写 `cmd.OutOrStdout()`，错误写 `cmd.ErrOrStderr()`。

**为什么**

> Cobra is a library providing a simple interface to create powerful modern CLI interfaces similar to git & go tools. ... Easy subcommand-based CLIs: app server, app fetch, etc.
>
> —— https://github.com/spf13/cobra

手写 `os.Args` 分发要在子命令、全局 `flag`、帮助文本上重复造轮子；`cobra` 把输出与参数注入命令，命令可直接在测试里执行。

**正例**

```go
cmd := &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, _ []string) error {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), "ok")
		return err
	},
}
```

**反例**

```go
switch os.Args[1] {
case "list":
	fmt.Println("ok")
}
```

**依据**

- https://github.com/spf13/cobra

**检测**：人工核对（核对是否以 cobra 组织子命令）

### 【SHOULD】TOOL-008 升级 Go 工具链后用 go fix 应用现代化改写。

- 归属：工具规约/命令行
- 起始版本：Go 1.26

把构建切到更新的 Go 版本后，在干净的 git 工作区跑 `go fix ./...`，改动单独成一个提交，便于逐处复核。

- 只看不改用 `go fix -diff ./...` 预览补丁
- 默认跑全部分析器；拆分改动时用 -<分析器> 只跑一个（如 -any），写 -<分析器>=false 则排除它
- 落在生成文件上的改写会被丢弃，这类改动要改的是生成器本身
- 按 GOOS / GOARCH 各跑一次，覆盖不同构建标签下的代码
- 包作者迁移自己的 API 时，在被替换的函数上写 //go:fix inline，由 inline 分析器落地

modernize 检查与本节工具指向同一批分析器，区别是 `golangci-lint` 在评审期提示，`go fix` 直接生成补丁。

**为什么**

> Fix is a tool executed by "go fix" to update Go programs that use old features of the language and library and rewrite them to use newer ones. After you update to a new Go release, fix helps make the necessary changes to your programs.
>
> —— https://pkg.go.dev/cmd/fix

语言与标准库每版都在补新写法，旧代码不改就逐版变成技术债；靠人工记忆与评审提出，没被想到的地方无人知道。`go fix` 把这批改写收成一条可复现的命令，跟着工具链升级跑一次即可。

**正例**

```console
$ go fix -diff ./...   # preview the patch first
$ go fix ./...          # then apply it
```

**反例**

```console
$ go vet ./...
# after upgrading the toolchain only static checks ran, so old idioms were never rewritten in bulk
```

**依据**

- https://pkg.go.dev/cmd/fix
- https://go.dev/blog/gofix
- https://go.dev/doc/go1.26

**检测**：人工核对（核对升级 Go 版本后是否跑过 go fix）

## （二）测试工具

### 【MUST】TOOL-005 测试断言统一使用 testify。

- 归属：工具规约/测试工具
- 起始版本：Go 1.0

断言用 `assert` 或 `require`，按失败是否终止用例选择：前置条件用 `require`，其余检查用 `assert`。测试文件里不手写 if 加 `t.Errorf`、`t.Fatalf` 做比较。

**为什么**

> The assert package provides some helpful methods that allow you to write better test code in Go. Prints friendly, easy to read failure descriptions.
>
> —— https://github.com/stretchr/testify#assert-package

手写比较的失败信息只有自定义文案，缺少期望与实际对照；`testify` 统一格式化差异并带出调用位置，排查时不必重读源码。

**正例**

```go
assert.Equal(t, want, got)
```

**反例**

```go
if got != want {
	t.Errorf("got %v, want %v", got, want)
}
```

**依据**

- https://github.com/stretchr/testify#assert-package

**检测**：grep 正则 \bt\.(Errorf|Fatalf)\(（命中处应改用 testify 断言）

## （三）代码生成

### 【MUST】TOOL-010 生成物不手工编辑，只由生成器产出。

- 归属：工具规约/代码生成
- 起始版本：Go 1.4

生成器与输入源是唯一可改的地方，改完重新生成，不在生成物上直接改。生成器自带 // Code generated ... DO NOT EDIT. 一类标记行时保留它，不删改；生成器不输出标记的，不必补造，识别交给 TOOL-011 的文件名约定。本仓库的 gdm/cmd/gdm-cli/wire/wire_gen.go 与 manual/ 都属此类，生成命令走 `make generate` 与 `make gen`。

**为什么**

> To convey to humans and machine tools that code is generated, generated source should have a line that matches the following regular expression (in Go syntax): ^// Code generated .* DO NOT EDIT\.$
>
> —— https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source

生成物是输入源与生成器的投影，直接改就与这个关系脱钩：下一次生成把改动覆盖掉，或者生成物与输入源长期不一致，成了一截不敢动的代码。官方给 Go 生成源码的识别手段是文件内标记行（下面引文），但标记要生成器主动输出，模板渲染、文档生成一类不写标记，这类生成物改靠文件名识别，见 TOOL-011。

**正例**

```go
// gdm/di/wire.go: the wiring change goes into the input source, then run make generate
var ProviderSet = wire.NewSet(jsonfile.New)
```

**反例**

```go
// gdm/cmd/gdm-cli/wire/wire_gen.go: the change is made on the generated file and is overwritten on the next run
repository := jsonfile.New()
```

**依据**

- https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source

**检测**：人工核对（核对改动是否落在生成器或输入源上）

### 【MUST】TOOL-011 生成文件命名带 gen 前缀或后缀，优先用后缀。

- 归属：工具规约/代码生成
- 起始版本：Go 1.0

生成器允许指定输出文件名时，文件名带 gen 前缀或后缀：

- 后缀优先，形如 wire_gen.go、order_repo_gen.go，按名称排序时与同类文件相邻。
- 前缀写成 gen_，形如 gen_version.go。
- 分隔符只用下划线，不写连字符。
- 文件名由生成器自己固定（`wire` 输出 wire_gen.go）或由平台决定时，按生成器默认走，不必改名。

**为什么**

> The name of the output file is created by replacing the .proto extension with .pb.go.
>
> —— https://protobuf.dev/reference/go/go-generated/

`protoc` 给生成文件加 .pb.go 后缀，`wire` 输出 wire_gen.go，都是在文件名上标明这是个生成物。名字由调用方指定时不加约定，生成物与手写文件混在同一层目录里，评审要逐个打开才知道哪个能改，改之前还得先翻生成命令。统一到 gen 前后缀后，目录列表本身就能作答；推荐后缀是因为它不打断原有的命名顺序，检索、排序与相邻关系仍按业务名前缀走。分隔符取自己的文件名约定，Go 工具链认的 `_test.go`、`source_windows_amd64.go` 都用下划线，标准库带连字符的文件只出现在 testdata 一类测试数据里。

**正例**

```go
// generated file names: wire_gen.go, mock_gen.go, clause_repo_gen.go
package wire
```

**反例**

```go
// generated file names: wire.go, mocks.go, clause_repo.go
package wire
```

**依据**

- https://protobuf.dev/reference/go/go-generated/
- https://pkg.go.dev/cmd/go#hdr-Build_constraints

**检测**：人工核对（核对生成物文件名是否带 gen 前缀或后缀）

