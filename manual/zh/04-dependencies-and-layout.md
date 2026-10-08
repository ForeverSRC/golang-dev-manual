# 四、依赖与工程结构

## （一）分层规约

### 【SHOULD】LAYER-001 按职责分层，依赖方向单向指向内层。

- 归属：依赖与工程结构/分层规约
- 起始版本：Go 1.0
- 主题：layering

领域模型层不含 IO，也不 import 其它层；应用服务层依赖领域模型，用接口声明所需的外部能力；适配器层实现这些接口，负责持久化、文件与网络；入口层只解析参数、渲染输出。外层 import 内层，内层代码里不出现外层包名。

**为什么**

> The rule to obey is that code pertaining to the inside part should not leak into the outside part.
>
> —— https://alistair.cockburn.us/hexagonal-architecture/

依赖反向时领域逻辑被 IO 细节绑住，跑一次测试也要带上真实文件或数据库；把外部能力放在内层声明的接口后面，替换实现不必改内层，测试可注入内存替身。

**正例**

```go
// inner layer declares the port, outer layer implements it
package service

type ClauseRepository interface {
	Load(ctx context.Context, path string) (*domain.Manual, error)
}
```

**反例**

```go
// inner layer depends directly on the outer adapter
package service

import "example/internal/repository/jsonfile"

type manualService struct {
	repo *jsonfile.Repository
}
```

**依据**

- https://alistair.cockburn.us/hexagonal-architecture/

**检测**：人工核对（核对包 import 方向；可用 golangci-lint 的 depguard 固化为规则）

## （二）包组织规约

### 【MUST】PKG-001 包名取导入路径的最后一段，与所在目录同名。

- 归属：依赖与工程结构/包组织规约
- 起始版本：Go 1.0
- 主题：naming

目录名与 package 子句写成同一个名字，多词目录同样连写。例外只有三类：

- 可执行文件目录声明 package main
- 包外测试声明「被测包名_test」
- 仅由生成代码导入的包可用下划线

**为什么**

> Client code uses the package path when importing the package. By convention, the last element of the package path is the package name:
>
> —— https://go.dev/blog/package-names

导入处写的是路径、引用处写的是包名，两者不一致时调用方要同时记住两个名字，按目录查找与按标识符查找会对不上；工具与日志只认路径，排查时得先反查目录才知道包名。

**正例**

```go
// directory gdm/internal/repository/jsonfile
package jsonfile
```

**反例**

```go
// directory gdm/internal/repository/jsonfile
package jsonrepo
```

**依据**

- https://go.dev/blog/package-names
- https://go.dev/wiki/CodeReviewComments#package-names

**检测**：人工核对（逐个目录核对 package 子句与目录名一致）

### 【MUST】PKG-003 import 分组排列，标准库一组、其余一组，组间空行。

- 归属：依赖与工程结构/包组织规约
- 起始版本：Go 1.0
- 主题：code-organization

组内按路径排序，不把不同来源混进同一组。本项目的 gci 配置分三组：standard、default、本仓库 module 前缀，交给 formatter 自动整理。

**为什么**

> Imports are organized in groups, with blank lines between them. The standard library packages are always in the first group.
>
> —— https://go.dev/wiki/CodeReviewComments#imports

分组把标准库与外部依赖在视觉上分开，评审时一眼能看出这次改动新加了哪个模块；混排时要逐行判断某个路径是不是标准库，依赖变动的范围也就看不出来。

**正例**

```go
import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)
```

**反例**

```go
import (
	"context"
	"github.com/spf13/cobra"
	"fmt"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)
```

**依据**

- https://go.dev/wiki/CodeReviewComments#imports

**检测**：golangci-lint gci（formatters 中启用 gci，分组顺序为 standard、default、本仓库 module 前缀）

### 【MUST】PKG-004 不使用 dot import。

- 归属：依赖与工程结构/包组织规约
- 起始版本：Go 1.0
- 主题：naming

唯一例外是包外测试因循环依赖无法与被测包同处一层时，可在 _test 文件里用 import . 引入被测包。常规代码一律显式写包名前缀。

**为什么**

