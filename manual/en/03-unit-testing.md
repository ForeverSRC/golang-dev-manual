# 3. Unit Testing

## (1) Test Naming and Structure

### 【MUST】TEST-005 Test files use external package tests, with the package name ending in _test.

- Category: Unit Testing/Test Naming and Structure
- Since: Go 1.0

A test file declares the package as "the package under test_test", writes cases only against the public API, does not access package-internal members, and does not drag internal implementation into tests just to pad coverage. Unexported logic is covered indirectly by cases exercising the package's public entry points.

**Why**

> Test files that declare a package with the suffix "_test" will be compiled as a separate package, and then linked and run with the main test binary.
>
> — https://pkg.go.dev/cmd/go#hdr-Test_packages

Same-package tests can reach unexported members, so cases get written tightly against the current implementation and a change to the internal structure turns the tests red first; external package tests force the public entry points, so what is tested is the outward contract, and room for refactoring is preserved.

**Good**

```go
package jsonfile_test

import "github.com/ForeverSRC/golang-dev-manual/gdm/internal/repository/jsonfile"
```

**Bad**

```go
package orderbook

func TestMatch(t *testing.T) {
	// in-package tests can construct the unexported matcher directly
	m := &matcher{}
	_ = m
}
```

**References**

- https://pkg.go.dev/cmd/go#hdr-Test_packages

**Detection**: golangci-lint testpackage

### 【SHOULD】TEST-001 Pure function tests use the native testing package, and component tests use the testify suite.

- Category: Unit Testing/Test Naming and Structure
- Since: Go 1.24

For functions that only map input to output, with no object under test or dependencies to construct, write `func TestXxx(t *testing.T)` plus assertions directly, organizing multiple inputs into a table with `t.Run`; for components that encapsulate behavior per type, need a constructor to produce the object under test, or need dependencies assembled and a shared fixture prepared, embed `suite.Suite`, put the object under test and its dependencies into struct fields assembled once in `SetupSuite`, and run everything through a single entry function that calls `suite.Run`.

**Why**

> The suite package provides functionality that you might be used to from more common object-oriented languages. With it, you can build a testing suite as a struct, build setup/teardown methods and testing methods on your struct, and run them with 'go test' as per normal.
>
> — https://github.com/stretchr/testify#suite-package

Pure functions have no dependencies and no state shared across cases, so wrapping them in a `suite` only adds three layers of boilerplate — the struct, the embedding, and the run entry point — and reading one case means jumping across three places; a component's object under test and dependencies are reused across many cases, and constructing them separately inside each case is repetitive and easy to miss when dependencies grow, while the `suite`'s fields and `SetupSuite` are exactly where this is held centrally.

**Good**

```go
type RepositorySuite struct {
	suite.Suite

	underTest *jsonfile.Repository
}

func (s *RepositorySuite) SetupSuite() {
	s.underTest = jsonfile.New()
}

func (s *RepositorySuite) TestLoad() {
	got, err := s.underTest.Load(s.T().Context(), "testdata/manual.json")
	s.Require().NoError(err)
	s.NotNil(got)
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
```

**Bad**

```go
func TestLoad(t *testing.T) {
	got, err := jsonfile.New().Load(t.Context(), "testdata/manual.json")
	require.NoError(t, err)
	assert.NotNil(t, got)
}

func TestLoadError(t *testing.T) {
	_, err := jsonfile.New().Load(t.Context(), "missing.json")
	require.Error(t, err)
}
```

**References**

- https://github.com/stretchr/testify#suite-package
- https://pkg.go.dev/github.com/stretchr/testify/suite

**Detection**: manual review (Check whether pure function cases are wrapped in a suite, and whether component cases assemble the object under test centrally with a suite.)

### 【SHOULD】TEST-002 Unit tests and integration tests are stored separately.

- Category: Unit Testing/Test Naming and Structure
- Since: Go 1.24

