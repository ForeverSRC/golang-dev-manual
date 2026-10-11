# 10. Appendix

## 1 Clause Index

| ID | Level | Category | Summary |
| --- | --- | --- | --- |
| ABSTRACT-001 | SHOULD | Interfaces and Design/Interface Definition and Placement | Introduce abstraction according to actual need, not spread out in advance. |
| CODE-001 | SHOULD | Errors and Logging/Error Codes | For external error codes, reuse the existing standard classification; do not build a separate numbering scheme. |
| CODE-002 | SHOULD | Errors and Logging/Error Codes | Error messages are not a contract; test errors only by code and type. |
| CODE-003 | SHOULD | Errors and Logging/Error Codes | When propagating an error across services, translate the error code; do not pass through the downstream raw error. |
| CODE-004 | MAY | Errors and Logging/Error Codes | Error strings are lowercase, unpunctuated; sentinel errors carry a package prefix. |
| COMMENT-001 | SHOULD | Programming Conventions/Comments | Comments record only what the code cannot express and never restate what the code does. |
| COMMENT-002 | SHOULD | Programming Conventions/Comments | A doc comment begins with the name being described and is written as a full sentence. |
| COMPOSE-001 | SHOULD | Interfaces and Design/Composition and Reuse | Express reusable capability through composition, not inheritance-style hierarchies. |
| COMPOSE-002 | SHOULD | Interfaces and Design/Composition and Reuse | Embedding must not change the outer type's zero value, copy semantics, or public interface. |
| CONC-001 | MUST | Programming Conventions/Concurrency | Before starting a goroutine, make its exit timing clear, and let the starter wait for it to finish. |
| CONC-002 | MUST | Programming Conventions/Concurrency | A channel is closed only by the sender, never by the receiver. |
| CONC-003 | SHOULD | Programming Conventions/Concurrency | A public API keeps synchronous semantics and leaves concurrency to the caller. |
| CONC-004 | SHOULD | Programming Conventions/Concurrency | Shared state and the lock protecting it live in the same struct and are read and written only through methods. |
| CONC-005 | SHOULD | Programming Conventions/Concurrency | Channel parameters in function signatures are typed as directional. |
| CONC-006 | SHOULD | Programming Conventions/Concurrency | Buffered channels with capacity above 1 carry a comment on sizing and overflow behavior. |
| CONC-007 | SHOULD | Programming Conventions/Concurrency | Selectable cancelable paths include a ctx.Done() branch. |
| CONC-008 | SHOULD | Programming Conventions/Concurrency | A select inside a for loop does not busy-wait via a default branch. |
| CONC-009 | SHOULD | Programming Conventions/Concurrency | Periodic ticks in loops use time.NewTicker with Stop. |
| CONC-010 | MUST | Programming Conventions/Concurrency | Wait with wg.Go, or place manual wg.Add before the go statement. |
| CONC-011 | SHOULD | Programming Conventions/Concurrency | Critical sections exclude I/O, network calls, and channel operations. |
| CONC-012 | MUST | Programming Conventions/Concurrency | Do not request Lock while holding RLock on the same RWMutex. |
| CONC-013 | SHOULD | Programming Conventions/Concurrency | Simple counters and flags use typed atomic types. |
| CONC-014 | SHOULD | Programming Conventions/Concurrency | Do not use sync.Map in write-heavy scenarios. |
| CONC-015 | SHOULD | Programming Conventions/Concurrency | Batch concurrent tasks set a concurrency limit. |
| CONC-016 | SHOULD | Programming Conventions/Concurrency | sync.Pool objects are Reset before being put back. |
| CONC-017 | SHOULD | Programming Conventions/Concurrency | Long-lived goroutine entry points recover and log panics. |
| CONST-001 | MUST | Programming Conventions/Constants and Enums | Constant names use MixedCaps, not all uppercase or a k prefix. |
| CONST-002 | SHOULD | Programming Conventions/Constants and Enums | Constants are named by their role, not by their value. |
| CONST-003 | MAY | Programming Conventions/Constants and Enums | Enumeration values that increase from one source are generated with iota. |
| CONST-004 | SHOULD | Programming Conventions/Constants and Enums | Leave the zero value of an iota enum as an Unknown sentinel, or start at iota+1. |
| CTRL-001 | SHOULD | Programming Conventions/Control Flow | Handle errors and edge cases first and return early, keeping the normal path unindented. |
| CTRL-002 | SHOULD | Programming Conventions/Control Flow | Do not explicitly copy the loop variable any longer. |
| CTRL-003 | MUST | Programming Conventions/Control Flow | Do not defer resource release inside a loop body. |
| CTRL-004 | MUST | Programming Conventions/Control Flow | HTTP handlers return immediately after writing an error response. |
| CTRL-005 | SHOULD | Programming Conventions/Control Flow | Inner scopes do not shadow the outer err with :=. |
| CTX-001 | SHOULD | Interfaces and Design/Passing context | Do not store context.Context in a struct; pass it as the first parameter of a function. |
| CTX-002 | MUST | Interfaces and Design/Passing context | Do not pass a nil context to a function; use context.TODO() as a placeholder. |
| CTX-003 | SHOULD | Interfaces and Design/Passing context | Use an unexported custom type for context.Value keys. |
| CTX-004 | SHOULD | Interfaces and Design/Passing context | Put only request-scoped data in context.Value, not optional parameters. |
| CTX-005 | MUST | Interfaces and Design/Passing context | The cancel returned by WithCancel and WithTimeout must be called. |
| CTX-006 | SHOULD | Interfaces and Design/Passing context | context.Background is created at entry points only; mid-chain calls do not substitute it for the ctx parameter. |
| CTX-007 | SHOULD | Interfaces and Design/Passing context | APIs offering *Context variants are always called through the variant with ctx. |
| CTX-008 | SHOULD | Interfaces and Design/Passing context | Post-request background tasks use context.WithoutCancel. |
| CTX-009 | SHOULD | Interfaces and Design/Passing context | Retry loops check ctx.Err each round and stop once canceled. |
| DATA-001 | SHOULD | Programming Conventions/Data Structures | Declare empty slices with var to get a nil slice. |
| DATA-002 | SHOULD | Programming Conventions/Data Structures | Struct literals are initialized with field names rather than listed by position. |
| DATA-003 | MUST | Programming Conventions/Data Structures | Type assertions use the comma-ok form, not the single-value form. |
| DATA-004 | MUST | Programming Conventions/Data Structures | Initialize a nil map before writing to it. |
| DATA-005 | MUST | Programming Conventions/Data Structures | Do not rely on map iteration order; sort keys when order matters. |
| DATA-006 | SHOULD | Programming Conventions/Data Structures | Check the value range before a narrowing integer conversion. |
| DATA-007 | MUST | Programming Conventions/Data Structures | Struct fields that participate in JSON serialization must be exported. |
| DB-001 | MUST | Programming Conventions/Database Access | Check rows.Err and Close rows after iterating sql.Rows. |
| DB-002 | MUST | Programming Conventions/Database Access | Wrap multi-statement writes in a transaction and defer Rollback right after Begin. |
| DB-003 | SHOULD | Programming Conventions/Database Access | Receive nullable columns into pointers or sql.Null. |
| DB-004 | MUST | Programming Conventions/Database Access | Match sql.ErrNoRows with errors.Is. |
| DEP-001 | MUST | Dependencies and Project Layout/Dependency Management | Leave go.mod and go.sum to the go command, and run go mod tidy before committing. |
| DEP-002 | MUST | Dependencies and Project Layout/Dependency Management | Commit go.sum to version control, and do not list it in .gitignore. |
| DEP-003 | SHOULD | Dependencies and Project Layout/Dependency Management | Do not keep replace directives pointing at local directories or private forks long term. |
| DEP-004 | SHOULD | Dependencies and Project Layout/Dependency Management | Declare tool dependencies in go.mod's tool directive rather than relying on global installs. |
| DEP-005 | SHOULD | Dependencies and Project Layout/Dependency Management | Do not bring in a third-party dependency for a small capability the standard library already covers. |
| DI-001 | SHOULD | Interfaces and Design/Dependency Injection and Wiring | Use wire for dependency wiring, generated at compile time. |
| DI-002 | SHOULD | Interfaces and Design/Dependency Injection and Wiring | Pass dependencies in explicitly through constructors, and do not read or write package-level mutable variables. |
| ERR-001 | MUST | Errors and Logging/Error Handling | Wrap with %w when the error chain must be preserved; do not use %v or %s. |
| ERR-003 | MUST | Errors and Logging/Error Handling | Test the kind of an error with errors.Is and errors.As, not with == or a type assertion. |
| ERR-004 | SHOULD | Errors and Logging/Error Handling | Return an error for expected failures; do not use panic. |
| ERR-005 | SHOULD | Errors and Logging/Error Handling | Return failures as an error or an ok value; do not signal them with a special value. |
| ERR-006 | MUST | Errors and Logging/Error Handling | Handle the same error only once: wrap and return it, or log it in place, never both. |
| ERR-007 | MUST | Errors and Logging/Error Handling | error return values are never discarded with _. |
| ERR-008 | MAY | Errors and Logging/Error Handling | Independent failures are merged with errors.Join. |
| ERR-009 | MUST | Errors and Logging/Error Handling | Library code and request paths never call os.Exit or log.Fatal. |
| ERR-010 | SHOULD | Errors and Logging/Error Handling | Close and Flush errors on writable resources must be checked. |
| ERR-011 | MUST | Errors and Logging/Error Handling | External error messages stay generic; details go to server-side logs only. |
| FUNC-001 | MUST | Programming Conventions/Functions and Methods | Exported functions and methods come at the front of a file, unexported ones at the back. |
| FUNC-002 | SHOULD | Programming Conventions/Functions and Methods | Use named results when the meaning of the return values is not self-evident. |
| FUNC-003 | SHOULD | Programming Conventions/Functions and Methods | Collect more than 4 arguments or more than 3 results into a struct. |
| FUNC-004 | MUST | Programming Conventions/Functions and Methods | Nil branches of interface-returning functions return nil explicitly. |
| IFACE-001 | SHOULD | Interfaces and Design/Interface Definition and Placement | Place interfaces according to the number of implementations: a single implementation shares the package with the implementation, while multiple implementations are defined by the consumer. |
| IFACE-002 | SHOULD | Interfaces and Design/Interface Definition and Placement | Declare a type's implementation of an interface with a compile-time assertion. |
| IFACE-003 | SHOULD | Interfaces and Design/Interface Definition and Placement | Do not define pointers to interfaces. |
| IFACE-004 | MAY | Interfaces and Design/Interface Definition and Placement | Keep interfaces small; larger contracts are composed from smaller ones. |
| IFACE-005 | SHOULD | Interfaces and Design/Interface Definition and Placement | Idiomatic method signatures match standard interfaces; do not invent variants. |
| IFACE-006 | SHOULD | Interfaces and Design/Interface Definition and Placement | Mark exported symbols Deprecated and keep them for one release cycle before deletion. |
| LAYER-001 | SHOULD | Dependencies and Project Layout/Layering | Layer code by responsibility so that dependencies point one-way toward the inner layers. |
| LAYOUT-001 | MUST | Dependencies and Project Layout/Project Layout | Put go.mod and the source code at the module root, without an extra src/ layer. |
| LAYOUT-002 | SHOULD | Dependencies and Project Layout/Project Layout | Put packages not meant to be exposed externally under internal/. |
| LAYOUT-003 | SHOULD | Dependencies and Project Layout/Project Layout | Keep executables together under cmd/, with one subdirectory per program. |
| LAYOUT-004 | MAY | Dependencies and Project Layout/Project Layout | Do not lay out a large project's directory skeleton in advance while the project is still small. |
| LOG-001 | MUST | Errors and Logging/Logging | Record business logs with log/slog; do not print directly with fmt or log. |
| LOG-002 | SHOULD | Errors and Logging/Logging | Write log fields as key-value pairs; do not splice them into the message text. |
| LOG-003 | SHOULD | Errors and Logging/Logging | Use the Context-carrying methods for logs on the request path. |
| LOG-004 | MUST | Errors and Logging/Logging | Logs never emit passwords, tokens, emails, or other PII and credentials. |
| MODERN-001 | MAY | Programming Conventions/Data Structures | Common operations on slices and maps prefer the slices and maps standard library packages over hand-written loops. |
| NAMING-001 | MUST | Programming Conventions/Naming | Package names use lowercase words written together, without underscores, camel case, or plurals. |
| NAMING-002 | MUST | Programming Conventions/Naming | Identifiers do not use underscores to separate words. |
| NAMING-003 | MUST | Programming Conventions/Naming | Initialisms keep a consistent case within identifiers. |
| NAMING-004 | SHOULD | Programming Conventions/Naming | The length of a name is proportional to the size of its scope. |
| NAMING-005 | SHOULD | Programming Conventions/Naming | Names do not repeat the type or the surrounding context. |
| NAMING-006 | SHOULD | Programming Conventions/Naming | Receiver names use a short abbreviation of the type name and stay consistent within the same type. |
| NAMING-007 | SHOULD | Programming Conventions/Naming | Boolean predicate methods are named with an Is, Has, or Can prefix. |
| NAMING-008 | SHOULD | Programming Conventions/Naming | Boolean variables and struct fields are named with an Is, Has, or Can prefix. |
| NAMING-009 | SHOULD | Programming Conventions/Naming | Single-method interfaces are named after the method plus an er suffix. |
| NAMING-010 | SHOULD | Programming Conventions/Naming | Function variants use fixed markers such as Must, f, and WithContext. |
| NAMING-011 | SHOULD | Programming Conventions/Naming | Identifiers do not shadow predeclared names such as len, cap, error, and min. |
| PERF-001 | SHOULD | Performance | Before optimizing, locate the hotspot with a benchmark or profile; do not change code on intuition. |
| PERF-002 | SHOULD | Performance | Use strconv, not fmt, to convert between primitive types and strings. |
| PERF-003 | SHOULD | Performance | Convert a fixed string to []byte once outside the loop and reuse the result. |
| PERF-004 | MUST | Performance | When the number of elements is known or can be estimated, specify the slice capacity with make. |
| PERF-005 | MUST | Performance | When the number of elements is known or can be estimated, pass a capacity hint to make(map). |
| PERF-006 | SHOULD | Performance | Use strings.Builder to concatenate strings in a loop or across many operations, and call Grow first when the capacity can be estimated. |
| PERF-007 | SHOULD | Performance | When keeping a small section of a large slice for a long time, use slices.Clone to cut the reference to the backing array. |
| PERF-008 | SHOULD | Performance | Pass arguments by value when the function only reads them through a dereference; do not pass pointers just to save a few bytes. |
| PERF-009 | SHOULD | Performance | Production code does not compare with reflect.DeepEqual. |
| PERF-010 | SHOULD | Performance | regexp and template compile once at package level. |
| PERF-011 | SHOULD | Performance | Unbounded inputs are processed streaming, not loaded whole with io.ReadAll. |
| PERF-012 | SHOULD | Performance | A reused http.Client sets Timeout explicitly and customizes its Transport. |
| PERF-013 | SHOULD | Performance | Handle resp.Body by version: Go 1.27+ reclaims on Close; earlier versions must drain it. |
| PERF-014 | SHOULD | Performance | Range over maps taking the value variable, not a second subscript lookup. |
| PERF-016 | SHOULD | Performance | sql.DB connection pool limits and lifetimes are configured explicitly. |
| PKG-001 | MUST | Dependencies and Project Layout/Package Organization | The package name is the last segment of the import path and matches the directory it lives in. |
| PKG-002 | SHOULD | Dependencies and Project Layout/Package Organization | Do not create packages with no informational content, such as util, common, or helper. |
| PKG-003 | MUST | Dependencies and Project Layout/Package Organization | Arrange imports in groups, one for the standard library and one for everything else, separated by a blank line. |
| PKG-004 | MUST | Dependencies and Project Layout/Package Organization | Do not use dot imports. |
| PKG-005 | SHOULD | Dependencies and Project Layout/Package Organization | Exported names do not repeat the package name, and constructors prefer New. |
| PKG-006 | SHOULD | Dependencies and Project Layout/Package Organization | Blank imports appear only in the main package or in test files. |
| PKG-007 | SHOULD | Dependencies and Project Layout/Package Organization | Give an import alias only when package names collide or the name does not match the last path segment. |
| SEC-001 | MUST | Security | SQL statements pass parameters through placeholders; do not concatenate parameter values into the statement. |
| SEC-002 | MUST | Security | Pass the executable and arguments of external commands to exec separately; do not build a command string. |
| SEC-003 | MUST | Security | User-controllable file paths are confined to an allowed directory. |
| SEC-004 | MUST | Security | Random numbers for security purposes come from crypto/rand, not math/rand. |
| SEC-005 | SHOULD | Security | Do not set InsecureSkipVerify to skip TLS certificate verification. |
| SEC-006 | MUST | Security | Do not use broken algorithms such as MD5, SHA-1, DES, and RC4 for security purposes. |
| SEC-007 | MUST | Security | Do not hard-code keys, passwords, or tokens in source code. |
| SEC-008 | SHOULD | Security | Grant only the required permissions when creating files and directories; do not use globally writable values. |
| SEC-009 | SHOULD | Security | Use html/template, not text/template, to output HTML. |
| SEC-010 | SHOULD | Security | HTTP services explicitly set the read header, read, and write timeouts. |
| SEC-011 | SHOULD | Security | Compare passwords, tokens, and MACs with constant-time functions, not comparisons that return early. |
| SEC-012 | MUST | Security | User-controllable SQL identifiers are validated against a whitelist before splicing. |
| SEC-013 | MUST | Security | Untrusted input is size-limited; request bodies go through http.MaxBytesReader. |
| SEC-014 | MUST | Security | Temporary files are created with os.CreateTemp. |
| SEC-015 | MUST | Security | tls.Config sets MinVersion explicitly to TLS 1.2 or above. |
| SEC-016 | MUST | Security | SSH connections verify host keys; InsecureIgnoreHostKey is forbidden. |
| SEC-017 | MUST | Security | Password storage uses slow hashes such as Argon2id and bcrypt. |
| SEC-018 | MUST | Security | Symmetric encryption uses authenticated modes; nonces are random and never reused. |
| SEC-019 | MUST | Security | Session cookies set HttpOnly, Secure, and SameSite. |
| SEC-020 | MUST | Security | JWT verification pins an algorithm whitelist; the token's alg field is not trusted. |
| SEC-021 | MUST | Security | Application code does not use unsafe to bypass memory safety. |
| SEC-022 | MUST | Security | Authorization uses the server-verified identity, never client request headers. |
| STD-001 | MUST | Programming Conventions/Standard Library Usage | time.Duration is never multiplied by time.Duration. |
| STD-002 | SHOULD | Programming Conventions/Standard Library Usage | Compare time.Time values with Equal, not ==. |
| STD-003 | SHOULD | Programming Conventions/Standard Library Usage | Strip prefixes and suffixes with strings.TrimPrefix and TrimSuffix. |
| STRUCT-001 | MUST | Interfaces and Design/Safe Struct Shape | Structs containing sync.Mutex or other synchronization primitives must not be copied by value; always pass or receive them by pointer. |
| STRUCT-002 | MAY | Interfaces and Design/Safe Struct Shape | For structs that should not support comparison, use an incomparable field to block == at compile time. |
| STRUCT-003 | SHOULD | Interfaces and Design/Safe Struct Shape | Methods of one type keep a consistent receiver type. |
| STRUCT-004 | SHOULD | Interfaces and Design/Safe Struct Shape | Exported fields participating in serialization carry field tags. |
| STYLE-001 | SHOULD | Programming Conventions/Formatting and Style | Use any instead of interface{}. |
| TEST-001 | SHOULD | Unit Testing/Test Naming and Structure | Pure function tests use the native testing package, and component tests use the testify suite. |
| TEST-002 | SHOULD | Unit Testing/Test Naming and Structure | Unit tests and integration tests are stored separately. |
| TEST-003 | SHOULD | Unit Testing/Assertions and Test Data | Compare the complete result in one shot; do not break it apart field by field or trim it into a subset. |
| TEST-004 | SHOULD | Unit Testing/Assertions and Test Data | Bulky test inputs and expected data go into the testdata directory. |
| TEST-005 | MUST | Unit Testing/Test Naming and Structure | Test files use external package tests, with the package name ending in _test. |
| TEST-006 | SHOULD | Unit Testing/Using Coverage | Test business logic and boundaries, and do not chase coverage numbers. |
| TEST-007 | SHOULD | Unit Testing/Assertions and Test Data | Assert on an error only by whether it is non-nil; do not compare error message text. |
| TEST-008 | MAY | Unit Testing/Test Naming and Structure | Test case names preferably follow the should-plus-outcome, when-plus-condition pattern. |
| TEST-010 | SHOULD | Unit Testing/Table-Driven and Case Design | Only cases that share the same check logic go into the same table. |
| TEST-011 | SHOULD | Unit Testing/Table-Driven and Case Design | The table holds only the fields that vary per case. |
| TEST-012 | MUST | Unit Testing/Table-Driven and Case Design | Write expected values independently; do not produce them with the function under test. |
| TEST-013 | SHOULD | Unit Testing/Table-Driven and Case Design | Create separate cases for boundaries and special inputs. |
| TEST-014 | SHOULD | Unit Testing/Table-Driven and Case Design | A defect fix comes with a reproducing case. |
| TEST-015 | MUST | Unit Testing/Mocks and Test Doubles | Mock generation uses go.uber.org/mock. |
| TEST-016 | MUST | Unit Testing/Mocks and Test Doubles | Do not call ctrl.Finish() after passing *testing.T. |
| TEST-017 | MUST | Unit Testing/Mocks and Test Doubles | Do not write redundant .Times(1). |
| TEST-018 | SHOULD | Unit Testing/Mocks and Test Doubles | Use gomock.InOrder() only when call order is genuinely required. |
| TEST-019 | SHOULD | Unit Testing/Mocks and Test Doubles | Write argument values directly in expectations, avoiding wrapping in Eq and bare Any(). |
| TEST-020 | SHOULD | Unit Testing/Table-Driven and Case Design | Table-driven cases run through t.Run as named subtests. |
| TEST-021 | SHOULD | Unit Testing/Test Naming and Structure | Tests do not depend on execution order; every case runs on its own. |
| TEST-022 | SHOULD | Unit Testing/Test Naming and Structure | Tests wait on synchronization primitives, never on time.Sleep. |
| TEST-023 | SHOULD | Unit Testing/Test Naming and Structure | Test files are named after the source file under test. |
| TEST-024 | SHOULD | Unit Testing/Test Naming and Structure | Test files are written into t.TempDir. |
| TEST-025 | SHOULD | Unit Testing/Test Naming and Structure | Tests take their ctx from t.Context. |
| TEST-026 | SHOULD | Unit Testing/Test Naming and Structure | Benchmarks use b.Loop, not hand-written b.N loops. |
| TEST-027 | SHOULD | Unit Testing/Test Naming and Structure | Benchmark comparisons use benchstat with enough samples for significance. |
| TEST-028 | MAY | Unit Testing/Test Naming and Structure | Benchmark functions live in a separate _bench_test.go file. |
| TEST-029 | SHOULD | Unit Testing/Assertions and Test Data | require guards preconditions; assert collects result checks. |
| TEST-030 | MAY | Unit Testing/Test Naming and Structure | Packages with goroutines use goleak to check for leaks. |
| TOOL-001 | MUST | Static Analysis/Configuration and Running | Use golangci-lint uniformly for static code analysis. |
| TOOL-002 | SHOULD | Static Analysis/Recommended Rules | Enable modernize in golangci-lint to keep up with modern Go idioms. |
| TOOL-003 | SHOULD | Static Analysis/Recommended Rules | Enable testifylint alongside testify when testify is used. |
| TOOL-004 | SHOULD | Tooling/Command Line | Command line programs organize subcommands with `cobra`. |
| TOOL-005 | MUST | Tooling/Testing Tools | Unit test assertions uniformly use `testify`. |
| TOOL-006 | MUST | Static Analysis/Required Rules | Enable at least the standard set of golangci-lint. |
| TOOL-007 | SHOULD | Static Analysis/Recommended Rules | Beyond the standard set, enable the commonly used supplementary analyzers. |
| TOOL-008 | SHOULD | Tooling/Command Line | After upgrading the Go toolchain, apply modernization rewrites with `go fix`. |
| TOOL-009 | MUST | Static Analysis/Configuration and Running | Name the checker and state the reason when suppressing a warning. |
| TOOL-010 | MUST | Tooling/Code Generation | Generated files are not edited by hand; they are produced only by the generator. |
| TOOL-011 | MUST | Tooling/Code Generation | Generated files are named with a gen prefix or suffix, preferring the suffix. |
| TOOL-012 | MUST | Tooling/Command Line | cobra programs write output through cmd.OutOrStdout and ErrOrStderr. |
| TOOL-013 | SHOULD | Tooling/Testing Tools | Tests and CI run with go test -race. |
| TOOL-014 | SHOULD | Tooling/Command Line | govulncheck runs as part of routine dependency vulnerability scanning. |
| TOOL-015 | SHOULD | Tooling/Command Line | Signal handling uses signal.NotifyContext into the context cancellation chain. |

