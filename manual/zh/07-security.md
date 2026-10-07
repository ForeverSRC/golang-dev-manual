# 七、安全规约

### 【MUST】SEC-001 SQL 语句用占位符传参，不把参数值拼进语句。

- 归属：安全规约
- 起始版本：Go 1.8

查询与写入的参数值经 `database/sql` 的占位符传入，不把值拼进 SQL 文本。表名、列名、排序字段这类占位符覆盖不到的标识符，先用白名单校验再拼接。

**为什么**

> You can avoid an SQL injection risk by providing SQL parameter values as `sql` package function arguments. Many functions in the `sql` package provide parameters for the SQL statement and for values to be used in that statement's parameters (others provide a parameter for a prepared statement and parameters).
>
> —— https://go.dev/doc/database/sql-injection

参数交给 `sql` 包后，驱动把语句与参数分开发送，参数值不会被当作 SQL 语法解析。改用 `fmt.Sprintf` 先把值拼进语句再交给 `Query`，整条 SQL 已成型，调用方传入的片段直接参与语法，`id` 传成 `1 OR 1=1` 就能绕过条件，官方文档把这处写法直接标成 SECURITY RISK。

**正例**

```go
rows, err := db.QueryContext(ctx, "SELECT id, name FROM user WHERE id = ?", id)
```

**反例**

```go
rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT id, name FROM user WHERE id = %s", id))
```

**依据**

- https://go.dev/doc/database/sql-injection
- https://go.dev/doc/database/prepared-statements

**检测**：golangci-lint gosec(G201, G202)（命中处确认参数值是否来自外部输入）

### 【MUST】SEC-002 外部命令的可执行文件与参数分别传给 exec，不拼命令串。

- 归属：安全规约
- 起始版本：Go 1.7

用 `exec.Command`、`exec.CommandContext` 时把可执行文件与每个参数作为独立参数传入，不借助 `sh -c`、`cmd /c` 拼接完整命令行。需要通配符展开或管道时改用 `filepath.Glob` 或在程序里实现；不得不经 shell 时，先对拼接的外部输入做转义。

**为什么**

> Unlike the "system" library call from C and other languages, the os/exec package intentionally does not invoke the system shell and does not expand any glob patterns or handle other expansions, pipelines, or redirections typically done by shells.
>
> —— https://pkg.go.dev/os/exec

`os/exec` 按参数数组启动进程，不走 shell，参数里的 `;`、`|`、`&&` 只当普通字符。改成 `sh -c` 拼接后，外部输入重新进入 shell 语法，一段 `; rm -rf` 就能改变实际执行的命令；官方文档在通配符展开处也提醒直接调 shell 要当心转义外部输入。

**正例**

```go
cmd := exec.CommandContext(ctx, "git", "log", "--oneline", "-n", n)
```

**反例**

```go
cmd := exec.CommandContext(ctx, "sh", "-c", "git log --oneline -n "+n)
```

**依据**

- https://pkg.go.dev/os/exec

**检测**：golangci-lint gosec(G204)（命中处确认是否把外部输入拼进了命令串）

### 【MUST】SEC-003 用户可控的文件路径限定在允许目录内。

- 归属：安全规约
- 起始版本：Go 1.24

外部输入参与构造文件路径时，先用 `os.OpenRoot` 在允许的根目录上打开 `*os.Root`，再经它的 `Open`、`Create`、`ReadFile` 等方法访问。用不上 `os.Root` 时，至少 `filepath.Clean` 后校验结果仍落在根目录前缀内，并单独处理符号链接。

**为什么**

> The os.OpenRoot function opens a directory and returns an os.Root. Methods on os.Root operate within the directory and do not permit paths that refer to locations outside the directory, including ones that follow symbolic links out of the directory.
>
> —— https://go.dev/doc/go1.24#directory-limited-filesystem-access

路径里的 `../` 与指向外部的符号链接能把 `os.Open` 引出预期目录，读到或覆盖根目录外的文件。`os.Root` 以根目录为边界，逐次操作校验最终解析结果，越界路径直接报错；手写校验要在 `Clean` 之后比对前缀，还要单独挡住符号链接，容易漏项。

**正例**

```go
root, err := os.OpenRoot(dataDir)
if err != nil {
	return err
}
defer root.Close()

f, err := root.Open(filepath.Join("reports", name))
```

**反例**

```go
f, err := os.Open(filepath.Join(dataDir, name))
```

**依据**

- https://go.dev/doc/go1.24#directory-limited-filesystem-access
- https://pkg.go.dev/os#Root

**检测**：golangci-lint gosec(G304)（命中处确认路径是否含外部输入、是否限定在根目录内）

### 【MUST】SEC-004 安全用途的随机数取 crypto/rand，不用 math/rand。

- 归属：安全规约
- 起始版本：Go 1.24

生成令牌、会话 ID、密钥、盐值、验证码等与安全相关的随机值时用 `crypto/rand` 的 `Read`、`Text` 等。`math/rand` 只用于模拟、抽样、洗牌一类不涉及安全的场景。

