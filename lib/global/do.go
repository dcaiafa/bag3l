package global

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

func do0(m *vm.VM, iter vm.Iterator, fn vm.Callable) error {
	for {
		v, err := m.IterNext(iter, iter.IterNRet())
		if err != nil {
			return err
		}
		if v == nil {
			return nil
		}
		_, err = m.Call(fn, v, 0)
		if err != nil {
			return err
		}
	}
}
