package lib

import (
	nitro "github.com/dcaiafa/bag3l"
	"github.com/dcaiafa/bag3l/internal/stub"
)

func next(vm *nitro.VM, args []nitro.Value, nret int) ([]nitro.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}

	iter, ok := args[0].(nitro.Iterator)
	if !ok {
		return nil, stub.InvalidArg(args, 0)
	}

	res, err := vm.IterNext(iter, nret)
	if err != nil {
		return nil, err
	}

	if res == nil {
		vm.IterClose(iter)
		return make([]nitro.Value, nret), nil
	}

	return res, nil
}
