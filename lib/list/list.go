package list

import "github.com/dcaiafa/bag3l/internal/vm"

//go:generate go run ../../internal/stub/stubgen list.stubgen

func append0(m *vm.VM, l *vm.List, v []vm.Value) (*vm.List, error) {
	for _, val := range v {
		l.Add(val)
	}
	return l, nil
}

func append_iter0(m *vm.VM, l *vm.List, iter vm.Iterator) (*vm.List, error) {
	for {
		v, err := m.IterNext(iter, 1)
		if err != nil {
			return nil, err
		}
		if v == nil {
			break
		}
		l.Add(v[0])
	}
	return l, nil
}

func find0(m *vm.VM, l *vm.List, v vm.Value) (vm.Value, error) {
	ndx := l.Find(v)
	if ndx == -1 {
		return nil, nil
	}
	return vm.NewInt(int64(ndx)), nil
}

func from_iter0(m *vm.VM, v vm.Iterator) (*vm.List, error) {
	list := vm.NewList()
	for {
		vals, err := m.IterNext(v, 1)
		if err != nil {
			return nil, err
		}
		if vals == nil {
			break
		}
		list.Add(vals[0])
	}
	return list, nil
}
