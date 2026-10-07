# 条款规范

本文件是 `gdm/data/` 条款的数据格式与写作规范。改条款数据或新增条款前先读它，写完按 §3 自检。上游依据清单见 `../REFERENCE.md`。

`gdm/data/` 是手册的单一数据源。人读的 `manual/` 与 Agent 用的 `gdm-cli` 都由它生成，两处内容不会分叉。

## 1 数据格式

### 文件划分

| 路径 | 内容 |
| --- | --- |
| `gdm/data/manual.json` | 顶层字段与目录 |
| `gdm/data/clauses/<章节 id>.json` | 一个章节一个文件，内容是该章节的条款数组，文件名取 `toc` 的 `id` |

- 加载时按 `toc` 顺序拼接各章节文件，条款的全局顺序即 toc 顺序。
- 新增章节：先在 `manual.json` 的 `toc` 加一项并给出 `id` 与 `name`，再建 `clauses/<id>.json`。
- `gdm/data/data.go` 以 `//go:embed manual.json clauses i18n` 把数据编进 `gdm-cli` 与 `gdm-gen` 二进制，两个命令都只读内嵌数据，运行时不依赖工作区。`go:embed` 的模式不能向上跨目录，只能显式列出 `gdm/data/` 的顶层条目，新增顶层条目时同步这一行。

### 语言覆盖层（gdm/data/i18n/<lang>/）

中文是事实来源，其他语言以覆盖层挂在它之上；来源语言也有一份覆盖层，只放文案。

| 路径 | 内容 |
| --- | --- |
| `gdm/data/i18n/<lang>/manual.json` | 索引页文案、附录文案、渲染标签、CLI 标签，非来源语言另含 `toc` 译文 |
| `gdm/data/i18n/<lang>/clauses/<章节 id>.json` | 该章节条款的 `summary`、`details`、`rationale` 与 `detect.note` 翻译，按条款 id 索引 |

- `toc` 译文按章节 id 与 `gdm/data/manual.json` 的 `toc` 挂接，每项给章节名与该章小节名，不依赖数组下标。
- `render` 里的模板字段用 `{v}` 作占位符，如 `"baseline_line": "基准版本：Go {v}"`。
- `render.number_style` 取 `chinese`（章节用「一、」）或 `arabic`（章节用「1.」）。
- `cli` 块给 `gdm-cli` 的字段标签：`details`、`why`、`quote`、`source`、`good`、`bad`、`references`、`detection`、`category`、`since_go`、`no_hits`、`uncovered`。
- 归属不落译文：条款用 `chapter` / `section` 两个 id 引用目录，展示值按语言的 `toc` 译文推导。
- 生成非来源语言时，`toc` 或条款翻译缺失即失败，报错列出缺失的章节、小节与条款 id。
- 来源语言不需要 `clauses/` 目录；其余语言的 `clauses/` 必须覆盖该章节全部条款。

### 顶层字段（gdm/data/manual.json）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `version` | string | 数据格式版本，格式变更时递增 |
| `go_baseline` | string | 本版条款的 Go 基准版本 |
| `toc` | array | 手册目录，决定章节顺序与手册文件划分 |

### toc

- 数组顺序即手册章节顺序。
- `id`：语言无关的 ASCII 标识（kebab-case），是 `clauses/<id>.json` 的文件名，也是条款归属引用的对象。
- `name`：来源语言的章节展示名，其他语言的展示名在覆盖层里按 `id` 给。
- `sections`：小节数组，可省略；省略时该章节不分子节。小节的 `id` 取简短 ASCII 短名，一次定下不随译文措辞调整，`name` 是来源语言的展示名。
- 没有 `sections` 的章节，条款只写 `chapter`。

### clauses

`gdm/data/clauses/<章节 id>.json` 是一个 JSON 数组，元素字段如下。

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | string | 是 | 条款编号，格式 `前缀-三位序号`，如 `NAMING-001`；全局唯一，同级别内按它升序 |
| `level` | string | 是 | 分级，取值 `MUST` / `SHOULD` / `MAY` |
| `chapter` | string | 是 | 所属章节 id，必须能在 `toc` 中找到 |
| `section` | string | 否 | 所属小节 id，必须能在该章节的 `sections` 中找到；省略表示归属整个章节 |
| `since_go` | string | 是 | 条款按本手册写法可落地的最低 Go 版本，取语言特性或标准库 API 的引入版本，如 `1.25`（sync.WaitGroup.Go）；第三方库的版本要求不计；无版本依赖的条款写 `1.0` |
| `summary` | string | 是 | 一句话条款，`list` 只输出它，须短且可判定 |
| `details` | string | 是 | 条款说明，写清边界与例外 |
| `rationale` | string | 是 | 为什么这么定，讲解 `quote` 引文，落到具体代价或踩坑 |
| `quote` | object | 是 | 上游原文摘录与出处，见下 |
| `examples` | object | 是 | 正反例，见下 |
| `sources` | array | 是 | 上游依据链接，只放与条款直接对应的官方地址 |
| `detect` | object | 否 | 检测方式，见下；省略表示暂不检测 |

### 手册内排序

