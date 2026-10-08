# 2. Errors and Logging

## (1) Error Handling

### 【MUST】ERR-001 Wrap with %w when the error chain must be preserved; do not use %v or %s.

- Category: Errors and Logging/Error Handling
- Since: Go 1.13
- Tags: error-handling

When the caller needs to test the error with `errors.Is` or `errors.As`, wrap it with %w and add this layer's context. Use %v only at the boundary that truly terminates the error chain (output to the outside, writing logs).

**Why**

> In Go 1.13, the fmt.Errorf function supports a new %w verb. When this verb is present, the error returned by fmt.Errorf will have an Unwrap method returning the argument of %w, which must be an error. In all other ways, %w is identical to %v.
>
> — https://go.dev/blog/go1.13-errors

%v cuts the error chain, so the caller can no longer test the error by type and can only compare strings, which breaks as soon as the wording changes.

**Good**

```go
if err := repo.Save(ctx, order); err != nil {
	return fmt.Errorf("save order %s: %w", order.ID, err)
}
```

**Bad**

```go
if err := repo.Save(ctx, order); err != nil {
	return fmt.Errorf("save order %s: %v", order.ID, err)
}
```

**References**

- https://go.dev/blog/go1.13-errors
- https://google.github.io/styleguide/go/decisions

**Detection**: golangci-lint errorlint

### 【MUST】ERR-003 Test the kind of an error with errors.Is and errors.As, not with == or a type assertion.

- Category: Errors and Logging/Error Handling
- Since: Go 1.13
- Tags: error-handling

A bare comparison and a type assertion only look at the outermost layer of the error chain: once the error is wrapped with %w, == no longer holds and `err.(*MyError)` fails outright. `errors.Is` compares layer by layer against sentinel errors, and `errors.As` matches layer by layer by error type; both unwrap the error chain. The sentinel errors and error types available for testing are exported by the package that produces the error; compare with == only when the error is known never to be wrapped.

**Why**

> In the simplest case, the errors.Is function behaves like a comparison to a sentinel error, and the errors.As function behaves like a type assertion. When operating on wrapped errors, however, these functions consider all the errors in a chain.
>
> — https://go.dev/blog/go1.13-errors

Once some layer changes a direct return into a %w wrap, == and type assertions silently stop matching, the error branch is never taken, and the problem usually surfaces only in production; errorlint covers both forms, and following it blocks them automatically.

**Good**

```go
if errors.Is(err, os.ErrNotExist) {
	// file does not exist
}
```

**Bad**

```go
if err == os.ErrNotExist {
	// file does not exist
}
```

**References**

- https://go.dev/blog/go1.13-errors
- https://pkg.go.dev/errors#Is

**Detection**: golangci-lint errorlint (Covers both the == comparison and the type assertion forms)

### 【MUST】ERR-006 Handle the same error only once: wrap and return it, or log it in place, never both.

- Category: Errors and Logging/Error Handling
- Since: Go 1.13
- Tags: error-handling

Handle an error at the layer that can decide what to do about it: if a caller above can still respond, add context and return it directly; if it has reached the boundary and no one else will take it, log it in place and let the flow degrade. Do not write a log entry and then return err; when the same place both logs and returns, the layer above usually logs it again, and one failure shows up as several entries in the logs.

**Why**

> You should only handle errors once. Handling an error means inspecting the error value, and making a single decision. ... But making more than one decision in response to a single error is also problematic.
>
> — https://dave.cheney.net/practical-go/presentations/qcon-china.html#_only_handle_an_error_once

Duplicate logging multiplies the number of log entries, so when investigating you first have to tell the real number of failures apart from the pile of duplicates; once logging is centralized in one place, the number of log entries corresponds one-to-one with the number of failures.

**Good**

```go
if err := save(order); err != nil {
	return fmt.Errorf("save order %s: %w", order.ID, err)
}
```

**Bad**

```go
if err := save(order); err != nil {
	slog.Error("save order failed", "error", err)
	return err
}
```

**References**

- https://dave.cheney.net/practical-go/presentations/qcon-china.html#_only_handle_an_error_once
- https://go.dev/wiki/CodeReviewComments#handle-errors

**Detection**: manual review (Check each if err != nil branch: writing a log and returning err in the same branch is a violation)

### 【SHOULD】ERR-004 Return an error for expected failures; do not use panic.

- Category: Errors and Logging/Error Handling
- Since: Go 1.16
- Tags: error-handling

Reserve panic for cases where the caller cannot continue, such as missing required configuration at startup or a branch that is logically unreachable being reached. Return an error for expected failures (invalid arguments, a resource that does not exist, an error returned by a downstream dependency) and let the caller decide how to handle it. Library code does not panic; recover is only a last resort at the process entry point.

