package global

import (
	"reflect"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func deep_equal0(m *vm.VM, a, b vm.Value) (bool, error) {
	ctx := &deepEqualContext{
		visiting: make(map[vm.Value]bool),
	}

	equal := ctx.Equal(a, b)
	return equal, nil
}

type deepEqualContext struct {
	visiting map[vm.Value]bool
}

func (e *deepEqualContext) Equal(a, b vm.Value) bool {
	if a == b {
		return true
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	switch a := a.(type) {
	case *vm.List:
		return e.arrayEqual(a, b.(*vm.List))

	case *vm.Map:
		return e.mapEqual(a, b.(*vm.Map))

	default:
		res, err := vm.EvalOp(vm.OpEq, a, b)
		return err == nil && res == vm.True
	}
}

func (e *deepEqualContext) arrayEqual(a, b *vm.List) bool {
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

func (e *deepEqualContext) mapEqual(a, b *vm.Map) bool {
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
	a.ForEach(func(k, va vm.Value) bool {
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
