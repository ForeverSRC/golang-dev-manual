# 10. Appendix

## 1 Clause Index

| ID | Level | Category | Summary |
| --- | --- | --- | --- |
| ABSTRACT-001 | SHOULD | Interfaces and Design/Interface Definition and Placement | Introduce abstraction according to actual need, not spread out in advance. |
| CODE-001 | SHOULD | Errors and Logging/Error Codes | For external error codes, reuse the existing standard classification; do not build a separate numbering scheme. |
| CODE-002 | SHOULD | Errors and Logging/Error Codes | Error messages are not a contract; test errors only by code and type. |
| CODE-003 | SHOULD | Errors and Logging/Error Codes | When propagating an error across services, translate the error code; do not pass through the downstream raw error. |
| COMMENT-001 | SHOULD | Programming Conventions/Comments | Comments record only what the code cannot express and never restate what the code does. |
| COMMENT-002 | SHOULD | Programming Conventions/Comments | A doc comment begins with the name being described and is written as a full sentence. |
| COMPOSE-001 | SHOULD | Interfaces and Design/Composition and Reuse | Express reusable capability through composition, not inheritance-style hierarchies. |
| COMPOSE-002 | SHOULD | Interfaces and Design/Composition and Reuse | Embedding must not change the outer type's zero value, copy semantics, or public interface. |
| CONC-001 | MUST | Programming Conventions/Concurrency | Before starting a goroutine, make its exit timing clear, and let the starter wait for it to finish. |
| CONC-002 | MUST | Programming Conventions/Concurrency | A channel is closed only by the sender, never by the receiver. |
| CONC-003 | SHOULD | Programming Conventions/Concurrency | A public API keeps synchronous semantics and leaves concurrency to the caller. |
| CONC-004 | SHOULD | Programming Conventions/Concurrency | Shared state and the lock protecting it live in the same struct and are read and written only through methods. |
| CONST-001 | MUST | Programming Conventions/Constants and Enums | Constant names use MixedCaps, not all uppercase or a k prefix. |
| CONST-002 | SHOULD | Programming Conventions/Constants and Enums | Constants are named by their role, not by their value. |
| CONST-003 | MAY | Programming Conventions/Constants and Enums | Enumeration values that increase from one source are generated with iota. |
| CTRL-001 | SHOULD | Programming Conventions/Control Flow | Handle errors and edge cases first and return early, keeping the normal path unindented. |
| CTRL-002 | SHOULD | Programming Conventions/Control Flow | Do not explicitly copy the loop variable any longer. |
| CTX-001 | SHOULD | Interfaces and Design/Passing context | Do not store context.Context in a struct; pass it as the first parameter of a function. |
| CTX-002 | MUST | Interfaces and Design/Passing context | Do not pass a nil context to a function; use context.TODO() as a placeholder. |
| CTX-003 | SHOULD | Interfaces and Design/Passing context | Use an unexported custom type for context.Value keys. |
| CTX-004 | SHOULD | Interfaces and Design/Passing context | Put only request-scoped data in context.Value, not optional parameters. |
| DATA-001 | SHOULD | Programming Conventions/Data Structures | Declare empty slices with var to get a nil slice. |
| DATA-002 | SHOULD | Programming Conventions/Data Structures | Struct literals are initialized with field names rather than listed by position. |
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
| FUNC-001 | MUST | Programming Conventions/Functions and Methods | Exported functions and methods come at the front of a file, unexported ones at the back. |
| FUNC-002 | SHOULD | Programming Conventions/Functions and Methods | Use named results when the meaning of the return values is not self-evident. |
| FUNC-003 | SHOULD | Programming Conventions/Functions and Methods | Collect more than 4 arguments or more than 3 results into a struct. |
| IFACE-001 | SHOULD | Interfaces and Design/Interface Definition and Placement | Place interfaces according to the number of implementations: a single implementation shares the package with the implementation, while multiple implementations are defined by the consumer. |
| IFACE-002 | SHOULD | Interfaces and Design/Interface Definition and Placement | Declare a type's implementation of an interface with a compile-time assertion. |
| IFACE-003 | SHOULD | Interfaces and Design/Interface Definition and Placement | Do not define pointers to interfaces. |
| LAYER-001 | SHOULD | Dependencies and Project Layout/Layering | Layer code by responsibility so that dependencies point one-way toward the inner layers. |
| LAYOUT-001 | MUST | Dependencies and Project Layout/Project Layout | Put go.mod and the source code at the module root, without an extra src/ layer. |
| LAYOUT-002 | SHOULD | Dependencies and Project Layout/Project Layout | Put packages not meant to be exposed externally under internal/. |
| LAYOUT-003 | SHOULD | Dependencies and Project Layout/Project Layout | Keep executables together under cmd/, with one subdirectory per program. |
| LAYOUT-004 | MAY | Dependencies and Project Layout/Project Layout | Do not lay out a large project's directory skeleton in advance while the project is still small. |
| LOG-001 | MUST | Errors and Logging/Logging | Record business logs with log/slog; do not print directly with fmt or log. |
| LOG-002 | SHOULD | Errors and Logging/Logging | Write log fields as key-value pairs; do not splice them into the message text. |
| LOG-003 | SHOULD | Errors and Logging/Logging | Use the Context-carrying methods for logs on the request path. |
| MODERN-001 | MAY | Programming Conventions/Data Structures | Common operations on slices and maps prefer the slices and maps standard library packages over hand-written loops. |
| NAMING-001 | MUST | Programming Conventions/Naming | Package names use lowercase words written together, without underscores, camel case, or plurals. |
| NAMING-002 | MUST | Programming Conventions/Naming | Identifiers do not use underscores to separate words. |
| NAMING-003 | MUST | Programming Conventions/Naming | Initialisms keep a consistent case within identifiers. |
| NAMING-004 | SHOULD | Programming Conventions/Naming | The length of a name is proportional to the size of its scope. |
| NAMING-005 | SHOULD | Programming Conventions/Naming | Names do not repeat the type or the surrounding context. |
| NAMING-006 | SHOULD | Programming Conventions/Naming | Receiver names use a short abbreviation of the type name and stay consistent within the same type. |
| PERF-001 | SHOULD | Performance | Before optimizing, locate the hotspot with a benchmark or profile; do not change code on intuition. |
| PERF-002 | SHOULD | Performance | Use strconv, not fmt, to convert between primitive types and strings. |
| PERF-003 | SHOULD | Performance | Convert a fixed string to []byte once outside the loop and reuse the result. |
| PERF-004 | MUST | Performance | When the number of elements is known or can be estimated, specify the slice capacity with make. |
| PERF-005 | MUST | Performance | When the number of elements is known or can be estimated, pass a capacity hint to make(map). |
| PERF-006 | SHOULD | Performance | Use strings.Builder to concatenate strings in a loop or across many operations, and call Grow first when the capacity can be estimated. |
| PERF-007 | SHOULD | Performance | When keeping a small section of a large slice for a long time, use slices.Clone to cut the reference to the backing array. |
| PERF-008 | SHOULD | Performance | Pass arguments by value when the function only reads them through a dereference; do not pass pointers just to save a few bytes. |
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
| STRUCT-001 | MUST | Interfaces and Design/Safe Struct Shape | Structs containing sync.Mutex or other synchronization primitives must not be copied by value; always pass or receive them by pointer. |
| STRUCT-002 | MAY | Interfaces and Design/Safe Struct Shape | For structs that should not support comparison, use an incomparable field to block == at compile time. |
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

