package bool

import "github.com/dcaiafa/bag3l/internal/vm"

//go:generate go run ../../internal/stub/stubgen bool.stubgen

func into0(m *vm.VM, v vm.Value) (bool, error) {
	return vm.CoerceToBool(v), nil
}
