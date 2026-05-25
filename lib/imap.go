package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func imap(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 2, 2); err != nil {
		return nil, err
	}
	inIter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}
	fn, err := getCallableArg(args, 1)
	if err != nil {
		return nil, err
	}

	mapIter := &mapIter{
		inIter: inIter,
		fn:     fn,
	}

	outIter := vm.NewIterator(mapIter.Next, mapIter.Close, 1)
	return []vm.Value{outIter}, nil
}

type mapIter struct {
	inIter vm.Iterator
	fn     vm.Value
}

func (i *mapIter) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	v, err := m.IterNext(i.inIter, i.inIter.IterNRet())
	if err != nil {
		i.Close(m)
		return nil, err
	}
	if v == nil {
		i.Close(m)
		return nil, nil
	}
	res, err := m.Call(i.fn, v, 1)
	if err != nil {
		i.Close(m)
		return nil, err
	}

	return []vm.Value{res[0]}, nil
}

func (i *mapIter) Close(m *vm.VM) error {
	m.IterClose(i.inIter)
	return nil
}
