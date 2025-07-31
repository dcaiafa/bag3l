package lib

import "github.com/dcaiafa/bag3l/internal/vm"

func args(m *vm.VM, args []vm.Value, nRet int) ([]vm.Value, error) {
	frameArgs := m.GetCallerArgs()
	argsCopy := make([]vm.Value, len(frameArgs))
	copy(argsCopy, frameArgs)
	res := vm.NewListWithSlice(argsCopy)
	return []vm.Value{res}, nil
}
