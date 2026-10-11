# 1. Programming Conventions

## (1) Naming

### 【MUST】NAMING-001 Package names use lowercase words written together, without underscores, camel case, or plurals.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

The package name is the prefix that callers use to reference identifiers, so only all-lowercase with no separators stays consistent. Write multiple words together, without underscores or camel case, and do not add plurals.

**Why**

> In Go, package names must be concise and use only lowercase letters and numbers (e.g., k8s, oauth2). Multi-word package names should remain unbroken and in all lowercase (e.g., tabwriter instead of tabWriter, TabWriter, or tab_writer).
>
> — https://google.github.io/styleguide/go/decisions#package-names

The package name appears at every external reference point, so inconsistent naming keeps raising the reading cost; tools such as stylecheck have established conventions for package names, and following them lets deviations be detected automatically.

**Good**

```go
package orderbook
```

**Bad**

```go
package order_book
```

**References**

- https://go.dev/wiki/CodeReviewComments#package-names
- https://google.github.io/styleguide/go/decisions

**Detection**: golangci-lint stylecheck(ST1003) (Check the package name identifier)

### 【MUST】NAMING-002 Identifiers do not use underscores to separate words.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

Variables, constants, types, functions, and package names all use MixedCaps, without underscores separating words. There are only three exceptions:

- Package names used only by generated code
- Test, Benchmark, and Example function names in test files, where underscores may group cases
- Low-level libraries that interact with the operating system or cgo and reuse existing identifiers

**Why**

> Names in Go should in general not contain underscores. There are three exceptions to this principle:
>
> — https://google.github.io/styleguide/go/decisions#underscores

Underscore naming comes from the habits of other languages, while the Go convention is MixedCaps; mixing the two styles in one repository gives the same kind of identifier two forms, adding a step to searching and reading.

**Good**

```go
func readConfig() {}
```

**Bad**

```go
func read_config() {}
```

**References**

- https://google.github.io/styleguide/go/decisions#underscores
- https://go.dev/wiki/CodeReviewComments#mixed-caps

**Detection**: golangci-lint stylecheck(ST1003) (Check the identifier naming style)

### 【MUST】NAMING-003 Initialisms keep a consistent case within identifiers.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

Initialisms such as URL, ID, and HTTP are either all uppercase or all lowercase, never written as Url or Id. The case follows how the initialism is written in English (such as iOS and gRPC), and the first letter is uppercased when the identifier must be exported.

**Why**

> Words in names that are initialisms or acronyms (e.g., URL and NATO) should have the same case. URL should appear as URL or url (as in urlPony, or URLPony), never as Url.
>
> — https://google.github.io/styleguide/go/decisions#initialisms

Mixed case makes one concept appear under several spellings such as URL, Url, and url, so searching and replacing by name both miss hits; tools have established rules for this kind of naming, and following the convention lets them be detected automatically.

**Good**

```go
func ServeHTTP(w http.ResponseWriter, r *http.Request) {}
```

**Bad**

```go
func ServeHttp(w http.ResponseWriter, r *http.Request) {}
```

**References**

- https://go.dev/wiki/CodeReviewComments#initialisms
- https://google.github.io/styleguide/go/decisions#initialisms

**Detection**: golangci-lint stylecheck(ST1003)

### 【SHOULD】NAMING-004 The length of a name is proportional to the size of its scope.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

A loop index or method receiver used within one or two lines needs only a single letter; a variable that spans an entire function or package-level scope uses a full word. When several similar concepts live in the same scope, add a qualifier to tell them apart.

**Why**

> The general rule of thumb is that the length of a name should be proportional to the size of its scope and inversely proportional to the number of times that it is used within that scope.
>
> — https://google.github.io/styleguide/go/decisions#variable-names

A short name saves noise in a short scope, but in a long scope it makes readers forget what it refers to by the time they reach the end; sizing the name to the scope keeps readers from having to scroll back to the declaration.

**Good**

```go
orderCount := 0
for _, order := range orders {
	if order.Paid {
		orderCount++
	}
}
```

**Bad**

```go
c := 0
for _, o := range orders {
	if o.Paid {
		c++
	}
}
```

**References**

- https://google.github.io/styleguide/go/decisions#variable-names
- https://go.dev/wiki/CodeReviewComments#variable-names

**Detection**: manual review (Decide by hand whether name length matches the scope)

### 【SHOULD】NAMING-005 Names do not repeat the type or the surrounding context.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

The declaration already gives the type, and the package name, method name, and file name already provide context, so the name does not repeat them; add a qualifier only when the type genuinely needs disambiguation.

- Use users, not userSlice or numUsers
- In the report package use Report, not ReportReport
- In a UserCount method use count, not userCount

**Why**

> Names that include information from their surrounding context often create extra noise without benefit. The package name, method name, type name, function name, import path, and even filename can all provide context that automatically qualifies all names within.
>
> — https://google.github.io/styleguide/go/decisions#external-context-vs-local-names

The repeated type and context add no discriminating power; they only lengthen the name and the code lines, and more lines wrap as a result. Once the redundancy is removed, the name still uniquely points to the intended value.

**Good**

```go
var users []User
```

**Bad**

```go
var userSlice []User
```

**References**

- https://google.github.io/styleguide/go/decisions#variable-name-vs-type
- https://google.github.io/styleguide/go/decisions#external-context-vs-local-names

**Detection**: manual review (Check whether the name repeats the type or the surrounding context)

### 【SHOULD】NAMING-006 Receiver names use a short abbreviation of the type name and stay consistent within the same type.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

A receiver name is usually one or two letters taken from an abbreviation of the type name, and every method of the same type uses the same receiver name. Do not use this or self, and do not leave a placeholder name when the receiver is unused.

**Why**

> The name of a method's receiver should be a reflection of its identity; often a one or two letter abbreviation of its type suffices (such as "c" or "cl" for "Client"). Don't use generic names such as "me", "this" or "self". Be consistent, too: if you call the receiver "c" in one method, don't call it "cl" in another.
>
> — https://go.dev/wiki/CodeReviewComments#receiver-names

The methods of one type are spread across several files, and using a different receiver name in each makes readers think they refer to different objects; this and self are habits carried over from other languages and are not the Go convention.

**Good**

```go
func (w *ReportWriter) Write(p []byte) (int, error) {}
```

**Bad**

```go
func (this *ReportWriter) Write(p []byte) (int, error) {}
```

**References**

- https://google.github.io/styleguide/go/decisions#receiver-names
- https://go.dev/wiki/CodeReviewComments#receiver-names

**Detection**: manual review (Check whether receiver names are consistent within the same type and whether this / self is used)

### 【SHOULD】NAMING-007 Boolean predicate methods are named with an Is, Has, or Can prefix.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

Methods returning bool pick a prefix by semantics so the name reads as a yes/no question.

- State checks use Is: IsHealthy, IsEmpty
- Possession checks use Has: HasExpired, HasPrefix
- Capability checks use Can: CanRetry, CanClose
- Existing stdlib predicates without a prefix (such as `strings.Contains`) keep their convention

**Why**

> it's neither idiomatic nor necessary to put Get into the getter's name
>
> — https://go.dev/doc/effective_go#Getters

The standard library's position on naming methods is to let the name carry the meaning directly, and predicates are no exception: IsConnected reads as a yes/no question at the call site, while the bare adjective Connected cannot be told apart from a connection object, so readers must jump to the definition to confirm it is a bool.

**Good**

```go
func (s *Server) IsHealthy() bool { return s.err == nil }
```

**Bad**

```go
func (s *Server) Healthy() bool { return s.err == nil }
```

**References**

- https://go.dev/doc/effective_go#Getters

**Detection**: manual review (Check that methods returning bool carry a predicate prefix)

### 【SHOULD】NAMING-008 Boolean variables and struct fields are named with an Is, Has, or Can prefix.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

Variables and fields of type bool use a linking-verb prefix to express a judgment.

- IsConnected, HasCache, CanRetry
- No bare adjectives or nouns: Connected, Cache, Retry
- Temporary boolean variables follow the same rule: isValid, hasMore

**Why**

> it's neither idiomatic nor necessary to put Get into the getter's name
>
> — https://go.dev/doc/effective_go#Getters

Bare adjectives collide with object names: cfg.Cache could mean a cache configuration or a boolean that enables caching; adding Is/Has/Can makes the type semantics self-evident in the name, so reading an if condition no longer requires a jump to the definition.

**Good**

```go
type Config struct { IsVerbose bool }
```

**Bad**

```go
type Config struct { Verbose bool }
```

**References**

- https://go.dev/doc/effective_go#Getters

