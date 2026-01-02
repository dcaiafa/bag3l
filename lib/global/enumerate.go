package global

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

func enumerate0(m *vm.VM, inner vm.Iterator) (vm.Iterator, error) {
	iter := &enumIter{
		iter: inner,
	}

	res := vm.NewIterator(iter.Next, iter.Close, iter.iter.IterNRet()+1)

	return res, nil
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
