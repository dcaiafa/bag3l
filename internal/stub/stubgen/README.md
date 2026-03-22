# stubgen

stubgen is a code generator that produces Go boilerplate for bag3l builtin
packages. It reads a `.stubgen` definition file and outputs a `.stubgen.go` file
containing:

- Wrapper functions that validate arguments, resolve overloads via type
  switching, apply defaults, and convert between VM and Go types.
- Struct types with `FromMap` methods for converting VM maps to Go structs.
- An `Exports` table that registers all functions and constants with the runtime.

## Architecture

```
.stubgen file
    |
    v
 Parser  (internal/stub/parser2)
    |
    v
  AST    (internal/stub/ast)
    |
    v  two passes: Check -> Emit
Analysis (internal/stub/analysis)
    |
    v
 gofmt
    |
    v
.stubgen.go file
```

**Parser** reads the `.stubgen` file and builds an AST.

**Check pass** walks the AST and populates an `Analysis` struct: resolves type
references, registers functions (merging overloads), structs, constants, and
imports.

**Emit pass** is a no-op on the AST side; after both passes complete, the
generator calls `Analysis.Emit()` which writes all Go code into a buffer.

**gofmt** formats the output before writing to disk.

Entry point: `internal/stub/stubgen/main.go`. The generator runs two AST passes
via `generator.go`, then emits and formats.

## `.stubgen` syntax

### Package declaration

Every `.stubgen` file starts with a package declaration. The name is quoted:

```
package "mypackage"
```

### Imports

Import a Go package with an alias. The alias is used in constant expressions
with a `#` prefix:

```
import gotime "time"
```

### Type declarations

Map a stubgen type name to a Go type. The Go type lives in the same package
unless a full import path is given:

```
type File *File
type Time Time
type Duration Duration
type Location *Location
```

The `*` indicates a reference (pointer) type. If the Go type is in another
package:

```
type Foo "github.com/example/pkg".Foo
```

### Struct declarations

Declare a struct that will be deserialized from a VM `Map`. stubgen generates
the Go struct and a `FromMap(*vm.Map) error` method:

```
struct OpenOptions {
  read Bool
  write Bool
  append Bool
  create Bool
}
```

Field names are snake_case in the definition and converted to PascalCase in
generated Go. Declaring a struct also registers it as a type, so it can be used
as a function parameter.

### Constant declarations

Export a named constant. The value is a Go expression wrapped in backticks.
Reference imported packages with `#alias.Symbol`:

