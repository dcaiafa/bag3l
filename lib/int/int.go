package int

import (
	"fmt"
	"strconv"

	"github.com/dcaiafa/bag3l/internal/vm"
)

//go:generate go run ../../internal/stub/stubgen int.stubgen

func parse0(m *vm.VM, v string, base int64) (int64, error) {
	return strconv.ParseInt(v, int(base), 64)
}

func into0(m *vm.VM, v vm.Value) (int64, error) {
	if intArg, ok := v.(vm.Int); ok {
		return intArg.Int64(), nil
	}

	inter, ok := v.(interface {
		ToInt() (vm.Int, error)
	})
	if !ok {
		return 0, vm.MakeNonRecoverableError(
			fmt.Errorf("%v is not convertible to Int",
				vm.TypeName(v)))
	}

	res, err := inter.ToInt()
	if err != nil {
		return 0, err
	}

	return res.Int64(), nil
}