## 2 Sources

- ABSTRACT-001: https://go.dev/wiki/CodeReviewComments#interfaces
- CODE-001: https://cloud.google.com/apis/design/errors
- CODE-002: https://cloud.google.com/apis/design/errors
- CODE-003: https://cloud.google.com/apis/design/errors
- COMMENT-001: https://google.github.io/styleguide/go/guide#clarity-rationale
- COMMENT-002: https://go.dev/wiki/CodeReviewComments#comment-sentences https://go.dev/doc/effective_go#commentary
- COMPOSE-001: https://go.dev/doc/faq#inheritance https://go.dev/doc/effective_go#embedding
- COMPOSE-002: https://github.com/uber-go/guide/blob/master/style.md#embedding-in-structs https://github.com/uber-go/guide/blob/master/style.md#avoid-embedding-types-in-public-structs
- CONC-001: https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes https://pkg.go.dev/golang.org/x/sync/errgroup
- CONC-002: https://go.dev/blog/pipelines
- CONC-003: https://go.dev/wiki/CodeReviewComments#synchronous-functions
- CONC-004: https://go.dev/wiki/MutexOrChannel https://pkg.go.dev/sync
- CONST-001: https://go.dev/wiki/CodeReviewComments#mixed-caps https://google.github.io/styleguide/go/decisions#constant-names
- CONST-002: https://google.github.io/styleguide/go/decisions#constant-names
- CONST-003: https://go.dev/doc/effective_go#constants https://go.dev/ref/spec#Constant_declarations
- CTRL-001: https://go.dev/wiki/CodeReviewComments#indent-error-flow
- CTRL-002: https://go.dev/doc/go1.22 https://go.dev/blog/loopvar-preview
- CTX-001: https://go.dev/wiki/CodeReviewComments#contexts https://pkg.go.dev/context
- CTX-002: https://pkg.go.dev/context#pkg-overview https://staticcheck.dev/docs/checks/#SA1012
- CTX-003: https://pkg.go.dev/context#Context.Value https://staticcheck.dev/docs/checks/#SA1029
- CTX-004: https://pkg.go.dev/context#pkg-overview
- DATA-001: https://go.dev/wiki/CodeReviewComments#declaring-empty-slices
- DATA-002: https://pkg.go.dev/cmd/vet https://github.com/uber-go/guide/blob/master/style.md#use-field-names-to-initialize-structs
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
- FUNC-001: https://github.com/uber-go/guide/blob/master/style.md#function-grouping-and-ordering https://github.com/manuelarte/funcorder
- FUNC-002: https://go.dev/wiki/CodeReviewComments#named-result-parameters
- FUNC-003: https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md https://go-critic.com/overview.html#toomanyresultschecker
- IFACE-001: https://go.dev/wiki/CodeReviewComments#interfaces
- IFACE-002: https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
- IFACE-003: https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces
- LAYER-001: https://alistair.cockburn.us/hexagonal-architecture/
- LAYOUT-001: https://github.com/golang-standards/project-layout#directories-you-shouldnt-have https://go.dev/doc/modules/layout
- LAYOUT-002: https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages https://go.dev/doc/go1.4#internalpackages
- LAYOUT-003: https://go.dev/doc/modules/layout#multiple-commands https://github.com/golang-standards/project-layout#cmd
- LAYOUT-004: https://github.com/golang-standards/project-layout#overview
- LOG-001: https://go.dev/blog/slog https://pkg.go.dev/log/slog
- LOG-002: https://go.dev/blog/slog
- LOG-003: https://pkg.go.dev/log/slog#Logger.InfoContext https://go.dev/blog/slog
- MODERN-001: https://pkg.go.dev/slices https://go.dev/doc/go1.21
- NAMING-001: https://go.dev/wiki/CodeReviewComments#package-names https://google.github.io/styleguide/go/decisions
- NAMING-002: https://google.github.io/styleguide/go/decisions#underscores https://go.dev/wiki/CodeReviewComments#mixed-caps
- NAMING-003: https://go.dev/wiki/CodeReviewComments#initialisms https://google.github.io/styleguide/go/decisions#initialisms
- NAMING-004: https://google.github.io/styleguide/go/decisions#variable-names https://go.dev/wiki/CodeReviewComments#variable-names
- NAMING-005: https://google.github.io/styleguide/go/decisions#variable-name-vs-type https://google.github.io/styleguide/go/decisions#external-context-vs-local-names
- NAMING-006: https://google.github.io/styleguide/go/decisions#receiver-names https://go.dev/wiki/CodeReviewComments#receiver-names
- PERF-001: https://go.dev/wiki/Performance
- PERF-002: https://github.com/uber-go/guide/blob/master/style.md#prefer-strconv-over-fmt
- PERF-003: https://github.com/uber-go/guide/blob/master/style.md#avoid-repeated-string-to-byte-conversions
- PERF-004: https://go.dev/blog/allocation-optimizations https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity
- PERF-005: https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity
- PERF-006: https://pkg.go.dev/strings#Builder
- PERF-007: https://go.dev/blog/slices-intro#a-possible-gotcha
- PERF-008: https://go.dev/wiki/CodeReviewComments#pass-values
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
- STRUCT-001: https://pkg.go.dev/cmd/vet https://google.github.io/styleguide/go/decisions
- STRUCT-002: https://github.com/protocolbuffers/protobuf-go/blob/master/internal/pragma/pragma.go https://go.dev/ref/spec#Comparison_operators
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