Tests whose input is constructed by the case itself and which neither read nor write real files or the network are unit tests, written directly inside that package; tests that need real files, real dependencies, or an end-to-end path are integration tests, named with the _it_test.go suffix, in the same package as the unit tests, and executed together by `go test`. Integrity checks on real data belong to integration tests.

**Why**

> Go test recompiles each package along with any files with names matching the file pattern "*_test.go".
>
> — https://pkg.go.dev/cmd/go#hdr-Test_packages

A failure in the two kinds of tests means different things: a red unit test is a logic error, while a red integration test may be an environment or data problem. Mixed together, environment fluctuation is taken for a code defect, and the investigation heads in the wrong direction from the start.

**Good**

```go
// unit test: input lives in testdata, no external resources
func TestLoad(t *testing.T) {
	got, err := jsonfile.New().Load(t.Context(), filepath.Join("testdata", "manual.json"))
	require.NoError(t, err)
}
```

**Bad**

```go
// reading real repo data in a unit test file, result depends on the working directory
func TestLoad(t *testing.T) {
	got, err := jsonfile.New().Load(t.Context(), "../../../data/manual.json")
	require.NoError(t, err)
}
```

**References**

- https://pkg.go.dev/cmd/go#hdr-Test_packages

**Detection**: manual review (Check whether tests that read real files, the network, or a database end with _it_test.go.)

### 【MAY】TEST-008 Test case names preferably follow the should-plus-outcome, when-plus-condition pattern.

- Category: Unit Testing/Test Naming and Structure
- Since: Go 1.0

Write the subtest name given to `t.Run` and the name field in a table as should followed by the expected outcome and when followed by the triggering condition, for example should return error when file is missing. When the scenario is clear at a glance, a short Chinese or English description works just as well; do not stretch the name just to fit the pattern.

**Why**

> When you use t.Run to create a subtest, the first argument is used as a descriptive name for the test. To ensure that test results are legible to humans reading the logs, choose subtest names that will remain useful and readable after escaping.
>
> — https://go.dev/wiki/TestComments#choose-human-readable-subtest-names

Subtest names go straight into `go test` output, and people judge which one failed from the name alone; upstream only requires the name to be readable, and putting the expected outcome and the triggering condition into the name is a stable way to achieve readability, so a failure can be located without rereading the case code.

**Good**

```go
{name: "should return error when file is missing"}
```

**Bad**

```go
{name: "error branch"}
```

**References**

- https://go.dev/wiki/TestComments#choose-human-readable-subtest-names

**Detection**: manual review (Check whether t.Run names and table name fields state the expected outcome and the triggering condition clearly.)

## (2) Table-Driven and Case Design

### 【MUST】TEST-012 Write expected values independently; do not produce them with the function under test.

- Category: Unit Testing/Table-Driven and Case Design
- Since: Go 1.0

Write the expected value by hand as a literal, or assemble it with a factory unrelated to the implementation under test. Do not call the function under test itself, nor the few steps it reuses internally, to produce the expected result. When the data is too large to hand-write, put it into an expected file under testdata and fix it in advance, rather than going back to the code under test to fill it in. When a reference implementation is used for comparison, it must be a separate implementation independent of the code under test.

**Why**

> Instead, construct the struct that you're expecting your function to return, and compare in one shot using diffs or deep comparisons.
>
> — https://go.dev/wiki/TestComments#compare-full-structures

Generating the expected value with the code under test makes both sides wrong together, leaving the case green forever; such a case verifies that something equals itself, changes along with a broken implementation, and catches no defect.

**Good**

```go
want := []domain.Clause{
	{ID: "NAMING-001", Level: domain.LevelMust, Category: "Programming Conventions/Naming"},
}
assert.Equal(t, want, manual.Filter([]string{"MUST"}, "Programming Conventions"))
```

**Bad**

```go
want := manual.MustClauses() // another function in the package under test
assert.Equal(t, want, manual.Filter([]string{"MUST"}, "Programming Conventions"))
```

**References**

