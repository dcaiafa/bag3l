package lib

import (
	"errors"

	"github.com/dcaiafa/bag3l/internal/vm"
)

var errTypeUsage = errors.New(
	`invalid usage. Expected type(any)`)

func typep(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 1 {
		return nil, errTypeUsage
	}

	res := vm.TypeName(args[0])

	return []vm.Value{vm.NewString(res)}, nil
}