**为什么**

> Package rand implements pseudo-random number generators suitable for tasks such as simulation, but it should not be used for security-sensitive work.
>
> —— https://pkg.go.dev/math/rand

`math/rand` 的输出由种子决定、可被推断，官方文档明说它不适合安全用途，并指出序列容易被猜出。用它生成令牌，攻击者拿到少量输出即可推算出后续值；`crypto/rand` 取操作系统提供的安全随机源，输出不可预测。

**正例**

```go
import "crypto/rand"

token := make([]byte, 32)
rand.Read(token)
```

**反例**

```go
import "math/rand"

token := make([]byte, 32)
rand.Read(token)
```

**依据**

- https://pkg.go.dev/math/rand
- https://pkg.go.dev/crypto/rand

**检测**：golangci-lint gosec(G404)（命中 math/rand 调用，确认是否用于安全场景）

### 【MUST】SEC-006 安全用途不用 MD5、SHA-1、DES、RC4 等已攻破算法。

- 归属：安全规约
- 起始版本：Go 1.0

口令哈希、签名、消息认证、密钥派生等安全场景不引入 `crypto/md5`、`crypto/sha1`、`crypto/des`、`crypto/rc4`、`golang.org/x/crypto/md4`、`golang.org/x/crypto/ripemd160`。这些算法在非安全场景仍可用，例如文件校验和与内容寻址。

**为什么**

> MD5 is cryptographically broken and should not be used for secure applications.
>
> —— https://pkg.go.dev/crypto/md5

这几个包的文档各自写明算法已被攻破、不应再用于安全场合。MD5 与 SHA-1 已有可行的碰撞构造，两份不同输入能哈希到同一值，用它们做完整性校验或签名，内容被换掉仍能通过；DES 与 RC4 的密钥空间和密钥流弱点可被实际破解。

**正例**

```go
sum := sha256.Sum256(data)
```

**反例**

```go
sum := md5.Sum(data)
```

**依据**

- https://pkg.go.dev/crypto/md5
- https://pkg.go.dev/crypto/sha1
- https://pkg.go.dev/crypto/des
- https://pkg.go.dev/crypto/rc4

**检测**：golangci-lint gosec(G401, G405, G501, G502, G503, G505)

### 【MUST】SEC-007 不把密钥、口令、令牌硬编码在源码里。

- 归属：安全规约
- 起始版本：Go 1.0

密钥、口令、访问令牌、连接串等敏感值经环境变量、配置文件或密钥管理服务注入，不写成源码里的字面量。测试凭证放测试夹具并明显区分，不进生产代码与版本库。

**为什么**

> G101 — Look for hardcoded credentials
>
> —— https://github.com/securego/gosec/blob/master/RULES.md

写进源码的凭证随提交进入版本历史，读过仓库的人都能拿到；换密钥要改代码重新发布，实际往往长期不换，一处泄露的影响面持续扩大。gosec 的 G101 按标识符名与字面量熵值识别这类赋值，命中处应改为从外部注入。

**正例**

```go
apiKey := os.Getenv("API_KEY")
```

**反例**

```go
apiKey := "sk-live-9f3c1d2e4b5a6c7d"
```

**依据**

- https://github.com/securego/gosec/blob/master/RULES.md

**检测**：golangci-lint gosec(G101)（命中处确认是否真实凭证，误报用 #nosec G101 标注理由）

### 【SHOULD】SEC-005 不设 InsecureSkipVerify 跳过 TLS 证书校验。

- 归属：安全规约
- 起始版本：Go 1.0

`tls.Config` 保持默认的证书链与主机名校验，`InsecureSkipVerify` 保持 false。自签证书场景把 CA 加进 `RootCAs`，或用 `VerifyPeerCertificate`、`VerifyConnection` 自定义校验。测试环境与一次性验证脚本里可临时开启，并写明理由。

**为什么**

> If InsecureSkipVerify is true, crypto/tls accepts any certificate presented by the server and any host name in that certificate. In this mode, TLS is susceptible to machine-in-the-middle attacks unless custom verification is used.
>
> —— https://pkg.go.dev/crypto/tls#Config

证书校验是 TLS 确认对端身份的依据，关掉后客户端接受任意证书与任意主机名，中间人用自签证书即可解密或篡改流量。跳过校验去连生产服务，加密通道只剩防被动窃听，主动攻击不再受阻。

**正例**

```go
cfg := &tls.Config{}
```

**反例**

```go
cfg := &tls.Config{InsecureSkipVerify: true}
```

**依据**

- https://pkg.go.dev/crypto/tls#Config

**检测**：golangci-lint gosec(G402)（命中处确认是否非测试代码仍跳过校验）

### 【SHOULD】SEC-008 创建文件与目录只给必需权限，不用全局可写的值。

- 归属：安全规约
- 起始版本：Go 1.16

`os.OpenFile`、`os.WriteFile`、`os.Mkdir`、`os.MkdirAll` 的权限参数按最小需要取，不用 `0777`、`0666` 这类对所有用户开放写权限的值，含敏感数据的文件用 `0600`。需要精确权限时在创建后用 `Chmod` 显式设置，进程 umask 会先过滤一遍传入值。

