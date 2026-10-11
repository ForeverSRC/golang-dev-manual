---
name: gdm-workflows
description: 用 gdm-cli 按 Golang 开发手册的条款核对代码、按主题检索条款、展开某条 Go 规约。当需要按级别或主题拉取 Go 条款清单、按主题搜索条款、展开某条规约的正反例，或对 Go 代码做手册级 review 时触发。
allowed-tools:
  - Bash
version: 0.2.0
---

# gdm-cli 工作流

`gdm-cli` 是 Golang 开发手册的检索与检测入口，条款数据内嵌在二进制里。手册正文在线发布在 <https://foreversrc.github.io/golang-dev-manual/>。

## 安装前置

```bash
go install github.com/ForeverSRC/golang-dev-manual/gdm/cmd/gdm-cli@latest
```

- 装到 `GOBIN`（默认 `$(go env GOPATH)/bin`），确认该目录在 `PATH` 上。
- 条款数据内嵌在二进制里，装完即用，不需要 clone 仓库。

## 命令

- 调用方式：`gdm-cli <子命令>`
- 子命令：
  - `list`：按级别、归属或主题过滤，只输出 `编号: 一句话`。
  - `search <查询词>`：按主题搜索条款，按相关度排序输出，同样是 `编号: 一句话`。
  - `tags`：列出横向主题词表，每个 tag 给一句话解释。
  - `explain <编号>`：展开单条条款的详情、正反例与依据。
  - `check <路径>`：对代码目录执行可自动化的检测，输出命中项。
- 语言：输出语言跟随 `--lang`，默认 `zh`；使用者用英文时传 `--lang en`，条款正文与字段标签一起切换到英文。
- 过滤选项：`--level` 取 `MUST`/`SHOULD`/`MAY`，`--category` 取归属 id（如 `programming-conventions` 或 `programming-conventions/naming`），`--tag` 取横向主题 id（逗号分隔、命中任一即可，如 `concurrency,testing`）；三者叠加取交集。`list` 与 `search` 都接受它们。
- 主题取值：不确定有哪些 tag、或看不出某个 tag 管什么时，先跑 `gdm-cli tags`，它按当前语言给出每个 tag 的一句话解释；`--tag` 与 `search` 的取值都来自这份词表。传非法值时（如 `list --tag x`）报错也会列出全部可选 tag。
- 机器可读输出：`list`、`search`、`tags`、`check` 加 `--json` 输出单个 JSON 文档，`search` 另带 `score` 与 `matched`。
- 其他：`--verbose` 让 `check` 列出未集成自动检测的条款。选项写在位置参数之前。

## review 流程

对当前改动做手册级 review 时，按 [references/review-workflow.md](references/review-workflow.md) 的分步命令执行：先缩小候选条款，逐条核对后只展开疑点，需要时机器检测。

## 纪律

- `list` 或 `search` 拿候选后不要 `explain` 全部条款，只展开疑点，避免整本手册进上下文。
- 逐条核对时保留编号，问题与编号一一对应。
- 手册只覆盖语言级条款，项目专属约定不在其中。