## 2 Sources

- ABSTRACT-001: https://go.dev/wiki/CodeReviewComments#interfaces
- CODE-001: https://cloud.google.com/apis/design/errors
- CODE-002: https://cloud.google.com/apis/design/errors
- CODE-003: https://cloud.google.com/apis/design/errors
- CODE-004: https://go.dev/wiki/CodeReviewComments#error-strings
- COMMENT-001: https://google.github.io/styleguide/go/guide#clarity-rationale
- COMMENT-002: https://go.dev/wiki/CodeReviewComments#comment-sentences https://go.dev/doc/effective_go#commentary
- COMPOSE-001: https://go.dev/doc/faq#inheritance https://go.dev/doc/effective_go#embedding
- COMPOSE-002: https://github.com/uber-go/guide/blob/master/style.md#embedding-in-structs https://github.com/uber-go/guide/blob/master/style.md#avoid-embedding-types-in-public-structs
- CONC-001: https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes https://pkg.go.dev/golang.org/x/sync/errgroup
- CONC-002: https://go.dev/blog/pipelines
- CONC-003: https://go.dev/wiki/CodeReviewComments#synchronous-functions
- CONC-004: https://go.dev/wiki/MutexOrChannel https://pkg.go.dev/sync
- CONC-005: https://go.dev/ref/spec#Channel_types
- CONC-006: https://github.com/uber-go/guide/blob/master/style.md#channel-size
- CONC-007: https://pkg.go.dev/context#Context.Done
- CONC-008: https://go.dev/ref/spec#Select_statements
- CONC-009: https://pkg.go.dev/time#Tick
- CONC-010: https://pkg.go.dev/sync#WaitGroup.Go https://go.dev/doc/go1.25
- CONC-011: https://go.dev/ref/mem
- CONC-012: https://pkg.go.dev/sync#RWMutex
- CONC-013: https://pkg.go.dev/sync/atomic#Int64
- CONC-014: https://pkg.go.dev/sync#Map
- CONC-015: https://pkg.go.dev/golang.org/x/sync/errgroup#Group.SetLimit
- CONC-016: https://pkg.go.dev/sync#Pool
- CONC-017: https://pkg.go.dev/net/http#Server
- CONST-001: https://go.dev/wiki/CodeReviewComments#mixed-caps https://google.github.io/styleguide/go/decisions#constant-names
- CONST-002: https://google.github.io/styleguide/go/decisions#constant-names
- CONST-003: https://go.dev/doc/effective_go#constants https://go.dev/ref/spec#Constant_declarations
- CONST-004: https://go.dev/ref/spec#Iota
- CTRL-001: https://go.dev/wiki/CodeReviewComments#indent-error-flow
- CTRL-002: https://go.dev/doc/go1.22 https://go.dev/blog/loopvar-preview
- CTRL-003: https://go.dev/doc/effective_go#defer
- CTRL-004: https://pkg.go.dev/net/http#ResponseWriter
- CTRL-005: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow
- CTX-001: https://go.dev/wiki/CodeReviewComments#contexts https://pkg.go.dev/context
- CTX-002: https://pkg.go.dev/context#pkg-overview https://staticcheck.dev/docs/checks/#SA1012
- CTX-003: https://pkg.go.dev/context#Context.Value https://staticcheck.dev/docs/checks/#SA1029
- CTX-004: https://pkg.go.dev/context#pkg-overview
- CTX-005: https://pkg.go.dev/context#WithCancel
- CTX-006: https://pkg.go.dev/context#Background
- CTX-007: https://pkg.go.dev/context
- CTX-008: https://pkg.go.dev/context#WithoutCancel
- CTX-009: https://pkg.go.dev/context
- DATA-001: https://go.dev/wiki/CodeReviewComments#declaring-empty-slices
- DATA-002: https://pkg.go.dev/cmd/vet https://github.com/uber-go/guide/blob/master/style.md#use-field-names-to-initialize-structs
- DATA-003: https://go.dev/ref/spec#Type_assertions
- DATA-004: https://go.dev/blog/maps
- DATA-005: https://go.dev/ref/spec#For_range
- DATA-006: https://go.dev/ref/spec#Conversions_to_and_from_integer_string_types
- DATA-007: https://pkg.go.dev/encoding/json#Unmarshal
- DB-001: https://pkg.go.dev/database/sql#Rows.Err
- DB-002: https://pkg.go.dev/database/sql#Tx
- DB-003: https://pkg.go.dev/database/sql#NullString
- DB-004: https://pkg.go.dev/database/sql#ErrNoRows
- DEP-001: https://go.dev/doc/modules/managing-dependencies#synchronizing https://go.dev/doc/modules/gomod-ref
- DEP-002: https://go.dev/doc/modules/managing-dependencies#enable_tracking
- DEP-003: https://go.dev/doc/modules/managing-dependencies#local_directory https://go.dev/doc/modules/gomod-ref#replace
- DEP-004: https://go.dev/doc/modules/gomod-ref#tool https://go.dev/doc/modules/managing-dependencies#tools
- DEP-005: https://go-proverbs.github.io/
- DI-001: https://github.com/google/wire
- DI-002: https://github.com/uber-go/guide/blob/master/style.md#avoid-mutable-globals
- ERR-001: https://go.dev/blog/go1.13-errors https://google.github.io/styleguide/go/decisions
- ERR-003: https://go.dev/blog/go1.13-errors https://pkg.go.dev/errors#Is
- ERR-004: https://go.dev/wiki/CodeReviewComments#dont-panic https://go.dev/doc/effective_go#errors
- ERR-005: https://go.dev/wiki/CodeReviewComments#in-band-errors
- ERR-006: https://dave.cheney.net/practical-go/presentations/qcon-china.html#_only_handle_an_error_once https://go.dev/wiki/CodeReviewComments#handle-errors
- ERR-007: https://pkg.go.dev/builtin#error
- ERR-008: https://pkg.go.dev/errors#Join
- ERR-009: https://pkg.go.dev/os#Exit
- ERR-010: https://pkg.go.dev/bufio#Writer.Flush
- ERR-011: https://docs.bearer.com/reference/rules/go_lang_information_leakage/
- FUNC-001: https://github.com/uber-go/guide/blob/master/style.md#function-grouping-and-ordering https://github.com/manuelarte/funcorder
- FUNC-002: https://go.dev/wiki/CodeReviewComments#named-result-parameters
- FUNC-003: https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md https://go-critic.com/overview.html#toomanyresultschecker
- FUNC-004: https://go.dev/doc/faq#nil_error
- IFACE-001: https://go.dev/wiki/CodeReviewComments#interfaces
- IFACE-002: https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
- IFACE-003: https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces
- IFACE-004: https://go-proverbs.github.io/
- IFACE-005: https://pkg.go.dev/fmt#Stringer
- IFACE-006: https://go.dev/doc/comment
- LAYER-001: https://alistair.cockburn.us/hexagonal-architecture/
- LAYOUT-001: https://github.com/golang-standards/project-layout#directories-you-shouldnt-have https://go.dev/doc/modules/layout
- LAYOUT-002: https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages https://go.dev/doc/go1.4#internalpackages
- LAYOUT-003: https://go.dev/doc/modules/layout#multiple-commands https://github.com/golang-standards/project-layout#cmd
- LAYOUT-004: https://github.com/golang-standards/project-layout#overview
- LOG-001: https://go.dev/blog/slog https://pkg.go.dev/log/slog
- LOG-002: https://go.dev/blog/slog
- LOG-003: https://pkg.go.dev/log/slog#Logger.InfoContext https://go.dev/blog/slog
- LOG-004: https://docs.bearer.com/reference/rules/go_lang_logger/
- MODERN-001: https://pkg.go.dev/slices https://go.dev/doc/go1.21
- NAMING-001: https://go.dev/wiki/CodeReviewComments#package-names https://google.github.io/styleguide/go/decisions
- NAMING-002: https://google.github.io/styleguide/go/decisions#underscores https://go.dev/wiki/CodeReviewComments#mixed-caps
- NAMING-003: https://go.dev/wiki/CodeReviewComments#initialisms https://google.github.io/styleguide/go/decisions#initialisms
- NAMING-004: https://google.github.io/styleguide/go/decisions#variable-names https://go.dev/wiki/CodeReviewComments#variable-names
- NAMING-005: https://google.github.io/styleguide/go/decisions#variable-name-vs-type https://google.github.io/styleguide/go/decisions#external-context-vs-local-names
- NAMING-006: https://google.github.io/styleguide/go/decisions#receiver-names https://go.dev/wiki/CodeReviewComments#receiver-names
- NAMING-007: https://go.dev/doc/effective_go#Getters
- NAMING-008: https://go.dev/doc/effective_go#Getters
- NAMING-009: https://go.dev/wiki/CodeReviewComments#interface-names
- NAMING-010: https://pkg.go.dev/regexp#MustCompile
- NAMING-011: https://go.dev/ref/spec#Predeclared_identifiers
- PERF-001: https://go.dev/wiki/Performance
- PERF-002: https://github.com/uber-go/guide/blob/master/style.md#prefer-strconv-over-fmt
- PERF-003: https://github.com/uber-go/guide/blob/master/style.md#avoid-repeated-string-to-byte-conversions
- PERF-004: https://go.dev/blog/allocation-optimizations https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity
- PERF-005: https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity
- PERF-006: https://pkg.go.dev/strings#Builder
- PERF-007: https://go.dev/blog/slices-intro#a-possible-gotcha
- PERF-008: https://go.dev/wiki/CodeReviewComments#pass-values
- PERF-009: https://pkg.go.dev/slices#Equal
- PERF-010: https://pkg.go.dev/regexp#MustCompile
- PERF-011: https://pkg.go.dev/io#Copy
- PERF-012: https://pkg.go.dev/net/http#Client
- PERF-013: https://go.dev/doc/go1.27 https://pkg.go.dev/net/http#Response.Body
- PERF-014: https://go.dev/ref/spec#For_range
- PERF-016: https://pkg.go.dev/database/sql#DB
- PKG-001: https://go.dev/blog/package-names https://go.dev/wiki/CodeReviewComments#package-names
- PKG-002: https://google.github.io/styleguide/go/decisions#package-names https://go.dev/blog/package-names
- PKG-003: https://go.dev/wiki/CodeReviewComments#imports
- PKG-004: https://go.dev/wiki/CodeReviewComments#import-dot
- PKG-005: https://go.dev/blog/package-names https://google.github.io/styleguide/go/decisions#package-names
- PKG-006: https://go.dev/wiki/CodeReviewComments#import-blank
- PKG-007: https://go.dev/wiki/CodeReviewComments#imports
- SEC-001: https://go.dev/doc/database/sql-injection https://go.dev/doc/database/prepared-statements
- SEC-002: https://pkg.go.dev/os/exec
- SEC-003: https://go.dev/doc/go1.24#directory-limited-filesystem-access https://pkg.go.dev/os#Root
- SEC-004: https://pkg.go.dev/math/rand https://pkg.go.dev/crypto/rand
- SEC-005: https://pkg.go.dev/crypto/tls#Config
- SEC-006: https://pkg.go.dev/crypto/md5 https://pkg.go.dev/crypto/sha1 https://pkg.go.dev/crypto/des https://pkg.go.dev/crypto/rc4
- SEC-007: https://github.com/securego/gosec/blob/master/RULES.md
- SEC-008: https://github.com/securego/gosec/blob/master/RULES.md
- SEC-009: https://pkg.go.dev/html/template
- SEC-010: https://pkg.go.dev/net/http#Server.ReadHeaderTimeout https://github.com/securego/gosec/blob/master/RULES.md
- SEC-011: https://pkg.go.dev/crypto/subtle#ConstantTimeCompare https://pkg.go.dev/crypto/hmac#Equal
- SEC-012: https://go.dev/doc/database/sql-injection
- SEC-013: https://pkg.go.dev/net/http#MaxBytesReader
- SEC-014: https://docs.bearer.com/reference/rules/go_gosec_filesystem_tempfile/ https://pkg.go.dev/os#CreateTemp
- SEC-015: https://docs.bearer.com/reference/rules/go_lang_missing_tls_minversion/ https://pkg.go.dev/crypto/tls#Config
- SEC-016: https://docs.bearer.com/reference/rules/go_gosec_crypto_insecure_ignore_host_key/
- SEC-017: https://docs.bearer.com/reference/rules/go_lang_weak_password_encryption_md5/ https://pkg.go.dev/golang.org/x/crypto/bcrypt
- SEC-018: https://pkg.go.dev/crypto/cipher#AEAD
- SEC-019: https://docs.bearer.com/reference/rules/go_lang_cookie_missing_http_only/ https://docs.bearer.com/reference/rules/go_lang_insecure_cookie/
- SEC-020: https://datatracker.ietf.org/doc/html/rfc8725#section-3.1
- SEC-021: https://pkg.go.dev/unsafe https://docs.bearer.com/reference/rules/go_gosec_unsafe_unsafe/
- SEC-022: https://www.rfc-editor.org/rfc/rfc7239.html#section-8.1
- STD-001: https://pkg.go.dev/time#Duration
- STD-002: https://pkg.go.dev/time#Time.Equal
- STD-003: https://pkg.go.dev/strings#Trim
- STRUCT-001: https://pkg.go.dev/cmd/vet https://google.github.io/styleguide/go/decisions
- STRUCT-002: https://github.com/protocolbuffers/protobuf-go/blob/master/internal/pragma/pragma.go https://go.dev/ref/spec#Comparison_operators
- STRUCT-003: https://go.dev/wiki/CodeReviewComments#receiver-type
- STRUCT-004: https://pkg.go.dev/encoding/json#Marshal
- STYLE-001: https://go.dev/doc/go1.18 https://google.github.io/styleguide/go/decisions
- TEST-001: https://github.com/stretchr/testify#suite-package https://pkg.go.dev/github.com/stretchr/testify/suite
- TEST-002: https://pkg.go.dev/cmd/go#hdr-Test_packages
- TEST-003: https://go.dev/wiki/TestComments#compare-full-structures
- TEST-004: https://pkg.go.dev/embed https://pkg.go.dev/cmd/go#hdr-Test_packages
- TEST-005: https://pkg.go.dev/cmd/go#hdr-Test_packages
- TEST-006: https://research.swtch.com/testing
- TEST-007: https://go.dev/wiki/TestComments#test-error-semantics https://go.dev/blog/go1.13-errors
- TEST-008: https://go.dev/wiki/TestComments#choose-human-readable-subtest-names
- TEST-010: https://go.dev/wiki/TestComments#table-driven-tests-vs-multiple-test-functions https://research.swtch.com/testing
- TEST-011: https://research.swtch.com/testing https://go.dev/wiki/TableDrivenTests
- TEST-012: https://go.dev/wiki/TestComments#compare-full-structures https://research.swtch.com/testing
- TEST-013: https://research.swtch.com/testing
- TEST-014: https://research.swtch.com/testing
- TEST-015: https://github.com/uber-go/mock
- TEST-016: https://pkg.go.dev/go.uber.org/mock/gomock#Controller.Finish
- TEST-017: https://pkg.go.dev/go.uber.org/mock/gomock#Call.Times https://github.com/uber-go/mock
- TEST-018: https://pkg.go.dev/go.uber.org/mock/gomock#InOrder
- TEST-019: https://pkg.go.dev/go.uber.org/mock/gomock#Matcher https://github.com/uber-go/mock/blob/main/gomock/call.go
- TEST-020: https://pkg.go.dev/testing#T.Run
- TEST-021: https://pkg.go.dev/cmd/go
- TEST-022: https://pkg.go.dev/time#Sleep
- TEST-023: https://pkg.go.dev/testing
- TEST-024: https://pkg.go.dev/testing#T.TempDir
- TEST-025: https://pkg.go.dev/testing#T.Context
- TEST-026: https://pkg.go.dev/testing#B.Loop
- TEST-027: https://pkg.go.dev/golang.org/x/perf/cmd/benchstat
- TEST-028: https://pkg.go.dev/testing
- TEST-029: https://pkg.go.dev/github.com/stretchr/testify/require
- TEST-030: https://pkg.go.dev/go.uber.org/goleak
- TOOL-001: https://golangci-lint.run/
- TOOL-002: https://golangci-lint.run/docs/linters/configuration/#modernize
- TOOL-003: https://github.com/Antonboom/testifylint
- TOOL-004: https://github.com/spf13/cobra
- TOOL-005: https://github.com/stretchr/testify#assert-package
- TOOL-006: https://golangci-lint.run/docs/welcome/quick-start/
- TOOL-007: https://golangci-lint.run/docs/linters/
- TOOL-008: https://pkg.go.dev/cmd/fix https://go.dev/blog/gofix https://go.dev/doc/go1.26
- TOOL-009: https://golangci-lint.run/docs/linters/false-positives/#nolint-directive https://golangci-lint.run/docs/linters/configuration/#nolintlint
- TOOL-010: https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source
- TOOL-011: https://protobuf.dev/reference/go/go-generated/ https://pkg.go.dev/cmd/go#hdr-Build_constraints
- TOOL-012: https://pkg.go.dev/github.com/spf13/cobra#Command.OutOrStdout
- TOOL-013: https://go.dev/doc/security/best-practices https://go.dev/blog/race-detector
- TOOL-014: https://go.dev/doc/security/best-practices https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
- TOOL-015: https://pkg.go.dev/os/signal#NotifyContext

