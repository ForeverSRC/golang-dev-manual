# 4. Dependencies and Project Layout

## (1) Layering

### 【SHOULD】LAYER-001 Layer code by responsibility so that dependencies point one-way toward the inner layers.

- Category: Dependencies and Project Layout/Layering
- Since: Go 1.0

The domain model layer contains no IO and does not import other layers; the application service layer depends on the domain model and declares the external capabilities it needs through interfaces; the adapter layer implements those interfaces and handles persistence, files, and network access; the entry layer only parses arguments and renders output. Outer layers import inner layers, and the names of outer packages never appear in inner-layer code.

**Why**

> The rule to obey is that code pertaining to the inside part should not leak into the outside part.
>
> — https://alistair.cockburn.us/hexagonal-architecture/

When dependencies run in reverse, domain logic is tied to IO details, and even a single test run has to bring along real files or a database; putting external capabilities behind interfaces declared in the inner layer means implementations can be swapped without changing the inner layer, and tests can inject in-memory doubles.

**Good**

```go
// inner layer declares the port, outer layer implements it
package service

type ClauseRepository interface {
	Load(ctx context.Context, path string) (*domain.Manual, error)
}
```

**Bad**

```go
// inner layer depends directly on the outer adapter
package service

import "example/internal/repository/jsonfile"

type manualService struct {
	repo *jsonfile.Repository
}
```

**References**

- https://alistair.cockburn.us/hexagonal-architecture/

