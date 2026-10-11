# 七、安全规约

### 【MUST】SEC-001 SQL 语句用占位符传参，不把参数值拼进语句。

- 归属：安全规约
- 起始版本：Go 1.8
- 主题：security

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
- 主题：security

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
- 主题：security

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
- 主题：security

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
- 主题：security

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
- 主题：security

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

### 【MUST】SEC-012 用户可控的 SQL 标识符用白名单校验后拼接。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

占位符只覆盖值，不覆盖列名、表名、排序方向等结构片段。

- 动态选择的列名、表名、`ORDER BY` 方向用白名单映射后再进语句
- 占位符参数化的对象是值，标识符不能走 `?` 传递
- 与 SEC-001（值参数化）互补，两者覆盖面不同

**为什么**

> That snippet could complete the SQL statement in unpredictable ways that are dangerous to your application.
>
> —— https://go.dev/doc/database/sql-injection

参数化把值与语句分离，但 ORDER BY 与表名必须是语句的一部分，直接拼接用户输入的标识符与拼接值同样可注入；白名单把「可选结构」限定在代码枚举的集合内，输入只做选择不做拼接。

**正例**

```go
var allowed = map[string]string{"name": "name", "created": "created_at"}
col, ok := allowed[req.Sort]
if !ok {
  return fmt.Errorf("invalid sort field")
}
q := "SELECT id, name FROM users ORDER BY " + col
```

**反例**

```go
q := "SELECT id, name FROM users ORDER BY " + req.Sort
```

**依据**

- https://go.dev/doc/database/sql-injection

**检测**：人工核对（核对动态 SQL 的标识符来源）

### 【MUST】SEC-013 不可信输入强制大小上限，请求体用 http.MaxBytesReader。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

读取路径的输入必须有上限。

- 请求体用 `http.MaxBytesReader` 包裹 `w` 与 `r.Body`
- 解压（gzip、zip）目标大小设上限，防解压炸弹
- 解码 JSON 等格式前先限读，再解析

**为什么**

> MaxBytesReader is similar to io.LimitReader but is intended for limiting the size of incoming request bodies.
>
> —— https://pkg.go.dev/net/http#MaxBytesReader

无上限的读取把内存峰值交给对端：小压缩文件可以展开成 GB 级明文，超大请求体直接耗尽进程内存；上限把「对端想给多大」变成「我们收多大」。

