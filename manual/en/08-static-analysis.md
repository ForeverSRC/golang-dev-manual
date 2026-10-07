# 8. Static Analysis

## (1) Configuration and Running

### 【MUST】TOOL-001 Use golangci-lint uniformly for static code analysis.

- Category: Static Analysis/Configuration and Running
- Since: Go 1.0

Keep a single .golangci.yml across the whole repository, and run the same command locally and in CI; the `golangci-lint` version is pinned by the tool directive in go.mod, not relying on a global installation on each machine.

**Why**

> Golangci-lint is a fast linters runner for Go. It runs linters in parallel, uses caching, supports YAML configuration, integrates with all major IDEs, and includes over a hundred linters.
>
> — https://golangci-lint.run/

Bare linters of different versions on different machines by different people produce inconsistent results, so reviewers cannot reach the same conclusion; only one centralized config plus a pinned version makes problems reproducible.

**Good**

```makefile
lint:
	golangci-lint run ./...
```

**Bad**

```makefile
lint:
	go vet ./...
	staticcheck ./...
```

**References**

- https://golangci-lint.run/

**Detection**: manual review (Check that the repository root has .golangci.yml and that the lint entry point invokes golangci-lint)

### 【MUST】TOOL-009 Name the checker and state the reason when suppressing a warning.

- Category: Static Analysis/Configuration and Running
- Since: Go 1.0

After confirming a false positive, write //nolint:<checker> // reason at the triggering site, where the checker name and the reason are both indispensable; do not write a bare //nolint or //nolint:all, since they do not indicate which check is being suppressed; also do not turn off checks wholesale at the top of a file or package. Enable nolintlint in the configuration, enforced by the three settings require-specific, require-explanation, and allow-unused.

**Why**

> You may add a comment explaining or justifying why a `nolint` directive is being used on the same line as the flag itself:
>
> — https://golangci-lint.run/docs/linters/false-positives/#nolint-directive

A suppression is a hole opened in the checkers; writing the checker and the reason makes clear during review which rule was bypassed and what the cost of bypassing it is; a bare //nolint or //nolint:all opens the hole on all checks, so analyzers added later are also let through and no one sees the problems anymore.

**Good**

```go
raw, err := os.ReadFile(path) //nolint:gosec // path is specified explicitly by the caller
```

**Bad**

```go
raw, err := os.ReadFile(path) //nolint
```

**References**

- https://golangci-lint.run/docs/linters/false-positives/#nolint-directive
- https://golangci-lint.run/docs/linters/configuration/#nolintlint

**Detection**: golangci-lint nolintlint (Configured as require-specific + require-explanation + allow-unused, blocking suppressions that do not name a checker, lack a reason, or are no longer effective)

## (2) Required Rules

### 【MUST】TOOL-006 Enable at least the standard set of golangci-lint.

- Category: Static Analysis/Required Rules
- Since: Go 1.0

The standard set is the five checkers enabled by default in `golangci-lint`; do not turn them off in the configuration, and do not keep only some of them:

- errcheck: unhandled error return values
- govet: suspicious constructs, equivalent to `go vet`
- ineffassign: assignments whose value is never read again
- staticcheck: the diagnostic rules of staticcheck
- unused: unused identifiers

**Why**

> Golangci-lint can be used with zero configuration. By default, the following linters are enabled: errcheck, govet, ineffassign, staticcheck, unused.
>
> — https://golangci-lint.run/docs/welcome/quick-start/

All five target defects that can be determined with certainty: missed errors, suspicious constructs that are wrong yet still compile, intermediate values computed but never used, and dead code that can never be reached. Turning off any one of them amounts to giving up machine interception for that category of defect, falling back to review by human memory.

**Good**

```yaml
linters:
  default: standard
```

**Bad**

```yaml
linters:
  default: none
  enable:
    - gocritic
```

**References**

- https://golangci-lint.run/docs/welcome/quick-start/

**Detection**: manual review (Check that .golangci.yml does not turn off the standard set)

## (3) Recommended Rules

### 【SHOULD】TOOL-002 Enable modernize in golangci-lint to keep up with modern Go idioms.

- Category: Static Analysis/Recommended Rules
- Since: Go 1.22

After a Go version upgrade, modernize points out the places where old code can be simplified; without it enabled, these idioms rely solely on human memory.

**Why**

> modernize: A suite of analyzers that suggest simplifications to Go code, using modern language and library features.
>
> — https://golangci-lint.run/docs/linters/configuration/#modernize

The standard library and language features keep filling in gaps, and old code becomes technical debt version by version; letting analyzers flag it in bulk costs less than refactoring after the fact.

**Good**

```go
for range n {
	handle()
}
```

**Bad**

```go
for i := 0; i < n; i++ {
	handle()
}
```

**References**

- https://golangci-lint.run/docs/linters/configuration/#modernize

**Detection**: manual review (Check that linters.enable in .golangci.yml includes modernize)

### 【SHOULD】TOOL-003 Enable testifylint alongside testify when testify is used.

- Category: Static Analysis/Recommended Rules
- Since: Go 1.0

testifylint blocks high-frequency misuses such as judging errors with `assert.Nil` and reversed argument order; projects that do not use `testify` need not enable it.

**Why**

> testify is the most popular Golang testing framework in recent years. But it has a terrible ambiguous API in places, and the purpose of this linter is to protect you from annoying mistakes.
>
> — https://github.com/Antonboom/testifylint

Some misuse of the testify API raises no compile-time error, for example `assert.Nil(t, err)` fails to judge correctly when err is an interface; handing this to a linter is more reliable than relying on reviewers' memory.

**Good**

```go
require.NoError(t, err)
```

**Bad**

```go
assert.Nil(t, err)
```

**References**

- https://github.com/Antonboom/testifylint

**Detection**: manual review (Check that linters.enable in .golangci.yml includes testifylint)

### 【SHOULD】TOOL-007 Beyond the standard set, enable the commonly used supplementary analyzers.

- Category: Static Analysis/Recommended Rules
- Since: Go 1.0

Beyond the standard set, enable the following analyzers:

- errorlint: writing problems related to error chains
- nilerr: checking err != nil yet returning nil
- bodyclose: an HTTP response body is not closed
- noctx: a request carries no `context.Context`
- nolintlint: a suppression does not name a checker, lacks a reason, or is no longer effective
- gosec: security risks
- gocritic: defect and performance diagnostics
- exhaustive: an enum switch is not exhaustive
- copyloopvar: redundant loop variable copies since Go 1.22
- testpackage: same-package tests
- funcorder: the order of exported and unexported functions
- misspell: spelling mistakes

See the separate clauses in the same section for modernize and testifylint. Those that do not fit the shape of the project may be left disabled, with the reason written into a configuration comment or the review record.

**Why**

> To see a list of supported linters and which linters are enabled/disabled: golangci-lint help linters
>
> — https://golangci-lint.run/docs/linters/

Nothing beyond the standard set is enabled by default; not turning it on means not having it. Each of these analyzers targets a specific class of defect or high-frequency deviation, with controllable false positives; the list includes only general items unrelated to any project, and project-specific checks are added within their own projects.

**Good**

```yaml
linters:
  default: standard
  enable:
    - errorlint
    - gosec
    - nolintlint
    - testpackage
```

**Bad**

```yaml
linters:
  default: standard
  # none of the supplementary analyzers is enabled
```

**References**

- https://golangci-lint.run/docs/linters/

**Detection**: manual review (Check whether the enable list covers the list of analyzers and whether disabled ones have a reason)

