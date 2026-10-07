---
name: gdm-workflows
description: 用 gdm-cli 按 Golang 开发手册的条款核对代码、查询某条 Go 规约。当需要按级别拉取 Go 条款清单、展开某条规约的正反例、或对 Go 代码做手册级 review 时触发。
allowed-tools:
  - Bash
version: 0.1.0
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
  - `list`：按级别或归属过滤，只输出 `编号: 一句话`。
  - `explain <编号>`：展开单条条款的详情、正反例与依据。
  - `check <路径>`：对代码目录执行可自动化的检测，输出命中项。
- 语言：输出语言跟随 `--lang`，默认 `zh`；使用者用英文时传 `--lang en`，条款正文与字段标签一起切换到英文。
- 通用选项：`--level`、`--category` 过滤 `list`，`--category` 取归属 id，如 `programming-conventions` 或 `programming-conventions/naming`；`--verbose` 让 `check` 列出未集成自动检测的条款。选项写在位置参数之前。

## review 流程

1. `list --level MUST` 拿到强制条款短清单。
2. 对当前改动逐条核对，命中疑点的记下编号。
3. 只对疑点编号调 `explain`，展开详情与正反例。
4. 需要机器检测时调 `check <改动目录> --verbose`，按命中项定位。
5. 报告问题时给出编号与条款一句话，便于回到手册查证。

## 纪律

- `list` 拿清单后不要 `explain` 全部条款，只展开疑点，避免整本手册进上下文。
- 逐条核对时保留编号，问题与编号一一对应。
- 手册只覆盖语言级条款，项目专属约定不在其中。
