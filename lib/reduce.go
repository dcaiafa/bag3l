package lib

import (
	"errors"

	"github.com/dcaiafa/bag3l/internal/vm"
)

var errReduceUsage = errors.New(
	`invalid usage. Expected reduce(iter, func)`)

func reduce(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 2 {
		return nil, errReduceUsage
	}

	iter, err := vm.MakeIterator(m, args[0])
	if err != nil {
		return nil, errReduceUsage
	}
	defer m.IterClose(iter)

	reducer, ok := args[1].(vm.Callable)
	if !ok {
		return nil, errReduceUsage
	}

	var accum vm.Value
	for {
		val, err := m.IterNext(iter, 1)
		if err != nil {
			return nil, err
		}
		if val == nil {
			break
		}
		res, err := m.Call(reducer, []vm.Value{accum, val[0]}, 1)
		if err != nil {
			return nil, err
		}
		accum = res[0]
	}

	res, err := m.Call(reducer, []vm.Value{accum}, 1)
	if err != nil {
		return nil, err
	}

	return []vm.Value{res[0]}, nil
}