- 小节内先按 `level` 排（MUST -> SHOULD -> MAY），同级别内按 `id` 升序。
- 章节顺序与小节顺序由 `toc` 决定，与 `level` 无关。
- 附录的总索引、依据来源与检测映射按 `id` 升序。

### quote

| 字段 | 说明 |
| --- | --- |
| `text` | 依据资料里的原文摘录，一到两句，逐字引用不改写 |
| `source` | 摘录出处，取 `sources` 里对应的那条链接 |

- 摘录只取能直接支撑条款的那一两句，不做大段引用。
- `rationale` 紧接着引文，用中文讲解它，二者在手册里合成「为什么」一节。

### examples

| 字段 | 说明 |
| --- | --- |
| `good` | 正例，最小可运行片段 |
| `bad` | 反例，与正例对照，只差被禁止的写法 |
| `lang` | 可选，示例代码块的语言标记，取值如 `go`、`yaml`、`makefile`、`console`；省略时按 `go` 渲染 |

### detect

| 字段 | 说明 |
| --- | --- |
| `tool` | 检测工具，取值 `grep` / `golangci-lint` / `manual` |
| `rule` | `golangci-lint` 时的规则名，如 `errorlint` |
| `pattern` | `grep` 时的正则，`gdm-cli check` 用它扫描 `.go` 文件 |
| `note` | 补充说明，如需要人工核对的部分 |

## 2 条款写作

### 分级口径

| 级别 | 含义 | 判定方式 |
| --- | --- | --- |
| `MUST` | 违反即错 | 应能被自动检测，或能被机械核对 |
| `SHOULD` | 默认遵循 | 例外需在代码或评审里写明理由 |
| `MAY` | 提供方向与惯用法 | 由人判断是否适用 |

### 字段写作要求

- 行内代码
  - 正文里出现的包名、标识符、类型名、方法名、命令名与文件名一律用反引号包住，示例代码块内不再包。
  - 反引号前后与中文之间留半角空格。
  - 英文译文里按普通英文词使用时（error、context、channel 一类）不加反引号。
- `summary`
  - 一句话说清要求，不超过 30 字。
  - 只写要求，不写理由。
  - 能被读到的人直接判定"是否违反"。
- `details`
  - 写清适用范围与例外。
  - 不重复 `summary` 的原话。
  - 需要罗列多项（检查器、场景、例外等）时用无序列表，一条一项；不要挤成一段长句。
- `rationale`
  - 讲清 `quote` 引文为什么支撑这条条款，落到具体代价或踩坑。
  - 禁止 "为了代码质量" 一类空话。
- `quote`
  - `text` 逐字引用依据资料原文，一到两句，够支撑条款即可，不做大段引用。
  - `source` 取 `sources` 里与引文对应的那条链接。
- `examples`
  - `good` 与 `bad` 只差被禁止的那一处写法，方便对照。
  - 用最小可运行片段，不引入无关业务。
- `sources`
  - 从 `../REFERENCE.md` 的清单里挑与条款直接对应的上游官方深链，清单外的不直接用。
  - 优先官方文档与官方博客，少放二手翻译。
- `detect`
  - 能交给工具就交给工具：能用 lint 写 `golangci-lint`，能用简单匹配写 `grep`。
  - 工具覆盖不到的部分用 `note` 说明，并考虑退回 `manual`。

### 先确认依据

1. 打开候选的上游原文，确认链接可用且直接对应条款。
2. 确认条款与项目无关、可被客观判定。
3. 找不到直接依据的，先不写，记为待查。

### JSON 模板

```json
{
  "id": "前缀-001",
  "level": "MUST",
  "chapter": "章节 id",
  "section": "小节 id",
  "since_go": "1.0",
  "summary": "一句话要求。",
  "details": "适用范围与例外。",
  "rationale": "为什么这么定，讲解下面的引文。",
  "quote": {
    "text": "依据资料里的原文摘录",
    "source": "https://引文对应的官方链接"
  },
  "examples": {
    "good": "正确写法",
    "bad": "错误写法"
  },
  "sources": [
    "https://对应的官方链接"
  ],
  "detect": {
    "tool": "golangci-lint",
    "rule": "规则名"
  }
}
```

## 3 自检与维护

自检清单：

- 编号唯一，`chapter` 能在 `gdm/data/manual.json` 的 `toc` 中找到，`section` 能在该章节的 `sections` 中找到。
- 级别选对：违反即错用 `MUST`，例外可解释用 `SHOULD`，方向性用 `MAY`。
- `summary` 不含理由，`rationale` 不含要求原话。
- `quote.text` 逐字来自 `quote.source`，两处链接可打开。
- 正反例只差被禁止的写法。
- `sources` 链接可打开且直接对应条款。
- 每条 `MUST` 都能回答"怎么检测"。

维护流程：

1. 改条款改对应章节的 `gdm/data/clauses/<章节 id>.json`；改目录或版本改 `gdm/data/manual.json`；改译文改 `gdm/data/i18n/<lang>/`。
2. 执行 `make gen` 重新生成 `manual/zh/` 与 `manual/en/`。
3. 执行 `make check` 确认代码与测试通过。
4. 提交时同时带上数据与 `manual/` 的改动。
