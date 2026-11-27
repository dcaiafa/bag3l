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
| `stack_depth.go` | `computeMaxStack` — static operand-stack bound per `Fn` |
| `stack_debug_on.go` / `stack_debug_off.go` | `debugStackCheck` build-tag switch |
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
- `minArgs` — the declared parameter count. `OpInitCallFrame` uses it to pad
  any parameters the caller omitted with `nil` so every parameter owns a
  stack slot (see "Calling convention" below).
- `maxStack` — the peak operand-stack depth above the frame's locals,
  computed by `computeMaxStack` when the emitter finishes the function
  (`Emitter.PopFn`). See "Stack overflow and growth" below.

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
    framePool     []*frame            // freelist for frame structs
    pendingErr    error               // error injected by SignalError
    contextStack  []context.Context
}
```

The operand stack is unified: arguments, locals, captures, and
temporaries all live on it. No instruction takes a pointer *into* the operand
stack; assignment writes into stack slots by index (see the store opcodes
below), which is what lets the backing array be reallocated when the stack
grows (see "Stack overflow and growth"). The hot fields (`frame`, `pkg`,
`globals`, `instrs`) are cached out of `frame.fn.pkg` and `frame.fn.instrs` whenever a
frame becomes current, so the interpreter loop only touches `co.*` and not
deeper pointer chains.

### Stack overflow and growth

A coroutine's stack starts at `initialStackSize` (256) slots. An
`ILIterator`'s stack starts out `nil` and is allocated on its first call,
sized to exactly what its frame needs. Both grow on demand.

Push opcodes store to `stack[sp]` without checking for room. That is safe
because each frame reserves everything it can need up front:

- **Static bound.** `computeMaxStack` (`stack_depth.go`) walks a function's
  bytecode once at compile time, following every branch, and records the peak
  operand depth above the locals in `fn.maxStack`. Each opcode has a fixed
  pop/push effect; `OpCall` additionally counts a transient peak of `nret`
  slots because results are written before being moved down over the
  callable. Catch handlers are entered at depth 1 (the error value). Where
  paths join with different depths the larger one wins. The pass panics on
  operand underflow or an unknown opcode, both of which are compiler bugs —
  **a new opcode must be added to its switch.**
- **Prologue reservation.** `OpInitCallFrame` calls
  `co.Reserve(padding + nLocals + fn.maxStack)` before touching the stack.
  In the common case that is one comparison per call and nothing per
  instruction.
- **Dynamic pushes.** The few pushes whose size is only known at run time
  call `co.Reserve(n)` themselves: argument expansion in `OpCall`
  (`f(list...)`), `VM.Call` (native code pushing arguments), `iterNext`, and
  `Run`.

When `Reserve` finds too little room it allocates a new backing array (at
least double, capped at `maxStackSize`), copies the live slots `[0, sp)`
and replaces `co.stack`. Consequences worth knowing:

- Nothing may keep writing through a slice of the stack across a call that
  can reach `Reserve`. The interpreter loop always goes through
  `m.co.stack`. The `args` slice handed to a native function may end up
  aliasing the old array if that function calls back into the VM; it stays
  readable and holds the right values, but it is no longer the live stack.
- While an `ILIterator` runs, `co.stack` *is* the iterator's stack, so growth
  happens there. `call` and `iterNext` store `co.stack` back into the
  iterator before switching to the caller's stack.
- Stacks never shrink.
- Slots above `sp` are not cleared when values are popped, so a stack that
  was once deep keeps those values reachable by the GC until the slots are
  overwritten.

Two limits turn runaway programs into `ErrStackOverflow`:

- `maxStackSize` (2^20 slots, 16 MB) per operand stack.
- `maxCallDepth` (10000 frames) per coroutine, checked in
  `OpInitCallFrame`. The slot limit alone does not bound recursion: every
  Bag3l call also recurses on the Go stack, and each `ILIterator` runs on
  its own operand stack, so a recursive generator would otherwise never run
  out of slots.

`ErrStackOverflow` is a `NonRecoverableError`: it cannot be caught by
`try`/`catch`.

Building with `-tags bag3ldebug` turns on `debugStackCheck`, which asserts
after every instruction that `sp` is within the frame's computed bound. Run
the test suite that way after changing opcodes or `computeMaxStack`;
`make test` runs it both with and without the tag.

### Frames

```go
type frame struct {
    nRet       int
    nArg       int             // arguments the caller actually provided
    argSlots   int             // parameter slots (>= nArg; omitted params padded)
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
    bp - argSlots                  bp                                      sp
```

`OpInitCallFrame` is emitted as the first instruction of every Bag3l
function. It first **pads omitted parameters**: if the caller provided fewer
arguments than the function declares (`fn.minArgs`), it pushes `nil` for each
missing parameter so every parameter owns a slot, and bumps `argSlots` to the
declared count. It then sets `bp = sp`, advances `sp` by `nLocals`, and zeros
the local slots. Subsequent loads use `bp` as the anchor:

| Opcode | Slot |
| --- | --- |
| `OpLoadArg n` / `OpLoadArgDeref n` | `bp - argSlots + n` |
| `OpLoadLocal n` / `OpLoadLocalDeref n` | `bp + n` |
| `OpLoadGlobal n` | `globals[n]` |
| `OpLoadCapture n` / `OpLoadCaptureBox n` | `frame.caps[n]` |
| `OpLoadLiteral lit, dep` | `pkg.Deps[dep].Literals[lit]` |

The `*Deref` loads read through a lifted variable's box; `OpLoadCaptureBox`
pushes the raw `ValueRef` box (used to build a nested closure's capture list),
whereas `OpLoadCapture` auto-dereferences it.

Assignment is the mirror image — direct stores by index, each popping the value
on top of the stack:

| Opcode | Effect |
| --- | --- |
| `OpStoreLocal n` / `OpStoreLocalDeref n` | `bp + n` (direct / through box) |
| `OpStoreArg n` / `OpStoreArgDeref n` | `bp - argSlots + n` (direct / through box) |
| `OpStoreGlobal n` | `globals[n]` |
| `OpStoreCapture n` | `*frame.caps[n].Ref` |
| `OpStoreIndex` | `container.SetIndex(key, value)`; pops `value, container, key` |

Because omitted parameters are padded to `nil`, a parameter is just a local
that happens to be pre-initialized from the caller's argument: it can always
be read or assigned, whether or not the caller provided it. `nArg` (the count
the caller actually passed) is kept distinct from `argSlots` (the padded width
used for addressing) so that `narg()` / `args()` still report the real
argument count. When the caller passes *more* arguments than the function
declares, `argSlots == nArg` and the extra arguments are reachable only via
`args()`.

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

`ValueRef` is `struct{ Ref *Value }`. It is a heap-allocated box used to give a
variable a stable storage location independent of the operand stack. A function
argument or local is "lifted" — boxed in a `ValueRef` stored in its slot — when
an inner closure or iterator captures it (so the closure and the outer frame
share one cell). The box is created by `OpInitLiftedLocal` (locals) or
`OpLiftArg` (args); lifted slots are read with `OpLoad*Deref`, written with
`OpStore*Deref`, and shared into a closure's `caps` via `OpLoadCaptureBox`.
Non-lifted variables live directly in their stack slot and are read/written
without boxing.

Assignment statements never push a pointer into the stack. The right-hand side
is fully evaluated first (leaving `v1..vN` on the stack); the lvalues are then
stored **right-to-left** so that each value is exposed on top of the stack just
as its target's container/key prefix is pushed above it, letting `OpStoreIndex`
find the value directly beneath the `container, key` pair. As a consequence, an
lvalue's container/index subexpressions are evaluated *after* the right-hand
side, and in reverse order.

The container/collection types:

- `*List` — backed by `[]Value`. Negative indices are supported by `Index`
  and `Slice`. `SetIndex` requires a non-negative in-range `Int` index.
- `*Map` — `map[Value]*mapNode` plus a doubly linked list of `mapNode`s, so
  iteration is in insertion order. `Put`/`SetIndex` insert at the tail
  if the key is new. `mapIter` walks the list using `GetNext`, which silently
  ends iteration if the current key was deleted mid-loop.
- `String` — immutable value type wrapping `string`. `SetIndex` always
  returns an error ("cannot modify str").

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

- `OpCall`: `op1` packs `(narg | expandFlag)`; `op2` is `nret`.
  Flags: `CallExpandFlag = 0x80000000`, `CallArgCountMask = 0x7FFFFFFF`.
- `OpNewIter`: `op1` packs `fnIndex (24 bits low) | iterNRet (8 bits high)`;
  `op2` is the capture count.
- `OpNewClosure`: `op1` is the function literal index; `op2` is the capture
  count.
- `OpLoadIndex`: `op2 & OptionalIndexFlag` (`0x0001`) marks an optional
  index (`?.` / `?[]`) so a missing key yields nil instead of an error.
- `OpLoadLiteral`: `op1` is the literal index; `op2` is the dependency index
  into `pkg.Deps` (allowing cross-package access to constants/functions).

## Calling convention

`call(callable, narg, nret)` handles all call shapes:

| Callable type | Behavior |
| --- | --- |
| `*Closure` | New frame with `fn = callable.fn`, `caps = callable.caps`; `runFrame` |
| `*Fn` | New frame with `fn = callable` |
| `*NativeIterator` | Treated as an external call with no captures |
| `*ILIterator` | Resumes the generator (see Iterators) |
| `Callable` | Dispatched through `callExtFn` |
| `nil` | `ErrCannotCallNil` |
| other | error: `cannot call <type>` |

`OpCall` sets up the call from the caller's stack:

1. Decode `narg`/`nret`/`expand`.
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
owns its own private operand stack (allocated lazily, see "Stack overflow and
growth") plus the resume state. Everything the iterator's body calls also
runs on that stack, not on the stack of whoever is consuming the iterator —
which is why iterator stacks have to be able to grow.

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
    bp         int                 // base pointer, carried across yields
    closed     bool
    running    bool                // body is executing; see below
}
```

On every `Next()`, the VM:

1. Saves the host coroutine's `stack`/`sp` (`rstack`, `rsp`).
2. Swaps in the iterator's `stack`/`sp` and pushes a frame whose `ip`,
   `bp`, `tryCatches`, `defers`, `nlocals`, and `caps` come from the iterator.
3. Runs `runFrame` until `OpIterYield` (yield) or `OpIterRet`/`OpRet` (done).
4. Copies the yielded values back to the host stack and restores the host
   `stack`/`sp`.

`OpIterYield` saves `ip`, `bp`, `tryCatches`, `defers`, `nlocals` back onto
the `ILIterator` and returns `nil` from the inner interpreter loop. `OpIterRet`
sets `ip = -1` to mark the generator finished.

Both ways of advancing an `ILIterator` — `iterNext` (used by `next`, `for`
and native code) and the `*ILIterator` branch of `call` (a direct `it()`) —
go through `resumeILIterator`, which does the steps above and reports
whether the body yielded. They differ only in how they report the end:
`iterNext` returns `false`, while `call` pushes `False` plus nils for the
remaining slots, whether the iterator was already finished or finished
during that call. Arguments passed to a direct call are ignored.

The saved `bp` is in practice always 0: the compiler wraps every yielding
function so that the iterator body itself takes no parameters (see
`Func.RunPass` in `internal/ast/func_base.go`), and it starts on an empty
stack.

An `ILIterator` is not re-entrant. Its saved `ip`/`sp`/`bp` are only
written at a yield, so they are stale for as long as the body runs, and a
second activation started from them would run on the same stack slots as the
first. `resumeILIterator` therefore sets `running` around `runFrame`, and
both it and `Close` fail with `ErrIteratorRunning` (not catchable) when the
flag is set. This is reachable two ways: the body resumes or closes its own
iterator, or the body blocks (`sleep`, I/O, `co.start`) and another coroutine
advances the same iterator. Sharing one generator between coroutines is
therefore an error unless the consumers never overlap.

`Close` marks the iterator closed and runs any pending `defer`s registered
in the iterator body in LIFO order. It is idempotent — calling `Close` again
is a no-op. Defers are cleared once they run, and a defer that errors aborts
the LIFO sweep (remaining defers do not run).

### `*NativeIterator`

Wraps a Go-side `*NativeFn` plus an optional `CloseFn`. Each call to `Next`
invokes the native function; returning an empty slice (`nil, nil`) signals
end-of-iteration and triggers `Close`.

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

**The rule: runtime errors are never catchable.** That covers stack
overflow, dereferencing or indexing nil, divide by zero, and every other
fault the runtime detects. How it is enforced today depends on where the
error comes from:

- An error returned by an instruction in the interpreter loop is wrapped
  with `Recoverable: false` by `resume`. New VM-level errors get the right
  behavior by just returning a plain `error`.
- An error returned by a **native function** is wrapped with
  `Recoverable: true` by `callExtFn` (and by `iterNext` for native
  iterators), so it *is* catchable unless the function wraps it with
  `MakeNonRecoverableError`. An error the VM hands to native code (e.g.
  `ErrStackOverflow` from `VM.Call`) must therefore be a
  `NonRecoverableError` to stay uncatchable on its way back out.

> **TODO (audit pending).** The rule also applies to native functions: a
> builtin that reports a runtime fault must not be catchable. This is a new
> rule and `lib/` has not been audited against it, so today most such errors
> are catchable by default. Which native errors count as runtime faults, and
> how to enforce it (per call site with `MakeNonRecoverableError`, or by
> flipping the default in `callExtFn`), is still to be decided.

**`try` assumes an empty operand stack.** Entering the catch block resets
`sp` to `bp + nLocals + 1`, discarding any temporaries the frame had pushed
when the error occurred. That is correct only because `try` is a statement
and statements start with nothing above the locals; nothing enforces it.
`computeMaxStack` makes the same assumption (handlers start at depth 1).

Defers run on the panic path as well: when a non-recoverable runtime error
escapes a frame, the frame's `defers` are invoked LIFO before propagation
continues. This preserves the "bugs can't be caught by `try`/`catch`"
stance — only the cleanup side runs, not user error-handling logic. If a
deferred closure itself errors, the new error is wrapped onto the
propagating one (`defer threw error: … while handling error: …`), though
this wrap is currently only visible at the immediately-enclosing frame —
any further wrap-cycle collapses back to the original `*RuntimeError`.

## Defer

`OpDefer` pops a `*Closure` and appends it to `frame.defers`. **It type-asserts
to `*Closure` without an `ok` check**; deferring a `*Fn` or `*NativeFn`
panics.

On normal frame return (`OpRet`), `resume` invokes `runDefers`, which calls
each deferred closure in LIFO order. On any error that escapes the frame —
recoverable or non-recoverable — `runDefers` also runs, and a defer-side
error is wrapped onto the propagating error message. On the catch path
(recoverable error with an active `tryCatch`), defers do *not* run; they
remain on the frame and run when it later returns or unwinds.

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
`Closer` that registered itself via `RegisterCloser`. `RegisterCloser`,
`UnregisterCloser`, and `shutdown`'s iteration over the closer set are all
guarded by `VM.mu`; `shutdown` snapshots the set under the lock and calls
`Close()` outside it, so a `Close` that re-enters `UnregisterCloser` does not
deadlock.

## Stack traces

`GetStackInfo()` snapshots the current `callStack` into `[]FrameInfo`
records. For Bag3l frames, the source file/line is reconstructed by
`getLocation(fn, ip)`, which scans `fn.locations` to find the largest entry
with `ip <= currentIp`. `wrapRuntimeError` captures a stack trace the first
time a runtime error is wrapped and reuses it on every re-wrap so the trace
reflects the throw site, not propagation sites.

`RuntimeError.Error()` prints at most 20 frames; for deeper stacks (runaway
recursion) it keeps the innermost and outermost 10 and reports how many were
omitted.

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

### Sub-optimal / fragile patterns

- **`getLocation` does a linear scan** of `fn.locations` on every
  stack-trace synthesis. Locations are stored sorted by `ip` and should
  use binary search.

- **`wrapRuntimeError` allocates a stack trace on every error**, including
  errors that will be caught by a nearby `try`. `GetStackInfo` allocates a
  `[]FrameInfo` proportional to `callStack` depth.

- **`coroutine.framePool` is per-coroutine.** Each spawned coroutine
  allocates its own pool, so coroutine-heavy workloads re-pay frame
  allocation costs per fiber.

- **`Bool.EvalOp` does not return `ErrOperationNotSupported`.** It returns
  a plain `fmt.Errorf("bool does not support this operation")`, so the
  `FallbackEvaluator` path inside `EvalOp` can never trigger for bool
  operands.

- **`RuntimeError.Index` errors on unknown fields** instead of returning
  `(nil, false, nil)`. That means `err?.unknown` (optional access) still
  throws, defeating the point of the optional-index flag for this type.

### Unconfirmed

Noticed while reading the code; none of these has been reproduced or ruled
out. Verify before relying on or fixing them.

- **`GetCallerArgs` indexes the current operand stack with another frame's
  `bp`.** Frames on `callStack` can live on different operand stacks (each
  running `ILIterator` switches `co.stack`), so the result is only right when
  the caller's frame is on the current one.

### Things that are not obvious from the code

- `OpIterYield`'s `op1` (the number of yielded values) is ignored by the
  interpreter, which derives the count from `sp`. It is read only by
  `computeMaxStack`, so emitting a wrong count would not fail at run time —
  it would silently skew the function's stack bound.

- The exact contract between `runDefers` and the frame error path: defers
  run on `OpRet`, on propagation regardless of recoverability, and at
  iterator end-of-life via `Close`. They do *not* run when an iterator
  yields (`OpIterYield` migrates the defers list from frame to iterator
  instead), nor on the catch path itself (defers stay on the frame and run
  when it later returns or unwinds). None of this is documented in code
  comments.

- Both `Fn.Call` and `Closure.Call` panic ("not called") because real
  dispatch is special-cased in `call()`. They exist purely to satisfy
  `Callable`; the dual role is easy to miss.

- `OpLoadCapture` vs `OpLoadCaptureBox`: the difference (auto-deref vs.
  raw `ValueRef` box) is only clear from the compiler emit sites —
  `OpLoadCaptureBox` is used solely to build a nested closure's capture list.

- `ILIterator.iterNRet` (declared yield arity) vs. the `nret` requested by the
  caller (passed into `iterNext` at the call site) are tracked separately and
  no opcode enforces consistency between them.
