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

### 【MUST】SEC-012 User-controllable SQL identifiers are validated against a whitelist before splicing.

- Category: Security
- Since: Go 1.0
- Tags: security

Placeholders cover values only, never structural fragments such as column names, table names, or sort directions.

- Dynamically chosen columns, tables, and `ORDER BY` directions map through a whitelist before entering the statement
- Parameterized placeholders apply to values; identifiers cannot travel through `?`
- Complementary to SEC-001 (value parameterization); their coverage differs

**Why**

> That snippet could complete the SQL statement in unpredictable ways that are dangerous to your application.
>
> — https://go.dev/doc/database/sql-injection

Parameterization separates values from the statement, but ORDER BY and table names must be part of the statement, so splicing user-controlled identifiers there is as injectable as splicing values; a whitelist restricts the selectable structure to code-enumerated sets, so input chooses rather than composes.

**Good**

```go
var allowed = map[string]string{"name": "name", "created": "created_at"}
col, ok := allowed[req.Sort]
if !ok {
  return fmt.Errorf("invalid sort field")
}
q := "SELECT id, name FROM users ORDER BY " + col
```

**Bad**

```go
q := "SELECT id, name FROM users ORDER BY " + req.Sort
```

**References**

- https://go.dev/doc/database/sql-injection

**Detection**: manual review (Check the identifier sources in dynamic SQL)

### 【MUST】SEC-013 Untrusted input is size-limited; request bodies go through http.MaxBytesReader.

- Category: Security
- Since: Go 1.0
- Tags: security

Every read path bounds its input.

- Wrap request bodies with `http.MaxBytesReader` over `w` and `r.Body`
- Decompression targets (gzip, zip) get size limits against decompression bombs
- Limit the read before decoding JSON and other formats

**Why**

> MaxBytesReader is similar to io.LimitReader but is intended for limiting the size of incoming request bodies.
>
> — https://pkg.go.dev/net/http#MaxBytesReader

Unbounded reads hand the memory peak to the peer: a tiny compressed file can expand into gigabytes of plaintext, and an oversized request body directly exhausts process memory; a limit turns how much the peer wants to send into how much we accept.

