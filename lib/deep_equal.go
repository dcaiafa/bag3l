package lib

import (
	"reflect"

	"github.com/dcaiafa/bag3l"
)

func deepEqual(vm *bag3l.VM, args []bag3l.Value, nret int) ([]bag3l.Value, error) {
	if len(args) < 2 {
		return nil, errNotEnoughArgs
	}

	ctx := &deepEqualContext{
		visiting: make(map[bag3l.Value]bool),
	}

	equal := ctx.Equal(args[0], args[1])

	return []bag3l.Value{bag3l.NewBool(equal)}, nil
}

type deepEqualContext struct {
	visiting map[bag3l.Value]bool
}

func (e *deepEqualContext) Equal(a, b bag3l.Value) bool {
	if a == b {
		return true
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	switch a := a.(type) {
	case *bag3l.Array:
		return e.arrayEqual(a, b.(*bag3l.Array))

	case *bag3l.Object:
		return e.mapEqual(a, b.(*bag3l.Object))

	default:
		res, err := bag3l.EvalOp(bag3l.OpEq, a, b)
		return err == nil && res == bag3l.True
	}
}

func (e *deepEqualContext) arrayEqual(a, b *bag3l.Array) bool {
	if a.Len() != b.Len() {
		return false
	}

	for i := 0; i < a.Len(); i++ {
		va := a.Get(i)
		vb := b.Get(i)

		if !e.Equal(va, vb) {
			return false
		}
	}

	return true
}

func (e *deepEqualContext) mapEqual(a, b *bag3l.Object) bool {
	if a.Len() != b.Len() {
		return false
	}

	if e.visiting[a] {
		// Cycle detected.
		return false
	}

	e.visiting[a] = true
	defer delete(e.visiting, a)

	isEqual := true
	a.ForEach(func(k, va bag3l.Value) bool {
		vb, ok := b.Get(k)
		if !ok {
			isEqual = false
			return false
		}
		if !e.Equal(va, vb) {
			isEqual = false
			return false
		}
		return true
	})

	return isEqual
}
