# 十、附录

## 1 条款总索引

| 编号 | 级别 | 归属 | 一句话 |
| --- | --- | --- | --- |
| ABSTRACT-001 | SHOULD | 接口与设计规约/接口定义与落位 | 抽象按实际需要引入，不提前铺开。 |
| CODE-001 | SHOULD | 错误与日志/错误码规约 | 对外错误码沿用既有标准分类，不另建编号体系。 |
| CODE-002 | SHOULD | 错误与日志/错误码规约 | 错误消息不作为契约，判断错误只看码与类型。 |
| CODE-003 | SHOULD | 错误与日志/错误码规约 | 跨服务传播错误时转换错误码，不透传下游原始错误。 |
| COMMENT-001 | SHOULD | 编程规约/注释规约 | 注释只写代码表达不出的东西，不复述代码在做什么。 |
| COMMENT-002 | SHOULD | 编程规约/注释规约 | 文档注释以被描述的名称开头，写成完整句子。 |
| COMPOSE-001 | SHOULD | 接口与设计规约/组合与复用 | 复用能力通过组合表达，不搭继承式层次。 |
| COMPOSE-002 | SHOULD | 接口与设计规约/组合与复用 | 嵌入不改变外层的零值、拷贝语义与公开接口。 |
| CONC-001 | MUST | 编程规约/并发规约 | 启动 goroutine 前明确退出时机，并由启动方等待结束。 |
| CONC-002 | MUST | 编程规约/并发规约 | channel 只由发送方关闭，接收方不关闭。 |
| CONC-003 | SHOULD | 编程规约/并发规约 | 对外 API 保持同步语义，由调用方决定是否并发。 |
| CONC-004 | SHOULD | 编程规约/并发规约 | 共享状态与保护它的锁放在同一结构体，仅经方法读写。 |
| CONST-001 | MUST | 编程规约/常量与枚举规约 | 常量名用 MixedCaps，不用全大写或 k 前缀。 |
| CONST-002 | SHOULD | 编程规约/常量与枚举规约 | 常量按用途命名，不复刻它的取值。 |
| CONST-003 | MAY | 编程规约/常量与枚举规约 | 同源递增的枚举值用 iota 生成。 |
| CTRL-001 | SHOULD | 编程规约/控制语句规约 | 先处理错误与边界并提前返回，正常路径保持零缩进。 |
| CTRL-002 | SHOULD | 编程规约/控制语句规约 | 不再显式拷贝循环变量。 |
| CTX-001 | SHOULD | 接口与设计规约/context 传递 | 不把 context.Context 存进结构体，作为函数第一个参数传递。 |
| CTX-002 | MUST | 接口与设计规约/context 传递 | 不向函数传 nil context，占位用 context.TODO()。 |
| CTX-003 | SHOULD | 接口与设计规约/context 传递 | context.Value 的 key 用未导出的自定义类型。 |
| CTX-004 | SHOULD | 接口与设计规约/context 传递 | context.Value 只放请求作用域数据，不放可选参数。 |
| DATA-001 | SHOULD | 编程规约/数据结构规约 | 声明空切片时用 var 得到 nil 切片。 |
| DATA-002 | SHOULD | 编程规约/数据结构规约 | 结构体字面量用字段名初始化，不按位置罗列。 |
| DEP-001 | MUST | 依赖与工程结构/依赖管理规约 | go.mod 与 go.sum 交给 go 命令维护，提交前运行 go mod tidy。 |
| DEP-002 | MUST | 依赖与工程结构/依赖管理规约 | go.sum 提交进版本库，不写进 .gitignore。 |
| DEP-003 | SHOULD | 依赖与工程结构/依赖管理规约 | 不长期保留指向本地目录或私有 fork 的 replace。 |
| DEP-004 | SHOULD | 依赖与工程结构/依赖管理规约 | 工具依赖写进 go.mod 的 tool 指令，不靠全局安装。 |
| DEP-005 | SHOULD | 依赖与工程结构/依赖管理规约 | 标准库已能覆盖的小功能，不引入第三方依赖。 |
| DI-001 | SHOULD | 接口与设计规约/依赖注入与装配 | 依赖装配用 wire 在编译期生成。 |
| DI-002 | SHOULD | 接口与设计规约/依赖注入与装配 | 依赖经构造函数显式传入，不读写包级可变变量。 |
| ERR-001 | MUST | 错误与日志/错误处理规约 | 需要保留错误链时用 %w 包装，不使用 %v 或 %s。 |
| ERR-003 | MUST | 错误与日志/错误处理规约 | 判断错误种类用 errors.Is 与 errors.As，不用 == 或类型断言。 |
| ERR-004 | SHOULD | 错误与日志/错误处理规约 | 可预期的失败返回 error，不用 panic。 |
| ERR-005 | SHOULD | 错误与日志/错误处理规约 | 失败用 error 或 ok 返回，不用特殊值表示。 |
| ERR-006 | MUST | 错误与日志/错误处理规约 | 同一个错误只处理一次：包装返回或就地记录，不两者都做。 |
| FUNC-001 | MUST | 编程规约/函数与方法规约 | 导出的函数与方法放在文件前部，未导出的放在后部。 |
| FUNC-002 | SHOULD | 编程规约/函数与方法规约 | 返回值含义不自明时用命名返回值。 |
| IFACE-001 | SHOULD | 接口与设计规约/接口定义与落位 | 接口按实现数量落位：单实现与实现同包，多实现由使用方定义。 |
| IFACE-002 | SHOULD | 接口与设计规约/接口定义与落位 | 用编译期断言声明类型对接口的实现。 |
| IFACE-003 | SHOULD | 接口与设计规约/接口定义与落位 | 不定义指向接口的指针。 |
| LAYER-001 | SHOULD | 依赖与工程结构/分层规约 | 按职责分层，依赖方向单向指向内层。 |
| LAYOUT-001 | MUST | 依赖与工程结构/项目布局规约 | module 根目录放 go.mod 与源码，不套 src/ 一层。 |
| LAYOUT-002 | SHOULD | 依赖与工程结构/项目布局规约 | 不打算对外暴露的包放进 internal/ 下。 |
| LAYOUT-003 | SHOULD | 依赖与工程结构/项目布局规约 | 可执行程序集中放在 cmd/ 下，一个程序一个子目录。 |
| LAYOUT-004 | MAY | 依赖与工程结构/项目布局规约 | 项目还小时不预先铺设大型项目的目录骨架。 |
| LOG-001 | MUST | 错误与日志/日志规约 | 业务日志用 log/slog 记录，不用 fmt 与 log 直接打印。 |
| LOG-002 | SHOULD | 错误与日志/日志规约 | 日志字段写成键值对，不拼进消息文本。 |
| LOG-003 | SHOULD | 错误与日志/日志规约 | 请求路径上的日志用带 Context 的方法。 |
| MODERN-001 | MAY | 编程规约/数据结构规约 | 切片与映射的常见操作优先使用 slices、maps 标准库，替代手写循环。 |
| NAMING-001 | MUST | 编程规约/命名规约 | 包名使用小写单词连写，不使用下划线、驼峰或复数。 |
| NAMING-002 | MUST | 编程规约/命名规约 | 标识符不使用下划线分隔单词。 |
| NAMING-003 | MUST | 编程规约/命名规约 | 首字母缩略词在标识符中大小写统一。 |
| NAMING-004 | SHOULD | 编程规约/命名规约 | 名称长度与作用域范围成正比。 |
| NAMING-005 | SHOULD | 编程规约/命名规约 | 名称不复刻类型与外围上下文信息。 |
| NAMING-006 | SHOULD | 编程规约/命名规约 | 接收者名用类型名的短缩写，同类型内保持一致。 |
| PERF-001 | SHOULD | 性能规约 | 优化前先用基准测试或 profile 定位热点，不凭直觉改动。 |
| PERF-002 | SHOULD | 性能规约 | 原始类型与字符串互转使用 strconv，不用 fmt。 |
| PERF-003 | SHOULD | 性能规约 | 固定字符串转 []byte 在循环外转换一次后复用。 |
| PERF-004 | MUST | 性能规约 | 元素数量已知或可估算时，用 make 指定切片容量。 |
| PERF-005 | MUST | 性能规约 | 元素数量已知或可估算时，给 make(map) 传容量提示。 |
| PERF-006 | SHOULD | 性能规约 | 循环或多次拼接字符串用 strings.Builder，容量可预估时先 Grow。 |
| PERF-007 | SHOULD | 性能规约 | 长期保留大切片的一小段时，用 slices.Clone 切断底层数组引用。 |
| PERF-008 | SHOULD | 性能规约 | 函数只用解引用读取参数时按值传参，不为省几次字节传指针。 |
| PKG-001 | MUST | 依赖与工程结构/包组织规约 | 包名取导入路径的最后一段，与所在目录同名。 |
| PKG-002 | SHOULD | 依赖与工程结构/包组织规约 | 不新建 util、common、helper 一类无信息量的包。 |
| PKG-003 | MUST | 依赖与工程结构/包组织规约 | import 分组排列，标准库一组、其余一组，组间空行。 |
| PKG-004 | MUST | 依赖与工程结构/包组织规约 | 不使用 dot import。 |
| PKG-005 | SHOULD | 依赖与工程结构/包组织规约 | 导出名不再重复包名，构造函数优先用 New。 |
| PKG-006 | SHOULD | 依赖与工程结构/包组织规约 | blank import 只出现在 main 包或测试文件里。 |
| PKG-007 | SHOULD | 依赖与工程结构/包组织规约 | 只在包名冲突或与路径末段不符时才起 import 别名。 |
| SEC-001 | MUST | 安全规约 | SQL 语句用占位符传参，不把参数值拼进语句。 |
| SEC-002 | MUST | 安全规约 | 外部命令的可执行文件与参数分别传给 exec，不拼命令串。 |
| SEC-003 | MUST | 安全规约 | 用户可控的文件路径限定在允许目录内。 |
| SEC-004 | MUST | 安全规约 | 安全用途的随机数取 crypto/rand，不用 math/rand。 |
| SEC-005 | SHOULD | 安全规约 | 不设 InsecureSkipVerify 跳过 TLS 证书校验。 |
| SEC-006 | MUST | 安全规约 | 安全用途不用 MD5、SHA-1、DES、RC4 等已攻破算法。 |
| SEC-007 | MUST | 安全规约 | 不把密钥、口令、令牌硬编码在源码里。 |
| SEC-008 | SHOULD | 安全规约 | 创建文件与目录只给必需权限，不用全局可写的值。 |
| SEC-009 | SHOULD | 安全规约 | 输出 HTML 用 html/template，不用 text/template。 |
| SEC-010 | SHOULD | 安全规约 | HTTP 服务显式设置读头、读体与写超时。 |
| SEC-011 | SHOULD | 安全规约 | 比较口令、令牌、MAC 用常量时间函数，不用提前返回的比较。 |
| STRUCT-001 | MUST | 接口与设计规约/结构体安全形状 | 含 sync.Mutex 等同步原语的结构体禁止值拷贝，一律用指针传递或接收。 |
| STRUCT-002 | MAY | 接口与设计规约/结构体安全形状 | 不应支持比较的结构体，用不可比较字段让编译期拦截 ==。 |
| STYLE-001 | SHOULD | 编程规约/格式化与代码风格规约 | 用 any 代替 interface{}。 |
| TEST-001 | SHOULD | 单元测试规约/测试命名与结构 | 纯函数测试用原生 testing，组件测试用 testify suite。 |
| TEST-002 | SHOULD | 单元测试规约/测试命名与结构 | 单元测试与集成测试分开存放。 |
| TEST-003 | SHOULD | 单元测试规约/断言与测试数据 | 一次比对完整结果，不逐字段拆散或裁剪成子集。 |
| TEST-004 | SHOULD | 单元测试规约/断言与测试数据 | 庞大的测试输入与期望数据放 testdata 目录。 |
| TEST-005 | MUST | 单元测试规约/测试命名与结构 | 测试文件用包外测试，包名以 _test 结尾。 |
| TEST-006 | SHOULD | 单元测试规约/覆盖率的使用 | 测业务逻辑与边界，不追求覆盖率数字。 |
| TEST-007 | SHOULD | 单元测试规约/断言与测试数据 | 断言错误只判是否为非 nil，不比对错误消息文本。 |
| TEST-008 | MAY | 单元测试规约/测试命名与结构 | 测试用例名优先写成 should 接结果、when 接条件的句式。 |
| TEST-010 | SHOULD | 单元测试规约/表驱动与用例设计 | 检查逻辑相同的用例才进同一张表。 |
| TEST-011 | SHOULD | 单元测试规约/表驱动与用例设计 | 表格只放随用例变化的字段。 |
| TEST-012 | MUST | 单元测试规约/表驱动与用例设计 | 期望值独立写出，不用被测函数生成。 |
| TEST-013 | SHOULD | 单元测试规约/表驱动与用例设计 | 边界与特殊输入单独建用例。 |
| TEST-014 | SHOULD | 单元测试规约/表驱动与用例设计 | 修缺陷的改动附上复现用例。 |
| TEST-015 | MUST | 单元测试规约/Mock 与测试替身 | Mock 生成用 go.uber.org/mock。 |
| TEST-016 | MUST | 单元测试规约/Mock 与测试替身 | 传入 *testing.T 后不调用 ctrl.Finish()。 |
| TEST-017 | MUST | 单元测试规约/Mock 与测试替身 | 不写冗余的 .Times(1)。 |
| TEST-018 | SHOULD | 单元测试规约/Mock 与测试替身 | 确有调用顺序要求时才用 gomock.InOrder()。 |
| TEST-019 | SHOULD | 单元测试规约/Mock 与测试替身 | 期望里直接写参数值，避免包裹 Eq 与裸 Any()。 |
| TOOL-001 | MUST | 静态检查规约/配置与运行 | 统一用 golangci-lint 做静态代码检查。 |
| TOOL-002 | SHOULD | 静态检查规约/推荐规则 | golangci-lint 中开启 modernize，跟进新版 Go 惯用法。 |
| TOOL-003 | SHOULD | 静态检查规约/推荐规则 | 使用 testify 时，同步开启 testifylint。 |
| TOOL-004 | SHOULD | 工具规约/命令行 | 命令行程序用 cobra 组织子命令。 |
| TOOL-005 | MUST | 工具规约/测试工具 | 测试断言统一使用 testify。 |
| TOOL-006 | MUST | 静态检查规约/必开规则 | 至少开启 golangci-lint 的 standard 集。 |
| TOOL-007 | SHOULD | 静态检查规约/推荐规则 | standard 集之外，推荐开启常用的补充分析器。 |
| TOOL-008 | SHOULD | 工具规约/命令行 | 升级 Go 工具链后用 go fix 应用现代化改写。 |
| TOOL-009 | MUST | 静态检查规约/配置与运行 | 抑制告警要指名检查器并写明理由。 |
| TOOL-010 | MUST | 工具规约/代码生成 | 生成物不手工编辑，只由生成器产出。 |
| TOOL-011 | MUST | 工具规约/代码生成 | 生成文件命名带 gen 前缀或后缀，优先用后缀。 |