- https://go.dev/wiki/TestComments#compare-full-structures
- https://research.swtch.com/testing

**Detection**: manual review (Check whether expected values are produced by functions or methods inside the package under test.)

### 【SHOULD】TEST-010 Only cases that share the same check logic go into the same table.

- Category: Unit Testing/Table-Driven and Case Design
- Since: Go 1.24

When multiple inputs share one set of check logic (build the input, call, compare the result), write it as a table walk, where adding a case means adding one row of data. When some cases are checked differently, split them into multiple test functions; only when the setup is the same but the checks differ may you write them as a sequence of flat subtests inside a single test function. Do not distinguish kinds of cases with conditional branches in the loop body.

**Why**

> When some test cases need to be checked using different logic from other test cases, it is more appropriate to write multiple test functions. ... If they have different logic but identical setup, a sequence of subtests within a single test function might also make sense.
>
> — https://go.dev/wiki/TestComments#table-driven-tests-vs-multiple-test-functions

The benefit of table-driven tests is that the check logic is written once. Cases with different checks forced into one table make the loop body decide which set of assertions to run by branching, forcing a reader to jump back and forth between data and logic for a single case, and changing one piece of logic means confirming which rows it affects.

**Good**

```go
func TestLoadMissingFile(t *testing.T) {
	_, err := jsonfile.New().Load(t.Context(), "missing.json")
	require.Error(t, err)
}

func TestLoadInvalidJSON(t *testing.T) {
	_, err := jsonfile.New().Load(t.Context(), "invalid.json")
	require.Error(t, err)
}
```

**Bad**

```go
func TestLoad(t *testing.T) {
	for _, path := range []string{"missing.json", "invalid.json"} {
		_, err := jsonfile.New().Load(t.Context(), path)
		if path == "missing.json" {
			require.True(t, os.IsNotExist(err))
			continue
		}
		require.Error(t, err)
	}
}
```

**References**

- https://go.dev/wiki/TestComments#table-driven-tests-vs-multiple-test-functions
- https://research.swtch.com/testing

**Detection**: manual review (Check whether a single loop body switches check logic with conditional branches.)

### 【SHOULD】TEST-011 The table holds only the fields that vary per case.

- Category: Unit Testing/Table-Driven and Case Design
- Since: Go 1.0

Keep only the triggering condition and the expected outcome as table fields; move request parameters shared by all cases and common setup operations out into variables outside the loop. When a few cases need different stub behavior, carry that difference in a field holding a closure, rather than adding a branch in the loop body or copying the whole stub configuration into every row.

**Why**

> Separate test cases from test logic.
>
> — https://research.swtch.com/testing

Repeating the same parameters on every row makes the table longer without adding information, and readers must compare rows one by one to find where the differences actually are; changing a shared parameter means editing row by row, and missing one row leaves a case testing something else.

**Good**

```go
query := &domain.Query{Category: "Programming Conventions"} // shared by all cases

tests := []struct {
	name  string
	level string
	want  int
}{
	{name: "should return all clauses when level is empty", level: "", want: 24},
	{name: "should return must clauses when level is MUST", level: "MUST", want: 12},
}
```

**Bad**

```go
tests := []struct {
	name  string
	query *domain.Query
	level string
	want  int
}{
	{name: "should return all clauses when level is empty", query: &domain.Query{Category: "Programming Conventions"}, level: "", want: 24},
	{name: "should return must clauses when level is MUST", query: &domain.Query{Category: "Programming Conventions"}, level: "MUST", want: 12},
}
```

**References**

- https://research.swtch.com/testing
- https://go.dev/wiki/TableDrivenTests

**Detection**: manual review (Check whether table fields are written with the same value by every case.)

### 【SHOULD】TEST-013 Create separate cases for boundaries and special inputs.

- Category: Unit Testing/Table-Driven and Case Design
- Since: Go 1.0

Beyond the normal path, create a case for each boundary input:

- empty collections and nil
- a single element, the first and last elements
- out-of-range indices, keys falling between two elements
- threshold numeric values

