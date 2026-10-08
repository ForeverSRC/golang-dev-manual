# 7. Security

### 【MUST】SEC-001 SQL statements pass parameters through placeholders; do not concatenate parameter values into the statement.

- Category: Security
- Since: Go 1.8
- Tags: security

Parameter values for queries and writes are passed through `database/sql` placeholders, not concatenated into the SQL text. Identifiers that placeholders cannot cover, such as table names, column names, and sort fields, are validated against a whitelist before being concatenated.

**Why**

> You can avoid an SQL injection risk by providing SQL parameter values as `sql` package function arguments. Many functions in the `sql` package provide parameters for the SQL statement and for values to be used in that statement's parameters (others provide a parameter for a prepared statement and parameters).
>
> — https://go.dev/doc/database/sql-injection

Once parameters are handed to the `sql` package, the driver sends the statement and the parameters separately, so parameter values are not parsed as SQL syntax. Switching to `fmt.Sprintf` to concatenate values into the statement before handing it to `Query` means the whole SQL text is already formed, and fragments supplied by the caller take part in the syntax directly; passing `1 OR 1=1` as `id` bypasses the condition, and the official documentation marks this usage as a SECURITY RISK.

**Good**

```go
rows, err := db.QueryContext(ctx, "SELECT id, name FROM user WHERE id = ?", id)
```

**Bad**

```go
rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT id, name FROM user WHERE id = %s", id))
```

**References**

- https://go.dev/doc/database/sql-injection
- https://go.dev/doc/database/prepared-statements

**Detection**: golangci-lint gosec(G201, G202) (At each hit, confirm whether the parameter value comes from external input.)

### 【MUST】SEC-002 Pass the executable and arguments of external commands to exec separately; do not build a command string.

- Category: Security
- Since: Go 1.7
- Tags: security

When using `exec.Command` or `exec.CommandContext`, pass the executable and each argument as separate arguments, and do not build a complete command line through `sh -c` or `cmd /c`. When glob expansion or pipelines are needed, use `filepath.Glob` or implement them in the program; when the shell cannot be avoided, escape the concatenated external input first.

**Why**

> Unlike the "system" library call from C and other languages, the os/exec package intentionally does not invoke the system shell and does not expand any glob patterns or handle other expansions, pipelines, or redirections typically done by shells.
>
> — https://pkg.go.dev/os/exec

`os/exec` starts a process with an argument array and does not go through the shell, so `;`, `|`, and `&&` in arguments are treated as ordinary characters. Once changed to concatenation with `sh -c`, external input re-enters shell syntax, and a fragment of `; rm -rf` can change the command actually executed; the official documentation also warns, in the glob expansion section, to be careful about escaping external input when invoking the shell directly.

**Good**

```go
cmd := exec.CommandContext(ctx, "git", "log", "--oneline", "-n", n)
```

**Bad**

```go
cmd := exec.CommandContext(ctx, "sh", "-c", "git log --oneline -n "+n)
```

**References**

- https://pkg.go.dev/os/exec

**Detection**: golangci-lint gosec(G204) (At each hit, confirm whether external input was concatenated into the command string.)

### 【MUST】SEC-003 User-controllable file paths are confined to an allowed directory.

- Category: Security
- Since: Go 1.24
- Tags: security

When external input takes part in constructing a file path, first use `os.OpenRoot` to open an `*os.Root` on the allowed root directory, then access files through its `Open`, `Create`, `ReadFile` and other methods. When `os.Root` does not apply, at least run `filepath.Clean` and verify that the result still falls within the root directory prefix, and handle symbolic links separately.

**Why**

> The os.OpenRoot function opens a directory and returns an os.Root. Methods on os.Root operate within the directory and do not permit paths that refer to locations outside the directory, including ones that follow symbolic links out of the directory.
>
> — https://go.dev/doc/go1.24#directory-limited-filesystem-access

