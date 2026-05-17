# Virtual Machine

The Bag3l VM is a stack-based bytecode interpreter implemented in Go in
`internal/vm`. It executes programs produced by `internal/compiler` and
provides the runtime services those programs rely on: values, calls, closures,
iterators, exceptions, `defer`, and cooperative coroutines.

This document describes the runtime architecture: the data structures, the
calling convention, the instruction set, and how the moving parts fit together.

## File map

| File | Role |
| --- | --- |
| `vm.go` | `VM` type, the main interpreter loop, call/iter dispatch |
| `coroutine.go` | Per-coroutine state: stack, call stack, frame pool, IP |
| `instr.go` | `Instr` struct and the `OpCode` enum |
| `emitter.go` | `Emitter` used by the compiler to produce bytecode |
| `program.go` | `Program` and `CompiledPackage` containers |
| `fn.go` | `Fn` — a compiled Bag3l function (bytecode + literals + metadata) |
| `closure.go` | `Closure` — `Fn` bound to a captured environment |
| `native_fn.go` | `NativeFn` — Go function exposed as a `Callable` value |
| `iterator.go` | `Iterator`, `ILIterator` (bytecode generator), `NativeIterator` |
| `value.go` | `Value` interface, `Operable`, `Indexable`, `ValueRef` |
| `bool.go` / `int.go` / `float.go` / `string.go` | Primitive value types |
| `list.go` / `object.go` | `List` and `Map` collection types |
| `regex.go` | `Regex` value wrapping `regexp.Regexp` |
| `reader.go` / `writer.go` | `Reader`/`Writer` interfaces used by I/O builtins |
| `runtime_error.go` | `RuntimeError`, sentinel errors, recoverability flag |
| `location.go` | Line-table entry (`Location`) for stack-trace mapping |
| `util.go` | Small helpers |
| `nil.go` | Empty placeholder file |

## Program model

A `Program` is a slice of `*CompiledPackage`, ordered so that
`Program.Packages[len-1]` is the main package
(`Program.MainPackage()`). `cmd/bag3l` constructs a `Program` via the
compiler and hands it to `vm.NewVM(prog).Run(args)`.

A `CompiledPackage` carries everything needed to execute and link a package:

- `Literals` — every constant the package can refer to by index. Functions
  (`*Fn`), strings, ints and floats all live here. The compiler deduplicates
  strings via `Emitter.stringMap`.
- `Deps` — the dependency packages (referenced by index from `OpLoadLiteral`).
- `numGlobals` / `globals` — the globals array. `NewVM` allocates one slice
  of size `numGlobals` per package up front.
- `params` — declared script parameters (`@param`) used by `SetParam`.
- `MainFnNdx` — index in `Literals` of the package's top-level function.
- `Metadata` — `@info` / `@param` / `@flag` annotations consumed by the CLI.

A `Fn` holds:

- `pkg` — back-pointer to its owning package (used to resolve `Literals` and
  `globals` when a frame is pushed).
- `name` — index into `pkg.Literals` of the function's display name string.
- `instrs` — the function's bytecode.
- `locations` — sorted ascending by `ip`; each entry says "from this `ip`
  onward, source line is X in file Y". Used by `getLocation` to build stack
  traces (binary search would be more apt; the current code does a linear
  scan).
- `minArgs` — the declared minimum-arg count. **Currently never enforced by
  the VM** (see "Bugs and quirks" below).

A `Closure` is just `{fn, caps}` — an `Fn` plus a list of `ValueRef` captures.

## Coroutine and stack model

Execution state is held in a `coroutine`:

```go
type coroutine struct {
    fiber         *fiber.Fiber
    callStack     []*frame
    frame         *frame              // top of callStack (hot path)
    pkg           *CompiledPackage    // current frame's package
    stack         []Value             // operand stack
    globals       []Value             // current package's globals
    instrs        []Instr             // current frame's bytecode
    sp            int                 // operand stack pointer
    ip            int                 // instruction pointer
    preAllocStack [stackSize]Value    // 1000-slot inline backing array
    framePool     []*frame            // freelist for frame structs
    pendingErr    error               // error injected by SignalError
    contextStack  []context.Context
}
```

