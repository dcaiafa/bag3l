package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func takeWhile(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 2, 3); err != nil {
		return nil, err
	}

	inIter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	whileFunc, err := getCallableArg(args, 1)
	if err != nil {
		return nil, err
	}

	inclusive := false
	if len(args) == 3 {
		inclusive, err = getBoolArg(args, 2)
		if err != nil {
			return nil, err
		}
	}

	takeIter := &takeWhileIterator{
		inIter:      inIter,
		includeLast: inclusive,
		whileFunc:   whileFunc,
	}

	return []vm.Value{vm.NewIterator(takeIter.Next, takeIter.Close, inIter.IterNRet())}, nil
}

type takeWhileIterator struct {
	inIter      vm.Iterator
	includeLast bool
	whileFunc   vm.Callable
}

func (i *takeWhileIterator) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if i.inIter.IsClosed() {
		return nil, nil
	}

	v, err := m.IterNext(i.inIter, i.inIter.IterNRet())
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	res, err := m.Call(i.whileFunc, v, 1)
	if err != nil {
		return nil, err
	}

	stop := !vm.CoerceToBool(res[0])
	if stop {
		// We will no longer take anything from this iterator.
		// The iterator contract requires that we close it.
		err := m.IterClose(i.inIter)
		if err != nil {
			return nil, err
		}
		if i.includeLast {
			return v, nil
		}

		return nil, nil
	}

	return v, nil
}

func (i *takeWhileIterator) Close(vm *vm.VM) error {
	return vm.IterClose(i.inIter)
}
