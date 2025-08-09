package lib

import (
	"github.com/dcaiafa/bag3l/internal/stub"
	"github.com/dcaiafa/bag3l/internal/vm"
)

func close_iter(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 2 {
		return nil, errInvalidNumberOfArgs
	}

	iter, ok := args[0].(vm.Iterator)
	if !ok {
		return nil, stub.InvalidArg(args, 0)
	}
	throwOnError, err := getBoolArg(args, 1)
	if err != nil {
		return nil, err
	}

	err = m.IterClose(iter)
	if err != nil && throwOnError {
		return nil, err
	}

	return nil, nil
}