**Why**

> Don't use panic for normal error handling. Use error and multiple return values.
>
> — https://go.dev/wiki/CodeReviewComments#dont-panic

panic bypasses the caller's error branches and turns a failure that could have been degraded in place into an interruption of the whole call stack; recover can only catch a panic in the same goroutine, while a panic in another goroutine terminates the process directly and cannot be caught.

**Good**

```go
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parse(data)
}
```

**Bad**

```go
func Load(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return parse(data)
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#dont-panic
- https://go.dev/doc/effective_go#errors

**Detection**: manual review (Check whether panic is used for expected failures)

### 【SHOULD】ERR-005 Return failures as an error or an ok value; do not signal them with a special value.

- Category: Errors and Logging/Error Handling
- Since: Go 1.0
- Tags: error-handling

Do not use -1, an empty string or a zero value to express "not found" or "failed"; return an error instead, or add an ok return value when no reason needs to be given. The caller must check ok or err before using the result, and a missed check fails to compile on the spot. This clause does not apply when nil or 0 in the return values is itself a valid result.

**Why**

> Go's support for multiple return values provides a better solution. Instead of requiring clients to check for an in-band error value, a function should return an additional value to indicate whether its other return values are valid.
>
> — https://go.dev/wiki/CodeReviewComments#in-band-errors

A special value lives in the return values alongside the function's legitimate results, so when the caller misses the check the error flows on through the computation and blows up somewhere else, pointing the investigation at the wrong place; multiple return values turn a missed check into a compile error.

**Good**

```go
func Lookup(id string) (Order, bool) {
	// ...
}

order, ok := Lookup(id)
if !ok {
	return ErrNotFound
}
```

**Bad**

```go
func Lookup(id string) Order {
	return Order{} // return zero value when not found
}