## 3 Detection Map

| ID | Detection |
| --- | --- |
| ABSTRACT-001 | manual review (Verify whether the interface or factory already has a second implementation, a replacement, or a testing need.) |
| CODE-001 | manual review (Check whether the external error codes place self-made numbers alongside the standard classification) |
| CODE-002 | manual review (Check whether there is a branch that tests by error message text) |
| CODE-003 | manual review (Check whether an error returned to the outside comes directly from a downstream service) |
| COMMENT-001 | manual review (Check whether a comment restates what the code itself already shows) |
| COMMENT-002 | manual review (Check whether a doc comment begins with the name and forms a full sentence) |
| COMPOSE-001 | manual review (Verify whether embedded types play the role of a base class (Base / Abstract naming, multi-level embedding).) |
| COMPOSE-002 | manual review (Verify whether embedding changes the outer zero value, copy semantics, or exposed unrelated methods (especially for synchronization primitives).) |
| CONC-001 | manual review (Check that each go statement has a matching waiter and exit path) |
| CONC-002 | manual review (Check that every close is on the sender side) |
| CONC-003 | manual review (Check whether a public function returns a channel that leaves the caller waiting) |
| CONC-004 | manual review (Check whether the lock-protected field is exported and whether it is reachable only through methods) |
| CONST-001 | golangci-lint stylecheck(ST1003) |
| CONST-002 | manual review (Check whether the constant name merely restates the value) |
| CONST-003 | manual review (Check whether enumerations that increase from one source use hand-written numbers) |
| CTRL-001 | manual review (Check whether the normal path is indented by an else) |
| CTRL-002 | golangci-lint copyloopvar |
| CTX-001 | golangci-lint containedctx (The first-parameter convention requires manual verification.) |
| CTX-002 | golangci-lint staticcheck(SA1012) (SA1012 detects calls that pass a nil context to a function.) |
| CTX-003 | golangci-lint staticcheck(SA1029) (SA1029 detects the use of a built-in type as a context key.) |
| CTX-004 | manual review (Verify whether what is put in context.Value is request-scoped data.) |
| DATA-001 | manual review (Outputting a JSON array [] is an exception) |
| DATA-002 | golangci-lint govet(composites) (govet covers only positional literals from another package, and the composites check must be enabled) |
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
| FUNC-001 | golangci-lint funcorder |
| FUNC-002 | manual review (Check whether results are named when their meaning is not self-evident) |
| FUNC-003 | golangci-lint revive(argument-limit, function-result-limit) |
| IFACE-001 | manual review (Verify whether the package holding the interface matches the number of implementations: for a single implementation, whether it shares the package with the implementation; for multiple implementations, whether it is in the consumer's package.) |
| IFACE-002 | manual review (Verify whether exported types that implement an interface per an API contract have a compile-time assertion.) |
| IFACE-003 | manual review (Verify whether function signatures and struct fields contain *interface types.) |
| LAYER-001 | manual review (Check the direction of package imports; this can be locked down as a rule with golangci-lint's depguard) |
| LAYOUT-001 | manual review (Check that the directory containing go.mod is the repository root and that the source is not placed under src/) |
| LAYOUT-002 | manual review (Check whether packages used only by this module live under internal/) |
| LAYOUT-003 | manual review (Check that main packages are gathered under cmd/ and that directory names match the executable names) |
| LAYOUT-004 | manual review (Check that every repository directory has real content and a clear responsibility) |
| LOG-001 | grep regex \b(fmt\.Print|log\.Print) (Check by hand whether the match is a business log) |
| LOG-002 | manual review (Check whether log values are spliced into the message string) |
| LOG-003 | manual review (Check whether logs on the request path pass in ctx) |
| MODERN-001 | golangci-lint modernize (go fix also suggests the equivalent modern form) |
| NAMING-001 | golangci-lint stylecheck(ST1003) (Check the package name identifier) |
| NAMING-002 | golangci-lint stylecheck(ST1003) (Check the identifier naming style) |
| NAMING-003 | golangci-lint stylecheck(ST1003) |
| NAMING-004 | manual review (Decide by hand whether name length matches the scope) |
| NAMING-005 | manual review (Check whether the name repeats the type or the surrounding context) |
| NAMING-006 | manual review (Check whether receiver names are consistent within the same type and whether this / self is used) |
| PERF-001 | manual review (Check whether a performance change is accompanied by benchmark results or profile conclusions.) |
| PERF-002 | manual review (Check whether conversions between a single primitive type and a string use fmt.) |
| PERF-003 | manual review (Check whether the loop contains conversions between string and []byte for a fixed string.) |
| PERF-004 | manual review (Check whether the make call for the slice an append targets carries a capacity argument.) |
| PERF-005 | manual review (Check whether a map written entry by entry carries a capacity hint.) |
| PERF-006 | manual review (Check whether the loop contains string concatenation with +=.) |
| PERF-007 | manual review (Check whether a sub-slice cut from a large buffer and kept for a long time is copied.) |
| PERF-008 | manual review (Check whether a parameter that is only read through a dereference uses a pointer.) |
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
| STRUCT-001 | golangci-lint govet(copylocks) |
| STRUCT-002 | manual review (Verify whether exported structs that intentionally prohibit comparison have embedded an incomparable field.) |
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