> Except for this one case, do not use import . in your programs. It makes the programs much harder to read because it is unclear whether a name like Quux is a top-level identifier in the current package or in an imported package.
>
> —— https://go.dev/wiki/CodeReviewComments#import-dot

dot import 把另一个包的导出名直接拉进当前作用域，读到一个名字时无法从写法判断它来自哪个包，查定义只能先猜再试；省下的是几个字符的前缀，付出的是每次阅读都要多一次来源判断。

**正例**

```go
import "github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"

var m domain.Manual
```

**反例**

```go
import . "github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"

var m Manual
```

**依据**

- https://go.dev/wiki/CodeReviewComments#import-dot

**检测**：grep 正则 ^\s*(import\s+)?\.\s+"（命中处核对是否为循环依赖场景下的 _test 文件）

### 【SHOULD】PKG-002 不新建 util、common、helper 一类无信息量的包。

- 归属：依赖与工程结构/包组织规约
- 起始版本：Go 1.0
- 主题：naming

包名要让调用方看出领域边界。名字里只有「工具」或「通用」含义的包没有归属标准，谁都能往里放。已有的通用逻辑按真实职责拆成具名包，例如集合操作归 stringset、重试归 retry。

**为什么**

> Avoid uninformative package names like `util`, `utility`, `common`, `helper`, `model`, `testhelper`, and so on that would tempt users of the package to rename it when importing.
>
> —— https://google.github.io/styleguide/go/decisions#package-names

没有边界就没有拆分依据，这类包随时间累积不同来源的依赖，编译面变大、导入冲突变多；等到要拆时，调用点已经散落各处，迁移只能逐个改 import。

**正例**

```go
package stringset
```

**反例**

```go
package util
```

**依据**

- https://google.github.io/styleguide/go/decisions#package-names
- https://go.dev/blog/package-names

**检测**：grep 正则 ^package\s+(util|utility|common|misc|helper|helpers|model)\b（命中处按其真实职责拆成具名包）

### 【SHOULD】PKG-005 导出名不再重复包名，构造函数优先用 New。

- 归属：依赖与工程结构/包组织规约
- 起始版本：Go 1.0
- 主题：naming

调用方已用包名作前缀，导出类型与方法名里重复它只增加长度。包内只有一个类型、或构造函数返回与包同名的类型时用 `New`；返回包内其它类型时才在名字里带上该类型。

**为什么**

> Since client code uses the package name as a prefix when referring to the package contents, the names for those contents need not repeat the package name. The HTTP server provided by the `http` package is called `Server`, not `HTTPServer`.
>
> —— https://go.dev/blog/package-names

`widget.NewWidget` 在调用点展开就是 `widget.NewWidget()`，一半字符是同一信息；只有一个同名类型时用 `New`，名字更短，读法与 `time.Now`、`list.New` 这类标准库入口也一致。

**正例**

```go
package manual

func New() *Manual
```

**反例**

```go
package manual

func NewManual() *Manual
```

**依据**

- https://go.dev/blog/package-names
- https://google.github.io/styleguide/go/decisions#package-names

**检测**：人工核对（核对导出名是否重复其所在包名）

### 【SHOULD】PKG-006 blank import 只出现在 main 包或测试文件里。

- 归属：依赖与工程结构/包组织规约
- 起始版本：Go 1.0
- 主题：naming

只为一处副作用（注册驱动、初始化全局表）而导入的包，按约定只放在程序入口或需要它的测试中。库代码里出现 blank import，会把这份副作用强加给所有导入它的调用方。

**为什么**

> Packages that are imported only for their side effects (using the syntax `import _ "pkg"`) should only be imported in the main package of a program, or in tests that require them.
>
> —— https://go.dev/wiki/CodeReviewComments#import-blank

blank import 的副作用在编译期就生效，库包带上它，调用方没有退出的选项；集中在 main 包与测试，副作用的作用范围才与实际使用它的位置一致。

**正例**

```go
// gdm/cmd/gdm-cli/main.go
import _ "net/http/pprof"
```

**反例**

```go
// gdm/internal/server/server.go
import _ "net/http/pprof"
```

**依据**

- https://go.dev/wiki/CodeReviewComments#import-blank