**为什么**

> File and directory permission rules can be configured with stricter maximum permissions
>
> —— https://github.com/securego/gosec/blob/master/RULES.md

创建时给 0777，文件对所有用户可写，本机其他账户或共享环境里的进程都能替换内容；写进临时目录、或被当作配置与脚本读取时，可被用于提权或注入。gosec 的 G301、G302、G306 对准这几类创建与改权限调用，阈值可按项目收紧。

**正例**

```go
err := os.WriteFile(path, data, 0o600)
```

**反例**

```go
err := os.WriteFile(path, data, 0o777)
```

**依据**

- https://github.com/securego/gosec/blob/master/RULES.md

**检测**：golangci-lint gosec(G301, G302, G306)（阈值可在 .golangci.yml 的 gosec 配置中按需收紧）

### 【SHOULD】SEC-009 输出 HTML 用 html/template，不用 text/template。

- 归属：安全规约
- 起始版本：Go 1.0

生成 HTML、XML、JS 等需要转义的文本时用 `html/template`，它按数据所处上下文自动转义。`text/template` 只用于纯文本输出；只能用 `text/template` 产出 HTML 时，输出前自行对数据转义。

**为什么**

> Package template (html/template) implements data-driven templates for generating HTML output safe against code injection. It provides the same interface as text/template and should be used instead of text/template whenever the output is HTML.
>
> —— https://pkg.go.dev/html/template

`text/template` 把数据原样写进输出，数据里一段 `<script>` 会作为标签被浏览器执行，构成跨站脚本。`html/template` 按数据落在标签、属性、URL、JS 字符串哪个位置套用对应转义，模板代码不动，换包即挡住注入。

**正例**

```go
import "html/template"

t, err := template.New("page").Parse(pageTemplate)
```

**反例**

```go
import "text/template"

t, err := template.New("page").Parse(pageTemplate)
```

**依据**

- https://pkg.go.dev/html/template

**检测**：golangci-lint gosec(G203)（G203 命中把未转义数据写进 HTML 模板的用法）

### 【SHOULD】SEC-010 HTTP 服务显式设置读头、读体与写超时。

- 归属：安全规约
- 起始版本：Go 1.8

`http.Server` 至少设置 `ReadHeaderTimeout`，按业务再配 `ReadTimeout`、`WriteTimeout`、`IdleTimeout`。用 `http.ListenAndServe` 一类不带超时的便捷函数时改走自定义 `http.Server`。

**为什么**

> If ReadHeaderTimeout is zero, the value of ReadTimeout is used. If both are zero, there is no timeout.
>
> —— https://pkg.go.dev/net/http#Server.ReadHeaderTimeout

连接建立后不把请求头发完，服务端会一直占着这条连接等下去；并发连接累积后正常请求排不到，形成慢速攻击。`ReadHeaderTimeout` 给读请求头设上限，官方文档说明该值为零且 `ReadTimeout` 也为零时等于不超时，零值不设就落在这一状态。

**正例**

```go
srv := &http.Server{
	Addr:              ":8080",
	Handler:           mux,
	ReadHeaderTimeout: 5 * time.Second,
	ReadTimeout:       15 * time.Second,
	WriteTimeout:      15 * time.Second,
}
```

**反例**

```go
http.ListenAndServe(":8080", mux)
```

**依据**

- https://pkg.go.dev/net/http#Server.ReadHeaderTimeout
- https://github.com/securego/gosec/blob/master/RULES.md

**检测**：golangci-lint gosec(G112, G114)（G112 检出未设 ReadHeaderTimeout，G114 检出无超时支持的 serve 函数）

### 【SHOULD】SEC-011 比较口令、令牌、MAC 用常量时间函数，不用提前返回的比较。

- 归属：安全规约
- 起始版本：Go 1.0

校验 MAC、令牌、签名等秘密值时用 `hmac.Equal`，或 `crypto/subtle` 的 `ConstantTimeCompare`。`bytes.Equal`、`==`、`strings.Compare` 在首个不同处提前返回，不用于秘密比较。

**为什么**

> The time taken is a function of the length of the slices and is independent of the contents.
>
> —— https://pkg.go.dev/crypto/subtle#ConstantTimeCompare

普通比较遇到第一个不同字节就返回，攻击者能逐字节试探：猜对前缀时耗时更长，按耗时筛选即可在有限次数内还原出秘密。`ConstantTimeCompare` 的耗时只与长度相关，掐掉这条时间侧信道。

**正例**

```go
if subtle.ConstantTimeCompare(got, want) != 1 {
	return errUnauthorized
}
```

**反例**

```go
if !bytes.Equal(got, want) {
	return errUnauthorized
}
```

**依据**

- https://pkg.go.dev/crypto/subtle#ConstantTimeCompare
- https://pkg.go.dev/crypto/hmac#Equal

**检测**：人工核对（核对秘密比较处是否用了常量时间函数）

