# 6. Performance

### 【MUST】PERF-004 When the number of elements is known or can be estimated, specify the slice capacity with make.

- Category: Performance
- Since: Go 1.0

Use `make([]T, 0, n)` to initialize a slice that will be appended to, with n set to the known length or a reasonable upper bound. It is not required when the count cannot be estimated at all, or the slice is usually very small or short-lived.

**Why**

> If your guess for the number of tasks was a good one, then there's only one allocation site in this program. The make call allocates a slice backing store of the correct size, and append never has to do any reallocation.
>
> — https://go.dev/blog/allocation-optimizations

Without a specified capacity, append starts from capacity 1 and grows step by step, allocating a new backing array and copying the old elements each round, with the old array turning into garbage right after. Once the capacity is accurate, only one allocation remains and both growth and copying are removed; when the capacity is a constant, from Go 1.26 this allocation can also land on the stack, doing away with the heap allocation and its corresponding GC burden as well.

**Good**

```go
tasks := make([]task, 0, len(items))
for _, it := range items {
	tasks = append(tasks, toTask(it))
}
```

**Bad**

```go
var tasks []task
for _, it := range items {
	tasks = append(tasks, toTask(it))
}
```

**References**

- https://go.dev/blog/allocation-optimizations
- https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity

**Detection**: manual review (Check whether the make call for the slice an append targets carries a capacity argument.)

### 【MUST】PERF-005 When the number of elements is known or can be estimated, pass a capacity hint to make(map).

- Category: Performance
- Since: Go 1.0

Use `make(map[K]V, n)`, with n set to the known element count or a reasonable upper bound. It is not required when the count is unknown, or the map usually holds only a few elements. When the element set is fixed, write it out in one map literal rather than calling make first and assigning entry by entry.

**Why**

> Providing a capacity hint to `make()` tries to right-size the map at initialization time, which reduces the need for growing the map and allocations as elements are added to the map.
>
> — https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity

Growing a map means re-bucketing and moving the existing keys and values into new buckets, so one growth allocates and moves more than a slice does. A capacity hint brings the initial bucket count close to the actual size and reduces the number of growths. The Uber text also states that it is only an approximation and does not guarantee a single allocation, so there is no need to force a value when the size clearly cannot be estimated.

**Good**

```go
m := make(map[string]int, len(items))
for _, it := range items {
	m[it.Name] = it.Count
}
```

**Bad**

```go
m := make(map[string]int)
for _, it := range items {
	m[it.Name] = it.Count
}
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#prefer-specifying-container-capacity

**Detection**: manual review (Check whether a map written entry by entry carries a capacity hint.)

### 【SHOULD】PERF-001 Before optimizing, locate the hotspot with a benchmark or profile; do not change code on intuition.

- Category: Performance
- Since: Go 1.24

Before changing performance-related code, use `go test -bench`, `go tool pprof`, or `runtime/pprof` to get data and confirm which section the hotspot falls in. A rewrite without data to support it does not enter the code: it may land on a path that is rarely taken, or sacrifice the readability of the common path for a rare case.

**Why**

> These tools can help you to identify various types of hotspots (CPU, IO, memory), hotspots are the places that you need to concentrate on in order to significantly improve performance.
>
> — https://go.dev/wiki/Performance

Performance problems only surface at runtime; reading code cannot tell which section is actually executed most. Optimizing on intuition often lands on the path you assume is slow but that is in fact rarely taken, so the change accomplishes nothing and leaves behind code that sacrificed readability for optimization. A benchmark and a profile give each function's share of time and allocations, concentrating the limited changes on the real hotspot.

**Good**

```go
func BenchmarkRender(b *testing.B) {
	for b.Loop() {
		_ = render(m)
	}
}

// go test -bench=. -cpuprofile=cpu.out
// go tool pprof -top cpu.out
```

**Bad**

```go
// no measurement at all, guessed 'this should be slow' and rewrote render's loop into a hand-written buffer
```

**References**

- https://go.dev/wiki/Performance

**Detection**: manual review (Check whether a performance change is accompanied by benchmark results or profile conclusions.)

### 【SHOULD】PERF-002 Use strconv, not fmt, to convert between primitive types and strings.

- Category: Performance
- Since: Go 1.0

Convert between int, uint, float, bool, and strings with `strconv` functions such as `Itoa`, `Atoi`, `FormatInt`, `ParseFloat`, and `FormatBool`. Use `fmt.Sprintf` only when multiple values and a complex format (such as a `%v` combination, width, and precision) are genuinely needed.

**Why**

> When converting primitives to/from strings, `strconv` is faster than `fmt`.
>
> — https://github.com/uber-go/guide/blob/master/style.md#prefer-strconv-over-fmt

`fmt`'s formatting goes through interfaces and reflection on a general parsing path, allocating more and checking more than `strconv`'s dedicated functions. In Uber's benchmark the same conversion costs 143 ns/op and two allocations with `fmt.Sprint`, against 64.2 ns/op and one allocation with `strconv.Itoa`; on a high-frequency path the gap is amplified as it accumulates with the call count.

**Good**

```go
s := strconv.Itoa(n)
```

**Bad**

```go
s := fmt.Sprint(n)
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#prefer-strconv-over-fmt

**Detection**: manual review (Check whether conversions between a single primitive type and a string use fmt.)