order := Lookup(id) // the caller cannot tell not-found from a zero value
```

**References**

- https://go.dev/wiki/CodeReviewComments#in-band-errors

**Detection**: manual review (Check whether -1, an empty value or a zero value is used to express failure)

## (2) Error Codes

### 【SHOULD】CODE-001 For external error codes, reuse the existing standard classification; do not build a separate numbering scheme.

- Category: Errors and Logging/Error Codes
- Since: Go 1.0
- Tags: api-design

When a service returns errors to the outside, express the broad class with the standard classification of its protocol (HTTP status codes, gRPC status codes) and put the finer reason in the error details or in a custom error type. When the business really has a broad class that the standard classification does not cover, add a reason field on top of the standard code; do not invent a numbering scheme that sits alongside the standard codes.

**Why**

> Google APIs must use the canonical error codes defined by google.rpc.Code. Individual APIs must avoid defining additional error codes, since developers are very unlikely to write logic to handle a large number of error codes.
>
> — https://cloud.google.com/apis/design/errors

A self-made numbering scheme forces every caller to write a mapping and a test for each set of codes, and once there are many codes callers in practice handle only the broad class, so the fine-grained part is wasted effort; the standard classification comes with mappings in clients for every language and with established retry semantics, so reusing it is direct.

**Good**

```go
http.Error(w, "order not found", http.StatusNotFound)
```

**Bad**

```go
http.Error(w, "1001", http.StatusOK)
```

**References**

- https://cloud.google.com/apis/design/errors

**Detection**: manual review (Check whether the external error codes place self-made numbers alongside the standard classification)

### 【SHOULD】CODE-002 Error messages are not a contract; test errors only by code and type.

- Category: Errors and Logging/Error Codes
- Since: Go 1.13
- Tags: api-design, error-handling

Error messages address people and their wording can change at any time: callers, tests and alerting rules should all test by error code or error type, not by message text. When localized text must be shown to users, put the text in a separate field or in the error details, kept apart from the message shown to developers.

**Why**

> Error messages are not part of the API surface. They are subject to changes without notice. Application code must not have a hard dependency on error messages.
>
> — https://cloud.google.com/apis/design/errors

A branch that tests by message text silently stops matching after the wording is changed, for example when a parameter name is added and it falls into the default branch; error codes and error types are stable, and testing the kind of an error should look at them.

**Good**

```go
if errors.Is(err, ErrNotFound) {
	// branch by error type
}
```

**Bad**

```go
if strings.Contains(err.Error(), "not found") {
	// branch by message text
}
```

**References**

- https://cloud.google.com/apis/design/errors

**Detection**: manual review (Check whether there is a branch that tests by error message text)

### 【SHOULD】CODE-003 When propagating an error across services, translate the error code; do not pass through the downstream raw error.

- Category: Errors and Logging/Error Codes
- Since: Go 1.0
- Tags: api-design

After receiving an error from a downstream service, judge its meaning and then return it to the outside under this service's error codes: a downstream internal fault becomes unavailable or internal, and an argument error is kept only when it was truly caused by this service's caller. Write the downstream raw error and error chain into the local logs, not into the external response.

**Why**

> If your API service depends on other services, you should not blindly propagate errors from those services to your clients. When translating errors, we suggest the following: Hide implementation details and confidential information. Adjust the party responsible for the error.
>
> — https://cloud.google.com/apis/design/errors

Passing it through exposes the downstream service name, parameter names and internal structure to the external caller, and it also attributes the fault to the wrong party: the downstream reports an invalid argument, but the outside sees an invalid argument of this service; after translating, the outside sees only this service's semantics, and the details stay in its own logs.

**Good**

```go
if err := callOrder(ctx); err != nil {
	// map the downstream failure to unavailable, write the original error to the local log
	return status.Error(codes.Unavailable, "order service unavailable")
}
```

**Bad**

```go
if err := callOrder(ctx); err != nil {
	return err // pass through the downstream error directly
}
```

**References**

- https://cloud.google.com/apis/design/errors

**Detection**: manual review (Check whether an error returned to the outside comes directly from a downstream service)

## (3) Logging

### 【MUST】LOG-001 Record business logs with log/slog; do not print directly with fmt or log.

- Category: Errors and Logging/Logging
- Since: Go 1.21
- Tags: logging

Route structured logging uniformly through `log/slog`: assemble the `Handler` with `slog.New` at the process entry point, use `JSONHandler` in production, and use an injected `*slog.Logger` or `slog.Default()` everywhere else. Use `fmt.Print` and the standard `log` package only for the output of one-off scripts; when adopting a third-party logging library, funnel it into the same output through a custom `Handler`.

**Why**

> Structured logs use key-value pairs so they can be parsed, filtered, searched, and analyzed quickly and reliably.
>
> — https://go.dev/blog/slog

A whole line of text printed by `fmt` or `log` cannot be parsed by field, so the collector can only store it as one block, and searching by order_id or aggregating by latency requires writing separate parsing rules; the key-value pairs and JSON output of `slog` can be recognized directly by the collector.

**Good**

```go
slog.Info("order created", "order_id", order.ID)
```

**Bad**

```go
fmt.Printf("order created: %s\n", order.ID)
```

**References**

- https://go.dev/blog/slog
- https://pkg.go.dev/log/slog

**Detection**: grep regex \b(fmt\.Print|log\.Print) (Check by hand whether the match is a business log)

### 【SHOULD】LOG-002 Write log fields as key-value pairs; do not splice them into the message text.

- Category: Errors and Logging/Logging
- Since: Go 1.21
- Tags: logging

Carry every value that accompanies a log entry (ID, count, latency, error) as a key-value argument, and write only what happened in the message. Use the same key name for the same meaning throughout the repository; pass an error value as `slog.Any("error", err)` rather than writing it into the message.

**Why**

> Unlike with the log package, we can easily add key-value pairs to our output by writing them after the message.
>
> — https://go.dev/blog/slog

Once a value is spliced into the message, the collector can only fuzzy-match the whole line and can neither search nor aggregate precisely by field; a key-value pair becomes a separate field at the collector and can be filtered and counted directly.

**Good**

```go
slog.Info("order created", "order_id", order.ID, "amount", order.Amount)
```

**Bad**

```go
slog.Info(fmt.Sprintf("order %s created, amount %d", order.ID, order.Amount))
```

**References**

- https://go.dev/blog/slog

**Detection**: manual review (Check whether log values are spliced into the message string)

### 【SHOULD】LOG-003 Use the Context-carrying methods for logs on the request path.

- Category: Errors and Logging/Logging
- Since: Go 1.21
- Tags: context, logging

When handling a request, pass ctx together with `slog` methods such as `InfoContext` and `ErrorContext`, so the handler can fill in fields that run through the whole request, such as the request ID and trace identifier. Use the methods without Context on paths that have no request context, such as startup and shutdown.

**Why**

> As the call to LogAttrs shows, you can pass a context.Context to some log functions so a handler can extract context information like trace IDs. (Canceling the context does not prevent the log entry from being written.)
>
> — https://go.dev/blog/slog

When a log entry carries no ctx, the several entries produced by one request cannot be tied together automatically at the collector, and investigating first requires pulling a batch by time window and telling them apart by hand; with ctx, the request-level fields are filled in uniformly by the handler, with no need to write them out at every place.

**Good**

```go
slog.ErrorContext(ctx, "order save failed", "error", err)
```

**Bad**

```go
slog.Error("order save failed", "error", err)
```

**References**

- https://pkg.go.dev/log/slog#Logger.InfoContext
- https://go.dev/blog/slog

**Detection**: manual review (Check whether logs on the request path pass in ctx)

