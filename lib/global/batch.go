package global

import (
	"github.com/dcaiafa/bag3l/internal/vm"
)

func batch0(m *vm.VM, inIter vm.Iterator, n int64) (vm.Iterator, error) {
	batchIter := &batchIter{
		inIter: inIter,
		n:      int(n),
	}

	outIter := vm.NewIterator(
		batchIter.Next, batchIter.Close, 1)

	return outIter, nil
}

type batchIter struct {
	inIter vm.Iterator
	n      int
}

func (i *batchIter) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if i.inIter.IsClosed() {
		return nil, nil
	}

	var b = make([]vm.Value, 0, i.n)
	for len(b) < i.n {
		v, err := m.IterNext(i.inIter, 1)
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

	return []vm.Value{vm.NewListWithSlice(b)}, nil
}

func (i *batchIter) Close(m *vm.VM) error {
	m.IterClose(i.inIter)
	return nil
}
