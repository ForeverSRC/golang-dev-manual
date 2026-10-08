# gdm 子树规则

> 父项目规则见 `../AGENTS.md`。本文件只约束 `gdm/` 子树，不重复项目级规则。
> 改条款数据或写条款前先读 `./CLAUSES.md`，上游依据从 `../REFERENCE.md` 挑。

## 作用范围

`gdm/` 是手册的 Go 工具链与数据源：`gdm/data/` 存唯一事实来源，`gdm-cli` 对外检索与检测条款，`gdm-gen` 把数据渲染成 `../manual/` 下的双语 markdown。

## 目录约定

| 目录 | 角色 |
|------|------|
| `gdm/data/` | 条款数据与 `//go:embed` 载体，编进两个命令的二进制 |
| `gdm/di/` | wire 装配，端口与实现的绑定集中在 `ProviderSet`；`RealSet` 给生产入口提供内嵌数据树，`ITSet` 给集成测试提供 fixture 数据树 |
| `gdm/internal/ittest/` | 集成测试的 fixture 数据树，`//go:embed testdata` 后按 `fs.FS` 提供 |
| `gdm/cmd/gdm-cli/`、`gdm/cmd/gdm-gen/` | 两个入口，各自的 `wire/` 放生产注入器与容器，集成测试的注入器放同级 `wireit/`（当前只有 gdm-cli） |
| `gdm/internal/domain/` | 领域类型与领域方法，不做 IO |
| `gdm/internal/repository/jsonfile/` | 读条款与语言覆盖层 |
| `gdm/internal/adapter/filesystem/` | 扫描 Go 源码、写手册文件 |
| `gdm/internal/service/` | 应用服务，跨层端口接口定义在 `ports.go` |
| `gdm/internal/api/clihandler/` | 入站适配，两个命令共用 |
| `gdm/internal/server/` | cobra 命令树 |

依赖单向：`cmd`、`server` -> `api` -> `service` -> `repository`、`adapter` -> `domain`。新增一个能力先加 `service/ports.go` 的接口与 `di` 的绑定，再跑 `make generate`。

## 不要放进这里

- 条款正文与译文 -> `gdm/data/`
- 手册 markdown -> `../manual/`（生成物，禁止手工编辑）
- 上游依据清单 -> `../REFERENCE.md`
- 手册的使用工作流 -> `../skills/gdm-workflows/`
