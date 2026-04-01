package lib

import (
	"fmt"

	"github.com/dcaiafa/bag3l/internal/vm"
)

// $format - used by compiler to implement format expressions.
func format(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if len(args) != 2 {
		return nil, errInvalidNumberOfArgs
	}

	v := args[0]

	var vany any
	switch v := v.(type) {
	case vm.Int:
		vany = v.Int64()
	case vm.Float:
		vany = v.Float64()
	case vm.Bool:
		vany = v.Bool()
	default:
		if v == nil {
			vany = "<nil>"
		} else {
			vany = v.String()
		}
	}

	f, err := getStringArg(args, 1)
	if err != nil {
		return nil, err
	}

	res := fmt.Sprintf(f, vany)

	return []vm.Value{vm.NewString(res)}, nil
}