**检测**：grep 正则 ^\s*(import\s+)?_\s+"（命中处核对所在文件是否为 main 包或 _test 文件）

### 【SHOULD】PKG-007 只在包名冲突或与路径末段不符时才起 import 别名。

- 归属：依赖与工程结构/包组织规约
- 起始版本：Go 1.0
- 主题：naming

冲突时改本地或项目内的那个导入，保留上游的规范名。别名本身要遵守包名规则：全小写、不带下划线。

**为什么**

> Avoid renaming imports except to avoid a name collision; good package names should not require renaming. In the event of collision, prefer to rename the most local or project-specific import.
>
> —— https://go.dev/wiki/CodeReviewComments#imports

别名是额外一层映射，读代码时要在别名与原包名之间换算，同一路径在不同文件里起了不同别名更无法按名字检索；只在真冲突时改名，需要记的映射最少。

**正例**

```go
import (
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	"github.com/go-playground/validator/v10"
)
```

**反例**

```go
import (
	d "github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	v "github.com/go-playground/validator/v10"
)
```

**依据**

- https://go.dev/wiki/CodeReviewComments#imports

**检测**：人工核对（核对每个别名是否用于解决同名冲突或包名与路径不符）

## （三）依赖管理规约

### 【MUST】DEP-001 go.mod 与 go.sum 交给 go 命令维护，提交前运行 go mod tidy。

- 归属：依赖与工程结构/依赖管理规约
- 起始版本：Go 1.11
- 主题：tooling

增删依赖用 `go get`，改 require 与 replace 用 `go mod edit`，整理用 `go mod tidy`，go.sum 只由工具写入。手工改版本号或校验和会与依赖图、模块缓存对不上。

**为什么**

> To keep your managed dependency set tidy, use the `go mod tidy` command. Using the set of packages imported in your code, this command edits your go.mod file to add modules that are necessary but missing. It also removes unused modules that don't provide any relevant packages.
>
> —— https://go.dev/doc/modules/managing-dependencies#synchronizing

依赖图里含大量间接依赖，人手工维护只能管到直接依赖那一层；`go mod tidy` 按实际 import 反推，缺的补上、不再用的删掉，go.sum 随之校正，构建结果才与源码一致。

**正例**

```go
$ go get github.com/spf13/cobra@v1.10.2
$ go mod tidy
```

**反例**

```go
// adding a dependency: hand-write a line into both go.mod and go.sum
require github.com/spf13/cobra v1.10.2
```

**依据**

- https://go.dev/doc/modules/managing-dependencies#synchronizing
- https://go.dev/doc/modules/gomod-ref

**检测**：人工核对（提交前运行 go mod tidy，确认 go.mod 与 go.sum 无未提交改动）

### 【MUST】DEP-002 go.sum 提交进版本库，不写进 .gitignore。

- 归属：依赖与工程结构/依赖管理规约
- 起始版本：Go 1.11
- 主题：tooling

go.mod 与 go.sum 一同提交。私有模块拉取失败时用 GOPRIVATE、GONOSUMDB 一类环境变量放开校验范围，不删除也不忽略 go.sum。

**为什么**

> Include the go.mod and go.sum files in your repository with your code.
>
> —— https://go.dev/doc/modules/managing-dependencies#enable_tracking

go.sum 记录每个模块内容的校验和，是拉取时的完整性依据；不入库时 CI 与同事机器各自重新计算，上游标签被移动或内容被改都发现不了，构建结果也无法跨机器比对。

**正例**

```go
// repository root
go.mod
go.sum
```

**反例**

```go
// .gitignore
go.sum
```

**依据**

- https://go.dev/doc/modules/managing-dependencies#enable_tracking

**检测**：人工核对（核对仓库含 go.sum 且 .gitignore 未忽略它）

### 【SHOULD】DEP-003 不长期保留指向本地目录或私有 fork 的 replace。

- 归属：依赖与工程结构/依赖管理规约
- 起始版本：Go 1.11
- 主题：tooling

本地联调可用 replace 指到同级目录或自己的 fork，合入前移除。确需长期替换时，把原因与移除条件写进 go.mod 注释，并定期复查上游是否已合并。

**为什么**