The operand stack is unified: arguments, locals, captures-references, and
temporaries all live on it. The hot fields (`frame`, `pkg`, `globals`,
`instrs`) are cached out of `frame.fn.pkg` and `frame.fn.instrs` whenever a
frame becomes current, so the interpreter loop only touches `co.*` and not
deeper pointer chains.

`preAllocStack` is `stackSize = 1000`. **There is no overflow check** on stack
pushes — running past the end of the inline array panics with an
out-of-bounds index.

### Frames

```go
type frame struct {
    nRet       int
    nArg       int
    nLocals    int
    iter       *ILIterator     // non-nil if this frame is an iterator body
    fn         *Fn             // nil if this frame is an external call
    extFn      Callable        // non-nil for external function frames
    caps       []ValueRef      // closure / iterator captures
    tryCatches []tryCatch
    defers     []*Closure
    pipeline   bool
    ip         int             // saved across nested PushFrame
    bp         int             // base pointer into the operand stack
}
```

Frames are allocated lazily and recycled through `coroutine.framePool`.
`PushFrame` saves the previous frame's `ip` and switches the cached `pkg`,
`globals`, `instrs`, and `ip`; `PopFrame` zeros the frame, returns it to the
pool, and restores the previous frame's view.

The operand stack layout for a Bag3l function frame at steady state:

```
... | argN-1 | argN-2 | ... | arg0 | local0 | local1 | ... | <temporaries>
    ^                              ^                                       ^
    bp - nArg                      bp                                      sp
```

`OpInitCallFrame` is emitted as the first instruction of every Bag3l
function. It sets `bp = sp`, advances `sp` by `nLocals`, and zeros the local
slots. Subsequent loads use `bp` as the anchor:

| Opcode | Slot |
| --- | --- |
| `OpLoadArg n` / `OpLoadArgRef n` / `OpLoadArgDeref n` | `bp - nArg + n` |
| `OpLoadLocal n` / `OpLoadLocalRef n` / `OpLoadLocalDeref n` | `bp + n` |
| `OpLoadGlobal n` / `OpLoadGlobalRef n` | `globals[n]` |
| `OpLoadCapture n` / `OpLoadCaptureRef n` | `frame.caps[n]` |
| `OpLoadLiteral lit, dep` | `pkg.Deps[dep].Literals[lit]` |

## Values

Every runtime value implements `Value`:

```go
type Value interface {
    String() string
    Type() string
    Traits() Traits
}
```

`Traits` is a small bitset; only `TraitEq` is used today. Types that opt in to
`TraitEq` (`Int`, `Float`, `String`, `Bool`, plus all numeric/struct types
that implement `Operable`) get dispatched through `Operable.EvalOp` for all
binary operators. Types without `TraitEq` fall back to Go interface equality
(`==`) for `OpEq`/`OpNE`, which is pointer equality for reference types like
`*Map` and `*List`.

`OpNE` is always implemented as "negate `OpEq`" — `EvalOp` rewrites `OpNE` to
`OpEq`, calls the operand, and inverts the resulting `Bool`. There is no
separate `EvalOp(OpNE, …)` entry point.