Every branch of an error path has its own case, not just one of them. The set of cases should be able to answer which input combinations have not been thought of yet.

**Why**

> Look for special cases.
>
> — https://research.swtch.com/testing

Coverage shows that code was executed, not that boundaries were considered. Defects such as empty input and out-of-range access appear only at boundaries, which cases along the normal path never reach; a field left unassigned also passes silently because it equals the zero value of the expected value.

**Good**

```go
tests := []struct {
	name   string
	levels []string
}{
	{name: "should return empty when levels is nil", levels: nil},
	{name: "should return empty when levels is empty", levels: []string{}},
	{name: "should match single element when levels has one", levels: []string{"MUST"}},
}
```

**Bad**

```go
tests := []struct {
	name   string
	levels []string
}{
	{name: "should match levels", levels: []string{"MUST", "SHOULD"}},
}
```

**References**

- https://research.swtch.com/testing

**Detection**: manual review (Check whether boundary inputs such as empty collections, single elements, the first and last elements, and out-of-range access are covered.)

### 【SHOULD】TEST-014 A defect fix comes with a reproducing case.

- Category: Unit Testing/Table-Driven and Case Design
- Since: Go 1.0

When fixing a defect, add a case that reproduces it in the same commit, with input fixed to the conditions that trigger the defect, failing before the fix and passing after it. The triggering conditions stay in the case, not in a review note or a record of manual steps.

**Why**

> If you didn't add a test, you didn't fix the bug.
>
> — https://research.swtch.com/testing

Without a reproducing case, the triggering conditions exist only within that one investigation, and the same defect reappears in later refactoring; upstream puts it bluntly — if you didn't add a test, you didn't fix the bug.

**Good**

```go
// bug: case mismatch when filtering by id caused a missed match
func TestFilterByID(t *testing.T) {
	want := []domain.Clause{{ID: "NAMING-001"}}
	assert.Equal(t, want, manual.FilterByID("naming-001"))
}
```

**Bad**

```go
// the test input also passed before the fix, the trigger condition was not pinned down
func TestFilterByID(t *testing.T) {
	assert.NotEmpty(t, manual.FilterByID("NAMING-001"))
}
```

**References**

- https://research.swtch.com/testing

**Detection**: manual review (Check whether a commit that fixes a defect carries a case matching the triggering conditions.)

## (3) Assertions and Test Data

### 【SHOULD】TEST-003 Compare the complete result in one shot; do not break it apart field by field or trim it into a subset.

- Category: Unit Testing/Assertions and Test Data
- Since: Go 1.0

Write the expected value in one shot, as a literal or through an independent constructor, spelling out fields expected to be zero values as well, and then compare it against the method's complete return value: if it returns a slice, compare the slice itself, neither mapping it first into a subset such as a list of IDs, nor splitting it into one assertion per field. Only when the result contains uncontrollable content (timestamps, random IDs, floating-point precision loss) or fields that do not support equality comparison should you assert field by field or switch to a dedicated comparison helper. When multiple values are returned, compare them one by one; there is no need to wrap them into a struct first.

**Why**

> If your function returns a struct, don't write test code that performs an individual comparison for each field of the struct. Instead, construct the struct that you're expecting your function to return, and compare in one shot using diffs or deep comparisons. The same rule applies to arrays and maps.
>
> — https://go.dev/wiki/TestComments#compare-full-structures

Once assertions are broken apart or narrowed, the missing part is invisible from the code: with per-field assertions, a newly added field does not turn old cases red, and once projected into a subset, fields such as level and category drop out of the comparison entirely, so cases still pass even when the code under test returns wrong content; comparing the whole result exposes the omission on the spot and brings field additions and removals into scope automatically.

**Good**

```go
got := manual.Filter([]string{"MUST"}, "Programming Conventions")
want := []domain.Clause{{ID: "A-001", Level: domain.LevelMust, Category: "Programming Conventions/Naming"}}
assert.Equal(t, want, got)
```