> When you use the replace directive, Go tools don't authenticate external modules as described in Adding a dependency.
>
> —— https://go.dev/doc/modules/managing-dependencies#local_directory

replace 会让被替换模块跳过校验和验证，来源不再可核对；本地路径只存在于本人机器上，别人拿到这份 go.mod 直接构建失败，指向 fork 则长期落后于上游修复。

**正例**

```go
// send the needed change upstream first, then bump the version after it is merged
require example.com/lib v1.4.0
```

**反例**

```go
// keep a modified copy locally and pin it with replace long-term
require example.com/lib v1.4.0

replace example.com/lib => ../lib
```

**依据**

- https://go.dev/doc/modules/managing-dependencies#local_directory
- https://go.dev/doc/modules/gomod-ref#replace

**检测**：人工核对（核对 go.mod 里的 replace 未长期指向本地目录，且写明了移除条件）

### 【SHOULD】DEP-004 工具依赖写进 go.mod 的 tool 指令，不靠全局安装。

- 归属：依赖与工程结构/依赖管理规约
- 起始版本：Go 1.24
- 主题：tooling

命令行工具用 tool 指令声明包路径，用 `go tool <name>` 运行，版本由 go.mod 与 go.sum 固定。全局 `go install` 得到的版本随机器与执行时间变化，CI 与本地无法对齐。

**为什么**

> Adds a package as a dependency of the current module, and makes it available to run with `go tool` when the current working directory is within this module.
>
> —— https://go.dev/doc/modules/gomod-ref#tool

全局安装的版本由执行 install 时的 @latest 决定，同一命令在不同机器上可能命中不同版本，检查结果无法复现；写进 tool 指令后版本进入依赖图，谁都跑同一份。

**正例**

```go
// go.mod
tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint

$ go tool golangci-lint run ./...
```

**反例**

```go
$ go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
$ golangci-lint run ./...
```

**依据**

- https://go.dev/doc/modules/gomod-ref#tool
- https://go.dev/doc/modules/managing-dependencies#tools

**检测**：人工核对（核对仓库用到的命令行工具经 go.mod 的 tool 指令声明）

### 【SHOULD】DEP-005 标准库已能覆盖的小功能，不引入第三方依赖。

- 归属：依赖与工程结构/依赖管理规约
- 起始版本：Go 1.0
- 主题：tooling

一行代码可写、标准库已有等价能力的，自己写；通用且非平凡的能力（CLI 框架、测试断言、结构化日志）确有收益再引。引入前先确认标准库没有现成能力。

**为什么**

> A little copying is better than a little dependency.
>
> —— https://go-proverbs.github.io/

每引一个模块，它的发版节奏、安全修复与 API 变更都进入自己的维护面；为一个小函数引入整个模块，收益通常低于日后跟进版本升级的成本。

**正例**

```go
n, err := strconv.Atoi(s)
```

**反例**

```go
n, err := convert.ToInt(s) // pull in the convert module for a single type conversion
```

**依据**

- https://go-proverbs.github.io/

**检测**：人工核对（核对新增依赖是否为标准库已能覆盖的小功能）

## （四）项目布局规约

### 【MUST】LAYOUT-001 module 根目录放 go.mod 与源码，不套 src/ 一层。

- 归属：依赖与工程结构/项目布局规约
- 起始版本：Go 1.11
- 主题：layering

导入路径从 module 根起算，源码与其上层目录在同一层级展开。GOPATH 时代或 Java 习惯带来的 src/ 会让仓库路径与 import path 之间多出一段无意义前缀。

**为什么**

> Some Go projects do have a `src` folder, but it usually happens when the devs came from the Java world where it's a common pattern. If you can help yourself try not to adopt this Java pattern.
>
> —— https://github.com/golang-standards/project-layout#directories-you-shouldnt-have

go 命令以 go.mod 所在目录为模块根，import path 由该目录往下的相对路径拼成；根下再套 src/，所有导入、`go install` 的包路径与文档示例都要多写一层，还与仓库实际目录对不上。

**正例**

```go
project-root/
  go.mod
  main.go
  internal/
```

**反例**

```go
project-root/
  src/
    go.mod
    main.go
```

**依据**

