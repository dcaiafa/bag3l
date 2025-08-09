package lib

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

func iterate(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	iter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	if iter.IsClosed() {
		return nil, vm.ErrIteratorClosed
	}

	return []vm.Value{iter}, nil
}
