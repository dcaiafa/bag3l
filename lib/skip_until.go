package lib

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

func skipUntil(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
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

	skipIter := &skipUntilIterator{
		inIter:    inIter,
		untilFunc: untilFunc,
	}

	return []vm.Value{vm.NewIterator(skipIter.Next, skipIter.Close, inIter.IterNRet())}, nil
}

type skipUntilIterator struct {
	inIter    vm.Iterator
	untilFunc vm.Callable
	open      bool
}

func (i *skipUntilIterator) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	for {
		v, err := m.IterNext(i.inIter, i.inIter.IterNRet())
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, nil
		}
		if !i.open {
			res, err := m.Call(i.untilFunc, v, 1)
			if err != nil {
				return nil, err
			}
			i.open = vm.CoerceToBool(res[0])
			continue
		}
		return v, nil
	}
}

func (i *skipUntilIterator) Close(vm *vm.VM) error {
	return vm.IterClose(i.inIter)
}