`ValueRef` is `struct{ Ref *Value }`. It is itself a `Value`, which means
references are first-class values that flow through the stack like any other.
A function argument or local is "lifted" by storing a `ValueRef` into its
slot whenever something needs a stable reference to it (assignment via
`OpStore`, capture into a closure via `OpCaptureLocal`/`OpCaptureArg`, an
inner iterator's view of an outer variable).

The container/collection types:

- `*List` — backed by `[]Value`. Negative indices are supported by `Index`
  and `Slice`. `IndexRef` requires the index to be in range.
- `*Map` — `map[Value]*mapNode` plus a doubly linked list of `mapNode`s, so
  iteration is in insertion order. `Put` and `IndexRef` insert at the tail
  if the key is new. `mapIter` walks the list using `GetNext`, which silently
  ends iteration if the current key was deleted mid-loop.
- `String` — immutable value type wrapping `string`. `IndexRef` returns
  an error and a zero `ValueRef`; the caller must inspect the error before
  dereferencing.

## Instruction format

```go
type Instr struct {
    opc OpCode
    op2 uint16
    op1 uint32
}
```

A fixed 8-byte triple. Most opcodes use only `op1`; a few pack additional
fields:

- `OpCall`: `op1` packs `(narg | expandFlag | pipelineFlag)`; `op2` is `nret`.
  Flags: `CallExpandFlag = 0x80000000`, `CallPipelineFlag = 0x40000000`,
  `CallArgCountMask = 0x3FFFFFFF`.
- `OpNewIter`: `op1` packs `fnIndex (24 bits low) | iterNRet (8 bits high)`;
  `op2` is the capture count.
- `OpNewClosure`: `op1` is the function literal index; `op2` is the capture
  count.
- `OpObjectGet`: `op2 & OptionalIndexFlag` (`0x0001`) marks an optional
  index (`?.` / `?[]`) so a missing key yields nil instead of an error.
- `OpLoadLiteral`: `op1` is the literal index; `op2` is the dependency index
  into `pkg.Deps` (allowing cross-package access to constants/functions).

## Calling convention

`call(callable, narg, nret, pipeline)` handles all call shapes:

| Callable type | Behavior |
| --- | --- |
| `*Closure` | New frame with `fn = callable.fn`, `caps = callable.caps`; `runFrame` |
| `*Fn` | New frame with `fn = callable`; **`pipeline` is hardcoded to `true`** |
| `*NativeIterator` | Treated as an external call with no captures |
| `*ILIterator` | Resumes the generator (see Iterators) |
| `Callable` | Dispatched through `callExtFn` |
| `nil` | `ErrCannotCallNil` |
| other | error: `cannot call <type>` |

`OpCall` sets up the call from the caller's stack:

1. Decode `narg`/`nret`/`expand`/`pipeline`.
2. If `expand` is set, the top-of-stack value is replaced with its contents:
   `nil` is dropped (narg decremented), `*List` is spread, anything else is
   an error.
3. The callable sits at `sp - narg - 1`; the arguments occupy the slots above
   it.
4. `call` runs to completion. On return, `nret` values are copied down to
   `(originalSp - narg - 1)`, overwriting the callable's slot, and `sp` is
   set so that exactly `nret` values sit at the top.

For Bag3l functions, `call → runFrame → resume` runs the callee's bytecode
synchronously; the interpreter loop is re-entered with the new top frame.
This means the Go call stack mirrors the Bag3l call stack — deep Bag3l
recursion consumes Go stack.

`callExtFn` synthesizes a frame whose `extFn` field carries the callable so
that `GetFrameInfo`/`GetCallerArgs` can still introspect it. The native
function receives an `args []Value` slice that is a *view into the operand
stack*, so mutating the slice's backing array would corrupt the stack. The
function returns `[]Value` rets; they must be at least `nret` long.

Return values for native callables are copied via:

```go
copy(m.co.stack[m.co.sp - narg:], rets[:nret])
m.co.sp = m.co.sp - narg + nret
```

That is, the *callable slot is not overwritten* by `callExtFn`; instead,
`OpCall`'s outer copy moves the rets down by one more slot to consume the
callable too.

## Iterators

There are two iterator implementations behind a common `Iterator` interface:

```go
type Iterator interface {
    Value
    IterNRet() int
    Close(vm *VM) error
    IsClosed() bool
    isIterator()
}
```

### `*ILIterator` — bytecode-defined generators

Created by `OpNewIter` from an `Fn` whose body uses `yield`. The iterator
owns its own private operand stack (`preAllocStack`) plus the resume state:

```go
type ILIterator struct {
    fn         *Fn
    captures   []ValueRef
    iterNRet   int                 // number of values per yield
    tryCatches []tryCatch
    defers     []*Closure
    stack      []Value
    nlocals    int
    ip         int                 // -1 means finished
    sp         int
    closed     bool

    preAllocStack [stackSize]Value
}
```

On every `Next()`, the VM:

1. Saves the host coroutine's `stack`/`sp` (`rstack`, `rsp`).
2. Swaps in the iterator's `stack`/`sp` and pushes a frame whose `ip`,
   `tryCatches`, `defers`, `nlocals`, and `caps` come from the iterator.
3. Runs `runFrame` until `OpIterYield` (yield) or `OpIterRet`/`OpRet` (done).
4. Copies the yielded values back to the host stack and restores the host
   `stack`/`sp`.

`OpIterYield` saves `ip`, `tryCatches`, `defers`, `nlocals` back onto the
`ILIterator` and returns `nil` from the inner interpreter loop. `OpIterRet`
sets `ip = -1` to mark the generator finished. Calling a finished iterator
returns `False` plus nils for the remaining slots.

`Close` only marks the iterator closed; it does *not* run any pending
`defer`s in the iterator body. There is no way for an iterator-local defer to
run a single time at end-of-life (see "Bugs and quirks").

### `*NativeIterator`

Wraps a Go-side `*NativeFn` plus an optional `CloseFn`. Each call to `Next`
invokes the native function; returning an empty slice (`nil, nil`) signals
end-of-iteration and triggers `Close`.

### `OpNext`

The bytecode for `for ... in iter { ... }` looks like:

```
... <ref1> <ref2> ... <refN> <iter>      // sp top
OpNext jumpEndAddr, N
```

`OpNext` calls `iterNext` for `N` return values. On success it stores the
yielded values into the `N` ValueRefs sitting below the iterator on the
stack; on end-of-iteration it jumps to `jumpEndAddr`. Either way it pops back
to `sp - N - 1` (it leaves the iterator on the stack for the next round, or
the jump path pops it implicitly via `sp = rsp`).

### `OpMakeIter`

Coerces the top of stack into an `Iterator`: pass-through for things that
already implement `Iterator`; call `MakeIterator()` for things that implement
`Iterable`; install an empty iterator for `nil`; error otherwise.

## Try / catch / throw

- `OpBeginTry addr` pushes a `tryCatch{CatchAddr: addr}` onto the frame's
  `tryCatches`.
- `OpEndTry endAddr` pops the most recent `tryCatch` and jumps past the
  catch block.
- `OpThrow` pops the top of stack. If it's already a `*RuntimeError`, throw
  it as-is; otherwise wrap it in a new `*RuntimeError` with
  `Recoverable: true` and `ErrValue` set to the thrown value.

When an instruction returns an `error`, `resume` catches it:

1. `wrapRuntimeError(err, recoverable=false)` either reuses an existing
   `*RuntimeError` (preserving its `Recoverable` flag) or wraps the raw
   error with `Recoverable: false`. A wrapped `*NonRecoverableError` always
   forces `Recoverable: false`.
2. If the error is non-recoverable, OR the frame has no `tryCatch` on the
   stack, the frame is popped and the error propagates up.
3. Otherwise, the catch block is entered: `ip = catchAddr`,
   `sp = bp + nLocals + 1`, and the `*RuntimeError` itself is pushed as the
   single value at the top of stack. The catch block reads its fields via
   `RuntimeError.Index` (today only `error` is exposed).

**Recoverability is intentional and asymmetric.** `throw x` is the
language's `error`-equivalent: a recoverable condition the program chose to
signal, catchable by `try`/`catch`, and defers run before propagation.
Runtime-level errors (divide-by-zero, "type not indexable", etc.) are the
panic-equivalent: a bug or invariant violation that leaves the program in
an undefined state, so `try`/`catch` deliberately cannot catch them.
Missing-data conditions that are *not* bugs have their own non-throwing
syntax (`?.`, `?[]`) so users do not reach for `try` as a substitute.

Defers, however, **should still run on the panic path** to release
resources cleanly (file handles, locks, etc.); the current implementation
skips them. See "Bugs and quirks".

## Defer

`OpDefer` pops a `*Closure` and appends it to `frame.defers`. **It type-asserts
to `*Closure` without an `ok` check**; deferring a `*Fn` or `*NativeFn`
panics.

On normal frame return (`OpRet`), `resume` invokes `runDefers`, which calls
each deferred closure in LIFO order. On a *recoverable* error that escapes
the frame, `runDefers` also runs (and any error from a defer is wrapped into
the propagating error message). On a non-recoverable error, defers are
skipped.

## Pipelines

Pipelines are signaled with the `CallPipelineFlag` bit set on `OpCall`'s
`op1`, which propagates to `frame.pipeline`. The VM does not itself
interpret pipelines — it just exposes them so native functions can ask:

- `vm.IsPipeline()` — is the current frame a pipeline call?
- `vm.IsCallerPipeline()` — was the caller invoked as a pipeline?

Native I/O builtins (`io`, `file`, `exec`) use these to decide whether to
chain `Reader`/`Writer`s together vs. produce a one-shot value.

## Coroutines and the fiber scheduler

`VM` is backed by a `fiber.Scheduler`. Each Bag3l coroutine corresponds to
one `*fiber.Fiber` which in turn owns a real goroutine. The scheduler is
*cooperative*: at most one fiber runs logical code at a time; fibers yield
by calling `Block` or `SwitchToNew`.

- `VM.Run(args)` creates the initial coroutine, starts the main function, and
  blocks until the scheduler drains.
- `VM.StartCoroutine(callable)` creates a new coroutine and uses
  `Scheduler.SwitchToNew` to run it immediately, parking the current fiber.
- `VM.Block(f)` is the bridge for blocking native operations. The active
  fiber goes into the `blocked` queue while `f(ctx)` runs concurrently on
  this fiber's goroutine; other ready fibers can be scheduled in the
  meantime.

The single `VM.co` pointer always points to the *currently active* coroutine.
Each entry point that may change which fiber is running is responsible for
re-establishing `m.co`:

- `newCoroutine`'s startup closure assigns `m.co = co` before it begins
  executing.
- `StartCoroutine` saves `prev := m.co`, switches, then restores
  `m.co = prev` after coming back.
- `Block` does the same with `self := m.co`.

The interpreter loop itself never reschedules implicitly — only operations
that go through `Block`/`StartCoroutine`/`SwitchToNew` can move execution to
a different fiber. The loop has no preemption point of its own (apart from
`processInterrupt` checks).

## Interrupts and shutdown

`VM.SignalError(err)` is thread-safe and is the public way to push an error
into the running program from outside the interpreter (e.g. Ctrl-C). It
stores the error in `m.injectedErr`, sets the `interrupt` atomic flag, and
cancels every blocked fiber's context via `Scheduler.CancelBlocked`.

`resumeWithoutRecovery` calls `processInterrupt` at the top of every
iteration and on exit. When the interrupt flag is set, the injected error is
copied to *every* coroutine's `pendingErr` so that whichever fiber resumes
next surfaces it.

`VM.shutdown` runs once after `Run` returns. It sets `shuttingDown = true`
(visible via `VM.ShuttingDown()`) and calls `Close()` on every
`Closer` that registered itself via `RegisterCloser`. `RegisterCloser` /
`UnregisterCloser` are **not synchronized** — they must only be called from
the active fiber.

## Stack traces

`GetStackInfo()` snapshots the current `callStack` into `[]FrameInfo`
records. For Bag3l frames, the source file/line is reconstructed by
`getLocation(fn, ip)`, which scans `fn.locations` to find the largest entry
with `ip <= currentIp`. `wrapRuntimeError` captures a stack trace the first
time a runtime error is wrapped and reuses it on every re-wrap so the trace
reflects the throw site, not propagation sites.

## Compilation entry points

The compiler talks to the runtime only through `Emitter`:

- `NewFn`, `PushFn`, `PopFn`, `SetFuncMinArgs` — define functions and push
  them as the current emit target.
- `NewLabel`, `ResolveLabel`, `EmitJump` — backward/forward jumps. A
  forward-referenced label collects every `instr.op1` slot that needs to be
  patched and writes the resolved address into all of them when
  `ResolveLabel` is called. `EmitJump` panics if asked to jump across
  function boundaries.
- `Emit(pos, op, op1, op2)` — append an instruction. Each emit also updates
  the line table: if `pos.Filename` changed, a new filename literal is
  added; if filename or line number differs from the last entry, a new
  `Location` is appended pointing at the just-emitted instruction's `ip`.
- `AddLiteral` — append to the package's literal pool, deduplicating strings.
- `SetGlobalCount`, `AddGlobalParam` — declare the package's global slots
  and named parameters.
- `ToCompiledPackage` — finalize and hand back the `*CompiledPackage`.

The resulting `Program` carries one `CompiledPackage` per imported module
plus the main package. `VM.NewVM` allocates per-package `globals` slices
before `Run` begins.

## Bugs and quirks

This section catalogues current rough edges discovered while documenting the
VM. They are observable behaviors, not design decisions; treat them as
candidates for fixing rather than as contracts to rely on.

### Confirmed bugs

- **`defer` in iterator bodies runs on every `yield`.** `OpIterYield` saves
  `frame.defers` onto the `ILIterator` and returns `nil` from the inner
  loop. The outer `resume` treats that as a normal return and calls
  `runDefers` before `PopFrame`. Verified with a script that yields three
  times and prints `"defer ran"` four times. `Close` on the iterator also
  does not run defers, so there is no way for an iterator-local `defer` to
  execute exactly once at end-of-life.

- **Two divide-by-zero sentinel errors exist.** `int.go` defines
  `ErrDivByZero` (used by `lib/time` and its tests); `runtime_error.go`
  defines `ErrDivideByZero` (used by `int.go` and `float.go` for the
  division operator). `errors.Is(err, vm.ErrDivByZero)` will not match
  integer division-by-zero errors.

- **Defers do not run on non-recoverable errors.** `resume` only invokes
  `runDefers` on the normal-return path and on the recoverable-error
  propagation path; a non-recoverable error pops the frame without
  running defers. Since runtime errors (divide-by-zero, indexing nil,
  etc.) are intentionally non-recoverable — see "Try / catch / throw" —
  any `defer file.close()` / `defer lock.release()` in the affected
  frame is silently skipped on the way out. The intentional "bugs can't
  be caught" stance is fine, but defers should still unwind for resource
  cleanup (Go's panic+defer model).

- **`OpCall` hardcodes `pipeline=true` for `*Fn`.** The `*Closure` arm of
  `call` honors the `pipeline` parameter; the `*Fn` arm forces it to
  `true`. Bare `*Fn` calls (cross-package function-literal calls) therefore
  always report `pipeline=true` to `IsPipeline` / `IsCallerPipeline`.

- **`OpDefer` panics on non-`*Closure` operands.** The handler does a bare
  type assertion (`.(*Closure)`) with no `ok` check. `defer somefn` where
  `somefn` resolves to `*Fn` or `*NativeFn` panics the Go process instead
  of erroring at the Bag3l level.

### Sub-optimal / fragile patterns

- **No operand-stack overflow check.** `preAllocStack` is a fixed
  `[stackSize]Value` (1000 slots). Push opcodes (`OpDup`, `OpLoad*`, etc.)
  blindly index `stack[sp]`, so deep recursion or pathological local counts
  panic with an out-of-bounds index rather than producing a runtime error.

- **`getLocation` does a linear scan** of `fn.locations` on every
  stack-trace synthesis. Locations are stored sorted by `ip` and should
  use binary search.

- **`wrapRuntimeError` allocates a stack trace on every error**, including
  errors that will be caught by a nearby `try`. `GetStackInfo` allocates a
  `[]FrameInfo` proportional to `callStack` depth.

- **`OpUnaryMinus` has a redundant nil check.** `EvalOp` already rejects
  nil operands for any operation other than `OpEq`/`OpNE`; the manual
  `if term == nil` in the opcode handler is dead.

- **`RegisterCloser` / `UnregisterCloser` are not synchronized.** `VM.mu`
  exists but isn't taken. Works today because only the active fiber calls
  them, but a native function spawning its own goroutine that registers a
  closer would race silently.

- **`coroutine.framePool` is per-coroutine.** Each spawned coroutine
  allocates its own pool, so coroutine-heavy workloads re-pay frame
  allocation costs per fiber.

- **`Fn.minArgs` is wired up but not enforced.** `Emitter.SetFuncMinArgs`
  stores it on the function, but no opcode validates it. `OpLoadArgRef`
  even carries a `// TODO: min arg count` comment. Calls that omit
  required arguments silently observe `nil` for the missing slots.

- **`ILIterator` resume implicitly assumes `bp == 0`.** The first
  invocation runs `OpInitCallFrame` with `sp=0`, setting `bp=0`. On
  subsequent resumes the resume frame's `bp` defaults to zero (from the
  pooled `frame{}`), and "happens to work" because the original value was
  also zero. Storing `bp` on the iterator (alongside `sp`/`ip`/`nlocals`)
  would make this explicit rather than incidental.

- **Naming inconsistency: `OpNewObject` for the map type.** The type was
  renamed from `Object` to `Map` (commit `Rename object => map`), but the
  opcode and several internal identifiers still use the old name.

- **`Bool.EvalOp` does not return `ErrOperationNotSupported`.** It returns
  a plain `fmt.Errorf("bool does not support this operation")`, so the
  `FallbackEvaluator` path inside `EvalOp` can never trigger for bool
  operands.

- **`RuntimeError.Index` errors on unknown fields** instead of returning
  `(nil, false, nil)`. That means `err?.unknown` (optional access) still
  throws, defeating the point of the optional-index flag for this type.

- **`String.IndexRef` returns `NewValueRef(nil)` plus an error.** Any
  caller that dereferences the returned `ValueRef` before checking the
  error nil-panics. The error is always returned, so callers that respect
  it are safe — but the zero-value `ValueRef{}` would be a safer sentinel.

### Things that are not obvious from the code

- The exact contract between `runDefers` and `Recoverable` — defers run on
  `OpRet` and on recoverable propagation, but are skipped on
  non-recoverable propagation, and as noted above also (incorrectly) run
  on every `OpIterYield`. None of this is documented in code comments.

- `frame.pipeline` semantics are only inferable from the native builtins
  that read `IsPipeline` / `IsCallerPipeline`.

- Both `Fn.Call` and `Closure.Call` panic ("not called") because real
  dispatch is special-cased in `call()`. They exist purely to satisfy
  `Callable`; the dual role is easy to miss.

- `OpLoadCapture` vs `OpLoadCaptureRef`: the difference (auto-deref vs.
  raw `ValueRef`) is only clear from the compiler emit sites.

- `ILIterator.iterNRet` (declared yield arity) vs. the `nret` parameter
  `OpNext` passes (requested arity at the call site) are tracked
  separately and no opcode enforces consistency between them.
