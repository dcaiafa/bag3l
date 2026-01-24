package global

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

func filter0(m *vm.VM, iter vm.Iterator, test vm.Callable) (vm.Iterator, error) {
	filterIter := &filterIter{
		inIter: iter,
		test:   test,
	}

	outIter := vm.NewIterator(filterIter.Next, filterIter.Close, iter.IterNRet())

	return outIter, nil
}

type filterIter struct {
	inIter vm.Iterator
	test   vm.Callable
}

func (i *filterIter) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	for {
		v, err := m.IterNext(i.inIter, i.inIter.IterNRet())
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, nil
		}
		res, err := m.Call(i.test, v, 1)
		if err != nil {
			m.IterClose(i.inIter)
			return nil, err
		}
		if vm.CoerceToBool(res[0]) {
			return v, nil
		}
	}
}

func (i *filterIter) Close(m *vm.VM) error {
	m.IterClose(i.inIter)
	return nil
}
