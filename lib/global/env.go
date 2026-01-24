package global

import (
	"os"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func env0(m *vm.VM, name string) (vm.Value, error) {
	if name == "" {
		allEnv := os.Environ()
		res := make([]vm.Value, len(allEnv))
		for i, e := range allEnv {
			res[i] = vm.NewString(e)
		}
		return vm.NewListWithSlice(res), nil
	}
	return vm.NewString(os.Getenv(name)), nil
}
