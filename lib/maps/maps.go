package maps

import (
	"fmt"

	"github.com/dcaiafa/bag3l/internal/vm"
)

//go:generate go run ../../internal/stub/stubgen maps.stubgen

func clone0(m *vm.VM, mp *vm.Map) (*vm.Map, error) {
	return mp.Clone(), nil
}

func update0(m *vm.VM, mp *vm.Map, other *vm.Map) (*vm.Map, error) {
	other.ForEach(func(k, v vm.Value) bool {
		mp.Put(k, v)
		return true
	})
	return mp, nil
}

func update1(m *vm.VM, mp *vm.Map, f vm.Callable) (*vm.Map, error) {
	res, err := m.Call(f, []vm.Value{mp}, 1)
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
	return update0(m, mp, other)
}

func union0(m *vm.VM, mp *vm.Map, other *vm.Map) (*vm.Map, error) {
	return update0(m, mp.Clone(), other)
}

func union1(m *vm.VM, mp *vm.Map, f vm.Callable) (*vm.Map, error) {
	return update1(m, mp.Clone(), f)
}

func delete0(m *vm.VM, mp *vm.Map, k vm.Value) (*vm.Map, error) {
	mp.Delete(k)
	return mp, nil
}

func make0(m *vm.VM, iter vm.Iterator, f vm.Callable) (*vm.Map, error) {
	mp := vm.NewMap()

	if f != nil {
		for {
			v, err := m.IterNext(iter, iter.IterNRet())
			if err != nil {
				return nil, err
			}
			if v == nil {
				break
			}
			res, err := m.Call(f, v, 1)
			if err != nil {
				return nil, err
			}

			larg, ok := res[0].(*vm.List)
			if !ok {
				return nil, fmt.Errorf(
					"conversion func must return \"List\"; instead it returned %v",
					vm.TypeName(res[0]))
			}
			if larg.Len() != 2 {
				return nil, fmt.Errorf(
					"conversion func must return a list with 2 elements (key and value); "+
						"instead it returned a list with %v elements",
					larg.Len())
			}
			mp.Put(larg.Get(0), larg.Get(1))
		}
	} else {
		for {
			v, err := m.IterNext(iter, 1)
			if err != nil {
				return nil, err
			}
			if v == nil {
				break
			}
			mp.Put(v[0], vm.True)
		}
	}
	return mp, nil
}
