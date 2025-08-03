package lib

import (
	"fmt"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func mod(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 2 {
		return nil, errInvalidNumberOfArgs
	}

	a, ok := args[0].(vm.Operable)
	if !ok {
		return nil, fmt.Errorf("argument does not support operation")
	}

	res, err := a.EvalOp(vm.OpMod, args[1])
	if err != nil {
		return nil, err
	}

	return []vm.Value{res}, nil
}
