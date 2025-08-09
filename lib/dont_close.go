package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func dontClose(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	iter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	res := vm.NewIterator(
		func(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
			return m.IterNext(iter, nret)
		},
		nil, // Don't close
		iter.IterNRet(),
	)

	return []vm.Value{res}, nil
}
