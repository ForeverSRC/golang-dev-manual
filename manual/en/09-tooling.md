# 9. Tooling

## (1) Command Line

### 【SHOULD】TOOL-004 Command line programs organize subcommands with `cobra`.

- Category: Tooling/Command Line
- Since: Go 1.0

Attach subcommands to the root command and put business logic in each command's `RunE`; write output to `cmd.OutOrStdout()` and errors to `cmd.ErrOrStderr()`.

**Why**

> Cobra is a library providing a simple interface to create powerful modern CLI interfaces similar to git & go tools. ... Easy subcommand-based CLIs: app server, app fetch, etc.
>
> — https://github.com/spf13/cobra

Hand-rolled `os.Args` dispatch has to reinvent subcommands, global flags, and help text; `cobra` injects output and arguments into commands, so commands can be executed directly in tests.

**Good**

```go
cmd := &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, _ []string) error {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), "ok")
		return err
	},
}
```

**Bad**

```go
switch os.Args[1] {
case "list":
	fmt.Println("ok")
}
```

**References**

- https://github.com/spf13/cobra

**Detection**: manual review (Check whether subcommands are organized with `cobra`)

### 【SHOULD】TOOL-008 After upgrading the Go toolchain, apply modernization rewrites with `go fix`.

- Category: Tooling/Command Line
- Since: Go 1.26

After switching the build to a newer Go version, run `go fix ./...` in a clean git working tree and keep the changes in a separate commit for item-by-item review.

- Preview the patch without applying it with `go fix -diff ./...`
- Run all analyzers by default; when splitting changes, use `-<analyzer>` to run only one (for example `-any`), and `-<analyzer>=false` to exclude it
- Rewrites that land on generated files are discarded; such changes belong in the generator itself
- Run once per GOOS / GOARCH to cover code under different build tags
- When a package author migrates their own API, write `//go:fix inline` on the replaced function and let the inline analyzer apply it

The modernize check and the tool in this section point to the same set of analyzers; the difference is that `golangci-lint` reports during review, while `go fix` produces the patch directly.

**Why**

> Fix is a tool executed by "go fix" to update Go programs that use old features of the language and library and rewrite them to use newer ones. After you update to a new Go release, fix helps make the necessary changes to your programs.
>
> — https://pkg.go.dev/cmd/fix

Each release adds new idioms to the language and standard library, and old code that is not rewritten turns into technical debt release by release; relying on human memory and review to raise them leaves the spots nobody thought of unknown. `go fix` collects these rewrites into one reproducible command that only needs to be run once alongside a toolchain upgrade.

**Good**

```console
$ go fix -diff ./...   # preview the patch first
$ go fix ./...          # then apply it
```

**Bad**

```console
$ go vet ./...
# after upgrading the toolchain only static checks ran, so old idioms were never rewritten in bulk
```

**References**

- https://pkg.go.dev/cmd/fix
- https://go.dev/blog/gofix
- https://go.dev/doc/go1.26

**Detection**: manual review (Check whether `go fix` was run after the Go version upgrade)

## (2) Testing Tools

### 【MUST】TOOL-005 Unit test assertions uniformly use `testify`.

- Category: Tooling/Testing Tools
- Since: Go 1.0

Use `assert` or `require` for assertions, choosing by whether a failure should stop the test case: use `require` for preconditions and `assert` for all other checks. Do not hand-write `if` statements with `t.Errorf` or `t.Fatalf` for comparisons in test files.

**Why**

> The assert package provides some helpful methods that allow you to write better test code in Go. Prints friendly, easy to read failure descriptions.
>
> — https://github.com/stretchr/testify#assert-package

Hand-written comparisons produce failure messages that are only custom text and lack an expected-versus-actual comparison; `testify` formats diffs uniformly and includes the call site, so debugging does not require rereading the source.

**Good**

```go
assert.Equal(t, want, got)
```

**Bad**

```go
if got != want {
	t.Errorf("got %v, want %v", got, want)
}
```

**References**

- https://github.com/stretchr/testify#assert-package

**Detection**: grep regex \bt\.(Errorf|Fatalf)\( (Matches should be replaced with `testify` assertions)

## (3) Code Generation

### 【MUST】TOOL-010 Generated files are not edited by hand; they are produced only by the generator.

- Category: Tooling/Code Generation
- Since: Go 1.4

The generator and its input sources are the only places to change; regenerate after changing them and never edit the generated file directly. When the generator emits a marker line such as `// Code generated ... DO NOT EDIT.`, keep it intact and do not delete or alter it; when the generator emits no marker, do not fabricate one, since identification is handled by the file name convention in TOOL-011. This repository's gdm/cmd/gdm-cli/wire/wire_gen.go and manual/ both fall into this category, and their generation commands are `make generate` and `make gen`.

**Why**

> To convey to humans and machine tools that code is generated, generated source should have a line that matches the following regular expression (in Go syntax): ^// Code generated .* DO NOT EDIT\.$
>
> — https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source

A generated file is a projection of its input sources and generator, and editing it directly breaks that relationship: the next generation overwrites the change, or the generated file drifts from the input sources over time and becomes a piece of code nobody dares to touch. The official way to identify generated Go source is an in-file marker line (quoted below), but the marker must be emitted by the generator; template rendering and document generation do not write one, so those generated files are identified by file name instead, see TOOL-011.

**Good**

```go
// gdm/di/wire.go: the wiring change goes into the input source, then run make generate
var ProviderSet = wire.NewSet(jsonfile.New)
```

**Bad**

```go
// gdm/cmd/gdm-cli/wire/wire_gen.go: the change is made on the generated file and is overwritten on the next run
repository := jsonfile.New()
```

**References**

- https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source

**Detection**: manual review (Check whether changes land on the generator or on the input sources)

### 【MUST】TOOL-011 Generated files are named with a gen prefix or suffix, preferring the suffix.

- Category: Tooling/Code Generation
- Since: Go 1.0

When the generator allows specifying the output file name, name the file with a gen prefix or suffix:

- Prefer the suffix, as in wire_gen.go and order_repo_gen.go, so that name sorting keeps it next to files of the same kind.
- Write the prefix as gen_, as in gen_version.go.
- Use only underscores as separators, never hyphens.
- When the file name is fixed by the generator itself (`wire` outputs wire_gen.go) or decided by the platform, follow the generator default and do not rename it.

**Why**

> The name of the output file is created by replacing the .proto extension with .pb.go.
>
> — https://protobuf.dev/reference/go/go-generated/

`protoc` adds the .pb.go suffix to generated files and `wire` outputs wire_gen.go, both marking the file name as generated. When the name is chosen by the caller with no convention, generated files and hand-written files mix in the same directory, so a reviewer must open each one to learn which can be edited and must look up the generation command before editing. With a uniform gen prefix or suffix, the directory listing itself answers that; the suffix is recommended because it does not break the existing naming order, leaving retrieval, sorting, and adjacency driven by the business-name prefix. The separator follows its own file name convention: the `_test.go` and `source_windows_amd64.go` forms recognized by the Go toolchain both use underscores, and hyphenated files in the standard library appear only in test data such as testdata.

**Good**

```go
// generated file names: wire_gen.go, mock_gen.go, clause_repo_gen.go
package wire
```

**Bad**

```go
// generated file names: wire.go, mocks.go, clause_repo.go
package wire
```

**References**

- https://protobuf.dev/reference/go/go-generated/
- https://pkg.go.dev/cmd/go#hdr-Build_constraints

**Detection**: manual review (Check whether generated file names carry a gen prefix or suffix)

