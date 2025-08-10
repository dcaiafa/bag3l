package lib

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

func skipWhile(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
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

	exclusive := false
	if len(args) == 3 {
		exclusive, err = getBoolArg(args, 2)
		if err != nil {
			return nil, err
		}
	}

	skipIter := &skipWhileIterator{
		inIter:    inIter,
		skipOne:   exclusive,
		whichFunc: whileFunc,
	}

	return []vm.Value{vm.NewIterator(skipIter.Next, skipIter.Close, inIter.IterNRet())}, nil
}

type skipWhileIterator struct {
	inIter    vm.Iterator
	skipOne   bool
	whichFunc vm.Callable
	open      bool
}

func (i *skipWhileIterator) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	for {
		v, err := m.IterNext(i.inIter, i.inIter.IterNRet())
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, nil
		}
		if !i.open {
			res, err := m.Call(i.whichFunc, v, 1)
			if err != nil {
				return nil, err
			}
			i.open = !vm.CoerceToBool(res[0])
		}
		if !i.open {
			continue
		}
		if i.skipOne {
			i.skipOne = false
			continue
		}

		return v, nil
	}
}

func (i *skipWhileIterator) Close(vm *vm.VM) error {
	return vm.IterClose(i.inIter)
}
