package maps

import (
	"fmt"

	nitro "github.com/dcaiafa/bag3l"
	"github.com/dcaiafa/bag3l/internal/vm"
)

//go:generate stubgen maps.stubgen

func clone0(vm *vm.VM, m *vm.Map) (*vm.Map, error) {
	return m.Clone(), nil
}

func update0(vm *vm.VM, m *vm.Map, other *vm.Map) (*vm.Map, error) {
	other.ForEach(func(k, v nitro.Value) bool {
		m.Put(k, v)
		return true
	})
	return m, nil
}

func update1(theVM *vm.VM, m *vm.Map, f vm.Callable) (*vm.Map, error) {
	res, err := theVM.Call(f, []vm.Value{m}, 1)
	if err != nil {
		return nil, err
	}
	var ok bool
	other, ok := res[0].(*vm.Map)
	if !ok {
		return nil, fmt.Errorf(
			"func expected to return \"Map\", but returned %q instead",
			vm.TypeName(res[0]))
	}
	return update0(theVM, m, other)
}

func delete0(vm *vm.VM, m *vm.Map, k vm.Value) (*vm.Map, error) {
	m.Delete(k)
	return m, nil
}

func make0(vm *nitro.VM, iter vm.Iterator, f vm.Callable) (*vm.Map, error) {
	m := nitro.NewObject()

	if f != nil {
		for {
			v, err := vm.IterNext(iter, iter.IterNRet())
			if err != nil {
				return nil, err
			}
			if v == nil {
				break
			}
			res, err := vm.Call(f, v, 1)
			if err != nil {
				return nil, err
			}

			larg, ok := res[0].(*nitro.Array)
			if !ok {
				return nil, fmt.Errorf(
					"conversion func must return \"List\"; instead it returned %v",
					nitro.TypeName(res[0]))
			}
			if larg.Len() != 2 {
				return nil, fmt.Errorf(
					"conversion func must return a list with 2 elements (key and value); "+
						"instead it returned a list with %v elements",
					larg.Len())
			}
			m.Put(larg.Get(0), larg.Get(1))
		}
	} else {
		for {
			v, err := vm.IterNext(iter, 1)
			if err != nil {
				return nil, err
			}
			if v == nil {
				break
			}
			m.Put(v[0], nitro.True)
		}
	}
	return m, nil
}
