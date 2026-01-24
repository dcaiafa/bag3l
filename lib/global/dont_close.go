package global

import "github.com/dcaiafa/bag3l/internal/vm"

func dont_close0(m *vm.VM, iter vm.Iterator) (vm.Iterator, error) {
	res := vm.NewIterator(
		func(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
			return m.IterNext(iter, nret)
		},
		nil, // Don't close.
		iter.IterNRet(),
	)

	return res, nil
}
