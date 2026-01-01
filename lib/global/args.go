package global

import "github.com/dcaiafa/bag3l/internal/vm"

func args0(m *vm.VM) (vm.Value, error) {
	frameArgs := m.GetCallerArgs()
	argsCopy := make([]vm.Value, len(frameArgs))
	copy(argsCopy, frameArgs)
	res := vm.NewListWithSlice(argsCopy)
	return res, nil
}
