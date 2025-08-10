package lib

import (
	"fmt"

	"github.com/dcaiafa/bag3l/internal/stub"
	"github.com/dcaiafa/bag3l/internal/vm"
)

func index(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 2 {
		return nil, errInvalidNumberOfArgs
	}

	indexable, ok := args[0].(vm.Indexable)
	if !ok {
		return nil, stub.InvalidArgErr(
			args, 0, fmt.Errorf("%v is not indexable", vm.TypeName(args[0])))
	}

	key := args[1]

	res, ok, err := indexable.Index(key)
	if err != nil {
		return nil, err
	}

	if !ok {
		return []vm.Value{nil, vm.False}, nil
	}

	return []vm.Value{res, vm.True}, nil
}
