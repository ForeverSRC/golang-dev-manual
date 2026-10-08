# 5. Interfaces and Design

## (1) Interface Definition and Placement

### 【SHOULD】ABSTRACT-001 Introduce abstraction according to actual need, not spread out in advance.

- Category: Interfaces and Design/Interface Definition and Placement
- Since: Go 1.0
- Tags: api-design

Introduce interfaces, factories, and registries only when a second implementation appears, a test double is needed, or there is a genuine variation point to isolate; when there is only one implementation and no replacement or testing need, use a concrete type first. Layering and ports count as legitimate abstraction when there is a genuine isolation need, and this clause does not restrict them.

**Why**

> Do not define interfaces before they are used: without a realistic example of usage, it is too difficult to see whether an interface is even necessary, let alone what methods it ought to contain.
>
> — https://go.dev/wiki/CodeReviewComments#interfaces

Abstraction laid out in advance adds a layer of indirection and maintenance cost, and when the real requirement arrives its direction often differs from the preset, forcing the interface to be changed again; abstracting once real use cases appear lets the interface be defined by the consumer, so its shape fits.

**Good**

```go
type Loader struct {
	path string
}

func (l *Loader) Load() (*domain.Manual, error) {
	return readJSON(l.path)
}
```

**Bad**

```go
// only JSON exists today, yet an interface, a factory, and a format registry are defined up front for a future YAML source
type Source interface {
	Load() (*domain.Manual, error)
}

type SourceFactory interface {
	New(format string) (Source, error)
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#interfaces

**Detection**: manual review (Verify whether the interface or factory already has a second implementation, a replacement, or a testing need.)

### 【SHOULD】IFACE-001 Place interfaces according to the number of implementations: a single implementation shares the package with the implementation, while multiple implementations are defined by the consumer.

- Category: Interfaces and Design/Interface Definition and Placement
- Since: Go 1.0
- Tags: api-design, layering

Placement depends on how many implementations are expected:

- Only ever a single implementation: put the interface and the implementation in the same package, with the interface adjacent to the implementation, so that reading the code does not require cross-package comparison; export the interface and the constructor, and keep the implementation struct unexported.
- Multiple implementations expected, or replacement needed (storage, external services, notifications, and the like): put the interface in the package of the consumer, selecting only the methods the consumer actually needs; the implementation package exports the concrete type and constructor, and writes a compile-time assertion per IFACE-002 to put the promise into code.
- The interface signature is already fixed by the standard library or a third party (`io.Reader` and the like): this clause does not apply, and the implementer simply satisfies it as is.

**Why**

> Go interfaces generally belong in the package that uses values of the interface type, not the package that implements those values. The implementing package should return concrete (usually pointer or struct) types: that way, new methods can be added to implementations without requiring extensive refactoring.
>
> — https://go.dev/wiki/CodeReviewComments#interfaces

When there is a single implementation, splitting the interface into another package forces the reader to compare the implementation's context back and forth between two packages, and a change to one method signature must be aligned in both places; keeping the interface and implementation in the same package lets a single file show the whole contract and behavior. The common upstream Go wiki statement is that "interfaces belong to the package that uses values of the interface type" (quoted below), which targets dependencies that will be replaced or used by multiple parties, and its orientation differs from the single-implementation branch of this clause—the single-implementation branch follows local reading habits and is this manual's own convention; the multiple-implementation branch follows the upstream statement, in which case the interface contains only the methods the consumer needs, and the implementer can add methods without going back to change the interface.

**Good**

```go
// single implementation: interface and implementation live in the same package, contract and behavior are readable in one place
package service

type ManualService interface {
	List(ctx context.Context) ([]Clause, error)
}

type manualService struct{ repo ClauseRepository }
```

**Bad**

```go
// single implementation, yet the interface lives in the consumer package, so reading the implementation means cross-package comparison
package clihandler