**Detection**: manual review (Check that bool fields and variables carry a prefix)

### 【SHOULD】NAMING-009 Single-method interfaces are named after the method plus an er suffix.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

A single-method interface forms an agent noun from its method name plus er.

- Reader, Writer, Closer, Stringer, CloseNotifier
- Multi-method interfaces use concrete domain nouns: Conn, Pool, Store
- When er does not fit the method name, a close variant is allowed, such as Recognizer for Recognize

**Why**

> By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier.
>
> — https://go.dev/wiki/CodeReviewComments#interface-names

The er suffix ties the interface name directly to its single method, so the contract is visible from the name; the standard library follows this convention throughout, and deviating from it produces two naming styles within one family of interfaces.

**Good**

```go
type Reader interface { Read(p []byte) (n int, err error) }
```

**Bad**

```go
type DataFetcher interface { Read(p []byte) (n int, err error) }
```

**References**

- https://go.dev/wiki/CodeReviewComments#interface-names

**Detection**: manual review (Check that single-method interface type names correspond to their method)

### 【SHOULD】NAMING-010 Function variants use fixed markers such as Must, f, and WithContext.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

Variants of one function distinguish behavior through fixed prefixes and suffixes.

- Must prefix: panics on failure, as in `regexp.MustCompile`
- f suffix: takes a format string and variadic args, as in `errors.Errorf`
- WithContext suffix: carries a context, as in `http.NewRequestWithContext`
- Markers stay consistent with the standard library; do not invent new ones

**Why**

> MustCompile is like Compile but panics if the expression cannot be parsed. It simplifies safe initialization of global variables holding compiled regular expressions.
>
> — https://pkg.go.dev/regexp#MustCompile

Callers in code review and IDE hints often see only the function name, not the signature: the Must prefix says this panics, the f suffix says the argument is a format string, WithContext says cancellation propagates through ctx; if variant behavior is not in the name, misuse surfaces only when the implementation is read.

**Good**

```go
var re = regexp.MustCompile(`^v\d+$`)
```

**Bad**

```go
func MustParseTime(s string) time.Time { t, _ := time.Parse(time.RFC3339, s); return t }
```

**References**

- https://pkg.go.dev/regexp#MustCompile

**Detection**: manual review (Check that variant functions keep naming markers consistent with the standard library)

### 【SHOULD】NAMING-011 Identifiers do not shadow predeclared names such as len, cap, error, and min.

- Category: Programming Conventions/Naming
- Since: Go 1.0
- Tags: naming

Local variables, parameters, and fields do not reuse predeclared identifiers.

- Predeclared names include len, cap, new, make, min, max, clear, error, any
- Misusing a shadowed name produces no compile error; the semantics change silently
- When a same-name declaration is truly needed (for example generated code), state the reason in review

**Why**

> The following identifiers are implicitly declared in the universe block:
>
> — https://go.dev/ref/spec#Predeclared_identifiers

Predeclared identifiers are the fallback names of every scope; shadowing them locally turns min(x, y) into a user function call and error into a local variable, with no compiler hint — only unusual behavior. The predeclared checker turns this into a mechanical check.

**Good**

```go
func dedup(items []int) []int { maxSeen := 0; _ = maxSeen; return items }
```

**Bad**

```go
func dedup(items []int) []int { max := 0; _ = max; return items }
```

**References**

- https://go.dev/ref/spec#Predeclared_identifiers

**Detection**: golangci-lint predeclared (Check that local identifiers do not shadow predeclared names)

## (2) Constants and Enums

### 【MUST】CONST-001 Constant names use MixedCaps, not all uppercase or a k prefix.

- Category: Programming Conventions/Constants and Enums
- Since: Go 1.0
- Tags: naming

An exported constant starts with an uppercase letter and an unexported one with a lowercase letter, with each word capitalized in between. Do not use all-uppercase forms such as MAX_LENGTH, and do not use the k prefix of kMaxLength.

**Why**

> Constant names must use MixedCaps like all other names in Go. ... Do not use non-MixedCaps constant names or constants with a K prefix.
>
> — https://google.github.io/styleguide/go/decisions#constant-names

All-uppercase names and the k prefix come from other languages, while every identifier in Go shares one MixedCaps convention; giving constants their own style puts two naming styles in the same repository.

**Good**

```go
const MaxPacketSize = 512
```

**Bad**

```go
const MAX_PACKET_SIZE = 512
```

**References**

- https://go.dev/wiki/CodeReviewComments#mixed-caps
- https://google.github.io/styleguide/go/decisions#constant-names

**Detection**: golangci-lint stylecheck(ST1003)

### 【SHOULD】CONST-002 Constants are named by their role, not by their value.

- Category: Programming Conventions/Constants and Enums
- Since: Go 1.0
- Tags: naming

A constant name says what it stands for, not what it equals; when a value has no meaning beyond the number itself, there is no need to define it as a constant. When the same value serves different roles, name each one separately.

**Why**

> Name constants based on their role, not their values. If a constant does not have a role apart from its value, then it is unnecessary to define it as a constant.
>
> — https://google.github.io/styleguide/go/decisions#constant-names

Naming by value (Twelve = 12, ThreeSeconds = 3) makes the name wrong once the value changes, and it says nothing about the purpose at the call site; naming by role is what carries the intent at the call site.

**Good**

```go
const maxRetries = 3
```

**Bad**

```go
const three = 3
```

**References**

- https://google.github.io/styleguide/go/decisions#constant-names

**Detection**: manual review (Check whether the constant name merely restates the value)

### 【SHOULD】CONST-004 Leave the zero value of an iota enum as an Unknown sentinel, or start at iota+1.

- Category: Programming Conventions/Constants and Enums
- Since: Go 1.0
- Tags: constants, naming

The zero value of an enum type carries no business meaning.

- Define StatusUnknown first, or start the first business member at iota+1
- Zero serves as the uninitialized/invalid sentinel so validation can separate it from real states
- Exception: types whose zero value has clear business meaning (such as GenderMale at 0) may keep it, confirmed in review

**Why**

> Within a const declaration a name not explicitly initialized is implicitly initialized to iota, which is the line number of the corresponding ConstSpec.
>
> — https://go.dev/ref/spec#Iota

The zero value of var s Status silently equals the first const member, so an uninitialized field passes validation as the first business state, and debugging cannot tell a default from a real state; reserving zero for Unknown makes uninitialized states surface explicitly at validation.

**Good**

```go
type Status int
const (
  StatusUnknown Status = iota
  StatusActive
  StatusClosed
)
```

**Bad**

```go
type Status int
const (
  StatusActive Status = iota
  StatusClosed
)
```

**References**

- https://go.dev/ref/spec#Iota

**Detection**: manual review (Check whether the first iota member occupies the zero value)

### 【MAY】CONST-003 Enumeration values that increase from one source are generated with iota.

- Category: Programming Conventions/Constants and Enums
- Since: Go 1.0
- Tags: constants

A set of enumerations from one source that increase in order is generated with iota, avoiding hand-written numbers that are forgotten when a value is inserted midway. The value of iota depends on its line number within the const block, so inserting a line into the block means rechecking the following values.

**Why**

> In Go, enumerated constants are created using the iota enumerator. Since iota can be part of an expression and expressions can be implicitly repeated, it is easy to build intricate sets of values.
>
> — https://go.dev/doc/effective_go#constants

With hand-written numbers, inserting an enumeration value in the middle means renumbering every later entry, and missing one produces a duplicate; iota generates values by line, so an inserted line never conflicts with existing values.

**Good**

```go
const (
	StatusPending = iota
	StatusRunning
	StatusDone
)
```

**Bad**

```go
const (
	StatusPending = 0
	StatusRunning = 1
	StatusDone    = 1
)
```

**References**

- https://go.dev/doc/effective_go#constants
- https://go.dev/ref/spec#Constant_declarations

**Detection**: manual review (Check whether enumerations that increase from one source use hand-written numbers)

## (3) Formatting and Style

### 【SHOULD】STYLE-001 Use any instead of interface{}.

- Category: Programming Conventions/Formatting and Style
- Since: Go 1.18
- Tags: modernization

Go 1.18 introduced any as an alias for interface{}, and the two are entirely equivalent. New code writes any uniformly; when touching an interface{} in old code, replace it along the way.

**Why**

> The new predeclared identifier any is an alias for the empty interface. It may be used instead of interface{}.
>
> — https://go.dev/doc/go1.18

any is the officially recommended form, its meaning is direct, and it matches the form used in generic constraints.

**Good**

```go
func printValue(v any) {
	fmt.Println(v)
}
```

**Bad**

