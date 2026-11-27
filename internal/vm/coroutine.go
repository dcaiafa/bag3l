package vm

import (
	"context"

	"github.com/dcaiafa/bag3l/internal/fiber"
)

type coroutine struct {
	fiber        *fiber.Fiber
	callStack    []*frame
	frame        *frame
	pkg          *CompiledPackage
	stack        []Value
	globals      []Value
	instrs       []Instr
	sp           int
	ip           int
	framePool    []*frame
	pendingErr   error
	contextStack []context.Context
}

// Reserve makes sure the current operand stack has room for n more values
// above sp, growing it if necessary. It returns ErrStackOverflow if that
// would take the stack past maxStackSize.
//
// Bytecode pushes are covered by the single Reserve in OpInitCallFrame; other
// callers are pushes whose size is only known at run time.
//
// Growing replaces co.stack, so a slice of the stack taken before a call that
// can reach Reserve must not be written through afterwards.
func (co *coroutine) Reserve(n int) error {
	if co.sp+n <= len(co.stack) {
		return nil
	}
	return co.growStack(co.sp + n)
}

func (co *coroutine) growStack(need int) error {
	if need > maxStackSize {
		return ErrStackOverflow
	}
	// An iterator's stack starts out nil, so its first frame gets exactly what
	// it needs. After that, double.
	size := min(max(need, 2*len(co.stack)), maxStackSize)
	stack := make([]Value, size)
	copy(stack, co.stack[:co.sp])
	co.stack = stack
	return nil
}

func (co *coroutine) TopContext() context.Context {
	return co.contextStack[len(co.contextStack)-1]
}

func (co *coroutine) PushFrame(frame *frame) {
	if len(co.callStack) > 0 {
		co.frame.ip = co.ip
	}
	co.callStack = append(co.callStack, frame)
	co.frame = frame
	if co.frame.fn != nil {
		co.pkg = co.frame.fn.pkg
		co.globals = co.pkg.globals
		co.instrs = co.frame.fn.instrs
		co.ip = co.frame.ip
	} else {
		co.instrs = nil
	}
}

func (co *coroutine) PopFrame() {
	f := co.callStack[len(co.callStack)-1]
	*f = frame{}
	co.framePool = append(co.framePool, f)

	co.callStack[len(co.callStack)-1] = nil
	co.callStack = co.callStack[:len(co.callStack)-1]

	if len(co.callStack) > 0 {
		co.frame = co.callStack[len(co.callStack)-1]
		if co.frame.fn != nil {
			co.pkg = co.frame.fn.pkg
			co.globals = co.pkg.globals
			co.instrs = co.frame.fn.instrs
			co.ip = co.frame.ip
		}
	}
}

func (co *coroutine) NewFrame() *frame {
	if len(co.framePool) == 0 {
		return &frame{}
	}
	f := co.framePool[len(co.framePool)-1]
	co.framePool = co.framePool[:len(co.framePool)-1]
	return f
}

func (co *coroutine) PushContext(ctx context.Context) {
	co.contextStack = append(co.contextStack, ctx)
}

func (co *coroutine) PopContext() context.Context {
	n := len(co.contextStack) - 1
	popped := co.contextStack[n]
	co.contextStack[n] = nil
	co.contextStack = co.contextStack[:n]
	return popped
}