```
const HOUR Duration = `NewDuration(#gotime.Hour)`
```

### Function declarations

```
func name(param Type) ReturnType
func name(p1 Type, p2 Type) (Ret1, Ret2)
func name()
```

**Parameter features:**

| Feature              | Syntax                             |
|----------------------|------------------------------------|
| Required parameter   | `name Type`                        |
| Optional with default| `name Type = \`goExpr\``           |
| Nullable             | `name Type?`                       |
| Variadic             | `name ...Type` (must be last)      |

Default values are Go expressions in backticks. Common patterns:

```
func parse(v Str, layout Str = `"2006-01-02T15:04:05Z07:00"`) Time
func create_temp(pattern Str = `""`, dir Str = `""`) File
func open(name Str, opt OpenOptions = `nil`) File
func from_unix(sec Int, nano Int = `0`) Time
```

Default parameters must be trailing (all defaults at the end of the parameter
list).

**Overloading:** declare the same function name multiple times with different
parameter types. Each declaration becomes a separate signature and stubgen
builds a DFA-based dispatcher:

```
func replace(s Str, old Str, rep Str, n Int = `"-1"`) Str
func replace(s Str, old Regex, rep Str) Str
func replace(s Str, old Regex, rep Callable) Str

func stat(f File) Map
func stat(name Str) Map
```

## Built-in types

These types are pre-registered and available in every `.stubgen` file:

| stubgen type | VM Go type     | Convenience Go type | Ref |
|--------------|----------------|---------------------|-----|
| `Any`        | `vm.Value`     | `vm.Value`          | no  |
| `Bool`       | `vm.Bool`      | `bool`              | no  |
| `Callable`   | `vm.Callable`  | `vm.Callable`       | no  |
| `Float`      | `vm.Float`     | `float64`           | no  |
| `Int`        | `vm.Int`       | `int64`             | no  |
| `Iter`       | `vm.Iterator`  | `vm.Iterator`       | no  |
| `List`       | `*vm.List`     | `*vm.List`          | yes |
| `Map`        | `*vm.Map`      | `*vm.Map`           | yes |
| `Reader`     | `vm.Reader`    | `vm.Reader`         | no  |
| `Regex`      | `*vm.Regex`    | `*vm.Regex`         | yes |
| `Str`        | `vm.String`    | `string`            | no  |
| `Writer`     | `vm.Writer`    | `vm.Writer`         | no  |

**Convenience types** are used in your implementation functions. stubgen
generates the conversion code automatically:

- `Str` parameters arrive as `string` (via `.String()`), return values are
  wrapped with `vm.NewString()`.
- `Int` parameters arrive as `int64` (via `.Int64()`), wrapped with
  `vm.NewInt()`.
- `Float` parameters arrive as `float64` (via `.Float64()`), wrapped with
  `vm.NewFloat()`.
- `Bool` parameters arrive as `bool` (via `.Bool()`), wrapped with
  `vm.NewBool()`.
- `Iter` parameters arrive as `vm.Iterator` (via `stub.MustMakeIter()`),
  accepting both `Iterable` and `Iterator` values.
- `Reader` parameters arrive as `vm.Reader` (via `stub.MustMakeReader()`),
  accepting both `Reader` and `Readable` values.
- All other types pass through unchanged.

## Writing implementation functions

For each function declared in the `.stubgen` file, you write Go implementation
functions with a numeric suffix corresponding to the overload index (0-based, in
declaration order).

### Function signature

```go
func name0(vm *vm.VM, <params...>) (<returns...>, error) {
```

- First parameter is always `*vm.VM`.
- Subsequent parameters use convenience types (see table above).
- Last return value is always `error`.
- If the stubgen declaration has no return type, only return `error`.

### Examples

Given:

```
func batch(iter Iter, n Int) Iter
```

Implement:

```go
func batch0(m *vm.VM, iter vm.Iterator, n int64) (vm.Iterator, error) {
    // ...
}
```

Given overloads:

```
func stat(f File) Map
func stat(name Str) Map
```

Implement:

```go
func stat0(vm *vm.VM, f *File) (*vm.Map, error) {
    // File variant.
}

func stat1(vm *vm.VM, name string) (*vm.Map, error) {
    // String variant.
}
```

Given a struct parameter:

```
func open(name Str, opt OpenOptions = `nil`) File
```

The struct parameter arrives as a pointer (or nil when the default is `nil`):

```go
func open0(vm *vm.VM, name string, opts *OpenOptions) (*File, error) {
    if opts == nil {
        // Use defaults.
    }
    // opts.Read, opts.Write, opts.Append, etc.
}
```

### Functions with no return value

```
func close(v Any)
```

```go
func close0(m *vm.VM, v vm.Value) error {
    // Return only error.
}
```

### Multiple return values

```
func foo(x Int) (Str, Int)
```

```go
func foo0(m *vm.VM, x int64) (string, int64, error) {
    return "hello", 42, nil
}
```

## Generated code

For each function, stubgen generates a wrapper `_name` with the VM calling
convention signature:

```go
func _name(vm *VM, args []Value, nret int) ([]Value, error)
```

The wrapper:

1. Checks argument count against the minimum required.
2. Type-switches on each argument position using a DFA built from all
   overload signatures.
3. Applies default values for missing optional parameters.
4. Converts arguments from VM types to convenience types.
5. Calls the numbered implementation function (e.g., `name0`, `name1`).
6. Converts return values back to VM types.
7. Returns `ErrInsufficientArgs`, `ErrTooManyArgs`, or `InvalidArg` for
   mismatches.

An `Exports` variable is generated at the end:

```go
var Exports = export.Exports{
    {N: "batch",  T: export.Func,  F: _batch},
    {N: "filter", T: export.Func,  F: _filter},
    {N: "HOUR",   T: export.Value, V: NewDuration(gotime.Hour)},
}
```

## Workflow: adding a new builtin

1. **Edit the `.stubgen` file.** Add your function, type, struct, or constant
   declaration.

2. **Run code generation.** From the package directory:

   ```
   go generate
   ```

   This invokes the `//go:generate stubgen <file>.stubgen` directive in the
   package's `.go` file, producing the updated `.stubgen.go`.

3. **Implement the function.** Create or edit a `.go` file in the same package.
   Write the `name0` (and `name1`, `name2`, ... for overloads) functions
   following the conventions above.

4. **Build and test.**

### Example: adding `to_title` to the `str` package

**Step 1** -- add to `lib/str/str.stubgen`:

```
func to_title(s Str) Str
```

**Step 2** -- regenerate:

```
cd lib/str && go generate
```

**Step 3** -- implement in `lib/str/str.go`:

```go
func to_title0(m *vm.VM, s string) (string, error) {
    return strings.ToTitle(s), nil
}
```

### Example: adding an overloaded function with struct parameter

**Step 1** -- define struct and functions in your `.stubgen`:

```
struct QueryOptions {
  timeout Int
  retries Int
}

func query(url Str) Str
func query(url Str, opts QueryOptions) Str
```

**Step 2** -- regenerate. stubgen produces the `QueryOptions` struct with
`FromMap`, plus a dispatcher that routes to `query0` or `query1`.

**Step 3** -- implement:

```go
func query0(m *vm.VM, url string) (string, error) {
    return query1(m, url, nil)
}

func query1(m *vm.VM, url string, opts *QueryOptions) (string, error) {
    timeout := 30 // default
    if opts != nil {
        timeout = int(opts.Timeout)
    }
    // ...
}
```

## The `//go:generate` directive

Every package that uses stubgen needs a `//go:generate` directive in one of its
`.go` files (typically the main package file):

```go
//go:generate stubgen mypackage.stubgen
```

If `stubgen` is not on your `PATH`, use a relative path to the tool:

```go
//go:generate go run ../../../internal/stub/stubgen mypackage.stubgen
```

## File organization

A typical stubgen-based package looks like:

```
lib/mypackage/
  mypackage.go          # //go:generate directive, types, main implementations
  mypackage.stubgen      # stubgen declarations
  mypackage.stubgen.go   # generated (do not edit)
  feature_a.go           # additional implementation files
  feature_b.go           # ...
```

The `.stubgen.go` file is auto-generated and should not be edited by hand. All
implementation functions can live in any `.go` file within the package.