```go
func printValue(v interface{}) {
	fmt.Println(v)
}
```

**References**

- https://go.dev/doc/go1.18
- https://google.github.io/styleguide/go/decisions

**Detection**: golangci-lint modernize(any) (The same analyzer can rewrite it directly with go fix)

## (4) Functions and Methods

### 【MUST】FUNC-001 Exported functions and methods come at the front of a file, unexported ones at the back.

- Category: Programming Conventions/Functions and Methods
- Since: Go 1.0
- Tags: code-organization

Within one file, exported functions come before unexported functions, and the exported methods of a type come before its unexported methods. This clause governs only the relative order of exported and unexported, not where types and constructors go.

**Why**

> Functions should be sorted in rough call order. Functions in a file should be grouped by receiver. Therefore, exported functions should appear first in a file, after struct, const, var definitions.
>
> — https://github.com/uber-go/guide/blob/master/style.md#function-grouping-and-ordering

Reading a file from top to bottom, what the package offers comes first and implementation details unfold only when needed; when the public API sits between private helpers, readers must skip past implementations they do not care about to find the entry point. With a fixed order, both review and tooling can decide it directly.

**Good**

```go
func Load(path string) (*Manual, error) {
	return readJSON(path)
}

func readJSON(path string) (*Manual, error) {
	// ...
}
```

**Bad**