**Bad**

```go
got := manual.Filter([]string{"MUST"}, "Programming Conventions")
assert.Len(t, got, 1)
ids := make([]string, 0, len(got))
for _, c := range got {
	ids = append(ids, c.ID)
}
assert.Equal(t, []string{"A-001"}, ids)
```

**References**

- https://go.dev/wiki/TestComments#compare-full-structures

**Detection**: manual review (Check whether assertions are split field by field, or whether the result is trimmed into a subset before comparison.)

### 【SHOULD】TEST-004 Bulky test inputs and expected data go into the testdata directory.

- Category: Unit Testing/Assertions and Test Data
- Since: Go 1.16

When the input payload or the expected result is long, write it as a separate file under testdata/, kept apart from the code; the case embeds the file into the test binary with //go:embed and then reads it, rather than locating the file by joining paths from the working directory. When cases need to be driven in bulk by file name, embed an `embed.FS` and walk it with `fs.Glob`. When the object under test itself takes a file path as input (a file loader or a directory scanner, for example), pass the relative path under testdata directly. When expected results change intentionally along with the implementation, inject a -update flag to overwrite the expected file with the actual output, then review it by hand via git diff.

**Why**

> Go source files that import "embed" can use the //go:embed directive to initialize a variable of type string, []byte, or FS with the contents of files read from the package directory or subdirectories at compile time.
>
> — https://pkg.go.dev/embed

A large block of JSON written inside a Go string drowns the case logic in data, and changing one piece of data means touching the code; once it becomes a paired file, the data is readable and diffable. Reading a file by joining paths from the working directory makes it unfindable from another working directory (running the case directly in the IDE or starting `go test` from the repository root are both such cases); //go:embed pins the content into the binary at compile time, so the read result does not vary with the runtime environment.

**Good**

```go
import _ "embed"

//go:embed testdata/manual.json
var manualJSON []byte
```

**Bad**

```go
raw, err := os.ReadFile(filepath.Join("testdata", "manual.json"))
require.NoError(t, err)
```

**References**

- https://pkg.go.dev/embed
- https://pkg.go.dev/cmd/go#hdr-Test_packages

**Detection**: manual review (Check whether cases read testdata by joining paths from the working directory.)

### 【SHOULD】TEST-007 Assert on an error only by whether it is non-nil; do not compare error message text.

- Category: Unit Testing/Assertions and Test Data
- Since: Go 1.24

On a failure path, first assert whether an error was returned; when the kind of error needs to be distinguished, use `errors.Is` for sentinels and `errors.As` for error types, rather than string comparison of the message, and do not construct an identical error for value comparison. Only when the error message itself is part of the contract under test (for example, it must carry the input parameter name) should you assert on the message as a string, and only on such properties as are unaffected by wording.

**Why**

> don't use string comparison to check what type of error your function returns. ... It's OK to use string comparisons to check that error messages coming from the package under test satisfy some property, for example, that it includes the parameter name.
>
> — https://go.dev/wiki/TestComments#test-error-semantics

Error messages are meant for humans; judging the kind of error by string comparison makes the case turn red the moment the wording changes, and such a case catches no real error; to distinguish kinds of errors, the code under test must expose sentinel errors or error types, and the test then judges with `errors.Is` / `errors.As`.

**Good**

```go
_, err := jsonfile.New().Load(t.Context(), "missing.json")
require.Error(t, err)
```

**Bad**

```go
_, err := jsonfile.New().Load(t.Context(), "missing.json")
require.ErrorContains(t, err, "load")
```

**References**

- https://go.dev/wiki/TestComments#test-error-semantics
- https://go.dev/blog/go1.13-errors

**Detection**: manual review (Check whether cases judge the kind of error by error message text, and whether sentinels and error types are checked with errors.Is / errors.As.)

## (4) Mocks and Test Doubles

### 【MUST】TEST-015 Mock generation uses go.uber.org/mock.

