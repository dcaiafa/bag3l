package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func start(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) > 1 {
		return nil, errTooManyArgs
	} else if len(args) < 1 {
		return nil, errNotEnoughArgs
	}

	err := m.StartCoroutine(args[0])
	if err != nil {
		return nil, err
	}

	return nil, nil
}
