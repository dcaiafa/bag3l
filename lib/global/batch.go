package global

import (
	nitro "github.com/dcaiafa/bag3l"
	"github.com/dcaiafa/bag3l/internal/vm"
)

func batch0(m *vm.VM, inIter vm.Iterator, n int64) (vm.Iterator, error) {
	batchIter := &batchIter{
		inIter: inIter,
		n:      int(n),
	}

	outIter := nitro.NewIterator(
		batchIter.Next, batchIter.Close, 1)

	return outIter, nil
}

type batchIter struct {
	inIter nitro.Iterator
	n      int
}

func (i *batchIter) Next(vm *nitro.VM, args []nitro.Value, nRet int) ([]nitro.Value, error) {
	if i.inIter.IsClosed() {
		return nil, nil
	}

	var b = make([]nitro.Value, 0, i.n)
	for len(b) < i.n {
		v, err := vm.IterNext(i.inIter, 1)
		if err != nil {
			return nil, err
		}
		if v == nil {
			break
		}
		b = append(b, v[0])
	}

	if len(b) == 0 {
		return nil, nil
	}

	return []nitro.Value{nitro.NewArrayFromSlice(b)}, nil
}

func (i *batchIter) Close(vm *nitro.VM) error {
	vm.IterClose(i.inIter)
	return nil
}
