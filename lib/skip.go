package lib

import (
	"io"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func skip(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 2, 2); err != nil {
		return nil, err
	}

	if vm.IsIterable(args[0]) {
		return skipIter(m, args, nRet)
	} else {
		return skipReader(m, args, nRet)
	}
}

func skipIter(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 2, 2); err != nil {
		return nil, err
	}

	inIter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	skip, err := getIntArg(args, 1)
	if err != nil {
		return nil, err
	}

	skipIter := &skipIterator{inIter: inIter, skip: int(skip)}

	return []vm.Value{vm.NewIterator(skipIter.Next, skipIter.Close, inIter.IterNRet())}, nil
}

func skipReader(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 2, 2); err != nil {
		return nil, err
	}

	inReader, err := getReaderArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	skip, err := getIntArg(args, 1)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(io.Discard, io.LimitReader(inReader, skip))
	if err != nil {
		return nil, err
	}

	return []vm.Value{inReader}, nil
}

type skipIterator struct {
	inIter vm.Iterator
	skip   int
}

func (i *skipIterator) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	for {
		v, err := m.IterNext(i.inIter, i.inIter.IterNRet())
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, nil
		}
		if i.skip > 0 {
			i.skip--
			continue
		}
		return v, nil
	}
}

func (i *skipIterator) Close(m *vm.VM) error {
	return m.IterClose(i.inIter)
}