**正例**

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, "payload too large or invalid", http.StatusRequestEntityTooLarge)
  return
}
```

**反例**

```go
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
}
```

**依据**

- https://pkg.go.dev/net/http#MaxBytesReader

**检测**：人工核对（核对解压与读入路径的上限）

### 【MUST】SEC-014 临时文件用 os.CreateTemp 创建。

- 归属：安全规约
- 起始版本：Go 1.16
- 主题：security

临时文件名由系统生成，不手工拼固定路径。

- `os.CreateTemp(dir, pattern)` 生成随机名，并发创建不冲突
- 固定路径可被抢占（symlink 攻击）与并发互写
- 用完即删，清理失败要记录

**为什么**

> Permissive temporary file creation.
>
> —— https://docs.bearer.com/reference/rules/go_gosec_filesystem_tempfile/

固定路径的临时文件名可预测：攻击者可先创建同名符号链接劫持写入目标，或两个实例互踩同一个文件；CreateTemp 由内核保证名字唯一，抢占面消失。

**正例**

```go
f, err := os.CreateTemp("", "report-*.tmp")
if err != nil {
  return err
}
defer os.Remove(f.Name())
```

**反例**

```go
f, err := os.Create("/tmp/report.tmp")
```

**依据**

- https://docs.bearer.com/reference/rules/go_gosec_filesystem_tempfile/
- https://pkg.go.dev/os#CreateTemp

**检测**：golangci-lint gosec(G302, G304)

### 【MUST】SEC-015 tls.Config 显式设置 MinVersion 为 TLS 1.2 以上。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

TLS 最低版本不依赖库默认值。

- 自建 `tls.Config` 时显式写 `MinVersion: tls.VersionTLS12`
- 库默认值随版本变化，显式声明让安全基线固定
- gosec G402 可机检缺省与过旧版本

**为什么**

> Missing TLS MinVersion.
>
> —— https://docs.bearer.com/reference/rules/go_lang_missing_tls_minversion/

TLS 1.0/1.1 有已公开的攻击（BEAST、POODLE 类），且 crypto/tls 的默认最低版本随 Go 版本演进；不显式设置时，安全基线实际由编译时选的 Go 版本决定，升级降级都不受控。

**正例**

```go
cfg := &tls.Config{
  MinVersion: tls.VersionTLS12,
}
```

**反例**

```go
cfg := &tls.Config{}
```

**依据**

- https://docs.bearer.com/reference/rules/go_lang_missing_tls_minversion/
- https://pkg.go.dev/crypto/tls#Config

**检测**：golangci-lint gosec(G402)

### 【MUST】SEC-016 SSH 连接必须校验主机密钥，禁用 InsecureIgnoreHostKey。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

SSH 客户端必须验证对端身份。

- `HostKeyCallback` 用 `ssh.FixedHostKey` 或 known_hosts 校验
- `ssh.InsecureIgnoreHostKey` 放弃中间人防护，仅限测试
- 校验失败的连接必须终止

**为什么**

> Missing verification of host keys.
>
> —— https://docs.bearer.com/reference/rules/go_gosec_crypto_insecure_ignore_host_key/

跳过主机密钥校验后，客户端无法分辨真实服务器与中间人：凭据、命令与传输数据全部暴露给拦截方；gosec G106 与 Bearer 规则都把它列为高危，仅测试环境可豁免。

**正例**

```go
config := &ssh.ClientConfig{
  HostKeyCallback: ssh.FixedHostKey(hostKey),
}
```

**反例**

```go
config := &ssh.ClientConfig{
  HostKeyCallback: ssh.InsecureIgnoreHostKey(),
}
```

**依据**

- https://docs.bearer.com/reference/rules/go_gosec_crypto_insecure_ignore_host_key/

**检测**：golangci-lint gosec(G106)

### 【MUST】SEC-017 口令存储用 Argon2id、bcrypt 等慢哈希。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

口令哈希算法必须为暴力破解设计成本。

- 口令用 `golang.org/x/crypto/argon2` 或 `golang.org/x/crypto/bcrypt`
- 单轮 MD5/SHA-256 允许 GPU 每秒数十亿次尝试
- 哈希参数（cost、内存）按硬件定期复核

**为什么**

> Usage of weak hashing library on a password (MD5).
>
> —— https://docs.bearer.com/reference/rules/go_lang_weak_password_encryption_md5/

通用哈希函数为速度设计，口令哈希要为速度设计开销：Argon2id/bcrypt 的计算成本让离线爆破从秒级变成年级别；数据库泄露时这一层直接决定口令是否等于明文泄露。

**正例**

```go
hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
```

**反例**

```go
sum := md5.Sum([]byte(pwd))
```

**依据**

- https://docs.bearer.com/reference/rules/go_lang_weak_password_encryption_md5/
- https://pkg.go.dev/golang.org/x/crypto/bcrypt

**检测**：人工核对（核对口令存储的哈希算法）

### 【MUST】SEC-018 对称加密用认证模式，nonce 每次随机不复用。

- 归属：安全规约
- 起始版本：Go 1.2
- 主题：security

加密模式必须同时提供机密性与完整性。

- 用 AES-GCM 等 AEAD 模式，禁用 ECB 与未认证模式裸用
- nonce 每次随机生成，同一密钥下不复用
- Go 1.24+ 可用 `cipher.NewGCMWithRandomNonce` 让库管理 nonce

**为什么**

> The nonce must be NonceSize() bytes long and unique for all time, for a given key.
>
> —— https://pkg.go.dev/crypto/cipher#AEAD

无认证模式可被篡改而不被发现，ECB 还会泄露明文模式；GCM 下同一密钥复用 nonce 会同时泄露密钥流并破坏认证，官方文档把 nonce 唯一性写成必须项。

**正例**

```go
gcm, err := cipher.NewGCM(block)
if err != nil {
  return err
}
nonce := make([]byte, gcm.NonceSize())
if _, err := rand.Read(nonce); err != nil {
  return err
}
ct := gcm.Seal(nil, nonce, plaintext, nil)
```

**反例**

```go
iv := make([]byte, 16)
ct := gcm.Seal(nil, iv, plaintext, nil)
```

**依据**

- https://pkg.go.dev/crypto/cipher#AEAD

**检测**：人工核对（核对加密模式与 nonce 生成）

### 【MUST】SEC-019 会话 cookie 设置 HttpOnly、Secure 与 SameSite。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

认证类 cookie 的三个属性缺一不可。

- `HttpOnly` 阻止脚本读取，`Secure` 强制仅 HTTPS 传输，`SameSite` 限制跨站携带
- 会话、令牌、认证类 cookie 全部适用
- 非敏感的功能 cookie 可豁免并在评审中说明

**为什么**

> Missing HTTP Only option in cookie configuration.
>
> —— https://docs.bearer.com/reference/rules/go_lang_cookie_missing_http_only/

缺 HttpOnly 时 XSS 脚本可直接读走会话令牌；缺 Secure 时令牌随明文 HTTP 传输；缺 SameSite 时跨站请求自动携带 cookie 构成 CSRF，三项分别封堵一类攻击路径。

**正例**

```go
http.SetCookie(w, &http.Cookie{
  Name: "session",
  Value: tok,
  HttpOnly: true,
  Secure: true,
  SameSite: http.SameSiteLaxMode,
})
```

**反例**

```go
http.SetCookie(w, &http.Cookie{
  Name: "session",
  Value: tok,
})
```

**依据**

- https://docs.bearer.com/reference/rules/go_lang_cookie_missing_http_only/
- https://docs.bearer.com/reference/rules/go_lang_insecure_cookie/

**检测**：人工核对（核对 cookie 属性设置）

### 【MUST】SEC-020 JWT 校验固定签名算法白名单，不信任令牌 alg 字段。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

可接受的签名算法由代码枚举。

- 解析器显式配置算法白名单（如只允许 `RS256`）
- 不按令牌自带的 `alg` 字段选择验证算法
- `none` 算法除明确受控场景外一律拒绝

**为什么**

> Libraries MUST enable the caller to specify a supported set of algorithms and MUST NOT use any other algorithms when performing cryptographic operations.
>
> —— https://datatracker.ietf.org/doc/html/rfc8725#section-3.1

信任令牌 alg 字段即接受攻击者改写算法：改为 none 绕过验签，RS256 换 HS256 用公钥当 HMAC 密钥伪造签名，两者都是已公开利用的攻击；算法白名单让令牌声称的算法不在集合内即拒绝。

**正例**

```go
tok, err := jwt.Parse(tokenStr, keyFunc, jwt.WithValidMethods([]string{"RS256"}))
```

**反例**

```go
tok, err := jwt.Parse(tokenStr, keyFunc)
```

**依据**

- https://datatracker.ietf.org/doc/html/rfc8725#section-3.1

**检测**：人工核对（核对 JWT 解析配置）

### 【MUST】SEC-021 应用代码不使用 unsafe 绕过内存安全。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

unsafe 只属于经过专门审查的低层库。

- 业务代码不 import `unsafe` 绕过类型与边界检查
- 确需零拷贝等场景封装在经审的低层包内并写明不变量
- 引入 unsafe 的评审按最高级别对待

**为什么**

> Package unsafe contains operations that step around the type safety of Go programs.
>
> —— https://pkg.go.dev/unsafe

Go 的内存安全由编译器与 GC 保证，unsafe 绕过这套保护后，越界与别名问题从「编译失败或 panic」变成「未定义行为」，且逃过所有静态检查；收益只有极低层场景才配得上这个代价。

**正例**

```go
func toBytes(s string) []byte { return []byte(s) }
```

**反例**

```go
func toBytes(s string) []byte {
  return unsafe.Slice(unsafe.StringData(s), len(s))
}
```

**依据**

- https://pkg.go.dev/unsafe
- https://docs.bearer.com/reference/rules/go_gosec_unsafe_unsafe/

**检测**：grep 正则 \bunsafe\.

### 【SHOULD】SEC-005 不设 InsecureSkipVerify 跳过 TLS 证书校验。

- 归属：安全规约
- 起始版本：Go 1.0
- 主题：security

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
- 主题：security

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
- 主题：security

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
- 主题：security

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
- 主题：security

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

