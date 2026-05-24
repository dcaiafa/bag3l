package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func first(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) > 1 {
		return nil, errTooManyArgs
	}

	iter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	defer m.IterClose(iter)

	v, err := m.IterNext(iter, iter.IterNRet())
	if err != nil {
		return nil, err
	}

	if v == nil {
		return make([]vm.Value, nRet), nil
	}

	return v, nil
}