**Good**

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, "payload too large or invalid", http.StatusRequestEntityTooLarge)
  return
}
```

**Bad**

```go
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
}
```

**References**

- https://pkg.go.dev/net/http#MaxBytesReader

**Detection**: manual review (Check size limits on decompression and read paths)

### 【MUST】SEC-014 Temporary files are created with os.CreateTemp.

- Category: Security
- Since: Go 1.16
- Tags: security

Temporary file names come from the system, not hand-spliced fixed paths.

- `os.CreateTemp(dir, pattern)` generates a random name that never collides under concurrency
- Fixed paths can be hijacked (symlink attacks) and written by concurrent instances
- Delete after use, and log cleanup failures

**Why**

> Permissive temporary file creation.
>
> — https://docs.bearer.com/reference/rules/go_gosec_filesystem_tempfile/

A predictable temporary file name lets an attacker pre-create a same-name symlink to hijack the write target, or two instances stomp the same file; CreateTemp has the kernel guarantee name uniqueness, removing the hijack surface.

**Good**

```go
f, err := os.CreateTemp("", "report-*.tmp")
if err != nil {
  return err
}
defer os.Remove(f.Name())
```

**Bad**

```go
f, err := os.Create("/tmp/report.tmp")
```

**References**

- https://docs.bearer.com/reference/rules/go_gosec_filesystem_tempfile/
- https://pkg.go.dev/os#CreateTemp

**Detection**: golangci-lint gosec(G302, G304) (Check temporary file creation sites)

### 【MUST】SEC-015 tls.Config sets MinVersion explicitly to TLS 1.2 or above.

- Category: Security
- Since: Go 1.0
- Tags: security

The TLS floor never depends on library defaults.

- Hand-built `tls.Config` writes `MinVersion: tls.VersionTLS12` explicitly
- Library defaults change across versions; explicit declaration pins the security baseline
- gosec G402 detects missing or outdated versions

**Why**

> Missing TLS MinVersion.
>
> — https://docs.bearer.com/reference/rules/go_lang_missing_tls_minversion/

TLS 1.0/1.1 carry published attacks (BEAST, POODLE families), and crypto/tls's default minimum version evolves across Go releases; without an explicit setting the security baseline is actually decided by the Go version chosen at build time, uncontrolled through upgrades and downgrades alike.

**Good**

```go
cfg := &tls.Config{
  MinVersion: tls.VersionTLS12,
}
```

**Bad**

```go
cfg := &tls.Config{}
```

**References**

- https://docs.bearer.com/reference/rules/go_lang_missing_tls_minversion/
- https://pkg.go.dev/crypto/tls#Config

**Detection**: golangci-lint gosec(G402) (Check MinVersion on tls.Config literals)

### 【MUST】SEC-016 SSH connections verify host keys; InsecureIgnoreHostKey is forbidden.

- Category: Security
- Since: Go 1.0
- Tags: security

SSH clients must verify the peer's identity.

- `HostKeyCallback` uses `ssh.FixedHostKey` or known_hosts verification
- `ssh.InsecureIgnoreHostKey` abandons man-in-the-middle protection and is test-only
- Connections failing verification must terminate

**Why**

> Missing verification of host keys.
>
> — https://docs.bearer.com/reference/rules/go_gosec_crypto_insecure_ignore_host_key/

Skipping host key verification leaves the client unable to tell the real server from a man in the middle: credentials, commands, and transferred data are all exposed to the interceptor; gosec G106 and the Bearer rule both rank it high severity, with only test environments exempt.

**Good**

```go
config := &ssh.ClientConfig{
  HostKeyCallback: ssh.FixedHostKey(hostKey),
}
```

**Bad**

```go
config := &ssh.ClientConfig{
  HostKeyCallback: ssh.InsecureIgnoreHostKey(),
}
```

**References**

- https://docs.bearer.com/reference/rules/go_gosec_crypto_insecure_ignore_host_key/

**Detection**: golangci-lint gosec(G106) (Check HostKeyCallback configurations)

### 【MUST】SEC-017 Password storage uses slow hashes such as Argon2id and bcrypt.

- Category: Security
- Since: Go 1.0
- Tags: security

Password hashing algorithms must price out brute force.

- Use `golang.org/x/crypto/argon2` or `golang.org/x/crypto/bcrypt` for passwords
- Single-round MD5/SHA-256 allows billions of guesses per second on GPUs
- Hash parameters (cost, memory) are re-reviewed as hardware evolves

**Why**

> Usage of weak hashing library on a password (MD5).
>
> — https://docs.bearer.com/reference/rules/go_lang_weak_password_encryption_md5/

General hash functions are designed for speed, while password hashing must be designed for cost: the work factor of Argon2id/bcrypt moves offline cracking from seconds to years; when the database leaks, this layer decides whether passwords leak as plaintext.

**Good**

```go
hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
```

**Bad**

```go
sum := md5.Sum([]byte(pwd))
```

**References**

- https://docs.bearer.com/reference/rules/go_lang_weak_password_encryption_md5/
- https://pkg.go.dev/golang.org/x/crypto/bcrypt

**Detection**: manual review (Check the hashing algorithm for password storage)

### 【MUST】SEC-018 Symmetric encryption uses authenticated modes; nonces are random and never reused.

- Category: Security
- Since: Go 1.2
- Tags: security

Encryption modes must provide confidentiality and integrity together.

- Use AEAD modes such as AES-GCM; ECB and bare unauthenticated modes are forbidden
- Generate a fresh random nonce per encryption; never reuse under the same key
- Go 1.24+ offers `cipher.NewGCMWithRandomNonce` so the library manages the nonce

**Why**

> The nonce must be NonceSize() bytes long and unique for all time, for a given key.
>
> — https://pkg.go.dev/crypto/cipher#AEAD

Unauthenticated modes can be tampered with undetected, and ECB leaks plaintext patterns; reusing a nonce under one GCM key leaks the keystream and breaks authentication at once — the official documentation writes nonce uniqueness as a requirement.

**Good**

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

**Bad**

```go
iv := make([]byte, 16)
ct := gcm.Seal(nil, iv, plaintext, nil)
```

**References**

- https://pkg.go.dev/crypto/cipher#AEAD

**Detection**: manual review (Check the encryption mode and nonce generation)

### 【MUST】SEC-019 Session cookies set HttpOnly, Secure, and SameSite.

- Category: Security
- Since: Go 1.0
- Tags: security

All three attributes are required for auth cookies.

- `HttpOnly` blocks script reads, `Secure` forces HTTPS transport, `SameSite` restricts cross-site carrying
- Applies to all session, token, and auth cookies
- Non-sensitive functional cookies may be waived with a review note

**Why**

> Missing HTTP Only option in cookie configuration.
>
> — https://docs.bearer.com/reference/rules/go_lang_cookie_missing_http_only/

Without HttpOnly an XSS script reads the session token directly; without Secure the token travels over plaintext HTTP; without SameSite cross-site requests carry the cookie automatically, forming CSRF — each attribute closes one attack path.

**Good**

```go
http.SetCookie(w, &http.Cookie{
  Name: "session",
  Value: tok,
  HttpOnly: true,
  Secure: true,
  SameSite: http.SameSiteLaxMode,
})
```

**Bad**

```go
http.SetCookie(w, &http.Cookie{
  Name: "session",
  Value: tok,
})
```

**References**

- https://docs.bearer.com/reference/rules/go_lang_cookie_missing_http_only/
- https://docs.bearer.com/reference/rules/go_lang_insecure_cookie/

**Detection**: manual review (Check cookie attribute settings)

### 【MUST】SEC-020 JWT verification pins an algorithm whitelist; the token's alg field is not trusted.

- Category: Security
- Since: Go 1.0
- Tags: security

Acceptable signing algorithms are enumerated in code.

- Parsers explicitly configure an algorithm whitelist (such as `RS256` only)
- The verification algorithm is never chosen from the token's own `alg` field
- The `none` algorithm is rejected except in explicitly controlled scenarios

**Why**

> Libraries MUST enable the caller to specify a supported set of algorithms and MUST NOT use any other algorithms when performing cryptographic operations.
>
> — https://datatracker.ietf.org/doc/html/rfc8725#section-3.1

Trusting the token's alg field accepts attacker-chosen algorithms: none bypasses verification, and RS256 switched to HS256 forges signatures using the public key as the HMAC secret — both are publicly exploited attacks; a whitelist rejects any token claiming an algorithm outside the set.

**Good**

```go
tok, err := jwt.Parse(tokenStr, keyFunc, jwt.WithValidMethods([]string{"RS256"}))
```

**Bad**

```go
tok, err := jwt.Parse(tokenStr, keyFunc)
```

**References**

- https://datatracker.ietf.org/doc/html/rfc8725#section-3.1

**Detection**: manual review (Check JWT parser configurations)

### 【MUST】SEC-021 Application code does not use unsafe to bypass memory safety.

- Category: Security
- Since: Go 1.0
- Tags: security

unsafe belongs to specially reviewed low-level libraries.

- Business code does not import `unsafe` to bypass type and bounds checks
- Zero-copy scenarios are wrapped in reviewed low-level packages with invariants documented
- Introducing unsafe gets the highest review severity

**Why**

> Package unsafe contains operations that step around the type safety of Go programs.
>
> — https://pkg.go.dev/unsafe

Go's memory safety rests on the compiler and GC; bypassing them with unsafe turns out-of-bounds and aliasing problems from compile failures or panics into undefined behavior, invisible to every static check — the payoff only suits the lowest-level scenarios.

**Good**

```go
func toBytes(s string) []byte { return []byte(s) }
```

**Bad**

```go
func toBytes(s string) []byte {
  return unsafe.Slice(unsafe.StringData(s), len(s))
}
```

**References**

- https://pkg.go.dev/unsafe
- https://docs.bearer.com/reference/rules/go_gosec_unsafe_unsafe/

**Detection**: grep regex \bunsafe\. (Grep for unsafe package usage)

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