```go
func readJSON(path string) (*Manual, error) {
	// ...
}

func Load(path string) (*Manual, error) {
	return readJSON(path)
}
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#function-grouping-and-ordering
- https://github.com/manuelarte/funcorder

**Detection**: golangci-lint funcorder

### 【MUST】FUNC-004 Nil branches of interface-returning functions return nil explicitly.

- Category: Programming Conventions/Functions and Methods
- Since: Go 1.0
- Tags: error-handling

When the return type is an interface (error included), do not return a typed nil pointer variable.

- Return nil directly when there is no error
- When a concrete error type is needed, construct it and decide the return value by emptiness
- Applies to every interface return value; error is the most common case

**Why**

> If we store a nil pointer of type *int inside an interface value, the inner type will be *int regardless of the value of the pointer: (*int)(nil). Thus, the interface value will be non-nil even when the pointer value is inside.
>
> — https://go.dev/doc/faq#nil_error

An interface value carries a (type, value) pair; a typed nil pointer wrapped in it makes the interface non-nil, so the caller's err != nil check holds and enters the error branch, panicking later when nil fields are dereferenced. Returning nil explicitly removes the trap at the source.

**Good**

```go
func validate(c Config) error {
  if c.Host == "" {
    return ErrMissingHost
  }
  return nil
}
```

**Bad**

```go
func validate(c Config) error {
  var vErr *ValidationError
  if c.Host == "" {
    vErr = &ValidationError{Field: "host"}
  }
  return vErr
}
```

**References**

- https://go.dev/doc/faq#nil_error

**Detection**: manual review (Check the nil branches of interface-returning functions)

### 【SHOULD】FUNC-002 Use named results when the meaning of the return values is not self-evident.

- Category: Programming Conventions/Functions and Methods
- Since: Go 1.0
- Tags: naming

When an unnamed result already conveys the meaning, do not add a name; when several values of the same type are returned, or the meaning of a result cannot be seen from the signature, name the result parameters so that the signature and godoc explain themselves. The test is whether the meaning is clear, independent of function length.

- Name them when returning two float64 values: `(lat, long float64, err error)`
- Do not name them when the type already conveys the meaning: `(*Manual, error)`

**Why**

> On the other hand, if a function returns two or three parameters of the same type, or if the meaning of a result isn't clear from context, adding names may be useful in some contexts. ... Clarity of docs is always more important than saving a line or two in your function.
>
> — https://go.dev/wiki/CodeReviewComments#named-result-parameters

When two values of the same type are returned, the signature alone cannot tell which is which, so readers must dig into the implementation or documentation; naming lets the signature state the meaning itself. Naming a result whose meaning is already clear only adds noise to godoc.

**Good**

```go
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return x, y
}
```

**Bad**

```go
func split(sum int) (int, int) {
	x := sum * 4 / 9
	y := sum - x
	return x, y
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#named-result-parameters

**Detection**: manual review (Check whether results are named when their meaning is not self-evident)

### 【SHOULD】FUNC-003 Collect more than 4 arguments or more than 3 results into a struct.

- Category: Programming Conventions/Functions and Methods
- Since: Go 1.0
- Tags: api-design

When a function takes 5 or more arguments, or returns 4 or more results, gather the values that travel together into a struct and let callers pass them by field name. The counting and the way to collect them:

- Counting: arguments are counted by name, `ctx` included and the method receiver excluded; results are counted by name, `error` included.
- Arguments: business arguments go into one argument struct, while `ctx` stays first in the signature and out of the struct, consistent with `CTX-001`.
- Results: the results go into one result struct, while `error` stays in the last result position and out of the struct.
- Exceptions: a signature fixed by an external contract is left as is, such as interface implementations, `go:generate` output, and mock output.
- Division of labor with `FUNC-002`: name the results when their meaning is not self-evident; collect them into a struct only when the count exceeds the limit.

**Why**

> Warns when a function receives more parameters than the maximum set by the rule's configuration. ... Specifies the maximum number of results a function can return. Functions returning too many results can be hard to understand/use.
>
> — https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md

The two rules quoted above state configurable limits; Go itself only says to keep the number of parameters reasonable and gives no figure. The numbers come from tooling and the community: revive defaults `argument-limit` to 8 and `function-result-limit` to 3, and go-critic defaults `tooManyResultsChecker` to 5 with that check off by default. This manual sets its own limits of 4 arguments and 3 results. The reason is readability at the call site: by the fifth positional argument, readers can no longer tell which value carries which meaning and must scroll back to the signature to match them one by one; the same happens with a fourth result. Collected into a struct, every value carries a field name, the call site explains itself, and adding or removing a field later cannot silently shift positions.

**Good**

```go
func (s *Service) List(ctx context.Context, q ListQuery) ([]ClauseView, error) {
	return s.repo.Find(ctx, q)
}
```

**Bad**

```go
func (s *Service) List(ctx context.Context, lang string, levels []string, category string, tags []string) ([]ClauseView, error) {
	return s.repo.Find(ctx, lang, levels, category, tags)
}
```

**References**

- https://github.com/mgechev/revive/blob/master/RULES_DESCRIPTIONS.md
- https://go-critic.com/overview.html#toomanyresultschecker

**Detection**: golangci-lint revive(argument-limit, function-result-limit)

## (5) Data Structures

### 【MUST】DATA-003 Type assertions use the comma-ok form, not the single-value form.

- Category: Programming Conventions/Data Structures
- Since: Go 1.0
- Tags: data-structures

When the assertion may not hold, take the ok value and branch.

- Use `v, ok := x.(T)` and branch on ok
- The single-value form `v := x.(T)` is limited to cases where a failed assertion is a program error, with a comment
- Type switches are not covered by this clause

**Why**

> If it does not hold, a run-time panic occurs.
>
> — https://go.dev/ref/spec#Type_assertions

The single-value form panics on a type mismatch, turning the predictable branch of a wrong type into a process crash; comma-ok leaves the decision to the code, letting the caller take a downgrade or error path.

**Good**

```go
v, ok := x.(string)
if !ok {
  return fmt.Errorf("unexpected type %T", x)
}
```

**Bad**

```go
v := x.(string)
```

**References**

- https://go.dev/ref/spec#Type_assertions

**Detection**: golangci-lint forcetypeassert (Check the form of every type assertion)

### 【MUST】DATA-004 Initialize a nil map before writing to it.

- Category: Programming Conventions/Data Structures
- Since: Go 1.0
- Tags: data-structures

Confirm the map has been initialized with make or a literal before writing.

- Reading a nil map is legal and returns zero values
- Writing to a nil map panics
- Maps inside structs can be lazily initialized in the constructor or a method so the zero value stays usable

**Why**

> A nil map is equivalent to an empty map except no elements may be added.
>
> — https://go.dev/blog/maps

The zero value of `var m map[string]int` can be read but not written; any write path that reaches it panics with no compile-time hint at all. Initialize up front, or make the type's zero value usable — pick one and be explicit.

**Good**

```go
m := make(map[string]int)
m["k"] = 1
```

**Bad**

```go
var m map[string]int
m["k"] = 1
```

**References**

- https://go.dev/blog/maps

**Detection**: manual review (Check that map write paths operate on initialized maps)

### 【MUST】DATA-005 Do not rely on map iteration order; sort keys when order matters.

- Category: Programming Conventions/Data Structures
- Since: Go 1.0
- Tags: data-structures

The output order of a range over a map is not a contract.

- Serialization, comparison, and display that need stable output sort the keys first
- Order-dependent logic fails intermittently in CI and production alike
- The runtime randomizes map iteration; this is a language guarantee, not an implementation detail

**Why**

> The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next.
>
> — https://go.dev/ref/spec#For_range

The language specification explicitly makes no promise about iteration order, so order-dependent code can differ between two runs of the same process, and a passing test only means randomness favored it; explicit sorting is the only way to establish a contract.

**Good**

```go
keys := slices.Sorted(maps.Keys(m))
for _, k := range keys {
  fmt.Println(k, m[k])
}
```

**Bad**

```go
for k, v := range m {
  fmt.Println(k, v)
}
```

**References**

- https://go.dev/ref/spec#For_range

**Detection**: manual review (Check that logic after ranging a map does not depend on order)

### 【MUST】DATA-007 Struct fields that participate in JSON serialization must be exported.

- Category: Programming Conventions/Data Structures
- Since: Go 1.0
- Tags: data-structures

Fields that must appear in JSON use exported identifiers; unexported fields are silently ignored.

- Unexported fields cannot be encoded or decoded by `encoding/json`
- Pair with a field tag for the key mapping (see STRUCT-004)
- Numbers parsed into `interface{}` are always `float64`; large integers go through `json.Number` or typed structs

**Why**

> Unmarshal matches incoming object keys to the keys used by Marshal (either the struct field name or its tag), preferring an exact match but also accepting a case-insensitive match.
>
> — https://pkg.go.dev/encoding/json#Unmarshal

encoding/json cannot access unexported fields, so production JSON silently loses fields — and local tests that never verify the contract will not catch it; exporting with a tag pins the contract at the struct definition.

**Good**

```go
type User struct {
  Name string `json:"name"`
}
```

**Bad**

```go
type User struct {
  name string `json:"name"`
}
```

**References**

- https://pkg.go.dev/encoding/json#Unmarshal

**Detection**: manual review (Compare struct fields with the serialization contract and check unexported fields)

### 【SHOULD】DATA-001 Declare empty slices with var to get a nil slice.

- Category: Programming Conventions/Data Structures
- Since: Go 1.0
- Tags: data-structures

Empty slices are uniformly declared with var, and a non-nil empty slice is used only when genuinely needed: in JSON serialization a nil slice outputs null while an empty slice outputs []. In interface design, do not treat a nil slice and a non-nil empty slice as two states.

**Why**

> When declaring an empty slice, prefer var t []string over t := []string{}. The former declares a nil slice value, while the latter is non-nil but zero-length.
>
> — https://go.dev/wiki/CodeReviewComments#declaring-empty-slices

Both have len and cap of 0, so they are functionally equivalent yet written two ways; settling on nil means readers need not stop to work out why the other form was used.

**Good**

```go
var ids []string
```

**Bad**

```go
ids := []string{}
```

**References**

- https://go.dev/wiki/CodeReviewComments#declaring-empty-slices

**Detection**: manual review (Outputting a JSON array [] is an exception)

### 【SHOULD】DATA-002 Struct literals are initialized with field names rather than listed by position.

- Category: Programming Conventions/Data Structures
- Since: Go 1.0
- Tags: data-structures

A literal writes field name: value for each field. Small private structs with few fields and a stable order, and test table data, may be exceptions; a struct from another package always uses field names, because adding, removing, or reordering fields silently misplaces positional entries.

**Why**

> You should almost always specify field names when initializing structs. This is now enforced by go vet.
>
> — https://github.com/uber-go/guide/blob/master/style.md#use-field-names-to-initialize-structs

The positional form turns field order into an implicit contract, so after the other side adjusts the struct, positions of the same type are silently mismatched with no compile error; the field name form binds each value to its field so the change surfaces at once. The composites check in `go vet` covers positional literals from another package.

**Good**

```go
user := User{
	Name: "alice",
	Age:  30,
}
```

**Bad**

```go
user := User{"alice", 30}
```

**References**

- https://pkg.go.dev/cmd/vet
- https://github.com/uber-go/guide/blob/master/style.md#use-field-names-to-initialize-structs

**Detection**: golangci-lint govet(composites) (govet covers only positional literals from another package, and the composites check must be enabled)

### 【SHOULD】DATA-006 Check the value range before a narrowing integer conversion.

- Category: Programming Conventions/Data Structures
- Since: Go 1.0
- Tags: data-structures, security

Narrowing conversions such as int64 to int32 or uint to int need explicit validation.

- Verify the value fits the target type before converting
- Conversions driven by external input must be checked (the `gosec G115` case)
- Bit-level scenarios with clearly meaningless high bits may be waived with a comment

**Why**

> When converting between integer types, if the value is a signed integer, it is sign extended to infinite precision and then truncated to fit in the result type's size; if the value is an unsigned integer, it is zero-extended and then truncated to fit in the result type's size.
>
> — https://go.dev/ref/spec#Conversions_to_and_from_integer_string_types

Narrowing conversions truncate and wrap silently: three billion in int64 becomes a negative int32 with no error anywhere; when external input drives the conversion, the truncated result becomes a business bug or a security hole.

**Good**

```go
if v > math.MaxInt32 {
  return fmt.Errorf("value %d overflows int32", v)
}
n := int32(v)
```

**Bad**

```go
n := int32(v)
```

**References**

- https://go.dev/ref/spec#Conversions_to_and_from_integer_string_types

**Detection**: golangci-lint gosec(G115) (Check the range before narrowing conversions)

### 【MAY】MODERN-001 Common operations on slices and maps prefer the slices and maps standard library packages over hand-written loops.

- Category: Programming Conventions/Data Structures
- Since: Go 1.21
- Tags: modernization

`slices` and `maps` provide high-frequency operations such as `Contains`, `Sort`, `Clone`, and `Keys`. When the standard library can express the operation, stop writing loops by hand, which reduces boundary mistakes and unifies the style.

**Why**

> The new slices package provides many common operations on slices, using generic functions that work with slices of any element type.
>
> — https://go.dev/doc/go1.21

Hand-written loops tend to go wrong on empty slices, duplicate elements, and sort stability, and the standard library already covers these boundaries.

**Good**

```go
if slices.Contains(ids, target) {
	return true
}
```

**Bad**

```go
for _, id := range ids {
	if id == target {
		return true
	}
}
```

**References**

- https://pkg.go.dev/slices
- https://go.dev/doc/go1.21

**Detection**: golangci-lint modernize (go fix also suggests the equivalent modern form)

## (6) Standard Library Usage

### 【MUST】STD-001 time.Duration is never multiplied by time.Duration.

- Category: Programming Conventions/Standard Library Usage
- Since: Go 1.0
- Tags: data-structures

Duration arithmetic multiplies only by scalars.

- Use `d * n` (n an integer) for multiples
- `d * d` and `d * (x * time.Second)` are unit errors
- Build durations from variables with scalar multiplication such as `time.Duration(n) * time.Second`

**Why**

> A Duration represents the elapsed time between two instants as an int64 nanosecond count.
>
> — https://pkg.go.dev/time#Duration

A Duration is a nanosecond count, so multiplying two durations yields square nanoseconds, off by many orders of magnitude from intuition: `5 * time.Second * time.Second` is not 25 seconds but 2.5 trillion nanoseconds; the durationcheck checker turns this into a mechanical check.

**Good**

```go
timeout := 5 * time.Second
result := timeout * attempts
```

**Bad**

```go
timeout := 5 * time.Second * time.Second
```

**References**

- https://pkg.go.dev/time#Duration

**Detection**: golangci-lint durationcheck (Check every multiplication involving time.Duration)

### 【SHOULD】STD-002 Compare time.Time values with Equal, not ==.

- Category: Programming Conventions/Standard Library Usage
- Since: Go 1.0
- Tags: data-structures

Use Equal or Compare to test for the same instant.

- `==` also compares wall clock, location, and the monotonic clock reading
- Times parsed from strings or received over the network carry no monotonic reading, but mixed sources make `==` unreliable
- Zero checks use `IsZero`; ordering uses `Before`/`After`

**Why**

> Equal reports whether t and u represent the same time instant. Two times can be equal even if they are in different locations.
>
> — https://pkg.go.dev/time#Time.Equal

time.Time carries a monotonic clock reading internally, and == folds that reading into the comparison, so the same wall-clock instant can compare unequal; Equal compares only the time instant, matching business expectations.

**Good**

```go
if t1.Equal(t2) {
  return true
}
```

**Bad**

```go
if t1 == t2 {
  return true
}
```

**References**

- https://pkg.go.dev/time#Time.Equal

**Detection**: manual review (Check the expressions used to compare time.Time values)

### 【SHOULD】STD-003 Strip prefixes and suffixes with strings.TrimPrefix and TrimSuffix.

- Category: Programming Conventions/Standard Library Usage
- Since: Go 1.0
- Tags: data-structures

Do not use the strings.Trim family for substring trimming.

- `strings.TrimPrefix`/`strings.TrimSuffix` remove one exact substring: `strings.TrimPrefix(s, "api/")`
- `strings.Trim`/`strings.TrimLeft`/`strings.TrimRight` take a set of characters and repeatedly strip all matching ones
- Mixing the two is a common source of over-trimming

**Why**

> Trim returns a slice of the string s with all leading and trailing Unicode code points contained in cutset removed.
>
> — https://pkg.go.dev/strings#Trim

Trim's second argument is a set of characters, not a substring: strings.Trim(s, "ab") removes every leading and trailing a and b, which differs entirely from the intent of removing the prefix ab; the set semantics hide behind the function name, and misuse shows up as occasional over-trimming.

**Good**

```go
path := strings.TrimPrefix(raw, "/api/")
```

**Bad**

```go
path := strings.Trim(raw, "/api/")
```

**References**

- https://pkg.go.dev/strings#Trim

**Detection**: manual review (Check that the trimming intent matches the API used)

## (7) Database Access

### 【MUST】DB-001 Check rows.Err and Close rows after iterating sql.Rows.

- Category: Programming Conventions/Database Access
- Since: Go 1.0
- Tags: data-structures, error-handling

The rows lifecycle spans the whole iteration.

- Check `rows.Err()` after the `for rows.Scan` loop ends
- Close rows when done, or `defer rows.Close()` right after the query
- Statements that return no rows (INSERT/UPDATE/DELETE) use `Exec`, not `Query`

**Why**

> Err returns the error, if any, that was encountered during iteration. Err should be checked explicitly after the iteration is done.
>
> — https://pkg.go.dev/database/sql#Rows.Err

Without Close, a connection stays occupied and is never returned to the pool; without Err, an I/O interruption mid-iteration is swallowed (on a network blip the loop quietly reads fewer rows); Err's error never escapes through Scan, so it must be checked explicitly.

**Good**

```go
rows, err := db.QueryContext(ctx, q)
if err != nil {
  return err
}
defer rows.Close()
for rows.Next() {
  if err := rows.Scan(&id); err != nil {
    return err
  }
}
return rows.Err()
```

**Bad**

```go
rows, err := db.QueryContext(ctx, q)
if err != nil {
  return err
}
for rows.Next() {
  if err := rows.Scan(&id); err != nil {
    return err
  }
}
return nil
```

**References**

- https://pkg.go.dev/database/sql#Rows.Err

**Detection**: golangci-lint rowserrcheck, sqlclosecheck (Check that rows.Err and rows.Close appear after iteration)

### 【MUST】DB-002 Wrap multi-statement writes in a transaction and defer Rollback right after Begin.

- Category: Programming Conventions/Database Access
- Since: Go 1.0
- Tags: data-structures, error-handling

Transactions need a single cleanup path.

- Follow `tx := db.BeginTx(ctx, nil)` immediately with `defer tx.Rollback()`
- Rollback after a successful Commit returns `ErrTxDone`, which is harmless
- Panic or early-return paths leave no transaction hanging

**Why**

> After calling Commit or Rollback, all operations on the transaction fail with ErrTxDone.
>
> — https://pkg.go.dev/database/sql#Tx

Without defer Rollback, a transaction hangs until connection reclamation on panic or early return, holding row locks the whole time; calling Rollback after Commit only returns ErrTxDone, so as a safety net it is always safe.

**Good**

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
  return err
}
defer tx.Rollback()
if _, err := tx.ExecContext(ctx, q1); err != nil {
  return err
}
return tx.Commit()
```

**Bad**

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
  return err
}
if _, err := tx.ExecContext(ctx, q1); err != nil {
  return err
}
return tx.Commit()
```

