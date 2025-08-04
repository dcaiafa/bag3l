package lib

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

func enumerate(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	inner, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	iter := &enumIter{
		iter: inner,
	}

	res := vm.NewIterator(iter.Next, iter.Close, 1)

	return []vm.Value{res}, nil
}

type enumIter struct {
	iter  vm.Iterator
	index int
}

func (i *enumIter) Next(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	res, err := m.IterNext(i.iter, nret-1)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	idx := i.index
	i.index++

	res = append([]vm.Value{vm.NewInt(int64(idx))}, res...)

	return res, nil
}

func (i *enumIter) Close(vm *vm.VM) error {
	return i.iter.Close(vm)
}