- Category: Unit Testing/Mocks and Test Doubles
- Since: Go 1.0

Only `go.uber.org/mock` and its `mockgen` appear in go.mod and in the generated files, with no `github.com/golang/mock` introduced; existing old references are migrated to the uber version as well.

**Why**

> This project originates from Google's golang/mock repo. Unfortunately, Google no longer maintains this project, and given the heavy usage of gomock project within Uber, we've decided to fork and maintain this going forward at Uber.
>
> — https://github.com/uber-go/mock

`github.com/golang/mock` has been archived by Google and is no longer updated, no longer following Go releases; the fork Uber took over maintains API compatibility and is the version currently in use.

**Good**

```go
import "go.uber.org/mock/gomock"
```

**Bad**

```go
import "github.com/golang/mock/gomock"
```

**References**

- https://github.com/uber-go/mock

**Detection**: grep regex github\.com/golang/mock (Switch the matches to go.uber.org/mock; go.mod is out of scanning scope and needs manual checking.)

### 【MUST】TEST-016 Do not call ctrl.Finish() after passing *testing.T.

- Category: Unit Testing/Mocks and Test Doubles
- Since: Go 1.14

After `gomock.NewController(t)` is given a `*testing.T`, the controller validates expectations automatically when the test and its subtests finish, so `ctrl.Finish()` is no longer written in the code. When a self-built `TestReporter` is used, or what is passed in is not a `*testing.T`, you still have to ensure the validation timing yourself.

**Why**

> Note: If you pass a *testing.T into NewController, you no longer need to call ctrl.Finish() in your test methods.
>
> — https://pkg.go.dev/go.uber.org/mock/gomock#Controller.Finish

The `gomock` documentation states that after passing a `*testing.T` there is no need to call `Finish`; calling it by hand is the pre-Go 1.14 way, and the duplicate validation also cuts off early before the subtests have finished.

**Good**

```go
ctrl := gomock.NewController(t)
repo := NewMockRepository(ctrl)
```

**Bad**

```go
ctrl := gomock.NewController(t)
defer ctrl.Finish()
```

**References**

- https://pkg.go.dev/go.uber.org/mock/gomock#Controller.Finish

**Detection**: grep regex \.Finish\(\) (Confirm at each match whether a *testing.T was passed in.)

### 【MUST】TEST-017 Do not write redundant .Times(1).

- Category: Unit Testing/Mocks and Test Doubles
- Since: Go 1.0

A `gomock` expectation requires exactly one call by default, and `Times(1)` does not change the behavior. Write `Times(0)` only to guard against an unexpected call, or `Times(n)` when a specific multiple is required.

**Why**

> Times declares the exact number of times a function call is expected to be executed.
>
> — https://pkg.go.dev/go.uber.org/mock/gomock#Call.Times

Hanging a `Times(1)` at the end of every `EXPECT` forces readers to confirm one by one whether it is a deliberate constraint or added casually; the constraint already holds by default, so writing it out is only noise.

**Good**

```go
repo.EXPECT().Load(gomock.Any(), "a.json").Return(manual, nil)
```

**Bad**

```go
repo.EXPECT().Load(gomock.Any(), "a.json").Return(manual, nil).Times(1)
```

**References**

- https://pkg.go.dev/go.uber.org/mock/gomock#Call.Times
- https://github.com/uber-go/mock

**Detection**: grep regex \.Times\(1\) (Delete it at each match; one call is the default.)

### 【SHOULD】TEST-018 Use gomock.InOrder() only when call order is genuinely required.

- Category: Unit Testing/Mocks and Test Doubles
- Since: Go 1.0

By default `gomock` does not constrain the order in which expectations are called. Only in scenarios where the wrong order changes the result, such as lock, decrement, and release with their ordering dependencies, should you pin the order with `InOrder` or `After`; do not use it when the same set of expectations is called concurrently, to avoid occasional failures caused by varying scheduling order.

**Why**

