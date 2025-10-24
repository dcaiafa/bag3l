package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func takeUntil(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 2, 3); err != nil {
		return nil, err
	}

	inIter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	untilFunc, err := getCallableArg(args, 1)
	if err != nil {
		return nil, err
	}

	takeIter := &takeUntilIterator{
		inIter:    inIter,
		untilFunc: untilFunc,
	}

	return []vm.Value{vm.NewIterator(takeIter.Next, takeIter.Close, inIter.IterNRet())}, nil
}

type takeUntilIterator struct {
	inIter    vm.Iterator
	untilFunc vm.Callable
	stopNext  bool
}

func (i *takeUntilIterator) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if i.inIter.IsClosed() {
		return nil, nil
	}

	if i.stopNext {
		err := m.IterClose(i.inIter)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	v, err := m.IterNext(i.inIter, i.inIter.IterNRet())
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	res, err := m.Call(i.untilFunc, v, 1)
	if err != nil {
		return nil, err
	}

	i.stopNext = vm.CoerceToBool(res[0])

	return v, nil
}

func (i *takeUntilIterator) Close(vm *vm.VM) error {
	return vm.IterClose(i.inIter)
}