`../` in a path and symbolic links pointing outside can lead `os.Open` out of the intended directory, reading or overwriting files outside the root directory. `os.Root` takes the root directory as its boundary and validates the final resolved result on every operation, reporting an error immediately for out-of-bounds paths; a hand-written check must compare prefixes after `Clean` and also block symbolic links on its own, which easily misses cases.

**Good**

```go
root, err := os.OpenRoot(dataDir)
if err != nil {
	return err
}
defer root.Close()

f, err := root.Open(filepath.Join("reports", name))
```

**Bad**

```go
f, err := os.Open(filepath.Join(dataDir, name))
```

**References**

- https://go.dev/doc/go1.24#directory-limited-filesystem-access
- https://pkg.go.dev/os#Root

**Detection**: golangci-lint gosec(G304) (At each hit, confirm whether the path contains external input and whether it is confined within the root directory.)

### 【MUST】SEC-004 Random numbers for security purposes come from crypto/rand, not math/rand.

- Category: Security
- Since: Go 1.24
- Tags: security

When generating random values related to security, such as tokens, session IDs, keys, salts, and verification codes, use `Read`, `Text` and similar functions from `crypto/rand`. `math/rand` is used only for scenarios that do not involve security, such as simulation, sampling, and shuffling.

**Why**

> Package rand implements pseudo-random number generators suitable for tasks such as simulation, but it should not be used for security-sensitive work.
>
> — https://pkg.go.dev/math/rand

The output of `math/rand` is determined by the seed and can be inferred; the official documentation states explicitly that it is unsuitable for security purposes and points out that the sequence is easy to guess. Using it to generate tokens lets an attacker derive subsequent values from a small amount of output; `crypto/rand` draws on the secure random source provided by the operating system, and its output is unpredictable.

**Good**

```go
import "crypto/rand"

token := make([]byte, 32)
rand.Read(token)
```

**Bad**

```go
import "math/rand"

token := make([]byte, 32)
rand.Read(token)
```

**References**

- https://pkg.go.dev/math/rand
- https://pkg.go.dev/crypto/rand

**Detection**: golangci-lint gosec(G404) (A math/rand call is hit; confirm whether it is used in a security scenario.)

### 【MUST】SEC-006 Do not use broken algorithms such as MD5, SHA-1, DES, and RC4 for security purposes.

- Category: Security
- Since: Go 1.0
- Tags: security

Security scenarios such as password hashing, signatures, message authentication, and key derivation must not introduce `crypto/md5`, `crypto/sha1`, `crypto/des`, `crypto/rc4`, `golang.org/x/crypto/md4`, or `golang.org/x/crypto/ripemd160`. These algorithms remain usable in non-security scenarios, such as file checksums and content addressing.

**Why**

> MD5 is cryptographically broken and should not be used for secure applications.
>
> — https://pkg.go.dev/crypto/md5

The documentation of each of these packages states that the algorithm has been broken and should no longer be used in security contexts. Practical collision constructions exist for MD5 and SHA-1, so two different inputs can hash to the same value; when they are used for integrity checks or signatures, content that has been replaced still passes. The key space and keystream weaknesses of DES and RC4 can be broken in practice.

**Good**

```go
sum := sha256.Sum256(data)
```

**Bad**

```go
sum := md5.Sum(data)
```

**References**

- https://pkg.go.dev/crypto/md5
- https://pkg.go.dev/crypto/sha1
- https://pkg.go.dev/crypto/des
- https://pkg.go.dev/crypto/rc4

**Detection**: golangci-lint gosec(G401, G405, G501, G502, G503, G505)

### 【MUST】SEC-007 Do not hard-code keys, passwords, or tokens in source code.

- Category: Security
- Since: Go 1.0
- Tags: security

Sensitive values such as keys, passwords, access tokens, and connection strings are injected through environment variables, configuration files, or a secret management service, not written as literals in source code. Test credentials go into test fixtures, clearly distinguished, and do not enter production code or the version repository.