**References**

- https://pkg.go.dev/database/sql#Tx

**Detection**: manual review (Check that transactions defer Rollback right after Begin)

### 【MUST】DB-004 Match sql.ErrNoRows with errors.Is.

- Category: Programming Conventions/Database Access
- Since: Go 1.0
- Tags: error-handling

Not found is a domain result, not an ordinary error.

- Match with `errors.Is(err, sql.ErrNoRows)` and translate into domain semantics (such as a not-found marker)
- Direct `==` comparison breaks once the error is wrapped
- `QueryRow`'s single-row queries are the main source of ErrNoRows

**Why**

> ErrNoRows is returned by Scan when QueryRow doesn't return a row. In this case, QueryRow will return this error.
>
> — https://pkg.go.dev/database/sql#ErrNoRows

Treating ErrNoRows as an ordinary error leaves upstream unable to tell no such row from a database outage, which often shows up as a 5xx for a request that simply did not exist; matching with errors.Is and translating into a not-found result completes the error-chain semantics.

**Good**

```go
err := db.QueryRowContext(ctx, q, id).Scan(&name)
if errors.Is(err, sql.ErrNoRows) {
  return ErrNotFound
}
if err != nil {
  return err
}
```

**Bad**

```go
err := db.QueryRowContext(ctx, q, id).Scan(&name)
if err == sql.ErrNoRows {
  return ErrNotFound
}
if err != nil {
  return err
}
```

**References**

- https://pkg.go.dev/database/sql#ErrNoRows

**Detection**: manual review (Check how ErrNoRows is matched)

### 【SHOULD】DB-003 Receive nullable columns into pointers or sql.Null.

- Category: Programming Conventions/Database Access
- Since: Go 1.22
- Tags: data-structures

NULL and the Go zero value are different semantics.

- When a column is nullable, Scan into `*T` or `sql.Null[T]` (Go 1.22+ generics)
- Scanning NULL into a bare basic type fails outright
- When the schema changes (new nullable columns), re-check the Scan targets

**Why**

> NullString represents a string that may be null. NullString provides a convenient way of dealing with NULL string values.
>
> — https://pkg.go.dev/database/sql#NullString

Scanning NULL into a bare type fails with a conversion error, and the failure lands at the runtime data boundary, triggered by downstream DDL changes; pointers and sql.Null[T] encode nullable in the type, matching the database's semantics.

**Good**

```go
var nick *string
if err := rows.Scan(&id, &nick); err != nil {
  return err
}
```

**Bad**

```go
var nick string
if err := rows.Scan(&id, &nick); err != nil {
  return err
}
```

**References**

- https://pkg.go.dev/database/sql#NullString

**Detection**: manual review (Compare DDL nullable columns with the Scan target types)

## (8) Concurrency

### 【MUST】CONC-001 Before starting a goroutine, make its exit timing clear, and let the starter wait for it to finish.

- Category: Programming Conventions/Concurrency
- Since: Go 1.25
- Tags: concurrency

Before writing go, settle three things: when it ends, who waits for it to end, and who handles an error. The starter gathers it with `sync.WaitGroup` or `errgroup` in the same function, and uses `errgroup` when the subtask can return an error. Do not write a fire-and-forget with no waiter.

**Why**

> When you spawn goroutines, make it clear when - or whether - they exit. Goroutines can leak by blocking on channel sends or receives: the garbage collector will not terminate a goroutine even if the channels it is blocked on are unreachable.
>
> — https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes

A blocked goroutine is not reclaimed by garbage collection and keeps holding memory together with the references on its stack; with no waiter, the errors and panics of a background task have no owner and surface only in production.

**Good**

```go
var wg sync.WaitGroup
for _, item := range items {
	wg.Go(func() {
		process(item)
	})
}
wg.Wait()
```

**Bad**