- https://github.com/golang-standards/project-layout#directories-you-shouldnt-have
- https://go.dev/doc/modules/layout

**检测**：人工核对（核对 go.mod 所在目录即仓库根，源码未放在 src/ 之下）

### 【SHOULD】LAYOUT-002 不打算对外暴露的包放进 internal/ 下。

- 归属：依赖与工程结构/项目布局规约
- 起始版本：Go 1.4
- 主题：layering

internal 目录的可见范围是它所在目录的子树，仓库只有一个 module 时顶层 internal/ 覆盖全仓库。实现细节、适配器与仅本程序使用的类型都放进去；确需对外复用的包才放在 internal 之外。

**为什么**

> Initially, it's recommended placing such packages into a directory named `internal`; this prevents other modules from depending on packages we don't necessarily want to expose and support for external uses. Since other projects cannot import code from our `internal` directory, we're free to refactor its API and generally move things around without breaking external users.
>
> —— https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages

放在 internal 外的包会被别的模块按固定路径 import，一旦有人用上，改包名、改签名都成了破坏性变更；放进 internal 后编译器直接拦住外部导入，内部结构可以自由调整。

**正例**

```go
// external modules cannot import it
import "github.com/ForeverSRC/golang-dev-manual/gdm/internal/repository/jsonfile"
```

**反例**

```go
// the implementation package sits at the module root, so its API is locked once another module depends on it
import "github.com/ForeverSRC/golang-dev-manual/repository/jsonfile"
```

**依据**

- https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages
- https://go.dev/doc/go1.4#internalpackages

**检测**：人工核对（核对仅本模块使用的包是否位于 internal/ 下）

### 【SHOULD】LAYOUT-003 可执行程序集中放在 cmd/ 下，一个程序一个子目录。

- 归属：依赖与工程结构/项目布局规约
- 起始版本：Go 1.0
- 主题：layering

同一仓库既有可导入的包又有可执行程序时，main 包统一收进 cmd/<程序名>/，目录名与生成的可执行文件同名。只有单个程序、且仓库里没有可导入包时，可以把 main.go 放在模块根。

**为什么**

> A common convention is placing all commands in a repository into a `cmd` directory; while this isn't strictly necessary in a repository that consists only of commands, it's very useful in a mixed repository that has both commands and importable packages, as we will discuss next.
>
> —— https://go.dev/doc/modules/layout#multiple-commands

main 包混在模块根或业务目录里，import 路径与可执行文件名对不上，`go install` 时要先找文件所在目录；集中到 cmd/ 后，包与命令在结构上分开，每个程序的入口一目了然。

**正例**

```go
project-root/
  cmd/
    gdm/
      main.go
  internal/
    service/
```

**反例**

```go
project-root/
  gdm.go
  internal/
    server/
      main.go
```

**依据**

- https://go.dev/doc/modules/layout#multiple-commands
- https://github.com/golang-standards/project-layout#cmd

**检测**：人工核对（核对 main 包集中在 cmd/ 下，目录名与可执行文件同名）

### 【MAY】LAYOUT-004 项目还小时不预先铺设大型项目的目录骨架。

- 归属：依赖与工程结构/项目布局规约
- 起始版本：Go 1.11
- 主题：layering

cmd、internal、pkg、configs、scripts 在确有对应内容时再建，空目录或只有一个占位文件的目录不建。目录随代码增长而增加，名字取自真实内容。

**为什么**

> If you are trying to learn Go or if you are building a PoC or a simple project for yourself this project layout is an overkill. Start with something really simple instead (a single `main.go` file and `go.mod` is more than enough).
>
> —— https://github.com/golang-standards/project-layout#overview

空目录不承载信息，却让人以为存在对应的结构约定，后来者会往预留的 pkg/、util/ 里塞本可归入具体包的东西；等真实职责出现再建，目录名才能反映内容。

**正例**

```go
project-root/
  go.mod
  main.go
```

**反例**

```go
project-root/
  cmd/
  pkg/
  configs/
  scripts/
  internal/
  // each directory holds only a .gitkeep
```

**依据**

- https://github.com/golang-standards/project-layout#overview

**检测**：人工核对（核对仓库目录是否都有实际内容与明确职责）