**Why**

> G101 — Look for hardcoded credentials
>
> — https://github.com/securego/gosec/blob/master/RULES.md

Credentials written into source code enter version history with the commit, so anyone who reads the repository can obtain them; rotating a key requires changing code and releasing again, which in practice often does not happen for a long time, and the impact of a single leak keeps expanding. gosec's G101 identifies such assignments by identifier name and literal entropy, and each hit should be changed to injection from outside.

**Good**

```go
apiKey := os.Getenv("API_KEY")
```

**Bad**

```go
apiKey := "sk-live-9f3c1d2e4b5a6c7d"
```

**References**

- https://github.com/securego/gosec/blob/master/RULES.md

**Detection**: golangci-lint gosec(G101) (At each hit, confirm whether it is a real credential; for a false positive, annotate the reason with #nosec G101.)

### 【SHOULD】SEC-005 Do not set InsecureSkipVerify to skip TLS certificate verification.

- Category: Security
- Since: Go 1.0
- Tags: security

Keep `tls.Config` at its default certificate chain and host name verification, with `InsecureSkipVerify` left false. For self-signed certificate scenarios, add the CA to `RootCAs`, or use `VerifyPeerCertificate` or `VerifyConnection` for custom verification. It may be enabled temporarily in test environments and one-off verification scripts, with the reason stated.

**Why**

> If InsecureSkipVerify is true, crypto/tls accepts any certificate presented by the server and any host name in that certificate. In this mode, TLS is susceptible to machine-in-the-middle attacks unless custom verification is used.
>
> — https://pkg.go.dev/crypto/tls#Config

Certificate verification is the basis on which TLS confirms the identity of the peer; once it is turned off, the client accepts any certificate and any host name, and a machine-in-the-middle can decrypt or tamper with traffic with a self-signed certificate. When verification is skipped to connect to a production service, the encrypted channel only guards against passive eavesdropping, and active attacks are no longer impeded.

**Good**

```go
cfg := &tls.Config{}
```

**Bad**

```go
cfg := &tls.Config{InsecureSkipVerify: true}
```

**References**

- https://pkg.go.dev/crypto/tls#Config

**Detection**: golangci-lint gosec(G402) (At each hit, confirm whether non-test code still skips verification.)

### 【SHOULD】SEC-008 Grant only the required permissions when creating files and directories; do not use globally writable values.

- Category: Security
- Since: Go 1.16
- Tags: security

The permission arguments of `os.OpenFile`, `os.WriteFile`, `os.Mkdir`, and `os.MkdirAll` are chosen as the minimum needed; do not use values such as `0777` and `0666` that open write permission to all users, and use `0600` for files containing sensitive data. When exact permissions are needed, set them explicitly with `Chmod` after creation; the process umask filters the passed value first.

**Why**

> File and directory permission rules can be configured with stricter maximum permissions
>
> — https://github.com/securego/gosec/blob/master/RULES.md

Granting 0777 at creation makes the file writable by all users, so other accounts on the machine or processes in a shared environment can replace the content; when it is written into a temporary directory, or read as configuration or a script, this can be used for privilege escalation or injection. gosec's G301, G302, and G306 target these creation and permission-change calls, and the thresholds can be tightened per project.

**Good**

```go
err := os.WriteFile(path, data, 0o600)
```

**Bad**

```go
err := os.WriteFile(path, data, 0o777)
```

**References**

- https://github.com/securego/gosec/blob/master/RULES.md

**Detection**: golangci-lint gosec(G301, G302, G306) (The thresholds can be tightened as needed in the gosec configuration in .golangci.yml.)

### 【SHOULD】SEC-009 Use html/template, not text/template, to output HTML.

- Category: Security
- Since: Go 1.0
- Tags: security

When generating text that requires escaping, such as HTML, XML, and JS, use `html/template`, which escapes automatically according to the context the data appears in. `text/template` is used only for plain text output; when HTML can only be produced with `text/template`, escape the data yourself before output.

**Why**

> Package template (html/template) implements data-driven templates for generating HTML output safe against code injection. It provides the same interface as text/template and should be used instead of text/template whenever the output is HTML.
>
> — https://pkg.go.dev/html/template

`text/template` writes data into the output as-is, so a `<script>` fragment in the data is executed by the browser as a tag, constituting cross-site scripting. `html/template` applies the corresponding escaping according to whether the data falls in a tag, an attribute, a URL, or a JS string, with the template code unchanged; switching the package blocks the injection.

**Good**

```go
import "html/template"

t, err := template.New("page").Parse(pageTemplate)
```

**Bad**

```go
import "text/template"

t, err := template.New("page").Parse(pageTemplate)
```

**References**

- https://pkg.go.dev/html/template

**Detection**: golangci-lint gosec(G203) (A hit on G203 writes unescaped data into an HTML template.)

### 【SHOULD】SEC-010 HTTP services explicitly set the read header, read, and write timeouts.

- Category: Security
- Since: Go 1.8
- Tags: security

`http.Server` should at least set `ReadHeaderTimeout`, and configure `ReadTimeout`, `WriteTimeout`, and `IdleTimeout` according to the business needs. When using convenience functions without timeouts such as `http.ListenAndServe`, switch to a custom `http.Server`.

**Why**

> If ReadHeaderTimeout is zero, the value of ReadTimeout is used. If both are zero, there is no timeout.
>
> — https://pkg.go.dev/net/http#Server.ReadHeaderTimeout

If a connection does not finish sending its request headers after being established, the server keeps holding that connection waiting; once concurrent connections accumulate, normal requests cannot get through, forming a slow attack. `ReadHeaderTimeout` sets an upper bound on reading request headers, and the official documentation states that when this value is zero and `ReadTimeout` is also zero there is no timeout — leaving the zero value unset falls into that state.

**Good**

```go
srv := &http.Server{
	Addr:              ":8080",
	Handler:           mux,
	ReadHeaderTimeout: 5 * time.Second,
	ReadTimeout:       15 * time.Second,
	WriteTimeout:      15 * time.Second,
}
```

**Bad**

```go
http.ListenAndServe(":8080", mux)
```

**References**

- https://pkg.go.dev/net/http#Server.ReadHeaderTimeout
- https://github.com/securego/gosec/blob/master/RULES.md

**Detection**: golangci-lint gosec(G112, G114) (G112 detects a missing ReadHeaderTimeout; G114 detects a serve function without timeout support.)

### 【SHOULD】SEC-011 Compare passwords, tokens, and MACs with constant-time functions, not comparisons that return early.

- Category: Security
- Since: Go 1.0
- Tags: security

When verifying secret values such as MACs, tokens, and signatures, use `hmac.Equal`, or `ConstantTimeCompare` from `crypto/subtle`. `bytes.Equal`, `==`, and `strings.Compare` return early at the first difference and are not used for secret comparison.

**Why**

> The time taken is a function of the length of the slices and is independent of the contents.
>
> — https://pkg.go.dev/crypto/subtle#ConstantTimeCompare

An ordinary comparison returns at the first differing byte, letting an attacker probe byte by byte: a correct prefix takes longer, so filtering by elapsed time can recover the secret within a limited number of attempts. `ConstantTimeCompare` takes time related only to length, cutting off this timing side channel.

**Good**

```go
if subtle.ConstantTimeCompare(got, want) != 1 {
	return errUnauthorized
}
```

**Bad**

```go
if !bytes.Equal(got, want) {
	return errUnauthorized
}
```

**References**

- https://pkg.go.dev/crypto/subtle#ConstantTimeCompare
- https://pkg.go.dev/crypto/hmac#Equal

**Detection**: manual review (Check whether secret comparison sites use a constant-time function.)

