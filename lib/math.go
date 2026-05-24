package lib

import (
	"math"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func mathTrunc(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	if err := expectArgCount(args, 1, 1); err != nil {
		return nil, err
	}
	v, err := getFloatArg(args, 0)
	if err != nil {
		return nil, err
	}
	res := math.Trunc(v)
	return []vm.Value{vm.NewFloat(res)}, nil
}
