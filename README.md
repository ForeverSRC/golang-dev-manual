# Golang 开发手册

[English](README.en.md) · [在线阅读](https://foreversrc.github.io/golang-dev-manual/)

![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![Go 1.27.1](https://img.shields.io/badge/Go-1.27.1-00ADD8.svg)

个人维护的 Golang 开发手册：9 个章节、113 条语言级条款，分 `MUST` / `SHOULD` / `MAY` 三级，每条带一句话、说明、依据讲解、上游引文、正反例与参考链接。条款都可客观判定，能直接转成 lint 规则或 review 项。手册提供中英两版 markdown，并以静态站点发布。

## 设计初衷

- 把日常开发里散落的 Go 编码经验整理成一份可核对的规范。
- 形态参照《阿里巴巴 Java 开发手册》：三级分级、编号条款、正反例。
- 使用者是作者本人与其 Agent，不做组织级推行。
- 对人友好：可以当学习材料逐章阅读。
- 对 Agent 友好：review 代码时能低成本拿到三级条款清单、逐条核对，上下文消耗小。这一条是重点。
- 在 AI 编程时代仍然有用，原因有两个：
  - 纠正 AI 的默认偏差：手册里有一批条款，是在纠正 AI 写的代码时补进去的，错误断言比对了错误文本、断言把结果投影成子集、导出与未导出的次序混乱、给代码本身能读出的东西加注释。这些偏差都不难判定，但 AI 默认会犯，人也不一定每次都盯得住；固化成有条目、可检测的条款，比每次靠临时提醒可靠。
  - 让 AI 吸收最新的程序写法、遵循最佳实践写代码：这与 go fix 和 modernize 两条条款的用意一致。

## 与其他规范的关系

- 本手册是在 Effective Go、Go Code Review Comments、Google Go Style、Uber Go Style Guide 基础上按个人取舍整理的条款清单。
- 收录标准是短、可判定、有直接对应的上游依据。
- 本手册未覆盖的内容，以上游那几份规范为准。

## 快速开始

安装 CLI：

```bash
go install github.com/ForeverSRC/golang-dev-manual/gdm/cmd/gdm-cli@latest
```

```bash
gdm-cli list --level MUST     # 每条条款一行
gdm-cli explain NAMING-001    # 展开详情、正反例与依据
gdm-cli check ./some/package  # 执行可自动化的检测
```

条款数据已编入二进制，装好的 `gdm-cli` 不需要 clone 仓库。

配套 skill `skills/gdm-workflows/` 把 CLI 包装成 Agent 工作流，可用 `npx skills` 安装：

```bash
npx skills add ForeverSRC/golang-dev-manual
```

## 从源码构建

贡献者从签出的仓库出发：

```bash
make build     # 构建 bin/gdm-cli（对外 CLI）与 bin/gdm-gen（手册生成器）
make gen       # 由 gdm/data/ 重新生成 manual/zh/ 与 manual/en/
make check     # go vet + go test
make lint      # golangci-lint
make generate  # 重新生成 wire 装配代码
```

不构建时可以直接 `go run ./gdm/cmd/gdm-cli <子命令>`。

## CLI 用法

输出语言由 `--lang` 决定，取值 `zh`（默认）与 `en`，条款正文与字段标签一起切换。

`gdm-cli list` 每条条款输出一行，先读全量清单再决定展开哪几条。

```console
$ gdm-cli list --level MUST
NAMING-001: 包名使用小写单词连写，不使用下划线、驼峰或复数。
NAMING-002: 标识符不使用下划线分隔单词。
```

`gdm-cli explain` 展开一条或多条条款的详情、正反例与依据。

```console
$ gdm-cli explain NAMING-001
【MUST】NAMING-001 包名使用小写单词连写，不使用下划线、驼峰或复数。

归属: 编程规约/命名规约 | 起始版本: 1.0

说明:
包名是调用方引用标识符的前缀，全小写、无分隔符才能保持一致。……

为什么:
引文: In Go, package names must be concise and use only lowercase letters and numbers……
出处: https://google.github.io/styleguide/go/decisions#package-names

正例: 接受的写法，一段 Go 代码
反例: 被拒的写法，与正例只差一处

依据:
- https://go.dev/wiki/CodeReviewComments#package-names

检测: golangci-lint stylecheck(ST1003)（检查包名标识符）
```

`gdm-cli check` 对代码目录执行基于 grep 的检测，输出命中项。

```console
$ gdm-cli check ./gdm --verbose
未发现命中。

以下条款没有自动检测：
  [NAMING-001] golangci-lint stylecheck(ST1003)（检查包名标识符）
```

`list` 另接受：

- `--level`：按级别过滤，如 `MUST,SHOULD`。
- `--category`：按归属 id 过滤，如 `programming-conventions` 或 `programming-conventions/naming`；取值非法时报错并列出全部可选 id。

## 仓库结构

- `manual/zh/`、`manual/en/`：生成的双语 markdown 手册，已入库并直接作为站点数据源，请勿手工修改。
- `gdm/`：Go 工具链。
  - `gdm/data/`：唯一事实来源。`gdm/data/data.go` 把它编入 `gdm-cli` 与 `gdm-gen` 二进制。
    - `gdm/data/manual.json`：版本、Go 基准与目录结构。
    - `gdm/data/clauses/<章节 id>.json`：按章节拆分的条款，一个章节一个文件。
    - `gdm/data/i18n/<lang>/`：语言覆盖层，中文是事实来源，其余语言挂翻译。
  - `gdm/cmd/gdm-cli`：对外发布的 CLI，分为 `internal/domain`、`internal/service`、`internal/repository`、`internal/adapter`、`internal/api`、`internal/server`，装配在 `di/`。
  - `gdm/cmd/gdm-gen`：手册生成器，只由维护者与 CI 构建。
  - `gdm/CLAUSES.md`：条款的数据格式、写作口径与自检维护。
- `REFERENCE.md`：条款依据的上游材料清单。
- `skills/gdm-workflows/`：配套 skill，把 CLI 包装成 Agent 工作流。
- `website/`：静态站点的 MkDocs 配置。

## 条款格式

| 字段 | 含义 |
| --- | --- |
| `id` | 条款编号，`前缀-三位序号`，全局唯一 |
| `level` | `MUST` / `SHOULD` / `MAY` |
| `chapter` | 章节 id，必须能在 `toc` 中找到 |
| `section` | 小节 id，可省略，省略时归属整个章节 |
| `since_go` | 条款可落地的最低 Go 版本 |
| `summary` | 一句话条款，`gdm-cli list` 只输出它 |
| `details` | 适用范围与例外 |
| `rationale` | 为什么这么定 |
| `quote` | 上游原文摘录与出处 |
| `examples` | `good` 与 `bad` 最小片段 |
| `sources` | 上游依据链接 |
| `detect` | 检测方式，可省略 |

分级口径：

- `MUST`：违反即错，应能被自动检测或机械核对。
- `SHOULD`：默认遵循，例外需写明理由。
- `MAY`：提供方向与惯用法，由人判断是否适用。

`detect.tool` 有三种取值：

- `golangci-lint`：规则名写在 `detect.rule`。
- `grep`：正则在 `detect.pattern`，由 `gdm-cli check` 扫描 `.go` 文件时使用。
- `manual`：人工核对，具体内容写在 `detect.note`。

## 如何贡献

- 条款有缺陷、缺条款、参考链接失效或归属有误，用 issue 说明。
- 增改条款的 PR 只动 `gdm/data/`，并在同一提交里执行 `make gen` 重新生成 `manual/`。
- 提交前执行 `make check` 与 `make lint`。
- 项目专属约定不在范围内，组织级推行也不在范围内。

## 许可与致谢

- 采用 MIT 许可，覆盖 CLI 代码与条款内容。
- 条款里逐字引用的上游原文版权归各出处，已在 `quote.source` 与手册附录中标明。
- 本手册依据的上游规范：Effective Go、Go Code Review Comments、Google Go Style、Uber Go Style Guide。