> By default, expected calls are not enforced to run in any particular order. Call order dependency can be enforced by use of InOrder and/or Call.After.
>
> — https://pkg.go.dev/go.uber.org/mock/gomock#InOrder

Wrapping order-independent operations in `InOrder` makes the case turn red after a single refactor that changes no behavior; when multiple Goroutines call in, the order is decided by the scheduler, and the case goes red and green by turns.

**Good**

```go
repo.EXPECT().Load(gomock.Any(), "a.json").Return(manual, nil)
repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
```

**Bad**

```go
gomock.InOrder(
	repo.EXPECT().Load(gomock.Any(), "a.json").Return(manual, nil),
	repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil),
)
```

**References**

- https://pkg.go.dev/go.uber.org/mock/gomock#InOrder

**Detection**: grep regex gomock\.InOrder (Confirm at each match whether the call order affects the result.)

### 【SHOULD】TEST-019 Write argument values directly in expectations, avoiding wrapping in Eq and bare Any().

- Category: Unit Testing/Mocks and Test Doubles
- Since: Go 1.0

Arguments that are not matchers match by equality, so writing the literal directly suffices, with no extra layer of `gomock.Eq`; use `AssignableToTypeOf` to constrain the type, and `Cond`, `Len`, or `Nil` to match by condition. Use `Any()` only for values that are not asserted on, such as `context.Context`, and state the reason for relaxing the match there.

**Why**

> A Matcher is a representation of a class of values. It is used to represent the valid or expected arguments to a mocked method. ... Any returns a matcher that always matches.
>
> — https://pkg.go.dev/go.uber.org/mock/gomock#Matcher

A bare `Any()` only checks that the method was called, so a wrong argument still passes; writing the argument into the expectation lets argument errors at the call site be caught, and the case shows what this call should look like. An extra layer of `Eq` changes no matching semantics and only adds noise.

**Good**

```go
// first arg ctx is not asserted, matching is relaxed
repo.EXPECT().Load(gomock.Any(), "a.json").Return(manual, nil)
repo.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Manual{})).Return(nil)
```

**Bad**

```go
repo.EXPECT().Load(gomock.Any(), gomock.Eq("a.json")).Return(manual, nil)
repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
```

**References**

- https://pkg.go.dev/go.uber.org/mock/gomock#Matcher
- https://github.com/uber-go/mock/blob/main/gomock/call.go

**Detection**: grep regex gomock\.Any\(\) (Except for context.Context, switch the matches to Eq or AssignableToTypeOf.)

## (5) Using Coverage

### 【SHOULD】TEST-006 Test business logic and boundaries, and do not chase coverage numbers.

- Category: Unit Testing/Using Coverage
- Since: Go 1.0

Cases target business rules, branch boundaries, and error paths, and coverage numbers are not an acceptance condition. Utility types and pure algorithm types are required to reach 100% coverage, still achieved through external package tests, without accessing package-internal members because of the coverage requirement.

**Why**

> Use test coverage to find untested code. Coverage is no substitute for thought.
>
> — https://research.swtch.com/testing

Aiming at a coverage number produces cases that only call code without asserting, and tests written as copies of the implementation just to bring internals into coverage; the test itself is the explanation of that logic, and cases assembled to pad a number convey no business semantics, so the number rises while fewer real errors are caught.

**Good**

```go
func TestFilter(t *testing.T) {
	// business rule: filter by level combined with category
	got := manual.Filter([]string{"MUST"}, "Programming Conventions")
	want := []domain.Clause{{ID: "A-001", Level: domain.LevelMust, Category: "Programming Conventions/Naming"}}
	assert.Equal(t, want, got)
}
```

**Bad**

```go
func TestFilter(t *testing.T) {
	// only exercises the lines, without asserting the result
	manual.Filter([]string{"MUST"}, "")
	manual.Filter(nil, "Programming Conventions")
}
```

**References**

- https://research.swtch.com/testing

**Detection**: manual review (Check whether cases assert and whether they target business semantics; coverage numbers are not the basis for judgment.)

