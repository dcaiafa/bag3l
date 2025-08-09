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

	// 'next' returns all the values from the iterator plus a boolean reporting
	// success. Example:
	//   a, b, ok = next(it)
	// In this case, nret will be 3, but we are going to request only 2 values
	// from the actual iterator. We need to protect from the case where nret is
	// zero. Example:
	//   next(it)
	// As we can't request negative values from the iterator.
	// N.B. nret is the minimum number of values expected by the VM. It is OK to
	// return more values than requested. They will be discarded by the VM.
	if nret == 0 {
		nret = 1
	}

	res, err := vm.IterNext(iter, nret-1)
	if err != nil {
		return nil, err
	}

	if res == nil {
		vm.IterClose(iter)
		res = make([]nitro.Value, nret)
		res[len(res)-1] = nitro.False
	} else {
		res = append(res, nitro.True)
	}

	return res, nil
}