type ManualService interface {
	List(ctx context.Context) ([]Clause, error)
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#interfaces

**Detection**: manual review (Verify whether the package holding the interface matches the number of implementations: for a single implementation, whether it shares the package with the implementation; for multiple implementations, whether it is in the consumer's package.)

### 【SHOULD】IFACE-002 Declare a type's implementation of an interface with a compile-time assertion.

- Category: Interfaces and Design/Interface Definition and Placement
- Since: Go 1.0
- Tags: api-design

Write `var _ Interface = zero value of the type` near the type definition, using nil for pointers, slices, and maps, and an empty struct literal for structs. When the interface is defined by the consumer, this clause is used as a pair with IFACE-001: the consumer declares the interface and the implementation package writes the assertion, and if either is missing the practice does not hold. Applicable scenarios:

- An interface defined by the consumer per IFACE-001 in a multiple-implementation scenario
- An exported type that must implement a standard library or third-party interface per an API contract
- A group of types that together implement the same interface
- Failing to satisfy the interface would make callers fail to compile outright

**Why**

> Verify interface compliance at compile time where appropriate. ... The statement `var _ http.Handler = (*Handler)(nil)` will fail to compile if `*Handler` ever stops matching the `http.Handler` interface.
>
> — https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance

The assertion hands the judgment of "whether the type satisfies the interface" to the compiler: if the interface changes or a method signature no longer matches, the build fails on the spot and the person making the change knows immediately; without an assertion, a mismatch surfaces only when the value is assigned to the interface type somewhere, often far from the change.

**Good**

```go
type Handler struct{}

var _ http.Handler = (*Handler)(nil)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
```

**Bad**

```go
type Handler struct{}

// without the assertion, the moment *Handler stops satisfying http.Handler only surfaces at the call site
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance

**Detection**: manual review (Verify whether exported types that implement an interface per an API contract have a compile-time assertion.)

### 【SHOULD】IFACE-003 Do not define pointers to interfaces.

- Category: Interfaces and Design/Interface Definition and Placement
- Since: Go 1.0
- Tags: api-design

Use interface values themselves for parameters, return values, and fields, not `*Interface`. An interface value already contains the type information and a data pointer internally, and whether the underlying data is a value or a pointer is decided by the concrete type; when methods need to modify the underlying data, the implementer expresses that with a pointer receiver.

**Why**

> You almost never need a pointer to an interface. You should be passing interfaces as values—the underlying data can still be a pointer.
>
> — https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces

`*Interface` adds a layer of dereferencing at call sites and makes nil checks more complex: an interface value that is nil and an interface whose internal data pointer is nil are two states, and once mixed together the emptiness condition is either written wrong or guessed. Passing the interface value itself expresses the same intent with a single meaning for emptiness.

**Good**

```go
func Save(w io.Writer) error
```

**Bad**

```go
func Save(w *io.Writer) error
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces

**Detection**: manual review (Verify whether function signatures and struct fields contain *interface types.)

## (2) Composition and Reuse

### 【SHOULD】COMPOSE-001 Express reusable capability through composition, not inheritance-style hierarchies.

- Category: Interfaces and Design/Composition and Reuse
- Since: Go 1.0
- Tags: api-design

To reuse a type's capability, hold it explicitly as a field; use anonymous embedding only when its methods also need to be exposed. Naming a type Base, Abstract, or the like and then embedding it simulates an inheritance hierarchy and should be avoided. Interfaces likewise compose small interfaces by embedding, rather than building one all-encompassing interface.

**Why**

> Go does not provide the typical, type-driven notion of subclassing, but it does have the ability to "borrow" pieces of an implementation by embedding types within a struct or interface. ... There's an important way in which embedding differs from subclassing. When we embed a type, the methods of that type become methods of the outer type, but when they are invoked the receiver of the method is the inner type, not the outer one.
>
> — https://go.dev/doc/effective_go#embedding

A composition relationship written on a field can be located directly for replacement, testing, and reading; Go's embedding has no override rules like Java's, and which layer a method lands on is inferred by name resolution, so simulating inheritance leaves readers unable to tell which implementation is actually called.

**Good**

```go
type Server struct {
	logger *slog.Logger
	store  Store
}
```

**Bad**

```go
type Server struct {
	BaseServer
}
```

**References**

- https://go.dev/doc/faq#inheritance
- https://go.dev/doc/effective_go#embedding

**Detection**: manual review (Verify whether embedded types play the role of a base class (Base / Abstract naming, multi-level embedding).)

### 【SHOULD】COMPOSE-002 Embedding must not change the outer type's zero value, copy semantics, or public interface.

- Category: Interfaces and Design/Composition and Reuse
- Since: Go 1.0
- Tags: api-design

Use anonymous embedding only when the inner type's methods should become the outer type's methods; otherwise use a named field. After embedding, the outer zero value must still be usable, the copy semantics unchanged, and no unrelated methods additionally exposed. Synchronization primitives such as `sync.Mutex` always use a named field, since anonymous embedding hangs `Lock` and `Unlock` onto the outer API.

**Why**

> Embedding should provide tangible benefit, like adding or augmenting functionality in a semantically-appropriate way. It should do this with zero adverse user-facing effects.
>
> — https://github.com/uber-go/guide/blob/master/style.md#embedding-in-structs

Once embedding picks an exported type, any method added to the inner type later becomes a public method of the outer type, and replacing the outer type or swapping the inner one is a breaking change; embedding a pointer type also makes the outer zero value panic as soon as it is called. A named field exposes only the set of methods you want to expose, and the outer shape is yours to decide.

**Good**

```go
type Server struct {
	mu     sync.Mutex
	logger *slog.Logger
}
```

**Bad**

```go
type Server struct {
	sync.Mutex
	*bytes.Buffer
}
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#embedding-in-structs
- https://github.com/uber-go/guide/blob/master/style.md#avoid-embedding-types-in-public-structs

**Detection**: manual review (Verify whether embedding changes the outer zero value, copy semantics, or exposed unrelated methods (especially for synchronization primitives).)

## (3) Dependency Injection and Wiring

### 【SHOULD】DI-001 Use wire for dependency wiring, generated at compile time.

- Category: Interfaces and Design/Dependency Injection and Wiring
- Since: Go 1.0
- Tags: dependency-injection

Concentrate each layer's constructors and interface bindings in a single ProviderSet, and write interface-to-implementation bindings with `wire.Bind`; the entry point only calls the generated `Initialize` function. The `wire` version is pinned via the tool directive in go.mod; after changing the wiring, regenerate wire_gen.go with `go generate`.

**Why**

> Wire is a code generation tool that automates connecting components using dependency injection. ... Because Wire operates without runtime state or reflection, code written to be used with Wire is useful even for hand-written initialization.
>
> — https://github.com/google/wire

Runtime reflective injection defers a mistyped field name or a missing binding until startup or runtime to surface; compile-time generation reports errors directly when a binding is missing or duplicated, and the wiring relationships can also be checked line by line in the generated code.

**Good**

```go
var ProviderSet = wire.NewSet(
	jsonfile.New,
	wire.Bind(new(service.ClauseRepository), new(*jsonfile.Repository)),
)
```

**Bad**

```go
container := dig.New()
if err := container.Provide(jsonfile.New); err != nil {
	panic(err)
}
```

**References**

- https://github.com/google/wire

**Detection**: manual review (Verify that the repository has wire_gen.go and that port bindings are concentrated in the ProviderSet.)

### 【SHOULD】DI-002 Pass dependencies in explicitly through constructors, and do not read or write package-level mutable variables.

- Category: Interfaces and Design/Dependency Injection and Wiring
- Since: Go 1.0
- Tags: dependency-injection

Dependencies such as time sources, random sources, clients, and configuration are passed in as constructor parameters or fields, and package-level variables hold only read-only constants. When a test needs to replace behavior, pass a double rather than changing global state. The same applies to function pointers and values of other types.

**Why**

> Avoid mutating global variables, instead opting for dependency injection. This applies to function pointers as well as other kinds of values.
>
> — https://github.com/uber-go/guide/blob/master/style.md#avoid-mutable-globals

Mutable globals hide dependencies outside the call chain: reading the code does not reveal which time source or client this logic uses, a test must save the old value before replacing it and restore it afterward, and parallel test cases interfere with each other. Passed in as a parameter, the dependency is written on the constructor signature, making replacement and reading direct.

**Good**

```go
type Signer struct {
	now func() time.Time
}

func NewSigner(now func() time.Time) *Signer {
	return &Signer{now: now}
}
```

**Bad**

```go
var now = time.Now

func Sign(msg string) string {
	return signWithTime(msg, now())
}
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#avoid-mutable-globals

**Detection**: manual review (Verify whether package-level variables are ever assigned and whether dependencies are passed in through constructors.)

## (4) Passing context

### 【MUST】CTX-002 Do not pass a nil context to a function; use context.TODO() as a placeholder.

- Category: Interfaces and Design/Passing context
- Since: Go 1.7
- Tags: context

When the upstream context is not available yet, pass `context.TODO()`; use `context.Background()` at program entry, in initialization, and in tests. Do not treat nil as the expression of "no context".

**Why**

> Do not pass a nil Context, even if a function permits it. Pass context.TODO if you are unsure about which Context to use.
>
> — https://pkg.go.dev/context#pkg-overview

nil is a legal interface value and passes at compile time; the recipient panics outright as soon as it calls `Done` or `Value`, and whether it panics depends on its internal implementation, which is invisible from the function signature. Passing `context.TODO()` has the same semantics and is always safe, and static analysis can also use it to find such omissions.

**Good**

```go
ctx := context.TODO()
return svc.Load(ctx, path)
```

**Bad**

```go
return svc.Load(nil, path)
```

**References**

- https://pkg.go.dev/context#pkg-overview
- https://staticcheck.dev/docs/checks/#SA1012

**Detection**: golangci-lint staticcheck(SA1012) (SA1012 detects calls that pass a nil context to a function.)

### 【SHOULD】CTX-001 Do not store context.Context in a struct; pass it as the first parameter of a function.

- Category: Interfaces and Design/Passing context
- Since: Go 1.7
- Tags: context

A context carries the cancellation, timeout, and values of a single request, with a lifetime matching the request. When passing it across layers, make the context the first parameter and name it ctx consistently.

**Why**

> Most functions that use a Context should accept it as their first parameter. ... Don't add a Context member to a struct type; instead add a ctx parameter to each method on that type that needs to pass it along.
>
> — https://go.dev/wiki/CodeReviewComments#contexts

Once stored in a struct, the context lifetime is detached from the request, the framework's cancellation propagation can no longer cover it, leaks easily form, and the cancellation boundary cannot be expressed on the interface.

**Good**

```go
func (s *Service) Get(ctx context.Context, id string) (*Order, error) {
	return s.repo.Find(ctx, id)
}
```

**Bad**

```go
type Service struct {
	ctx  context.Context
	repo Repository
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#contexts
- https://pkg.go.dev/context

**Detection**: golangci-lint containedctx (The first-parameter convention requires manual verification.)

### 【SHOULD】CTX-003 Use an unexported custom type for context.Value keys.

- Category: Interfaces and Design/Passing context
- Since: Go 1.7
- Tags: context

Each package defines an unexported type such as `type key int` and unexported constants for its own value retrieval, then provides type-safe accessor functions so that callers do not touch the key directly. Do not use built-in types such as string or int as keys.

**Why**

> A key can be any type that supports equality; packages should define keys as an unexported type to avoid collisions.
>
> — https://pkg.go.dev/context#Context.Value

When a built-in type is used as the key, different packages collide into the same slot as long as they use the same string or integer value, overwriting each other's reads and writes with no compile-time hint; an unexported type can only be constructed by its own package, so cross-package collisions are impossible, values can only be retrieved through functions the package provides, and the returned type is guaranteed by that function.

**Good**

```go
type key int

const userKey key = 0

func WithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}
```

**Bad**

```go
ctx = context.WithValue(ctx, "user", u)
```

**References**

- https://pkg.go.dev/context#Context.Value
- https://staticcheck.dev/docs/checks/#SA1029

**Detection**: golangci-lint staticcheck(SA1029) (SA1029 detects the use of a built-in type as a context key.)

### 【SHOULD】CTX-004 Put only request-scoped data in context.Value, not optional parameters.

- Category: Interfaces and Design/Passing context
- Since: Go 1.7
- Tags: context

Only request data that transits processes and APIs, such as user identity, tracing information, and deadlines, goes into the context; a function's own optional parameters and configuration items go through explicit parameters or an options struct.

**Why**

> Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.
>
> — https://pkg.go.dev/context#pkg-overview

A parameter stuffed into the context disappears from the function signature, so reading the signature does not reveal which inputs the function depends on; a caller who omits it only gets nil at runtime, and the point of failure is far from the real cause. With explicit parameters or an options struct, dependencies are written on the signature, and both the compiler and reviewers can check them on the spot.

**Good**

```go
// request-scoped data: user identity goes into context
ctx := context.WithValue(r.Context(), userKey, user)

// optional parameters: template and render options are passed via the function signature
func Render(ctx context.Context, tmpl *Template, opts Options) error {
	u, _ := ctx.Value(userKey).(*User)
	return tmpl.Execute(u, opts)
}
```

**Bad**

```go
// optional parameters: the template is stuffed into context and invisible in the function signature
ctx = context.WithValue(ctx, tmplKey, tmpl)

func Render(ctx context.Context) error {
	tmpl, _ := ctx.Value(tmplKey).(*Template)
	return tmpl.Execute(nil, Options{})
}
```

**References**

- https://pkg.go.dev/context#pkg-overview

**Detection**: manual review (Verify whether what is put in context.Value is request-scoped data.)

## (5) Safe Struct Shape

### 【MUST】STRUCT-001 Structs containing sync.Mutex or other synchronization primitives must not be copied by value; always pass or receive them by pointer.

- Category: Interfaces and Design/Safe Struct Shape
- Since: Go 1.0
- Tags: concurrency

Once a synchronization primitive is copied, the copies no longer share the same lock, and the mutual exclusion semantics silently fail. Define methods on structs that contain a lock with a pointer receiver, and state in the doc comment that they must not be copied.

**Why**

> copylocks check for locks erroneously passed by value
>
> — https://pkg.go.dev/cmd/vet

Value copies are legal at compile time, and the copylocks check in `go vet` is the main line of defense; once it is missed, the problem surfaces only under concurrency pressure, making it costly to track down.

**Good**

```go
type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
```

**Bad**

```go
func (c Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
```

**References**

- https://pkg.go.dev/cmd/vet
- https://google.github.io/styleguide/go/decisions

**Detection**: golangci-lint govet(copylocks)

### 【MAY】STRUCT-002 For structs that should not support comparison, use an incomparable field to block == at compile time.

- Category: Interfaces and Design/Safe Struct Shape
- Since: Go 1.0
- Tags: api-design

For structs that contain internal state, caches, or pointer fields and whose == cannot compare business identity, embed a helper type named `[0]func()`; from then on, both == and using it as a map key no longer compile. This affects only comparison syntax and takes no extra memory.

**Why**

> DoNotCompare can be embedded in a struct to prevent comparability.
>
> — https://github.com/protocolbuffers/protobuf-go/blob/master/internal/pragma/pragma.go

When all fields are comparable, == is legal at compile time and compares field-by-field content; such comparison is incidental rather than intentional and holds only under certain data combinations. After embedding an incomparable field, misuse fails to compile outright, at the cost of a zero-length array that does not change the struct's size or the meaning of its fields.

**Good**

```go
type DoNotCompare [0]func()

type Session struct {
	DoNotCompare
	id string
}
```

**Bad**

```go
type Session struct {
	id   string
	seen bool
}

// all fields are comparable, so == compiles, but it compares internal state rather than business identity
if s1 == s2 {
	// ...
}
```

**References**

- https://github.com/protocolbuffers/protobuf-go/blob/master/internal/pragma/pragma.go
- https://go.dev/ref/spec#Comparison_operators

**Detection**: manual review (Verify whether exported structs that intentionally prohibit comparison have embedded an incomparable field.)

