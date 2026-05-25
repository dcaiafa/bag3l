package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func in(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	if len(args) < 2 {
		return nil, errNotEnoughArgs
	}

	v := args[0]

	found := false
	for _, arg := range args[1:] {
		if v == arg {
			found = true
			break
		}
	}

	return []vm.Value{vm.NewBool(found)}, nil
}