## 2 依据来源

- ABSTRACT-001：https://go.dev/wiki/CodeReviewComments#interfaces
- CODE-001：https://cloud.google.com/apis/design/errors
- CODE-002：https://cloud.google.com/apis/design/errors
- CODE-003：https://cloud.google.com/apis/design/errors
- COMMENT-001：https://google.github.io/styleguide/go/guide#clarity-rationale
- COMMENT-002：https://go.dev/wiki/CodeReviewComments#comment-sentences https://go.dev/doc/effective_go#commentary
- COMPOSE-001：https://go.dev/doc/faq#inheritance https://go.dev/doc/effective_go#embedding
- COMPOSE-002：https://github.com/uber-go/guide/blob/master/style.md#embedding-in-structs https://github.com/uber-go/guide/blob/master/style.md#avoid-embedding-types-in-public-structs
- CONC-001：https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes https://pkg.go.dev/golang.org/x/sync/errgroup
- CONC-002：https://go.dev/blog/pipelines
- CONC-003：https://go.dev/wiki/CodeReviewComments#synchronous-functions
- CONC-004：https://go.dev/wiki/MutexOrChannel https://pkg.go.dev/sync
- CONST-001：https://go.dev/wiki/CodeReviewComments#mixed-caps https://google.github.io/styleguide/go/decisions#constant-names
- CONST-002：https://google.github.io/styleguide/go/decisions#constant-names
- CONST-003：https://go.dev/doc/effective_go#constants https://go.dev/ref/spec#Constant_declarations
- CTRL-001：https://go.dev/wiki/CodeReviewComments#indent-error-flow
- CTRL-002：https://go.dev/doc/go1.22 https://go.dev/blog/loopvar-preview
- CTX-001：https://go.dev/wiki/CodeReviewComments#contexts https://pkg.go.dev/context
- CTX-002：https://pkg.go.dev/context#pkg-overview https://staticcheck.dev/docs/checks/#SA1012
- CTX-003：https://pkg.go.dev/context#Context.Value https://staticcheck.dev/docs/checks/#SA1029
- CTX-004：https://pkg.go.dev/context#pkg-overview
- DATA-001：https://go.dev/wiki/CodeReviewComments#declaring-empty-slices
- DATA-002：https://pkg.go.dev/cmd/vet https://github.com/uber-go/guide/blob/master/style.md#use-field-names-to-initialize-structs
- DEP-001：https://go.dev/doc/modules/managing-dependencies#synchronizing https://go.dev/doc/modules/gomod-ref
- DEP-002：https://go.dev/doc/modules/managing-dependencies#enable_tracking
- DEP-003：https://go.dev/doc/modules/managing-dependencies#local_directory https://go.dev/doc/modules/gomod-ref#replace
- DEP-004：https://go.dev/doc/modules/gomod-ref#tool https://go.dev/doc/modules/managing-dependencies#tools
- DEP-005：https://go-proverbs.github.io/
- DI-001：https://github.com/google/wire
- DI-002：https://github.com/uber-go/guide/blob/master/style.md#avoid-mutable-globals
- ERR-001：https://go.dev/blog/go1.13-errors https://google.github.io/styleguide/go/decisions
- ERR-003：https://go.dev/blog/go1.13-errors https://pkg.go.dev/errors#Is
- ERR-004：https://go.dev/wiki/CodeReviewComments#dont-panic https://go.dev/doc/effective_go#errors
- ERR-005：https://go.dev/wiki/CodeReviewComments#in-band-errors
- ERR-006：https://dave.cheney.net/practical-go/presentations/qcon-china.html#_only_handle_an_error_once https://go.dev/wiki/CodeReviewComments#handle-errors
- FUNC-001：https://github.com/uber-go/guide/blob/master/style.md#function-grouping-and-ordering https://github.com/manuelarte/funcorder
- FUNC-002：https://go.dev/wiki/CodeReviewComments#named-result-parameters
- IFACE-001：https://go.dev/wiki/CodeReviewComments#interfaces
- IFACE-002：https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
- IFACE-003：https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces
- LAYER-001：https://alistair.cockburn.us/hexagonal-architecture/
- LAYOUT-001：https://github.com/golang-standards/project-layout#directories-you-shouldnt-have https://go.dev/doc/modules/layout
- LAYOUT-002：https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages https://go.dev/doc/go1.4#internalpackages
- LAYOUT-003：https://go.dev/doc/modules/layout#multiple-commands https://github.com/golang-standards/project-layout#cmd
- LAYOUT-004：https://github.com/golang-standards/project-layout#overview
- LOG-001：https://go.dev/blog/slog https://pkg.go.dev/log/slog
- LOG-002：https://go.dev/blog/slog
- LOG-003：https://pkg.go.dev/log/slog#Logger.InfoContext https://go.dev/blog/slog
- MODERN-001：https://pkg.go.dev/slices https://go.dev/doc/go1.21
- NAMING-001：https://go.dev/wiki/CodeReviewComments#package-names https://google.github.io/styleguide/go/decisions
- NAMING-002：https://google.github.io/styleguide/go/decisions#underscores https://go.dev/wiki/CodeReviewComments#mixed-caps
- NAMING-003：https://go.dev/wiki/CodeReviewComments#initialisms https://google.github.io/styleguide/go/decisions#initialisms
- NAMING-004：https://google.github.io/styleguide/go/decisions#variable-names https://go.dev/wiki/CodeReviewComments#variable-names
- NAMING-005：https://google.github.io/styleguide/go/decisions#variable-name-vs-type https://google.github.io/styleguide/go/decisions#external-context-vs-local-names
- NAMING-006：https://google.github.io/styleguide/go/decisions#receiver-names https://go.dev/wiki/CodeReviewComments#receiver-names
- PERF-001：https://go.dev/wiki/Performance
- PERF-002：https://github.com/uber-go/guide/blob/master/style.md#prefer-strconv-over-fmt
- PERF-003：https://github.com/uber-go/guide/blob/master/style.md#avoid-repeated-string-to-byte-conversions
- PERF-004：https://go.dev/blog/allocation-optimizations https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity
- PERF-005：https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity
- PERF-006：https://pkg.go.dev/strings#Builder
- PERF-007：https://go.dev/blog/slices-intro#a-possible-gotcha
- PERF-008：https://go.dev/wiki/CodeReviewComments#pass-values
- PKG-001：https://go.dev/blog/package-names https://go.dev/wiki/CodeReviewComments#package-names
- PKG-002：https://google.github.io/styleguide/go/decisions#package-names https://go.dev/blog/package-names
- PKG-003：https://go.dev/wiki/CodeReviewComments#imports
- PKG-004：https://go.dev/wiki/CodeReviewComments#import-dot
- PKG-005：https://go.dev/blog/package-names https://google.github.io/styleguide/go/decisions#package-names
- PKG-006：https://go.dev/wiki/CodeReviewComments#import-blank
- PKG-007：https://go.dev/wiki/CodeReviewComments#imports
- SEC-001：https://go.dev/doc/database/sql-injection https://go.dev/doc/database/prepared-statements
- SEC-002：https://pkg.go.dev/os/exec
- SEC-003：https://go.dev/doc/go1.24#directory-limited-filesystem-access https://pkg.go.dev/os#Root
- SEC-004：https://pkg.go.dev/math/rand https://pkg.go.dev/crypto/rand
- SEC-005：https://pkg.go.dev/crypto/tls#Config
- SEC-006：https://pkg.go.dev/crypto/md5 https://pkg.go.dev/crypto/sha1 https://pkg.go.dev/crypto/des https://pkg.go.dev/crypto/rc4
- SEC-007：https://github.com/securego/gosec/blob/master/RULES.md
- SEC-008：https://github.com/securego/gosec/blob/master/RULES.md
- SEC-009：https://pkg.go.dev/html/template
- SEC-010：https://pkg.go.dev/net/http#Server.ReadHeaderTimeout https://github.com/securego/gosec/blob/master/RULES.md
- SEC-011：https://pkg.go.dev/crypto/subtle#ConstantTimeCompare https://pkg.go.dev/crypto/hmac#Equal
- STRUCT-001：https://pkg.go.dev/cmd/vet https://google.github.io/styleguide/go/decisions
- STRUCT-002：https://github.com/protocolbuffers/protobuf-go/blob/master/internal/pragma/pragma.go https://go.dev/ref/spec#Comparison_operators
- STYLE-001：https://go.dev/doc/go1.18 https://google.github.io/styleguide/go/decisions
- TEST-001：https://github.com/stretchr/testify#suite-package https://pkg.go.dev/github.com/stretchr/testify/suite
- TEST-002：https://pkg.go.dev/cmd/go#hdr-Test_packages
- TEST-003：https://go.dev/wiki/TestComments#compare-full-structures
- TEST-004：https://pkg.go.dev/embed https://pkg.go.dev/cmd/go#hdr-Test_packages
- TEST-005：https://pkg.go.dev/cmd/go#hdr-Test_packages
- TEST-006：https://research.swtch.com/testing
- TEST-007：https://go.dev/wiki/TestComments#test-error-semantics https://go.dev/blog/go1.13-errors
- TEST-008：https://go.dev/wiki/TestComments#choose-human-readable-subtest-names
- TEST-010：https://go.dev/wiki/TestComments#table-driven-tests-vs-multiple-test-functions https://research.swtch.com/testing
- TEST-011：https://research.swtch.com/testing https://go.dev/wiki/TableDrivenTests
- TEST-012：https://go.dev/wiki/TestComments#compare-full-structures https://research.swtch.com/testing
- TEST-013：https://research.swtch.com/testing
- TEST-014：https://research.swtch.com/testing
- TEST-015：https://github.com/uber-go/mock
- TEST-016：https://pkg.go.dev/go.uber.org/mock/gomock#Controller.Finish
- TEST-017：https://pkg.go.dev/go.uber.org/mock/gomock#Call.Times https://github.com/uber-go/mock
- TEST-018：https://pkg.go.dev/go.uber.org/mock/gomock#InOrder
- TEST-019：https://pkg.go.dev/go.uber.org/mock/gomock#Matcher https://github.com/uber-go/mock/blob/main/gomock/call.go
- TOOL-001：https://golangci-lint.run/
- TOOL-002：https://golangci-lint.run/docs/linters/configuration/#modernize
- TOOL-003：https://github.com/Antonboom/testifylint
- TOOL-004：https://github.com/spf13/cobra
- TOOL-005：https://github.com/stretchr/testify#assert-package
- TOOL-006：https://golangci-lint.run/docs/welcome/quick-start/
- TOOL-007：https://golangci-lint.run/docs/linters/
- TOOL-008：https://pkg.go.dev/cmd/fix https://go.dev/blog/gofix https://go.dev/doc/go1.26
- TOOL-009：https://golangci-lint.run/docs/linters/false-positives/#nolint-directive https://golangci-lint.run/docs/linters/configuration/#nolintlint
- TOOL-010：https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source
- TOOL-011：https://protobuf.dev/reference/go/go-generated/ https://pkg.go.dev/cmd/go#hdr-Build_constraints

