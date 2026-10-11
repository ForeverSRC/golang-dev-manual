# review 流程

对一次代码改动做手册级 review 的分步操作，每步给出可直接执行的命令。

## 步骤 1：看主题词表，锁定 tag

```bash
gdm-cli tags
```

- 输出横向主题词表，每个主题给一句话解释，例如 `concurrency` 管 goroutine 退出、channel、同步原语与结构体拷贝
- 按改动内容挑出相关主题：goroutine、channel、锁相关改动取 `concurrency`，测试代码改动取 `testing`
- 一次可传多个主题，逗号分隔，命中任一即可
- 拿不准某个主题词管什么时，以 `tags` 输出的一句话解释为准，不要靠猜

## 步骤 2：拉取候选条款

按主题过滤：

```bash
gdm-cli list --tag concurrency
```

只核对强制条款时补 `--level`，与 `--tag` 叠加取交集：

```bash
gdm-cli list --tag concurrency --level MUST
```

改动里有具体关键词时用搜索，结果按相关度排序：

```bash
gdm-cli search context 取消
```

- 输出只有 `编号: 一句话`，例如 `CONC-002: channel 只由发送方关闭，接收方不关闭。`
- `--level` 逗号分隔可一次取多级，如 `--level MUST,SHOULD`
- 编号前缀对应章节，如 `CONC` 是并发、`ERR` 是错误处理

## 步骤 3：逐条核对

- 拿候选清单对照当前改动，判断每条是否命中
- 命中疑点的记下编号，此步不展开任何条款详情

## 步骤 4：展开疑点

只对疑点编号调 `explain`，支持一次传多个编号：

```bash
gdm-cli explain CONC-001 CONC-002
```

- 输出含分级、归属、起始版本、说明、为什么、正例、反例与依据链接
- 对照正反例确认：确认违反的进报告，确认不违反的丢弃
- 没有疑点时 review 到此结束

## 步骤 5：机器检测（可选）

```bash
gdm-cli check --verbose <改动目录>
```

- 输出分两段：命中项在前，`--verbose` 追加的未集成自动检测清单在后
- 未检测清单逐条标注检测方式，如 `golangci-lint stylecheck(ST1003)` 或 `人工核对`
- 命中项定位到具体代码位置复核，确认违反的进报告
- 未检测清单里与本次改动相关的条款，回到步骤 3、4 人工核对

## 步骤 6：报告

- 每条问题给出编号与条款一句话，例如 `ERR-001: 需要保留错误链时用 %w 包装，不使用 %v 或 %s。`
- 问题与编号一一对应，不合并同类项
- 凭编号可用 `explain` 回到条款详情查证
