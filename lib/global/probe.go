package global

import (
	"github.com/dcaiafa/bag3l/internal/vm"
	libio "github.com/dcaiafa/bag3l/lib/io"
)

func probe0(m *vm.VM, args []vm.Value) (vm.Value, error) {
	probeValue := args[0]
	otherArgs := make([]vm.Value, len(args)-1)
	copy(otherArgs, args[1:])

	if inIter, ok := probeValue.(vm.Iterator); ok {
		probeIter := &probeIter{
			inIter:    inIter,
			otherArgs: otherArgs,
		}
		return vm.NewIterator(probeIter.Next, probeIter.Close, inIter.IterNRet()), nil
	}

	printArgs := append(otherArgs, probeValue)
	err := basePrint(libio.Stderr(m), m, printArgs)
	if err != nil {
		return nil, err
	}

	return probeValue, nil
}

type probeIter struct {
	inIter    vm.Iterator
	otherArgs []vm.Value
}

func (i *probeIter) Next(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	v, err := m.IterNext(i.inIter, i.inIter.IterNRet())
	if err != nil {
		i.Close(m)
		return nil, err
	}
	if v == nil {
		i.Close(m)
		return nil, nil
	}
	printArgs := append(i.otherArgs, v...)
	basePrint(libio.Stderr(m), m, printArgs)

	return v, nil
}

func (i *probeIter) Close(m *vm.VM) error {
	m.IterClose(i.inIter)
	return nil
}