## 3 检测映射

| 编号 | 检测方式 |
| --- | --- |
| ABSTRACT-001 | 人工核对（核对接口或工厂是否已有第二个实现、替换或测试需求） |
| CODE-001 | 人工核对（核对对外错误码是否与标准分类并存了自造编号） |
| CODE-002 | 人工核对（核对是否有按错误消息文本判断的分支） |
| CODE-003 | 人工核对（核对对外返回的错误是否直接来自下游） |
| COMMENT-001 | 人工核对（核对注释是否在复述代码本身能读出的内容） |
| COMMENT-002 | 人工核对（核对文档注释是否以名称开头并写成完整句子） |
| COMPOSE-001 | 人工核对（核对内嵌类型是否承担基类角色（Base / Abstract 命名、多层嵌入）） |
| COMPOSE-002 | 人工核对（核对嵌入是否改变外层零值、拷贝语义或公开无关方法（同步原语尤其）） |
| CONC-001 | 人工核对（核对每个 go 是否有对应的等待方与退出路径） |
| CONC-002 | 人工核对（核对 close 是否都在发送方一侧） |
| CONC-003 | 人工核对（核对对外函数是否返回 channel 让调用方等待） |
| CONC-004 | 人工核对（核对受锁保护的字段是否导出、是否只能经方法访问） |
| CONST-001 | golangci-lint stylecheck(ST1003) |
| CONST-002 | 人工核对（核对常量名是否只是在复述取值） |
| CONST-003 | 人工核对（核对同源递增的枚举是否手写了数字） |
| CTRL-001 | 人工核对（核对正常路径是否被 else 缩进） |
| CTRL-002 | golangci-lint copyloopvar |
| CTX-001 | golangci-lint containedctx（首参约定需人工核对） |
| CTX-002 | golangci-lint staticcheck(SA1012)（SA1012 检出向函数传 nil context 的调用） |
| CTX-003 | golangci-lint staticcheck(SA1029)（SA1029 检出用内建类型作 context key） |
| CTX-004 | 人工核对（核对 context.Value 里放的是否为请求作用域数据） |
| DATA-001 | 人工核对（需要输出 JSON 数组 [] 时属例外） |
| DATA-002 | golangci-lint govet(composites)（govet 只覆盖跨包位置字面量，需启用 composites 检查） |
| DEP-001 | 人工核对（提交前运行 go mod tidy，确认 go.mod 与 go.sum 无未提交改动） |
| DEP-002 | 人工核对（核对仓库含 go.sum 且 .gitignore 未忽略它） |
| DEP-003 | 人工核对（核对 go.mod 里的 replace 未长期指向本地目录，且写明了移除条件） |
| DEP-004 | 人工核对（核对仓库用到的命令行工具经 go.mod 的 tool 指令声明） |
| DEP-005 | 人工核对（核对新增依赖是否为标准库已能覆盖的小功能） |
| DI-001 | 人工核对（核对仓库有 wire_gen.go 且端口绑定集中在 ProviderSet） |
| DI-002 | 人工核对（核对包级变量是否存在赋值，依赖是否从构造函数传入） |
| ERR-001 | golangci-lint errorlint |
| ERR-003 | golangci-lint errorlint（覆盖 == 比较与类型断言两类写法） |
| ERR-004 | 人工核对（核对 panic 是否用在可预期的失败上） |
| ERR-005 | 人工核对（核对是否用 -1、空值或零值表达失败） |
| ERR-006 | 人工核对（逐个 if err != nil 分支核对：同一分支内既写日志又 return err 即违反） |
| FUNC-001 | golangci-lint funcorder |
| FUNC-002 | 人工核对（核对返回值含义不自明时是否命名） |
| IFACE-001 | 人工核对（核对接口所在包与实现数量是否对应：单实现看是否与实现同包，多实现看是否在使用方包） |
| IFACE-002 | 人工核对（核对按 API 契约实现接口的导出类型处是否有编译期断言） |
| IFACE-003 | 人工核对（核对函数签名与结构体字段是否出现 *接口类型） |
| LAYER-001 | 人工核对（核对包 import 方向；可用 golangci-lint 的 depguard 固化为规则） |
| LAYOUT-001 | 人工核对（核对 go.mod 所在目录即仓库根，源码未放在 src/ 之下） |
| LAYOUT-002 | 人工核对（核对仅本模块使用的包是否位于 internal/ 下） |
| LAYOUT-003 | 人工核对（核对 main 包集中在 cmd/ 下，目录名与可执行文件同名） |
| LAYOUT-004 | 人工核对（核对仓库目录是否都有实际内容与明确职责） |
| LOG-001 | grep 正则 \b(fmt\.Print|log\.Print)（命中处人工核对是否为业务日志） |
| LOG-002 | 人工核对（核对日志取值是否拼进了消息字符串） |
| LOG-003 | 人工核对（核对请求路径上的日志是否传入 ctx） |
| MODERN-001 | golangci-lint modernize（go fix 也会提示等价的现代写法） |
| NAMING-001 | golangci-lint stylecheck(ST1003)（检查包名标识符） |
| NAMING-002 | golangci-lint stylecheck(ST1003)（检查标识符命名风格） |
| NAMING-003 | golangci-lint stylecheck(ST1003) |
| NAMING-004 | 人工核对（人工判断名称长度是否与作用域匹配） |
| NAMING-005 | 人工核对（核对名称是否重复了类型或外围上下文） |
| NAMING-006 | 人工核对（核对同一类型的接收者名是否一致、是否用了 this / self） |
| PERF-001 | 人工核对（核对性能改动是否附基准测试结果或 profile 结论） |
| PERF-002 | 人工核对（核对单一原始类型与字符串互转处是否用了 fmt） |
| PERF-003 | 人工核对（核对循环内是否有对固定字符串的 string 与 []byte 互转） |
| PERF-004 | 人工核对（核对 append 目标切片的 make 是否带了容量参数） |
| PERF-005 | 人工核对（核对逐条写入的 map 是否带容量提示） |
| PERF-006 | 人工核对（核对循环内是否有字符串 += 拼接） |
| PERF-007 | 人工核对（核对从大缓冲区截取并长期保存的子切片是否做了复制） |
| PERF-008 | 人工核对（核对只做解引用读取的参数是否用了指针） |
| PKG-001 | 人工核对（逐个目录核对 package 子句与目录名一致） |
| PKG-002 | grep 正则 ^package\s+(util|utility|common|misc|helper|helpers|model)\b（命中处按其真实职责拆成具名包） |
| PKG-003 | golangci-lint gci（formatters 中启用 gci，分组顺序为 standard、default、本仓库 module 前缀） |
| PKG-004 | grep 正则 ^\s*(import\s+)?\.\s+"（命中处核对是否为循环依赖场景下的 _test 文件） |
| PKG-005 | 人工核对（核对导出名是否重复其所在包名） |
| PKG-006 | grep 正则 ^\s*(import\s+)?_\s+"（命中处核对所在文件是否为 main 包或 _test 文件） |
| PKG-007 | 人工核对（核对每个别名是否用于解决同名冲突或包名与路径不符） |
| SEC-001 | golangci-lint gosec(G201, G202)（命中处确认参数值是否来自外部输入） |
| SEC-002 | golangci-lint gosec(G204)（命中处确认是否把外部输入拼进了命令串） |
| SEC-003 | golangci-lint gosec(G304)（命中处确认路径是否含外部输入、是否限定在根目录内） |
| SEC-004 | golangci-lint gosec(G404)（命中 math/rand 调用，确认是否用于安全场景） |
| SEC-005 | golangci-lint gosec(G402)（命中处确认是否非测试代码仍跳过校验） |
| SEC-006 | golangci-lint gosec(G401, G405, G501, G502, G503, G505) |
| SEC-007 | golangci-lint gosec(G101)（命中处确认是否真实凭证，误报用 #nosec G101 标注理由） |
| SEC-008 | golangci-lint gosec(G301, G302, G306)（阈值可在 .golangci.yml 的 gosec 配置中按需收紧） |
| SEC-009 | golangci-lint gosec(G203)（G203 命中把未转义数据写进 HTML 模板的用法） |
| SEC-010 | golangci-lint gosec(G112, G114)（G112 检出未设 ReadHeaderTimeout，G114 检出无超时支持的 serve 函数） |
| SEC-011 | 人工核对（核对秘密比较处是否用了常量时间函数） |
| STRUCT-001 | golangci-lint govet(copylocks) |
| STRUCT-002 | 人工核对（核对有意禁止比较的导出结构体是否嵌入了不可比较字段） |
| STYLE-001 | golangci-lint modernize(any)（同一分析器 go fix 可直接改写） |
| TEST-001 | 人工核对（核对纯函数用例是否套了 suite、组件用例是否用 suite 集中组装被测对象） |
| TEST-002 | 人工核对（核对读取真实文件、网络或数据库的测试是否以 _it_test.go 结尾） |
| TEST-003 | 人工核对（核对断言是否拆成逐字段，或把结果裁剪成子集后再比） |
| TEST-004 | 人工核对（核对用例读取 testdata 是否靠运行目录拼路径） |
| TEST-005 | golangci-lint testpackage |
| TEST-006 | 人工核对（核对用例是否有断言、是否对准业务语义；覆盖率数字不作为判定依据） |
| TEST-007 | 人工核对（核对用例是否用错误消息文本判断错误种类；哨兵与错误类型是否用 errors.Is / errors.As） |
| TEST-008 | 人工核对（核对 t.Run 与表格 name 是否写清期望结果与触发条件） |
| TEST-010 | 人工核对（核对同一循环体内是否用条件分支切换检查方式） |
| TEST-011 | 人工核对（核对表格字段是否被所有用例写成同一份值） |
| TEST-012 | 人工核对（核对期望值是否由被测包内的函数或方法产出） |
| TEST-013 | 人工核对（核对是否覆盖空集合、单元素、首尾与越界等边界输入） |
| TEST-014 | 人工核对（核对修缺陷的提交是否带上触发条件对应的用例） |
| TEST-015 | grep 正则 github\.com/golang/mock（命中处改用 go.uber.org/mock；go.mod 不在扫描范围内，需人工核对） |
| TEST-016 | grep 正则 \.Finish\(\)（命中处确认是否传入了 *testing.T） |
| TEST-017 | grep 正则 \.Times\(1\)（命中处删除即可，默认即一次） |
| TEST-018 | grep 正则 gomock\.InOrder（命中处确认调用顺序是否影响结果） |
| TEST-019 | grep 正则 gomock\.Any\(\)（除 context.Context 外，命中处改用 Eq 或 AssignableToTypeOf） |
| TOOL-001 | 人工核对（核对仓库根有 .golangci.yml 且 lint 入口调用 golangci-lint） |
| TOOL-002 | 人工核对（核对 .golangci.yml 的 linters.enable 含 modernize） |
| TOOL-003 | 人工核对（核对 .golangci.yml 的 linters.enable 含 testifylint） |
| TOOL-004 | 人工核对（核对是否以 cobra 组织子命令） |
| TOOL-005 | grep 正则 \bt\.(Errorf|Fatalf)\(（命中处应改用 testify 断言） |
| TOOL-006 | 人工核对（核对 .golangci.yml 未关闭 standard 集） |
| TOOL-007 | 人工核对（核对 enable 列表是否覆盖清单，未启用项是否有理由） |
| TOOL-008 | 人工核对（核对升级 Go 版本后是否跑过 go fix） |
| TOOL-009 | golangci-lint nolintlint（配置为 require-specific + require-explanation + allow-unused，拦下不指名检查器、缺理由与已失效的抑制） |
| TOOL-010 | 人工核对（核对改动是否落在生成器或输入源上） |
| TOOL-011 | 人工核对（核对生成物文件名是否带 gen 前缀或后缀） |