## 3 Detection Map

| ID | Detection |
| --- | --- |
| ABSTRACT-001 | manual review (Verify whether the interface or factory already has a second implementation, a replacement, or a testing need.) |
| CODE-001 | manual review (Check whether the external error codes place self-made numbers alongside the standard classification) |
| CODE-002 | manual review (Check whether there is a branch that tests by error message text) |
| CODE-003 | manual review (Check whether an error returned to the outside comes directly from a downstream service) |
| CODE-004 | manual review (Check the message text of errors.New and Errorf) |
| COMMENT-001 | manual review (Check whether a comment restates what the code itself already shows) |
| COMMENT-002 | manual review (Check whether a doc comment begins with the name and forms a full sentence) |
| COMPOSE-001 | manual review (Verify whether embedded types play the role of a base class (Base / Abstract naming, multi-level embedding).) |
| COMPOSE-002 | manual review (Verify whether embedding changes the outer zero value, copy semantics, or exposed unrelated methods (especially for synchronization primitives).) |
| CONC-001 | manual review (Check that each go statement has a matching waiter and exit path) |
| CONC-002 | manual review (Check that every close is on the sender side) |
| CONC-003 | manual review (Check whether a public function returns a channel that leaves the caller waiting) |
| CONC-004 | manual review (Check whether the lock-protected field is exported and whether it is reachable only through methods) |
| CONC-005 | manual review (Check the direction of channel parameters in exported functions) |
| CONC-006 | grep regex make\(chan[^)]*,\s*([^)1][^)]*)\) (Grep for make(chan, n) with n greater than 1 and no adjacent comment) |
| CONC-007 | manual review (Check the case set of long-lived goroutine selects) |
| CONC-008 | manual review (Check the default branch of selects inside loops) |
| CONC-009 | manual review (Check for time.After and time.Tick inside loop bodies) |
| CONC-010 | manual review (Check the relative position of wg.Add and the go statement) |
| CONC-011 | manual review (Check for blocking calls between Lock and Unlock) |
| CONC-012 | manual review (Check for nested RLock and Lock holds) |
| CONC-013 | manual review (Check the form of atomic operations) |
| CONC-014 | manual review (Check the read/write ratio and key-set pattern of sync.Map usage) |
| CONC-015 | manual review (Check that batch task fan-out sets a limit) |
| CONC-016 | manual review (Check that Pool.Put sites Reset first) |
| CONC-017 | manual review (Check that long-lived goroutine entry points recover and log the stack) |
| CONST-001 | golangci-lint stylecheck(ST1003) |
| CONST-002 | manual review (Check whether the constant name merely restates the value) |
| CONST-003 | manual review (Check whether enumerations that increase from one source use hand-written numbers) |
| CONST-004 | manual review (Check whether the first iota member occupies the zero value) |
| CTRL-001 | manual review (Check whether the normal path is indented by an else) |
| CTRL-002 | golangci-lint copyloopvar |
| CTRL-003 | manual review (Grep for defer statements inside loop bodies) |
| CTRL-004 | manual review (Check that handlers return after http.Error/w.Write) |
| CTRL-005 | golangci-lint govet(shadow) (Check for := declarations that shadow an outer err) |
| CTX-001 | golangci-lint containedctx (The first-parameter convention requires manual verification.) |
| CTX-002 | golangci-lint staticcheck(SA1012) (SA1012 detects calls that pass a nil context to a function.) |
| CTX-003 | golangci-lint staticcheck(SA1029) (SA1029 detects the use of a built-in type as a context key.) |
| CTX-004 | manual review (Verify whether what is put in context.Value is request-scoped data.) |
| CTX-005 | golangci-lint govet(lostcancel) (Check that every cancel returned by With* is called) |
| CTX-006 | manual review (Check for context.Background and TODO calls outside entry points) |
| CTX-007 | golangci-lint noctx (Check for bare calls where a *Context variant exists) |
| CTX-008 | manual review (Check the ctx source of background tasks started in handlers) |
| CTX-009 | manual review (Check the exit conditions of retry loops) |
| DATA-001 | manual review (Outputting a JSON array [] is an exception) |
| DATA-002 | golangci-lint govet(composites) (govet covers only positional literals from another package, and the composites check must be enabled) |
| DATA-003 | golangci-lint forcetypeassert (Check the form of every type assertion) |
| DATA-004 | manual review (Check that map write paths operate on initialized maps) |
| DATA-005 | manual review (Check that logic after ranging a map does not depend on order) |
| DATA-006 | golangci-lint gosec(G115) (Check the range before narrowing conversions) |
| DATA-007 | manual review (Compare struct fields with the serialization contract and check unexported fields) |
| DB-001 | golangci-lint rowserrcheck, sqlclosecheck (Check that rows.Err and rows.Close appear after iteration) |
| DB-002 | manual review (Check that transactions defer Rollback right after Begin) |
| DB-003 | manual review (Compare DDL nullable columns with the Scan target types) |
| DB-004 | manual review (Check how ErrNoRows is matched) |
| DEP-001 | manual review (Run go mod tidy before committing, and confirm go.mod and go.sum have no uncommitted changes) |
| DEP-002 | manual review (Check that the repository contains go.sum and that .gitignore does not ignore it) |
| DEP-003 | manual review (Check that replace directives in go.mod do not point at local directories long term and that removal conditions are stated) |
| DEP-004 | manual review (Check that the command-line tools used by the repository are declared through go.mod's tool directive) |
| DEP-005 | manual review (Check whether a newly added dependency is a small capability the standard library already covers) |
| DI-001 | manual review (Verify that the repository has wire_gen.go and that port bindings are concentrated in the ProviderSet.) |
| DI-002 | manual review (Verify whether package-level variables are ever assigned and whether dependencies are passed in through constructors.) |
| ERR-001 | golangci-lint errorlint |
| ERR-003 | golangci-lint errorlint (Covers both the == comparison and the type assertion forms) |
| ERR-004 | manual review (Check whether panic is used for expected failures) |
| ERR-005 | manual review (Check whether -1, an empty value or a zero value is used to express failure) |
| ERR-006 | manual review (Check each if err != nil branch: writing a log and returning err in the same branch is a violation) |
| ERR-007 | golangci-lint errcheck (Check for error values discarded with _) |
| ERR-008 | manual review (Check how multi-failure paths are merged) |
| ERR-009 | manual review (Check for os.Exit and log.Fatal calls outside main) |
| ERR-010 | manual review (Check error handling at write-path resource closing) |
| ERR-011 | manual review (Check how external error responses are constructed) |
| FUNC-001 | golangci-lint funcorder |
| FUNC-002 | manual review (Check whether results are named when their meaning is not self-evident) |
| FUNC-003 | golangci-lint revive(argument-limit, function-result-limit) |
| FUNC-004 | manual review (Check the nil branches of interface-returning functions) |
| IFACE-001 | manual review (Verify whether the package holding the interface matches the number of implementations: for a single implementation, whether it shares the package with the implementation; for multiple implementations, whether it is in the consumer's package.) |
| IFACE-002 | manual review (Verify whether exported types that implement an interface per an API contract have a compile-time assertion.) |
| IFACE-003 | manual review (Verify whether function signatures and struct fields contain *interface types.) |
| IFACE-004 | manual review (Check the method count of interface definitions) |
| IFACE-005 | manual review (Check that method signatures matching standard interfaces are exact) |
| IFACE-006 | manual review (Check git history and comments before deleting exported symbols) |
| LAYER-001 | manual review (Check the direction of package imports; this can be locked down as a rule with golangci-lint's depguard) |
| LAYOUT-001 | manual review (Check that the directory containing go.mod is the repository root and that the source is not placed under src/) |
| LAYOUT-002 | manual review (Check whether packages used only by this module live under internal/) |
| LAYOUT-003 | manual review (Check that main packages are gathered under cmd/ and that directory names match the executable names) |
| LAYOUT-004 | manual review (Check that every repository directory has real content and a clear responsibility) |
| LOG-001 | grep regex \b(fmt\.Print|log\.Print) (Check by hand whether the match is a business log) |
| LOG-002 | manual review (Check whether log values are spliced into the message string) |
| LOG-003 | manual review (Check whether logs on the request path pass in ctx) |
| LOG-004 | manual review (Check log fields for sensitive information) |
| MODERN-001 | golangci-lint modernize (go fix also suggests the equivalent modern form) |
| NAMING-001 | golangci-lint stylecheck(ST1003) (Check the package name identifier) |
| NAMING-002 | golangci-lint stylecheck(ST1003) (Check the identifier naming style) |
| NAMING-003 | golangci-lint stylecheck(ST1003) |
| NAMING-004 | manual review (Decide by hand whether name length matches the scope) |
| NAMING-005 | manual review (Check whether the name repeats the type or the surrounding context) |
| NAMING-006 | manual review (Check whether receiver names are consistent within the same type and whether this / self is used) |
| NAMING-007 | manual review (Check that methods returning bool carry a predicate prefix) |
| NAMING-008 | manual review (Check that bool fields and variables carry a prefix) |
| NAMING-009 | manual review (Check that single-method interface type names correspond to their method) |
| NAMING-010 | manual review (Check that variant functions keep naming markers consistent with the standard library) |
| NAMING-011 | golangci-lint predeclared (Check that local identifiers do not shadow predeclared names) |
| PERF-001 | manual review (Check whether a performance change is accompanied by benchmark results or profile conclusions.) |
| PERF-002 | manual review (Check whether conversions between a single primitive type and a string use fmt.) |
| PERF-003 | manual review (Check whether the loop contains conversions between string and []byte for a fixed string.) |
| PERF-004 | manual review (Check whether the make call for the slice an append targets carries a capacity argument.) |
| PERF-005 | manual review (Check whether a map written entry by entry carries a capacity hint.) |
| PERF-006 | manual review (Check whether the loop contains string concatenation with +=.) |
| PERF-007 | manual review (Check whether a sub-slice cut from a large buffer and kept for a long time is copied.) |
| PERF-008 | manual review (Check whether a parameter that is only read through a dereference uses a pointer.) |
| PERF-009 | manual review (Check DeepEqual calls on production paths) |
| PERF-010 | manual review (Check for regexp.Compile and template.Parse inside function bodies) |
| PERF-011 | manual review (Check whether io.ReadAll inputs are bounded) |
| PERF-012 | manual review (Check the configuration of package-level http.Clients) |
| PERF-013 | manual review (Check body handling in HTTP client code) |
| PERF-014 | manual review (Check map range loop bodies) |
| PERF-016 | manual review (Check pool configuration after sql.Open) |
| PKG-001 | manual review (Check directory by directory that the package clause matches the directory name) |
| PKG-002 | grep regex ^package\s+(util|utility|common|misc|helper|helpers|model)\b (Split each match into named packages according to its real responsibility) |
| PKG-003 | golangci-lint gci (Enable gci in formatters, with group order standard, default, and the module prefix of this repository) |
| PKG-004 | grep regex ^\s*(import\s+)?\.\s+" (For each match, check whether it is a _test file in an import-cycle scenario) |
| PKG-005 | manual review (Check whether exported names repeat the package name they live in) |
| PKG-006 | grep regex ^\s*(import\s+)?_\s+" (For each match, check whether the containing file is in the main package or is a _test file) |
| PKG-007 | manual review (Check whether each alias resolves a name collision or a mismatch between package name and path) |
| SEC-001 | golangci-lint gosec(G201, G202) (At each hit, confirm whether the parameter value comes from external input.) |
| SEC-002 | golangci-lint gosec(G204) (At each hit, confirm whether external input was concatenated into the command string.) |
| SEC-003 | golangci-lint gosec(G304) (At each hit, confirm whether the path contains external input and whether it is confined within the root directory.) |
| SEC-004 | golangci-lint gosec(G404) (A math/rand call is hit; confirm whether it is used in a security scenario.) |
| SEC-005 | golangci-lint gosec(G402) (At each hit, confirm whether non-test code still skips verification.) |
| SEC-006 | golangci-lint gosec(G401, G405, G501, G502, G503, G505) |
| SEC-007 | golangci-lint gosec(G101) (At each hit, confirm whether it is a real credential; for a false positive, annotate the reason with #nosec G101.) |
| SEC-008 | golangci-lint gosec(G301, G302, G306) (The thresholds can be tightened as needed in the gosec configuration in .golangci.yml.) |
| SEC-009 | golangci-lint gosec(G203) (A hit on G203 writes unescaped data into an HTML template.) |
| SEC-010 | golangci-lint gosec(G112, G114) (G112 detects a missing ReadHeaderTimeout; G114 detects a serve function without timeout support.) |
| SEC-011 | manual review (Check whether secret comparison sites use a constant-time function.) |
| SEC-012 | manual review (Check the identifier sources in dynamic SQL) |
| SEC-013 | manual review (Check size limits on decompression and read paths) |
| SEC-014 | golangci-lint gosec(G302, G304) (Check temporary file creation sites) |
| SEC-015 | golangci-lint gosec(G402) (Check MinVersion on tls.Config literals) |
| SEC-016 | golangci-lint gosec(G106) (Check HostKeyCallback configurations) |
| SEC-017 | manual review (Check the hashing algorithm for password storage) |
| SEC-018 | manual review (Check the encryption mode and nonce generation) |
| SEC-019 | manual review (Check cookie attribute settings) |
| SEC-020 | manual review (Check JWT parser configurations) |
| SEC-021 | grep regex \bunsafe\. (Grep for unsafe package usage) |
| SEC-022 | manual review (Check where the auth middleware takes identity from) |
| STD-001 | golangci-lint durationcheck (Check every multiplication involving time.Duration) |
| STD-002 | manual review (Check the expressions used to compare time.Time values) |
| STD-003 | manual review (Check that the trimming intent matches the API used) |
| STRUCT-001 | golangci-lint govet(copylocks) |
| STRUCT-002 | manual review (Verify whether exported structs that intentionally prohibit comparison have embedded an incomparable field.) |
| STRUCT-003 | manual review (Check receiver declarations across methods of the same type) |
| STRUCT-004 | golangci-lint structtag (Check that serialized exported fields carry tags) |
| STYLE-001 | golangci-lint modernize(any) (The same analyzer can rewrite it directly with go fix) |
| TEST-001 | manual review (Check whether pure function cases are wrapped in a suite, and whether component cases assemble the object under test centrally with a suite.) |
| TEST-002 | manual review (Check whether tests that read real files, the network, or a database end with _it_test.go.) |
| TEST-003 | manual review (Check whether assertions are split field by field, or whether the result is trimmed into a subset before comparison.) |
| TEST-004 | manual review (Check whether cases read testdata by joining paths from the working directory.) |
| TEST-005 | golangci-lint testpackage |
| TEST-006 | manual review (Check whether cases assert and whether they target business semantics; coverage numbers are not the basis for judgment.) |
| TEST-007 | manual review (Check whether cases judge the kind of error by error message text, and whether sentinels and error types are checked with errors.Is / errors.As.) |
| TEST-008 | manual review (Check whether t.Run names and table name fields state the expected outcome and the triggering condition clearly.) |
| TEST-010 | manual review (Check whether a single loop body switches check logic with conditional branches.) |
| TEST-011 | manual review (Check whether table fields are written with the same value by every case.) |
| TEST-012 | manual review (Check whether expected values are produced by functions or methods inside the package under test.) |
| TEST-013 | manual review (Check whether boundary inputs such as empty collections, single elements, the first and last elements, and out-of-range access are covered.) |
| TEST-014 | manual review (Check whether a commit that fixes a defect carries a case matching the triggering conditions.) |
| TEST-015 | grep regex github\.com/golang/mock (Switch the matches to go.uber.org/mock; go.mod is out of scanning scope and needs manual checking.) |
| TEST-016 | grep regex \.Finish\(\) (Confirm at each match whether a *testing.T was passed in.) |
| TEST-017 | grep regex \.Times\(1\) (Delete it at each match; one call is the default.) |
| TEST-018 | grep regex gomock\.InOrder (Confirm at each match whether the call order affects the result.) |
| TEST-019 | grep regex gomock\.Any\(\) (Except for context.Context, switch the matches to Eq or AssignableToTypeOf.) |
| TEST-020 | manual review (Check that table-driven tests run through t.Run) |
| TEST-021 | manual review (Check whether cases depend on data left by other cases) |
| TEST-022 | manual review (Check the purpose of time.Sleep in tests) |
| TEST-023 | manual review (Check the naming mapping between test files and source files) |
| TEST-024 | manual review (Check temporary directory creation in tests) |
| TEST-025 | manual review (Check the ctx source inside tests) |
| TEST-026 | manual review (Check the loop form of benchmark functions) |
| TEST-027 | manual review (Check the sampling and judgment behind performance conclusions) |
| TEST-028 | manual review (Check which file benchmarks belong to) |
| TEST-029 | golangci-lint testifylint (Check the require/assert split in test bodies) |
| TEST-030 | manual review (Check that TestMain wires in the leak check together with its exemption list) |
| TOOL-001 | manual review (Check that the repository root has .golangci.yml and that the lint entry point invokes golangci-lint) |
| TOOL-002 | manual review (Check that linters.enable in .golangci.yml includes modernize) |
| TOOL-003 | manual review (Check that linters.enable in .golangci.yml includes testifylint) |
| TOOL-004 | manual review (Check whether subcommands are organized with `cobra`) |
| TOOL-005 | grep regex \bt\.(Errorf|Fatalf)\( (Matches should be replaced with `testify` assertions) |
| TOOL-006 | manual review (Check that .golangci.yml does not turn off the standard set) |
| TOOL-007 | manual review (Check whether the enable list covers the list of analyzers and whether disabled ones have a reason) |
| TOOL-008 | manual review (Check whether `go fix` was run after the Go version upgrade) |
| TOOL-009 | golangci-lint nolintlint (Configured as require-specific + require-explanation + allow-unused, blocking suppressions that do not name a checker, lack a reason, or are no longer effective) |
| TOOL-010 | manual review (Check whether changes land on the generator or on the input sources) |
| TOOL-011 | manual review (Check whether generated file names carry a gen prefix or suffix) |
| TOOL-012 | manual review (Check output writing in cobra commands) |
| TOOL-013 | manual review (Check CI configuration) |
| TOOL-014 | manual review (Check CI configuration) |
| TOOL-015 | manual review (Check signal handling in main) |
