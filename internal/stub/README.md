# stubgen

stubgen generates Go boilerplate code for writing bag3l builtin packages. It handles argument parsing, type checking, function overload resolution, and creates the export table.

## Usage

Add a go:generate directive to your package:

```go
//go:generate stubgen mypackage.stubgen
```

Then run `go generate` to create `mypackage.stubgen.go`.

## Language Reference

### Package Declaration

Every stubgen file must start with a package declaration:

```
package "packagename"
```

### Imports

Import Go packages for use in constant expressions:

```
import gotime "time"
import libtime "github.com/dcaiafa/bag3l/lib/time"
```

The alias is used with the `#` prefix in constant expressions.

### Type Declarations

Map a bag3l type name to a Go type:

```
type File *File           # Pointer to local File struct
type Time Time            # Local Time struct (non-pointer)
type Duration Duration    # Local Duration struct
```

For external types, use the full package path:

```
type Time "github.com/dcaiafa/nitro/lib/time".Time
```

### Struct Declarations

Declare structs that convert from bag3l Maps:

```
struct OpenOptions {
  read Bool
  write Bool
  append Bool
  create Bool
}
```

This generates a Go struct with a `FromMap` method that parses a `*vm.Map` into the struct. Field names in bag3l use snake_case while the generated Go struct uses PascalCase.

### Constant Declarations

Export constant values:

```
const HOUR Duration = `NewDuration(#gotime.Hour)`
const MINUTE Duration = `NewDuration(#gotime.Minute)`
```

The value is a raw string containing a Go expression. Use `#alias` to reference imported packages.

### Function Declarations

Declare exported functions:

```
func create(name Str) File
func now() Time
func parse(v Str, layout Str = `"2006-01-02T15:04:05Z07:00"`) Time
```

#### Parameters

- **Required**: `name Type`
- **Optional with default**: `name Type = \`default_expr\``
- **Variadic**: `name ...Type`
- **Nullable**: `name Type?` (the `?` allows nil values)

#### Return Types

- **Single return**: `func foo() Type`
- **Multiple returns**: `func foo() (Type1, Type2)`
- **No return**: `func foo()`

### Function Overloading

Multiple functions with the same name but different signatures create overloads:

```
func stat(f File) Map
func stat(name Str) Map

func replace(s Str, old Str, rep Str, n Int = "-1") Str
func replace(s Str, old Regex, rep Str) Str
func replace(s Str, old Regex, rep Callable) Str
```

The generated code dispatches to the correct implementation based on argument types at runtime.

### Built-in Types

| stubgen Type | Go Type |
|--------------|---------|
| `Str` | `string` (converted from `vm.String`) |
| `Int` | `int64` (converted from `vm.Int`) |
| `Bool` | `bool` (converted from `vm.Bool`) |
| `Any` | `vm.Value` |
| `Map` | `*vm.Map` |
| `Iter` | `vm.Iterator` |
| `Reader` | `vm.Reader` |
| `Writer` | `vm.Writer` |
| `Callable` | `vm.Callable` |
| `Regex` | `vm.Regex` |

Custom types declared with `type` are used directly.

## Implementation Convention

For each function declared in the stubgen file, implement Go functions with a numeric suffix:

**stubgen:**
```
func stat(f File) Map
func stat(name Str) Map
```

**Go implementation:**
```go
func stat0(vm *vm.VM, f *File) (*vm.Map, error) {
    // Implementation for stat(File)
}

func stat1(vm *vm.VM, name string) (*vm.Map, error) {
    // Implementation for stat(Str)
}
```

Rules:
1. The first parameter is always `*vm.VM`
2. Overloads are numbered starting from 0
3. Return an error as the last return value
4. The VM types are converted to Go types (e.g., `Str` becomes `string`)

## Generated Code

stubgen produces:

1. **Wrapper functions** (`_funcname`) that:
   - Validate argument count
   - Type-switch on arguments for overload resolution
   - Apply default values for optional parameters
   - Convert between VM types and Go types
   - Call the user-implemented function

2. **Struct types** with `FromMap` methods for struct declarations

3. **Exports variable** containing all exported functions and constants:
   ```go
   var Exports = export.Exports{
       {N: "create", T: export.Func, F: _create},
       {N: "HOUR", T: export.Value, V: NewDuration(time.Hour)},
   }
   ```

## Example

**time.stubgen:**
```
package "time"

import gotime "time"

type Time Time
type Duration Duration

const HOUR Duration = `NewDuration(#gotime.Hour)`
const SECOND Duration = `NewDuration(#gotime.Second)`

func now() Time
func parse(v Str, layout Str = `"2006-01-02T15:04:05Z07:00"`) Time
func format(t Time, layout Str = `"2006-01-02T15:04:05Z07:00"`) Str
```

**time.go:**
```go
package time

import "time"

//go:generate stubgen time.stubgen

type Time struct {
    time time.Time
}

func NewTime(t time.Time) Time {
    return Time{time: t}
}

func now0(m *vm.VM) (Time, error) {
    return NewTime(time.Now()), nil
}

func parse0(m *vm.VM, v string, layout string) (Time, error) {
    t, err := time.Parse(layout, v)
    if err != nil {
        return Time{}, err
    }
    return NewTime(t), nil
}

func format0(m *vm.VM, t Time, layout string) (string, error) {
    return t.time.Format(layout), nil
}
```
