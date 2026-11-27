package vm

import (
	"errors"
	"fmt"
	"strings"
)

type NonRecoverableError struct {
	wrapped error
}

func (e *NonRecoverableError) Error() string {
	return e.wrapped.Error()
}

func (e *NonRecoverableError) Unwrap() error {
	return e.wrapped
}

func MakeNonRecoverableError(err error) error {
	return &NonRecoverableError{wrapped: err}
}

type RuntimeError struct {
	Err         error
	ErrValue    Value
	Stack       []FrameInfo
	Recoverable bool
}

var _ error = (*RuntimeError)(nil)

func wrapRuntimeError(vm *VM, err *error, recoverable bool) *RuntimeError {
	if *err == nil {
		return nil
	}
	var rerr *RuntimeError
	if !errors.As(*err, &rerr) {
		rerr = &RuntimeError{
			Err:         *err,
			Recoverable: recoverable,
		}
		// A base NonRecoverableError error overrides `recoverable`.
		var nre *NonRecoverableError
		if errors.As(*err, &nre) {
			rerr.Recoverable = false
		}
	}
	if rerr.Stack == nil {
		rerr.Stack = vm.GetStackInfo()
	}
	*err = rerr
	return rerr
}

func (e *RuntimeError) Index(k Value) (Value, bool, error) {
	k, ok := k.(String)
	if !ok {
		return nil, false, fmt.Errorf("RuntimeError does not have field %v", k)
	}

	switch k.String() {
	case "error":
		errValue := e.ErrValue
		if errValue == nil {
			errValue = NewString(e.Err.Error())
		}
		return errValue, true, nil

	default:
		return nil, false, fmt.Errorf("RuntimeError does not have field %v", k)
	}
}

func (e *RuntimeError) SetIndex(k, value Value) error {
	return fmt.Errorf("RuntimeError is read-only")
}

func (e *RuntimeError) String() string { return e.Error() }
func (e *RuntimeError) Type() string   { return "error" }
func (e *RuntimeError) Traits() Traits { return TraitNone }

func (e *RuntimeError) Message() string {
	if e.Err != nil {
		return e.Err.Error()
	} else if e.ErrValue != nil {
		return e.ErrValue.String()
	} else {
		return "<no message>"
	}
}

func (e *RuntimeError) Error() string {
	str := strings.Builder{}
	str.WriteString(e.Message())

	// Elide the middle of very deep stacks (e.g. runaway recursion).
	const maxShown = 20
	for i, f := range e.Stack {
		if len(e.Stack) > maxShown && i == maxShown/2 {
			fmt.Fprintf(&str, "\n ... %v frames omitted ...", len(e.Stack)-maxShown)
		}
		if len(e.Stack) > maxShown && i >= maxShown/2 && i < len(e.Stack)-maxShown/2 {
			continue
		}
		fmt.Fprintf(&str, "\n %v", f.String())
	}

	return str.String()
}

func (e *RuntimeError) Unwrap() error {
	return e.Err
}

var ErrCannotCallNil = errors.New("cannot evaluate function call because target is nil")
var ErrIsNotIterable = errors.New("is not iterable")
var ErrIsNotReadable = errors.New("is not readable")
var ErrDivideByZero = errors.New("divide by zero")

// ErrStackOverflow is returned when a coroutine exhausts its operand stack or
// exceeds the maximum call depth. It cannot be caught by try/catch.
var ErrStackOverflow = MakeNonRecoverableError(errors.New("stack overflow"))