```go
for _, item := range items {
	go process(item)
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes
- https://pkg.go.dev/golang.org/x/sync/errgroup

**Detection**: manual review (Check that each go statement has a matching waiter and exit path)

### 【MUST】CONC-002 A channel is closed only by the sender, never by the receiver.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

The close is bound to the producer's exit path, usually right after all sends complete. With multiple producers, do not close the data channel directly; use a separate stop signal together with context to coordinate. Cancellation always goes through context, and closing a channel is not used to signal cancellation.

**Why**

> Sends on a closed channel panic, so it's important to ensure all sends are done before calling close. ... stages close their outbound channels when all the send operations are done.
>
> — https://go.dev/blog/pipelines

Sending on a closed channel panics, and closing from the receiver side or closing twice blows up at run time; keeping the right to close with the sender, so whoever writes is responsible for wrapping up, makes the boundary unique.

**Good**

```go
ch := make(chan int)
go func() {
	defer close(ch)
	for _, v := range vals {
		ch <- v
	}
}()
```

**Bad**

```go
ch := make(chan int)
go func() {
	for _, v := range vals {
		ch <- v
	}
}()
<-ch
close(ch)
```

**References**

- https://go.dev/blog/pipelines

**Detection**: manual review (Check that every close is on the sender side)

### 【MUST】CONC-010 Wait with wg.Go, or place manual wg.Add before the go statement.

- Category: Programming Conventions/Concurrency
- Since: Go 1.25
- Tags: concurrency

On Go 1.25+, prefer wg.Go(f) over the manual Add/go/Done pairing.

- `wg.Go(f)` encapsulates Add, goroutine start, and Done in one step, eliminating pairing errors
- With manual use, `wg.Add(1)` must precede the go statement; inside the goroutine it races with Wait
- Use `errgroup` when collecting per-task errors; Go below 1.25 still pairs manually

**Why**

> Go calls the given function in a new goroutine. It blocks until all goroutines have exited.
>
> — https://pkg.go.dev/sync#WaitGroup.Go

When Add runs inside the goroutine, the main goroutine's Wait may pass the zero-count check before Add executes, and the wait becomes a no-op; wg.Go has the standard library guarantee the pairing order, removing this class of intermittent race at the root.

**Good**

```go
var wg sync.WaitGroup
for _, t := range tasks {
  wg.Go(func() { handle(t) })
}
wg.Wait()
```

**Bad**

```go
var wg sync.WaitGroup
for _, t := range tasks {
  go func() {
    wg.Add(1)
    defer wg.Done()
    handle(t)
  }()
}
wg.Wait()
```

**References**

- https://pkg.go.dev/sync#WaitGroup.Go
- https://go.dev/doc/go1.25

**Detection**: manual review (Check the relative position of wg.Add and the go statement)

### 【MUST】CONC-012 Do not request Lock while holding RLock on the same RWMutex.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

An RWMutex cannot upgrade; read and write lock requests must be one-directional.

- For writes while holding a read lock, `RUnlock` first, then `Lock`
- For read-modify-write atomicity, use `Lock` directly
- Recursive read locking is equally forbidden

**Why**

> If a goroutine holds a RWMutex for reading and another goroutine might call Lock, no goroutine should expect to be able to acquire a read lock until the initial read lock is released.
>
> — https://pkg.go.dev/sync#RWMutex

A pending write lock blocks subsequent read locks, so requesting the write lock while holding the read lock deadlocks on yourself: you are waiting for your own read lock to release; the documentation explicitly forbids recursive read locking too, so upgrade semantics must split into two steps or use a write lock throughout.

**Good**

```go
mu.RLock()
v := m[key]
mu.RUnlock()
if needUpdate(v) {
  mu.Lock()
  m[key] = compute()
  mu.Unlock()
}
```

**Bad**

```go
mu.RLock()
if needUpdate(m[key]) {
  mu.Lock()
  m[key] = compute()
  mu.Unlock()
}
mu.RUnlock()
```

**References**

- https://pkg.go.dev/sync#RWMutex

**Detection**: manual review (Check for nested RLock and Lock holds)

### 【SHOULD】CONC-003 A public API keeps synchronous semantics and leaves concurrency to the caller.

- Category: Programming Conventions/Concurrency
- Since: Go 1.7
- Tags: api-design, concurrency

Internally it may start several goroutines to fetch concurrently, but it gathers them before returning and hands back ordinary results and an error rather than a channel that leaves the caller to wait.

**Why**

> Prefer synchronous functions - functions which return their results directly or finish any callbacks or channel ops before returning - over asynchronous ones. ... If callers need more concurrency, they can add it easily by calling the function from a separate goroutine. But it is quite difficult - sometimes impossible - to remove unnecessary concurrency at the caller side.
>
> — https://go.dev/wiki/CodeReviewComments#synchronous-functions

An API returning a channel pushes scheduling and lifecycle management onto every caller, who each rewrite the waiting and cancellation and easily miss something; a synchronous interface keeps the concurrency inside the implementation so callers use it like an ordinary function.

**Good**

```go
func Fetch(ctx context.Context) (Data, error) {
	// concurrency inside, joined before returning
}
```

**Bad**

```go
func FetchAsync() <-chan Data {
	// the caller waits on its own
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#synchronous-functions

**Detection**: manual review (Check whether a public function returns a channel that leaves the caller waiting)

### 【SHOULD】CONC-004 Shared state and the lock protecting it live in the same struct and are read and written only through methods.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

The protected field is unexported, reads and writes go through locked methods, and the lock is not exposed outside the struct. Choose a lock or a channel by scenario: use `sync.Mutex` to protect shared data such as a cache or state, and use a channel to transfer data ownership, distribute tasks, or pass asynchronous results.

**Why**

> Use whichever is most expressive and/or most simple. A common Go newbie mistake is to over-use channels and goroutines just because it's possible, and/or because it's fun. Don't be afraid to use a sync.Mutex if that fits your problem best.
>
> — https://go.dev/wiki/MutexOrChannel

When the state and the lock are separate, a caller can bypass the method and change the field directly, leaving the lock useless; once they are together, the only entry for reads and writes is the method, and the protected range can be checked line by line.

**Good**

```go
type counter struct {
	mu sync.Mutex
	n  int
}

func (c *counter) inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
```

**Bad**

```go
type counter struct {
	mu sync.Mutex
	N  int // exported field, callers can bypass the lock and mutate it directly
}
```

**References**

- https://go.dev/wiki/MutexOrChannel
- https://pkg.go.dev/sync

**Detection**: manual review (Check whether the lock-protected field is exported and whether it is reachable only through methods)

### 【SHOULD】CONC-005 Channel parameters in function signatures are typed as directional.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

Pick chan<- or <-chan by the parameter's actual use.

- Send-only parameters use `chan<- T`, receive-only use `<-chan T`
- Directional types stop erroneous sends and closes at compile time
- Return values and struct fields follow the same rule

**Why**

> The optional <- operator specifies the channel direction, send or receive. If a direction is given, the channel is directional, otherwise it is bidirectional.
>
> — https://go.dev/ref/spec#Channel_types

A bidirectional chan T places no compile-time constraint on erroneous sends or closes, and closing a producer's channel makes consumers panic; directional types put the usage boundary in the signature, so misuse fails compilation instead.

**Good**

```go
func consume(in <-chan int) { for v := range in { process(v) } }
```

**Bad**

```go
func consume(in chan int) { for v := range in { process(v) } }
```

**References**

- https://go.dev/ref/spec#Channel_types

**Detection**: manual review (Check the direction of channel parameters in exported functions)

### 【SHOULD】CONC-006 Buffered channels with capacity above 1 carry a comment on sizing and overflow behavior.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

Buffer capacity is a concurrency design decision, not an arbitrary number.

- Unbuffered by default; buffers above 1 must note the sizing rationale and what happens when full (block, drop, or degrade)
- Large buffers mask backpressure, and queued work grows unbounded
- Producer/consumer rate alignment notes go next to the make statement

**Why**

> Channels should usually be size one or give an unbuffered implementation. A channel size of zero or one is sufficient in most scenarios.
>
> — https://github.com/uber-go/guide/blob/master/style.md#channel-size

An unbuffered channel provides natural backpressure, while a buffered one delays the signal of overproduction until after the backlog builds; an unannotated capacity cannot be told apart from a casual number, so review cannot judge whether the full-buffer behavior is safe.

**Good**

```go
// 容量 64 对齐下游批处理批次；写满时生产端阻塞，保证背压
jobs := make(chan Job, 64)
```

**Bad**

```go
jobs := make(chan Job, 1000)
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#channel-size

**Detection**: grep regex make\(chan[^)]*,\s*([^)1][^)]*)\) (Grep for make(chan, n) with n greater than 1 and no adjacent comment)

### 【SHOULD】CONC-007 Selectable cancelable paths include a ctx.Done() branch.

- Category: Programming Conventions/Concurrency
- Since: Go 1.7
- Tags: concurrency, context

Long-lived goroutines need a cancellation exit in their select loop.

- Put `ctx.Done()` in the select case set and exit on cancellation
- Only one-shot selects with certain termination may omit it
- Goroutines without a cancellation branch cannot exit when the service shuts down

**Why**

> Done returns a channel that's closed when work done on behalf of this context should be canceled. Done is provided for use in select statements
>
> — https://pkg.go.dev/context#Context.Done

A select blocks on multiple channels; if the cancellation signal is not in the case set, the goroutine keeps blocking after the caller cancels, and every such goroutine leaks during graceful shutdown.

**Good**

```go
for {
  select {
  case <-ctx.Done():
    return
  case v := <-in:
    process(v)
  }
}
```

**Bad**

```go
for {
  select {
  case v := <-in:
    process(v)
  }
}
```

**References**

- https://pkg.go.dev/context#Context.Done

**Detection**: manual review (Check the case set of long-lived goroutine selects)

### 【SHOULD】CONC-008 A select inside a for loop does not busy-wait via a default branch.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

The default branch of an in-loop select spins the loop.

- Drop `default` when waiting for data, letting select block
- For non-blocking polling, sleep or back off inside the `default` branch
- A receiver sets a closed channel to nil to exit the multiplexing

**Why**

> Execution of a "select" statement proceeds in several steps: ... If multiple cases can proceed, a uniform pseudo-random choice is made
>
> — https://go.dev/ref/spec#Select_statements

`for { select { default: } }` retries immediately whenever no case is ready, so the spinning loop keeps one CPU core busy doing nothing; a closed channel is always ready and consumes a core the same way once mixed into a select, and setting it to nil is the standard way to exit.

**Good**

```go
for {
  select {
  case <-ctx.Done():
    return
  case v, ok := <-in:
    if !ok {
      in = nil
      continue
    }
    process(v)
  }
}
```

**Bad**

```go
for {
  select {
  case v := <-in:
    process(v)
  default:
  }
}
```

**References**

- https://go.dev/ref/spec#Select_statements

**Detection**: manual review (Check the default branch of selects inside loops)

### 【SHOULD】CONC-009 Periodic ticks in loops use time.NewTicker with Stop.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

Periodic tasks that can stop avoid time.After and time.Tick.

- Each `time.After` in a loop allocates a fresh timer; old ones are unreclaimable until fired
- The ticker behind `time.Tick` cannot be stopped and leaks forever
- Use `ticker := time.NewTicker(d)` with `defer ticker.Stop()`

**Why**

> While Tick is useful for clients that have no need to shut down the Ticker, be aware that without a way to shut it down the underlying Ticker can never be recovered by the garbage collector; it "leaks".
>
> — https://pkg.go.dev/time#Tick

time.After allocates a new timer per loop iteration that the GC cannot reclaim until it fires, showing up as memory and timer buildup in long loops; the official documentation marks time.Tick as leaking outright, and NewTicker plus Stop is the only path with stop semantics.

**Good**

```go
ticker := time.NewTicker(interval)
defer ticker.Stop()
for {
  select {
  case <-ticker.C:
    tick()
  case <-ctx.Done():
    return
  }
}
```

**Bad**

```go
for {
  select {
  case <-time.After(interval):
    tick()
  case <-ctx.Done():
    return
  }
}
```

**References**

- https://pkg.go.dev/time#Tick

**Detection**: manual review (Check for time.After and time.Tick inside loop bodies)

### 【SHOULD】CONC-011 Critical sections exclude I/O, network calls, and channel operations.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

Lock hold time covers memory operations only.

- Blocking calls over network, disk, and databases stay outside the lock
- When lock-held data is needed for I/O, copy the result out, then unlock
- Channel sends and receives can block indefinitely and move out as well

**Why**

> Programs that modify data being simultaneously accessed by multiple goroutines must serialize such access.
>
> — https://go.dev/ref/mem

Making a network call while holding a lock stretches the hold time to the network round trip, and milliseconds of I/O serialize every goroutine waiting on that lock; moving I/O outside brings the hold time back to nanoseconds, so concurrency does not decay with external latency.

**Good**

```go
resp, err := http.Get(url)
if err != nil {
  return err
}
mu.Lock()
cache[key] = resp
mu.Unlock()
```

**Bad**

```go
mu.Lock()
resp, err := http.Get(url)
if err != nil {
  mu.Unlock()
  return err
}
cache[key] = resp
mu.Unlock()
```

**References**

- https://go.dev/ref/mem

**Detection**: manual review (Check for blocking calls between Lock and Unlock)

### 【SHOULD】CONC-013 Simple counters and flags use typed atomic types.

- Category: Programming Conventions/Concurrency
- Since: Go 1.19
- Tags: concurrency

atomic.Int64, atomic.Bool, and siblings replace the function-style atomic calls.

- Single-value state such as counters, switches, and version numbers uses `atomic.Int64`/`atomic.Int32`/`atomic.Bool`/`atomic.Pointer`
- Do not use the function form `atomic.AddInt64(&counter, 1)`
- Compound state still uses a mutex

**Why**

> Int64 is an atomically operated int64 value.
>
> — https://pkg.go.dev/sync/atomic#Int64

Typed atomics put this field is atomically accessed into the type itself, so passing a plain int64 fails at compile time; with function-style calls the target field is type-identical to any other field, and new access points can easily miss the atomic operation.

**Good**

```go
var counter atomic.Int64
counter.Add(1)
```

**Bad**

```go
var counter int64
atomic.AddInt64(&counter, 1)
```

**References**

- https://pkg.go.dev/sync/atomic#Int64

**Detection**: manual review (Check the form of atomic operations)

### 【SHOULD】CONC-014 Do not use sync.Map in write-heavy scenarios.

- Category: Programming Conventions/Concurrency
- Since: Go 1.9
- Tags: concurrency, data-structures

sync.Map serves two access patterns only.

- Keys written once and read many times (grow-only caches)
- Multiple goroutines operating on disjoint key sets
- Everything else uses `map` with `RWMutex`; `sync.Map` is slower on write-heavy paths and type-unsafe

**Why**

> The Map type is optimized for two common use cases: (1) when the entry for a given key is only ever written once but read many times, as in caches that only grow, or (2) when multiple goroutines read, write, and overwrite entries for disjoint keys.
>
> — https://pkg.go.dev/sync#Map

sync.Map trades a read-optimized internal structure for lock-free reads, and its maintenance cost exceeds map+RWMutex under heavy writes; the official documentation names exactly two optimized modes, so using it beyond them has no performance basis and loses compile-time type checking.

**Good**

```go
mu.Lock()
m[key] = value
mu.Unlock()
```

**Bad**

```go
sm.Store(key, value)
```

**References**

- https://pkg.go.dev/sync#Map

**Detection**: manual review (Check the read/write ratio and key-set pattern of sync.Map usage)

### 【SHOULD】CONC-015 Batch concurrent tasks set a concurrency limit.

- Category: Programming Conventions/Concurrency
- Since: Go 1.0
- Tags: concurrency

Unbounded goroutine startup at input scale exhausts resources.

- Limit with `errgroup.Group`'s `SetLimit` or a buffered-channel semaphore
- The limit follows downstream capacity (connection pool, downstream QPS) and is annotated
- One-off small tasks (bounded fan-out input) may be waived with a note

**Why**

> SetLimit limits the number of active goroutines in this group to at most n.
>
> — https://pkg.go.dev/golang.org/x/sync/errgroup#Group.SetLimit

Every goroutine occupies stack memory and scheduling cost, so starting one goroutine per item at a scale of ten thousand exhausts memory and downstream connections quickly; a fixed limit keeps concurrency bounded, and errgroup also collects errors in one place.

**Good**

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(8)
for _, t := range tasks {
  g.Go(func() error { return handle(ctx, t) })
}
err := g.Wait()
```

**Bad**

```go
for _, t := range tasks {
  go handle(ctx, t)
}
```

**References**

- https://pkg.go.dev/golang.org/x/sync/errgroup#Group.SetLimit

**Detection**: manual review (Check that batch task fan-out sets a limit)

### 【SHOULD】CONC-016 sync.Pool objects are Reset before being put back.

- Category: Programming Conventions/Concurrency
- Since: Go 1.3
- Tags: concurrency, performance

Objects returned to the pool must clear their internal state.

- `Reset` slices, maps, and pointer fields so stale data does not leak into the next `Get`
- Do not assume `Get` returns a zero-valued object
- Document the `Reset` obligation on pooled types

**Why**

> Pool's purpose is to cache allocated but unused items for later reuse, relieving pressure on the garbage collector.
>
> — https://pkg.go.dev/sync#Pool

The pool caches memory without resetting state, and a dirty object carries the previous user's data into the next Get, leaking user data across requests or corrupting logic; Get means usable memory, not a zero value, so the reset obligation sits with the caller.

**Good**

```go
buf := pool.Get().(*bytes.Buffer)
defer func() {
  buf.Reset()
  pool.Put(buf)
}()
```

**Bad**

```go
buf := pool.Get().(*bytes.Buffer)
defer pool.Put(buf)
```

**References**

- https://pkg.go.dev/sync#Pool

**Detection**: manual review (Check that Pool.Put sites Reset first)

## (9) Control Flow

### 【MUST】CTRL-003 Do not defer resource release inside a loop body.

- Category: Programming Conventions/Control Flow
- Since: Go 1.0
- Tags: control-flow, data-structures

defer binds to the function, not the loop iteration.

- Close resources opened in a loop within the same iteration, or extract the loop body into a function
- deferring in a long loop releases everything only when the function returns, exhausting file descriptors and connections
- For cleanup-at-the-end semantics, collect explicitly and close in one place

**Why**

> The deferred call's arguments are evaluated immediately, but the function call is not executed until the surrounding function returns.
>
> — https://go.dev/doc/effective_go#defer

deferred calls run when the surrounding function returns, so a thousand iterations defer a thousand releases, and handles only accumulate for the lifetime of the function; on servers this shows up as fd exhaustion and dead pool connections in long loops.

**Good**

```go
for _, name := range names {
  func() {
    f, err := os.Open(name)
    if err != nil {
      return
    }
    defer f.Close()
    process(f)
  }()
}
```

**Bad**

```go
for _, name := range names {
  f, err := os.Open(name)
  if err != nil {
    return
  }
  defer f.Close()
  process(f)
}
```

**References**

- https://go.dev/doc/effective_go#defer

**Detection**: manual review (Grep for defer statements inside loop bodies)

### 【MUST】CTRL-004 HTTP handlers return immediately after writing an error response.

- Category: Programming Conventions/Control Flow
- Since: Go 1.0
- Tags: control-flow, error-handling

Business logic must not continue after `http.Error` or `w.Write`.

- The error branch returns right after writing the response
- Continuing may call `WriteHeader` again, triggering the superfluous WriteHeader warning
- The same applies to middleware that rewrites responses

**Why**

> Changing the header map after a call to WriteHeader (or Write) has no effect.
>
> — https://pkg.go.dev/net/http#ResponseWriter

Code that runs after the response was written either calls WriteHeader a second time (the header is already sent, changes are ignored, and a warning is logged) or writes sensitive error details into the body; an immediate return closes the error branch.

**Good**

```go
if err := parse(r); err != nil {
  http.Error(w, "bad request", http.StatusBadRequest)
  return
}
```

**Bad**

```go
if err := parse(r); err != nil {
  http.Error(w, "bad request", http.StatusBadRequest)
}
process(r)
```

**References**

- https://pkg.go.dev/net/http#ResponseWriter

**Detection**: manual review (Check that handlers return after http.Error/w.Write)

### 【SHOULD】CTRL-001 Handle errors and edge cases first and return early, keeping the normal path unindented.

- Category: Programming Conventions/Control Flow
- Since: Go 1.0
- Tags: control-flow

Rewrite the form that wraps the success branch of an if in an else into one that checks for the error and returns, with the normal logic following; when an if with a short variable declaration defines a variable still used later, move the declaration before the if.

**Why**

> Try to keep the normal code path at a minimal indentation, and indent the error handling, dealing with it first. This improves the readability of the code by permitting visually scanning the normal path quickly.
>
> — https://go.dev/wiki/CodeReviewComments#indent-error-flow

The normal path is the part readers look at most, and wrapping it in an else adds one indent level and one negation of the condition; with an early return, the normal path stays at the far left and reads straight through from top to bottom.

**Good**

```go
if err := save(); err != nil {
	return err
}
process()
```

**Bad**

```go
if err := save(); err == nil {
	process()
} else {
	return err
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#indent-error-flow

**Detection**: manual review (Check whether the normal path is indented by an else)

### 【SHOULD】CTRL-002 Do not explicitly copy the loop variable any longer.

- Category: Programming Conventions/Control Flow
- Since: Go 1.22
- Tags: modernization

Since Go 1.22 the iteration variable of a for loop is independent on each iteration, and a closure or goroutine captures the value of that iteration, so writing `v := v` or passing the variable as an argument into the goroutine is no longer needed. This clause applies to modules whose go.mod declares Go 1.22 or above.

**Why**

> Previously, the variables declared by a "for" loop were created once and updated by each iteration. In Go 1.22, each iteration of the loop creates new variables, to avoid accidental sharing bugs.
>
> — https://go.dev/doc/go1.22

The explicit copy was the way to avoid a shared variable before Go 1.22, and keeping it after the language was fixed only makes readers think protection is still needed here; letting copyloopvar detect it cleans it up uniformly.

**Good**

```go
fns := make([]func(), 0, len(items))
for _, item := range items {
	fns = append(fns, func() {
		process(item)
	})
}
```

**Bad**

```go
fns := make([]func(), 0, len(items))
for _, item := range items {
	item := item
	fns = append(fns, func() {
		process(item)
	})
}
```

**References**

- https://go.dev/doc/go1.22
- https://go.dev/blog/loopvar-preview

**Detection**: golangci-lint copyloopvar

### 【SHOULD】CTRL-005 Inner scopes do not shadow the outer err with :=.

- Category: Programming Conventions/Control Flow
- Since: Go 1.0
- Tags: control-flow, error-handling

Keep one identifier for an error variable across scopes.

- The err in `if v, err := f(); err != nil` is unrelated to the outer err, and assignments do not escape
- To carry an inner result out, declare first, then assign
- The shadow checker must be enabled explicitly; waive noisy hits with nolint

**Why**

> check for shadowed variables
>
> — https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow

:= declares a new variable in the new scope, so the err inside an if statement is a fresh one — checking it leaves the outer err still nil, silently bypassing error handling; such defects surface neither at compile time nor at runtime, only in review or via the shadow checker.

**Good**

```go
var err error
for _, p := range paths {
  data, err := os.ReadFile(p)
  if err != nil {
    return err
  }
  _ = data
}
return err
```

**Bad**

```go
var err error
for _, p := range paths {
  data, err := os.ReadFile(p)
  if err != nil {
    return err
  }
  _ = data
}
return nil
```

**References**

- https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow

**Detection**: golangci-lint govet(shadow) (Check for := declarations that shadow an outer err)

## (10) Comments

### 【SHOULD】COMMENT-001 Comments record only what the code cannot express and never restate what the code does.

- Category: Programming Conventions/Comments
- Since: Go 1.21
- Tags: comments

What naming and structure already make clear needs no comment. Only these cases are worth a comment:

- The reason the code is written this way cannot be seen from the code itself
- The code implements a well-known algorithm
- The code applies a mathematical formula

A comment explains why it is done this way and does not restate what the code does.

**Why**

> It is often better for comments to explain why something is done, not what the code is doing.
>
> — https://google.github.io/styleguide/go/guide#clarity-rationale

The code itself already reads as what it does, so writing it again means maintaining two explanations that can diverge at any time, and readers must decide which one is correct. Reserving comments for what the code cannot say is what makes comments worth trusting.

**Good**

```go
// Make a copy: the caller may keep reusing the incoming slice.
entries := slices.Clone(in)
```

**Bad**

```go
// copy the slice
entries := slices.Clone(in)
```

**References**

- https://google.github.io/styleguide/go/guide#clarity-rationale

**Detection**: manual review (Check whether a comment restates what the code itself already shows)

### 【SHOULD】COMMENT-002 A doc comment begins with the name being described and is written as a full sentence.

- Category: Programming Conventions/Comments
- Since: Go 1.0
- Tags: comments

When a doc comment is needed (package comments, and exported identifiers that genuinely warrant explanation), it begins with the name followed by one sentence and ends with a period. This clause governs only how a doc comment is written; it does not require a comment on every exported identifier, and whether to write one is decided by the stance of COMMENT-001.

**Why**

> Comments documenting declarations should be full sentences, even if that seems a little redundant. This approach makes them format well when extracted into godoc documentation. Comments should begin with the name of the thing being described and end in a period.
>
> — https://go.dev/wiki/CodeReviewComments#comment-sentences

When a comment begins with the name, the summary is a full sentence on the godoc index page and the name link is recognized correctly; without the name the summary reads like half a sentence and the link is easily missed.

**Good**

```go
// Load reads and parses the clause data file.
func Load(path string) (*Manual, error) {
```

**Bad**

```go
// read the clause data
func Load(path string) (*Manual, error) {
```

**References**

- https://go.dev/wiki/CodeReviewComments#comment-sentences
- https://go.dev/doc/effective_go#commentary

**Detection**: manual review (Check whether a doc comment begins with the name and forms a full sentence)