**Detection**: manual review (Check the direction of package imports; this can be locked down as a rule with golangci-lint's depguard)

## (2) Package Organization

### 【MUST】PKG-001 The package name is the last segment of the import path and matches the directory it lives in.

- Category: Dependencies and Project Layout/Package Organization
- Since: Go 1.0

The directory name and the package clause are written as the same name, and multi-word directories are likewise written without separators. There are only three exceptions:

- a directory for an executable declares package main
- an external test package declares "<package under test>_test"
- a package imported only by generated code may use an underscore

**Why**

> Client code uses the package path when importing the package. By convention, the last element of the package path is the package name:
>
> — https://go.dev/blog/package-names

Imports use the path while references use the package name, so when the two disagree callers have to remember both names, and searching by directory no longer lines up with searching by identifier; tools and logs only recognize paths, so during troubleshooting you have to look the directory back up before you can learn the package name.

**Good**

```go
// directory gdm/internal/repository/jsonfile
package jsonfile
```

**Bad**

```go
// directory gdm/internal/repository/jsonfile
package jsonrepo
```

**References**

- https://go.dev/blog/package-names
- https://go.dev/wiki/CodeReviewComments#package-names

**Detection**: manual review (Check directory by directory that the package clause matches the directory name)

### 【MUST】PKG-003 Arrange imports in groups, one for the standard library and one for everything else, separated by a blank line.

- Category: Dependencies and Project Layout/Package Organization
- Since: Go 1.0

Sort by path within each group, and do not mix imports from different sources into the same group. This project's gci configuration defines three groups: standard, default, and the module prefix of this repository, left to the formatter to organize automatically.

**Why**

> Imports are organized in groups, with blank lines between them. The standard library packages are always in the first group.
>
> — https://go.dev/wiki/CodeReviewComments#imports

Grouping separates the standard library from external dependencies visually, so during review you can see at a glance which module a change added; when everything is mixed together you have to judge line by line whether a path is from the standard library, and the scope of the dependency change becomes invisible.

**Good**

```go
import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)
```

**Bad**

```go
import (
	"context"
	"github.com/spf13/cobra"
	"fmt"
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
)
```

**References**

- https://go.dev/wiki/CodeReviewComments#imports

**Detection**: golangci-lint gci (Enable gci in formatters, with group order standard, default, and the module prefix of this repository)

### 【MUST】PKG-004 Do not use dot imports.

- Category: Dependencies and Project Layout/Package Organization
- Since: Go 1.0

The only exception is when an external test package cannot sit at the same level as the package under test because of an import cycle, in which case the package under test may be brought in with import . in a _test file. Ordinary code always writes the package name prefix explicitly.

**Why**

> Except for this one case, do not use import . in your programs. It makes the programs much harder to read because it is unclear whether a name like Quux is a top-level identifier in the current package or in an imported package.
>
> — https://go.dev/wiki/CodeReviewComments#import-dot

A dot import pulls another package's exported names directly into the current scope, so when you read a name you cannot tell from its form which package it comes from, and finding its definition means guessing first and then trying; what you save is a prefix of a few characters, and what you pay is one extra judgment about origin on every read.

**Good**

```go
import "github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"

var m domain.Manual
```

**Bad**

```go
import . "github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"

var m Manual
```

**References**

- https://go.dev/wiki/CodeReviewComments#import-dot

**Detection**: grep regex ^\s*(import\s+)?\.\s+" (For each match, check whether it is a _test file in an import-cycle scenario)

### 【SHOULD】PKG-002 Do not create packages with no informational content, such as util, common, or helper.

- Category: Dependencies and Project Layout/Package Organization
- Since: Go 1.0

A package name should let callers see the domain boundary. A package whose name carries only the meaning of "utility" or "general" has no criteria for what belongs in it, so anything can be put there. Existing general-purpose logic should be split into named packages by its real responsibility, for example set operations in stringset and retries in retry.

**Why**

> Avoid uninformative package names like `util`, `utility`, `common`, `helper`, `model`, `testhelper`, and so on that would tempt users of the package to rename it when importing.
>
> — https://google.github.io/styleguide/go/decisions#package-names

Without a boundary there is no basis for splitting, so such packages accumulate dependencies from different sources over time, widening the compilation surface and multiplying import conflicts; by the time a split is needed, call sites are already scattered everywhere and the migration can only proceed import by import.

**Good**

```go
package stringset
```

**Bad**

```go
package util
```

**References**

- https://google.github.io/styleguide/go/decisions#package-names
- https://go.dev/blog/package-names

**Detection**: grep regex ^package\s+(util|utility|common|misc|helper|helpers|model)\b (Split each match into named packages according to its real responsibility)

### 【SHOULD】PKG-005 Exported names do not repeat the package name, and constructors prefer New.

- Category: Dependencies and Project Layout/Package Organization
- Since: Go 1.0

Callers already use the package name as a prefix, so repeating it in an exported type or method name only adds length. Use New when the package contains a single type, or when the constructor returns the type that shares the package's name; only when it returns another type in the package should that type appear in the name.

**Why**

> Since client code uses the package name as a prefix when referring to the package contents, the names for those contents need not repeat the package name. The HTTP server provided by the `http` package is called `Server`, not `HTTPServer`.
>
> — https://go.dev/blog/package-names

At the call site widget.NewWidget expands to widget.NewWidget(), where half the characters carry the same information; when there is a single type sharing the name, New is shorter and reads the same way as standard library entry points such as time.Now and list.New.

**Good**

```go
package manual

func New() *Manual
```

**Bad**

```go
package manual

func NewManual() *Manual
```

**References**

- https://go.dev/blog/package-names
- https://google.github.io/styleguide/go/decisions#package-names

**Detection**: manual review (Check whether exported names repeat the package name they live in)

### 【SHOULD】PKG-006 Blank imports appear only in the main package or in test files.

- Category: Dependencies and Project Layout/Package Organization
- Since: Go 1.0

A package imported solely for one side effect (registering a driver, initializing a global table) is by convention placed only at the program entry point or in the tests that need it. A blank import in library code imposes that side effect on every caller that imports it.

**Why**

> Packages that are imported only for their side effects (using the syntax `import _ "pkg"`) should only be imported in the main package of a program, or in tests that require them.
>
> — https://go.dev/wiki/CodeReviewComments#import-blank

A blank import's side effect takes hold at compile time, so when a library package carries one, callers have no way to opt out; concentrating them in the main package and tests keeps the side effect's reach aligned with where it is actually used.

**Good**

```go
// gdm/cmd/gdm-cli/main.go
import _ "net/http/pprof"
```

**Bad**

```go
// gdm/internal/server/server.go
import _ "net/http/pprof"
```

**References**

- https://go.dev/wiki/CodeReviewComments#import-blank

**Detection**: grep regex ^\s*(import\s+)?_\s+" (For each match, check whether the containing file is in the main package or is a _test file)

### 【SHOULD】PKG-007 Give an import alias only when package names collide or the name does not match the last path segment.

- Category: Dependencies and Project Layout/Package Organization
- Since: Go 1.0

On a collision, rename the local or project-internal import and keep the upstream canonical name. The alias itself must follow the package name rules: all lowercase, with no underscores.

**Why**

> Avoid renaming imports except to avoid a name collision; good package names should not require renaming. In the event of collision, prefer to rename the most local or project-specific import.
>
> — https://go.dev/wiki/CodeReviewComments#imports

An alias is an extra layer of mapping, so reading the code requires converting between the alias and the original package name, and when the same path gets different aliases in different files it can no longer be searched by name; renaming only on a genuine collision keeps the mappings to remember to a minimum.

**Good**

```go
import (
	"github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	"github.com/go-playground/validator/v10"
)
```

**Bad**

```go
import (
	d "github.com/ForeverSRC/golang-dev-manual/gdm/internal/domain"
	v "github.com/go-playground/validator/v10"
)
```

**References**

- https://go.dev/wiki/CodeReviewComments#imports

**Detection**: manual review (Check whether each alias resolves a name collision or a mismatch between package name and path)

## (3) Dependency Management

### 【MUST】DEP-001 Leave go.mod and go.sum to the go command, and run go mod tidy before committing.

- Category: Dependencies and Project Layout/Dependency Management
- Since: Go 1.11

Add or remove dependencies with go get, change require and replace with go mod edit, tidy with go mod tidy, and let only the tooling write go.sum. Editing version numbers or checksums by hand will not match the dependency graph or the module cache.

**Why**

> To keep your managed dependency set tidy, use the `go mod tidy` command. Using the set of packages imported in your code, this command edits your go.mod file to add modules that are necessary but missing. It also removes unused modules that don't provide any relevant packages.
>
> — https://go.dev/doc/modules/managing-dependencies#synchronizing

The dependency graph contains many indirect dependencies, and manual maintenance can only reach the direct dependency layer; go mod tidy works backward from the actual imports, adding what is missing and removing what is no longer used, with go.sum corrected accordingly, so that the build result matches the source.

**Good**

```go
$ go get github.com/spf13/cobra@v1.10.2
$ go mod tidy
```

**Bad**

```go
// adding a dependency: hand-write a line into both go.mod and go.sum
require github.com/spf13/cobra v1.10.2
```

**References**

- https://go.dev/doc/modules/managing-dependencies#synchronizing
- https://go.dev/doc/modules/gomod-ref

**Detection**: manual review (Run go mod tidy before committing, and confirm go.mod and go.sum have no uncommitted changes)

### 【MUST】DEP-002 Commit go.sum to version control, and do not list it in .gitignore.

- Category: Dependencies and Project Layout/Dependency Management
- Since: Go 1.11

Commit go.mod and go.sum together. When fetching private modules fails, widen the verification scope with environment variables such as GOPRIVATE and GONOSUMDB rather than deleting or ignoring go.sum.

**Why**

> Include the go.mod and go.sum files in your repository with your code.
>
> — https://go.dev/doc/modules/managing-dependencies#enable_tracking

go.sum records the checksum of each module's content and is the basis for integrity during fetching; when it is not committed, CI and colleagues' machines each recompute it, changes from a moved upstream tag or altered content go undetected, and build results cannot be compared across machines.

**Good**

```go
// repository root
go.mod
go.sum
```

**Bad**

```go
// .gitignore
go.sum
```

**References**

- https://go.dev/doc/modules/managing-dependencies#enable_tracking

**Detection**: manual review (Check that the repository contains go.sum and that .gitignore does not ignore it)

### 【SHOULD】DEP-003 Do not keep replace directives pointing at local directories or private forks long term.

- Category: Dependencies and Project Layout/Dependency Management
- Since: Go 1.11

For local integration work, replace may point at a sibling directory or your own fork, and must be removed before merging. When a long-term replacement is genuinely needed, write the reason and the removal condition into a comment in go.mod, and periodically check whether upstream has merged the change.

**Why**

> When you use the replace directive, Go tools don't authenticate external modules as described in Adding a dependency.
>
> — https://go.dev/doc/modules/managing-dependencies#local_directory

A replace directive makes the substituted module skip checksum verification, so its origin can no longer be verified; a local path exists only on your own machine, so anyone else receiving this go.mod fails to build outright, while pointing at a fork stays permanently behind upstream fixes.

**Good**

```go
// send the needed change upstream first, then bump the version after it is merged
require example.com/lib v1.4.0
```

**Bad**

```go
// keep a modified copy locally and pin it with replace long-term
require example.com/lib v1.4.0

replace example.com/lib => ../lib
```

**References**

- https://go.dev/doc/modules/managing-dependencies#local_directory
- https://go.dev/doc/modules/gomod-ref#replace

**Detection**: manual review (Check that replace directives in go.mod do not point at local directories long term and that removal conditions are stated)

### 【SHOULD】DEP-004 Declare tool dependencies in go.mod's tool directive rather than relying on global installs.

- Category: Dependencies and Project Layout/Dependency Management
- Since: Go 1.24

Command-line tools declare their package path with the tool directive and run via go tool <name>, with the version pinned by go.mod and go.sum. A version obtained from a global go install varies by machine and by when it was run, so CI and local environments cannot be brought into line.

**Why**

> Adds a package as a dependency of the current module, and makes it available to run with `go tool` when the current working directory is within this module.
>
> — https://go.dev/doc/modules/gomod-ref#tool

The version installed globally is determined by @latest at the moment install runs, so the same command may hit different versions on different machines and check results cannot be reproduced; once it is written into the tool directive the version enters the dependency graph, and everyone runs the same one.

**Good**

```go
// go.mod
tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint

$ go tool golangci-lint run ./...
```

**Bad**

```go
$ go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
$ golangci-lint run ./...
```

**References**

- https://go.dev/doc/modules/gomod-ref#tool
- https://go.dev/doc/modules/managing-dependencies#tools

**Detection**: manual review (Check that the command-line tools used by the repository are declared through go.mod's tool directive)

### 【SHOULD】DEP-005 Do not bring in a third-party dependency for a small capability the standard library already covers.

- Category: Dependencies and Project Layout/Dependency Management
- Since: Go 1.0

Write it yourself when it fits in one line of code or the standard library already has equivalent capability; general-purpose and non-trivial capabilities (CLI frameworks, test assertions, structured logging) are worth pulling in only when the benefit is real. Before adding anything, first confirm the standard library has no ready-made capability.

**Why**

> A little copying is better than a little dependency.
>
> — https://go-proverbs.github.io/

Every module you pull in brings its release cadence, security fixes, and API changes into your own maintenance surface; pulling in a whole module for one small function usually yields less benefit than the future cost of tracking version upgrades.

**Good**

```go
n, err := strconv.Atoi(s)
```

**Bad**

```go
n, err := convert.ToInt(s) // pull in the convert module for a single type conversion
```

**References**

- https://go-proverbs.github.io/

**Detection**: manual review (Check whether a newly added dependency is a small capability the standard library already covers)

## (4) Project Layout

### 【MUST】LAYOUT-001 Put go.mod and the source code at the module root, without an extra src/ layer.

- Category: Dependencies and Project Layout/Project Layout
- Since: Go 1.11

Import paths are measured from the module root, and the source tree unfolds at the same level as its parent directories. A src/ carried over from the GOPATH era or from Java habits inserts a meaningless prefix between the repository path and the import path.

**Why**

> Some Go projects do have a `src` folder, but it usually happens when the devs came from the Java world where it's a common pattern. If you can help yourself try not to adopt this Java pattern.
>
> — https://github.com/golang-standards/project-layout#directories-you-shouldnt-have

The go command treats the directory containing go.mod as the module root, and the import path is assembled from the relative path below that directory; nesting another src/ under the root forces every import, go install package path, and documentation example to add an extra layer, and it does not line up with the repository's actual directories.

**Good**

```go
project-root/
  go.mod
  main.go
  internal/
```

**Bad**

```go
project-root/
  src/
    go.mod
    main.go
```

**References**

- https://github.com/golang-standards/project-layout#directories-you-shouldnt-have
- https://go.dev/doc/modules/layout

**Detection**: manual review (Check that the directory containing go.mod is the repository root and that the source is not placed under src/)

### 【SHOULD】LAYOUT-002 Put packages not meant to be exposed externally under internal/.

- Category: Dependencies and Project Layout/Project Layout
- Since: Go 1.4

An internal directory is visible only within the subtree of the directory that contains it, and when the repository has a single module the top-level internal/ covers the whole repository. Implementation details, adapters, and types used only by this program all go inside; only packages genuinely meant for external reuse are placed outside internal.

**Why**

> Initially, it's recommended placing such packages into a directory named `internal`; this prevents other modules from depending on packages we don't necessarily want to expose and support for external uses. Since other projects cannot import code from our `internal` directory, we're free to refactor its API and generally move things around without breaking external users.
>
> — https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages

A package placed outside internal can be imported by other modules at a fixed path, and once someone uses it, changing the package name or a signature becomes a breaking change; once it is inside internal the compiler blocks external imports outright, leaving the internal structure free to adjust.

**Good**

```go
// external modules cannot import it
import "github.com/ForeverSRC/golang-dev-manual/gdm/internal/repository/jsonfile"
```

**Bad**

```go
// the implementation package sits at the module root, so its API is locked once another module depends on it
import "github.com/ForeverSRC/golang-dev-manual/repository/jsonfile"
```

**References**

- https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages
- https://go.dev/doc/go1.4#internalpackages

**Detection**: manual review (Check whether packages used only by this module live under internal/)

### 【SHOULD】LAYOUT-003 Keep executables together under cmd/, with one subdirectory per program.

- Category: Dependencies and Project Layout/Project Layout
- Since: Go 1.0

When a repository holds both importable packages and executables, main packages are collected uniformly into cmd/<program name>/, with the directory name matching the generated executable. Only when there is a single program and the repository contains no importable packages may main.go be placed at the module root.

**Why**

> A common convention is placing all commands in a repository into a `cmd` directory; while this isn't strictly necessary in a repository that consists only of commands, it's very useful in a mixed repository that has both commands and importable packages, as we will discuss next.
>
> — https://go.dev/doc/modules/layout#multiple-commands

When main packages are mixed into the module root or business directories, the import path does not match the executable name and go install first requires locating the directory holding the file; once they are gathered into cmd/, packages and commands are separated structurally and each program's entry point is obvious at a glance.

**Good**

```go
project-root/
  cmd/
    gdm/
      main.go
  internal/
    service/
```

**Bad**

```go
project-root/
  gdm.go
  internal/
    server/
      main.go
```

**References**

- https://go.dev/doc/modules/layout#multiple-commands
- https://github.com/golang-standards/project-layout#cmd

**Detection**: manual review (Check that main packages are gathered under cmd/ and that directory names match the executable names)

### 【MAY】LAYOUT-004 Do not lay out a large project's directory skeleton in advance while the project is still small.

- Category: Dependencies and Project Layout/Project Layout
- Since: Go 1.11

Create cmd, internal, pkg, configs, and scripts only once there is real content for them, and do not create empty directories or directories holding only a placeholder file. Directories grow along with the code, and their names come from real content.

**Why**

> If you are trying to learn Go or if you are building a PoC or a simple project for yourself this project layout is an overkill. Start with something really simple instead (a single `main.go` file and `go.mod` is more than enough).
>
> — https://github.com/golang-standards/project-layout#overview

An empty directory carries no information yet makes people assume a corresponding structural convention exists, so newcomers stuff things that could belong to a specific package into the reserved pkg/ and util/; creating directories once the real responsibility appears is what lets the names reflect the content.

**Good**

```go
project-root/
  go.mod
  main.go
```

**Bad**

```go
project-root/
  cmd/
  pkg/
  configs/
  scripts/
  internal/
  // each directory holds only a .gitkeep
```

**References**

- https://github.com/golang-standards/project-layout#overview

**Detection**: manual review (Check that every repository directory has real content and a clear responsibility)

