# 1. Programming Conventions

## (1) Naming

### 【MUST】NAMING-001 Package names use lowercase words written together, without underscores, camel case, or plurals.

- Category: Programming Conventions/Naming
- Since: Go 1.0

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

## (2) Constants and Enums

### 【MUST】CONST-001 Constant names use MixedCaps, not all uppercase or a k prefix.

- Category: Programming Conventions/Constants and Enums
- Since: Go 1.0

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

### 【MAY】CONST-003 Enumeration values that increase from one source are generated with iota.

- Category: Programming Conventions/Constants and Enums
- Since: Go 1.0

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

### 【SHOULD】FUNC-002 Use named results when the meaning of the return values is not self-evident.

- Category: Programming Conventions/Functions and Methods
- Since: Go 1.0

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

## (5) Data Structures

### 【SHOULD】DATA-001 Declare empty slices with var to get a nil slice.

- Category: Programming Conventions/Data Structures
- Since: Go 1.0

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

A literal writes field name: value for each field. Small private structs with few fields and a stable order, and test table data, may be exceptions; a struct from another package always uses field names, because adding, removing, or reordering fields silently misplaces positional entries.

**Why**

> You should almost always specify field names when initializing structs. This is now enforced by go vet.
>
> — https://github.com/uber-go/guide/blob/master/style.md#use-field-names-to-initialize-structs

The positional form turns field order into an implicit contract, so after the other side adjusts the struct, positions of the same type are silently mismatched with no compile error; the field name form binds each value to its field so the change surfaces at once. The composites check in go vet covers positional literals from another package.

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

### 【MAY】MODERN-001 Common operations on slices and maps prefer the slices and maps standard library packages over hand-written loops.

- Category: Programming Conventions/Data Structures
- Since: Go 1.21

slices and maps provide high-frequency operations such as Contains, Sort, Clone, and Keys. When the standard library can express the operation, stop writing loops by hand, which reduces boundary mistakes and unifies the style.

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

## (6) Concurrency

### 【MUST】CONC-001 Before starting a goroutine, make its exit timing clear, and let the starter wait for it to finish.

- Category: Programming Conventions/Concurrency
- Since: Go 1.25

Before writing go, settle three things: when it ends, who waits for it to end, and who handles an error. The starter gathers it with sync.WaitGroup or errgroup in the same function, and uses errgroup when the subtask can return an error. Do not write a fire-and-forget with no waiter.

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

### 【SHOULD】CONC-003 A public API keeps synchronous semantics and leaves concurrency to the caller.

- Category: Programming Conventions/Concurrency
- Since: Go 1.7

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

The protected field is unexported, reads and writes go through locked methods, and the lock is not exposed outside the struct. Choose a lock or a channel by scenario: use sync.Mutex to protect shared data such as a cache or state, and use a channel to transfer data ownership, distribute tasks, or pass asynchronous results.

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

## (7) Control Flow

### 【SHOULD】CTRL-001 Handle errors and edge cases first and return early, keeping the normal path unindented.

- Category: Programming Conventions/Control Flow
- Since: Go 1.0

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

Since Go 1.22 the iteration variable of a for loop is independent on each iteration, and a closure or goroutine captures the value of that iteration, so writing v := v or passing the variable as an argument into the goroutine is no longer needed. This clause applies to modules whose go.mod declares Go 1.22 or above.

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

## (8) Comments

### 【SHOULD】COMMENT-001 Comments record only what the code cannot express and never restate what the code does.

- Category: Programming Conventions/Comments
- Since: Go 1.21

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

