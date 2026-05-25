package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func flatten(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	if len(args) > 1 {
		return nil, errTooManyArgs
	}

	inIter, err := getIterArg(m, args, 0)
	if err != nil {
		return nil, err
	}

	flattenIter := &flattenIter{
		first: inIter,
	}

	outIter := vm.NewIterator(flattenIter.Next, flattenIter.Close, 1)

	return []vm.Value{outIter}, nil
}

type flattenIter struct {
	first  vm.Iterator
	second vm.Iterator
}

func (i *flattenIter) Next(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if i.second != nil {
		v, err := m.IterNext(i.second, 1)
		if err != nil {
			i.Close(m)
			return nil, err
		}
		if v == nil {
			i.second = nil
			return i.Next(m, args, nRet)
		}
		return v, nil
	}

	v, err := m.IterNext(i.first, 1)
	if err != nil {
		i.Close(m)
		return nil, err
	}

	if v == nil {
		i.Close(m)
		return nil, nil
	}

	i.second, err = vm.MakeIterator(m, v[0])
	if err != nil {
		return v, nil
	}

	return i.Next(m, args, nRet)
}

func (i *flattenIter) Close(m *vm.VM) error {
	m.IterClose(i.first)
	if i.second != nil {
		m.IterClose(i.second)
	}
	return nil
}