### 【SHOULD】PERF-003 Convert a fixed string to []byte once outside the loop and reuse the result.

- Category: Performance
- Since: Go 1.22

When a loop or hot path repeatedly converts the same string literal or variable to `[]byte`, hoist the conversion out of the loop and reuse the result of that single conversion. When the string's content changes each iteration, or a slice that can be modified independently is needed, keep converting inside the loop.

**Why**

> Do not create byte slices from a fixed string repeatedly. Instead, perform the conversion once and capture the result.
>
> — https://github.com/uber-go/guide/blob/master/style.md#avoid-repeated-string-to-byte-conversions

Converting between string and `[]byte` allocates new memory and copies byte by byte. A fixed string converted inside a loop produces a slice that is discarded immediately each iteration, with allocations on the order of the loop count and the same amount of garbage added for the GC. In Uber's benchmark one conversion per iteration costs 22.2 ns/op, dropping to 3.25 ns/op once hoisted out of the loop.

**Good**

```go
data := []byte("Hello world")
for range n {
	w.Write(data)
}
```

**Bad**

```go
for range n {
	w.Write([]byte("Hello world"))
}
```

**References**

- https://github.com/uber-go/guide/blob/master/style.md#avoid-repeated-string-to-byte-conversions

**Detection**: manual review (Check whether the loop contains conversions between string and []byte for a fixed string.)

### 【SHOULD】PERF-006 Use strings.Builder to concatenate strings in a loop or across many operations, and call Grow first when the capacity can be estimated.

- Category: Performance
- Since: Go 1.10

Use `strings.Builder` when building a string piece by piece (inside a loop, across many branches, with an indefinite count), and call `Grow` first when the result length can be estimated. Keep using `+` for a small, fixed inline concatenation. Use `bytes.Buffer` when a read-write buffer or `Bytes` access is needed.

**Why**

> A Builder is used to efficiently build a string using Builder.Write methods. It minimizes memory copying. The zero value is ready to use. Do not copy a non-zero Builder.
>
> — https://pkg.go.dev/strings#Builder

Strings are immutable, so every += builds a new string and copies the whole accumulated content; inside a loop the repeated moving grows with the accumulated length, making the total copied close to quadratic. `Builder` writes into one mutable buffer and produces the result only once at retrieval; calling `Grow` before writing removes even the buffer growth.

**Good**

```go
var b strings.Builder
b.Grow(total)
for _, part := range parts {
	b.WriteString(part)
}
s := b.String()
```

**Bad**

```go
var s string
for _, part := range parts {
	s += part
}
```

**References**

- https://pkg.go.dev/strings#Builder

**Detection**: manual review (Check whether the loop contains string concatenation with +=.)

### 【SHOULD】PERF-007 When keeping a small section of a large slice for a long time, use slices.Clone to cut the reference to the backing array.

- Category: Performance
- Since: Go 1.21

`s[i:j]` shares the backing array with the original slice, so while the sub-slice is still referenced the whole array cannot be reclaimed. When the sub-slice lives long, or the original slice is far larger than the section needed, use `slices.Clone` (since Go 1.21) or `make` plus `copy` to produce an independent slice. When the sub-slice is released as the function returns, or covers most of the elements anyway, slicing directly is enough.

**Why**

> As mentioned earlier, re-slicing a slice doesn't make a copy of the underlying array. The full array will be kept in memory until it is no longer referenced. Occasionally this can cause the program to hold all the data in memory when only a small piece of it is needed.
>
> — https://go.dev/blog/slices-intro#a-possible-gotcha

When a small section is cut out of a large file or large response and kept for a long time, the sub-slice still points at the original backing array, the GC treats the whole array as still referenced, and the process's resident memory far exceeds what is actually needed. `Clone` copies out an independent array containing only the needed elements, and the original array can be reclaimed once it loses its references.

**Good**

```go
line := slices.Clone(buf[start:end])
return line
```

**Bad**

```go
line := buf[start:end]
return line
```

**References**

- https://go.dev/blog/slices-intro#a-possible-gotcha

**Detection**: manual review (Check whether a sub-slice cut from a large buffer and kept for a long time is copied.)

### 【SHOULD】PERF-008 Pass arguments by value when the function only reads them through a dereference; do not pass pointers just to save a few bytes.

- Category: Performance
- Since: Go 1.0

A string, an interface value, and a small struct are themselves only a machine word in size, so passing by value and passing a pointer cost about the same to copy. Pass by value when the function body only reads through `*p` and never writes. Pass a pointer when the caller's data must be modified, or the struct is clearly large.

**Why**

> Don't pass pointers as function arguments just to save a few bytes. If a function refers to its argument x only as *x throughout, then the argument shouldn't be a pointer.
>
> — https://go.dev/wiki/CodeReviewComments#pass-values

A string and an interface value are fixed-size headers whose copy cost equals that of passing a pointer, with one less dereference. Passing `*string` or `*io.Reader` to save bytes adds a layer of indirection, and when the pointer escapes the function it can push the data onto the heap, actually increasing allocations.

**Good**

```go
func printName(name string) {
	fmt.Println(name)
}
```

**Bad**

```go
func printName(name *string) {
	fmt.Println(*name)
}
```

**References**

- https://go.dev/wiki/CodeReviewComments#pass-values

**Detection**: manual review (Check whether a parameter that is only read through a dereference uses a pointer.)

